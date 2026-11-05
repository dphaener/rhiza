# question-9bf069b9-1

question · `01M6MM63FWGGJ17BEQYE2DX5K3`

## session

9bf069b9-df44-4e1b-b29f-01a1a941fa85

## statement

Change set 01M6MM3 (guides show stdin form for run findings) serves no intent yet. Which intent should it serve? Recent run-reliability fixes (e.g. #259 scoped tests, #267 run abandon) serve milestone-25, whose point is the daemon building and shipping the backlog on its own; a stalled permission prompt in a run is exactly what blocks that. / `rhiza run findings` and `rhiza run verdict` both already accept `-` for stdin. The build, review and plan-review guides show `<file>` for both. Which commands should the guides switch to the stdin heredoc form? / How should the change be tested?

## Links

### In

- answers [stdin-guides-ruling](../decision/stdin-guides-ruling.md)

## History

- `01M6MM3MAXFA2ARM9GAF4RV1ZW` Fix Agents in runs hit Claude Code permission prompts despite bypassPermissions when recording findings: the build, review and plan-review guides show 'rhiza run findings <file>', so agents write mktemp-and-echo one-liners that Claude Code's shell parser flags ('Dangerous rm operation not ruled out'), stalling the session up to 2 minutes until the prompt auto-blocks (seen 2026-11-05 in run 01M6MJW2 verifier and model check). The guides should show the stdin form ('rhiza run findings --check X - <<EOF') so no temp file or compound shell is needed.
