"""Genera la versione PDF del corso a partire dal sito già costruito (cartella site/).

Il PDF contiene solo la parte teorica: restano fuori esercizi, laboratori, soluzioni, quiz, progetti,
verifiche, simulazioni d'esame e contenuti multimediali (diagrammi, immagini, video).

Uso:
    mkdocs build --strict
    python scripts/build_pdf.py            # scrive site/assets/pdf/racoon-cloud-devops-academy.pdf

Richiede Playwright con Chromium (pip install playwright && python -m playwright install chromium).
"""

import datetime
import functools
import html
import http.server
import os
import sys
import threading

from mkdocs.config import load_config
from playwright.sync_api import sync_playwright

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SITE = os.path.join(ROOT, "site")
OUT_NAME = "racoon-cloud-devops-academy.pdf"
OUT = os.path.join(SITE, "assets", "pdf", OUT_NAME)

# Pagine escluse per nome: sono esercitazioni o strumenti interattivi.
EXCLUDED = {"progetto.md", "verifica.md", "simulazione.md", "progressi.md"}

# Pulizia eseguita in ogni pagina: restituisce l'HTML del solo contenuto teorico.
CLEAN_JS = """
(site) => {
  const root = document.querySelector("article.md-content__inner");
  const drop = /^(Esercizi|Esercizio|Laboratorio|Laboratori|Verifica rapida|E adesso\\?)$/i;
  root.querySelectorAll("h2").forEach((h) => {
    const title = h.textContent.replace("¶", "").trim();
    if (!drop.test(title)) return;
    let n = h.nextElementSibling;
    while (n && n.tagName !== "H2") { const next = n.nextElementSibling; n.remove(); n = next; }
    h.remove();
  });
  root.querySelectorAll([
    ".md-content__button", ".headerlink", ".md-clipboard", ".md-source-file", ".academy-done",
    ".esercizio", ".laboratorio", "details.success",
    ".academy-quiz", "script", ".exam-sim",
    ".mermaid", "img", "svg.twemoji", "video", "audio", "iframe", "figure",
  ].join(",")).forEach((e) => e.remove());
  root.querySelectorAll("details").forEach((d) => d.setAttribute("open", ""));
  // Le schede (Debian/Ubuntu, RHEL, AWS, Google Cloud) diventano blocchi in sequenza.
  root.querySelectorAll(".tabbed-set").forEach((set) => {
    const labels = [...set.querySelectorAll(":scope > .tabbed-labels > label")].map((l) => l.textContent.trim());
    const blocks = [...set.querySelectorAll(":scope > .tabbed-content > .tabbed-block")];
    const wrap = document.createElement("div");
    wrap.className = "pdf-tabs";
    blocks.forEach((b, i) => {
      const t = document.createElement("p");
      t.className = "pdf-tab-label";
      t.textContent = labels[i] || "";
      wrap.appendChild(t);
      const d = document.createElement("div");
      d.innerHTML = b.innerHTML;
      wrap.appendChild(d);
    });
    set.replaceWith(wrap);
  });
  // I collegamenti puntano al sito pubblicato; gli id vengono tolti per evitare duplicati.
  root.querySelectorAll("a[href]").forEach((a) => {
    a.setAttribute("href", a.href.replace(location.origin + "/", site));
  });
  root.querySelectorAll("[id]").forEach((e) => e.removeAttribute("id"));
  return root.innerHTML;
}
"""

PRINT_CSS = """
@page { size: A4; margin: 18mm 16mm 20mm 16mm; }
body { background: #fff; }
.md-typeset { font-size: 10pt; line-height: 1.55; }
.md-typeset pre > code { white-space: pre-wrap; word-break: break-word; }
.md-typeset table:not([class]) { font-size: 0.85em; display: table; }
.md-typeset h1, .md-typeset h2, .md-typeset h3 { break-after: avoid; }
.md-typeset .admonition, .md-typeset details, .md-typeset pre, .md-typeset tr { break-inside: avoid; }
.pdf-page { break-before: page; }
.pdf-tab-label { font-weight: 700; margin: 1em 0 0.3em; border-bottom: 1px solid #ccc; }
.pdf-cover { height: 240mm; display: flex; flex-direction: column; justify-content: center; }
.pdf-cover h1 { font-size: 2.6em; margin-bottom: 0.3em; }
.pdf-toc ul { list-style: none; padding-left: 0; }
.pdf-toc li li { padding-left: 1.2em; }
.pdf-toc a { text-decoration: none; color: inherit; }
"""


def flatten(nav, section=None):
    """Restituisce (sezione, titolo, file) nell'ordine del menu."""
    for item in nav:
        if isinstance(item, str):
            yield section, None, item
        elif isinstance(item, dict):
            for title, value in item.items():
                if isinstance(value, list):
                    yield from flatten(value, title if section is None else section)
                else:
                    yield (title if section is None else section), title, value


def page_url(md):
    if md == "index.md":
        return ""
    if md.endswith("/index.md"):
        return md[: -len("index.md")]
    return md[: -len(".md")] + "/"


def main():
    if not os.path.isdir(SITE):
        sys.exit("Cartella site/ assente: esegui prima mkdocs build.")
    config = load_config(os.path.join(ROOT, "mkdocs.yml"))
    site_url = config["site_url"]
    pages = [
        (section, md) for section, _, md in flatten(config["nav"])
        if md != "index.md" and os.path.basename(md) not in EXCLUDED
    ]

    class Quiet(http.server.SimpleHTTPRequestHandler):
        def log_message(self, *args):
            pass

    handler = functools.partial(Quiet, directory=SITE)
    server = http.server.ThreadingHTTPServer(("127.0.0.1", 0), handler)
    threading.Thread(target=server.serve_forever, daemon=True).start()
    base = f"http://127.0.0.1:{server.server_address[1]}/"

    with sync_playwright() as p:
        browser = p.chromium.launch(executable_path=os.environ.get("CHROMIUM_PATH") or None)
        page = browser.new_page()
        page.route("**/*", lambda r: r.continue_() if r.request.url.startswith(base) else r.abort())

        sections, toc = [], []
        styles = None
        for i, (section, md) in enumerate(pages):
            page.goto(base + page_url(md), wait_until="load")
            if styles is None:
                styles = page.eval_on_selector_all(
                    "link[rel=stylesheet]", "ls => ls.map(l => l.href)")
            title = page.inner_text("article.md-content__inner h1").replace("¶", "").strip()
            body = page.evaluate(CLEAN_JS, site_url)
            anchor = f"p{i}"
            sections.append(f'<section class="pdf-page" id="{anchor}">{body}</section>')
            if not toc or toc[-1][0] != section:
                toc.append((section, []))
            toc[-1][1].append((anchor, title))

        toc_html = "".join(
            f"<li><strong>{html.escape(s or '')}</strong><ul>"
            + "".join(f'<li><a href="#{a}">{html.escape(t)}</a></li>' for a, t in items)
            + "</ul></li>"
            for s, items in toc
        )
        today = datetime.date.today().strftime("%d/%m/%Y")
        doc = f"""<!doctype html><html lang="it"><head><meta charset="utf-8">
<title>{html.escape(config['site_name'])}</title>
{''.join(f'<link rel="stylesheet" href="{h}">' for h in styles)}
<style>{PRINT_CSS}</style></head>
<body><div class="md-typeset">
<div class="pdf-cover">
  <h1>{html.escape(config['site_name'])}</h1>
  <p>Corso gratuito per diventare Cloud / DevOps Engineer partendo da zero.</p>
  <p>Versione PDF con la sola parte teorica: esercizi, laboratori, quiz, progetti e simulazioni
  d'esame sono disponibili sul sito, insieme a diagrammi e contenuti interattivi.</p>
  <p><a href="{site_url}">{site_url}</a><br>Generato il {today}</p>
</div>
<section class="pdf-page pdf-toc"><h1>Indice</h1><ul>{toc_html}</ul></section>
{''.join(sections)}
</div></body></html>"""

        tmp = os.path.join(SITE, "_pdf.html")
        with open(tmp, "w", encoding="utf-8") as f:
            f.write(doc)
        page.goto(base + "_pdf.html", wait_until="load")
        os.makedirs(os.path.dirname(OUT), exist_ok=True)
        page.pdf(
            path=OUT,
            format="A4",
            print_background=True,
            prefer_css_page_size=True,
            display_header_footer=True,
            header_template="<span></span>",
            footer_template=(
                '<div style="font-size:8px;width:100%;text-align:center;color:#666">'
                '<span class="pageNumber"></span> / <span class="totalPages"></span></div>'
            ),
        )
        os.remove(tmp)
        browser.close()
    server.shutdown()
    print(f"PDF: {os.path.relpath(OUT, ROOT)} ({len(pages)} pagine del sito, "
          f"{os.path.getsize(OUT) // 1024} KB)")


if __name__ == "__main__":
    main()
