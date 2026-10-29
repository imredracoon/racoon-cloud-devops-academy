# 8.12 Cluster con kubeadm

!!! obiettivi "Obiettivi"
    - Preparare i nodi Linux per Kubernetes (Debian/Ubuntu e RHEL).
    - Creare un cluster con `kubeadm init` e `kubeadm join` e installare un plugin di rete.
    - Aggiornare il cluster di una versione minore, nodo per nodo.
    - Salvare e ripristinare etcd.
    - Controllare certificati e componenti statici del control plane.

## Perché imparare kubeadm

Nel lavoro userai spesso Kubernetes **gestito** (EKS, GKE), dove il control plane è del provider. Ma la CKA chiede di
saper **installare, aggiornare e riparare** un cluster, e capirne le viscere ti rende molto più efficace quando
qualcosa si rompe. **kubeadm** è lo strumento ufficiale per creare cluster conformi: genera certificati, scrive i
manifest dei componenti del control plane come *static pod*, configura il kubelet.

!!! info "Cosa ti serve"
    Due o tre macchine virtuali Linux con almeno 2 CPU e 2 GB di RAM ciascuna (4 GB per il control plane è più
    comodo), in rete tra loro: per esempio VM locali create con KVM/libvirt come quelle del corso, oppure istanze EC2 o Compute Engine
    (ricorda di spegnerle a fine lavoro). kind non basta per questa lezione, ma la parte su etcd e certificati la
    puoi provare anche lì.

## 1. Preparare ogni nodo

Esegui questi passi su **tutti** i nodi, come root o con `sudo`.

### Requisiti di sistema

```bash
# Moduli del kernel e parametri di rete per il traffico dei Pod
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF
sudo modprobe overlay
sudo modprobe br_netfilter

cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF
sudo sysctl --system

# Swap: il kubelet di default rifiuta di partire se lo swap è attivo
sudo swapoff -a
sudo sed -i '/\sswap\s/ s/^/#/' /etc/fstab
```

- Hostname e indirizzo MAC devono essere **unici** per ogni nodo (attenzione alle VM clonate).
- Il sistema deve usare **cgroup v2** (predefinito su Ubuntu 22.04+, Debian 12+, RHEL 9+): dalla versione 1.35 il
  kubelet non parte su cgroup v1. Verifica con `stat -fc %T /sys/fs/cgroup` (deve rispondere `cgroup2fs`).
- Le porte del control plane (6443 API server, 2379-2380 etcd, 10250 kubelet, 10257 e 10259) e dei worker (10250,
  30000-32767 NodePort) devono essere raggiungibili tra i nodi.

### Container runtime e pacchetti di Kubernetes

I pacchetti ufficiali sono su `pkgs.k8s.io`, con un repository **per ogni versione minore**: per passare alla 1.38
dovrai cambiare repository.

=== "Debian/Ubuntu"

    ```bash
    # containerd dal repository della distribuzione
    sudo apt-get update
    sudo apt-get install -y containerd apt-transport-https ca-certificates curl gpg

    # Repository di Kubernetes 1.37
    sudo mkdir -p -m 755 /etc/apt/keyrings
    curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.37/deb/Release.key \
      | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
    echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.37/deb/ /' \
      | sudo tee /etc/apt/sources.list.d/kubernetes.list

    sudo apt-get update
    sudo apt-get install -y kubelet kubeadm kubectl
    sudo apt-mark hold kubelet kubeadm kubectl     # niente aggiornamenti accidentali
    sudo systemctl enable --now kubelet
    ```

=== "RHEL (Rocky/AlmaLinux)"

    ```bash
    # containerd dal repository di Docker (pacchetto containerd.io)
    sudo dnf install -y dnf-plugins-core
    sudo dnf config-manager --add-repo https://download.docker.com/linux/rhel/docker-ce.repo
    sudo dnf install -y containerd.io

    # SELinux in modalità permissiva, come indica la guida ufficiale di kubeadm
    sudo setenforce 0
    sudo sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

    # Repository di Kubernetes 1.37 (exclude: i pacchetti si aggiornano solo su richiesta esplicita)
    cat <<EOF | sudo tee /etc/yum.repos.d/kubernetes.repo
    [kubernetes]
    name=Kubernetes
    baseurl=https://pkgs.k8s.io/core:/stable:/v1.37/rpm/
    enabled=1
    gpgcheck=1
    gpgkey=https://pkgs.k8s.io/core:/stable:/v1.37/rpm/repodata/repomd.xml.key
    exclude=kubelet kubeadm kubectl cri-tools kubernetes-cni
    EOF

    sudo dnf install -y kubelet kubeadm kubectl --setopt=disable_excludes=kubernetes
    sudo systemctl enable --now kubelet

    # firewalld: apri le porte (sul control plane; sui worker 10250 e 30000-32767/tcp)
    sudo firewall-cmd --permanent --add-port={6443,2379-2380,10250,10257,10259}/tcp
    sudo firewall-cmd --reload
    ```

    Sui sistemi RHEL il CNI scelto può richiedere porte aggiuntive (per esempio BGP 179/tcp e VXLAN 4789/udp per
    Calico): controlla la sua documentazione.

### Configurare containerd

Kubelet e runtime devono usare lo stesso **driver dei cgroup**: con systemd come init, `systemd`. kubeadm imposta già
il kubelet su `systemd`; resta containerd:

```bash
sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml > /dev/null
sudo sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml
sudo systemctl restart containerd
sudo systemctl enable containerd
```

Il kubelet ora riparte in continuazione: è normale, aspetta le istruzioni di kubeadm.

## 2. Inizializzare il control plane

Sul primo nodo:

```bash
sudo kubeadm init --pod-network-cidr=192.168.0.0/16
```

kubeadm controlla i requisiti (*preflight*), genera la CA e i certificati in `/etc/kubernetes/pki`, scrive i
manifest di etcd, API server, controller manager e scheduler in `/etc/kubernetes/manifests`, avvia tutto e stampa il
comando per unire i worker. Poi:

```bash
mkdir -p $HOME/.kube
sudo cp /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
kubectl get nodes          # NotReady: manca il plugin di rete
```

### Il plugin di rete

Senza CNI i Pod non ricevono indirizzi e CoreDNS resta `Pending`. Con **Calico** (che applica anche le Network
Policy), seguendo la sua guida ufficiale:

```bash
kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.33.0/manifests/v3_projectcalico_org.yaml
kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.33.0/manifests/tigera-operator.yaml
kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.33.0/manifests/custom-resources-bpf.yaml
kubectl get pods -A -w
```

Il CIDR dei Pod nel file `custom-resources` deve coincidere con `--pod-network-cidr` (qui il predefinito
`192.168.0.0/16`, che non deve sovrapporsi alla rete delle VM: se le tue VM sono in `192.168.x.x`, usa per esempio
`10.244.0.0/16` in entrambi i punti).

## 3. Unire i worker

Su ogni worker, il comando stampato da `kubeadm init`:

```bash
sudo kubeadm join 10.0.0.10:6443 --token abcdef.0123456789abcdef \
  --discovery-token-ca-cert-hash sha256:...
```

Il token scade dopo 24 ore. Per generarne uno nuovo con il comando completo:

```bash
kubeadm token create --print-join-command
```

## 4. Aggiornare il cluster

Si sale **una versione minore alla volta** (1.36 → 1.37, mai 1.35 → 1.37), prima il control plane, poi i worker
uno per uno. Esempio da 1.36 a 1.37 su Debian/Ubuntu (su RHEL stessi passi con `dnf` e
`--disableexcludes=kubernetes`):

```bash
# --- Sul control plane ---
# 1. repository della nuova minore
sudo sed -i 's/v1.36/v1.37/' /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
apt-cache madison kubeadm | head -3                  # versioni disponibili

# 2. kubeadm
sudo apt-mark unhold kubeadm
sudo apt-get install -y kubeadm='1.37.2-*'
sudo apt-mark hold kubeadm
sudo kubeadm upgrade plan
sudo kubeadm upgrade apply v1.37.2                   # aggiorna i componenti del control plane

# 3. kubelet e kubectl, con il nodo svuotato
kubectl drain cp1 --ignore-daemonsets
sudo apt-mark unhold kubelet kubectl
sudo apt-get install -y kubelet='1.37.2-*' kubectl='1.37.2-*'
sudo apt-mark hold kubelet kubectl
sudo systemctl daemon-reload && sudo systemctl restart kubelet
kubectl uncordon cp1

# --- Su ogni worker, uno alla volta ---
# repository e kubeadm come sopra, poi:
sudo kubeadm upgrade node
kubectl drain worker1 --ignore-daemonsets            # da una macchina con kubectl
# kubelet e kubectl come sopra, daemon-reload, restart
kubectl uncordon worker1
```

Regole di compatibilità (*version skew*): il kubelet può essere fino a **tre** versioni minori più vecchio
dell'API server, mai più nuovo; kubectl al massimo una minore di distanza. Prima di aggiornare leggi sempre le note
di rilascio: le API deprecate vengono rimosse.

## 5. Backup e ripristino di etcd

etcd contiene **tutto** lo stato. Un backup regolare è obbligatorio, e all'esame è quasi certo.

```bash
# Backup (sul control plane; i percorsi dei certificati sono nel manifest di etcd)
sudo ETCDCTL_API=3 etcdctl --endpoints=https://127.0.0.1:2379 \
  --cacert=/etc/kubernetes/pki/etcd/ca.crt \
  --cert=/etc/kubernetes/pki/etcd/server.crt \
  --key=/etc/kubernetes/pki/etcd/server.key \
  snapshot save /var/backups/etcd-$(date +%F).db

sudo etcdutl snapshot status /var/backups/etcd-2026-10-28.db -w table
```

Ripristino:

```bash
# 1. ripristina lo snapshot in una NUOVA cartella
sudo etcdutl snapshot restore /var/backups/etcd-2026-10-28.db --data-dir /var/lib/etcd-ripristino

# 2. fai puntare etcd alla nuova cartella: nel manifest statico cambia il volume hostPath
sudo vi /etc/kubernetes/manifests/etcd.yaml
#   volumes:
#   - hostPath:
#       path: /var/lib/etcd-ripristino      # era /var/lib/etcd
#     name: etcd-data

# 3. il kubelet ricrea il Pod di etcd da solo; attendi che l'API server torni
kubectl get pods -n kube-system
```

Dalla versione 3.5 di etcd i comandi offline (`snapshot status`, `snapshot restore`) sono in `etcdutl`; `etcdctl`
resta il client per i comandi verso il server, come `snapshot save`. Se `etcdctl` non è installato sul nodo, lo trovi
dentro il Pod di etcd.

## Static pod e certificati

I componenti del control plane sono **static pod**: il kubelet li legge dalla cartella indicata da `staticPodPath`
nella sua configurazione (`/var/lib/kubelet/config.yaml`), di solito `/etc/kubernetes/manifests`. Modificare uno di
quei file riavvia il componente; spostarlo fuori dalla cartella lo ferma. Gli oggetti che vedi con
`kubectl get pods -n kube-system` sono solo **copie speculari** (*mirror pod*): eliminarli non serve.

```bash
sudo kubeadm certs check-expiration     # i certificati scadono dopo un anno
sudo kubeadm certs renew all            # rinnovo (anche automatico a ogni kubeadm upgrade)
```

## Alta disponibilità, in breve

Un control plane in alta disponibilità ha tre (o cinque) nodi dietro un **load balancer** per l'API server:
`kubeadm init --control-plane-endpoint=lb.example.com:6443 --upload-certs`, poi `kubeadm join ... --control-plane`
sugli altri. etcd può stare sugli stessi nodi (*stacked*, la scelta comune) o su macchine separate (*external*).

## Laboratorio

!!! laboratorio "Laboratorio · Backup di etcd in kind"
    Il nodo di control plane di kind è creato proprio con kubeadm: puoi esplorarlo.

    ```bash
    docker exec -it academy-control-plane bash
    ls /etc/kubernetes/manifests /etc/kubernetes/pki
    kubeadm certs check-expiration
    grep -E 'staticPodPath|cgroupDriver' /var/lib/kubelet/config.yaml
    exit
    ```

    Poi esegui un backup di etcd usando `etcdctl` dentro il Pod di etcd:

    ```bash
    kubectl -n kube-system exec etcd-academy-control-plane -- etcdctl \
      --endpoints=https://127.0.0.1:2379 \
      --cacert=/etc/kubernetes/pki/etcd/ca.crt \
      --cert=/etc/kubernetes/pki/etcd/server.crt \
      --key=/etc/kubernetes/pki/etcd/server.key \
      snapshot save /var/lib/etcd/backup.db
    docker exec academy-control-plane ls -lh /var/lib/etcd/backup.db
    ```

    1. Crea un namespace `prima-del-disastro`, poi esegui il backup.
    2. Elimina il namespace.
    3. Ripristina il backup e verifica che il namespace sia tornato.

??? success "Soluzione"
    ```bash
    kubectl create namespace prima-del-disastro
    # ...backup con il comando qui sopra...
    kubectl delete namespace prima-del-disastro
    ```
    L'immagine di etcd contiene anche `etcdutl`, e il container vede la cartella `/var/lib/etcd` del nodo:
    ripristina lì dentro, in una nuova sottocartella, poi fai puntare il volume del manifest statico alla nuova
    cartella (solo la riga `path:` del volume `hostPath`, non `mountPath`):
    ```bash
    kubectl -n kube-system exec etcd-academy-control-plane -- \
      etcdutl snapshot restore /var/lib/etcd/backup.db --data-dir /var/lib/etcd/ripristino
    docker exec academy-control-plane sed -i \
      's#^\( *\)path: /var/lib/etcd$#\1path: /var/lib/etcd/ripristino#' /etc/kubernetes/manifests/etcd.yaml
    ```
    Dopo circa un minuto il kubelet ricrea etcd e l'API server si riconnette (nel frattempo `kubectl` può dare
    errori di connessione: è normale). Poi `kubectl get ns prima-del-disastro` mostra il namespace. Il ripristino
    riporta **tutto** il cluster a quel momento: in produzione si coordina con cura, perché tutto ciò che è stato
    creato dopo il backup sparisce.

!!! esercizio "Esercizio · L'ordine dell'aggiornamento"
    Metti in ordine i passi per aggiornare un worker da 1.36 a 1.37: (a) `kubectl uncordon`; (b) installare il
    nuovo kubeadm; (c) `kubeadm upgrade node`; (d) `kubectl drain`; (e) installare kubelet e kubectl nuovi e
    riavviare il kubelet; (f) cambiare il repository dei pacchetti.

??? success "Soluzione"
    **f, b, c, d, e, a**. Il control plane deve essere già stato aggiornato. `kubeadm upgrade node` sul worker
    aggiorna la configurazione locale del kubelet; il drain serve prima di riavviare il kubelet con la nuova
    versione.

## Nei tre livelli

=== "Livello 1 · Un PC"

    Il laboratorio con kind. Con 16 GB di RAM sulla workstation puoi fare anche la costruzione del livello 2 con le
    VM in locale: i passi sono identici, senza `srv1`.

=== "Livello 2 · Workstation + server"

    **Il cluster sulle VM di `srv1`.**

    1. Su `srv1` crea le VM (`./academy-vm crea k8s-cp`, poi `k8s-w1` e `k8s-w2`) e aggiorna `~/.ssh/config` della
       workstation con `./academy-vm ssh-config srv1`.
    2. Su tutte e tre esegui il [passo 1](#1-preparare-ogni-nodo) (scheda Debian/Ubuntu). Ancora meglio, scrivilo come ruolo
       Ansible `kubernetes_nodo`: lo riuserai.
    3. Su `k8s-cp`, poiché le VM sono in `192.168.56.0/24`, usa un CIDR dei Pod diverso, e aggiungi `127.0.0.1` ai
       nomi del certificato dell'API server, per usare `kubectl` dalla workstation con un tunnel:

        ```bash
        sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --apiserver-cert-extra-sans=127.0.0.1
        ```

        Nel file `custom-resources` di Calico imposta lo stesso `10.244.0.0/16`.

    4. Unisci `k8s-w1` e `k8s-w2`.
    5. **kubectl dalla workstation.** Copia `/etc/kubernetes/admin.conf` di `k8s-cp` in `~/.kube/homelab.conf` sulla
       workstation, cambia la riga `server:` in `https://127.0.0.1:6443` e apri il tunnel:

        ```bash
        ssh -N -L 6443:192.168.56.31:6443 srv1 &
        export KUBECONFIG=~/.kube/homelab.conf
        kubectl get nodes -o wide
        ```

    **Collaudo**

    | Prova | Risultato atteso |
    |---|---|
    | `kubectl get nodes -o wide` | tre nodi `Ready`, con indirizzi `192.168.56.31-33` |
    | `kubectl get pods -A` | tutti `Running`, compresi Calico e CoreDNS |
    | `kubectl create deployment prova --image=nginx:1.29 --replicas=4`, poi `kubectl get pods -o wide` | Pod distribuiti sui due worker |
    | `kubectl run t --rm -it --image=busybox:1.37 -- nslookup kubernetes.default` | risponde il DNS del cluster |
    | `ssh srv1 virsh shutdown k8s-w2` | dopo circa un minuto il nodo è `NotReady`; dopo 5 minuti i Pod vengono ricreati sull'altro worker |
    | il backup di etcd del laboratorio, su `k8s-cp` con `etcdctl` installato o tramite il Pod di etcd | snapshot salvato e ripristino riuscito |
    | l'aggiornamento di una versione minore, se disponibile ([passo 4](#4-aggiornare-il-cluster)) | `kubectl get nodes` mostra la nuova versione |

    **Servizi `LoadBalancer` con MetalLB.** Su un cluster fuori dal cloud nessuno assegna gli indirizzi ai servizi
    `LoadBalancer`: lo fa **MetalLB**, che annuncia sulla rete locale un indirizzo preso da un intervallo che gli
    assegni (fuori dal DHCP: le VM usano `100-199`).

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.15.2/config/manifests/metallb-native.yaml
    kubectl -n metallb-system wait --for=condition=Available deployment --all --timeout=5m
    ```

    ```yaml title="metallb.yaml"
    apiVersion: metallb.io/v1beta1
    kind: IPAddressPool
    metadata:
      name: laboratorio
      namespace: metallb-system
    spec:
      addresses: ["192.168.56.200-192.168.56.220"]
    ---
    apiVersion: metallb.io/v1beta1
    kind: L2Advertisement
    metadata:
      name: laboratorio
      namespace: metallb-system
    spec:
      ipAddressPools: [laboratorio]
    ```

    **Storage NFS.** Su `srv1` esporta `/srv/condivisa/k8s` a `192.168.56.0/24` ([lezione 1.6](../fase-1/06-storage.md#filesystem-di-rete-nfs-e-autofs)), installa
    `nfs-common` sulle VM e il *provisioner* che crea un volume per ogni PVC:

    ```bash
    helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/
    helm install nfs nfs-subdir-external-provisioner/nfs-subdir-external-provisioner \
      -n nfs --create-namespace \
      --set nfs.server=192.168.56.1 --set nfs.path=/srv/condivisa/k8s \
      --set storageClass.name=nfs --set storageClass.defaultClass=true
    ```

    | Prova | Risultato atteso |
    |---|---|
    | `kubectl expose deployment prova --port=80 --type=LoadBalancer`, poi `kubectl get svc prova` | `EXTERNAL-IP` `192.168.56.200` |
    | `ssh srv1 curl -sI http://192.168.56.200` | `200` da Nginx |
    | un PVC con `storageClassName: nfs` usato da un Pod | `Bound`; su `srv1` compare una cartella in `/srv/condivisa/k8s` |
    | il Pod spostato sull'altro worker (`kubectl drain` del primo) | i file scritti nel volume ci sono ancora: lo storage è condiviso, non locale al nodo |

=== "Livello 3 · Mini data center"

    **Il cluster sui nodi fisici.** `nodo2` diventa il control plane, `nodo3` il worker.

    !!! warning "Prima di iniziare: liberare i nodi"
        Su `nodo2` e `nodo3` girano gli stack delle fasi 3-7. Fai il **backup** dei database (lo ripristinerai nel
        cluster con il progetto), poi la scelta più pulita è **reinstallare** i due nodi e applicare il playbook
        `sito.yml` della [Fase 6](../fase-6/progetto.md#nei-tre-livelli): è anche la prova che il tuo codice ricostruisce davvero
        il data center. In alternativa, ferma e disinstalla Docker e Podman e i loro servizi. Il monitoraggio della
        Fase 7 si sposta nel cluster con il progetto.

    1. Il [passo 1](#1-preparare-ogni-nodo) su `nodo2` (scheda Debian/Ubuntu) e su `nodo3` (scheda RHEL).
    2. **Firewall.** Il perimetro del data center è `nodo1`: sui nodi del cluster fidati della rete del laboratorio e
       di quella dei Pod. Su `nodo2`, nella catena `input` della tabella `homelab`,
       `ip saddr { 10.10.10.0/24, 10.244.0.0/16 } accept`; su `nodo3`
       `sudo firewall-cmd --permanent --zone=trusted --add-source=10.10.10.0/24 --add-source=10.244.0.0/16` e
       `--reload`.
    3. Su `nodo2`: `sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --apiserver-cert-extra-sans=127.0.0.1`,
       Calico con lo stesso CIDR, poi `kubeadm join` su `nodo3`.
    4. Con due soli nodi, permetti al control plane di ospitare i Pod:
       `kubectl taint nodes nodo2 node-role.kubernetes.io/control-plane:NoSchedule-`.
    5. kubectl dalla workstation come nel livello 2, con il tunnel `ssh -N -L 6443:10.10.10.12:6443 nodo1 &`.
    6. **MetalLB** come nel livello 2, con l'intervallo `10.10.10.200-10.10.10.220` (fuori dal DHCP di `nodo1`).
    7. **Storage NFS** come nel livello 2, con `nfs.server=nodo1.lab`, la cartella `/srv/condivisa/k8s` esportata a
       `10.10.10.0/24`, e i pacchetti client su `nodo2` e `nodo3` (dalla [lezione 1.6](../fase-1/06-storage.md#nei-tre-livelli) ci sono già).
    8. **Pubblicazione verso casa.** Il Gateway della [lezione 8.9](09-ingress.md#nei-tre-livelli) riceverà da MetalLB un indirizzo,
       per esempio `10.10.10.200`: nel firewall di `nodo1` le porte 80 e 443 vanno inoltrate lì invece che a `nodo2`.

    **Collaudo**

    | Prova | Risultato atteso |
    |---|---|
    | `kubectl get nodes -o wide` | `nodo2` (control plane) e `nodo3` `Ready`, con `OS-IMAGE` Ubuntu e AlmaLinux |
    | `kubectl create deployment prova --image=nginx:1.29 --replicas=4`, poi `-o wide` | Pod su tutti e due i nodi |
    | `kubectl expose deployment prova --port=80 --type=LoadBalancer` | `EXTERNAL-IP` `10.10.10.200`; `ssh nodo1 curl -sI http://10.10.10.200` risponde |
    | **stacca il cavo di rete di `nodo3`** | `nodo3` `NotReady` dopo circa un minuto; i Pod tornano su `nodo2` dopo circa 5 minuti; il servizio risponde ancora |
    | ricollega il cavo | `nodo3` torna `Ready` |
    | PVC `nfs` e `kubectl drain nodo3` | il Pod riparte su `nodo2` con gli stessi dati |
    | backup di etcd su `nodo2`, copiato su NFS e fuori sede ([lezione 4.6](../fase-4/06-storage.md#nei-tre-livelli)) | lo snapshot è nel bucket |
    | riavvio di tutto il data center (prima `nodo1`, poi gli altri) | il cluster torna `Ready` senza interventi |

## Verifica rapida

<script type="application/json" class="quiz-data">
{"questions":[
 {"q":"Dove si trovano i manifest dei componenti del control plane in un cluster kubeadm?","options":["/etc/kubernetes/manifests","/var/lib/etcd","/etc/systemd/system","~/.kube"],"answer":[0],"explanation":"Sono static pod letti dal kubelet."},
 {"q":"Quale comando crea un nuovo token con il comando completo per unire un worker?","options":["kubeadm join --new","kubeadm token create --print-join-command","kubectl create token","kubeadm init phase join"],"answer":[1],"explanation":"I token durano 24 ore."},
 {"q":"Da 1.35 vuoi arrivare a 1.37. Come procedi?","options":["Direttamente a 1.37","Prima a 1.36, poi a 1.37","Reinstalli il cluster","Aggiorni solo i worker"],"answer":[1],"explanation":"Una minore alla volta."},
 {"q":"Quale comando salva uno snapshot di etcd?","options":["etcdctl snapshot save","etcdutl backup","kubectl backup etcd","kubeadm etcd save"],"answer":[0],"explanation":"Il ripristino offline si fa con etcdutl snapshot restore."},
 {"q":"Quali passi sono necessari su ogni nodo prima di kubeadm init o join?","options":["Container runtime installato e configurato con driver systemd","Moduli overlay e br_netfilter e ip_forward attivo","Swap disattivato (o kubelet configurato per tollerarlo)","Docker Engine installato","Un Ingress controller"],"answer":[0,1,2],"explanation":"Docker Engine non serve: basta containerd."},
 {"q":"Dopo kubeadm init i nodi sono NotReady e CoreDNS è Pending. Cosa manca?","options":["Un plugin di rete (CNI)","metrics-server","Un Ingress","Un'altra replica di etcd"],"answer":[0],"explanation":"Per esempio Calico o Cilium."},
 {"q":"Quale comando mostra la scadenza dei certificati del cluster?","options":["kubeadm certs check-expiration","openssl list","kubectl get certs","kubeadm token list"],"answer":[0],"explanation":"I certificati kubeadm durano un anno e si rinnovano a ogni upgrade."}
]}
</script>

## Fonti

- [Kubernetes: Installing kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/)
- [Kubernetes: Container Runtimes](https://kubernetes.io/docs/setup/production-environment/container-runtimes/)
- [Kubernetes: Creating a cluster with kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)
- [Kubernetes: Upgrading kubeadm clusters](https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-upgrade/) e [Version Skew Policy](https://kubernetes.io/releases/version-skew-policy/)
- [Kubernetes: Operating etcd clusters for Kubernetes](https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/)
- [Kubernetes: Certificate Management with kubeadm](https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-certs/)
- [Kubernetes: Ports and Protocols](https://kubernetes.io/docs/reference/networking/ports-and-protocols/)
- [Calico: Quickstart for Calico on Kubernetes](https://docs.tigera.io/calico/latest/getting-started/kubernetes/k8s-single-node)
- [Docker: Install Docker Engine on RHEL (repository containerd.io)](https://docs.docker.com/engine/install/rhel/)
