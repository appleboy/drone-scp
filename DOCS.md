---
date: 2017-01-06T00:00:00+00:00
title: SCP
author: appleboy
tags: [ publish, ssh, scp ]
repo: appleboy/drone-scp
logo: term.svg
image: appleboy/drone-scp
---

The SCP plugin copy files and artifacts to target host machine via SSH. The below pipeline configuration demonstrates simple usage:

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
: plain text of user private key

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
