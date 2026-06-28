package utils

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const defaulBucketSize = 20
const refreshBucketSize = 5
const defaultTickInterval = 500

type Limiter struct {
	bucketSize        int
	currentBucketSize int
	interval          time.Duration

	mu     sync.Mutex
	cancel context.CancelFunc
}

func NewLimiter(ctx context.Context, bucketSize int, interval time.Duration) *Limiter {
	ctx, cancel := context.WithCancel(ctx)

	if bucketSize == 0 {
		bucketSize = defaulBucketSize
	}

	if interval == 0{
		interval = defaultTickInterval
	}

	l := &Limiter{
		bucketSize:        bucketSize,
		currentBucketSize: bucketSize,
		interval:          interval,
	}

	go l.refreshBucket(ctx)

	l.cancel = cancel

	return l
}

func (limiter *Limiter) Stop() {
	limiter.cancel()
}

func (limiter *Limiter) refreshBucket(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(limiter.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("release all")
				return
			case <-ticker.C:
				// add mutex, because CurrentBucketSize access
				// is shared among decrementBucket and GetBucketSize
				limiter.mu.Lock()
				if limiter.currentBucketSize < limiter.bucketSize {
					currentTickingSize := limiter.currentBucketSize + refreshBucketSize

					if currentTickingSize > limiter.bucketSize {
						limiter.currentBucketSize = limiter.bucketSize
					} else {
						limiter.currentBucketSize = currentTickingSize
					}
				}
				limiter.mu.Unlock()
			}
		}
	}()
}

func (limiter *Limiter) decrementBucket() {
	// add mutex, because currentBucketSize access
	// is shared among GetBucketSize and refreshBucket
	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	limiter.currentBucketSize--
}

func (limiter *Limiter) GetBucketSize() int {
	// add mutex, because currentBucketSize access
	// is shared among decrementBucket and refreshBucket
	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	return limiter.currentBucketSize
}