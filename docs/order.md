# [JPL - JSON Processing Language](index.md) - Operator Order

The first operator in the list is the first operator that the interpreter attempts to parse.
When parsing one operator fails, the interpreter attempts the next operator in the list.

Usually, all operators also fork into the next operator in the list.
The only exception are some operators that fork into another pipe operator or subpipe operator.
The subpipe operator handles the same operations as the pipe operator, but does not allow output concatenation in its operands, because it is used in contexts where the `,` token is already reserved for value separation.

The last operator in the list is executed first.
When multiple operators share the same execution priority, they occur grouped together in this list.
Those operators are executed from left to right.

- `$ | $` (pipe) <- subpipe

---

- `$, $` (not available in subpipes and subroutes)

---

- `try $ catch $` <- subroute

---

- `$ or $`

---

- `$ and $`

---

- `$ == $`
- `$ != $`

---

- `$ < $`
- `$ <= $`
- `$ > $`
- `$ >= $`

---

- `not $`

---

- `$ ?`

---

- `$ - $`
- `$ + $`

---

- `$ % $`
- `$ / $` (divide by)
- `$ * $`

---

- `$ ?? $`

---

- `- $`

---

- `if $pipe then $pipe elif $pipe then $pipe else $pipe end`

---

- `true`
- `false`
- `null`

---

- `number`

---

- `func variable-name(variable-name, ...): $subroute`

---

- `func (variable-name, ...): $subroute`

---

- `variable-name`
- `variable-name<.path> = $subroute`
- `variable-name<.path> |= $subroute`
- `variable-name<.path> -= $subroute`
- `variable-name<.path> += $subroute`
- `variable-name<.path> %= $subroute`
- `variable-name<.path> /= $subroute`
- `variable-name<.path> *= $subroute`
- `variable-name<.path> ?= $subroute`
- `$.path = $subroute`
- `$.path |= $subroute`
- `$.path -= $subroute`
- `$.path += $subroute`
- `$.path %= $subroute`
- `$.path /= $subroute`
- `$.path *= $subroute`
- `$.path ?= $subroute`
- `.`
- `.path`
- `$<.path>.key`
- `$<.path>[$pipe]`
- `$<.path>[$pipe:$pipe]`
- `$<.path>[]`
- `$<.path>($subpipe, ...)`

---

- `{ ($pipe): $subpipe, "": $subpipe, variable-name: $subpipe, ... }`
- `[ $pipe ]`
- `""`, `"\($pipe)"`

---

- `($pipe)`
