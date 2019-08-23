package sigctx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type (
	// Logger is based on the standard logger proposals discussed in detail, linked below
	// https://docs.google.com/document/d/1oTjtY49y8iSxmM9YBaz2NrZIlaXtQsq3nQMd-E0HtwM/edit#
	Logger interface {
		// Log is a flexible log function described in the standard logger proposals.
		Log(...interface{}) error
	}

	noOpLogger struct{}
)

var _logger Logger = &noOpLogger{}

func (*noOpLogger) Log(_ ...interface{}) error {
	return nil
}

// SetLogger sets a Logger interface for sigctx to use.
func SetLogger(logger Logger) {
	if logger != nil {
		_logger = logger
	}
}

// FromContext wraps the provided context.Context with a context.CancelFunc
// that has cancel signal monitoring attached to it and returns the
// context.Context.
func FromContext(ctx context.Context) context.Context {
	ctx, _ = WithCancel(ctx)
	return ctx
}

// WithCancel wraps the provided context.Context with a context.CancelFunc
// that has cancel signal monitoring attached to it.
func WithCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		sig := <-sc
		_logger.Log(sig.String() + ` signal received`)
		cancel()
	}()

	return ctx, cancel
}

// StartWith starts a blocking function with a context.Context and
// context.CancelFunc from WithCancel.
func StartWith(f func(ctx context.Context) error) error {
	return StartWithContext(context.Background(), f)
}

// StartWithContext starts a blocking function, it creates a new
// context.Context and context.CancelFunc from WithCancel using
// the provided context.Context as a base.
func StartWithContext(ctx context.Context, f func(ctx context.Context) error) error {
	ctx, cancel := WithCancel(ctx)
	errCh := make(chan error, 1)
	defer close(errCh)
	go func() {
		defer cancel()
		err := f(ctx)
		if err != nil {
			errCh <- err
			return
		}
		_logger.Log(`f finished, cancelling context`)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			return err
		}
	}
}
