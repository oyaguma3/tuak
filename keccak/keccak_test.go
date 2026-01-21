package keccak

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
	"tuak/testvectors"
)

func TestPermuteF1600Vectors(t *testing.T) {
	data, err := testvectors.LoadKeccakVectors()
	if err != nil {
		t.Fatalf("LoadKeccakVectors: %v", err)
	}
	for _, v := range data.KeccakF1600 {
		in := decodeHex(t, v.In)
		want := decodeHex(t, v.Out)
		got, err := PermuteF1600(in)
		if err != nil {
			t.Fatalf("PermuteF1600: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("vector %d mismatch", v.ID)
		}
	}
}

func TestPermuteF1600VectorsFrom35233JSON(t *testing.T) {
	vectors, err := loadKeccakVectors35233JSON()
	if err != nil {
		t.Fatalf("load 35233 json vectors: %v", err)
	}
	if len(vectors) == 0 {
		t.Fatalf("no keccak vectors parsed from 35233 json")
	}
	for _, v := range vectors {
		in := decodeHex(t, v.in)
		want := decodeHex(t, v.out)
		got, err := PermuteF1600(in)
		if err != nil {
			t.Fatalf("PermuteF1600: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("vector %d mismatch", v.id)
		}
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

type keccakVector struct {
	id  int
	in  string
	out string
}

func loadKeccakVectors35233JSON() ([]keccakVector, error) {
	content, err := os.ReadFile("../testdata/ts35233_keccak.json")
	if err != nil {
		return nil, err
	}
	var parsed struct {
		KeccakF1600 []struct {
			ID  int    `json:"id"`
			In  string `json:"in"`
			Out string `json:"out"`
		} `json:"keccak_f1600"`
	}
	if err := json.Unmarshal(content, &parsed); err != nil {
		return nil, err
	}
	out := make([]keccakVector, 0, len(parsed.KeccakF1600))
	for _, v := range parsed.KeccakF1600 {
		out = append(out, keccakVector{id: v.ID, in: v.In, out: v.Out})
	}
	return out, nil
}
