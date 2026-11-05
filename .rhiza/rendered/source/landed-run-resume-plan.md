# landed-run-resume-plan

source · `01M6J9JS9HF8E155BWW7KFSAWE`

## quotes

```json
[]
```

## raw_sha256

206bc1c15e3513ece3ece6ab8dbf5c20cadac9d17d9b8ce29b66be9d9abde503

## session

0a9726dd-c836-4069-833f-01a1a4935119

## summary

Autonomous plan run, no owner. The planner read run resume's guard (internal/api/run.go, already exempting plan runs on a landed record), the daemon's Resumable accepted check (internal/daemon/pick.go, daemonPoll1), the inbox's skipped listing (internal/api/inbox.go) and ship-change.toml, and found rhiza ship land skips a change set already landed. It raised one question (served intent, plan-run scope given draft 01M6H5VQ, and the rule for which landed runs resume). The arbiter decided (landed-run-resume-scope): serve m4-runs-hold-up; cover non-plan runs only, leaving plan runs to 01M6H5VQ; a non-plan run on a landed change set resumes, by hand and by the daemon, when its journal shows it reached the land step; others stay refused with rhiza run abandon as next. Leaves planned: run resume continues such a run, the daemon resumes it headless, and the inbox stops listing it as skipped.

## title

Plan of 01M6J965: resume runs whose change set landed

## Links

### In

- derived-from [landed-run-resume-scope](../decision/landed-run-resume-scope.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
