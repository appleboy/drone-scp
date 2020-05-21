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
- name: scp files
  image: appleboy/drone-scp
  settings:
    host: example.com
    username: foo
    password: bar
    port: 22
    target: /var/www/deploy/${DRONE_REPO_OWNER}/${DRONE_REPO_NAME}
    source: release.tar.gz
```

Example configuration with multiple source and target folder:

```diff
  - name: scp files
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
  - name: scp files
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
  - name: scp files
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
  - name: scp files
    image: appleboy/drone-scp
    settings:
      target: /home/deploy/web
      source: release.tar.gz
+     rm: true
```

Example for remove the specified number of leading path elements:

```diff
  - name: scp files
    image: appleboy/drone-scp
    settings:
      host: example.com
      target: /home/deploy/web
      source: dist/release.tar.gz
+     strip_components: 1
```

Example configuration using ｀SSHProxyCommand｀:

```diff
  - name: scp files
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
  - name: scp files
    image: appleboy/drone-scp
    settings:
      host:
        - example1.com
        - example2.com
      user: ubuntu
      port: 22
-     password: 1234
+     password:
+       from_secret: ssh_password
      target: /home/deploy/web
      source:
        - release/*.tar.gz
```

Example configuration using command timeout:

```diff
  - name: scp files
    image: appleboy/drone-scp
    settings:
      host:
      - example1.com
      - example2.com
      user: ubuntu
      password:
      from_secret: ssh_password
      port: 22
-     command_timeout: 120
+     command_timeout: 2m
      target: /home/deploy/web
      source:
      - release/*.tar.gz
```

Example configuration for ignore list:

```diff
  - name: scp files
    image: appleboy/drone-scp
    settings:
      host:
        - example1.com
        - example2.com
      user: ubuntu
      password:
        from_secret: ssh_password
      port: 22
      command_timeout: 2m
      target: /home/deploy/web
      source:
+       - !release/README.md
        - release/*
```

Example configuration for passphrase which protecting a private key:

```diff
  - name: scp files
    image: appleboy/drone-scp
    settings:
      host:
        - example1.com
        - example2.com
      user: ubuntu
+     key:
+       from_secret: ssh_key
+     passphrase: 1234
      port: 22
      command_timeout: 2m
      target: /home/deploy/web
      source:
        - release/*
```

## Parameter Reference

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

passphrase
: The purpose of the passphrase is usually to encrypt the private key.

fingerprint
: fingerprint SHA256 of the host public key, default is to skip verification

target
: folder path of target host

source
: source lists you want to copy

rm
: remove target folder before copy files and artifacts

timeout
: Timeout is the maximum amount of time for the ssh connection to establish, default is 30 seconds.

command_timeout
: Command timeout is the maximum amount of time for the execute commands, default is 10 minutes.

strip_components
: remove the specified number of leading path elements

tar_tmp_path
: temporary path for tar file on the dest host

tar_exec
: alternative `tar` executable to on the dest host

overwrite
: use `--overwrite` flag with tar

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

proxy_passphrase
: The purpose of the passphrase is usually to encrypt the private key.

proxy_fingerprint
: fingerprint SHA256 of the host public key, default is to skip verification

## Template Reference

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
