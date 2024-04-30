# [JPL - JSON Processing Language](index.md) - Builtin functions

## `map(f)`

Maps the input to an array by calling `f` for each field.
`map(f)` is equivalent to `[.[] | f()]`. In fact, this is how it is defined.

Example: `[1, 2, 3] | map(func(): .+1)` -> `[2, 3, 4]`

## `mapValues(f)`

Maps the input's values by calling `f` for each field.
This differs from `map` in that it does not change the type of the input. However, other than with `map`, if `f` returns multiple values, also multiple results are returned.
`mapValues(f)` is defined as `.[] |= f()`.

Example: `{ a: 1, b: 2 } | mapValues(func(): .+1)` -> `{ "a": 2, "b": 3 }`

## `select(f)`

Only returns the input if `f` returns a truthy value.

Example: `1, true | select(func(): . == 1)` -> `1`

## `reduce(f, initialValue)`

Reduces the input array by recursively calling `f` for each field. `f` takes each field as its input and two arguments: The first is the cumulated sum and the second is the current field's index.

Example: `[1, 2, 3] | reduce(func(sum): sum+., 0)` -> `6`

## `while(cond, f)`

Calls `f` while `cond` returns a truthy value and returns all results.
Both `cond` and `f` take the result of the previous iteration as their input.

Example: `1 | while(func(): .<=3, func(): .+1)` -> `1, 2, 3`

## `until(cond, f)`

Calls `f` until `cond` returns a truthy value and returns the single final result. Both `cond` and `f` take the result of the previous iteration as their input.

Example: `1 | until(func(): .>=3, func(): .+1)` -> `3`

## `range(from, to, step ?? 1)`

Returns a range of numbers in the half-open interval `[from, to)` with an increment of `step` (which defaults to `1`).

Example: `range(1, 10)` -> `1, 2, 3, 4, 5, 6, 7, 8, 9`

## `toEntries()`

Returns an array of all entries of the input object or array. Each entry is represented as an object with a `key` field and a `value` field. The `key` is a string for objects and a number for arrays.

Example: `{ a: 1, b: 2 } | toEntries()` -> `[{ "key": "a", "value": 1 }, { "key": "b", "value": 2 }]`

## `fromEntries()`

Creates an object based on the input array. Each entry of the array is expected to be an object with a key and a value. The function accepts key, Key, name, Name, value and Value as keys.

Example: `[{ key: "a", value: 1 }, { key: "b", value: 2 }] | fromEntries()` -> `{ "a": 1, "b": 2 }`

## `withEntries(f)`

`withEntries(f)` is a shorthand for `toEntries() | map(f) | fromEntries()`. It is useful for doing some operation to all keys and values of an object.

Example: `{ a: 1, b: 2} | withEntries(func(): .key |= "_"+.)` -> `{ "_a": 1, "_b": 2 }`

## `add()`

Adds together all values of the input array. This might mean sum, concatenate or merge depending on the types of the elements of the input array - the rules are the same as those for the `+` operator.

If the input is an empty array, `add` produces `null`.

Example: `[1, 2, 3], ["a", "b", "c"] | add()` -> `6, "abc"`

## `join(sep)`

Joins the array of elements given as input, using `sep` as separator.

Example: `"abc" / "" | join(", ")` -> `"a, b, c"`

## `sort()`, `sortBy(f)`

The sort functions sorts its input, which must be an array. Values are sorted in the following order:

- `null`
- `false`
- `true`
- numbers
- strings, in alphabetical order (by unicode codepoint value)
- arrays, in lexical order
- objects

The ordering for objects is a little complex: first they're compared by comparing their sets of keys (as arrays in sorted order), and if their keys are equal then the values are compared key by key.

`sortBy(f)` compares two elements by comparing the result of `f` on each element.

Example: `[2, 1] | sort()` -> `[1, 2]`

## `group()`, `groupBy(f)`

The group function groups its input array into separate arrays, and produces all of these arrays as elements of a larger sorted array.

`groupBy(f)` groups the elements by the result of `f` on each element.

Example: `{ a: 2, b: 1, c: 2 } | toEntries() | groupBy(func(): .value)` -> `[[{ "key": "b", "value": 1 }], [{ "key": "a", "value": 2 }, { "key": "c", "value": 2 }]]`

## `unique()`, `uniqueBy(f)`

Removes duplicate values from the input array. The resulting array is sorted.

`uniqueBy(f)` determines which elements are duplicates by comparing the results of `f` on each element.

Example: `[1, 1, 2, 3] | unique()` -> `[1, 2, 3]`

## `recurse(cond ?? func(): .!=null)`, `recurseBy(f, cond ?? func(): .!=null)`

The `recurse` function allows you to search through a recursive structure, and extract interesting data from all levels.

Example: `{ name: "/", files: [{ name: "/a", files: [{ name: "/a/a.txt" }, null] }, { name: "/b.txt" }] } | recurseBy(func(): .files?[]?) | .name` -> `"/", "/a", "/a/a.txt", "/b.txt"`

## `reverse()`

Reverses the order of the input array.

Example: `[1, 3, 2] | reverse()` | `[2, 3, 1]`

## `min()`, `minBy(f)`, `max()`,`maxBy(f)`

Returns the smallest (`min`) or largest (`max`) value of the input array.

`minBy(f)` and `maxBy(f)` select the result based on the result of `f` on each element.

Example: `[3, 2, 4] | min()` -> `2`

## `first(f)`

Returns the first result of `f`.

Example: `first(func(): range(0, 3))` -> `0`

## `nth(which, f)`

Returns the nth result of `f`, according to `which`. If `which` is a function, `which` is called with the number of results as its input.

Example: `nth(1, func(): range(0, 3))` -> `1`

## `last(f)`

Returns the last result of `f`.

Example: `last(func(): range(0, 3))` -> `2`

## `isEmpty(f)`

Returns `true` if `f` does not return any results, `false` otherwise.

Example: `isEmpty(func(): (true | strings()))` -> `true`

## `all(cond)`, `allBy(f, cond)`

Returns `true` if `cond` returns a truthy value for all entries of the input array, `false` otherwise.

`allBy(f, cond)` checks the results of `f` instead of the input.

Example: `[1, 2, 3] | all(func(): . > 0)` -> `true`

## `any(cond)`, `anyBy(f, cond)`

Returns `true` if `cond` returns a truthy value for at least one entry of the input array, `false` otherwise.

`anyBy(f, cond)` checks the results of `f` instead of the input.

Example: `[1, 2, 3] | any(func(): . > 1)` -> `true`

## `startsWith(token)`

Returns `true` if the input string starts with `token`, `false` otherwise.

Example: `"Hello, World" | startsWith("Hello")` -> `true`

## `endsWith(token)`

Returns `true` if the input string ends with `token`, `false` otherwise.

Example: `"Hello, World" | endsWith("World")` -> `true`

## `contains(token)`

Returns `true` if the input string, array or object contains `token`, `false` otherwise.
A token is contained in a string if it is a substring.
It is contained in an array or object, if it equals one of the entries (or values in case of an object).

Example: `"Hello, World" | contains(",")` -> `true`

## `trim()`

Trims whitespace surrounding the input string.

Example: `" Hello, World\n" | trim()` -> `"Hello, World"`

## `trimStart()`

Trims whitespace from the start of the input string.

Example: `" Hello, World\n" | trimStart()` -> `"Hello, World\n"`

## `trimEnd()`

Trims whitespace from the end of the input string.

Example: `" Hello, World\n" | trimEnd()` -> `" Hello, World"`

## `toNumber()`

Parses the input string as a number.

Example: `"1.5" | toNumber()` -> `1.5`

## `toString()`

Converts the input to a string. Strings are returned unchanged, whereas all other types are converted to JSON.

Example: `"string", { a: 1, b: 2 } | toString()` -> `"string", "{\"a\":1,\"b\":2}"`

## `toJSON()`

Converts the input to a JSON string. This differs from `toString` in that strings are encoded to JSON strings.

Example: `"string", { a: 1, b: 2 } | toJSON()` -> `"\"string\"", "{\"a\":1,\"b\":2}"`

## `fromJSON()`

Parses the input string as JSON.

Example: `"{\"a\":1,\"b\":2}" | fromJSON()` -> `{ "a": 1, "b": 2 }`

## `has(key)`

Returns `true` if the input has a field for specified key, `false` otherwise. Arrays and objects are supported.

Example: `{ a: 1 } | has("a")` -> `true`

## `in(target)`

Returns `true` if the specified target has a field for the input as its key, `false` otherwise. Arrays and objects are supported.

Example: `"a" | in({ a: 1 })` -> `true`

## `keys()`

Returns all keys or indices of the provided input. Arrays and objects are supported.

Example: `{ a: 1, b: 2 } | keys()` -> `["a", "b"]`

## `length()`

Returns the length of the provided input. Arrays, strings and objects are supported. For `null`, `0` is returned.

Example: `[1, 2, 3] | length()` -> `3`

## `type()`

Returns the type of the provided input.

Example: `null, func():., true, 1, "string", [], {} | type()` -> `"null", "function", "boolean", "number", "string", "array", "object"`

## `error()`

Throws a new error with the provided input as message. The message may be of any type.

Example: `try ("test" | error()) catch .` -> `"test"`

## `void()`

Returns nothing.

Example: `[1, 2, 3] | map(void)` -> `[]`

## `hasContent()`

Returns `true` if the provided input is neither `null`, a function, nor an empty array, object or string, `false` otherwise.

## `now()`

Returns the number of milliseconds elapsed since midnight, January 1, 1970 Universal Coordinated Time (UTC).

## Type selectors

The following type selectors return only those inputs that match specific types.

- `arrays()` - arrays
- `objects()` - objects
- `iterables()` - arrays and objects
- `scalars()` - anything but arrays and objects
- `booleans()` - booleans
- `numbers()` - numbers
- `strings()` - strings
- `nulls()` - null
- `functions()` - functions
- `nullLikes()` - null and functions
- `values()` - anything but null and functions
- `contents()` - anything but null, functions, empty arrays, empty objects and empty strings

Example: `[1, true, "test", {}] | map(func(): (numbers, strings)())` -> `[1, "test"]`

## Mathematical functions

### `pow(exp)`

Raises the numeric input by the power of `exp`.

Example: `4 | pow(2)` -> `16`

### `sqrt()`

Returns the square route of the numeric input.

Example: `4 | sqrt()` -> `2`

### `exp()`

Returns e (the base of natural logarithms) raised to the power of the numeric input.

Example: `0 | exp()` -> `1`

### `log()`

Returns the natural logarithm (base e) of the numeric input.

Example: `1 | log()` -> `0`

### `log10()`

Returns the base 10 logarithm of the numeric input.

Example: `10 | log10()` -> `1`

### `sin()`

Returns the sine of the numeric input.

Example: `0 | sin()` -> `0`

### `cos()`

Returns the cosine of the numeric input.

Example: `0 | cos()` -> `1`

### `tan()`

Returns the tangent of the numeric input.

Example: `0 | tan()` -> `0`

### `asin()`

Returns the arcsine of the numeric input.

Example: `0 | asin()` -> `0`

### `acos()`

Returns the arc cosine of the numeric input.

Example: `1 | acos()` -> `0`

### `atan()`

Returns the arctangent of the numeric input.

Example: `0 | atan()` -> `0`

### `ceil()`

Returns the smallest integer greater than or equal to the numeric input.

Example: `1.5, -1.5 | ceil()` -> `2, -1`

### `floor()`

Returns the greatest integer less than or equal to the numeric input.

Example: `1.5, -1.5 | floor()` -> `1, -2`

### `round()`

Returns the numeric input rounded to the nearest integer.

Example: `1.5, -1.5 | round()` -> `2, -1`

### `trunc()`

Returns the integral part of the numeric input, removing any fractional digits.

Example: `1.5, -1.5 | trunc()` -> `1, -1`

### `abs()`

Returns the absolute value of the numeric input (the value without regard to whether it is positive or negative).

Example: `-1 | abs()` -> `1`

## Date and time

See [Date and time](builtins-date.md)
