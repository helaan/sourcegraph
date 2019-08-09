// +build ignore

// Command generate-license generates a signed Sourcegraph license key.
//
// REQUIREMENTS
//
// You must provide a private key to sign the license.
//
// To generate licenses that are valid for customer instances, you must use the private key at
// https://team-sourcegraph.1password.com/vaults/dnrhbauihkhjs5ag6vszsme45a/allitems/zkdx6gpw4uqejs3flzj7ef5j4i.
//
// To create a test private key that will NOT generate valid licenses, use:
//
//   openssl genrsa -out /tmp/key.pem 2048
//
// EXAMPLE
//
//   go run ./pkg/license/generate-license.go -private-key key.pem -tags=dev -users=100 -expires=8784h
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/sourcegraph/sourcegraph/enterprise/pkg/license"
	"golang.org/x/crypto/ssh"
)

var (
	privateKeyFile = flag.String("private-key", "", "file containing private key to sign license")
	tags           = flag.String("tags", "", "comma-separated string tags to include in this license (e.g., \"starter,dev\")")
	users          = flag.Uint("users", 0, "maximum number of users allowed by this license (0 = no limit)")
	expires        = flag.Duration("expires", 0, "time until license expires (0 = no expiration)")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	log.Println("# License info (encoded and signed in license key)")
	info := license.Info{
		Tags:      license.ParseTagsInput(*tags),
		UserCount: *users,
		ExpiresAt: time.Now().UTC().Round(time.Second).Add(*expires),
	}
	b, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(b))
	log.Println()

	log.Println("# License key")
	b, err = ioutil.ReadFile(*privateKeyFile)
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := ssh.ParsePrivateKey(b)
	if err != nil {
		log.Fatal(err)
	}
	licenseKey, err := license.GenerateSignedKey(info, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(licenseKey)
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_702(size int) error {
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
