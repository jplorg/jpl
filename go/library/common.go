package library

type iteratorIndex[Value any] struct {
	max, current int
	values       []Value
}
