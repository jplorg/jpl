package jpl

type JPLDataType string

const JPLT_NULL = JPLDataType("null")
const JPLT_FUNCTION = JPLDataType("function")
const JPLT_BOOLEAN = JPLDataType("boolean")
const JPLT_NUMBER = JPLDataType("number")
const JPLT_STRING = JPLDataType("string")
const JPLT_ARRAY = JPLDataType("array")
const JPLT_OBJECT = JPLDataType("object")

// Order which applies when comparing values with different types
var TypeOrder = []JPLDataType{JPLT_NULL, JPLT_FUNCTION, JPLT_BOOLEAN, JPLT_NUMBER, JPLT_STRING, JPLT_ARRAY, JPLT_OBJECT}
