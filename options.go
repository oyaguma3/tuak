package tuak

// Options configures TUAK parameters.
type Options struct {
	KLength          int
	MACLength        int
	RESLength        int
	CKLength         int
	IKLength         int
	KeccakIterations int
	DebugHook        DebugHook
}

// Option configures TUAK parameters.
type Option func(*Options)

// WithKLength sets the K length in bits.
func WithKLength(bits int) Option {
	return func(o *Options) {
		o.KLength = bits
	}
}

// WithMACLength sets the MAC length in bits.
func WithMACLength(bits int) Option {
	return func(o *Options) {
		o.MACLength = bits
	}
}

// WithRESLength sets the RES length in bits.
func WithRESLength(bits int) Option {
	return func(o *Options) {
		o.RESLength = bits
	}
}

// WithCKLength sets the CK length in bits.
func WithCKLength(bits int) Option {
	return func(o *Options) {
		o.CKLength = bits
	}
}

// WithIKLength sets the IK length in bits.
func WithIKLength(bits int) Option {
	return func(o *Options) {
		o.IKLength = bits
	}
}

// WithKeccakIterations sets the number of Keccak permutations.
func WithKeccakIterations(n int) Option {
	return func(o *Options) {
		o.KeccakIterations = n
	}
}

// WithDebugHook enables debug callbacks with intermediate buffers.
func WithDebugHook(h DebugHook) Option {
	return func(o *Options) {
		o.DebugHook = h
	}
}

func applyOptions(k []byte, opts []Option) Options {
	out := Options{
		KeccakIterations: 1,
	}
	for _, opt := range opts {
		opt(&out)
	}
	if out.KLength == 0 && len(k) > 0 {
		out.KLength = len(k) * 8
	}
	return out
}
