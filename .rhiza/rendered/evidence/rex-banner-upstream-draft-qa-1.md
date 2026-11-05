# rex-banner-upstream-draft-qa-1

evidence · `01M6JV7PP42WNXJFK2H8REF79R` · actor [claude-code-eaefca50](../actor/claude-code-eaefca50.md)

## result

PASS

## statement

QA at e2cbef27: docs/rex-upstream/notification-banner.md holds the VM's rex version (matches `rex version` on the VM now: 0.0.1793645657-gd018d2, commit d018d2de), rhiza's exact blocked bytes ending in ST, Claude Code's captured blocked record (state, kind, msg, no title), the VM host setup, the chime-no-banner symptom and blank OWNER fields for the Mac. Its Claude Code record shape matches what real Rex on the VM gave me today for a claude-code blocked report, and rhiza's report as I captured it from real Rex (app=rhiza, title, msg) matches the doc's decoded title and msg format. rhiza show rex-banner-report-ready --json gives a question naming docs/rex-upstream/notification-banner.md and asking the owner to check the Mac, fill the fields and file it if the banner is still missing; it was written by a record-kind change set (01M6JPC8F), so it outlives this one, and has no answer.

## Links

### Out

- verifies [rex-banner-upstream-draft](../capability/rhiza/orchestration/notifications/rex-banner-upstream-draft.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
