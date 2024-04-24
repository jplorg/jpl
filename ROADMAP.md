# Roadmap

## Additional luxon based classes

- durations
- intervals?

## Encoding

- highlightLocation does not handle unicode codepoints and \r\n correctly
- also the template function specifies that it looks at unicode codepoints but doesn't
- because of that, the whole codebase should be checked for wrong unicode usage
  -> template also doesn't handle this correctly (substring)

## Cross compatibility

- the JS implementation must be updated to be compatible with missing params, otherwise there won't be intercompatibility between Golang generated instructions and the JS runtime
- There should be an official API for exporting and importing instructions from / to a program with versioning to enable compatibility checks - the version should eventually also be included in the REPL when using !i

## REPL

- support multiline input
  - automatically expand to the next line if the input is incomplete
  - handle shift+enter
- add flag for running from a provided program definition
- CLI flags and better output formatting for input/output mode? (`echo '"some jpl"' | jpl-repl -i >file.json`)
