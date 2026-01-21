# TUAK (Go)

This repository provides a Go implementation of the 3GPP TUAK algorithm set,
including Keccak-f[1600] and the functions f1, f1*, f2, f3, f4, f5, f5*.
The public API is shaped to feel similar to `github.com/wmnsk/milenage`.
This codebase is developed with AI coding assistance using GPT-5.2-Codex.

Note on bit/byte ordering:
- External inputs/outputs are MSB-first as defined in 3GPP specs.
- Internal IN/OUT for Keccak are represented as 200 bytes with LSB-first bit
  ordering per TS 35.232.

## Usage

Compute TOPc and run f1/f1*/f2345/f5*:

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

## Debugging

You can capture intermediate IN/OUT buffers:

```go
t, _ := tuak.NewWithTOPc(
	k, topc, rand, sqn, amf,
	tuak.WithMACLength(64),
	tuak.WithDebugHook(func(label string, data []byte) {
		fmt.Printf("%s:\n%s\n", label, tuak.DebugHexBytes(data))
	}),
)
```

## Tests

Run all tests:

```sh
GOCACHE=/tmp/go-build go test ./...
```

The test vectors are stored under `testdata/`.

Regenerate testdata from reference text:

```sh
scripts/generate_testdata.py
```
