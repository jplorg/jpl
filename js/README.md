# [JPL - JSON Processing Language](../docs/index.md) - JavaScript implementation

## Code example

```js
import jpl from 'jpli';

(async () => {
  const inputs = [null];

  const results = await jpl.run('"Hello, 🌎!"', inputs);

  console.log(...results);
})();
```

## REPL

The package provides a CLI REPL, which can be used as a language playground.

```sh
> jpl # or `npm start`

Welcome to JPL v1.0.0.
Type ":h" for more information.

> "Hello, \('🌎', 'JPL')!"
"Hello, 🌎!", "Hello, JPL!"
>
```

For debugging purposes or to be able to dive into how a code snippet is interpreted, the REPL provides a `:i` directive.

```sh
> :i "Hello, 🌎!"
[
  {
    "op": "\"\"",
    "params": {
      "value": "Hello, 🌎!"
    }
  }
]
>
```

## Extending JPL

TODO: inform about the runtime API, functions, JPLTypes and different error types
