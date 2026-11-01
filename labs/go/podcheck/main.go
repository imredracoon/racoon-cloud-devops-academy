// podcheck elenca i Pod che hanno problemi (lezione G.9 del percorso Go).
//
//	go run . -n tasks            un namespace
//	go run . -A                  tutti i namespace
//	go run . -A -restarts 3      segnala anche i container riavviati almeno 3 volte
package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	home, _ := os.UserHomeDir()
	kubeconfig := flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "file kubeconfig")
	namespace := flag.String("n", "default", "namespace da controllare")
	tutti := flag.Bool("A", false, "controlla tutti i namespace")
	soglia := flag.Uint("restarts", 5, "riavvii oltre i quali un container è sospetto")
	flag.Parse()

	if *tutti {
		*namespace = "" // stringa vuota = tutti i namespace
	}

	clientset, err := connetti(*kubeconfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "errore:", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	pods, err := clientset.CoreV1().Pods(*namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "errore:", err)
		os.Exit(1)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "NAMESPACE\tPOD\tPROBLEMA")
	trovati := 0
	limite := int32(min(*soglia, math.MaxInt32)) //nolint:gosec // il valore è limitato a MaxInt32
	for _, p := range pods.Items {
		for _, problema := range problemi(p, limite) {
			fmt.Fprintf(tw, "%s\t%s\t%s\n", p.Namespace, p.Name, problema)
			trovati++
		}
	}
	_ = tw.Flush()
	fmt.Printf("\n%d Pod controllati, %d problemi trovati\n", len(pods.Items), trovati)
	if trovati > 0 {
		os.Exit(2)
	}
}

// connetti usa la configurazione interna se gira dentro un Pod, altrimenti il kubeconfig.
func connetti(kubeconfig string) (*kubernetes.Clientset, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("configurazione del cluster: %w", err)
		}
	}
	return kubernetes.NewForConfig(cfg)
}

// problemi restituisce, in parole semplici, cosa non va in un Pod.
func problemi(p corev1.Pod, soglia int32) []string {
	var out []string
	if p.Status.Phase == corev1.PodPending || p.Status.Phase == corev1.PodFailed {
		out = append(out, "fase "+string(p.Status.Phase))
	}
	for _, c := range p.Status.ContainerStatuses {
		if w := c.State.Waiting; w != nil && w.Reason != "ContainerCreating" {
			out = append(out, fmt.Sprintf("container %s in attesa: %s", c.Name, w.Reason))
		}
		if c.RestartCount >= soglia {
			out = append(out, fmt.Sprintf("container %s riavviato %d volte", c.Name, c.RestartCount))
		}
		if !c.Ready && p.Status.Phase == corev1.PodRunning {
			out = append(out, fmt.Sprintf("container %s non pronto", c.Name))
		}
	}
	return out
}
