# [JPL - JSON Processing Language](index.md) - Specification

## How to use this document

Regexp expressions are used in this document to describe specific formats.
`\d` is used for decimal digits and is equivalent to `[0-9]`.

When describing definitions, other definitions are referenced like `{definition-name}` and expressions which must result in a specific type are referred like `{expression -> type}`.
This form is not part of the language but is only used for this documentation.

## Basic concepts

A program is built by one or more expressions, separated by the pipe operator.

Each expression takes multiple inputs, processes them one by one, and can also generate any number of outputs afterwards.
All outputs are concatenated and passed into the next expression.
Expressions that reference variables also run for each single value of each variable.
Array constructors can be used to combine multiple outputs into a single array.

JPL expects all expressions (including (native) function calls) to be idempotent, meaning that an expression should always produce the same output when run in the same environment (input, variables).
This allows subsequent calls to be cached which has large impact on the overall performance.

Indices are zero based, meaning that the first item in an array has index `0`, the second item has index `1`, and so on.

## Pipe operators

- `. | .`: `{expression} | {expression}`

The pipe operator separates two expressions.
Each output from the expression on the left side is fed into the expression on the right side.

## Variables

Variables are referenced by variable identifiers.
A variable can be one of the following:

- `.` (identity selector)
- `var1_name`: `[a-zA-Z_$][a-zA-Z_$\d]*` (variable selector)

- Note that variable and function declarations in JPL use scopes.
- Special language reserved verbs including `true`, `false`, `null` can not be used as variable names.

## Variable definitions

Variables can be defined by using the following syntax:

- `var1_name = .`: `{variable-selector} = {expression}`

- The expression used to determine the value is executed in its own scope.
- A variable definition returns its input as its output, not the variable.

## Identity selector

The identity selector references the input of the current expression.

## Variable identifiers

- Variable identifiers have the following form: `.|([a-zA-Z_$][a-zA-Z_$\d]*)`
- A special form is the identity selector `.`, where duplicate `.` are condensed into a single `.`, e.g. when using a field access operator: `.a`, which accesses field `a` from the identity variable, instead of `..a`

- Variable identifiers MUST NOT start with `\d`, as this would collide with the number syntax.

## Object field access operator

- `var1_name.some_field`: `{expression -> object}.{field-name}`
- `var1_name["special% field"]`: `{expression -> object}[{expression -> string}]` (generic index)

- Field names may only include simple names in the form of `[a-zA-Z_$][a-zA-Z_$\d]*`. The generic index form is to be used for complex or dynamic field names.
- The expression used in the generic index form to determine the field name is executed in its own scope.

## Array / String field access operator

- `var1_name[1]`: `{expression -> array}[{expression -> number}]` (generic index)

- Negative indices can be used to reference values from the end of the array, -1 referring to the last element, -2 referring to the next to last element, and so on.
- The expression used in the generic index form to determine the index is executed in its own scope.
- When used on strings, the operator returns the unicode character at the specified index.

## Array / String slice

- `var1_name[1:2]`: `{expression -> array|string}[{expression -> number}:{expression -> number}]`

- Used to return a subarray of an array or substring of a string. The array returned by `.[10:15]` will be of length 5, containing the elements from index 10 (inclusive) to index 15 (exclusive). Either index may be negative (in which case it counts backwards from the end of the array) or omitted (in which case it refers to the start or end of the array). Omitting an index is equivalent to setting it to `null`.
- The expressions used to determine the start and end indices are executed in their own scopes.

## Array / String / Object value iterator

- `var1_name[]`: `{expression -> array|object}[]`

- Used to return all entries of an array or string.
- Can also be used on objects to return all values of the object.
- The values are returned as separate outputs.

## Optional field access operators

- `var1_name.some_field?`: `{object-field-access}?`
- `var1_name[1]?`: `{generic-field-access}?`

- Used to omit errors that occur when the input does not have the expected data type.

## Output concatenation operator

- `., .`, `{expression}, {expression}`

- Used for concatenating the outputs of two expressions.
- Each expression creates its own scope.

## Grouping operator

- `(.)`: `({expression})`

- Used for logically grouping expressions together.

## String literals

- `"some text with \"quotes\""`, `'some text with "quotes"'`
- either `"` or `'` (single quote) as boundaries
- `\` (backslash) as escape sequence, e.g. `\'` => `'`, `\\` => `\` (see [String escape sequences](#string-escape-sequences))
- Strings may only contain characters with a char code `>= 0x20`, excluding `CR`, `LF` and `TAB` (among others). [Escape sequences](#string-escape-sequences) can be used as a replacement.

## Multiline strings

- ```JPL
  `a string which can contain
  literal newlines and tabs`
  ```

- `` ` `` (backtick) as boundaries
- Multiline strings are similar to regular string literals, but are also allowed to include `CR`, `LF` and `TAB`.
- In multiline strings, newline symbols (`CR`, `LF` or `CRLF`) can be escaped to ignore them.

  ```JPL
  `a multiline string \
  without newline`
  ```

## String interpolation

- `"this \(.) is interpolated"`: `"\({expression})"`

- Whatever the expression returns will be interpolated into the string
- Can also be used in single quoted and multiline strings

## String escape sequences

The following sequences are valid inside string literals:

|   HEX Code    |     Escape sequence      | Char                                              |
| :-----------: | :----------------------: | :------------------------------------------------ |
|     0x08      |           `\b`           | `BS` (backspace)                                  |
|     0x09      |           `\t`           | `TAB`                                             |
|     0x0a      |           `\n`           | `LF` (line feed)                                  |
|     0x0c      |           `\f`           | `FF` (form feed)                                  |
|     0x0d      |           `\r`           | `CR` (carriage return)                            |
|     0x22      |           `\"`           | `"`                                               |
|     0x27      |           `\'`           | `'`                                               |
|     0x2f      |           `\/`           | `/` (forward slash)                               |
|     0x5c      |           `\\`           | `\` (backslash)                                   |
|     0x60      |         `` \` ``         | `` ` `` (backtick)                                |
| 0x00 - 0xFFFF | `\u<hex><hex><hex><hex>` | Unicode sequence, e.g. `\u0040` equals `@`        |
|               |     `\(expression)`      | See [String interpolation](#string-interpolation) |

## Numbers

- `-3.7`: `-?\d+(\.\d*)?([eE][-+]?\d+)?`
- Numbers MUST NOT start with a decimal point, as this would collide with the identity selector

## Booleans

- `true`
- `false`

## Null

- `null`

## Object constructors

- `{ key: . }`: `{ {field-name}: {expression}, ... }`
- `{ key }`: `{ {field-name}, ... }`
- `{ key? }`: `{ {field-name}?, ... }`
- `{ "key": . }`: `{ {string}: {expression}, ... }`
- `{ (.key): . }`: `{ ({expression}): {expression}, ... }`
- `{ (.key)?: . }`: `{ ({expression})?: {expression}, ... }`

- Used to generate JSON objects.
- The expressions used for determining a field's key and value are executed in their own scopes.
- If an expression for determining a field's key or value returns multiple results, an object for each output is generated.
- If no value is specified for a field name, the value is taken from a variable with the same name. In this case, `?` can be appended to omit errors that occur if the variable does not exist.
- Keys can also be specified with quotes, enabling more complex names, escape sequences and string interpolation.
- When a key is surrounded by parentheses, an expression can be used. In this case, `?` can be appended to omit errors that occur when the result of the expression is not a string.
- Objects are unordered.

## Array constructors

- `[ . ]`: `[ {expression} ]`

- Used to generate JSON arrays.
- The expression used for determining the values is executed in its own scope.
- If the expression for determining the value returns multiple results, one single array is generated including all outputs.

## Mathematical operations

### Addition

- `1 + 2`: `{expression} + {expression}`

- Numbers are added by normal arithmetic
- Arrays are added by being concatenated into a larger array
- Strings are added by being joined into a larger string
- Objects are added by merging, that is, inserting all key-value pairs from both objects into a single combined object. If both objects contain a value for the same key, the object on the right wins. (For recursive merge use the `*` operator)

- `null` can be added to any value, and returns the other value unchanged.

### Subtraction

- `1 - 2`: `{expression} - {expression}`

- Numbers are subtracted by normal arithmetic
- Subtracting one array from another removes all occurrences of the second array's elements from the first array
- Subtracting one string from another produces the first string, stripping out all occurrences of the second string
- Subtracting a string from an object removes the entry from the object where the key matches the string.
- Subtracting an array from an object removes all entries from the object where the value occurs in the array.

### Subtraction, Multiplication, Division, Remainder

- `1 * 2`: `{expression} * {expression}`
- `1 / 2`: `{expression} / {expression}`
- `1 % 2`: `{expression} % {expression}`

- Numbers are calculated by normal arithmetic
- Multiplying a string by a number produces the concatenation of that string that many times, `"x" * 0` produces `null`.
- Multiplying two objects will merge them recursively. This works like addition but if both objects contain a value for the same key, and the values are objects, the two are merged with the same strategy.
- Dividing a string by another splits the first using the second as separators.
- Division by zero raises an error

## Conditionals

- `==`, `!=`

- `==` returns `true`, if both values are equal in terms of type and value (like JavaScript's `===` comparator).
- `!=` is "not equal" and returns the opposite of `==`.

## Null coalescence

- `. ?? true`: `{expression} ?? {expression}`

- Returns the first value if it is not `null`, otherwise the second value.

## Assignment

- `.a.b = 1`: `{expression}.path = {expression}` (set)
- `.a.b |= . + 1`: `{expression}.path |= {expression}` (update)
- `.a.b += 1`: `{expression}.path += {expression}`
- `.a.b -= 1`: `{expression}.path -= {expression}`
- `.a.b *= 1`: `{expression}.path *= {expression}`
- `.a.b /= 1`: `{expression}.path /= {expression}`
- `.a.b %= 1`: `{expression}.path %= {expression}`
- `.a.b ?= 1`: `{expression}.path ?= {expression}` (null coalescence)
- `.a[1:2][] |= . * .`

- Assignment in JPL works different than in most other programming languages. JPL is a truely immutable language and has no concept of references. Thus, when assigning a value to another, a new value is created without modifying the existing ones.
- All values that can be inferred from the specified path are updated like specified by the operator. When the path targets multiple values, e.g. when using a value iterator or slice operator, each individual value is updated sequentially to create one big resulting value.
- When the right operand produces multiple outputs, for each output a new value is created. If the path targets multiple values, that means that in this situation there will be produced multiple results for every each targeted value.
- The update operator `|=` differs from the set operator (`=`) in that the value that is to be updated is used as the right operands input. This means, that for instance `{a:1,b:3} | .a = .b` produces `{a:3,b:3}`, whereas `{a:1,b:3} | .a |= .b` produces an error, because the program attempts to access `.b` from the value of `a`, which is a number.
- The path is part of the operator and must not be grouped. This however allows more complex application of the assignment operators, for example `(.a.b).c.d = 1` can be used to extract the value at `.a.b` and then run an assignment `.c.d = 1` on the result. This is equivalent to `.a.b | .c.d = 1`, but allows for accessing the whole input in the right operand, like e.g. `(.a.b).c.d = .a`.

## Variable assignment

`a.b = 1`: `{variable-selector}.path = {expression}`

- All assignment operators can also be used for variable assignment. `a.b = 1` is a shorthand for `a = (a).b = 1`
- Like with simple variable definitions, the output of the variable assignment operator is its input, not the variable.

## If

- `if A then B end`
- `if A then B else C end`
- `if A then B elif C then D else E end`

- Will act like `B` if `A` produces a value other than `false` or `null`, but act like `C` otherwise.
- `if A then B end` is the same as `if A then B else . end`.
- If the condition `A` produces multiple results, then `B` is evaluated once for each result that is not `false` or `null`, and `C` is evaluated once for each `false` or `null`.

## Comparison operators

- `>`, `>=`, `<=`, `<`

- Used to return wether the left argument is greater than, greater than or equal, less than or equal to or less than their right argument (respectively).
- Comparison is done in the following order (small to large):
- `null`
- functions
- `false`
- `true`
- numbers
- strings, in alphabetical order (by unicode codepoint value)
- arrays, in lexical order
- objects

- Objects are compared using the following rules: first they're compared by comparing their sets of keys (as arrays in sorted order), and if their keys are equal then the values are compared key by key

## Boolean operators

- `and`, `or`, `not`

- The standard of truth is the same as with if expressions (`false` and `null` are considered as "false values", anything else is a "true value")
- If an operand of one of these operators produces multiple results, the operator itself will produce a result for each input.

## Error handling

- `try .test`: `try {expression}`
- `try .test catch .` : `try {expression} catch {expression}`

- Used to catch any errors occuring in the specified expression
- If `catch` is omitted, no outputs are returned if an error is caught

## Error suppression

- `true + 1 ?`: `{expression}?`

- The `?` operator is shorthand for `try {expression}`

## Function definitions

- `func (a, b): a + b`: `func (arg1, arg2, ...): {expression}` (anonymous function definition)
- `func add(a, b): a + b`: `{variable-selector} = func {variable-selector}(arg1, arg2, ...): {expression}` (named function definition)
- `add = func (a, b): a + b`: `{variable-selector} = func (arg1, arg2, ...): {expression}`

- Even when no arguments are used, the parentheses must be specified
- A function is executed in its own scope
- Functions are in fact variables and thus can also be passed as arguments into other functions
- A named function definition is in fact a shorthand for an anonymous function definition in combination with a variable definition, thus a named function definition returns its input as its output, not the function

## Function calls

- `add(1, 2)`: `{variable-selector}({expression}, {expression})`
- `add->(input, 1, 2)`: `{variable-selector}->({expression}, {expression}, {expression})` (bound function call)

- Even when no arguments are provided, the parentheses must be added to the statement
- The expressions used to determine the function's arguments run in their own scopes
- When calling anonymous functions, the function definition must be wrapped in parentheses, e.g. `(func (x): x + 1)(1)`
- Like with all other types, a variable can reference multiple functions at once, in which case all functions are called and all resulting outputs are being concatenated
- When prefixing the argument array with `->`, the leftmost argument is bound to the function input, and the remaining arguments become the function's arguments. The following example returns `true`: `func plus(x): . + x | (1 | plus(2)) == plus->(1, 2)`

## Comments

- `/* comment */`
- `# line comment`

## Further reading

- [Builtin functions](builtins.md)
- [Libraries](libraries.md)
- [Operator priority](order.md)
