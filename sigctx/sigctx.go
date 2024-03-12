package sigctx

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Context calls signal.NotifyContext with the following signals:
// syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt and os.Kill
func Context(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, syscall.SIGABRT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
}

// With calls Context and Do in succession.
func With(f func(ctx context.Context) error) error {
	ctx, cancel := Context(context.Background())
	defer cancel()
	return Do(ctx, cancel, f)
}

// WithContext calls Context and Do in succession using the given context as the parent context.
func WithContext(ctx context.Context, f func(ctx context.Context) error) error {
	ctx, cancel := Context(ctx)
	defer cancel()
	return Do(ctx, cancel, f)
}

// Do calls f in a new goroutine and waits for it to complete or the context to be cancelled.
// It's expected that the provided context have signals attached to it using Context or signal.NotifyContext.
func Do(ctx context.Context, cancel context.CancelFunc, f func(ctx context.Context) error) error {
	errCh := make(chan error, 1)
	defer close(errCh)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		err := f(ctx)
		if err != nil {
			errCh <- err
			return
		}
	}()

	select {
	case <-ctx.Done():
		break
	case err := <-errCh:
		return err
	}

	wg.Wait()
	return nil
}
