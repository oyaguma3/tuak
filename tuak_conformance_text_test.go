package tuak

import (
	"bytes"
	"fmt"
	"testing"

	"tuak/testvectors"
)

func TestConformanceVectorsFrom35233Text(t *testing.T) {
	vectors, err := loadConformanceVectors35233Text()
	if err != nil {
		t.Fatalf("load 35233 vectors: %v", err)
	}
	if len(vectors) == 0 {
		t.Fatalf("no vectors parsed from 35233 text")
	}

	for _, v := range vectors {
		v := v
		t.Run(fmt.Sprintf("set%d", v.ID), func(t *testing.T) {
			k := decodeHex(t, v.K)
			top := decodeHex(t, v.Top)
			topc := decodeHex(t, v.Topc)
			rand := decodeHex(t, v.Rand)
			sqn := decodeHex(t, v.SQN)
			amf := decodeHex(t, v.AMF)

			wantF1 := decodeHex(t, v.F1)
			wantF1Star := decodeHex(t, v.F1Star)
			wantF2 := decodeHex(t, v.F2)
			wantF3 := decodeHex(t, v.F3)
			wantF4 := decodeHex(t, v.F4)
			wantF5 := decodeHex(t, v.F5)
			wantF5Star := decodeHex(t, v.F5Star)

			gotTopc, err := ComputeTOPc(k, top, optionsFromVector(v)...)
			if err != nil {
				t.Fatalf("ComputeTOPc: %v", err)
			}
			if !bytes.Equal(gotTopc, topc) {
				t.Fatalf("TOPc mismatch")
			}

			tuak, err := NewWithTOPc(k, topc, rand, sqn, amf, optionsFromVector(v)...)
			if err != nil {
				t.Fatalf("NewWithTOPc: %v", err)
			}

			gotF1, err := tuak.F1()
			if err != nil {
				t.Fatalf("F1: %v", err)
			}
			if !bytes.Equal(gotF1, wantF1) {
				t.Fatalf("f1 mismatch")
			}

			gotF1Star, err := tuak.F1Star()
			if err != nil {
				t.Fatalf("F1*: %v", err)
			}
			if !bytes.Equal(gotF1Star, wantF1Star) {
				t.Fatalf("f1* mismatch")
			}

			res, ck, ik, ak, err := tuak.F2345()
			if err != nil {
				t.Fatalf("F2345: %v", err)
			}
			if !bytes.Equal(res, wantF2) {
				t.Fatalf("f2 mismatch")
			}
			if !bytes.Equal(ck, wantF3) {
				t.Fatalf("f3 mismatch")
			}
			if !bytes.Equal(ik, wantF4) {
				t.Fatalf("f4 mismatch")
			}
			if !bytes.Equal(ak, wantF5) {
				t.Fatalf("f5 mismatch")
			}

			gotF5Star, err := tuak.F5Star()
			if err != nil {
				t.Fatalf("F5*: %v", err)
			}
			if !bytes.Equal(gotF5Star, wantF5Star) {
				t.Fatalf("f5* mismatch")
			}
		})
	}
}

func loadConformanceVectors35233Text() ([]testvectors.TUAKVector, error) {
	data, err := testvectors.LoadTUAKVectorsText()
	if err != nil {
		return nil, err
	}
	return data.Tests, nil
}
