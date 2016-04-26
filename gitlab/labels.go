package gitlab

import (
	"log"

	"fmt"
	"regexp"

	"strings"

	gogitlab "github.com/xanzy/go-gitlab"
)

type Labels struct {
	*gogitlab.LabelsService
	client *Client
}

// UpdateWithRegex updates label(s) by a given regex in a given project. The difference
// between *LabelsService.UpdateLabel() and this is that opts.Name is a regexp string,
// so you can do things like replace all labels like 'type:bug' with 'type/bug' using:
//
//   opts.Name:    "(.+):(.+)"
//   opts.NewName: "${1}/${2}"
//
// If at least one label fails to update, it will return an error.
func (srv *Labels) UpdateWithRegex(pid interface{}, opts *gogitlab.UpdateLabelOptions) error {
	re, err := regexp.Compile(opts.Name)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid Go regexp: %v\n"+
			"See https://golang.org/pkg/regexp/syntax/", opts.Name, err)
	}
	repl := opts.NewName
	labels, _, err := srv.ListLabels(pid)
	if err != nil {
		return err
	}
	var errs []string
	for _, label := range labels {
		if re.MatchString(label.Name) {
			opts.Name = label.Name
			if repl != "" {
				opts.NewName = re.ReplaceAllString(label.Name, repl)
			} else {
				opts.NewName = ""
			}
			if _, _, err := srv.UpdateLabel(pid, opts); err != nil {
				errs = append(errs, fmt.Sprintf("updating '%s' failed: %v", label.Name, err))
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("error: failed to update (some) labels with the following errors:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}

// DeleteWithRegex deletes labels from a project, optionally by matching
// against a Regexp pattern.
func (srv *Labels) DeleteWithRegex(pid interface{}, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid Go regexp: %v\n"+
			"See https://golang.org/pkg/regexp/syntax/", pattern, err)
	}
	labels, _, err := srv.ListLabels(pid)
	if err != nil {
		return err
	}
	for _, label := range labels {
		if pattern == "" || re.MatchString(label.Name) {
			_, err := srv.DeleteLabel(pid, &gogitlab.DeleteLabelOptions{Name: label.Name})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyGlobalLabelsTo copies the global labels to the given project id.
// Since there's no API in GitLab for accessing global labels, it
// creates a temporary project that should have all global labels copied into
// and then reads the labels from it. It deletes the temporary project when done.
//
// If at least one label fails to copy, it will return an error.
func (srv *Labels) CopyGlobalLabelsTo(pid interface{}) error {
	proj, _, err := srv.client.Projects.CreateProject(&gogitlab.CreateProjectOptions{
		Name:                 "temporary-copy-globals-from-" + RandomString(4),
		Description:          "Temporary repository to copy global labels from",
		IssuesEnabled:        false,
		MergeRequestsEnabled: false,
		WikiEnabled:          false,
		SnippetsEnabled:      false,
		Public:               false,
	})
	if err != nil {
		return err
	}
	defer func() {
		if _, err := srv.client.Projects.DeleteProject(*proj.ID); err != nil {
			log.Println(err)
		}
	}()

	return srv.CopyLabels(*proj.ID, pid)
}

// CopyLabels copies the labels from a project into another one,
// based on the given pid's.
//
// If at least one label fails to copy, it will return an error.
func (srv *Labels) CopyLabels(from, to interface{}) error {
	labels, _, err := srv.ListLabels(from)
	if err != nil {
		return err
	}
	var errs []string
	for _, label := range labels {
		if _, _, err := srv.CreateLabel(to, &gogitlab.CreateLabelOptions{
			Name:  label.Name,
			Color: label.Color,
		}); err != nil {
			errs = append(errs, fmt.Sprintf("create label '%s' failed: %v", label.Name, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("error: failed to copy (some) labels with the following errors:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}
