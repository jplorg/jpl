package builtins

import (
	interpreter "github.com/jplorg/jpl/go/v2/interpreter"
	jpl "github.com/jplorg/jpl/go/v2/jpl"
	library "github.com/jplorg/jpl/go/v2/library"
)

const builtins = `
# There are additional builtins that are implemented directly in the corresponding languages.
# Also the functions in the "internals" namespace are not exposed to the public and are only to be used by the builtins below.

func map(f): ([.[] | f()])
| func mapValues(f): (.[] |= f())
| func select(f): (if f() then . else void() end)
| func reduce(f, s): (
  l = length()
  | func i(i, n, c): (if n < l then i(i, n + 1, f->(.[n], c, n)) else c end)
  | i(i, 0, s)
)
| func while(cond, f): (
  func i(i): (if cond() then ., i->(f(), i) else void() end)
  | i(i)
)
| func until(cond, f): (
  func i(i): (if cond() then . else i->(f(), i) end)
  | i(i)
)
| func range(from, to, step): (
  t = type->(from) | if t != "number" then error->("cannot use \(t) as a number") end
  | t = type->(to) | if t != "number" then error->("cannot use \(t) as a number") end
  | s = (step | if . != 0 and . != null then abs() else 1 end)
  | while->(
    from,
    if from <= to then func (): (. < to) else func (): (. > to) end,
    if from <= to then func (): (. + s) else func (): (. - s) end
  )
)


# type selectors:
| func isArray(): (type() == "array")
| func arrays(): (select(isArray))
| func isObject(): (type() == "object")
| func objects(): (select(isObject))
| func isIterable(): (contains->(["array", "object"], type()))
| func iterables(): (select(isIterable))
| func isScalar(): (not contains->(["array", "object"], type()))
| func scalars(): (select(isScalar))
| func isBoolean(): (type() == "boolean")
| func booleans(): (select(isBoolean))
| func isNumber(): (type() == "number")
| func numbers(): (select(isNumber))
| func isString(): (type() == "string")
| func strings(): (select(isString))
| func isNull(): (type() == "null")
| func nulls(): (select(isNull))
| func isFunction(): (type() == "function")
| func functions(): (select(isFunction))
| func isNullLike(): (contains->(["null", "function"], type()))
| func nullLikes(): (select(isNullLike))
| func isValue(): (not contains->(["null", "function"], type()))
| func values(): (select(isValue))
| func isContent(): (isValue() and not contains->([[], {}, ""], .))
| func contents(): (select(isContent))


| func toEntries(): (o = . | keys() | map(func (): ({ key: ., value: o[.] })))
| func fromEntries(): (
  reduce(
    func (s): (s + { (.key ?? .Key ?? .name ?? .Name): if has("value") then .value else .Value end }),
    {}
  )
)
| func withEntries(f): (toEntries() | map(f) | fromEntries())

| func add(): (reduce(func (sum): (sum + .)))
| func join(sep): (
  reduce(
    func (sum): (
      if sum == null then "" else sum + sep end
      + (. | toString())
    )
  ) ?? ""
)
| func sortBy(f): ([.[] | [[f()], .]] | internals.sortEntries() | [.[][1]])
| func sort(): (sortBy(func (): (.)))
| func groupBy(f): (
  key = 0 | value = 1
  | [.[] | [[f()], [.]]] | internals.sortEntries()
  | if . == []
  then .
  else
    reduce->(
      .[1:],
      func (sum): (if sum[-1][key] == .[key] then (sum)[-1][value] += .[value] else sum + [.] end),
      .[:1]
    )
    | [.[][value]]
  end
)
| func group(): (groupBy(func (): (.)))
| func uniqueBy(f): (groupBy(f) | map(func (): (.[0])))
| func unique(): (group() | map(func (): (.[0])))
| func recurseBy(f, cond): (
  c = cond ?? func (): (. != null)
  | func r(r): (., (f() | select(c) | r(r)))
  | r(r)
)
| func recurse(cond): (recurseBy(func (): ((arrays(), objects()) | .[]), cond))
| func reverse(): ([.[length() - 1 - range(0, length())]])
| func minBy(f): (
  key = 0 | value = 1
  | [.[] | [[f()], .]]
  | reduce->(.[1:], func (sum): (if .[key] < sum[key] then . else sum end), .[0])
  | .[value]
)
| func min(): (reduce->(.[1:], func (sum): (if . < sum then . else sum end), .[0]))
| func maxBy(f): (
  key = 0 | value = 1
  | [.[] | [[f()], .]]
  | reduce->(.[1:], func (sum): (if .[key] > sum[key] then . else sum end), .[0])
  | .[value]
)
| func max(): (reduce->(.[1:], func (sum): (if . > sum then . else sum end), .[0]))
| func first(f): ([f()][0])
| func nth(which, f): (
  [f()]
  | .[if which | type() == "number" then which else length() | which() end]
)
| func last(f): ([f()][-1])
| func isEmpty(f): (first(func (): ((f() | false), true)))
| func allBy(f, cond): (isEmpty(func (): (f() | cond() and void())))
| func all(cond): (allBy(func (): (.[]), cond))
| func anyBy(f, cond): (not isEmpty(func (): (f() | cond() or void())))
| func any(cond): (anyBy(func (): (.[]), cond))

| func getPath(path): (
  fields = (path | if type() == 'string' then . / '.' end)
  | reduce->(fields, func (sum): sum[if type->(sum) == 'array' then toNumber() else toString() end], .)
)
| func updatePath(path, update): (
  fields = (path | if type() == 'string' then (. / '.')[] |= (v = . | try toNumber() catch v) end)
  | reduce->(
    reverse->(fields),
    func (sum): (
      field = .
      | func (): (
        if type() == 'array'
        then (
          input = .
          | try toNumber->(field) catch toString->(field)
          | if type() == 'number'
          then input[.] |= sum()
          else {}[.] |= sum()
          end
        ) elif type() == 'object'
        then (
          .[toString->(field)] |= sum()
        ) elif type->(field) == 'number'
        then (
          [][field] |= sum()
        ) else (
          {}[toString->(field)] |= sum()
        )
        end
      )
    ),
    update
  )->(.)
)
| func setPath(path, value): (
  updatePath(path, func (): value)
)
`

func getBuiltins() map[string]any {
	result := library.CopyMap(native)
	program, err := interpreter.SystemInterpreter.Parse(builtins, &jpl.JPLInterpreterConfig{
		Runtime: jpl.JPLRuntimeOptions{
			Vars: library.MergeMaps(map[string]any{"internals": internals}, native),
			AdjustResult: jpl.JPLScopedPiperFunc(func(_ any, scope jpl.JPLRuntimeScope) ([]any, jpl.JPLError) {
				result = library.ApplyObject(result, library.FilteredObjectEntries(scope.Vars(), "internals"))
				return nil, nil
			}),
		},
	})
	if err != nil {
		panic(err)
	}
	_, err = program.Run([]any{nil}, nil)
	if err != nil {
		panic(err)
	}
	return result
}

var Builtins = getBuiltins()
