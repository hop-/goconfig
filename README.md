# GoConfig

![CI](https://github.com/hop-/goconfig/workflows/CI/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/hop-/goconfig)](https://goreportcard.com/report/github.com/hop-/goconfig)[![Go Reference](https://pkg.go.dev/badge/github.com/hop-/goconfig.svg)](https://pkg.go.dev/github.com/hop-/goconfig)

A Go port of Node.js config package (which uses json files to configure application)

From original Library:
> Library organizes hierarchical configurations for your app deployments.
> It lets you define a set of default parameters, and extend them for different deployment environments (development, qa, staging, production, etc.).

## Installation

```shell
go get github.com/hop-/goconfig
```

## Configuration Files

All config files are in `HOST_CONFIG_DIR` directory, default is `config`.

The `HOST_ENV` environment variable defines the application deployment environment.

`$HOST_CONFIG_DIR/default.json`:

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

Override configurations for production when `HOST_ENV` is `production`.

`$HOST_CONFIG_DIR/production.json`:

```json
{
  "Customer": {
    "credit": {
      "initialDays": 30
    }
  }
}
```

### Custom Environment Variables

You can map environment variables to config keys using `$HOST_CONFIG_DIR/custom-environment-variables.json`:

```json
{
  "Customer": {
    "db": {
      "host": "DB_HOST",
      "port": "DB_PORT"
    }
  }
}
```

If the environment variable is set, its value will override the config value.

## Usage

### Load

Load all configurations. **Must be called once before using any other function** (`Get`, `GetAny`, `Has`, `Extract`). Calling those functions before `Load` will panic.

```go
import "github.com/hop-/goconfig"

func main() {
  if err := goconfig.Load(); err != nil {
    // handle error
  }
}
```

### Get

Retrieve a config value by dot-separated path and deserialize it into a typed value. Returns `(*T, error)` — the value is a pointer. Returns an error if the key does not exist.

```go
import "github.com/hop-/goconfig"

host, err := goconfig.Get[string]("Customer.db.host")
if err != nil {
  // handle error
}
fmt.Println(*host) // dereference the pointer

port, err := goconfig.Get[int]("Customer.db.port")
if err != nil {
  // handle error
}

type DbConfig struct {
  Host   string `json:"host"`
  Port   int    `json:"port"`
  DbName string `json:"dbName"`
}

db, err := goconfig.Get[DbConfig]("Customer.db")
if err != nil {
  // handle error
}
fmt.Println(db.Host) // db is *DbConfig
```

### GetAny

Retrieve a config value as `any`. Returns `nil` if the key does not exist. Passing an empty string returns the entire config as `map[string]any`.

```go
import "github.com/hop-/goconfig"

val := goconfig.GetAny("Customer.db.host") // nil if key missing

all := goconfig.GetAny("") // returns entire config as map[string]any
```

### Has

Check whether a config key exists. Returns `true` if the key exists, `false` otherwise.

```go
import "github.com/hop-/goconfig"

if goconfig.Has("Customer.db.host") {
  // key exists
}
```

### Extract

Extract config values into an annotated struct using the `goconfig` struct tag. The tag value is the dot-separated config path relative to the given root. Fields are required by default; missing required fields return a `*ConfigError`. Use `,optional` to make a field optional, or `,required` to be explicit.

```go
import "github.com/hop-/goconfig"

type AppConfig struct {
  DbHost      string `goconfig:"Customer.db.host"`
  DbPort      int    `goconfig:"Customer.db.port"`
  InitialDays int    `goconfig:"Customer.credit.initialDays"`
  OptionalKey string `goconfig:"some.optional.key,optional"`
  ExplicitKey string `goconfig:"some.required.key,required"`
}

cfg := &AppConfig{}
if err := goconfig.Extract("", cfg); err != nil {
  // handle error — a required key is missing
}
```

Nested struct fields (both value and pointer types) are supported. Each nested struct is extracted recursively using its own `goconfig` tags relative to the parent path:

```go
type DbConfig struct {
  Host   string `goconfig:"host"`
  Port   int    `goconfig:"port"`
  DbName string `goconfig:"dbName"`
}

type CustomerConfig struct {
  Db    DbConfig  `goconfig:"db"`  // value type
  DbPtr *DbConfig `goconfig:"db"`  // pointer type — must be pre-allocated
}

type AppConfig struct {
  Customer CustomerConfig `goconfig:"Customer"`
}

cfg := &AppConfig{Customer: CustomerConfig{DbPtr: &DbConfig{}}}
if err := goconfig.Extract("", cfg); err != nil {
  // handle error
}
```

A non-empty first argument scopes extraction to that config path:

```go
dbCfg := &DbConfig{}
if err := goconfig.Extract("Customer.db", dbCfg); err != nil {
  // handle error
}
```

Fields without a `goconfig` tag are silently skipped.

### Error Handling

All errors returned by this package are of type `*ConfigError`:

```go
type ConfigError struct {
  Message string
}

func (e *ConfigError) Error() string
```

You can type-assert to inspect them:

```go
if err := goconfig.Extract("", cfg); err != nil {
  if cfgErr, ok := err.(*goconfig.ConfigError); ok {
    fmt.Println(cfgErr.Message)
  }
}
```
