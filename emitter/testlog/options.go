package testlog

// Options holds option values.
type Options struct {
	shortfile bool
}

// Option sets an option for the emitter.
//
// Options are applied in the order specified.
type Option func(*Options)

// WithShortFile only prints the file name instead of the entire path
func WithShortFile() Option {
	return func(o *Options) {
		o.shortfile = true
	}
}
