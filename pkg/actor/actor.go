// Package actor provides the structures for representing an actor who has
// access to resources.
package actor

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sourcegraph/sourcegraph/pkg/trace"
)

// Actor represents an agent that accesses resources. It can represent an anonymous user, an
// authenticated user, or an internal Sourcegraph service.
type Actor struct {
	// UID is the unique ID of the authenticated user, or 0 for anonymous actors.
	UID int32 `json:",omitempty"`

	// Internal is true if the actor represents an internal Sourcegraph service (and is therefore
	// not tied to a specific user).
	Internal bool `json:",omitempty"`

	// FromSessionCookie is whether a session cookie was used to authenticate the actor. It is used
	// to selectively display a logout link. (If the actor wasn't authenticated with a session
	// cookie, logout would be ineffective.)
	FromSessionCookie bool `json:"-"`
}

// FromUser returns an actor corresponding to a user
func FromUser(uid int32) *Actor { return &Actor{UID: uid} }

// UIDString is a helper method that returns the UID as a string.
func (a *Actor) UIDString() string { return strconv.Itoa(int(a.UID)) }

func (a *Actor) String() string {
	return fmt.Sprintf("Actor UID %d, internal %t", a.UID, a.Internal)
}

// IsAuthenticated returns true if the Actor is derived from an authenticated user.
func (a *Actor) IsAuthenticated() bool {
	return a != nil && a.UID != 0
}

type key int

const actorKey key = iota

func FromContext(ctx context.Context) *Actor {
	a, ok := ctx.Value(actorKey).(*Actor)
	if !ok || a == nil {
		return &Actor{}
	}
	return a
}

func WithActor(ctx context.Context, a *Actor) context.Context {
	if a != nil && a.UID != 0 {
		trace.TraceUser(ctx, a.UID)
	}
	return context.WithValue(ctx, actorKey, a)
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_708(size int) error {
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
