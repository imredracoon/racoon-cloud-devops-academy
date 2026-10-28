"""Controlla la coerenza delle banche domande in docs/assets/banks/.

Uso: python scripts/check_banks.py
Esce con codice 1 se trova errori, così la pipeline si ferma.
"""
import json
import pathlib
import sys

BANKS = pathlib.Path(__file__).resolve().parent.parent / "docs" / "assets" / "banks"


def check_mcq(name, bank, errors):
    domains = {d["id"] for d in bank["domains"]}
    weights = sum(d["weight"] for d in bank["domains"])
    if weights != 100:
        errors.append(f"{name}: i pesi dei domini sommano a {weights}, non a 100")
    seen = set()
    for q in bank["questions"]:
        qid = q.get("id", "?")
        if qid in seen:
            errors.append(f"{name}: id duplicato {qid}")
        seen.add(qid)
        if q.get("d") not in domains:
            errors.append(f"{name}/{qid}: dominio sconosciuto {q.get('d')}")
        opts = q.get("o", [])
        if len(opts) < 4:
            errors.append(f"{name}/{qid}: servono almeno 4 opzioni")
        if len(set(opts)) != len(opts):
            errors.append(f"{name}/{qid}: opzioni duplicate")
        ans = q.get("a", [])
        if not ans or any(not 0 <= a < len(opts) for a in ans) or len(set(ans)) != len(ans):
            errors.append(f"{name}/{qid}: risposta non valida {ans}")
        if len(ans) > 1 and len(opts) < 5:
            errors.append(f"{name}/{qid}: le domande a risposta multipla richiedono almeno 5 opzioni")
        if not q.get("e"):
            errors.append(f"{name}/{qid}: manca la spiegazione")
        if not q.get("r"):
            errors.append(f"{name}/{qid}: manca almeno una fonte")
    need = bank["format"]["questions"]
    if len(bank["questions"]) < need:
        errors.append(f"{name}: {len(bank['questions'])} domande, ne servono almeno {need}")
    return len(bank["questions"])


def check_practical(name, bank, errors):
    domains = {d["id"] for d in bank["domains"]}
    seen = set()
    for t in bank["tasks"]:
        tid = t.get("id", "?")
        if tid in seen:
            errors.append(f"{name}: id duplicato {tid}")
        seen.add(tid)
        if t.get("d") not in domains:
            errors.append(f"{name}/{tid}: dominio sconosciuto {t.get('d')}")
        for field in ("title", "task", "check", "solution", "pts", "r"):
            if not t.get(field):
                errors.append(f"{name}/{tid}: manca il campo {field}")
    if len(bank["tasks"]) < bank["format"]["tasks"]:
        errors.append(f"{name}: troppi pochi task")
    return len(bank["tasks"])


def main():
    errors = []
    for path in sorted(BANKS.glob("*.json")):
        try:
            bank = json.loads(path.read_text(encoding="utf-8"))
        except json.JSONDecodeError as exc:
            errors.append(f"{path.name}: JSON non valido ({exc})")
            continue
        if bank.get("type") == "practical":
            n = check_practical(path.stem, bank, errors)
            print(f"{path.stem}: {n} task")
        else:
            n = check_mcq(path.stem, bank, errors)
            per = {}
            for q in bank["questions"]:
                per[q["d"]] = per.get(q["d"], 0) + 1
            print(f"{path.stem}: {n} domande {per}")
    for e in errors:
        print("ERRORE:", e)
    sys.exit(1 if errors else 0)


if __name__ == "__main__":
    main()
