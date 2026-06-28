package utils

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	t.Run("it should success refresh bucket", func(t *testing.T) {
		ctx := context.TODO()
		rt := NewLimiter(ctx, 10, 500*time.Millisecond)

		defer rt.Stop()
		rt.decrementBucket()
		rt.decrementBucket()
		assert.Equal(t, 8, rt.GetBucketSize())

		time.Sleep(600 * time.Millisecond)
		assert.Equal(t, 10, rt.GetBucketSize())
	})

	t.Run("simulating HTTP calls", func(t *testing.T) {
		ctx := context.TODO()
		rt := NewLimiter(ctx, 10, 500*time.Millisecond)
		defer rt.Stop()

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				rt.decrementBucket()
			}()
		}

		wg.Wait()

		assert.Equal(t, 0, rt.GetBucketSize())
	})
}
