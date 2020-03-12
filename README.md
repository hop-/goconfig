# GoConfig

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

All config files are in config directory

default.json:
```json
{
  // Customer configurations
  "Customer": {
    "db": {
      "host": "localhost",
      "port": 27017,
      "dbName": "customers"
    },
    "credit": {
      "initialLimit": 100,
      // Set low for development
      "initialDays": 1
    }
  }
}
```

Override some configurations for production.
production.json:
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

func main() {
  if err := goconfig.Load(); err != nil {
    log.Info(err.Error())
  }
  
  consumer := goconfig.Get("Consumer")

  consumerStructured := struct {
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
  goconfig.GetObject("Consumer", &consumerStructured)
}
```
