# manual-resume-refuses-unaccepted-qa-1

evidence · `01M62ZXR6G66Q0F720GXAYPDSV` · actor [claude-code-ea5d2a63](../actor/claude-code-ea5d2a63.md)

## result

PASS

## statement

At 2da4c0b with the built binary in a scratch repo, rhiza run resume on interrupted runs whose change sets are landed, rejected and withdrawn each exits 3 with 'run X is not resumed: its change set <id> is <state>, no longer accepted' and 'next: rhiza run abandon X'; --json gives code conflict with the same next; cmp shows each journal unchanged.

## Links

### Out

- verifies [manual-resume-refuses-unaccepted](../capability/rhiza/workflows/workflow-runs/manual-resume-refuses-unaccepted.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
