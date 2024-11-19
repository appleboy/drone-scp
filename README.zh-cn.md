# drone-scp

[繁體中文](README.zh-tw.md) | [English](README.md)

[![GoDoc](https://godoc.org/github.com/appleboy/drone-scp?status.svg)](https://godoc.org/github.com/appleboy/drone-scp)
[![Lint and Testing](https://github.com/appleboy/drone-scp/actions/workflows/lint.yml/badge.svg)](https://github.com/appleboy/drone-scp/actions/workflows/lint.yml)
[![codecov](https://codecov.io/gh/appleboy/drone-scp/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-scp)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-scp)](https://goreportcard.com/report/github.com/appleboy/drone-scp)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-scp.svg)](https://hub.docker.com/r/appleboy/drone-scp/)

复制文件和工件通过 SSH 使用二进制文件、docker 或 [Drone CI](http://docs.drone.io/)。

[English](README.md) | [繁體中文](README.zh-tw.md)

## 功能

* [x] 支持例程。
* [x] 支持来源列表中的通配符模式。
* [x] 支持将文件发送到多个主机。
* [x] 支持将文件发送到主机上的多个目标文件夹。
* [x] 支持从绝对路径或原始主体加载 ssh 密钥。
* [x] 支持 SSH ProxyCommand。
