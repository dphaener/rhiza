# wait-stops-on-75-qa-1

evidence · `01M6AJB2QT7497K5BR0339CZDW` · actor [claude-code-6255e6c3](../actor/claude-code-6255e6c3.md)

## result

PASS

## statement

QA by hand at af6e4b54 with a scratch build, in a scratch repo with RHIZA_* and REX_* unset: a hold workflow with one wait step (until = sh ../cond.sh, interval 50ms), where cond.sh exits 1 on poll 1, on poll 2 prints 'stop: first line' to stdout and 'stop: rebase conflict in code.txt' to stderr and exits 75, then exits 0. rhiza run start --workflow hold logged 'step gate: stopped, rebase conflict in code.txt' after 2 polls; run status showed the run stopped, current step gate, detail 'rebase conflict in code.txt' (the last stop: line, stderr included) and no step outcome. rhiza run resume <run> logged 'continuing visit 1', polled a third time, passed and the run succeeded with gate: pass.

## Links

### Out

- verifies [wait-stops-on-75](../capability/rhiza/workflows/workflow-definitions/wait-step/wait-stops-on-75.md)

## History

- `01M6A0JJDJ0FXD58K4WV88FV0A` A ship-change wait settles a PR that conflicts after main moves by rebasing it itself, and stops for the owner on a code conflict
