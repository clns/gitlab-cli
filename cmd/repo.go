package cmd

import (
	"net/url"

	"fmt"
	"strings"

	"log"

	"github.com/clns/gitlab-cli/gitlab"
	"github.com/howeyc/gopass"
	"github.com/spf13/viper"
	gogitlab "github.com/xanzy/go-gitlab"
)

// Repo represents a cli repository.
type Repo struct {
	Client  *gitlab.Client
	Project *gogitlab.Project
	Name    string
	Url_    string `mapstructure:"url"`
	URL     *url.URL
	Token   string `mapstructure:"token"`
}

type repoMap struct {
	URL   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

func LoadFromConfig(namepath string) (*Repo, error) {
	r := LoadFromConfigNoInit(namepath)
	if err := r.initialize(); err != nil {
		return nil, err
	}
	return r, nil
}

func LoadFromConfigNoInit(namepath string) *Repo {
	key := "repos." + namepath
	r := &Repo{
		Url_:  viper.GetString("_url"),
		Token: viper.GetString("_token"),
	}
	if viper.IsSet(key) {
		viper.UnmarshalKey(key, r)
		r.Name = namepath
	} else if namepath != "" {
		if r.URL, _ = url.Parse(r.Url_); r.URL != nil {
			r.URL.Path = namepath
		}
	}
	return r
}

func (r *Repo) String() string {
	return fmt.Sprintf(`%s
  url: %s
  token: %s`, r.Name, r.Url_, r.Token)
}

func (r *Repo) SaveToConfig() error {
	if r.Name == "" {
		return fmt.Errorf("cannot save to config without a name")
	}
	repos := make(map[string]*repoMap)
	for name, _ := range viper.GetStringMap("repos") {
		var rep *repoMap
		if err := viper.UnmarshalKey("repos."+name, &rep); err != nil {
			log.Println(err)
		}
		repos[name] = rep
	}
	repos[r.Name] = &repoMap{
		URL:   r.URL.String(),
		Token: r.Token,
	}
	viper.Set("repos", repos)

	return nil
}

func (r *Repo) initialize() error {
	var err error
	r.URL, err = url.Parse(r.Url_)
	if err != nil {
		return fmt.Errorf("invalid repo url: %v", err)
	}
	if r.URL.String() == "" {
		return fmt.Errorf("empty repo url")
	}
	if r.URL.Path == "" || strings.Index(r.URL.Path, "/") == -1 {
		return fmt.Errorf("invalid or no repo path specified")
	}
	if r.Client, err = r.client(); err != nil {
		return fmt.Errorf("failed to get GitLab client for repo '%s': %v", r.URL, err)
	}
	r.Token = r.Client.Token
	if r.Project, err = r.project(); err != nil {
		return fmt.Errorf("failed to get GitLab project '%s': %v", r.URL, err)
	}
	return nil
}

func (r *Repo) client() (*gitlab.Client, error) {
	u := *r.URL
	u.Path = ""
	if r.Token == "" && user != "" {
		if password == "" {
			fmt.Print("Password: ")
			pwd, _ := gopass.GetPasswdMasked()
			password = string(pwd)
		}
		return gitlab.NewClientForUser(&u, user, password)
	}
	return gitlab.NewClient(&u, r.Token)
}

func (r *Repo) project() (*gogitlab.Project, error) {
	proj, err := r.Client.Projects.ByPath(r.URL.Path)
	if err != nil {
		return nil, err
	}
	return proj, nil
}
