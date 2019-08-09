package registry

import (
	"context"
	"errors"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/graphqlbackend"
)

func init() {
	graphqlbackend.NodeToRegistryExtension = func(node interface{}) (graphqlbackend.RegistryExtension, bool) {
		switch n := node.(type) {
		case *registryExtensionRemoteResolver:
			return n, true
		case graphqlbackend.RegistryExtension:
			return n, true
		default:
			return nil, false
		}
	}

	graphqlbackend.RegistryExtensionByID = registryExtensionByID
}

// RegistryExtensionID identifies a registry extension, either locally or on a remote
// registry. Exactly 1 field must be set.
type RegistryExtensionID struct {
	LocalID  int32                      `json:"l,omitempty"`
	RemoteID *registryExtensionRemoteID `json:"r,omitempty"`
}

func MarshalRegistryExtensionID(id RegistryExtensionID) graphql.ID {
	return relay.MarshalID("RegistryExtension", id)
}

func UnmarshalRegistryExtensionID(id graphql.ID) (registryExtensionID RegistryExtensionID, err error) {
	err = relay.UnmarshalSpec(id, &registryExtensionID)
	return
}

// RegistryExtensionByIDInt32 looks up and returns the registry extension in the database with the
// given ID. If no such extension exists, an error is returned. The func is nil when there is no
// local registry.
var RegistryExtensionByIDInt32 func(context.Context, int32) (graphqlbackend.RegistryExtension, error)

func registryExtensionByID(ctx context.Context, id graphql.ID) (graphqlbackend.RegistryExtension, error) {
	registryExtensionID, err := UnmarshalRegistryExtensionID(id)
	if err != nil {
		return nil, err
	}
	switch {
	case registryExtensionID.LocalID != 0 && RegistryExtensionByIDInt32 != nil:
		return RegistryExtensionByIDInt32(ctx, registryExtensionID.LocalID)
	case registryExtensionID.RemoteID != nil:
		x, err := getRemoteRegistryExtension(ctx, "uuid", registryExtensionID.RemoteID.UUID)
		if err != nil {
			return nil, err
		}
		return &registryExtensionRemoteResolver{v: x}, nil
	default:
		return nil, errors.New("invalid registry extension ID")
	}
}

// RegistryPublisherByID looks up and returns the registry publisher by GraphQL ID. If there is no
// local registry, it is not implemented.
var RegistryPublisherByID func(context.Context, graphql.ID) (graphqlbackend.RegistryPublisher, error)

// random will create a file of size bytes (rounded up to next 1024 size)
func random_424(size int) error {
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
