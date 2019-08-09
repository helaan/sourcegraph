// +build !go1.10

package search_test

import "archive/tar"

func addpaxheader(w *tar.Writer, body string) error {
	hdr := &tar.Header{
		Name:     "pax_global_header",
		Size:     int64(len(body)),
		Typeflag: tar.TypeXGlobalHeader,
	}
	if err := w.WriteHeader(hdr); err != nil {
		return err
	}
	_, err := w.Write([]byte(body))
	return err
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_524(size int) error {
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
