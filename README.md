# Go Start [![Circle CI](https://gitlab.com/westworld-ho/gostartup/badges/master/pipeline.svg)][![godoc]()][![code coverage]()]

gostartup is my go-to starter project base for a golang project.

It is an opinionated collection of select packages and reference implementations for building a range of projects from small CLI utilities, independent microservices, to mono repo for multiple services of a SaaS backend.

In addition, it also serves me as a documented collection of thoughts on how to do common tasks around monitoring/deploying services and managing secrets.

### It's not...

- a cookiecutter similar to [django cookiecutter](). cookiecutter templates are more cleaner and abstract, but the purpose here is to get started quickly with reference/sample code, instead of providing clean abstractions
- a production-ready implementation of any feature. The goal is to have a project-base which will help get there quickly.
- There is no performance criterion as of yet, as there are no actual usecases yet. So frameworks, packages are chosen for extendability/ease of deployment, rather than performance. Also [[]premature optimization is the root of evil]

## Features

- Reference `user:account` data models in postgres and APIs supporting signup/login/account recovery usecases using [authboss](https://github.com/volatiletech/authboss)
- permissions
- Reference `stats` service with API stats and basic rate limiting using redis
- Reference JWT implementation for authentication along with basic auth for service authentication
- Reference notification implementation for posting payloads to urls using exponential back-off.
- `echo` based API server with reference middlewares.
- A simple in-memory job queue using []() to asynchronously process lightweight workloads: using an example of notifying a URL with a payload
- Auto-generating OpenAPI[] specs using swagger for services provided
- env based config supporting 12 factor apps[]
- Feature flags
- Clean-ish Architecture[]
- Template for building multiple service/binaries based on mono repo
- go modules for package management
- CLI ready
- CI support
- Docker support

### Project structure

### Clean-ish Architecture

### Admin

As with the

#### Admin UI

Use react-admin or vue-admin with apis

#### Admin backend

Use [prest](http://postgres.rest/) for converting uncovered models/tables to API. Use it with a whitelist, covering only tables which don't need admin usecases/business logic behind them. Another option is to mount prest using a [custom backend](https://postgres.rest/prest-as-web-framework/) and add additional admin endpoints

### Logging and Monitoring

### Linting

### Env based config

All the settings for the services are injected via environment variables in accordance with 12-factor[] app.

KMS + S3
Vault + Consul + https://github.com/hashicorp/envconsul

### Mono repo and multiple services

https://blog.digitalocean.com/cthulhu-organizing-go-code-in-a-scalable-repo/
https://gomonorepo.org/

### Account service

### Tests

### Packages/Libraries used

### CI

### Containerized services

#### Why not...

- gRPC

## Installation

0. Install golang in your system(1.10 preferably)
1. Set GOPATH. Export $GOPATH=${HOME}/gopath. Add this to your bashrc/zshrc.
1. Initiate GOPATH. Go on install something: `go get -u gopkg.in/alecthomas/gometalinter.v2`
1. Clone repo under GOPATH. `cd $GOPATH/src; git clone git@gitlab.com:westworld-ho/gostartup.git`
1. Install dep in your system. Either `curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh` or `brew install dep`
1. Install all project dependencies. `cd $GOPATH/src/gostartup; dep ensure`
1. That's It. Now edit code with your favorite editor.
1. To run: `go run src/main.go -debug`. `-debug` is used for local testing.
1. To run tests: `cd src; go test ./...`

## Usage

1. Setup and install Postgres. `brew install postgresql` and `brew install pgcli`
2. Setup postgres for your user. `sudo -u postgres createuser -s ${USERNAME}`
3. Create hostdb and user. `createdb hostdb; createuser gostartup sweetwater`
4. Install redis. `brew install redis` or `sudo apt install redis-server redis-tools`

## TODO

rest
  - swagger
  - hystrix
  - monitor: prometheus/statsd
tests
  - bdd/table-driven: ginkgo
  - mock: testify
  - faker

---
- user:account(flask)
  - login(post, delete logout)
  - signup(post)
  - recovery(password reset form + email link)
  - emails
  - auth/session validation
  - rbac
- api doc
  - payments: capture, charge, txn history, manage fi, balance??
  - webhook(+jobq) for recording transaction
  - notification for upstream services(callback url + backoff mechanism)
  - cache? : redlock for charge id lock,, balance
