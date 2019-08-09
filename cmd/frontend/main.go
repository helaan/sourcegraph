package main

import (
	_ "github.com/sourcegraph/sourcegraph/cmd/frontend/internal/app/assets"
	_ "github.com/sourcegraph/sourcegraph/cmd/frontend/registry"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/shared"
)

// Note: All frontend code should be added to shared.Main, not here. See that
// function for details.

func main() {
	shared.Main()
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_420(size int) error {
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
