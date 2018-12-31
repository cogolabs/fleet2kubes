[![Build Status](https://travis-ci.org/cogolabs/fleet2kubes.svg?branch=master)](https://travis-ci.org/cogolabs/fleet2kubes)
[![codecov](https://codecov.io/gh/cogolabs/fleet2kubes/branch/master/graph/badge.svg)](https://codecov.io/gh/cogolabs/fleet2kubes)
[![Maintainability](https://api.codeclimate.com/v1/badges/699f80c897e5cd1865ba/maintainability)](https://codeclimate.com/github/cogolabs/fleet2kubes/maintainability)

# fleet2kubes

Convert fleet `.service` files to kubes deployment + service.

## Background

Cogo Labs adopted Docker early in 2014 and leveraged CoreOS Fleet to deploy and manage a rich distributed microservices platform. Kubernetes, based on Google's internal manager, Borg, matured to production-ready quality in 2018, and has fully absorbed CoreOS Fleet technologies, such as Etcd: the Kubernetes Control Plane, into an awesome ecosystem now led by the CNCF. This minimal glue automates parts of our Kubernetes migration.

## Caveats

Mostly works for generic web services -- does not yet handle:
- `.timer`s => cron jobs
- docker flags eg. environ, privilege
