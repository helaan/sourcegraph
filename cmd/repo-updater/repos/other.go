package repos

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sourcegraph/sourcegraph/pkg/api"
	"github.com/sourcegraph/sourcegraph/pkg/conf/reposource"
	"github.com/sourcegraph/sourcegraph/pkg/httpcli"
	"github.com/sourcegraph/sourcegraph/pkg/jsonc"
	"github.com/sourcegraph/sourcegraph/schema"
)

// A OtherSource yields repositories from a single Other connection configured
// in Sourcegraph via the external services configuration.
type OtherSource struct {
	svc    *ExternalService
	conn   *schema.OtherExternalServiceConnection
	client httpcli.Doer
}

// NewOtherSource returns a new OtherSource from the given external service.
func NewOtherSource(svc *ExternalService, cf *httpcli.Factory) (*OtherSource, error) {
	var c schema.OtherExternalServiceConnection
	if err := jsonc.Unmarshal(svc.Config, &c); err != nil {
		return nil, errors.Wrapf(err, "external service id=%d config error", svc.ID)
	}

	cli, err := cf.Doer()
	if err != nil {
		return nil, err
	}

	return &OtherSource{svc: svc, conn: &c, client: cli}, nil
}

// ListRepos returns all Other repositories accessible to all connections configured
// in Sourcegraph via the external services configuration.
func (s OtherSource) ListRepos(ctx context.Context, results chan SourceResult) {
	if s.conn.ExperimentalFakehub {
		repos, err := s.fakehub(ctx)
		if err != nil {
			results <- SourceResult{Source: s, Err: err}
		}
		for _, r := range repos {
			results <- SourceResult{Source: s, Repo: r}
		}
		return
	}

	urls, err := s.cloneURLs()
	if err != nil {
		results <- SourceResult{Source: s, Err: err}
		return
	}

	urn := s.svc.URN()
	for _, u := range urls {
		r, err := s.otherRepoFromCloneURL(urn, u)
		if err != nil {
			results <- SourceResult{Source: s, Err: err}
			return
		}
		results <- SourceResult{Source: s, Repo: r}
	}
}

// ExternalServices returns a singleton slice containing the external service.
func (s OtherSource) ExternalServices() ExternalServices {
	return ExternalServices{s.svc}
}

func (s OtherSource) cloneURLs() ([]*url.URL, error) {
	if len(s.conn.Repos) == 0 {
		return nil, nil
	}

	var base *url.URL
	if s.conn.Url != "" {
		var err error
		if base, err = url.Parse(s.conn.Url); err != nil {
			return nil, err
		}
	}

	cloneURLs := make([]*url.URL, 0, len(s.conn.Repos))
	for _, repo := range s.conn.Repos {
		cloneURL, err := otherRepoCloneURL(base, repo)
		if err != nil {
			return nil, err
		}
		cloneURLs = append(cloneURLs, cloneURL)
	}

	return cloneURLs, nil
}

func otherRepoCloneURL(base *url.URL, repo string) (*url.URL, error) {
	if base == nil {
		return url.Parse(repo)
	}
	return base.Parse(repo)
}

func (s OtherSource) otherRepoFromCloneURL(urn string, u *url.URL) (*Repo, error) {
	repoURL := u.String()
	repoSource := reposource.Other{OtherExternalServiceConnection: s.conn}
	repoName, err := repoSource.CloneURLToRepoName(u.String())
	if err != nil {
		return nil, err
	}
	repoURI, err := repoSource.CloneURLToRepoURI(u.String())
	if err != nil {
		return nil, err
	}
	u.Path, u.RawQuery = "", ""
	serviceID := u.String()

	return &Repo{
		Name: string(repoName),
		URI:  repoURI,
		ExternalRepo: api.ExternalRepoSpec{
			ID:          string(repoName),
			ServiceType: "other",
			ServiceID:   serviceID,
		},
		Enabled: true,
		Sources: map[string]*SourceInfo{
			urn: {
				ID:       urn,
				CloneURL: repoURL,
			},
		},
	}, nil
}

func (s OtherSource) fakehub(ctx context.Context) ([]*Repo, error) {
	req, err := http.NewRequest("GET", s.conn.Url+"/v1/list-repos", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	var data struct {
		Items []*Repo
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response from fakehub")
	}

	urn := s.svc.URN()
	for _, r := range data.Items {
		// The only required field is URI
		if r.URI == "" {
			return nil, errors.Errorf("repo without URI returned from fakehub: %+v", r)
		}

		// Fields that fakehub isn't allowed to control
		r.Enabled = true
		r.ExternalRepo = api.ExternalRepoSpec{
			ID:          r.URI,
			ServiceType: "other",
			ServiceID:   s.conn.Url,
		}
		r.Sources = map[string]*SourceInfo{
			urn: {
				ID: urn,
				// TODO we should allow this to be set
				CloneURL: s.conn.Url + r.URI + "/.git",
			},
		}

		// The only required field left is Name
		if r.Name == "" {
			r.Name = r.URI
		}
	}

	return data.Items, nil
}
