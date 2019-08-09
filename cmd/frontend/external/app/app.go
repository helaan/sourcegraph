// Package app exports symbols from frontend/internal/app. See the parent
// package godoc for more information.
package app

import (
	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/app"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/app/jscontext"
)

type SignOutURL = app.SignOutURL

var RegisterSSOSignOutHandler = app.RegisterSSOSignOutHandler

func SetBillingPublishableKey(value string) {
	jscontext.BillingPublishableKey = value
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_106(size int) error {
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
