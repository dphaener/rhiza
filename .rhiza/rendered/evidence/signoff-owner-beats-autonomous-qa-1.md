# signoff-owner-beats-autonomous-qa-1

evidence · `01M73HWH1A03E10TW7KDZQ0AY5` · actor [claude-code-ead0f95a](../actor/claude-code-ead0f95a.md)

## result

PASS

## statement

QA cycle 1 at 375f9198: built the branch with go build -trimpath into a scratch dir and ran it by hand (not through testscript) in fresh scratch Rhiza roots, with RHIZA_* unset and only rex and claude faked. (1) A workflow with [autonomy] signoff = owner, rhiza run start --foreground --autonomous: log 'step signoff: waiting for the owner: workflow review declares autonomy signoff owner', run status state stopped, current signoff, signoff_waiting true, no signoffs and no signoff journal entry; rhiza run approve --foreground then succeeded with one signoff of rule owner. (2) The same workflow declaring signoff = auto, and one with no [autonomy] at all: each signed off under rule autonomous and succeeded. (3) Nested: a top-level workflow declaring signoff = owner whose only step runs an inner workflow holding the signoff step: inner/signoff waited, and the log named workflow review's declaration. (4) request-changes on an owner-declared run (max_visits raised): it went back to review (visit 2), then stopped at signoff waiting again, with only the owner's changes signoff in the journal and none of rule autonomous. run status showed 'Autonomy signoff: owner from workflow review'

## Links

### Out

- verifies [signoff-owner-beats-autonomous](../capability/rhiza/orchestration/accept-policy/signoff-owner-beats-autonomous.md)

## History

- `01M73FYE2PAG9E034MPBTTJ12D` Fix A workflow that declares 'autonomy signoff owner' is signed off autonomously anyway. The repo's docs workflow (#425) declares signoff = owner (rhiza workflow describe docs: 'autonomy signoff owner'), yet run 01M73CQ4D71CW038ZBNSQB1SJF (Concepts polish, PR #426) logged 'step signoff: signed off under rule autonomous' and went on to land; run 01M73CQ3TCG0TVWWR8ZV62VAE8 (README, PR #427) is headed the same way. Settings: rhiza.autonomy.signoff = auto (default, project); the daemon runs --autonomous. The milestone's Decisions require the owner to sign off each page, so a workflow's declared owner stage must win over the default/--autonomous (or the docs workflow is not honoured). Held by converting PRs #426 and #427 to draft so ship merge cannot merge them; the owner signs off on the PRs by hand.
