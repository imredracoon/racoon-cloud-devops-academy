/*
 * Racoon Cloud & DevOps Academy
 * - quiz interattivi dentro le lezioni (<script type="application/json" class="quiz-data">)
 * - pulsante "Segna come completata" per ogni lezione delle fasi
 * - pagina "I miei progressi"
 * Tutto lo stato resta nel browser di chi studia (localStorage).
 */
(function () {
  "use strict";

  var STORE = "academy.progress.v1";

  // ---------- utilità ----------
  function load(key, fallback) {
    try {
      var raw = localStorage.getItem(key);
      return raw ? JSON.parse(raw) : fallback;
    } catch (e) {
      return fallback;
    }
  }
  function save(key, value) {
    try { localStorage.setItem(key, JSON.stringify(value)); } catch (e) { /* storage non disponibile */ }
  }
  function esc(s) {
    return String(s).replace(/[&<>"']/g, function (c) {
      return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c];
    });
  }
  // Markdown minimo: blocchi ```, `codice`, **grassetto**, a capo.
  function md(s) {
    var parts = String(s).split("```");
    return parts.map(function (p, i) {
      if (i % 2 === 1) {
        var body = p.replace(/^[a-z]*\n/, "");
        return "<pre><code>" + esc(body.replace(/\n$/, "")) + "</code></pre>";
      }
      return esc(p)
        .replace(/`([^`]+)`/g, "<code>$1</code>")
        .replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>")
        .replace(/\n/g, "<br>");
    }).join("");
  }
  window.AcademyMD = md;
  function el(tag, cls, html) {
    var e = document.createElement(tag);
    if (cls) e.className = cls;
    if (html !== undefined) e.innerHTML = html;
    return e;
  }
  function pageKey(path) {
    var p = (path || location.pathname).replace(/index\.html$/, "");
    var m = p.match(/(fase-\d+\/[^/]*\/?)$/);
    if (!m) return null;
    var k = m[1].replace(/\/$/, "");
    return k.indexOf("/") === -1 ? null : k; // le pagine indice delle fasi non contano
  }

  // ---------- quiz in linea ----------
  function renderQuiz(script) {
    if (script.dataset.rendered) return;
    script.dataset.rendered = "1";
    var data;
    try { data = JSON.parse(script.textContent); } catch (e) {
      script.insertAdjacentHTML("afterend", '<p class="quiz-error">Quiz non valido.</p>');
      return;
    }
    var box = el("div", "academy-quiz");
    box.appendChild(el("p", "academy-quiz__title", esc(data.title || "Verifica rapida")));
    var correctCount = 0, answered = 0;
    var scoreLine = el("p", "academy-quiz__score");

    data.questions.forEach(function (q, qi) {
      var multi = q.answer.length > 1;
      var item = el("div", "academy-quiz__q");
      item.appendChild(el("p", "academy-quiz__text",
        "<strong>" + (qi + 1) + ".</strong> " + md(q.q) +
        (multi ? ' <em class="academy-quiz__hint">(scegli ' + q.answer.length + ")</em>" : "")));
      var name = "q" + Math.random().toString(36).slice(2);
      var list = el("div", "academy-quiz__opts");
      q.options.forEach(function (o, oi) {
        var lab = el("label", "academy-quiz__opt");
        lab.innerHTML = '<input type="' + (multi ? "checkbox" : "radio") + '" name="' + name +
          '" value="' + oi + '"> <span>' + md(o) + "</span>";
        list.appendChild(lab);
      });
      item.appendChild(list);
      var btn = el("button", "md-button academy-quiz__check", "Verifica");
      btn.type = "button";
      var fb = el("div", "academy-quiz__feedback");
      btn.addEventListener("click", function () {
        var chosen = Array.prototype.map.call(
          list.querySelectorAll("input:checked"), function (i) { return +i.value; });
        if (!chosen.length) { fb.innerHTML = "<p>Seleziona almeno una risposta.</p>"; return; }
        var ok = chosen.length === q.answer.length &&
          chosen.every(function (c) { return q.answer.indexOf(c) !== -1; });
        list.querySelectorAll("label").forEach(function (l, li) {
          l.classList.toggle("is-correct", q.answer.indexOf(li) !== -1);
          l.classList.toggle("is-wrong", chosen.indexOf(li) !== -1 && q.answer.indexOf(li) === -1);
          l.querySelector("input").disabled = true;
        });
        btn.disabled = true;
        answered++; if (ok) correctCount++;
        fb.className = "academy-quiz__feedback " + (ok ? "ok" : "ko");
        fb.innerHTML = "<p><strong>" + (ok ? "Corretto." : "Non corretto.") + "</strong> " +
          md(q.explanation || "") + "</p>";
        scoreLine.textContent = "Punteggio: " + correctCount + " su " + answered +
          " risposte date (" + data.questions.length + " domande in tutto).";
      });
      item.appendChild(btn);
      item.appendChild(fb);
      box.appendChild(item);
    });
    box.appendChild(scoreLine);
    script.parentNode.insertBefore(box, script.nextSibling);
  }

  // ---------- completamento lezioni ----------
  function addCompleteButton() {
    var key = pageKey();
    var content = document.querySelector("article.md-content__inner");
    if (!key || !content || content.querySelector(".academy-done")) return;
    var state = load(STORE, {});
    var wrap = el("div", "academy-done");
    var btn = el("button", "md-button");
    btn.type = "button";
    function paint() {
      var done = !!state[key];
      btn.textContent = done ? "Lezione completata ✓ (clic per annullare)" : "Segna come completata";
      btn.classList.toggle("md-button--primary", !done);
    }
    btn.addEventListener("click", function () {
      state = load(STORE, {});
      if (state[key]) delete state[key]; else state[key] = new Date().toISOString();
      save(STORE, state);
      paint();
    });
    paint();
    wrap.appendChild(btn);
    content.appendChild(wrap);
  }

  // ---------- pagina progressi ----------
  function renderProgress() {
    var host = document.getElementById("academy-progress");
    if (!host || host.dataset.rendered) return;
    host.dataset.rendered = "1";
    var state = load(STORE, {});
    var links = document.querySelectorAll(".md-nav--primary a.md-nav__link[href]");
    var phases = {}, order = [];
    Array.prototype.forEach.call(links, function (a) {
      var url = new URL(a.getAttribute("href"), location.href);
      var key = pageKey(url.pathname);
      if (!key) return;
      var phase = key.split("/")[0];
      if (!phases[phase]) { phases[phase] = { pages: {}, list: [] }; order.push(phase); }
      if (phases[phase].pages[key]) return;
      phases[phase].pages[key] = true;
      phases[phase].list.push({ key: key, title: a.textContent.trim(), href: url.href });
    });
    var total = 0, done = 0, html = "";
    order.sort().forEach(function (ph) {
      var list = phases[ph].list;
      var d = list.filter(function (p) { return state[p.key]; }).length;
      total += list.length; done += d;
      var pct = list.length ? Math.round(d / list.length * 100) : 0;
      html += '<details class="academy-progress__phase"><summary><strong>' +
        esc(ph.replace("fase-", "Fase ")) + "</strong> · " + d + "/" + list.length +
        ' <span class="academy-bar"><span style="width:' + pct + '%"></span></span></summary><ul>';
      list.forEach(function (p) {
        html += "<li>" + (state[p.key] ? "✅ " : "⬜ ") + '<a href="' + esc(p.href) + '">' +
          esc(p.title) + "</a></li>";
      });
      html += "</ul></details>";
    });
    var pctAll = total ? Math.round(done / total * 100) : 0;
    host.innerHTML = '<p class="academy-progress__total">Hai completato <strong>' + done +
      "</strong> pagine su <strong>" + total + "</strong> (" + pctAll + '%).</p>' +
      '<span class="academy-bar academy-bar--big"><span style="width:' + pctAll + '%"></span></span>' + html +
      '<p><button type="button" class="md-button" id="academy-export">Esporta progressi</button> ' +
      '<label class="md-button">Importa progressi<input type="file" id="academy-import" accept="application/json" hidden></label> ' +
      '<button type="button" class="md-button" id="academy-reset">Azzera tutto</button></p>';
    document.getElementById("academy-export").onclick = function () {
      var blob = new Blob([JSON.stringify({
        progress: load(STORE, {}), exams: load("academy.exams.v1", [])
      }, null, 2)], { type: "application/json" });
      var a = document.createElement("a");
      a.href = URL.createObjectURL(blob);
      a.download = "academy-progressi.json";
      a.click();
    };
    document.getElementById("academy-import").onchange = function (ev) {
      var f = ev.target.files[0]; if (!f) return;
      f.text().then(function (t) {
        try {
          var j = JSON.parse(t);
          if (j.progress) save(STORE, j.progress);
          if (j.exams) save("academy.exams.v1", j.exams);
          host.dataset.rendered = ""; renderProgress();
        } catch (e) { alert("File non valido."); }
      });
    };
    document.getElementById("academy-reset").onclick = function () {
      if (confirm("Vuoi davvero cancellare tutti i progressi salvati in questo browser?")) {
        save(STORE, {}); save("academy.exams.v1", []);
        host.dataset.rendered = ""; renderProgress();
      }
    };
    var hist = document.getElementById("academy-exam-history");
    if (hist) {
      var exams = load("academy.exams.v1", []);
      hist.innerHTML = exams.length ? "<table><thead><tr><th>Data</th><th>Esame</th><th>Modalità</th><th>Risultato</th></tr></thead><tbody>" +
        exams.slice().reverse().map(function (e) {
          return "<tr><td>" + esc(new Date(e.date).toLocaleString("it-IT")) + "</td><td>" + esc(e.exam) +
            "</td><td>" + esc(e.mode) + "</td><td>" + esc(e.pct + "% " + (e.passed ? "✅" : "❌")) + "</td></tr>";
        }).join("") + "</tbody></table>" : "<p>Nessuna simulazione svolta finora.</p>";
    }
  }

  function init() {
    document.querySelectorAll("script.quiz-data").forEach(renderQuiz);
    addCompleteButton();
    renderProgress();
  }

  if (window.document$ && typeof window.document$.subscribe === "function") {
    window.document$.subscribe(init);
  } else {
    document.addEventListener("DOMContentLoaded", init);
  }
})();
