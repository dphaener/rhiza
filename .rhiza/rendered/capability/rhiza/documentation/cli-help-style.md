# cli-help-style

capability · `01M71G7R77J51ATRS08H1X485G` · parent [documentation](../documentation.md)

## acceptance

A Go test over the registry fails on a description breaking the style's mechanical rules (capital first letter, no trailing period, defaults written as '; default <x>'), describe-golden is regenerated and TestScripts/describe-golden passes, rhiza docs output is regenerated and go test ./cmd/rhiza -run 'TestScripts/docs-reference' passes

## statement

Every command, argument and flag description (the desc tags of the request structs) follows one plain style: a capitalised sentence with no trailing period, saying what it does in everyday words, with defaults and choices in one consistent form, so rhiza --help and the generated docs/reference read alike

## Links

### Out

- decided-by [m31-polish-request](../../../decision/m31-polish-request.md)

### In

- verifies [cli-help-style-verified](../../../evidence/cli-help-style-verified.md)
- verifies [cli-help-style-qa-1](../../../evidence/cli-help-style-qa-1.md)
- verifies [pr-422-cli-help-wording](../../../evidence/pr-422-cli-help-wording.md)

## History

- `01M71G4RTMS00MVHJ8D9WMD273` CLI help wording: the desc tags of every request struct read plainly and in one style, so rhiza --help and the generated docs/reference are consistent
