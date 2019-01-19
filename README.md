[![Build Status](https://travis-ci.org/cogolabs/fleet2kubes.svg?branch=master)](https://travis-ci.org/cogolabs/fleet2kubes)
[![codecov](https://codecov.io/gh/cogolabs/fleet2kubes/branch/master/graph/badge.svg)](https://codecov.io/gh/cogolabs/fleet2kubes)
[![Go Report Card](https://goreportcard.com/badge/github.com/cogolabs/fleet2kubes)](https://goreportcard.com/report/github.com/cogolabs/fleet2kubes)
[![Maintainability](https://api.codeclimate.com/v1/badges/699f80c897e5cd1865ba/maintainability)](https://codeclimate.com/github/cogolabs/fleet2kubes/maintainability)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

# fleet2kubes

Converts fleet systemd units to kubernetes resources.

## Background

Cogo Labs adopted Docker early in 2014 and leveraged CoreOS Fleet to deploy and manage a rich distributed microservices platform. Kubernetes, based on Google's internal manager, Borg, matured to production-ready quality in 2018, and has fully absorbed CoreOS Fleet technologies, such as `etcd` (backend for the Kubernetes Control Plane), into an awesome ecosystem now led by the CNCF. This minimal glue automates parts of our Kubernetes migration.

## Install and Usage

The first command will install `f2k` to `~/go/bin/f2k`. For convenience, append `$GOPATH/bin` to your `$PATH`.

```sh
  go get -u github.com/cogolabs/fleet2kubes/cmd/f2k

  f2k [flags] my.service > my.yaml
```
