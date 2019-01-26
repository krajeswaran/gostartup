# Go Start

gostartup is my go-to starter pack/template for a golang API backend. 

The goal is to have a common, opiniated project base which I can use for my backend projects: from isolated microservices to entire backend of a SaaS app

[![pipeline status](https://gitlab.com/westworld-ho/gostartup/badges/master/pipeline.svg)](https://gitlab.com/westworld-ho/gostartup/commits/master)

### It's not...

* a production-ready implementation of anything useful. The goal is to have a project which is *near* production-ready and versatile enough to be used as a base for different projects
* an optimally performant system. There is no "performance" critiriea as of yet, I don't know the actual usecases. So frameworks, packages are chosen for extendability, rather than performance. Also [[]premature optimization is the root of evil]
* a cookiecutter similar to [](). cookiecutter templates are more cleaner and abstract, but the purpose here is to get started quickly with reference/sample code, instead of providing clean abstractions

## Features

* `echo` based API server with concrete API examples 
* Reference `user:account` data models in postgres and APIs supporting signup/login/account recovery usecases
* Reference `stats` service with API stats and basic rate limiting using redis 
* Reference JWT implementation for authentication along with basic auth for service authentication
* Reference Webhook implementation for listening to callback responses 
* []() A simple job queue to asynchronusly process lightweight workloads: like posting to a webhook or processing an incoming webhook response
* Auto-generating OpenAPI[] specs using swagger for services provided
* env based config supporting 12 factor apps[] 
* Feature flags
* Clean-ish Architecture[]
* Mono repo
* Template for building multiple service/binaries based on code
* `dep` for package management
* CI support
* Containerized services


### Project structure

### Clean-ish Architecture

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

* GraphQL?

* gRPC

## Installation

0. Install golang in your system(1.10 preferably)
1. Set GOPATH. Export $GOPATH=${HOME}/gopath. Add this to your bashrc/zshrc.
2. Initiate GOPATH. Go on install something: `go get -u gopkg.in/alecthomas/gometalinter.v2`
3. Clone repo under GOPATH. `cd $GOPATH/src; git clone git@gitlab.com:westworld-ho/gostartup.git`
4. Install dep in your system. Either `curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh` or `brew install dep`
5. Install all project dependencies. `cd  $GOPATH/src/gostartup; dep ensure`
6. That's It. Now edit code with your favorite editor. 
7. To run: `go run src/main.go -debug`. `-debug` is used for local testing.
8. To run tests: `cd src; go test ./...`

## Usage

1. Setup and install Postgres. `brew install postgresql` and `brew install pgcli`
2. Setup postgres for your user. `sudo -u postgres createuser -s ${USERNAME}`
3. Create hostdb and user. `createdb hostdb; createuser gostartup sweetwater`
4. Install redis. `brew install redis` or `sudo apt install redis-server redis-tools`

## TODO

[] env files
[] user service
[] stats service
[] jwt
[] tests
[] callback 
[] webhook
[] docker test
[] ci

https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use

