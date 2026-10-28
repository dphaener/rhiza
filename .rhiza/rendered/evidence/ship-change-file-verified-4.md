# ship-change-file-verified-4

evidence · `01M61C59Q8ZASW48G9JXV27MAC` · actor [claude-code-02fff339](../actor/claude-code-02fff339.md)

## result

PASS

## statement

At 7e954a9, go test ./cmd/rhiza -run 'TestScripts/(run-ship-change.*|describe-golden)$' -count=1 passes (run-ship-change, -fail, -rebase, describe-golden). run-ship-change runs the built-in ship-change on two leaves (second depends on first) with fake rex/claude/gh: both built and verified in order, build/check then render then review/code and review/review, PR #7 opened once with the review summary as body, run stopped at signoff until rhiza run approve, then landed, merged and closed out. New in this fix: at sign-off a fresh git clone of the pushed cs/ship-it branch passes rhiza check. With the old render step (HEAD~1's ship-change.toml) in a throwaway worktree the same script fails at that clone check: 'extra: render does not write .rhiza/rendered/actor/claude-code-590c1440.md'. So the test catches the uncommitted actor records. Read the diff: the render step now stages and commits .rhiza/log with .rhiza/rendered, as Decision render-commits-actor-records rules. The leaf's statement now lists the render step, as Decision ship-change-file-names-render rules. TestShipChange* in internal/workflow finds no shipit in ship-change.toml or its steps. make check exits 0.

## Links

### Out

- verifies [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
