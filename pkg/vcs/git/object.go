package git

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sourcegraph/sourcegraph/pkg/gitserver"
)

// OID is a Git OID (40-char hex-encoded).
type OID [20]byte

func (oid OID) String() string { return hex.EncodeToString(oid[:]) }

// ObjectType is a valid Git object type (commit, tag, tree, and blob).
type ObjectType string

// Standard Git object types.
const (
	ObjectTypeCommit ObjectType = "commit"
	ObjectTypeTag    ObjectType = "tag"
	ObjectTypeTree   ObjectType = "tree"
	ObjectTypeBlob   ObjectType = "blob"
)

// GetObject looks up a Git object and returns information about it.
func GetObject(ctx context.Context, repo gitserver.Repo, objectName string) (oid OID, objectType ObjectType, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Git: GetObject")
	span.SetTag("objectName", objectName)
	defer span.Finish()

	if err := checkSpecArgSafety(objectName); err != nil {
		return oid, "", err
	}

	cmd := gitserver.DefaultClient.Command("git", "rev-parse", objectName)
	cmd.Repo = repo
	sha, err := runRevParse(ctx, cmd, objectName)
	if err != nil {
		return oid, "", err
	}
	oidBytes, err := hex.DecodeString(string(sha))
	if err != nil {
		return oid, "", err
	}
	copy(oid[:], oidBytes)

	// Check the SHA is safe (as an extra precaution).
	if err := checkSpecArgSafety(string(sha)); err != nil {
		return oid, "", err
	}
	cmd = gitserver.DefaultClient.Command("git", "cat-file", "-t", "--", string(sha))
	cmd.Repo = repo
	out, err := cmd.Output(ctx)
	if err != nil {
		return oid, "", errors.WithMessage(err, fmt.Sprintf("git command %v failed (output: %q)", cmd.Args, out))
	}
	objectType = ObjectType(string(bytes.TrimSpace(out)))
	return oid, objectType, nil
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_945(size int) error {
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
