# plan-run-stops-open-qa-1

evidence · `01M69EAPK7HY0PEN9R3D1C441Z` · actor [claude-code-650867ac](../actor/claude-code-650867ac.md)

## result

PASS

## statement

QA at 29339ad6 with the built binary on PATH in a scratch repo and bare origin, the agent, rex and gh stood in by the plan-workflow-open.txtar shell fakes: rhiza run start <rec> --workflow plan ran interview, review, critic, check, then the gate printed the open question park-disk and the run ended 'stopped' at step gate (not succeeded); no gh call was made. run status --json showed state stopped, current gate, detail 'open questions: park-disk; answer each with a Decision, then rhiza run resume <run>'; rhiza inbox listed a kind open item with next 'rhiza run resume <run>'. Resuming before answering stopped again at the gate with no gh call. After adding Decision disk-ssd answering park-disk in the planned change set, rhiza run resume printed succeeded, gh.log had exactly one 'pr create', and the inbox was empty.

## Links

### Out

- verifies [plan-run-stops-open](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-run-stops-open.md)

## History

- `01M625W0AJZ6AMCCWGFXQW19FN` A plan run that ends with confirmed review findings or parked questions still open does not report success: it stops for the owner with them listed, so nothing has to be fixed by hand after a succeeded run.
