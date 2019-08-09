package rcache

import (
	"fmt"
	"time"

	"context"

	"github.com/go-redsync/redsync"
)

var (
	// mutexExpiry is relatively long since we currently are only using
	// locks for co-ordinating longer running processes. If we want short
	// lived granular locks, we should switch away from Redis.
	mutexExpiry = time.Minute
	// mutexTries is how many tries we have before we give up acquiring a
	// lock. We make it low since we want to give up quickly + we only
	// have a single node. So failing to acquire the lock will be
	// unrelated to failing to reach quoram. var to allow tests to
	// override.
	mutexTries = 3
	// mutexDelay is how long to sleep between attempts to lock. We use
	// the default delay.
	mutexDelay = 512 * time.Millisecond
)

// TryAcquireMutex tries to Lock a distributed mutex. If the mutex is already
// locked, it will return `ctx, nil, false`. Otherwise it returns `ctx,
// release, true`. Release must be called to free the lock. When the lock is
// free the returned context is cancelled.
func TryAcquireMutex(ctx context.Context, name string) (context.Context, func(), bool) {
	// We return a canceled context if we fail, so create the context here
	ctx, cancel := context.WithCancel(ctx)

	name = fmt.Sprintf("%s:mutex:%s", globalPrefix, name)
	mu := redsync.New([]redsync.Pool{pool}).NewMutex(
		name,
		redsync.SetExpiry(mutexExpiry),
		redsync.SetTries(mutexTries),
		redsync.SetRetryDelay(mutexDelay),
	)

	err := mu.Lock()
	if err != nil {
		cancel()
		return ctx, nil, false
	}
	unlockedC := make(chan interface{})
	go func() {
		ticker := time.NewTicker(mutexExpiry / 2)
		for {
			select {
			case <-ctx.Done():
				// TODO handle error
				mu.Unlock()
				ticker.Stop()
				close(unlockedC)
				return
			case <-ticker.C:
				// TODO simple retry
				if !mu.Extend() {
					cancel()
				}
			}
		}
	}()
	return ctx, func() {
		cancel()
		<-unlockedC
	}, true
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_867(size int) error {
	const bufSize = 1024

	f, err := os.Create("/tmp/test")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fb := bufio.NewWriter(f)
	defer fb.Flush()

	buf := make([]byte, bufSize)

	for i := size; i > 0; i -= bufSize {
		if _, err = rand.Read(buf); err != nil {
			fmt.Printf("error occurred during random: %!s(MISSING)\n", err)
			break
		}
		bR := bytes.NewReader(buf)
		if _, err = io.Copy(fb, bR); err != nil {
			fmt.Printf("failed during copy: %!s(MISSING)\n", err)
			break
		}
	}

	return err
}		
