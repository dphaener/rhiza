# cli-help-style-qa-1

evidence · `01M7301K9G9WV21XXD7JQA3NK8` · actor [claude-code-e7e28e5b](../actor/claude-code-e7e28e5b.md)

## result

PASS

## statement

QA at 88c941b6 with a binary built from the branch: rhiza request/show/root init/project add/change new/feed/daemon --help print the reworded descs (choices as 'One of ...'/'Either ...', defaults as '; default <x>', capitalised, no trailing period). rhiza docs in a scratch repo wrote cli.md byte-identical to docs/reference/cli.md. go test ./internal/api -run TestHelpStyle passes, and mutating RootInitRequest.Dir's desc to 'the directory (default: cwd).' made it fail with 4 errors. go test ./cmd/rhiza -run 'TestScripts/(describe-golden|docs-reference)' passes (5 scripts) with GOTMPDIR outside .git. Cosmetic: enum flags show the choices twice, e.g. request --kind 'One of feature, bug or milestone; ... (one of feature, bug, milestone)'

## Links

### Out

- verifies [cli-help-style](../capability/rhiza/documentation/cli-help-style.md)

## History

- `01M71G4RTMS00MVHJ8D9WMD273` CLI help wording: the desc tags of every request struct read plainly and in one style, so rhiza --help and the generated docs/reference are consistent
