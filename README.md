# grpool

grpool is thinnest goroutine pool library for go application

## Requirement

Go (>= 1.9)

## Install

```
go get github.com/hlts2/grpool
```

## Example

```go
package main

import (
        "context"
        "fmt"

        "github.com/hlts2/grpool"
        "github.com/hlts2/grpool/middlewares"
        "github.com/hlts2/grpool/middlewares/logger/zap"
        "github.com/hlts2/grpool/middlewares/recovery"
        "go.uber.org/zap"
)

const (
        errChBuffer = 1
        jobChBuffer = 1000000
        poolSize    = 256
)

func main() {
        z := zap.NewExample()

        g := grpool.New(
                grpool.WithError(make(chan error, errChBuffer)),
                grpool.WithPoolSize(poolSize),
                grpool.WithJobSize(jobChBuffer*2),
                grpool.WithInterceptor(
                        middlewares.ChainInterceptors(
                                grpool_recovery.Interceptor(
                                        func(p interface{}) {
                                                z.Info("recovery from panic")
                                        },
                                ),
                                grpool_zap.Interceptor(
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
gocache released under MIT license, refer [LICENSE](https://github.com/hlts2/grpool/blob/master/LICENSE) file.
