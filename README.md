# gitlab-cli [![Build Status](https://travis-ci.org/clns/gitlab-cli.svg?branch=master)](https://travis-ci.org/clns/gitlab-cli)

CLI commands for performing actions against GitLab repositories. The main reasons for building this tool is to be able to use it without any prerequisites and to deal with global labels, which GitLab API doesn't expose.

- [Usage](#usage)
  - [Labels](#labels)
  - [Specifying a repository](#specifying-a-repository)
  - [The config file](#the-config-file)
- [Install](#install)
  - [Update](#update)
- [Development](#development)

## Usage

See help for all available commands (`gitlab-cli -h`).

### Labels

##### Copy global labels into a target repository

> GitLab Limitation: Currently there's no way to [access global labels through the API](https://twitter.com/gitlab/status/724619173477924865), so this tool provides a workaround for copying them into a repository. Note that you should configure the global labels manually in GitLab.

```sh
gitlab-cli label copy -U https://gitlab.com/<USER>/<REPO> -t <TOKEN>
```

> Tip: To avoid specifying `-U` and `-t` every time you refer to a repository, you can use the config file to save the details of it. See [Specifying a repository](#specifying-a-repository).

##### Copy labels from one repository to another

```sh
gitlab-cli label copy -r <NAME> <GROUP>/<REPO>
```

> Tip: The above command copies labels between repositories on the same GitLab instance. To copy from/to a different GitLab instance, use the config file as explained in [Specifying a repository](#specifying-a-repository).

##### Update label(s) based on a regex match

```sh
gitlab-cli label update -r <NAME> --match <REGEX> --replace <REPLACE> --color <COLOR>
```

> Note: `<REGEX>` is a Go regex string as in <https://golang.org/pkg/regexp/syntax> and `<REPLACE>` is a replacement string as in <https://golang.org/pkg/regexp/#Regexp.FindAllString>.

##### Delete label(s) that match a regex

```sh
gitlab-cli label update -r <NAME> --regex <REGEX>
```

### TODO

Currently only the label commands are useful. Other commands can be added as needed. Feel free to open pull requests or issues.

### Specifying a repository

There are 2 ways to specify a repository:

1. By using the `--url (-U)` and `--token (-t)` flags (or `--user (-u)` and `--password (-p)` instead of token) with each command. This is the easiest to get started but requires a lot of typing.
2. By saving the repository details in the config file and referring to it by its saved name using `--repo (-r)` (e.g. `-r myrepo`)

Example:

```sh
gitlab-cli label copy -U https://git.my-site.com/my_group/my_repo -t ghs93hska
```

is the same as this, but the repo gets saved in the config file and we can refer to it later by its name:

```sh
gitlab-cli config repo save -r myrepo -U https://git.my-site.com/my_group/my_repo -t ghs93hska
gitlab-cli label copy -r myrepo
```

> Note: Some commands like [`label copy`](#copy-labels-from-one-repository-to-another) allow you to specify a repository by its path (e.g. `my_group/my_repo`), in which case the repository is considered on the same GitLab instance as the target repo.

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

## Install

1. Follow the instructions from the [releases page](https://github.com/clns/gitlab-cli/releases) and run the `curl` command, which the releases page specifies, in your terminal.

    > Note: If you get a "Permission denied" error, your `/usr/local/bin` directory probably isn't writable and you'll need to install Compose as the superuser. Run `sudo -i`, then the commands from the release page, then `exit`.

2. Test the installation.

    ```sh
    gitlab-cli version
    ```

### Update

To update, simply run the same command as for install from the [releases page](https://github.com/clns/gitlab-cli/releases). The existing binary will be overwritten by the latest version.

## Development

You'll need a [Go dev environment](https://golang.org/doc/install).

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
docker exec gitlab \
  sudo -u gitlab-psql \
    /opt/gitlab/embedded/bin/psql --port 5432 -h /var/opt/gitlab/postgresql -d gitlabhq_production -c " \
      INSERT INTO labels (title, color, template) VALUES ('feature', '#000000', true); \
      INSERT INTO labels (title, color, template) VALUES ('bug', '#ff0000', true); \
      UPDATE users SET authentication_token='secret' WHERE username='root';"

# Note: you may need to change GITLAB_URL to point to your docker container.
# 'http://docker' is for Docker beta for Windows. 
GITLAB_URL="http://docker:8055" GITLAB_TOKEN="secret" go test -v ./gitlab
```