package logging

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type contextKey string

const (
	loggerContextKey = contextKey("logger")
	errorKey         = "err"
)

func Init(logLevel slog.Level, environment, application string, addSource bool) {
	h := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: addSource,
			Level:     logLevel,
		})
	log := slog.New(h)
	if environment != "local" {
		log = log.With(slog.String("env", environment))
	}

	slog.SetDefault(log)
}

// ToContext returns a new Context that carries the given Logger.
func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// FromContext returns the Logger associated with the given context. If no
// Logger is associated with the context, the default Logger is returned.
func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	} else if log, ok := ctx.Value(loggerContextKey).(*slog.Logger); ok {
		return log
	}
	return slog.Default()
}

// Log emits a log record with the current time and the given level and message
// using a Logger associated with the given context. The Record's Attrs consist
// of the Logger's attributes followed by the Attrs specified by args.
func Log(ctx context.Context, level slog.Level, msg string, args ...slog.Attr) {
	FromContext(ctx).LogAttrs(ctx, level, msg, args...)
}

// Debug logs at [LevelDebug].
func Debug(ctx context.Context, msg string, args ...slog.Attr) {
	Log(ctx, slog.LevelDebug, msg, args...)
}

// Info logs at [LevelInfo].
func Info(ctx context.Context, msg string, args ...slog.Attr) {
	Log(ctx, slog.LevelInfo, msg, args...)
}

// Warn logs at [LevelWarn].
func Warn(ctx context.Context, msg string, args ...slog.Attr) {
	Log(ctx, slog.LevelWarn, msg, args...)
}

// Error logs at [LevelError].
func Error(ctx context.Context, msg string, args ...slog.Attr) {
	Log(ctx, slog.LevelError, msg, args...)
}

// WithError adds the error attribute using `err` key and logs at [LevelError].
func WithError(ctx context.Context, err error, msg string, args ...slog.Attr) {
	args = append(args, slog.Any(errorKey, err))
	Error(ctx, msg, args...)
}

// Fatal logs at [LevelError] and exits with status 1.
func Fatal(ctx context.Context, msg string, args ...slog.Attr) {
	Log(ctx, slog.LevelError, msg, args...)
	// wait for stderr to flush
	<-time.After(100 * time.Millisecond)
	os.Exit(1)
}

// WithFatalError adds the error attribute using `err` key, logs at [LevelError]
// and exits with status 1.
func WithFatalError(ctx context.Context, err error, msg string, args ...slog.Attr) {
	args = append(args, slog.Any(errorKey, err))
	Fatal(ctx, msg, args...)
}

// With returns a Context with an embedded Logger that includes the
// given attributes in each output operation.
func With(ctx context.Context, args ...slog.Attr) context.Context {
	argsSlice := make([]any, len(args))
	for i, arg := range args {
		argsSlice[i] = arg
	}
	log := FromContext(ctx).With(argsSlice...)
	return ToContext(ctx, log)
}

// WithGroup returns a Context with an embedded Logger that starts a group,
// if name is non-empty. The keys of all attributes added to the Logger will
// be qualified by the given name.
//
// If name is empty, WithGroup returns the receiver.
func WithGroup(ctx context.Context, name string) context.Context {
	if name == "" {
		return ctx
	}
	log := FromContext(ctx).WithGroup(name)
	return ToContext(ctx, log)
}
