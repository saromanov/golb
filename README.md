# golb
[![Go Report Card](https://goreportcard.com/badge/github.com/saromanov/golb)](https://goreportcard.com/report/github.com/saromanov/golb)
[![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/36ee1a51d3914831ad38546c85281e31)](https://www.codacy.com/app/saromanov/golb?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=saromanov/golb&amp;utm_campaign=Badge_Grade)
[![Build Status](https://travis-ci.org/saromanov/golb.svg?branch=master)](https://travis-ci.org/saromanov/golb)
[![Coverage Status](https://coveralls.io/repos/github/saromanov/golb/badge.svg?branch=master)](https://coveralls.io/github/saromanov/golb?branch=master)

Load balancer

## Table of Contents
*  [Getting Started](#getting-started)
    +  [Installing](#installing)
    +  [Quick start](#quick-start)
    +  [Metrics](#metrics)
    +  [Config](#config)

## Getting Started

### Installing

```sh
$ go build github.com/saromanov/golb/...
```

### Quick start

For quick test of the Golb, create one or more docker container with servers. Check "test" dir for create simple python server like
```sh
sudo docker build -t test/server .
sudo docker run test/server -p 7000:7000
sudo docker run test/server -p 7000:7001
```

Then, start golb instance
```sh
golb
```

After this, you can send test requests on the 
```sh
curl http://localhost:8080
```

### Metrics

```sh
curl http://localhost:8080/v1/metrics
```

### Config

```sh
curl http://localhost:8080/v1/config
```