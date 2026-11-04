# rex-banner-upstream-draft

capability · `01M6JBGRYGYVYDZ0VC01N3941H` · parent [notifications](../notifications.md)

## acceptance

docs/rex-upstream/notification-banner.md is committed holding the Rex server version on the VM (rex version), the exact OSC 7501 bytes the worktree's built rhiza writes for a blocked report (from the agent-tab-blocked-notice testscript or a scratch build's output), the OSC 7501 report Claude Code writes for a blocked agent tab as captured on the VM by rex events --json, the remote VM host setup and the reported symptom (chime, no banner), with clearly marked blank fields for the owner to fill on the Mac (Mac app version, whether a banner now shows, focus and notification-permission state); and rhiza show <question> --json gives an open question, written to the record with rhiza apply --record so it outlives this change set, whose statement names docs/rex-upstream/notification-banner.md and asks the owner to check the Mac, fill those fields and file it with Rex if the banner is still missing

## statement

After agent-tab-blocked-notice is built (it depends on it), the builder drafts a report for Rex upstream in the repo from what the VM can capture, leaving the Mac-side observations for the owner: the blocked report rhiza writes is taken from the worktree's built binary (the agent-tab-blocked-notice testscript or a scratch build outside $PATH), not an installed rhiza, since make install follows the merge. The builder raises an open question in the record for the owner naming the draft; the owner checks on the Mac whether a question in a VM agent tab now gets a banner, fills in the Mac fields and files the report only if it still only chimes

## Links

### Out

- depends-on [agent-tab-blocked-notice](agent-tab-blocked-notice.md)
- decided-by [headless-blocked-notice-in-agent-tab](../../../../decision/headless-blocked-notice-in-agent-tab.md)

### In

- verifies [rex-banner-upstream-draft-verified](../../../../evidence/rex-banner-upstream-draft-verified.md)
- verifies [rex-banner-upstream-draft-qa-1](../../../../evidence/rex-banner-upstream-draft-qa-1.md)
- verifies [pr-245-fix-questions-asked-in-agent](../../../../evidence/pr-245-fix-questions-asked-in-agent.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
