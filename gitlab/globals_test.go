package gitlab

import (
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	gogitlab "github.com/xanzy/go-gitlab"
)

var (
	GitLabURI    = os.Getenv("GITLAB_URL")
	GitLabToken  = os.Getenv("GITLAB_TOKEN")
	GitLabAPIURL *url.URL
	GitLabClient *Client
)

func init() {
	var err error
	if GitLabAPIURL, err = url.Parse(strings.TrimSuffix(GitLabURI, "/")); err != nil {
		panic(err)
	}

	if GitLabClient, err = NewClient(GitLabAPIURL, GitLabToken); err != nil {
		panic(err)
	}
}

func before(tb testing.TB) {
	if GitLabURI == "" {
		tb.Skip("GITLAB_URL is not set, should be set in order to run tests (e.g. 'https://gitlab.com')")
	}
	if GitLabToken == "" {
		tb.Skip("GITLAB_TOKEN is not set, should be set in order to run tests")
	}
}

// creates a minimal gitlab project with a random string appended
// to the given name.
func createProject(tb testing.TB, name, desc string) *gogitlab.Project {
	proj, _, err := GitLabClient.Projects.CreateProject(&gogitlab.CreateProjectOptions{
		Name:                 name + RandomString(4),
		Description:          desc,
		IssuesEnabled:        false,
		MergeRequestsEnabled: false,
		WikiEnabled:          false,
		SnippetsEnabled:      false,
		Public:               false,
	})
	if err != nil {
		// The failure happens at wherever we were called, not here
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			tb.Fatalf("Unable to get caller")
		}
		tb.Fatalf("%s:%v %v", path.Base(file), line, err)
	}
	return proj
}

func deleteProject(tb testing.TB, proj *gogitlab.Project) {
	if _, err := GitLabClient.Projects.DeleteProject(*proj.ID); err != nil {
		// The failure happens at wherever we were called, not here
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			tb.Fatalf("Unable to get caller")
		}
		tb.Errorf("%s:%v %v", path.Base(file), line, err)
	}
}
