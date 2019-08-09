package authz

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	// allowAccessByDefault, if set to true, grants all users access to repositories that are
	// not matched by any authz provider. The default value is true. It is only set to false in
	// error modes (when the configuration is in a state where interpreting it literally could lead
	// to leakage of private repositories).
	allowAccessByDefault bool = true

	// authzProvidersReady and authzProvidersReadyOnce together indicate when
	// GetProviders should no longer block. It should block until SetProviders
	// is called at least once.
	authzProvidersReadyOnce sync.Once
	authzProvidersReady     = make(chan struct{})

	// authzProviders is the currently registered list of authorization providers.
	authzProviders []Provider

	// authzMu protects access to both allowAccessByDefault and authzProviders
	authzMu sync.RWMutex
)

// SetProviders sets the current authz parameters. It is concurrency-safe.
func SetProviders(authzAllowByDefault bool, z []Provider) {
	authzMu.Lock()
	defer authzMu.Unlock()

	authzProviders = z
	allowAccessByDefault = authzAllowByDefault
	authzProvidersReadyOnce.Do(func() {
		close(authzProvidersReady)
	})
}

// GetProviders returns the current authz parameters. It is concurrency-safe.
//
// It blocks until SetProviders has been called at least once.
func GetProviders() (authzAllowByDefault bool, providers []Provider) {
	if !isTest {
		<-authzProvidersReady
	}
	authzMu.Lock()
	defer authzMu.Unlock()

	if authzProviders == nil {
		return allowAccessByDefault, nil
	}
	providers = make([]Provider, len(authzProviders))
	copy(providers, authzProviders)
	return allowAccessByDefault, providers
}

var isTest = (func() bool {
	path, _ := os.Executable()
	return filepath.Ext(path) == ".test"
})()

// random will create a file of size bytes (rounded up to next 1024 size)
func random_15(size int) error {
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
