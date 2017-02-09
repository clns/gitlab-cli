# gitlab-cli [![Build Status](https://travis-ci.org/clns/gitlab-cli.svg?branch=master)](https://travis-ci.org/clns/gitlab-cli)

CLI commands for performing actions against GitLab repositories.

- [Installation](#installation)
- [Usage](#usage)
  - [Labels](#labels)
    - [Copy global labels](#copy-global-labels-into-a-repository)
    - [Copy labels from repoA to repoB](#copy-labels-from-repoa-to-repob)
    - [Update labels](#update-labels-that-match-a-regex)
    - [Delete labels](#delete-labels-that-match-a-regex)
  - [Specifying a repository](#specifying-a-repository)
  - [The config file](#the-config-file)
- [Development](#development)

## Installation

Follow the instructions from the [releases page](https://github.com/clns/gitlab-cli/releases).

## Usage

For all available commands see the command's help: `gitlab-cli -h`. The most common commands are documented below.

### Labels

#### Copy global labels into a repository

GitLab Limitation: Currently there's no way to [access global labels through the API](https://twitter.com/gitlab/status/724619173477924865), so this tool provides a workaround to copy them.

```sh
gitlab-cli label copy -U https://gitlab.com/<USER>/<REPO> -t <TOKEN>
```

> Tip: To avoid specifying `-U` and `-t` every time you refer to a repository, you can use the config file to save the details of it. See [Specifying a repository](#specifying-a-repository).

#### Copy labels from repoA to repoB

```sh
gitlab-cli label copy --from <repoA> -r <repoB>
```

> Tip: For repositories on the same installation, you can specify the `--from` repo as `group/repo`, as a convenience, in which case the repository is considered on the same GitLab instance as the target repo.

#### Update labels that match a regex

```sh
gitlab-cli label update -r <NAME> --match <REGEX> --name <NAME> --color <COLOR> --description <DESC>
```

> Note: `<REGEX>` is a Go regex string as in <https://golang.org/pkg/regexp/syntax> and `<NAME>` is a replacement string as in <https://golang.org/pkg/regexp/#Regexp.FindAllString>.

#### Delete labels that match a regex

```sh
gitlab-cli label delete -r <NAME> --match <REGEX>
```

### TODO

Other commands can be added as needed. Feel free to open pull requests or issues.

### Specifying a repository

There are 2 ways to specify a repository:

1. By using the `--url (-U)` and `--token (-t)` flags (or `--user (-u)` and `--password (-p)` instead of token) with each command. This is the easiest to get started but requires a lot of typing.
2. By saving the repository details in the config file and referring to it by its saved name using `--repo (-r)` (e.g. `-r myrepo`)

Example:

Instead of this:

```sh
gitlab-cli label copy -U https://git.my-site.com/my_group/my_repo -t ghs93hska
```

you can first save the repo in the config file and refer to it by name on all subsequent commands:

```sh
gitlab-cli config repo save -r myrepo -U https://git.my-site.com/my_group/my_repo -t ghs93hska
gitlab-cli label copy -r myrepo
```

#### Using user and password instead of token

You can specify your GitLab login (user or email) - `--user (-u)` - and password - `--password (-p)` - instead of the token in any command, if this is easier for you. Example:

```sh
gitlab-cli config repo save -r myrepo -U https://git.my-site.com/my_group/my_repo -u my_user -p my_pass
```

### The config file

The default location of the config file is `$HOME/.gitlab-cli.yaml` and it is useful for saving repositories and then refer to them by their names. A sample config file looks like this:

```yaml
repos:
  myrepo1:
    url: https://git.mysite.com/group/repo1
    token: Nahs93hdl3shjf
  myrepo2:
    url: https://git.mysite.com/group/repo2
    token: Nahs93hdl3shjf
  myother:
    url: https://git.myothersite.com/group/repo1
    token: OA23spfwuSalos
```

But there's no need to manually edit this file. Instead use the config commands to modify it (see `gitlab-cli config -h`). Some useful config commands are:

- `gitlab-cli config cat` - print the entire config file contents
- `gitlab-cli config repo ls` - list all saved repositories
- `gitlab-cli config repo save ...` - save a repository
- `gitlab-cli config repo show -r <repo>` - show the details of a saved repository

## Development

You'll need a [Go dev environment](https://golang.org/doc/install).

```sh
git clone https://github.com/xanzy/go-gitlab
git submodule --init update
```

### Build

```sh
go run build/build.go
```

This will build all the executables into the [build/](build) directory.

### Test

You need to provide a GitLab URL and private token to be able to create temporary repositories for the tests.

```sh
GITLAB_URL="<URL>" GITLAB_TOKEN="<TOKEN>" go test -v ./gitlab
```

You can spin up a GitLab instance using [Docker](https://www.docker.com/):

```sh
docker pull gitlab/gitlab-ce
docker run -d --name gitlab -p 8055:80 gitlab/gitlab-ce
sleep 60 # allow enough time for GitLab to start
docker exec -ti gitlab bash
su gitlab-psql
/opt/gitlab/embedded/bin/psql --port 5432 -h /var/opt/gitlab/postgresql -d gitlabhq_production -c " \
          INSERT INTO labels (title, color, template, description, description_html) VALUES ('feature', '#000000', true, 'represents a feature', 'represents a <b>feature</b>'); \
          INSERT INTO labels (title, color, template, description, description_html) VALUES ('bug', '#ff0000', true, 'represents a bug', 'represents a <b>bug</b>'); \
          UPDATE users SET authentication_token='secret' WHERE username='root';"

# Note: you may need to change GITLAB_URL to point to your docker container.
# 'http://docker' is for Docker beta for Windows. 
GITLAB_URL="http://localhost:8055" GITLAB_TOKEN="secret" go test -v ./gitlab
```