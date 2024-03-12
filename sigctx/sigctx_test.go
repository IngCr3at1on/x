package sigctx_test

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ingcr3at1on/x/sigctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWith(t *testing.T) {
	t.Run("waits for function to complete", func(tt *testing.T) {
		var complete atomic.Bool
		sigctx.With(func(ctx context.Context) error {
			complete.Store(true)
			return nil
		})
		assert.True(tt, complete.Load())
	})

	t.Run("waits for signal", func(tt *testing.T) {
		var wg sync.WaitGroup
		go func() {
			wg.Done()
			sigctx.With(func(ctx context.Context) error {
				<-ctx.Done()
				return nil
			})
		}()

		p, err := os.FindProcess(os.Getpid())
		require.NoError(tt, err)

		err = p.Signal(os.Interrupt)
		require.NoError(tt, err)

		wg.Wait()
	})
}
