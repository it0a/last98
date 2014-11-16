[![Build Status](https://travis-ci.org/it0a/last98.svg?branch=master)](https://travis-ci.org/it0a/last98)

last98
======

# image uploader/gallery written in go #

## Install Guide ##

```
$ git clone https://github.com/it0a/last98 $GOPATH/src/last98
$ cd $GOPATH/src/last98
$ go install
```

## Heroku Install Guide ##

```
$ cd $GOPATH/src/last98
$ go get -u github.com/kr/godep
$ godep save
$ heroku create -b https://github.com/kr/heroku-buildpack-go.git
$ heroku addons:add heroku-postgresql
$ git push heroku master
```
