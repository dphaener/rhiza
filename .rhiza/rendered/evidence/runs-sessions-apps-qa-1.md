# runs-sessions-apps-qa-1

evidence · `01M641NAM16VJJAXQNJYBQTE5R` · actor [claude-code-04e62f34](../actor/claude-code-04e62f34.md)

## result

PASS

## statement

QA at 087eefc with a scratch build ($TMPDIR/rhiza-qa) on this repo's real .rhiza (RHIZA_ROOT/CHANGE/RUN unset): `rhiza tui app run-status --once` was cmp-identical to `rhiza run status --json` with no ref (latest run) and with ref 01M640PW (this review run, live checks and cycles); `rhiza tui app session-ls --once` was cmp-identical to `rhiza session ls --json` (197 KB, every session). An ambiguous ref (01M62) exits 2 for both, the app printing the error as text like the other read apps. In a real pty, `rhiza tui app run-status 01M640PW` and `rhiza tui app session-ls` without --once drew the JSON with a 'j/k scroll · q quit' footer, j scrolled and q quit with exit 0. `rhiza tui launch` lists run-status and session-ls beside roadmap. The .rhiza/log file list was identical before and after all runs, and git status stayed clean. go test of TestScripts/tui-runs-sessions-apps and ./internal/tui/... pass.

## Links

### Out

- verifies runs-sessions-apps (retired)

## History

- `01M62JMTT1Z8FP0SC619V3QS68` Read-only Rex apps for run status and session ls
