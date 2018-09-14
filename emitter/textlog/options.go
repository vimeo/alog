package textlog

const (
	fileFlag = 1 << iota
	shortfileFlag
	timeFlag
	utcFlag
)

// Options holds option values.
type Options struct {
	prefix  string
	datefmt string
	flags   uint
}

// Option sets an option for the emitter.
//
// Options are applied in the order specified.
type Option func(*Options)

// WithDateFormat sets the string format for timestamps using a layout string
// like the time package would take.
func WithDateFormat(layout string) Option {
	return func(o *Options) {
		o.datefmt = layout
		o.flags |= timeFlag
	}
}

// WithPrefix adds a set prefix to all lines.
func WithPrefix(prefix string) Option {
	return func(o *Options) { o.prefix = prefix }
}

// WithFile collects call information on each log line, like the log
// package's Llongfile flag.
//
// The alog.WithCaller() option also needs to be used when creating the Logger
// in order to have the file and line information added to the log entries.
func WithFile() Option {
	return func(o *Options) { o.flags |= fileFlag }
}

// WithShortFile is like WithFile, but only prints the file name
// instead of the entire path.
//
// The alog.WithCaller() option also needs to be used when creating the Logger
// in order to have the file and line information added to the log entries.
func WithShortFile() Option {
	return func(o *Options) { o.flags |= fileFlag | shortfileFlag }
}

// WithUTC sets timestamps to UTC.
func WithUTC() Option {
	return func(o *Options) { o.flags |= utcFlag }
}
