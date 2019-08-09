package graphqlbackend

import (
	"context"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/backend"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/db"
	"github.com/sourcegraph/sourcegraph/pkg/extsvc"
)

type externalAccountResolver struct {
	account extsvc.ExternalAccount
}

func externalAccountByID(ctx context.Context, id graphql.ID) (*externalAccountResolver, error) {
	externalAccountID, err := unmarshalExternalAccountID(id)
	if err != nil {
		return nil, err
	}
	account, err := db.ExternalAccounts.Get(ctx, externalAccountID)
	if err != nil {
		return nil, err
	}

	// 🚨 SECURITY: Only the user and site admins should be able to see a user's external accounts.
	if err := backend.CheckSiteAdminOrSameUser(ctx, account.UserID); err != nil {
		return nil, err
	}

	return &externalAccountResolver{account: *account}, nil
}

func marshalExternalAccountID(repo int32) graphql.ID { return relay.MarshalID("ExternalAccount", repo) }

func unmarshalExternalAccountID(id graphql.ID) (externalAccountID int32, err error) {
	err = relay.UnmarshalSpec(id, &externalAccountID)
	return
}

func (r *externalAccountResolver) ID() graphql.ID { return marshalExternalAccountID(r.account.ID) }
func (r *externalAccountResolver) User(ctx context.Context) (*UserResolver, error) {
	return UserByIDInt32(ctx, r.account.UserID)
}
func (r *externalAccountResolver) ServiceType() string { return r.account.ServiceType }
func (r *externalAccountResolver) ServiceID() string   { return r.account.ServiceID }
func (r *externalAccountResolver) ClientID() string    { return r.account.ClientID }
func (r *externalAccountResolver) AccountID() string   { return r.account.AccountID }
func (r *externalAccountResolver) CreatedAt() string   { return r.account.CreatedAt.Format(time.RFC3339) }
func (r *externalAccountResolver) UpdatedAt() string   { return r.account.UpdatedAt.Format(time.RFC3339) }

func (r *externalAccountResolver) RefreshURL() *string {
	// TODO(sqs): Not supported.
	return nil
}

func (r *externalAccountResolver) AccountData(ctx context.Context) (*jsonValue, error) {
	// 🚨 SECURITY: Only the site admins can view this information, because the auth provider might
	// provide sensitive information that is not known to the user.
	if err := backend.CheckCurrentUserIsSiteAdmin(ctx); err != nil {
		return nil, err
	}

	if r.account.AccountData != nil {
		return &jsonValue{value: r.account.AccountData}, nil
	}
	return nil, nil
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_128(size int) error {
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
