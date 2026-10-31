# command-step-keeps-output-qa-1

evidence · `01M69QFQW6PBQBDQZJG2DA54JF` · actor [claude-code-58a35149](../actor/claude-code-58a35149.md)

## result

PASS

## statement

Built rhiza at fc88da4a into $TMPDIR and ran it by hand in a scratch git repo (rhiza init, change new). A command step 'echo out1; echo boom >&2; exit 3' failed the run and run status --json gave the step detail "exit 3\nout1\nboom". A step printing 'line 1'..'line 500' then exit 1 kept 'exit 1' plus lines 481-500 (21 lines, 186 bytes). A 5000-char single line kept 2048 bytes; 3000 'é' kept a 2046-byte tail starting on a rune boundary. A step that exits 5 silently keeps 'exit 5'; 'nosuchcmd_qa' kept 'exit 127' plus the sh error; a stderr-to-/dev/null step kept only its stdout. A passing step that writes to stdout and stderr kept detail exactly 'exit 0'. go test ./internal/run -run TestCommandDetailKeepsRunnerOutput passed.

## Links

### Out

- verifies [command-step-keeps-output](../capability/rhiza/workflows/workflow-definitions/command-step/command-step-keeps-output.md)

## History

- `01M686RRCJSQ57M82FENWT2HNM` When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why
