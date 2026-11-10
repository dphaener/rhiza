# cli-help-style-verified

evidence · `01M72ZJB99EY84YZZG8H9JZA9Q` · actor [claude-code-0c0c1548](../actor/claude-code-0c0c1548.md)

## result

PASS

## statement

Verified at 5086a123. go test -trimpath ./internal/api -run TestHelpStyle passes; in a scratch git-archive copy, mutating RootInitRequest.Dir's desc to a lowercase start, a trailing period, and '(default: cwd)' each made TestHelpStyle fail with the matching message (capital, period, default form), so the test enforces all three mechanical rules over every operation summary and every desc reachable from Operations() request types. go test ./cmd/rhiza -run 'TestScripts/(describe-golden|docs-reference)' passes (describe-golden, docs-reference-cli, -config, -workflow, -stale), so describe-golden and docs/reference/cli.md are regenerated and current. Diff read: 15 descs reworded (choices as 'One of a, b or c' / 'Either a or b', defaults as '; default <x>'); remaining lowercase descs are only on response types, outside the leaf's request-struct scope. Note: rendered help appends the frame's own '(one of feature, bug, milestone)' after a desc already starting 'One of ...', a cosmetic repeat, not an acceptance failure

## Links

### Out

- verifies [cli-help-style](../capability/rhiza/documentation/cli-help-style.md)

## History

- `01M71G4RTMS00MVHJ8D9WMD273` CLI help wording: the desc tags of every request struct read plainly and in one style, so rhiza --help and the generated docs/reference are consistent
