# [JPL - JSON Processing Language](../docs/index.md) - JavaScript implementation

## Code example

```js
import jpl from '@jplorg/jpl';

(async () => {
  const inputs = [null];

  const results = await jpl.run('"Hello, 🌎!"', inputs);

  console.log(...results);
})();
```

## REPL

The package provides a CLI REPL, which can be used as a language playground.

```sh
> jpl-repl # or `npm --prefix js start`

Welcome to JPL v1.0.0.
Type ":h" for more information.

> "Hello, \('🌎', 'JPL')!"
"Hello, 🌎!", "Hello, JPL!"
>
```

## Extending JPL

TODO: inform about the runtime API, functions, JPLTypes and different error types
