/*
 * Racoon Cloud & DevOps Academy · simulatore d'esame
 *
 * <div class="exam-sim" data-bank="aws-clf"></div>
 * carica docs/assets/banks/aws-clf.json e offre:
 *  - banche a risposta multipla: esame completo a tempo, allenamento per dominio, ripasso degli errori
 *  - banche pratiche (CKA, LFCS): esame pratico a tempo con autovalutazione, allenamento su singolo task
 * Lo storico resta nel browser (localStorage).
 */
(function () {
  "use strict";

  var HISTORY = "academy.exams.v1";
  var scriptSrc = document.currentScript && document.currentScript.src;
  var ASSETS = scriptSrc ? new URL("../", scriptSrc).href : null; // .../assets/

  function md(s) { return window.AcademyMD ? window.AcademyMD(s) : String(s); }
  function esc(s) {
    return String(s).replace(/[&<>"']/g, function (c) {
      return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c];
    });
  }
  function load(k, f) { try { var r = localStorage.getItem(k); return r ? JSON.parse(r) : f; } catch (e) { return f; } }
  function save(k, v) { try { localStorage.setItem(k, JSON.stringify(v)); } catch (e) { /* ignorato */ } }
  function shuffle(a) {
    a = a.slice();
    for (var i = a.length - 1; i > 0; i--) {
      var j = Math.floor(Math.random() * (i + 1)); var t = a[i]; a[i] = a[j]; a[j] = t;
    }
    return a;
  }
  function fmtTime(sec) {
    var h = Math.floor(sec / 3600), m = Math.floor(sec % 3600 / 60), s = sec % 60;
    return (h ? h + ":" : "") + (m < 10 && h ? "0" : "") + m + ":" + (s < 10 ? "0" : "") + s;
  }
  var WORDS = ["zero", "one", "two", "three", "four", "five"];
  function multiHint(bank, n) {
    return bank.lang === "en" ? "(Choose " + (WORDS[n] || n) + ".)" : "(scegli " + n + ")";
  }
  function refs(list) {
    if (!list || !list.length) return "";
    return '<p class="exam-refs">Fonti: ' + list.map(function (u) {
      var label = u.replace(/^https?:\/\//, "").replace(/\/$/, "");
      if (label.length > 70) label = label.slice(0, 67) + "…";
      return '<a href="' + esc(u) + '" target="_blank" rel="noopener">' + esc(label) + "</a>";
    }).join(" · ") + "</p>";
  }

  // Estrazione proporzionale ai pesi dei domini (metodo dei resti maggiori).
  function weightedDraw(bank, items, n) {
    var byDomain = {};
    items.forEach(function (q) { (byDomain[q.d] = byDomain[q.d] || []).push(q); });
    var totalW = bank.domains.reduce(function (s, d) { return s + (byDomain[d.id] ? d.weight : 0); }, 0);
    var alloc = bank.domains.map(function (d) {
      var avail = (byDomain[d.id] || []).length;
      var exact = avail ? n * d.weight / totalW : 0;
      return { id: d.id, n: Math.min(avail, Math.floor(exact)), rest: exact - Math.floor(exact), avail: avail };
    });
    var assigned = alloc.reduce(function (s, a) { return s + a.n; }, 0);
    alloc.slice().sort(function (a, b) { return b.rest - a.rest; }).forEach(function (a) {
      if (assigned < n && a.n < a.avail) { a.n++; assigned++; }
    });
    // se una banca è piccola, riempie con quello che resta
    while (assigned < n) {
      var grew = false;
      alloc.forEach(function (a) { if (assigned < n && a.n < a.avail) { a.n++; assigned++; grew = true; } });
      if (!grew) break;
    }
    var out = [];
    alloc.forEach(function (a) { out = out.concat(shuffle(byDomain[a.id] || []).slice(0, a.n)); });
    return shuffle(out);
  }

  function domainName(bank, id) {
    for (var i = 0; i < bank.domains.length; i++) if (bank.domains[i].id === id) return bank.domains[i].name;
    return id;
  }

  function recordHistory(bank, mode, pct, passed) {
    var h = load(HISTORY, []);
    h.push({ date: new Date().toISOString(), exam: bank.title, mode: mode, pct: pct, passed: passed });
    save(HISTORY, h.slice(-200));
  }

  // =====================================================================
  // Banche a risposta multipla
  // =====================================================================
  function McqSim(root, bank) {
    this.root = root; this.bank = bank;
    this.wrongKey = "academy.wrong." + bank.id;
    this.home();
  }

  McqSim.prototype.home = function () {
    var b = this.bank, f = b.format, self = this;
    clearInterval(this.timer);
    var wrong = load(this.wrongKey, []);
    var domOpts = b.domains.map(function (d) {
      var c = b.questions.filter(function (q) { return q.d === d.id; }).length;
      return '<label><input type="checkbox" name="dom" value="' + esc(d.id) + '" checked> ' +
        esc(d.name) + " <small>(" + d.weight + "% · " + c + " domande in banca)</small></label>";
    }).join("");
    this.root.innerHTML =
      '<div class="exam-card">' +
      "<h3>" + esc(b.title) + "</h3>" +
      '<p class="exam-meta">Formato reale: <strong>' + f.questions + " domande</strong> in <strong>" +
      f.minutes + " minuti</strong> · soglia indicativa <strong>" + f.passPct + "%</strong> · " +
      b.questions.length + " domande originali in banca.</p>" +
      (f.note ? '<p class="exam-note">' + md(f.note) + "</p>" : "") +
      '<div class="exam-modes">' +
      '<div class="exam-mode"><h4>Esame completo</h4><p>Stesse condizioni dell\'esame reale: numero di domande, tempo, distribuzione per dominio. Le risposte si vedono solo alla fine.</p>' +
      '<button class="md-button md-button--primary" data-act="full">Inizia l\'esame</button></div>' +
      '<div class="exam-mode"><h4>Allenamento</h4><p>Scegli domini e numero di domande. Nessun timer, correzione e spiegazione subito dopo ogni risposta.</p>' +
      '<div class="exam-domains">' + domOpts + "</div>" +
      '<label>Domande: <input type="number" min="5" max="' + b.questions.length + '" value="20" class="exam-count"></label> ' +
      '<button class="md-button" data-act="practice">Inizia l\'allenamento</button></div>' +
      '<div class="exam-mode"><h4>Ripasso degli errori</h4><p>Ripropone le domande che hai sbagliato nelle sessioni precedenti (' + wrong.length + " al momento).</p>" +
      '<button class="md-button" data-act="wrong"' + (wrong.length ? "" : " disabled") + ">Ripassa</button></div>" +
      "</div></div>";
    this.root.querySelector('[data-act="full"]').onclick = function () {
      self.start(weightedDraw(b, b.questions, Math.min(f.questions, b.questions.length)), "esame", f.minutes);
    };
    this.root.querySelector('[data-act="practice"]').onclick = function () {
      var doms = Array.prototype.map.call(self.root.querySelectorAll('input[name="dom"]:checked'), function (i) { return i.value; });
      if (!doms.length) { alert("Scegli almeno un dominio."); return; }
      var n = Math.max(1, parseInt(self.root.querySelector(".exam-count").value, 10) || 20);
      var pool = b.questions.filter(function (q) { return doms.indexOf(q.d) !== -1; });
      self.start(weightedDraw(b, pool, Math.min(n, pool.length)), "allenamento", 0);
    };
    this.root.querySelector('[data-act="wrong"]').onclick = function () {
      var ids = load(self.wrongKey, []);
      var pool = b.questions.filter(function (q) { return ids.indexOf(q.id) !== -1; });
      self.start(shuffle(pool), "ripasso", 0);
    };
  };

  McqSim.prototype.start = function (qs, mode, minutes) {
    if (!qs.length) { alert("Nessuna domanda disponibile."); return; }
    var self = this;
    this.mode = mode;
    this.instant = mode !== "esame";
    this.items = qs.map(function (q) {
      var order = shuffle(q.o.map(function (_, i) { return i; }));
      return { q: q, order: order, chosen: [], flag: false, checked: false };
    });
    this.idx = 0;
    this.deadline = minutes ? Date.now() + minutes * 60000 : 0;
    this.root.innerHTML =
      '<div class="exam-run">' +
      '<div class="exam-bar"><span class="exam-pos"></span><span class="exam-timer"></span>' +
      '<button class="md-button exam-end">' + (this.instant ? "Termina" : "Consegna l'esame") + "</button></div>" +
      '<div class="exam-question"></div>' +
      '<div class="exam-nav"><button class="md-button exam-prev">← Precedente</button>' +
      '<label class="exam-flag"><input type="checkbox"> Segna per revisione</label>' +
      '<button class="md-button md-button--primary exam-next">Successiva →</button></div>' +
      '<div class="exam-palette"></div></div>';
    this.root.querySelector(".exam-prev").onclick = function () { self.go(self.idx - 1); };
    this.root.querySelector(".exam-next").onclick = function () { self.go(self.idx + 1); };
    this.root.querySelector(".exam-end").onclick = function () {
      var missing = self.items.filter(function (it) { return !it.chosen.length; }).length;
      var flagged = self.items.filter(function (it) { return it.flag; }).length;
      var msg = "Vuoi terminare?" + (missing ? "\nDomande senza risposta: " + missing + "." : "") +
        (flagged ? "\nDomande segnate per revisione: " + flagged + "." : "");
      if (confirm(msg)) self.finish();
    };
    this.root.querySelector(".exam-flag input").onchange = function (e) {
      self.items[self.idx].flag = e.target.checked; self.palette();
    };
    if (this.deadline) {
      this.tick();
      this.timer = setInterval(function () { self.tick(); }, 1000);
    } else {
      this.root.querySelector(".exam-timer").textContent = "Allenamento senza tempo";
    }
    this.go(0);
  };

  McqSim.prototype.tick = function () {
    var left = Math.max(0, Math.round((this.deadline - Date.now()) / 1000));
    var t = this.root.querySelector(".exam-timer");
    if (!t) { clearInterval(this.timer); return; }
    t.textContent = "Tempo rimasto " + fmtTime(left);
    t.classList.toggle("is-low", left < 300);
    if (left === 0) { clearInterval(this.timer); alert("Tempo scaduto: l'esame viene consegnato."); this.finish(); }
  };

  McqSim.prototype.go = function (i) {
    if (i < 0 || i >= this.items.length) return;
    this.idx = i;
    var self = this, it = this.items[i], q = it.q, multi = q.a.length > 1;
    var box = this.root.querySelector(".exam-question");
    var html = '<p class="exam-domain">' + esc(domainName(this.bank, q.d)) + "</p>" +
      '<div class="exam-text">' + md(q.q) + (multi ? " <strong>" + multiHint(this.bank, q.a.length) + "</strong>" : "") + "</div>" +
      '<div class="exam-opts">';
    it.order.forEach(function (oi, pos) {
      html += '<label class="exam-opt"><input type="' + (multi ? "checkbox" : "radio") + '" name="opt" value="' + oi + '"' +
        (it.chosen.indexOf(oi) !== -1 ? " checked" : "") + (it.checked ? " disabled" : "") + "> <span><b>" +
        String.fromCharCode(65 + pos) + ".</b> " + md(q.o[oi]) + "</span></label>";
    });
    html += "</div>";
    if (this.instant) html += '<button class="md-button exam-check"' + (it.checked ? " hidden" : "") + ">Verifica la risposta</button><div class=\"exam-feedback\"></div>";
    box.innerHTML = html;
    box.querySelectorAll("input").forEach(function (inp) {
      inp.onchange = function () {
        var vals = Array.prototype.map.call(box.querySelectorAll("input:checked"), function (x) { return +x.value; });
        if (multi && vals.length > q.a.length) { inp.checked = false; return; }
        it.chosen = vals; self.palette();
      };
    });
    if (this.instant) {
      var chk = box.querySelector(".exam-check");
      chk.onclick = function () {
        if (!it.chosen.length) return;
        it.checked = true; self.showFeedback(box, it); chk.hidden = true;
        box.querySelectorAll("input").forEach(function (x) { x.disabled = true; });
        self.palette();
      };
      if (it.checked) this.showFeedback(box, it);
    }
    this.root.querySelector(".exam-pos").textContent = "Domanda " + (i + 1) + " di " + this.items.length;
    this.root.querySelector(".exam-flag input").checked = it.flag;
    this.root.querySelector(".exam-prev").disabled = i === 0;
    this.root.querySelector(".exam-next").disabled = i === this.items.length - 1;
    this.palette();
    var top = this.root.getBoundingClientRect().top;
    if (top < 0) this.root.scrollIntoView({ behavior: "smooth" });
  };

  function isRight(it) {
    return it.chosen.length === it.q.a.length && it.chosen.every(function (c) { return it.q.a.indexOf(c) !== -1; });
  }

  McqSim.prototype.showFeedback = function (box, it) {
    var ok = isRight(it);
    box.querySelectorAll(".exam-opt").forEach(function (l) {
      var v = +l.querySelector("input").value;
      l.classList.toggle("is-correct", it.q.a.indexOf(v) !== -1);
      l.classList.toggle("is-wrong", it.chosen.indexOf(v) !== -1 && it.q.a.indexOf(v) === -1);
    });
    box.querySelector(".exam-feedback").innerHTML = '<div class="exam-expl ' + (ok ? "ok" : "ko") + '"><p><strong>' +
      (ok ? "Corretto." : "Non corretto.") + "</strong></p>" + md(it.q.e) + refs(it.q.r) + "</div>";
    this.trackWrong(it.q.id, ok);
  };

  McqSim.prototype.trackWrong = function (id, ok) {
    var w = load(this.wrongKey, []);
    var i = w.indexOf(id);
    if (!ok && i === -1) w.push(id);
    if (ok && i !== -1 && this.mode === "ripasso") w.splice(i, 1);
    save(this.wrongKey, w);
  };

  McqSim.prototype.palette = function () {
    var self = this, p = this.root.querySelector(".exam-palette");
    p.innerHTML = "";
    this.items.forEach(function (it, i) {
      var b = document.createElement("button");
      b.type = "button";
      b.textContent = i + 1;
      b.className = "exam-pal" + (it.chosen.length ? " is-answered" : "") + (it.flag ? " is-flagged" : "") +
        (i === self.idx ? " is-current" : "") + (it.checked ? (isRight(it) ? " is-ok" : " is-ko") : "");
      b.onclick = function () { self.go(i); };
      p.appendChild(b);
    });
  };

  McqSim.prototype.finish = function () {
    clearInterval(this.timer);
    var self = this, b = this.bank;
    var right = 0, per = {};
    this.items.forEach(function (it) {
      var ok = isRight(it);
      if (ok) right++;
      var d = per[it.q.d] = per[it.q.d] || { ok: 0, n: 0 };
      d.n++; if (ok) d.ok++;
      if (self.mode === "esame") self.trackWrong(it.q.id, ok);
    });
    var pct = Math.round(right / this.items.length * 100);
    var passed = pct >= b.format.passPct;
    recordHistory(b, this.mode, pct, passed);
    var rows = b.domains.filter(function (d) { return per[d.id]; }).map(function (d) {
      var x = per[d.id], p = Math.round(x.ok / x.n * 100);
      return "<tr><td>" + esc(d.name) + "</td><td>" + x.ok + "/" + x.n + "</td><td>" + p + "%</td></tr>";
    }).join("");
    var review = this.items.map(function (it, i) {
      var ok = isRight(it);
      var opts = it.order.map(function (oi, pos) {
        var cls = it.q.a.indexOf(oi) !== -1 ? "is-correct" : (it.chosen.indexOf(oi) !== -1 ? "is-wrong" : "");
        var mine = it.chosen.indexOf(oi) !== -1 ? " ← la tua risposta" : "";
        return '<li class="' + cls + '"><b>' + String.fromCharCode(65 + pos) + ".</b> " + md(it.q.o[oi]) + "<em>" + mine + "</em></li>";
      }).join("");
      return '<details class="exam-review ' + (ok ? "ok" : "ko") + '"><summary>' + (ok ? "✅" : "❌") + " Domanda " + (i + 1) +
        " · " + esc(domainName(b, it.q.d)) + (it.chosen.length ? "" : " · senza risposta") + "</summary>" +
        '<div class="exam-text">' + md(it.q.q) + "</div><ul>" + opts + '</ul><div class="exam-expl">' + md(it.q.e) + refs(it.q.r) + "</div></details>";
    }).join("");
    this.root.innerHTML = '<div class="exam-card exam-result ' + (passed ? "ok" : "ko") + '">' +
      "<h3>" + (passed ? "Superato" : "Non superato") + ": " + right + " risposte corrette su " + this.items.length + " (" + pct + "%)</h3>" +
      "<p>Soglia indicativa di questa simulazione: " + b.format.passPct + "%. " + md(b.format.scoring || "") + "</p>" +
      "<table><thead><tr><th>Dominio</th><th>Corrette</th><th>%</th></tr></thead><tbody>" + rows + "</tbody></table>" +
      '<p><button class="md-button md-button--primary exam-again">Torna al menu</button></p>' +
      "<h4>Revisione domanda per domanda</h4>" + review + "</div>";
    this.root.querySelector(".exam-again").onclick = function () { self.home(); };
    this.root.scrollIntoView({ behavior: "smooth" });
  };

  // =====================================================================
  // Banche pratiche (CKA, LFCS)
  // =====================================================================
  function PracticalSim(root, bank) { this.root = root; this.bank = bank; this.home(); }

  PracticalSim.prototype.home = function () {
    var b = this.bank, f = b.format, self = this;
    clearInterval(this.timer);
    this.root.innerHTML = '<div class="exam-card"><h3>' + esc(b.title) + "</h3>" +
      '<p class="exam-meta">Formato reale: esame <strong>pratico</strong> da riga di comando, <strong>' + f.minutes +
      " minuti</strong>, circa " + esc(f.tasksLabel || f.tasks) + " task · soglia indicativa <strong>" + f.passPct + "%</strong> · " +
      b.tasks.length + " task originali in banca.</p>" +
      (f.note ? '<div class="exam-note">' + md(f.note) + "</div>" : "") +
      '<div class="exam-modes">' +
      '<div class="exam-mode"><h4>Esame pratico</h4><p>' + f.tasks + " task estratti a caso con i pesi dei domini e un timer di " + f.minutes +
      " minuti. Li svolgi nel tuo laboratorio; alla consegna vedi verifiche e soluzioni e assegni i punti.</p>" +
      '<button class="md-button md-button--primary" data-act="full">Inizia l\'esame pratico</button></div>' +
      '<div class="exam-mode"><h4>Allenamento</h4><p>Un task a caso alla volta, con verifica e soluzione disponibili quando vuoi.</p>' +
      '<button class="md-button" data-act="one">Dammi un task</button></div></div></div>';
    this.root.querySelector('[data-act="full"]').onclick = function () {
      self.start(weightedDraw(b, b.tasks, Math.min(f.tasks, b.tasks.length)), f.minutes);
    };
    this.root.querySelector('[data-act="one"]').onclick = function () {
      self.single(shuffle(b.tasks)[0]);
    };
  };

  function taskBody(b, t, i) {
    return '<div class="exam-task" data-i="' + i + '"><p class="exam-domain">' + esc(domainName(b, t.d)) + " · " + t.pts + " punti</p>" +
      "<h4>" + (i !== undefined ? "Task " + (i + 1) + " · " : "") + esc(t.title) + "</h4>" +
      (t.ctx ? '<p class="exam-ctx">' + md(t.ctx) + "</p>" : "") +
      '<div class="exam-text">' + md(t.task) + "</div>";
  }
  function taskSolution(t) {
    return '<details class="exam-sol"><summary>Come verificare</summary>' + md(t.check) + "</details>" +
      '<details class="exam-sol"><summary>Soluzione spiegata</summary>' + md(t.solution) + refs(t.r) + "</details>";
  }

  PracticalSim.prototype.single = function (t) {
    var self = this;
    this.root.innerHTML = '<div class="exam-card">' + taskBody(this.bank, t) + "</div>" + taskSolution(t) +
      '<p><button class="md-button md-button--primary exam-next1">Un altro task</button> <button class="md-button exam-back">Menu</button></p></div>';
    this.root.querySelector(".exam-next1").onclick = function () { self.single(shuffle(self.bank.tasks)[0]); };
    this.root.querySelector(".exam-back").onclick = function () { self.home(); };
  };

  PracticalSim.prototype.start = function (tasks, minutes) {
    var self = this, b = this.bank;
    this.tasks = tasks;
    this.deadline = Date.now() + minutes * 60000;
    var list = tasks.map(function (t, i) {
      return taskBody(b, t, i) + '<label class="exam-flag"><input type="checkbox" class="t-done"> Svolto</label> ' +
        '<label class="exam-flag"><input type="checkbox" class="t-flag"> Da rivedere</label></div>';
    }).join("");
    this.root.innerHTML = '<div class="exam-run"><div class="exam-bar"><span class="exam-pos">' + tasks.length + " task · " +
      tasks.reduce(function (s, t) { return s + t.pts; }, 0) + ' punti</span><span class="exam-timer"></span>' +
      '<button class="md-button exam-end">Consegna</button></div>' + list + "</div>";
    this.root.querySelector(".exam-end").onclick = function () { if (confirm("Vuoi consegnare l'esame?")) self.finish(); };
    this.tick();
    this.timer = setInterval(function () { self.tick(); }, 1000);
  };

  PracticalSim.prototype.tick = McqSim.prototype.tick;

  PracticalSim.prototype.finish = function () {
    clearInterval(this.timer);
    var self = this, b = this.bank;
    var total = this.tasks.reduce(function (s, t) { return s + t.pts; }, 0);
    var html = '<div class="exam-card"><h3>Correzione</h3><p>Per ogni task esegui i comandi di verifica nel tuo laboratorio, ' +
      "poi indica il risultato. Sii severo: all'esame un task conta solo se il risultato è esattamente quello richiesto.</p>";
    this.tasks.forEach(function (t, i) {
      html += taskBody(b, t, i) + taskSolution(t) +
        '<p class="exam-grade">Risultato: <label><input type="radio" name="g' + i + '" value="1"> riuscito</label> ' +
        '<label><input type="radio" name="g' + i + '" value="0.5"> in parte</label> ' +
        '<label><input type="radio" name="g' + i + '" value="0" checked> non riuscito</label></p></div>';
    });
    html += '<p><button class="md-button md-button--primary exam-score">Calcola il punteggio</button></p><div class="exam-final"></div></div>';
    this.root.innerHTML = html;
    this.root.querySelector(".exam-score").onclick = function () {
      var got = 0;
      self.tasks.forEach(function (t, i) {
        var v = self.root.querySelector('input[name="g' + i + '"]:checked');
        got += t.pts * (v ? +v.value : 0);
      });
      var pct = Math.round(got / total * 100), passed = pct >= b.format.passPct;
      recordHistory(b, "esame pratico", pct, passed);
      self.root.querySelector(".exam-final").innerHTML = '<div class="exam-result ' + (passed ? "ok" : "ko") + '"><h3>' +
        (passed ? "Superato" : "Non superato") + ": " + got + " punti su " + total + " (" + pct + "%)</h3><p>Soglia indicativa: " +
        b.format.passPct + '%.</p><button class="md-button exam-back">Torna al menu</button></div>';
      self.root.querySelector(".exam-back").onclick = function () { self.home(); };
    };
    this.root.scrollIntoView({ behavior: "smooth" });
  };

  // =====================================================================
  function init() {
    document.querySelectorAll(".exam-sim").forEach(function (root) {
      if (root.dataset.ready) return;
      root.dataset.ready = "1";
      if (!ASSETS) { root.textContent = "Simulatore non disponibile."; return; }
      root.innerHTML = "<p>Caricamento della banca domande…</p>";
      fetch(ASSETS + "banks/" + root.dataset.bank + ".json", { cache: "no-cache" })
        .then(function (r) { if (!r.ok) throw new Error(r.status); return r.json(); })
        .then(function (bank) {
          if (bank.type === "practical") new PracticalSim(root, bank); else new McqSim(root, bank);
        })
        .catch(function () { root.innerHTML = "<p>Impossibile caricare la banca domande. Ricarica la pagina.</p>"; });
    });
  }

  if (window.document$ && typeof window.document$.subscribe === "function") window.document$.subscribe(init);
  else document.addEventListener("DOMContentLoaded", init);
})();
