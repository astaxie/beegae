## beegae

[beego](http://github.com/astaxie/beego) is a Go Framework inspired by tornado and sinatra.
[![beego Build Status](https://travis-ci.org/astaxie/beego.svg?branch=master)](https://travis-ci.org/astaxie/beego)
[![beego GoDoc](http://godoc.org/github.com/astaxie/beego?status.svg)](http://godoc.org/github.com/astaxie/beego)


beego is used for rapid development of RESTful APIs, web apps and backend services in Go.
It is inspired by Tornado, Sinatra and Flask. beego has some Go-specific features such as interfaces and struct embedding.

More info [beego.me](http://beego.me)

**beegae** is a port of beego intended to be used on Google's AppEngine. There are a few subtle differences between how beego and beegae initializes applications which you can see for yourself here [example](https://github.com/astaxie/beegae/tree/master/example)

The aim of this project is to keep as much of beego unchanged as possible in beegae.

## IMPORTANT UPDATE - Breaking Changes with beegae 1.6.1
There were changes made to the beego core that were updated in beegae for consistency between packages. Many of these changes are in relation to refactoring various objects in beego as well addressing
golint recommendations to variable/function names. This may result in minor tweaks to your code (e.g. UrlFor is now URLFor on controllers).

This package now supports both classic AppEngine runtime as well as flexible (formerly Managed VM) runtime.

## IMPORTANT UPDATE - Breaking Changes with beegae 1.5

With the latest update beegae now uses the [google.golang.org/appengine](https://godoc.org/google.golang.org/appengine) AppEngine package. You can read the differences from the classic `appengine` package on the new package's [repository](https://github.com/golang/appengine).
beegae users will need to update their apps to use this package as per the instructions on the repository. The `controller.AppEngineCtx` type was changed to `context.Context` so minimal changes to users' code should be necessary.

The session code was also updated to reflect these changes, most notably, the `context.Context` is now the first argument for the session interface functions. Additionally, the package was moved to be more in-line with how `beego` broke out its session implementations.

Like other session providers as part of `beego`, you need to do an blank import of the `appengine` session provider: `import _ "github.com/astaxie/beegae/session/appengine"`
This will register the `appengine` session provider so you may use it in your application. You may do this in your `main.go` file. Look at the [examples](https://github.com/astaxie/beegae/tree/master/example) for guidance.

My apologies for any inconvenience this brings to your code. The new package is the recommended approach to Go on AppEngine as it works on both classic AppEngine and Managed VMs, and so beegae was updated to support this recommendation.

As always, if there any bugs with the package please open an issue and/or submit a PR.

## Features

* Datastore + Memcached backed session store! [read more here](https://github.com/astaxie/beegae/tree/master/session#beegae-session) to learn how to use it.
* `AppEngineCtx` is part of the default Controller. View the included sessions package [documentation](https://github.com/astaxie/beegae/tree/master/session#beegae-session) for an example of using it!

## beego Documentation

* [English](http://beego.me/docs/intro/)
* [中文文档](http://beego.me/docs/intro/)
* [Русский](http://beego.me/docs/intro/)

## Community

* [http://beego.me/community](http://beego.me/community)
* Welcome to join us in Slack: [https://beego.slack.com](https://beego.slack.com), you can get invited from [here](https://github.com/beego/beedoc/issues/232)

## Getting Started

This will be a quick overview of how to setup the repository and get started with beegae on AppEngine. It is already assumed that you have setup the AppEngine SDK for Go correctly and setup your GOPATH correctly

* `# go get github.com/astaxie/beegae`
* `# go get github.com/beego/bee`
* `# cd $GOPATH/src`
* `# $GOPATH/bin/bee new hellogae`
* `# cd hellogae`
* Make a new file app.yaml and fill it as such:

```yaml
application: hellobeegae
version: 1
runtime: go
threadsafe: true
api_version: go1

handlers:
- url: /.*
  script: _go_app
```

or if you are using the flexible runtime:

```yaml
runtime: go
vm: true
threadsafe: true
api_version: 1

handlers:
- url: /.*
  script: _go_app
```
* `# gofmt -r '"github.com/astaxie/beego" -> "github.com/astaxie/beegae"' -w ./`
* `# gofmt -r 'beego -> beegae' -w ./`
* `# mkdir main && mv main.go main/ && mv app.yaml main/ && mv conf/ main/ && mv views/ main/ && cd main/`
* If you are using the classic AppeEngine runtime, open `main.go` and change `func main()` to `func init()`, otherwise skip this step.
* `# goapp serve`
* Done!

## LICENSE

beego source code is licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).
