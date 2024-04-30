# [JPL - JSON Processing Language](../docs/index.md) - Golang implementation

## Code example

```go
package main

import (
  "fmt"

  gojpl "github.com/jplorg/jpl/go"
)

func main() {
  inputs := []any{nil}

  results, err := gojpl.Run(`"Hello, 🌎!"`, inputs, nil)
  if err != nil {
    panic(err)
  }

  fmt.Println(results...)
}
```

## REPL

The package provides a CLI REPL, which can be used as a language playground.

```sh
> jpl-repl # or `go run github.com/jplorg/jpl/go/jpl-repl`

Welcome to JPL.
Type ":h" for more information.

> "Hello, \('🌎', 'JPL')!"
"Hello, 🌎!", "Hello, JPL!"
>
```

## Extending JPL

TODO: inform about the runtime API, functions, JPLTypes and different error types
