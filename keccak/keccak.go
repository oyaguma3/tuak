// Package keccak implements Keccak-f[1600] used by TUAK.
package keccak

import (
	"encoding/binary"
	"errors"
	"math/bits"
)

var errInvalidLength = errors.New("keccak: input must be 200 bytes")

// PermuteF1600 applies the Keccak-f[1600] permutation to a 200-byte state.
func PermuteF1600(in []byte) ([]byte, error) {
	if len(in) != 200 {
		return nil, errInvalidLength
	}

	var a [25]uint64
	for i := 0; i < 25; i++ {
		a[i] = binary.LittleEndian.Uint64(in[i*8 : i*8+8])
	}

	keccakF1600(&a)

	out := make([]byte, 200)
	for i := 0; i < 25; i++ {
		binary.LittleEndian.PutUint64(out[i*8:i*8+8], a[i])
	}
	return out, nil
}

func keccakF1600(a *[25]uint64) {
	var bc [5]uint64
	for round := 0; round < 24; round++ {
		for i := 0; i < 5; i++ {
			bc[i] = a[i] ^ a[i+5] ^ a[i+10] ^ a[i+15] ^ a[i+20]
		}
		for i := 0; i < 5; i++ {
			t := bc[(i+4)%5] ^ bits.RotateLeft64(bc[(i+1)%5], 1)
			a[i] ^= t
			a[i+5] ^= t
			a[i+10] ^= t
			a[i+15] ^= t
			a[i+20] ^= t
		}

		t := a[1]
		for i := 0; i < 24; i++ {
			j := keccakPiLane[i]
			bc[0] = a[j]
			a[j] = bits.RotateLeft64(t, int(keccakRotc[i]))
			t = bc[0]
		}

		for j := 0; j < 25; j += 5 {
			for i := 0; i < 5; i++ {
				bc[i] = a[j+i]
			}
			for i := 0; i < 5; i++ {
				a[j+i] ^= (^bc[(i+1)%5]) & bc[(i+2)%5]
			}
		}

		a[0] ^= keccakRoundConst[round]
	}
}

var keccakRotc = [24]uint64{
	1, 3, 6, 10, 15, 21, 28, 36, 45, 55, 2, 14,
	27, 41, 56, 8, 25, 43, 62, 18, 39, 61, 20, 44,
}

var keccakPiLane = [24]int{
	10, 7, 11, 17, 18, 3, 5, 16, 8, 21, 24, 4,
	15, 23, 19, 13, 12, 2, 20, 14, 22, 9, 6, 1,
}

var keccakRoundConst = [24]uint64{
	0x0000000000000001,
	0x0000000000008082,
	0x800000000000808A,
	0x8000000080008000,
	0x000000000000808B,
	0x0000000080000001,
	0x8000000080008081,
	0x8000000000008009,
	0x000000000000008A,
	0x0000000000000088,
	0x0000000080008009,
	0x000000008000000A,
	0x000000008000808B,
	0x800000000000008B,
	0x8000000000008089,
	0x8000000000008003,
	0x8000000000008002,
	0x8000000000000080,
	0x000000000000800A,
	0x800000008000000A,
	0x8000000080008081,
	0x8000000000008080,
	0x0000000080000001,
	0x8000000080008008,
}
