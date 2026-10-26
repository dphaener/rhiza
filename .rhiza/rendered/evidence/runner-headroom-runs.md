# runner-headroom-runs

evidence · `01M5WKQNHC7ZKWFTVEJWD2MWS5` · actor [claude-code-3e771a5f](../actor/claude-code-3e771a5f.md)

## result

Run 1: roadmap 9.94/9.96 ms, ready 9.82/9.94 ms. Run 2: roadmap 14.78/14.83 ms, ready 14.90/14.81 ms. Run 3: roadmap 13.36/13.87 ms, ready 13.58/13.45 ms (as-is/to-be). All under 25 ms, all runs successful. The runner machine varies about 2x between attempts, which is the 10 vs 14 ms spread.

## statement

On GitHub's runner roadmap and ready have a median under 25 ms in as-is and to-be on three consecutive CI runs of PR #28 at 3f74ee3, with readBudget unchanged at 50 ms.

## tests

Three consecutive attempts of workflow run 40366003546: https://github.com/dphaener/rhiza/actions/runs/40366003546, https://github.com/dphaener/rhiza/actions/runs/40366003546/attempts/2, https://github.com/dphaener/rhiza/actions/runs/40366003546/attempts/3. The reads are timed alone (the log header says 'alone after ...'; no run printed a note line). budget_test.go still has readBudget = 50 ms and no SetGCPercent exists in the code.

## Links

### Out

- verifies [runner-headroom](../capability/rhiza/views/fast-reads/runner-headroom.md)

## History

- `01M5WEV2063MZ6QQ2RHMPW49V4` Reads meet their budget on GitHub's runner with room to spare: roadmap and ready median under 25 ms there, and CI prints the read-budget timings on every run so the runner is measured, not guessed.
