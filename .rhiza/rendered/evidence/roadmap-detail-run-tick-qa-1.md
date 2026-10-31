# roadmap-detail-run-tick-qa-1

evidence · `01M6993FVJ9YHTKH809JXTMW5S` · actor [claude-code-84d1030f](../actor/claude-code-84d1030f.md)

## result

PASS

## statement

QA at 0f0c797 in a real Rex pane: in rhiza tui app roadmap, opening change set 01M686RWSEKXXBCNC7RZMT07HB's detail showed 'Latest run ship-change · state running · step build/build'. While the detail stayed open with no key press, that line changed to 'step review/review' when the real run moved on, so the open detail reloads its run section by itself.

## Links

### Out

- verifies [roadmap-detail-run-tick](../capability/rhiza/roadmap-detail-run-tick.md)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
