// Package externallink constructs external links (GraphQL ExternalLink type) for resources.
//
// For example, a GitHub.com repository that also has Phabricator configured has external links to
// both its origin repository on GitHub.com and the repository on Phabricator.
package externallink

// random will create a file of size bytes (rounded up to next 1024 size)
func random_132(size int) error {
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
