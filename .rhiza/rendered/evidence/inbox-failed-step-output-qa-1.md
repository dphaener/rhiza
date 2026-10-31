# inbox-failed-step-output-qa-1

evidence · `01M69QFTJ27CDR61MNWPFZADMM` · actor [claude-code-58a35149](../actor/claude-code-58a35149.md)

## result

PASS

## statement

With the built binary in a scratch repo, after the out workflow failed at command step b, rhiza inbox printed the failed item, '  step "b" ended with outcome "fail"', then '    exit 3', '    out1', '    boom' indented, then the resume --retry next line; rhiza inbox --json gave detail 'step "b" ended with outcome "fail"\nexit 3\nout1\nboom'. A later run failing on 'nosuchcmd_qa' showed 'exit 127' and the sh error indented the same way.

## Links

### Out

- verifies [inbox-failed-step-output](../capability/rhiza/orchestration/inbox/inbox-list/inbox-failed-step-output.md)

## History

- `01M686RRCJSQ57M82FENWT2HNM` When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why
