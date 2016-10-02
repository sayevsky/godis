Godis -- in-memory only key-value storage.
========
[![Build Status](https://travis-ci.org/sayevsky/godis.svg?branch=master)](https://travis-ci.org/sayevsky/godis.svg?branch=master)

This is a toy example of in-memory database. Just to show how golang is awesome.

Quick start
-----------
To build 

`go build`

Then to run

`./godis`

By default godis use `6380` port, but it possible to override by `-port=<your port>`. Type `./godis -h` to get list of parameters.

When godis is started it can store and retrieve some data. We can use golang client or `telnet` as example. For more information and examples of telnet commands to godis visit [wiki](https://github.com/sayevsky/godis/wiki)
