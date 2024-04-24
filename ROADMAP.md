# Roadmap

## Additional luxon based classes

- durations
- intervals?

## Unicode

- highlightlocation should look at runes rather than bytes

## Cross compatibility

- the JS implementation must be updated to be compatible with missing params, otherwise there won't be intercompatibility between Golang generated instructions and the JS runtime
- There should be an official API for exporting and importing instructions from / to a program with versioning to enable compatibility checks - the version should eventually also be included in the REPL when using !i
