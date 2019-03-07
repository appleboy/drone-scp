---
date: 2017-01-06T00:00:00+00:00
title: SCP
author: appleboy
tags: [ publish, ssh, scp ]
logo: term.svg
repo: appleboy/drone-scp
image: appleboy/drone-scp
---

The SCP plugin copy files and artifacts to target host machine via SSH. The below pipeline configuration demonstrates simple usage:

```yaml
pipeline:
  scp:
    image: appleboy/drone-scp
    settings:
      host: example.com
      target: /home/deploy/web
      source: release.tar.gz
```

Example configuration with custom username, password and port:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    settings:
      host: example.com
+     username: appleboy
+     password: 12345678
+     port: 4430
      target: /home/deploy/web
      source: release.tar.gz
```

Example configuration with multiple source and target folder:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    settings:
      host: example.com
      target:
+       - /home/deploy/web1
+       - /home/deploy/web2
      source:
+       - release_1.tar.gz
+       - release_2.tar.gz
```

Example configuration with multiple host:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    settings:
-     host: example.com
+     host:
+       - example1.com
+       - example2.com
      target: /home/deploy/web
      source: release.tar.gz
```

Example configuration with wildcard pattern of source list:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    settings:
      host:
        - example1.com
        - example2.com
      target: /home/deploy/web
      source:
-       - release/backend.tar.gz
-       - release/images.tar.gz
+       - release/*.tar.gz
```

Remove target folder before copy files and artifacts to target:

```diff
  scp:
    image: appleboy/drone-scp
    host: example.com
    settings:
      target: /home/deploy/web
      source: release.tar.gz
+     rm: true
```

Example for remove the specified number of leading path elements:

```diff
  scp:
    image: appleboy/drone-scp
    settings:
      host: example.com
      target: /home/deploy/web
      source: dist/release.tar.gz
+     strip_components: 1
```

Example configuration using ｀SSHProxyCommand｀:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    settings:
      host:
        - example1.com
        - example2.com
      target: /home/deploy/web
      source:
        - release/*.tar.gz
+     proxy_host: 10.130.33.145
+     proxy_user: ubuntu
+     proxy_port: 22
+     proxy_password: 1234
```

Example configuration using password from secrets:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    settings:
      host:
        - example1.com
        - example2.com
      user: ubuntu
      port: 22
-     password: 1234
+     password:
        from_secret: ssh_password
      target: /home/deploy/web
      source:
        - release/*.tar.gz
```

Example configuration using command timeout:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
    settings:
      host:
        - example1.com
        - example2.com
      user: ubuntu
      port: 22
-     command_timeout: 120
+     command_timeout: 2m
        from_secret: ssh_password
      target: /home/deploy/web
      source:
        - release/*.tar.gz
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

timeout
: timeout is the maximum amount of time for the TCP connection to establish

command_timeout
: timeout is the maximum amount of time for execute command

strip_components
: remove the specified number of leading path elements

proxy_host
: proxy hostname or IP

proxy_port
: ssh port of proxy host

proxy_username
: account for proxy host user

proxy_password
: password for proxy host user

proxy_key
: plain text of proxy private key

proxy_key_path
: key path of proxy private key

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
