# GoConfig

![CI](https://github.com/hop-/goconfig/workflows/CI/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/hop-/goconfig)](https://goreportcard.com/report/github.com/hop-/goconfig)[![Go Reference](https://pkg.go.dev/badge/github.com/hop-/goconfig.svg)](https://pkg.go.dev/github.com/hop-/goconfig)

A Go port of Node.js config package (which uses json files to configure application)

From original Library:
> Library organizes hierarchical configurations for your app deployments.
> It lets you define a set of default parameters, and extend them for different deployment environments (development, qa, staging, production, etc.).

## Installation

As a library

```shell
go get github.com/hop-/goconfig
```

Usage:

All config files are in `HOST_CONFIG_DIR` directory, default is 'config'

It is using `HOST_ENV` environment variable to define the application deployment environment

`$HOST_CONFIG_DIR`/default.json:

```json
{
  "Customer": {
    "db": {
      "host": "localhost",
      "port": 27017,
      "dbName": "customers"
    },
    "credit": {
      "initialLimit": 100,
      "initialDays": 1
    }
  }
}
```

Override some configurations for production when `HOST_ENV` is 'production'.

`$HOST_CONFIG_DIR`/production.json:

```json
{
  "Customer": {
    "credit": {
      "initialDays": 30
    }
  }
}
```

Use config in your code:

```go
import "github.com/hop-/goconfig"

type Consumer struct {
  Consumer struct {
    Db struct {
      host string
      port int
      dbName string
    }
    Credit struct {
      InitialLimit int
      InitialDays int
    }
  }
}

func main() {
  if err := goconfig.Load(); err != nil {
    // Some error handling
  }

  consumer, err := goconfig.Get[Consumer]("Consumer")
  if err != nil {
    // Some error handling
  }

  host, err := goconfig.Get[string]("Consumer.Db.host")
  if err != nil {
    // Some error handling
  }

  // Your code
}
```
