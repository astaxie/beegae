## beegae

[beego](http://github.com/astaxie/beego) is a Go Framework inspired by tornado and sinatra.

It is a simple & powerful web framework. Go to the project page to learn more about beego!

More info [beego.me](http://beego.me)

**beegae** is a port of beego intended to be used on Google's AppEngine. There are a few subtle differences between how beego and beegae initalizes applications which you can see for yourself here [example](https://github.com/astaxie/beegae/tree/master/example)

The aim of this project is to keep as much of beego unchanged as possible in beegae.

## Features

* Datastore + Memcached backed session store! [read more here](https://github.com/astaxie/beegae/tree/develop/session#beegae-session) to learn how to use it.
* `AppEngineCtx` is part of the default Controller. View the included sessions package [documentation](https://github.com/astaxie/beegae/tree/develop/session#beegae-session) for an example of using it!

## beego Documentation

* [English](http://beego.me/docs/intro/)
* [中文文档](http://beego.me/docs/intro/)

## Community

* [http://beego.me/community](http://beego.me/community)

## Getting Started

This will be a quick overview of how to setup the repository and get started with beegae on AppEngine. It is already assumed that you have setup the AppEngine SDK for Go correctly and setup your GOPATH correctly

* `# goapp get github.com/astaxie/beegae`
* `# goapp get github.com/beego/bee`
* `# cd $GOPATH/src`
* `# $GOPATH/bin/bee new hellogae`
* `# cd hellogae`
* Make a new file app.yaml and fill it as such:

```yaml
application: hellobeegae
version: 1
runtime: go
api_version: go1

handlers:
- url: /.*
  script: _go_app
```
* `# gofmt -r '"github.com/astaxie/beego" -> "github.com/astaxie/beegae"' -w ./`
* `# gofmt -r 'beego -> beegae' -w ./`
* `# mkdir main && mv main.go main/ && mv app.yaml main/ && mv conf/ main/ && mv views/ main/ && cd main/`
* Now open `main.go` and change `func main()` to `func init()`
* `# goapp serve`
* Done!

## LICENSE

beego is licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).
