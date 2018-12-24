# gorpool

gorpool is thinnest goroutine pool library for go application

## Requirement

Go (>= 1.9)

## Install

```
go get github.com/hlts2/gorpool
```

## Example

```go
package main

import (
        "context"
        "fmt"

        "github.com/hlts2/gorpool"
        "github.com/hlts2/gorpool/middlewares"
        "github.com/hlts2/gorpool/middlewares/logger/zap"
        "github.com/hlts2/gorpool/middlewares/recovery"
        "go.uber.org/zap"
)

const (
        errChBuffer = 1
        jobChBuffer = 1000000
        poolSize    = 256
)

func main() {
        z := zap.NewExample()

        g := gorpool.New(
                gorpool.WithError(make(chan error, errChBuffer)),
                gorpool.WithPoolSize(poolSize),
                gorpool.WithJobSize(jobChBuffer*2),
                gorpool.WithInterceptor(
                        middlewares.ChainInterceptors(
                                gorpool_recovery.Interceptor(
                                        func(p interface{}) {
                                                z.Info("recovery from panic")
                                        },
                                ),
                                gorpool_zap.Interceptor(
                                        z,
                                ),
                        ),
                ),
        ).Start(context.Background())

        for i := 0; i < jobChBuffer; i++ {
                // Register job.
                g.Add(func(context.Context) error {
                        z.Info("finish job")
                        return nil
                })
        }

        // Get error of job. If execution of all jobs is completed, exit Wait function.
        for err := range g.Wait() {
                if err != nil {
                        fmt.Println(err)
                }
        }

        // Register job again.
        for i := 0; i < jobChBuffer; i++ {
                g.Add(func(context.Context) error {
                        z.Info("finish job")
                        return nil
                })
        }

        for err := range g.Wait() {
                if err != nil {
                        fmt.Println(err)
                }
        }
}
```


## Author
[hlts2](https://github.com/hlts2)

## LICENSE
gorpool released under MIT license, refer [LICENSE](https://github.com/hlts2/gorpool/blob/master/LICENSE) file.
