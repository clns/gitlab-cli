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

	Projects *Projects
	Labels   *Labels
}

// NewClient returns a Client object that can be used to make API calls.
// If instead of token you have username and password, you should call
// SetTokenForUser() on the returned object, providing an empty token string.
func NewClient(uri *url.URL, token string) (*Client, error) {
	c := &Client{
		Client: getClient(token),
	}
	if err := c.Client.SetBaseURL(uri.String() + GitLabAPI); err != nil {
		return nil, err
	}

	c.Projects = &Projects{c.Client.Projects, c}
	c.Labels = &Labels{c.Client.Labels, c}

	return c, nil
}

// SetTokenForUser sets the token for the given user on the GitLab client.
func (c *Client) SetTokenForUser(user, pass string) error {
	sess, _, err := c.Client.Session.GetSession(&gogitlab.GetSessionOptions{
		Login:    user,
		Password: pass,
	})
	if err != nil {
		return err
	}
	c.SetToken(sess.PrivateToken)
	return nil
}

// SetToken sets the given token on the GitLab client.
func (c *Client) SetToken(token string) {
	c.Client = getClient(token)
}

// getClient returns a gitlab client with a timeout and https check disabled.
// Before using it, you should call SetBaseURL() to set the GitLab url.
func getClient(token string) *gogitlab.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}
	return gogitlab.NewClient(client, token)
}
