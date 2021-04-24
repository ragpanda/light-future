# light-future
A lightweight future, allows users to better abstract concurrency logic.

![Go Report Card](https://goreportcard.com/badge/github.com/ragpanda/light-future) ![Build status](https://travis-ci.com/ragpanda/light-future.svg?branch=master)



# Intro
Using Future tends to be more readable, asynchronous behavior is abstracted under the Future, and code tends to be more elegant

Example
```golang

    result, err := NewFuture(ctx, NewClosureRunnable(func(ctx context.Context) (interface{}, error) {
        b := a
        return b, nil
    })).
    Send().
    Await().
    Result()

```


# Benchmark

Compare using future and pass the result by the channel, you can compare this cost to determine whether to use

```
goos: linux
goarch: amd64
pkg: github.com/ragpanda/light-future
cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
BenchmarkFuture-8                        	   10000	    107847 ns/op	   44655 B/op	     998 allocs/op
BenchmarkGoroutineUsingChannelReturn-8   	   14350	     83810 ns/op	   15829 B/op	     398 allocs/op
PASS
ok  	github.com/ragpanda/light-future	3.142s
```
