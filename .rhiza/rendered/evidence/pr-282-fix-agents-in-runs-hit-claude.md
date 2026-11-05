# pr-282-fix-agents-in-runs-hit-claude

evidence · `01M6MWSM9GJ8YMFC7Z5ZF4XK8Y` · actor [dphaener](../actor/dphaener.md)

## head

91d4413a498425dfcef546143c8c23d8c5bc7096

## pr

https://github.com/dphaener/rhiza/pull/282

## result

1 check passed: check SUCCESS

## statement

Pull request #282 (https://github.com/dphaener/rhiza/pull/282) passed CI at 91d4413

## status

verified

## Links

### Out

- verifies [guides-stdin-findings](../capability/rhiza/workflows/built-in-workflows/guides-stdin-findings.md)

## History

- `01M6MM3MAXFA2ARM9GAF4RV1ZW` Fix Agents in runs hit Claude Code permission prompts despite bypassPermissions when recording findings: the build, review and plan-review guides show 'rhiza run findings <file>', so agents write mktemp-and-echo one-liners that Claude Code's shell parser flags ('Dangerous rm operation not ruled out'), stalling the session up to 2 minutes until the prompt auto-blocks (seen 2026-11-05 in run 01M6MJW2 verifier and model check). The guides should show the stdin form ('rhiza run findings --check X - <<EOF') so no temp file or compound shell is needed.
