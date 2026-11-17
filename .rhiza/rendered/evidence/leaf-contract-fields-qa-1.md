# leaf-contract-fields-qa-1

evidence · `01M7KW24SRJBBC41GD7JWZ1ST5` · actor [claude-code-ba01156e](../actor/claude-code-ba01156e.md)

## result

PASS

## statement

QA with a fresh build of cs/leaf-contract (go build -trimpath to $TMPDIR, scratch repo and root): rhiza add capability with files and must_not, rhiza set failure_waiver and rhiza apply with files/must_not all wrote and rhiza show read them back; rhiza add constraint with files refused exit 2 (known fields acceptance, statement); a capability field file refused exit 2 listing failure_waiver, files, must_not. A fake-agent run wrote brief.md listing Files:, Must not: and Failure waiver: under each leaf acceptance for the fields it has and none for a leaf without them.

## Links

### Out

- verifies [leaf-contract-fields](../capability/rhiza/workflows/built-in-workflows/plan-workflow/leaf-contract-fields.md)

## History

- `01M76XT3S20YMG07EA5KHDN37M` A leaf names the files it touches and what it must not do, its acceptance has at least one failure case or a waiver, and it is sized to one concern for one fresh session; plan review gains uncovered-intent, oversized-leaf and no-failure-case kinds, and mechanical checks move to review-check. Evidence: 100 of 160 verifier findings are acceptance points unmet or untested (prompt-gap rec 7).
