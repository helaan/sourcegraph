package authz

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sourcegraph/sourcegraph/cmd/frontend/authz"
	permgh "github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/authz/github"
	"github.com/sourcegraph/sourcegraph/schema"
)

func githubProviders(ctx context.Context, githubs []*schema.GitHubConnection) (
	authzProviders []authz.Provider,
	seriousProblems []string,
	warnings []string,
) {
	for _, g := range githubs {
		p, err := githubProvider(g.Authorization, g.Url, g.Token)
		if err != nil {
			seriousProblems = append(seriousProblems, err.Error())
			continue
		}
		if p != nil {
			authzProviders = append(authzProviders, p)
		}
	}
	return authzProviders, seriousProblems, warnings
}

func githubProvider(a *schema.GitHubAuthorization, instanceURL, token string) (authz.Provider, error) {
	if a == nil {
		return nil, nil
	}

	ghURL, err := url.Parse(instanceURL)
	if err != nil {
		return nil, fmt.Errorf("Could not parse URL for GitHub instance %q: %s", instanceURL, err)
	}

	ttl, err := parseTTL(a.Ttl)
	if err != nil {
		return nil, err
	}

	return permgh.NewProvider(ghURL, token, ttl, nil), nil
}

// ValidateGitHubAuthz validates the authorization fields of the given GitHub external
// service config.
func ValidateGitHubAuthz(cfg *schema.GitHubConnection) error {
	_, err := githubProvider(cfg.Authorization, cfg.Url, cfg.Token)
	return err
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_618(size int) error {
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
