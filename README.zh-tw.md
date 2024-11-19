# drone-scp

[簡體中文](README.zh-cn.md) | [English](README.md)

[![GoDoc](https://godoc.org/github.com/appleboy/drone-scp?status.svg)](https://godoc.org/github.com/appleboy/drone-scp)
[![Lint and Testing](https://github.com/appleboy/drone-scp/actions/workflows/lint.yml/badge.svg)](https://github.com/appleboy/drone-scp/actions/workflows/lint.yml)
[![codecov](https://codecov.io/gh/appleboy/drone-scp/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-scp)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-scp)](https://goreportcard.com/report/github.com/appleboy/drone-scp)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-scp.svg)](https://hub.docker.com/r/appleboy/drone-scp/)

複製檔案和工件通過 SSH 使用二進制檔案、docker 或 [Drone CI](http://docs.drone.io/)。

[English](README.md) | [简体中文](README.zh-cn.md)

## 功能

* [x] 支援例程。
* [x] 支援來源列表中的萬用字元模式。
* [x] 支援將檔案發送到多個主機。
* [x] 支援將檔案發送到主機上的多個目標資料夾。
* [x] 支援從絕對路徑或原始主體載入 ssh 金鑰。
* [x] 支援 SSH ProxyCommand。
