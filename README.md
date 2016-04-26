# gitlab-cli [![Build Status](https://travis-ci.org/clns/gitlab-cli.svg?branch=master)](https://travis-ci.org/clns/gitlab-cli)

Cli commands for performing actions against GitLab repositories.

- [Usage](#usage)
- [Install](#install)
- [Development](#development)

## Usage

See help for all available commands (`gitlab-cli -h`).

### Labels

##### Copy global labels into a target repository

> GitLab Limitation: Currently there's no way to [access global labels through the API](https://twitter.com/gitlab/status/724619173477924865), so this tool provides a workaround for copying them into a repository. Note that you should configure the global labels manually in GitLab.

```sh
gitlab-cli label copy -u https://gitlab.com/<USER>/<REPO> -t <TOKEN>
```

> Tip: To avoid specifying `-u` and `-t` every time you refer to a repository, you can save the details of it into the config file with `gitlab-cli config repo save -r <NAME> -u <URL> -t <TOKEN>`, then refer to it simply as `-r <NAME>`.

##### Copy labels from one repository to another

```sh
gitlab-cli label copy -r <NAME> <USER>/<SOURCE_REPO>
```

> Tip: The above command copies labels between repositories on the same GitLab instance. To copy from a different GitLab instance, first save the source repo in the config as explained above and specify its name as argument instead of the path.

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

Currently only the label commands are useful. Other commands can be added as needed.

## Install

1. Follow the instructions from the [releases page](https://github.com/clns/gitlab-cli/releases) and run the `curl` command, which the releases page specifies, in your terminal.

    > Note: If you get a "Permission denied" error, your `/usr/local/bin` directory probably isn't writable and you'll need to install Compose as the superuser. Run `sudo -i`, then the commands from the release page, then `exit`.

2. Test the installation.

    ```sh
    gitlab-cli version
    ```

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