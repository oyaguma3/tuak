package testvectors

import "testing"

func TestLoadKeccakVectors(t *testing.T) {
	data, err := LoadKeccakVectors()
	if err != nil {
		t.Fatalf("LoadKeccakVectors: %v", err)
	}
	if got := len(data.KeccakF1600); got != 6 {
		t.Fatalf("KeccakF1600 count = %d, want 6", got)
	}
	for _, v := range data.KeccakF1600 {
		if len(v.In)%2 != 0 || len(v.Out)%2 != 0 {
			t.Fatalf("vector %d has odd hex length", v.ID)
		}
		if v.InLen != len(v.In)/2 || v.OutLen != len(v.Out)/2 {
			t.Fatalf("vector %d has length mismatch", v.ID)
		}
	}
}

func TestLoadTUAKVectors(t *testing.T) {
	data, err := LoadTUAKVectors()
	if err != nil {
		t.Fatalf("LoadTUAKVectors: %v", err)
	}
	if got := len(data.Tests); got != 6 {
		t.Fatalf("TUAK test count = %d, want 6", got)
	}
	for _, v := range data.Tests {
		checkLen(t, v.ID, "k", v.K, v.Klength)
		checkLen(t, v.ID, "f1", v.F1, v.MAClength)
		checkLen(t, v.ID, "f1_star", v.F1Star, v.MAClength)
		checkLen(t, v.ID, "f2", v.F2, v.RESLength)
		checkLen(t, v.ID, "f3", v.F3, v.CKlength)
		checkLen(t, v.ID, "f4", v.F4, v.IKlength)
		if len(v.F5) != 12 || len(v.F5Star) != 12 {
			t.Fatalf("vector %d has invalid f5/f5* length", v.ID)
		}
	}
}

func checkLen(t *testing.T, id int, name, hex string, bits int) {
	t.Helper()
	if bits == 0 {
		return
	}
	if len(hex)*4 != bits {
		t.Fatalf("vector %d %s bit length = %d, want %d", id, name, len(hex)*4, bits)
	}
}
