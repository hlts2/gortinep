package gorpool_log

// Logger is the interface for log output
type Logger interface {
	Printf(format string, args ...interface{})
	Print(v ...interface{})
}
