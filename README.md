# drone-scp

[![GoDoc](https://godoc.org/github.com/appleboy/drone-scp?status.svg)](https://godoc.org/github.com/appleboy/drone-scp) [![Build Status](http://drone.wu-boy.com/api/badges/appleboy/drone-scp/status.svg)](http://drone.wu-boy.com/appleboy/drone-scp) [![codecov](https://codecov.io/gh/appleboy/drone-scp/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-scp) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-scp)](https://goreportcard.com/report/github.com/appleboy/drone-scp) [![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-scp.svg)](https://hub.docker.com/r/appleboy/drone-scp/) [![](https://images.microbadger.com/badges/image/appleboy/drone-scp.svg)](https://microbadger.com/images/appleboy/drone-scp "Get your own image badge on microbadger.com")

Copy files and artifacts via SSH using a binary, docker or [Drone CI](http://readme.drone.io/).

## Feature

* [x] Support routines.
* [x] Support wildcard pattern on source list.
* [x] Support send files to multiple host.
* [x] Support send files to multiple target folder on host.
* [x] Support load ssh key from absolute path or raw body.
* [x] Support SSH ProxyCommand.

```
     +--------+       +----------+      +-----------+
     | Laptop | <-->  | Jumphost | <--> | FooServer |
     +--------+       +----------+      +-----------+

                         OR

     +--------+       +----------+      +-----------+
     | Laptop | <-->  | Firewall | <--> | FooServer |
     +--------+       +----------+      +-----------+
     192.168.1.5       121.1.2.3         10.10.29.68
```

## Build or Download a binary

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/drone-scp/releases). Support the following OS type.

* Windows amd64/386
* Linux amd64/386
* Darwin amd64/386

With `Go` installed

```
$ go get -u -v github.com/appleboy/drone-scp
``` 

or build the binary with the following command:

```
$ make build
```

## Docker

Build the docker image with the following commands:

```
$ make docker
```

Please note incorrectly building the image for the correct x64 linux and with
CGO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-scp' not found or does not exist..
```

## Usage

There are three ways to send notification.

* [usage from binary](#usage-from-binary)
* [usage from docker](#usage-from-docker)
* [usage from drone ci](#usage-from-drone-ci)

<a name="usage-from-binary"></a>
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

```bash
$ ssh-add
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

<a name="usage-from-docker"></a>
### Usage from docker

#### Using public key

```bash
docker run --rm \
  -e SCP_HOST example.com \
  -e SCP_USERNAME xxxxxxx \
  -e SCP_PORT 22 \
  -e SCP_KEY_PATH "${HOME}/.ssh/id_rsa"
  -e SCP_SOURCE SOURCE_FILE_LIST \
  -e SCP_TARGET TARGET_FOLDER_PATH \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```

#### Using password

```diff
docker run --rm \
  -e SCP_HOST example.com \
  -e SCP_USERNAME xxxxxxx \
  -e SCP_PORT 22 \
+ -e SCP_PASSWORD "xxxxxxx"
  -e SCP_SOURCE SOURCE_FILE_LIST \
  -e SCP_TARGET TARGET_FOLDER_PATH \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```

#### Using ssh-agent

Start your local ssh agent:

```bash
eval `ssh-agent -s`
```

Import your local public key `~/.ssh/id_rsa`

```bash
$ ssh-add
```

You don't need to add `SCP_PASSWORD` or `SCP_KEY_PATH ` arguments.

```bash
docker run --rm \
  -e SCP_HOST example.com \
  -e SCP_USERNAME xxxxxxx \
  -e SCP_PORT 22 \
  -e SCP_SOURCE SOURCE_FILE_LIST \
  -e SCP_TARGET TARGET_FOLDER_PATH \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```

#### Send multiple source or target folder and hosts

```bash
docker run --rm \
  -e SCP_HOST example1.com,example2.com \
  -e SCP_USERNAME xxxxxxx \
  -e SCP_PASSWORD xxxxxxx \
  -e SCP_PORT 22 \
  -e SCP_SOURCE SOURCE_FILE_LIST_1,SOURCE_FILE_LIST_2 \
  -e SCP_TARGET TARGET_FOLDER_PATH_1,TARGET_FOLDER_PATH_2 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```

<a name="usage-from-drone-ci"></a>
### Usage from drone ci

Execute from the working directory:

```bash
docker run --rm \
  -e PLUGIN_HOST example.com \
  -e PLUGIN_USERNAME xxxxxxx \
  -e PLUGIN_PASSWORD xxxxxxx \
  -e PLUGIN_PORT xxxxxxx \
  -e PLUGIN_KEY "$(cat ${HOME}/.ssh/id_rsa)"
  -e PLUGIN_SOURCE SOURCE_FILE_LIST \
  -e PLUGIN_TARGET TARGET_FOLDER_PATH \
  -e PLUGIN_RM false \
  -e PLUGIN_DEBUG false \
  -e DRONE_REPO_OWNER appleboy \
  -e DRONE_REPO_NAME go-hello \
  -e DRONE_COMMIT_SHA e5e82b5eb3737205c25955dcc3dcacc839b7be52 \
  -e DRONE_COMMIT_BRANCH master \
  -e DRONE_COMMIT_AUTHOR appleboy \
  -e DRONE_BUILD_NUMBER 1 \
  -e DRONE_BUILD_STATUS success \
  -e DRONE_BUILD_LINK http://github.com/appleboy/go-hello \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```

You can get more [information](http://plugins.drone.io/appleboy/drone-scp/) about how to use scp in drone.

## Testing

Test the package with the following command:

```
$ make test
```
