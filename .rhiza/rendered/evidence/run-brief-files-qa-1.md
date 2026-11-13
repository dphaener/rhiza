# run-brief-files-qa-1

evidence · `01M7AWMP3KG7MEV6N5H7D6R1TG` · actor [claude-code-1707fde5](../actor/claude-code-1707fde5.md)

## result

PASS

## statement

QA: real runs wrote brief.md with sections in order (change set title/ID/intent, leaves with acceptance, decisions, constraints, check commands showing rhiza.ship.prePushCheck or none, contract ending Full contract: rhiza guide contract); identical bytes across sessions of one run; distinct directories build-1/builder-<leaf>-<attempt>, verifier-<leaf>-<attempt>, review-1/qa-check-1, qa-critic-1, fixer-<leaf>-<cycle>-<attempt>, verifier-<leaf>-<cycle>-<attempt>, fixer-<cycle>-<attempt>. A blocking file at the session path fails the start naming the path. handback.md writer has no trigger here; covered by go test ./internal/run (ok).

## Links

### Out

- verifies [run-brief-files](../capability/rhiza/workflows/workflow-runs/run-brief-files.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
