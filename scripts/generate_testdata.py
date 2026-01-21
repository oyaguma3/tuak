#!/usr/bin/env python3
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
REFERENCE = ROOT / "reference"
TESTDATA = ROOT / "testdata"


def write_json(path: Path, data: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(data, indent=2, sort_keys=True), encoding="utf-8")


def parse_35232_keccak(text: str) -> list[dict]:
    lines = text.splitlines()
    re_test = re.compile(r"^5\.\d+\s+Test set\s+(\d+)", re.IGNORECASE)
    re_hex = re.compile(r"\b[0-9a-fA-F]{2}\b")

    vectors = []
    current = None
    in_mode = False
    out_mode = False
    in_section = False

    for line in lines:
        s = line.strip()
        m = re_test.match(s)
        if m:
            if current:
                vectors.append(current)
            current = {"id": int(m.group(1)), "in": "", "out": ""}
            in_section = True
            in_mode = False
            out_mode = False
            continue
        if in_section and s.startswith("6"):
            in_section = False
            continue
        if not in_section or current is None:
            continue
        if s.startswith("IN"):
            in_mode = True
            out_mode = False
            continue
        if s.startswith("OUT"):
            out_mode = True
            in_mode = False
            continue
        if not (in_mode or out_mode):
            continue
        tokens = re_hex.findall(s)
        if not tokens:
            continue
        hexstr = "".join(t.lower() for t in tokens)
        if in_mode:
            current["in"] += hexstr
        else:
            current["out"] += hexstr

    if current:
        vectors.append(current)

    for v in vectors:
        v["in_len"] = len(v["in"]) // 2
        v["out_len"] = len(v["out"]) // 2
    return vectors


def parse_35233_tuak(text: str) -> list[dict]:
    lines = text.splitlines()
    re_test = re.compile(r"^6\.\d+\s+Test set\s+(\d+)", re.IGNORECASE)
    re_kv = re.compile(r"^([A-Za-z0-9\*]+):\s+([0-9a-fA-F]+)\b")

    vectors = []
    current = None
    in_binary = False

    for line in lines:
        s = line.strip()
        m = re_test.match(s)
        if m:
            if current:
                vectors.append(current)
            current = {"id": int(m.group(1))}
            in_binary = False
            continue
        if current is None:
            continue
        if s.startswith("Binary Format"):
            in_binary = True
            continue
        if in_binary:
            continue
        if "bits" in s or "KeccakIterations" in s:
            for key, field in [
                ("Klength", "klength"),
                ("MAClength", "maclength"),
                ("CKlength", "cklength"),
                ("IKlength", "iklength"),
                ("RESLength", "reslength"),
            ]:
                mlen = re.search(rf"\b{key}\b\s*=\s*(\d+)", s)
                if mlen:
                    current[field] = int(mlen.group(1))
            miter = re.search(r"\bKeccakIterations\b\s*=\s*(\d+)", s)
            if miter:
                current["keccak_iterations"] = int(miter.group(1))
            continue
        m = re_kv.match(s)
        if m:
            key = m.group(1)
            val = m.group(2).lower()
            key_map = {
                "K": "k",
                "RAND": "rand",
                "SQN": "sqn",
                "AMF": "amf",
                "TOP": "top",
                "TOPc": "topc",
                "f1": "f1",
                "f1*": "f1_star",
                "f2": "f2",
                "f3": "f3",
                "f4": "f4",
                "f5": "f5",
                "f5*": "f5_star",
            }
            if key in key_map and key_map[key] not in current:
                current[key_map[key]] = val

    if current:
        vectors.append(current)
    return vectors


def main() -> None:
    text_35232 = (REFERENCE / "35232-i00_chapter4-7.txt").read_text(encoding="utf-8")
    text_35233 = (REFERENCE / "35233-i00_chapter5-6.txt").read_text(encoding="utf-8")

    keccak = parse_35232_keccak(text_35232)
    write_json(
        TESTDATA / "ts35232_keccak.json",
        {
            "source": "3GPP TS 35.232 V18.0.0 (Implementers test data) - text extract",
            "keccak_f1600": keccak,
        },
    )

    # 35.233 uses the same Keccak vectors; keep a stable copy.
    write_json(
        TESTDATA / "ts35233_keccak.json",
        {
            "source": "3GPP TS 35.233 V18.0.0 (Design conformance test data) - text extract",
            "keccak_f1600": keccak,
        },
    )

    tuak = parse_35233_tuak(text_35233)
    write_json(
        TESTDATA / "ts35233_vectors_text.json",
        {
            "source": "3GPP TS 35.233 V18.0.0 (Design conformance test data) - text extract",
            "tests": tuak,
        },
    )


if __name__ == "__main__":
    main()
