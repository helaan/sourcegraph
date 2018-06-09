package conf

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sourcegraph/sourcegraph/schema"
)

// ContributeValidator adds the site configuration validator function to the validation process. It
// is called to validate site configuration. Any strings it returns are shown as validation
// problems.
//
// It may only be called at init time.
func ContributeValidator(f func(schema.SiteConfiguration) (problems []string)) {
	contributedValidators = append(contributedValidators, f)
}

var contributedValidators []func(schema.SiteConfiguration) []string

func validateCustomRaw(normalizedInput []byte) (problems []string, err error) {
	var cfg schema.SiteConfiguration
	if err := json.Unmarshal(normalizedInput, &cfg); err != nil {
		return nil, err
	}
	return validateCustom(cfg), nil
}

// validateCustom validates the site config using custom validation steps that are not
// able to be expressed in the JSON Schema.
func validateCustom(cfg schema.SiteConfiguration) (problems []string) {

	invalid := func(msg string) {
		problems = append(problems, msg)
	}

	// Auth provider config validation is contributed by the
	// github.com/sourcegraph/sourcegraph/cmd/frontend/internal/auth/... packages (using
	// ContributeValidator).

	{
		hasSMTP := cfg.EmailSmtp != nil
		hasSMTPAuth := cfg.EmailSmtp != nil && cfg.EmailSmtp.Authentication != "none"
		if hasSMTP && cfg.EmailAddress == "" {
			invalid(`should set email.address because email.smtp is set`)
		}
		if hasSMTPAuth && (cfg.EmailSmtp.Username == "" && cfg.EmailSmtp.Password == "") {
			invalid(`must set email.smtp username and password for email.smtp authentication`)
		}
	}

	{
		for _, phabCfg := range cfg.Phabricator {
			if len(phabCfg.Repos) == 0 && phabCfg.Token == "" {
				invalid(`each phabricator instance must have either "token" or "repos" set`)
			}
		}
	}

	for _, bbsCfg := range cfg.BitbucketServer {
		if bbsCfg.Token != "" && (bbsCfg.Username != "" || bbsCfg.Password != "") {
			invalid("for Bitbucket Server, specify either a token or a username/password, not both")
		} else if bbsCfg.Token == "" && bbsCfg.Username == "" && bbsCfg.Password == "" {
			invalid("for Bitbucket Server, you must specify either a token or a username/password to authenticate")
		}
	}

	for _, ghCfg := range cfg.Github {
		if ghCfg.PreemptivelyClone {
			invalid(`github[].preemptivelyClone is deprecated; use initialRepositoryEnablement instead`)
		}
	}

	for _, c := range cfg.Gitlab {
		if strings.Contains(c.Url, "example.com") {
			invalid(fmt.Sprintf(`invalid GitLab URL detected: %s (did you forget to remove "example.com"?)`, c.Url))
		}
	}

	if cfg.DisableExampleSearches {
		if cfg.DisableBuiltInSearches {
			invalid(`disableExampleSearches was renamed to disableBuiltInSearches, which is set; you can safely remove the disableExampleSearches setting`)
		} else {
			invalid(`disableExampleSearches was renamed to disableBuiltInSearches; use that instead`)
		}
	}

	if cfg.GithubEnterpriseAccessToken != "" || cfg.GithubEnterpriseCert != "" || cfg.GithubEnterpriseURL != "" {
		invalid(`githubEnterprise{AccessToken,Cert,URL} are deprecated; instead use {..., "github": [{"url": "https://github-enterprise.example.com", "token": "..."}], ...}`)
	}

	if cfg.GithubPersonalAccessToken != "" {
		invalid(`githubPersonalAccessToken is deprecated; instead use {..., "github": [{"url": "https://github.com", "token": "..."}], ....}`)
	}

	if cfg.GitOriginMap != "" {
		invalid(`gitOriginMap is deprecated; instead use code host configuration such as "github", "gitlab", "repos.list", documented at https://about.sourcegraph.com/docs/config/repositories`)
	}

	for _, f := range contributedValidators {
		problems = append(problems, f(cfg)...)
	}

	return problems
}

// TestValidator is an exported helper function for other packages to test their contributed
// validators (registered with ContributeValidator). It should only be called by tests.
func TestValidator(t interface {
	Errorf(format string, args ...interface{})
	Helper()
}, c schema.SiteConfiguration, f func(*schema.SiteConfiguration) []string, wantProblems []string) {
	t.Helper()
	problems := f(&c)
	wantSet := make(map[string]struct{}, len(wantProblems))
	for _, p := range wantProblems {
		wantSet[p] = struct{}{}
	}
	for _, p := range problems {
		var found bool
		for ps := range wantSet {
			if strings.Contains(p, ps) {
				delete(wantSet, ps)
				found = true
				break
			}
		}
		if !found {
			t.Errorf("got unexpected error %q", p)
		}
	}
	if len(wantSet) > 0 {
		t.Errorf("got no matches for expected error substrings %q", wantSet)
	}
}
