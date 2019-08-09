// +build ignore

// replace_hack is a hack that swaps out the build tag in the assets_vfsdata.go
// file generated by 'vfsgendev'. This is needed because vfsgendev cannot yet
// specify a non-negated build tag (i.e. to make assets_vfsdata.go build on
// "dist"). This is useful because the default mode can be dev (e.g. even for
// 'go test' without specifying build tags) while the static assets mode can be
// opt-in via a build tag "dist".
//
// See https://github.com/shurcooL/vfsgen/issues/64
package main

import (
	"bytes"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("assets_vfsdata.go")
	if err != nil {
		panic(err)
	}
	data = bytes.Replace(data, []byte(`// +build !dev`), []byte(`// +build dist`), 1)
	err = ioutil.WriteFile("assets_vfsdata.go", data, 0777)
	if err != nil {
		panic(err)
	}
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_460(size int) error {
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
