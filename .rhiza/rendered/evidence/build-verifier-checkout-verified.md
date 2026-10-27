# build-verifier-checkout-verified

evidence · `01M5ZMT5TWDSF7027FYF2PPWWW` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

Testscript build-verifier-checkout runs two leaves through fake rex and claude: each verifier logs dir=checkout, HEAD at its leaf's commit (= RHIZA_COMMIT), branch=HEAD (detached), under .git/rhiza/checkouts; both checkouts (one pass, one exit 1) are gone afterwards and git worktree list shows none; the passing verifier's Evidence is in the build worktree's .rhiza/log and the checkout stayed clean; the builder's staged, unstaged and untracked edits compare equal before and after. rhiza-root covers RHIZA_ROOT sending writes to the build worktree. TestLeavesStepVerifierCheckout finds the checkout removed on pass, fail and a failed removal; the resume and start-failure tests also remove it. All passed.

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/(build-verifier-checkout|rhiza-root)$'; ./internal/run -run 'TestLeavesStepVerifierCheckout|TestLeavesStepResumeRemovesEndedCheckouts|TestLeavesStepVerifierStartFailsRemovesCheckout' (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-verifier-checkout](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-verifier-checkout.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
