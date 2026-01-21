# TUAK (Go)

このリポジトリは、3GPP TUAK アルゴリズムセットの Go 実装です。
Keccak-f[1600] と、f1 / f1* / f2 / f3 / f4 / f5 / f5* を提供します。
API 形状は `github.com/wmnsk/milenage` と同様の使い勝手を意識しています。
本コードベースは GPT-5.2-Codex を用いた AI コーディングで開発しています。

ビット/バイト順について:
- 外部入力/出力は 3GPP 仕様の MSB-first 表記です。
- Keccak の IN/OUT は TS 35.232 に従い 200 バイト・LSB-first で表現します。

## 使い方

TOPc の導出と f1/f1*/f2345/f5* の例:

```go
topc, err := tuak.ComputeTOPc(k, top, tuak.WithKLength(128))
if err != nil {
	// handle err
}

t, err := tuak.NewWithTOPc(k, topc, rand, sqn, amf,
	tuak.WithKLength(128),
	tuak.WithMACLength(64),
	tuak.WithRESLength(32),
	tuak.WithCKLength(128),
	tuak.WithIKLength(128),
	tuak.WithKeccakIterations(1),
)
if err != nil {
	// handle err
}

macA, _ := t.F1()
macS, _ := t.F1Star()
res, ck, ik, ak, _ := t.F2345()
akStar, _ := t.F5Star()
```

## デバッグ

中間 IN/OUT を取得する場合:

```go
t, _ := tuak.NewWithTOPc(
	k, topc, rand, sqn, amf,
	tuak.WithMACLength(64),
	tuak.WithDebugHook(func(label string, data []byte) {
		fmt.Printf("%s:\n%s\n", label, tuak.DebugHexBytes(data))
	}),
)
```

## テスト

全テスト実行:

```sh
GOCACHE=/tmp/go-build go test ./...
```

テストベクトルは `testdata/` にあります。

テストデータの再生成:

```sh
scripts/generate_testdata.py
```
