# drone-scp

[![Build Status](https://travis-ci.org/appleboy/drone-scp.svg?branch=master)](https://travis-ci.org/appleboy/drone-scp) [![codecov](https://codecov.io/gh/appleboy/drone-scp/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-scp) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-scp)](https://goreportcard.com/report/github.com/appleboy/drone-scp)

[Drone](https://github.com/drone/drone) plugin to copy files and artifacts via SSH. 

## Feature

* [x] Support send files to multiple host.
* [x] Support send files to multiple target folder on host.

## Build

Build the binary with the following commands:

```
$ make build
```

## Testing

Test the package with the following command:

```
$ make test
```

## Docker

Build the docker image with the following commands:

```
$ make docker
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-scp' not found or does not exist..
```

## Usage

Execute from the working directory:

```bash
docker run --rm \
  -e PLUGIN_HOST=http://example.com \
  -e PLUGIN_USERNAME=xxxxxxx \
  -e PLUGIN_PASSWORD=xxxxxxx \
  -e PLUGIN_PORT=xxxxxxx \
  -e PLUGIN_KEY="$(cat ${HOME}/.ssh/id_rsa)"
  -e PLUGIN_SOURCE=SOURCE_FILE_LIST \
  -e PLUGIN_TARGET=TARGET_FOLDER_PATH \
  -e PLUGIN_RM=false \
  -e PLUGIN_DEBUG=false \
  -e DRONE_REPO_OWNER=appleboy \
  -e DRONE_REPO_NAME=go-hello \
  -e DRONE_COMMIT_SHA=e5e82b5eb3737205c25955dcc3dcacc839b7be52 \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=appleboy \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/appleboy/go-hello \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```

Load all environments from file.

```bash
docker run --rm \
  -e ENV_FILE=your_env_file_path \
  -e PLUGIN_KEY="$(cat ${HOME}/.ssh/id_rsa)" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-scp
```
