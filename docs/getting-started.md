# [JPL - JSON Processing Language](index.md) - Getting started

In JPL, each expression has exactly one input and produces any number of outputs, each of which in turn is fed individually into the next expression.
The current expression's input can be referred to by the identity selector `.`.

With this we can create the most basic JPL program...

```jpl
.
```

..., which simply outputs unchanged whatever was fed into the program (although admittedly this may not be the most useful one).

## Comments

There are two types of comments in JPL, line comments and block comments.

### Line comments

Line comments start with `#`, the rest of the line is a comment.

```jpl
# Full line comment
. # Line comment after code
```

### Block comments

Block comments can appear anywhere where whitespace is allowed. They start with `/*` and end with `*/` and can span multiple lines.

```jpl
/*
This is a block comment.
It can take up multiple lines.
*/
func add(summand): /* a comment in a line of code */ (
  . + summand
)
```

## Constants

Defining constant values in JPL is the same as writing JSON documents.
In fact, any valid JSON is also a valid JPL constant.

```jpl
# Booleans
true
false

# Null
null

# Strings
"a string"

# Numbers
1
-2
0.5
3.5e-3

# Objects
{ "name": "JPL", "longName": "JSON Processing Language" }

# Arrays
["value1", 2, true]
```

JPL, like JSON, does not distinguish between integers and floating point numbers.

## Strings

Strings are usually wrapped in double quotes `"` as boundaries, but you can also use single quotes `'` instead.

If a string contains its own boundaries, you must either [escape them](spec.md#string-escape-sequences), or switch to alternative boundaries.
Tabs and any newline symbols also need to be escaped, unless you are using [multiline strings](#multiline-strings).

```jpl
"string"
'also a string'
"string containing \"double quotes\" and 'single quotes'."
'string containing "double quotes" and \'single quotes\'.'
```

### Multiline strings

Multiline strings use backticks `` ` `` as their boundaries. Unlike with regular strings, you can use tabs and newline symbols without escaping them.

```jpl
`A string which contains a
newline symbol`
# -> "A string which contains a\nnewline symbol"
```

You can escape newline symbols in your multiline strings to ignore them.

```jpl
`A string which contains no \
newline symbols`
# -> "A string which contains no newline symbols"
```

### String interpolation

String interpolation enables you to put the result of a JPL expression into your string. A string interpolation starts with `\(` followed by any JPL expression and ends with `)`.
If the expression evaluates to multiple results, a string is created for each result.
Conversely, if the expression has no outputs, there is also no string created.

```jpl
"3 times 9 is \(3 * 9)"
# -> "3 times 9 is 27"

"multiple \("results", "strings")"
# -> "multiple results", "multiple strings"
```

## Mathematical operations

JPL provides a subset of mathematical operations.

```jpl
# Addition
3 + 2
# -> 5

# Subtraction
3 - 2
# -> 1

# Multiplication
3 * 2
# -> 6

# Division
3 / 2
# -> 1.5

# Remainder
3 % 2
# -> 1
```

## Conditions

The following expressions all result in `true`.

```jpl
1 == 1
1 != 2
1 < 2
3 > 2
1 <= 2
2 >= 2

not false

1 > 0 and -0 == 0
false or true
```

## Null coalescence

The null coalescing operator `??` can be used to provide default values.
If the left expression evaluates to `null`, the right expression's outputs are returned instead.

```jpl
1 ?? "default"
# -> 1

null ?? "default"
# -> "default"
```

The null coalescing operator can be combined with the `void()` builtin to suppress the output of `null`.

```jpl
1 ?? void()
# -> 1

null ?? void()
# -> <nothing>
```

There are also the builtins `isContent()` and `contents()`, which also omit other empty values than `null`.

## Combining expressions

In the previous examples, each line represented a single expression.
In a real program, however, all expressions are generally chained together.

### Pipes

The most common way to combine expressions is by using the pipe operator `|`, which takes the individual outputs of the left expression and passes them each as input to the right expression.

```jpl
1 + 2 | "the result is: \(.)"
# -> "the result is: 3"
```

### Output concatenation

To be able to produce multiple outputs, you can use the output concatenation operator `,`, which combines the outputs of both the left and right expression.

```jpl
1, 2 | "the result is: \(.)"
# -> "the result is: 1", "the result is: 2"
```

Note that `,` has a higher execution priority than `|`, so in the example above, `1, 2` is executed first and then the individual results are passed into `"the result is: \(.)"`.

### Grouping

You can group any expressions in `()` to change their execution order.

```jpl
1, (2 | "the result is: \(.)")
# -> 1, "the result is: 2"
```

## Variables

A variable can store an arbitrary value to be referenced later.

```jpl
myAge = 99 | "I am \(myAge) years old"
# -> "I am 99 years old"
```

Variables (as JPL in general) are immutable, which means that their value can't be changed after declaring it.
When redefining a variable, the new value only takes effect for subsequent expressions in the same scope.

In the following example, `a = 1` and `a = 2` have their own individual scopes (because they are combined with `,`), so when referring to `a` afterwards, the variable is not available anymore and the program fails.

```jpl
(a = 1, a = 2) | a
# ! JPLRuntimeError: ReferenceError - a is not defined
```

## Error handling

### Creating errors

If something goes wrong in a JPL program, an error is thrown.
You can throw an error yourself by using the `error()` builtin.

```jpl
error->("something went wrong")
```

There are two ways of handling errors in JPL, `try catch` and error suppression `?`.

### Try catch

When using `try` and `catch`, the expressions after `try` are executed and if they produce an error, the expression after `catch` is executed with the error as input.

```jpl
try
  error->("something went wrong")
catch
  "Oops! " + .
# -> "Oops! something went wrong"
```

### Error suppression

Error suppression can be used to suppress any errors. If the expression preceeding the `?` throws an error, the error is omitted and no outputs are returned.

```jpl
1 + "a"
# ! JPLRuntimeError: TypeError - number (1) and string ("a") cannot be added together

1 + "a" ?
# -> <nothing>
```

## If clauses

You can run code conditionally by using an `if` clause, which looks like `if A then B else C end`.
If the condition `A` evaluates to a truthy value (anything but `false` or `null`), `B` is executed, `C` otherwise.

```jpl
if 2 > 1
then
  "it is more"
else
  "it is less"
end
# -> "it is more"
```

`if A then B end` is a shorthand form of `if A then B else . end`.

If you want to test multiple conditions consecutively, you can add additional clauses with `elif`, which are tested if the previous conditions did not apply.

```jpl
2 |
if . == 1
then
  "one"
elif . == 2
then
  "two"
else
  toString()
end
# -> "two"
```

## Functions

Functions are declared using the `func` word.

There are two types of functions, named functions and anonymous functions.

### Named functions

Named functions are accessible using their name, similar to a variable.

```jpl
# Here, we define a named function called "greet"
func greet(): (
  "Hello!"
) |
# Now, we can call the function by its name
greet()
# -> "Hello!"
```

### Anonymous functions

Anonymous functions are returned as the expression's output instead, so you can access it afterwards using the identity operator `.`.

```jpl
# Here, we define an anonymous function
func (): (
  "Hello!"
) |
# We can access the function with the identity operator
.()
# -> "Hello!"
```

### Function arguments

A function can take an arbitrary number of arguments, which are available in the function body by their name.

```jpl
# This function takes two arguments, "name" and "surname"
func greet(name, surname): (
  "Hello, \(name) \(surname)!"
) |
greet("John", "Doe")
# -> "Hello, John Doe!"
```

### Function input

When calling a function, the current expression's input is passed into the function.

```jpl
func greet(): (
  "Hello, \(.)!"
) |
"John" | greet()
# -> "Hello, John!"
```

This means, that calling an anonymous function by accessing it from the identity operator passes the function into itself. The following program results in an endless loop.

```jpl
func(): (
  .()
) | .()
```

It is also possible to provide a function's input as its first argument.
This is done by using a bound function call `fn->(input, arg1, ...)`.

```jpl
func greet(surname): (
  "Hello, \(.) \(surname)!"
) |
greet->("John", "Doe")
# -> "Hello, John Doe!"
```

### Some more details about functions

A named function definition can also be written by putting an anonymous function into a variable.

```jpl
greet = func(): (
  "Hello!"
)
```

Functions can be passed as arguments and can be put into variables or JSON structures.
However, they are omitted from the program's result.
Comparing two functions returns `true`.

Even though the parentheses surrounding the function body are optional, it is recommended to include them for better readability.
If a function contains pipes `|` or output concatenation `,`, the parentheses are required.

Because function arguments are separated by `,`, you cannot directly use output concatenation `,` when calling a function.
Instead, the affected argument must be wrapped in parentheses.

```jpl
func greet(name, surname): (
  "Hello, \(name) \(surname)!"
) |
greet(("John", "Jane"), "Doe")
# -> "Hello, John Doe!", "Hello, Jane Doe!"
```

### Creating arrays

Like previously mentioned, arrays can be created by writing JSON.

```jpl
["Hello", "JPL"]
```

However, you can use any JPL expression inside of the array brackets.
All outputs are then combined into an array.

```jpl
["The result is \(1, 2, 3)"]
# -> ["The result is 1", "The result is 2", "The result is 3"]
```

Arrays can be used to handle all results of an expression at once.
In the following example, the text "no result" is returned if the expression in the array did not return any outputs.

```jpl
[
  # `numbers()` only returns numbers, so `true` is omitted
  true | numbers()
] | if . == [] then "no result" else .[] end
# -> "no result"
```

### Creating objects

Objects, like arrays, can be created by writing JSON.

```jpl
{
  "key": "value",
  "key2": 1
}
# -> { "key": "value", "key2": 1 }
```

The value part can be any JPL expression, but if it contains output concatenation `,`, it needs to be wrapped inside parentheses, as the comma is already used to separate object entries.
If the expression returns multiple outputs, an object is created for each output.
Conversely, if the expression returns no outputs, no objects are created.

```jpl
{
  "id": (1, 2),
  "time": now()
}
# -> { "id": 1, "time": 1714500000000 }, { "id": 2, "time": 1714500000000 }
```

### Object keys

If your object keys do not include any special symbols, you can omit the quotes.

```jpl
{ id: 1 }
# -> { "id": 1 }
```

You can also use any JPL expression to resolve a object key by wrapping the key in parentheses. However, the expression must only return strings, otherwise an error is thrown. You can use the `toString()` builtin to convert non string values.
If the expression returns multiple outputs, an object is created for each output.
Conversely, if the expression returns no outputs, no objects are created.

```jpl
{
  ("a", "b"): 1
}
# { "a": 1 }, { "b": 1 }
```

As it is a common case to put the value of a variable into an object with the variable's name as key, there is a shorthand for this.

```jpl
hello = "JPL" |
{ hello }
# -> { "hello": "JPL" }
```

## Accessing object entries

The form `object.name` extracts the value of `object` at the key `name`.
If the key contains special symbols, you need to use the generic index form `object["name"]` (or `.["name"]` for identity access).
This form can take any JPL expression, but must only return strings, otherwise an error is thrown. You can use the `toString()` builtin to convert non string values.

```jpl
variable = { a: 1 } | variable.a
# -> 1

{ a: 1 } | .a
# -> 1

{
  a: {
    b: 1
  }
} | .a.b
# -> 1

{
  a: {
    "123": 1
  }
} | .a["123"]
# -> 1

{
  "key with spaces": 1,
  "/^1": 2
} | .["key with spaces", "/^1"]
# -> 1, 2
```

## Accessing array entries

Entries of arrays can be accessed with `array[index]` (or `.[index]` for identity access).
`index` can be any JPL expression, but must only return numbers, otherwise an error is thrown. You can use the `toNumber()` builtin to parse a string as a number.

```jpl
["a", "b", "c"] | .[1]
# -> "b"

{
  a: [1, 2, 3]
} | .a[2]
# -> 3
```

## Selecting string characters

Selecting single unicode characters from a string works the same as accessing an array's index.

```jpl
"hello, world" | .[7, 11]
# -> "w", "d"
```

## Slicing arrays and strings

To extract a subset of an array or string, use the form `value[from:to]`.
`from` and `to` mark the inclusive start and exclusive end points for the slice.

```jpl
["a", "b", "c", "d"] | .[1:3]
# -> ["b", "c"]

"the quick brown fox jumps over the lazy dog" | .[4:(9, 15)]
# -> "quick", "quick brown"
```

If you omit one of the points (or set them to `null`), respectively the start or end of the input is selected.

```jpl
["an", "array", "of", "values"] | .[1:]
# -> ["array", "of", "values"]

"string-with-overhead"[:6]
# -> "string"
```

You can specify negative indices to look back from the end of the input.

```jpl
["a", "b", "c", "d"][-2:]
# -> ["c", "d"]

"cut off the last word"[:-4]
# -> "cut off the last "
```

## Iterating values

You can return all entries from an array or string by using the value iterator `value[]`.
This can also be used on an object to return all of its values.

```jpl
["a", "b", "c"] | .[]
# -> "a", "b", "c"

"string"[]
# -> "s", "t", "r", "i", "n", "g"

{
  a: 1,
  b: 2
} | .[]
# -> 1, 2
```

Array iteration is useful when you want to modify the values of an array.

```jpl
["a", "b", "c", "d"] |
[
  .[1:3][] | "value: \(.)"
]
# -> ["value: b", "value: c"]
```

## Assignments

If you want to modify nested values, you can use assignment operators.

```jpl
{
  a: [
    { b: 1 },
    { c: 2 }
  ]
} |
.a[1].c = 3
# -> { "a": [{ "b": 1 }, { "c": 3 }] }
```

Other than you might have expected, assignment does not change the original value, but instead creates a new one with the specified changes applied to it.

There is a number of assignment operators. Besides the update assignment operator, all operators apply the input of the whole expression to the right expression.
The update operator uses the value that is to be updated as the right expression's input.

```jpl
# Basic assignment - identity is the whole object
{ a: { b: 1 } } | .a.b = 2
# -> { "a": { "b": 2 } }

# Update assignment - identity is .a.b
{ a: { b: 1 } } | .a.b |= if . > 0 then . + 2 else . - 2 end
# -> { "a": { "b": 3 } }

# Mathematical assignment
{ a: { b: 1 } } | .a.b += 1
{ a: { b: 1 } } | .a.b -= 1
{ a: { b: 1 } } | .a.b *= 1
{ a: { b: 1 } } | .a.b /= 1
{ a: { b: 1 } } | .a.b %= 1

# Null coalescing assignment
{ a: { b: 1 } } | .a.c ?= 2
# -> { "a": { "b": 1, "c": 2 } }
```

An assignment can update multiple values at once.

```jpl
{ a: [1, 2, 3, 4] } | .a[] += 1
# -> { "a": [2, 3, 4, 5] }

{ a: [1, 2, 3, 4] } | .a[1:3][] |= if . < 3 then . * 3 else . / 2 end
# -> { "a": [1, 6, 1.5, 4] }
```

It is also possible to apply an assignment to a variable.
In this case, the result of the operation is assigned to the variable instead of returning it as output.
Thus, `variable.a = 1` is equivalent to `variable = (variable | .a = 1)`.

```jpl
value = { a: 1 } | value.a = 2 | value
# -> { "a": 2 }
```

If you don't want to overwrite the variable, you have to wrap it in parentheses or read the variable's value before doing the assignment.

```jpl
value = { a: 1 } | (value).a = 2, value
# -> { "a": 2 }, { "a": 1 }

value = { a: 1 } | value | .a = 2, value
# -> { "a": 2 }, { "a": 1 }
```

This can also be useful if you want to apply an assignment to a part of the input.

```jpl
{ a: { b: 1 } } | (.a).b = 2
# -> { "b": 2 }

{ a: { b: 1 } } | .a | .b = 2
# -> { "b": 2 }
```

## Further reading

- [Specification](spec.md)
- [Builtin functions](builtins.md)
- [Libraries](libraries.md)
- [Operator priority](order.md)
