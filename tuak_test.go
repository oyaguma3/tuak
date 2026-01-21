package tuak

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"

	"tuak/testvectors"
)

func TestComputeTOPcVectors(t *testing.T) {
	data, err := testvectors.LoadTUAKVectors()
	if err != nil {
		t.Fatalf("LoadTUAKVectors: %v", err)
	}
	for _, v := range data.Tests {
		k := decodeHex(t, v.K)
		top := decodeHex(t, v.Top)
		want := decodeHex(t, v.Topc)
		got, err := ComputeTOPc(k, top, optionsFromVector(v)...)
		if errors.Is(err, ErrNotImplemented) {
			t.Skip("ComputeTOPc not implemented")
		}
		if err != nil {
			t.Fatalf("ComputeTOPc: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("vector %d TOPc mismatch", v.ID)
		}
	}
}

func TestF1F1StarVectors(t *testing.T) {
	data, err := testvectors.LoadTUAKVectors()
	if err != nil {
		t.Fatalf("LoadTUAKVectors: %v", err)
	}
	for _, v := range data.Tests {
		tuak, err := newTUAKFromVector(t, v)
		if err != nil {
			t.Fatalf("NewWithTOPc: %v", err)
		}
		wantF1 := decodeHex(t, v.F1)
		wantF1Star := decodeHex(t, v.F1Star)

		gotF1, err := tuak.F1()
		if errors.Is(err, ErrNotImplemented) {
			t.Skip("F1 not implemented")
		}
		if err != nil {
			t.Fatalf("F1: %v", err)
		}
		if !bytes.Equal(gotF1, wantF1) {
			t.Fatalf("vector %d f1 mismatch", v.ID)
		}

		gotF1Star, err := tuak.F1Star()
		if errors.Is(err, ErrNotImplemented) {
			t.Skip("F1* not implemented")
		}
		if err != nil {
			t.Fatalf("F1*: %v", err)
		}
		if !bytes.Equal(gotF1Star, wantF1Star) {
			t.Fatalf("vector %d f1* mismatch", v.ID)
		}
	}
}

func TestF2345Vectors(t *testing.T) {
	data, err := testvectors.LoadTUAKVectors()
	if err != nil {
		t.Fatalf("LoadTUAKVectors: %v", err)
	}
	for _, v := range data.Tests {
		tuak, err := newTUAKFromVector(t, v)
		if err != nil {
			t.Fatalf("NewWithTOPc: %v", err)
		}
		wantF2 := decodeHex(t, v.F2)
		wantF3 := decodeHex(t, v.F3)
		wantF4 := decodeHex(t, v.F4)
		wantF5 := decodeHex(t, v.F5)

		res, ck, ik, ak, err := tuak.F2345()
		if errors.Is(err, ErrNotImplemented) {
			t.Skip("F2345 not implemented")
		}
		if err != nil {
			t.Fatalf("F2345: %v", err)
		}
		if !bytes.Equal(res, wantF2) {
			t.Fatalf("vector %d f2 mismatch", v.ID)
		}
		if !bytes.Equal(ck, wantF3) {
			t.Fatalf("vector %d f3 mismatch", v.ID)
		}
		if !bytes.Equal(ik, wantF4) {
			t.Fatalf("vector %d f4 mismatch", v.ID)
		}
		if !bytes.Equal(ak, wantF5) {
			t.Fatalf("vector %d f5 mismatch", v.ID)
		}
	}
}

func TestF5StarVectors(t *testing.T) {
	data, err := testvectors.LoadTUAKVectors()
	if err != nil {
		t.Fatalf("LoadTUAKVectors: %v", err)
	}
	for _, v := range data.Tests {
		tuak, err := newTUAKFromVector(t, v)
		if err != nil {
			t.Fatalf("NewWithTOPc: %v", err)
		}
		want := decodeHex(t, v.F5Star)
		got, err := tuak.F5Star()
		if errors.Is(err, ErrNotImplemented) {
			t.Skip("F5* not implemented")
		}
		if err != nil {
			t.Fatalf("F5*: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("vector %d f5* mismatch", v.ID)
		}
	}
}

func newTUAKFromVector(t *testing.T, v testvectors.TUAKVector) (*TUAK, error) {
	k := decodeHex(t, v.K)
	topc := decodeHex(t, v.Topc)
	rand := decodeHex(t, v.Rand)
	sqn := decodeHex(t, v.SQN)
	amf := decodeHex(t, v.AMF)
	return NewWithTOPc(k, topc, rand, sqn, amf, optionsFromVector(v)...)
}

func optionsFromVector(v testvectors.TUAKVector) []Option {
	return []Option{
		WithKLength(v.Klength),
		WithMACLength(v.MAClength),
		WithRESLength(v.RESLength),
		WithCKLength(v.CKlength),
		WithIKLength(v.IKlength),
		WithKeccakIterations(v.KeccakIterations),
	}
}

func decodeHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("decode hex: %v", err)
	}
	return b
}
