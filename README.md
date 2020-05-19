# go-redis-bug

Prerequisite:

1. Setup a local redis cluster (3x2) either manually or via utility
[script](http://download.redis.io/redis-stable/utils/create-cluster/create-cluster)
2. Update `redisAddr` in `main.go`

Build:
```sh
$ cat go.mod | grep "redis v"
# 	github.com/go-redis/redis v6.15.2+incompatible
# this version has the bug
$ go build
```

Reproducible almost always on my machine, YMMV:
```sh
$ while ./go-redis-bug; do :; done
  14% |████                               | (148019/999999, 19302 it/s) [8s:44s]2020/05/20 00:47:32
json.Unmarshal() failed: invalid character 'O' looking for beginning of value
```

Not reproducible on `6.15.7` which contains the [backported fix](https://github.com/go-redis/redis/pull/1124).
```sh
# update go-redis version
$ go get github.com/go-redis/redis@v6.15.7
go: github.com/go-redis/redis v6.15.7 => v6.15.7+incompatible
$ go build
# let this run for as long as you want...
$ while ./go-redis-bug; do :; done
 100% |███████████████████████████████████| (999999/999999, 24875 it/s) [40s:0s]
 100% |███████████████████████████████████| (999999/999999, 22628 it/s) [44s:0s]
 100% |███████████████████████████████████| (999999/999999, 16726 it/s) [59s:0s]
 100% |███████████████████████████████████| (999999/999999, 16778 it/s) [59s:0s]
 100% |██████████████████████████████████| (999999/999999, 16545 it/s) [1m0s:0s]
 100% |██████████████████████████████████| (999999/999999, 16152 it/s) [1m1s:0s]
  15% |█████                              | (150745/999999, 21537 it/s) [8s:39s]
```
