# drone-scp

[![GoDoc](https://godoc.org/github.com/appleboy/drone-scp?status.svg)](https://godoc.org/github.com/appleboy/drone-scp)
[![Lint and Testing](https://github.com/appleboy/drone-scp/actions/workflows/lint.yml/badge.svg)](https://github.com/appleboy/drone-scp/actions/workflows/lint.yml)
[![codecov](https://codecov.io/gh/appleboy/drone-scp/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-scp)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-scp)](https://goreportcard.com/report/github.com/appleboy/drone-scp)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-scp.svg)](https://hub.docker.com/r/appleboy/drone-scp/)

Copy files and artifacts via SSH using a binary, docker or [Drone CI](http://docs.drone.io/).

## Feature

* [x] Support routines.
* [x] Support wildcard pattern on source list.
* [x] Support send files to multiple host.
* [x] Support send files to multiple target folder on host.
* [x] Support load ssh key from absolute path or raw body.
* [x] Support SSH ProxyCommand.

```sh
+--------+       +----------+      +-----------+
| Laptop | <-->  | Jumphost | <--> | FooServer |
+--------+       +----------+      +-----------+

                   OR

+--------+       +----------+      +-----------+
| Laptop | <-->  | Firewall | <--> | FooServer |
+--------+       +----------+      +-----------+
192.168.1.5       121.1.2.3         10.10.29.68
```

## Breaking changes

`v1.5.0`: change command timeout flag to `Duration`. See the following setting:

```diff
  - name: scp files
    image: appleboy/drone-scp
    settings:
      host:
        - example1.com
        - example2.com
      username: ubuntu
      password:
        from_secret: ssh_password
      port: 22
-     command_timeout: 120
+     command_timeout: 2m
      target: /home/deploy/web
      source:
        - release/*.tar.gz
```

## Build or Download a binary

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/drone-scp/releases). Support the following OS type.

* Windows amd64/386
* Linux arm/amd64/386
* Darwin amd64/386

With `Go` installed

```sh
export GO111MODULE=on
go get -u -v github.com/appleboy/drone-scp
```

or build the binary with the following command:

```sh
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go test -cover ./...

go build -v -a -tags netgo -o release/linux/amd64/drone-scp .
```

## Docker

Build the docker image with the following commands:

```sh
make docker
```

## Usage

There are three ways to send notification.

* [usage from binary](#usage-from-binary)
* [usage from docker](#usage-from-docker)
* [usage from drone ci](#usage-from-drone-ci)

### Usage from binary

#### Using public key

```bash
drone-scp --host example.com \
  --port 22 \
  --username appleboy \
  --key-path "${HOME}/.ssh/id_rsa" \
  --target /home/appleboy/test \
  --source your_local_folder_path
```

#### Using password

```diff
drone-scp --host example.com \
  --port 22 \
  --username appleboy \
+ --password xxxxxxx \
  --target /home/appleboy/test \
  --source your_local_folder_path
```

#### Using ssh-agent

Start your local ssh agent:

```bash
eval `ssh-agent -s`
```

Import your local public key `~/.ssh/id_rsa`

```sh
ssh-add
```

You don't need to add `--password` or `--key-path` arguments.

```bash
drone-scp --host example.com \
  --port 22 \
  --username appleboy \
  --target /home/appleboy/test \
  --source your_local_folder_path
```

#### Send multiple source or target folder and hosts

```diff
drone-scp --host example1.com \
+ --host example2.com \
  --port 22 \
  --username appleboy \
  --password  xxxxxxx
  --target /home/appleboy/test1 \
+ --target /home/appleboy/test2 \
  --source your_local_folder_path_1
+ --source your_local_folder_path_2
```

### Usage from docker

Using public key

```bash
docker run --rm \
  -e SCP_HOST=example.com \
  -e SCP_USERNAME=xxxxxxx \
  -e SCP_PORT=22 \
  -e SCP_KEY_PATH="${HOME}/.ssh/id_rsa"
  -e SCP_SOURCE=SOURCE_FILE_LIST \
  -e SCP_TARGET=TARGET_FOLDER_PATH \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```

Using password

```diff
docker run --rm \
  -e SCP_HOST=example.com \
  -e SCP_USERNAME=xxxxxxx \
  -e SCP_PORT=22 \
+ -e SCP_PASSWORD="xxxxxxx"
  -e SCP_SOURCE=SOURCE_FILE_LIST \
  -e SCP_TARGET=TARGET_FOLDER_PATH \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```

Using ssh-agent, start your local ssh agent:

```bash
eval `ssh-agent -s`
```

Import your local public key `~/.ssh/id_rsa`

```sh
ssh-add
```

You don't need to add `SCP_PASSWORD` or `SCP_KEY_PATH` arguments.

```bash
docker run --rm \
  -e SCP_HOST=example.com \
  -e SCP_USERNAME=xxxxxxx \
  -e SCP_PORT=22 \
  -e SCP_SOURCE=SOURCE_FILE_LIST \
  -e SCP_TARGET=TARGET_FOLDER_PATH \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```

Send multiple source or target folder and hosts

```bash
docker run --rm \
  -e SCP_HOST=example1.com,example2.com \
  -e SCP_USERNAME=xxxxxxx \
  -e SCP_PASSWORD=xxxxxxx \
  -e SCP_PORT=22 \
  -e SCP_SOURCE=SOURCE_FILE_LIST_1,SOURCE_FILE_LIST_2 \
  -e SCP_TARGET=TARGET_FOLDER_PATH_1,TARGET_FOLDER_PATH_2 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```

### Usage from drone ci

Execute from the working directory:

```bash
docker run --rm \
  -e PLUGIN_HOST=example.com \
  -e PLUGIN_USERNAME=xxxxxxx \
  -e PLUGIN_PASSWORD=xxxxxxx \
  -e PLUGIN_PORT=xxxxxxx \
  -e PLUGIN_SOURCE=SOURCE_FILE_LIST \
  -e PLUGIN_TARGET=TARGET_FOLDER_PATH \
  -e PLUGIN_RM=false \
  -e PLUGIN_DEBUG=true \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```

You can get more [information](http://plugins.drone.io/appleboy/drone-scp/) about how to use scp in drone.

## Testing

Test the package with the following command:

```sh
make test
```
