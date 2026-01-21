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

### Parameters and lengths

TUAK requires fixed lengths per deployment. The implementation expects these
lengths in bits and validates input byte sizes:

- `K`: 128 or 256 bits (`len(K) == 16` or `32`)
- `TOP`/`TOPc`: 256 bits (`len == 32`)
- `RAND`: 128 bits (`len == 16`)
- `SQN`: 48 bits (`len == 6`)
- `AMF`: 16 bits (`len == 2`)

Output lengths (bits):

- `MAC` (f1/f1*): 64, 128, or 256
- `RES` (f2): 32, 64, 128, or 256
- `CK` (f3): 128 or 256
- `IK` (f4): 128 or 256
- `AK`/`AK*` (f5/f5*): 48 (always 6 bytes)

If `WithKLength` is omitted, it is inferred from `len(K)`; `KeccakIterations`
defaults to 1.

### API overview

- `ComputeTOPc(k, top, opts...)` derives TOPc from K and TOP.
- `New(k, top, rand, sqn, amf, opts...)` creates a context (TOPc computed as needed).
- `NewWithTOPc(k, topc, rand, sqn, amf, opts...)` creates a context with precomputed TOPc.
- `F1()` returns MAC-A (byte length = `MACLength/8`).
- `F1Star()` returns MAC-S (byte length = `MACLength/8`).
- `F2345()` returns `(RES, CK, IK, AK)` using `RESLength/CKLength/IKLength`.
- `F5Star()` returns AK* (always 6 bytes).

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
