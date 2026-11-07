# ship-tracks-own-pr-verified

evidence · `01M6TH0W1AJ1EVSW060095R1GN` · actor [claude-code-2d865243](../actor/claude-code-2d865243.md)

## result

PASS

## statement

At 006e559c, with fake gh only: go test -trimpath ./internal/api -run 'Ship|PRView|Plan' passes, including TestShipOpenPRRefusesForeignPR (open-pr refuses an OPEN, MERGED or CLOSED PR whose body names another change set or has no 'Change set' line, exit 3 naming #7, no pr create/edit, no ship state file) and TestShipCheckRefusesForeignPR (check with and without --settled, the wait step, refuses other-change-set, no-line and empty bodies in all three states, exit 3 naming #345, no state file); the existing ship tests, whose fake gh bodies name the run's change set, still pass, so an owned PR is tracked as before. The 32 testscripts the commit touched (run-ship-*, plan-workflow-*, ship-check-*, etc.) pass. Diff read: requirePROwned runs right after parsePRView in open-pr (before the merged-PR re-create path and after it) and in check, before any writeShipState; prOwnedBy matches the 'Change set <id>' line planSummary writes. An unfiltered go test ./internal/api/... hit the 600s default timeout on a machine at load ~17, which is unrelated to this leaf; the ship subset passed in 18s.

## Links

### Out

- verifies [ship-tracks-own-pr](../capability/rhiza/orchestration/shipping/ship-tracks-own-pr.md)

## History

- `01M6TBDRB2N5MTW4BW0K2NZCG6` Fix two change sets with the same title prefix share one worktree and branch: the 'usage governor's limits are static' bug request's change set 01M6T2DC got the title 'Fix the usage governor's', the same as 01M6T1ZR, so its plan run used worktree rhiza.fix-the-usage-governor-s and branch cs/fix-the-usage-governor-s; its plan entries rode in 01M6T1ZR's PR #345, its ship state pointed at #345, and its plan run failed at rebased-green ('#345 is MERGED, not open'). A change set's worktree and branch slug must be unique (e.g. add a short ID suffix on a collision), and a run must never commit or track another change set's PR
