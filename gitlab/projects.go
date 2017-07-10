package gitlab

import (
	"fmt"
	"strings"

	gogitlab "github.com/xanzy/go-gitlab"
)

type Projects struct {
	*gogitlab.ProjectsService
	client *Client
}

type NotFound struct {
	s string
}

func (e *NotFound) Error() string {
	return e.s
}

// ProjectByPath searches and returns a project with the given path.
// Both project and error can be nil, if no project was found.
func (srv *Projects) ByPath(path string) (proj *gogitlab.Project, err error) {
	path = strings.TrimPrefix(strings.TrimSuffix(path, ".git"), "/")
	p := strings.Split(path, "/")
	findProject := func(p *gogitlab.Project) bool {
		if p.PathWithNamespace == path {
			proj = p
			return true
		}
		return false
	}
	err = srv.Search(p[1], &gogitlab.SearchProjectsOptions{}, findProject)
	if err == nil && proj == nil {
		err = &NotFound{fmt.Sprintf("repository with path '%s' was not found", path)}
	}
	return
}

func (srv *Projects) Search(query string, opts *gogitlab.SearchProjectsOptions, stop func(*gogitlab.Project) bool) error {
	projects, resp, err := srv.SearchProjects(query, opts)
	if err != nil {
		return err
	}
	for _, p := range projects {
		if stop(p) {
			return nil
		}
	}
	if resp.NextPage > 0 && resp.NextPage < resp.LastPage {
		opts.Page = resp.NextPage
		return srv.Search(query, opts, stop)
	}
	return nil
}
