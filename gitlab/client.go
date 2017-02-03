package gitlab

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	gogitlab "github.com/xanzy/go-gitlab"
)

const GitLabAPI = "/api/v3/"

// Client is a wrapper for the go-gitlab.Client object that provides
// additional methods and initializes with a URL.
type Client struct {
	*gogitlab.Client
	Token string

	Projects *Projects
	Labels   *Labels
}

// NewClient returns a Client object that can be used to make API calls.
// If instead of token you have username and password, you should use
// NewClientForUser().
func NewClient(uri *url.URL, token string) (*Client, error) {
	c := &Client{
		Client: getClient(token),
		Token:  token,
	}
	if err := c.Client.SetBaseURL(uri.String() + GitLabAPI); err != nil {
		return nil, err
	}

	c.Projects = &Projects{c.Client.Projects, c}
	c.Labels = &Labels{c.Client.Labels, c}

	return c, nil
}

// NewClientForUser is the same as NewClient but uses an user instead
// of a private token to authenticate.
func NewClientForUser(uri *url.URL, user, pass string) (*Client, error) {
	c, err := NewClient(uri, "")
	if err != nil {
		return nil, err
	}
	t, err := c.getTokenForUser(user, pass)
	if err != nil {
		return nil, err
	}
	return NewClient(uri, t)
}

// getTokenForUser returns the token for the given user.
func (c *Client) getTokenForUser(user, pass string) (string, error) {
	sess, _, err := c.Client.Session.GetSession(&gogitlab.GetSessionOptions{
		Login:    &user,
		Password: &pass,
	})
	if err != nil {
		return "", err
	}
	return sess.PrivateToken, nil
}

// getClient returns a gitlab client with a timeout and https check disabled.
// Before using it, you should call SetBaseURL() to set the GitLab url.
func getClient(token string) *gogitlab.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 5 * time.Minute}
	return gogitlab.NewClient(client, token)
}
