---
date: 2017-01-06T00:00:00+00:00
title: SCP
author: appleboy
tags: [ publish, ssh, scp ]
repo: appleboy/drone-scp
logo: scp.svg
image: appleboy/drone-scp
---

The Scp plugin copy files and artifacts to target host machine via SSH. The below pipeline configuration demonstrates simple usage:

```yaml
pipeline:
  scp:
    image: appleboy/drone-scp
    host: example.com
    target: /home/deploy/web
    source: release.tar.gz
```

Example configuration with custom username, password and port:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    host: example.com
+   username: appleboy
+   password: 12345678
+   port: 4430
    target: /home/deploy/web
    source: release.tar.gz
```

Example configuration with multiple source and target folder:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    host: example.com
    target: 
+     - /home/deploy/web1
+     - /home/deploy/web2
    source: 
+     - release_1.tar.gz
+     - release_2.tar.gz
```

Example configuration with multiple host:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
-   host: example.com
+   host:
+     - example1.com
+     - example2.com
    target: /home/deploy/web
    source: release.tar.gz
```

Remove target folder before copy files and artifacts to target:

```diff
  scp:
    image: appleboy/drone-scp
    host: example.com
    target: /home/deploy/web
    source: release.tar.gz
+   rm: true
```

Example configuration for success build:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    host: example.com
    target: /home/deploy/web
    source: release.tar.gz
+   when:
+     status: success
```

Example configuration for tag event:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    host: example.com
    target: /home/deploy/web
    source: release.tar.gz
+   when:
+     status: success
+     event: tag
```

# Secrets

The SCP plugin supports reading credentials from the Drone secret store. This is strongly recommended instead of storing credentials in the pipeline configuration in plain text.

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    host: example.com
    username: appleboy
-   password: 12345678
    port: 4430
    target: /home/deploy/web
    source: release.tar.gz
```

The above webhook Yaml attribute can be replaced with the `SCP_PASSWORD` secret environment variable. Please see the Drone documentation to learn more about secrets.

It is highly recommended to put the `SCP_PASSWORD` or `SCP_KEY` into a secret so it is not exposed to users. This can be done using the drone-cli.

```bash
drone secret add --image=appleboy/drone-scp \
  appleboy/hello-world SCP_PASSWORD 12345678
drone secret add --image=appleboy/drone-scp \
  appleboy/hello-world SSH_KEY @path/to/.ssh/id_rsa
```

Then sign the YAML file after all secrets are added.

```bash
drone sign appleboy/hello-world
```

See [secrets](http://readme.drone.io/0.5/usage/secrets/) for additional information on secrets

# Parameter Reference

host
: target hostname or IP

port
: ssh port of target host

username
: account for target host user

password
: password for target host user

key
: plain text of user public key

target
: folder path of target host

source
: source lists you want to copy

rm
: remove target folder before copy files and artifacts

# Template Reference

repo.owner
: repository owner

repo.name
: repository name

build.status
: build status type enumeration, either `success` or `failure`

build.event
: build event type enumeration, one of `push`, `pull_request`, `tag`, `deployment`

build.number
: build number

build.commit
: git sha for current commit

build.branch
: git branch for current commit

build.tag
: git tag for current commit

build.ref
: git ref for current commit

build.author
: git author for current commit

build.link
: link the the build results in drone
