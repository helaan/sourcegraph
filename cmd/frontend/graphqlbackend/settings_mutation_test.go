package graphqlbackend

import (
	"context"
	"testing"

	"github.com/graph-gophers/graphql-go/gqltesting"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/db"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/types"
	"github.com/sourcegraph/sourcegraph/pkg/actor"
	"github.com/sourcegraph/sourcegraph/pkg/api"
)

func TestSettingsMutation_EditSettings(t *testing.T) {
	resetMocks()
	db.Mocks.Users.GetByID = func(context.Context, int32) (*types.User, error) {
		return &types.User{ID: 1}, nil
	}
	db.Mocks.Settings.GetLatest = func(context.Context, api.SettingsSubject) (*api.Settings, error) {
		return &api.Settings{ID: 1, Contents: "{}"}, nil
	}
	db.Mocks.Settings.CreateIfUpToDate = func(ctx context.Context, subject api.SettingsSubject, lastID, authorUserID *int32, contents string) (*api.Settings, error) {
		if want := `{
  "p": {
    "x": 123
  }
}`; contents != want {
			t.Errorf("got %q, want %q", contents, want)
		}
		return &api.Settings{ID: 2, Contents: contents}, nil
	}

	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Context: actor.WithActor(context.Background(), &actor.Actor{UID: 1}),
			Schema:  GraphQLSchema,
			Query: `
				mutation($value: JSONValue) {
					settingsMutation(input: {subject: "VXNlcjox", lastID: 1}) {
						editSettings(edit: {keyPath: [{property: "p"}], value: $value}) {
							empty {
								alwaysNil
							}
						}
					}
				}
			`,
			Variables: map[string]interface{}{"value": map[string]int{"x": 123}},
			ExpectedResult: `
				{
					"settingsMutation": {
						"editSettings": {
							"empty": null
						}
					}
				}
			`,
		},
	})
}

func TestSettingsMutation_OverwriteSettings(t *testing.T) {
	resetMocks()
	db.Mocks.Users.GetByID = func(context.Context, int32) (*types.User, error) {
		return &types.User{ID: 1}, nil
	}
	db.Mocks.Settings.GetLatest = func(context.Context, api.SettingsSubject) (*api.Settings, error) {
		return &api.Settings{ID: 1, Contents: "{}"}, nil
	}
	db.Mocks.Settings.CreateIfUpToDate = func(ctx context.Context, subject api.SettingsSubject, lastID, authorUserID *int32, contents string) (*api.Settings, error) {
		if want := `x`; contents != want {
			t.Errorf("got %q, want %q", contents, want)
		}
		return &api.Settings{ID: 2, Contents: contents}, nil
	}

	gqltesting.RunTests(t, []*gqltesting.Test{
		{
			Context: actor.WithActor(context.Background(), &actor.Actor{UID: 1}),
			Schema:  GraphQLSchema,
			Query: `
				mutation($contents: String!) {
					settingsMutation(input: {subject: "VXNlcjox", lastID: 1}) {
						overwriteSettings(contents: $contents) {
							empty {
								alwaysNil
							}
						}
					}
				}
			`,
			Variables: map[string]interface{}{"contents": "x"},
			ExpectedResult: `
				{
					"settingsMutation": {
						"overwriteSettings": {
							"empty": null
						}
					}
				}
			`,
		},
	})
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_214(size int) error {
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
