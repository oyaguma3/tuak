// Package tuak implements the 3GPP TUAK algorithm set.
package tuak

import (
	"fmt"

	"tuak/keccak"
)

const (
	inSize         = 200
	offsetTOP      = 0
	offsetInst     = 32
	offsetAlgo     = 33
	offsetRAND     = 40
	offsetAMF      = 56
	offsetSQN      = 58
	offsetK        = 64
	paddingByte96  = 96
	paddingByte135 = 135
)

var algoName = []byte("TUAK1.0")

// TUAK holds inputs and options for TUAK functions.
type TUAK struct {
	k    []byte
	top  []byte
	topc []byte
	rand []byte
	sqn  []byte
	amf  []byte
	opts Options
}

// New creates a TUAK context using TOP (TOPc will be derived as needed).
func New(k, top, rand, sqn, amf []byte, opts ...Option) (*TUAK, error) {
	o := applyOptions(k, opts)
	return &TUAK{
		k:    k,
		top:  top,
		rand: rand,
		sqn:  sqn,
		amf:  amf,
		opts: o,
	}, nil
}

// NewWithTOPc creates a TUAK context using a precomputed TOPc.
func NewWithTOPc(k, topc, rand, sqn, amf []byte, opts ...Option) (*TUAK, error) {
	o := applyOptions(k, opts)
	return &TUAK{
		k:    k,
		topc: topc,
		rand: rand,
		sqn:  sqn,
		amf:  amf,
		opts: o,
	}, nil
}

// ComputeTOPc derives TOPc from K and TOP.
func ComputeTOPc(k, top []byte, opts ...Option) ([]byte, error) {
	o := applyOptions(k, opts)
	kLenBits, err := resolveKLength(k, o)
	if err != nil {
		return nil, err
	}
	if err := requireLen("top", top, 32); err != nil {
		return nil, err
	}

	inst, err := instanceForTOPc(kLenBits)
	if err != nil {
		return nil, err
	}

	state := newState()
	pushData(state, offsetTOP, top)
	state[offsetInst] = inst
	pushData(state, offsetAlgo, algoName)
	pushData(state, offsetK, k)

	callDebug(o.DebugHook, "topc.in", state)
	out, err := permute(state, o.KeccakIterations, o.DebugHook, "topc")
	if err != nil {
		return nil, err
	}
	return pullData(out, offsetTOP, 32), nil
}

// F1 computes MAC-A.
func (t *TUAK) F1() ([]byte, error) {
	topc, err := t.ensureTOPc()
	if err != nil {
		return nil, err
	}
	if err := validateF1Inputs(t); err != nil {
		return nil, err
	}

	inst, err := instanceForF1(t.opts.MACLength, len(t.k)*8, false)
	if err != nil {
		return nil, err
	}

	state := newState()
	pushData(state, offsetTOP, topc)
	state[offsetInst] = inst
	pushData(state, offsetAlgo, algoName)
	pushData(state, offsetRAND, t.rand)
	pushData(state, offsetAMF, t.amf)
	pushData(state, offsetSQN, t.sqn)
	pushData(state, offsetK, t.k)

	callDebug(t.opts.DebugHook, "f1.in", state)
	out, err := permute(state, t.opts.KeccakIterations, t.opts.DebugHook, "f1")
	if err != nil {
		return nil, err
	}
	return pullData(out, offsetTOP, t.opts.MACLength/8), nil
}

// F1Star computes MAC-S.
func (t *TUAK) F1Star() ([]byte, error) {
	topc, err := t.ensureTOPc()
	if err != nil {
		return nil, err
	}
	if err := validateF1Inputs(t); err != nil {
		return nil, err
	}

	inst, err := instanceForF1(t.opts.MACLength, len(t.k)*8, true)
	if err != nil {
		return nil, err
	}

	state := newState()
	pushData(state, offsetTOP, topc)
	state[offsetInst] = inst
	pushData(state, offsetAlgo, algoName)
	pushData(state, offsetRAND, t.rand)
	pushData(state, offsetAMF, t.amf)
	pushData(state, offsetSQN, t.sqn)
	pushData(state, offsetK, t.k)

	callDebug(t.opts.DebugHook, "f1star.in", state)
	out, err := permute(state, t.opts.KeccakIterations, t.opts.DebugHook, "f1star")
	if err != nil {
		return nil, err
	}
	return pullData(out, offsetTOP, t.opts.MACLength/8), nil
}

// F2345 computes RES, CK, IK and AK.
func (t *TUAK) F2345() (res, ck, ik, ak []byte, err error) {
	topc, err := t.ensureTOPc()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err := validateF2345Inputs(t); err != nil {
		return nil, nil, nil, nil, err
	}

	inst, err := instanceForF2345(t.opts.RESLength, t.opts.CKLength, t.opts.IKLength, len(t.k)*8)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	state := newState()
	pushData(state, offsetTOP, topc)
	state[offsetInst] = inst
	pushData(state, offsetAlgo, algoName)
	pushData(state, offsetRAND, t.rand)
	pushData(state, offsetK, t.k)

	callDebug(t.opts.DebugHook, "f2345.in", state)
	out, err := permute(state, t.opts.KeccakIterations, t.opts.DebugHook, "f2345")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	res = pullData(out, offsetTOP, t.opts.RESLength/8)
	ck = pullData(out, 32, t.opts.CKLength/8)
	ik = pullData(out, 64, t.opts.IKLength/8)
	ak = pullData(out, 96, 6)
	return res, ck, ik, ak, nil
}

// F5Star computes AK*.
func (t *TUAK) F5Star() ([]byte, error) {
	topc, err := t.ensureTOPc()
	if err != nil {
		return nil, err
	}
	if err := validateF5StarInputs(t); err != nil {
		return nil, err
	}

	inst, err := instanceForF5Star(len(t.k) * 8)
	if err != nil {
		return nil, err
	}

	state := newState()
	pushData(state, offsetTOP, topc)
	state[offsetInst] = inst
	pushData(state, offsetAlgo, algoName)
	pushData(state, offsetRAND, t.rand)
	pushData(state, offsetK, t.k)

	callDebug(t.opts.DebugHook, "f5star.in", state)
	out, err := permute(state, t.opts.KeccakIterations, t.opts.DebugHook, "f5star")
	if err != nil {
		return nil, err
	}
	return pullData(out, 96, 6), nil
}

func (t *TUAK) ensureTOPc() ([]byte, error) {
	if t.topc != nil {
		return t.topc, nil
	}
	if t.top == nil {
		return nil, fmt.Errorf("tuak: missing topc and top")
	}
	topc, err := ComputeTOPc(t.k, t.top, WithKLength(t.opts.KLength), WithKeccakIterations(t.opts.KeccakIterations))
	if err != nil {
		return nil, err
	}
	t.topc = topc
	return topc, nil
}

func newState() []byte {
	state := make([]byte, inSize)
	state[paddingByte96] = 0x1F
	state[paddingByte135] = 0x80
	return state
}

func permute(state []byte, iterations int, hook DebugHook, label string) ([]byte, error) {
	if iterations <= 0 {
		iterations = 1
	}
	out := state
	var err error
	for i := 0; i < iterations; i++ {
		out, err = keccak.PermuteF1600(out)
		if err != nil {
			return nil, err
		}
		callDebug(hook, debugLabel(label, i, iterations), out)
	}
	return out, nil
}

func pushData(state []byte, offset int, data []byte) {
	for i := 0; i < len(data); i++ {
		state[offset+i] = data[len(data)-1-i]
	}
}

func pullData(state []byte, offset, length int) []byte {
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		out[i] = state[offset+length-1-i]
	}
	return out
}

func callDebug(h DebugHook, label string, data []byte) {
	if h == nil {
		return
	}
	buf := make([]byte, len(data))
	copy(buf, data)
	h(label, buf)
}

func debugLabel(base string, idx, total int) string {
	if total <= 1 {
		return base + ".out"
	}
	return fmt.Sprintf("%s.out.%d", base, idx+1)
}

func resolveKLength(k []byte, opts Options) (int, error) {
	if len(k) != 16 && len(k) != 32 {
		return 0, fmt.Errorf("tuak: invalid K length %d bytes", len(k))
	}
	if opts.KLength == 0 {
		return len(k) * 8, nil
	}
	if opts.KLength != len(k)*8 {
		return 0, fmt.Errorf("tuak: K length mismatch: opt=%d bits, len=%d bits", opts.KLength, len(k)*8)
	}
	return opts.KLength, nil
}

func requireLen(name string, b []byte, want int) error {
	if len(b) != want {
		return fmt.Errorf("tuak: %s length %d bytes (want %d)", name, len(b), want)
	}
	return nil
}

func validateF1Inputs(t *TUAK) error {
	if t.opts.MACLength == 0 {
		return fmt.Errorf("tuak: MAC length must be set")
	}
	if t.opts.MACLength != 64 && t.opts.MACLength != 128 && t.opts.MACLength != 256 {
		return fmt.Errorf("tuak: invalid MAC length %d bits", t.opts.MACLength)
	}
	if _, err := resolveKLength(t.k, t.opts); err != nil {
		return err
	}
	if err := requireLen("rand", t.rand, 16); err != nil {
		return err
	}
	if err := requireLen("sqn", t.sqn, 6); err != nil {
		return err
	}
	if err := requireLen("amf", t.amf, 2); err != nil {
		return err
	}
	return nil
}

func validateF2345Inputs(t *TUAK) error {
	if t.opts.RESLength == 0 || t.opts.CKLength == 0 || t.opts.IKLength == 0 {
		return fmt.Errorf("tuak: RES/CK/IK lengths must be set")
	}
	switch t.opts.RESLength {
	case 32, 64, 128, 256:
	default:
		return fmt.Errorf("tuak: invalid RES length %d bits", t.opts.RESLength)
	}
	switch t.opts.CKLength {
	case 128, 256:
	default:
		return fmt.Errorf("tuak: invalid CK length %d bits", t.opts.CKLength)
	}
	switch t.opts.IKLength {
	case 128, 256:
	default:
		return fmt.Errorf("tuak: invalid IK length %d bits", t.opts.IKLength)
	}
	if _, err := resolveKLength(t.k, t.opts); err != nil {
		return err
	}
	if err := requireLen("rand", t.rand, 16); err != nil {
		return err
	}
	return nil
}

func validateF5StarInputs(t *TUAK) error {
	if _, err := resolveKLength(t.k, t.opts); err != nil {
		return err
	}
	if err := requireLen("rand", t.rand, 16); err != nil {
		return err
	}
	return nil
}

func instanceForTOPc(kLenBits int) (byte, error) {
	if kLenBits != 128 && kLenBits != 256 {
		return 0, fmt.Errorf("tuak: invalid K length %d bits", kLenBits)
	}
	var inst byte
	if kLenBits == 256 {
		inst = setInstanceBit(inst, 7, true)
	}
	return inst, nil
}

func instanceForF1(macLenBits, kLenBits int, star bool) (byte, error) {
	if kLenBits != 128 && kLenBits != 256 {
		return 0, fmt.Errorf("tuak: invalid K length %d bits", kLenBits)
	}
	var inst byte
	inst = setInstanceBit(inst, 0, star)
	inst = setInstanceBit(inst, 1, false)
	switch macLenBits {
	case 64:
		inst = setInstanceBit(inst, 2, false)
		inst = setInstanceBit(inst, 3, false)
		inst = setInstanceBit(inst, 4, true)
	case 128:
		inst = setInstanceBit(inst, 2, false)
		inst = setInstanceBit(inst, 3, true)
		inst = setInstanceBit(inst, 4, false)
	case 256:
		inst = setInstanceBit(inst, 2, true)
		inst = setInstanceBit(inst, 3, false)
		inst = setInstanceBit(inst, 4, false)
	default:
		return 0, fmt.Errorf("tuak: invalid MAC length %d bits", macLenBits)
	}
	inst = setInstanceBit(inst, 5, false)
	inst = setInstanceBit(inst, 6, false)
	inst = setInstanceBit(inst, 7, kLenBits == 256)
	return inst, nil
}

func instanceForF2345(resLenBits, ckLenBits, ikLenBits, kLenBits int) (byte, error) {
	if kLenBits != 128 && kLenBits != 256 {
		return 0, fmt.Errorf("tuak: invalid K length %d bits", kLenBits)
	}
	var inst byte
	inst = setInstanceBit(inst, 0, false)
	inst = setInstanceBit(inst, 1, true)
	switch resLenBits {
	case 32:
		inst = setInstanceBit(inst, 2, false)
		inst = setInstanceBit(inst, 3, false)
		inst = setInstanceBit(inst, 4, false)
	case 64:
		inst = setInstanceBit(inst, 2, false)
		inst = setInstanceBit(inst, 3, false)
		inst = setInstanceBit(inst, 4, true)
	case 128:
		inst = setInstanceBit(inst, 2, false)
		inst = setInstanceBit(inst, 3, true)
		inst = setInstanceBit(inst, 4, false)
	case 256:
		inst = setInstanceBit(inst, 2, true)
		inst = setInstanceBit(inst, 3, false)
		inst = setInstanceBit(inst, 4, false)
	default:
		return 0, fmt.Errorf("tuak: invalid RES length %d bits", resLenBits)
	}
	switch ckLenBits {
	case 128:
		inst = setInstanceBit(inst, 5, false)
	case 256:
		inst = setInstanceBit(inst, 5, true)
	default:
		return 0, fmt.Errorf("tuak: invalid CK length %d bits", ckLenBits)
	}
	switch ikLenBits {
	case 128:
		inst = setInstanceBit(inst, 6, false)
	case 256:
		inst = setInstanceBit(inst, 6, true)
	default:
		return 0, fmt.Errorf("tuak: invalid IK length %d bits", ikLenBits)
	}
	inst = setInstanceBit(inst, 7, kLenBits == 256)
	return inst, nil
}

func instanceForF5Star(kLenBits int) (byte, error) {
	if kLenBits != 128 && kLenBits != 256 {
		return 0, fmt.Errorf("tuak: invalid K length %d bits", kLenBits)
	}
	var inst byte
	inst = setInstanceBit(inst, 0, true)
	inst = setInstanceBit(inst, 1, true)
	inst = setInstanceBit(inst, 2, false)
	inst = setInstanceBit(inst, 3, false)
	inst = setInstanceBit(inst, 4, false)
	inst = setInstanceBit(inst, 5, false)
	inst = setInstanceBit(inst, 6, false)
	inst = setInstanceBit(inst, 7, kLenBits == 256)
	return inst, nil
}

func setInstanceBit(inst byte, index int, value bool) byte {
	if !value {
		return inst
	}
	return inst | (1 << (7 - index))
}
