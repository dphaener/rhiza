# m2-interviews-evidence

evidence · `01M5WX85FR5KSERNTM6PCXNCX8`

## result

PR #26's tests (internal/adapter/claudecode, internal/api/session_interview_test.go) pass on main. Known limit: when someone other than the owner answers in the tab, as the orchestrator did in the proof run, the transcript can't tell, so such answers must be labelled by whoever gives them.

## statement

Interview answers become Decisions with the right provenance: change set 7 (PR #26, 20fdee5) reads AskUserQuestion results from the session transcript; an answer that is exactly an option label is recorded as the agent's wording picked by the owner, and only other text and notes are quoted as the owner's.

## Links

### Out

- relates-to [m2-interviews-recorded](../intent/anyone-builds/milestone-2/m2-interviews-recorded.md) — meets this done point

## History

- `01M5WX861EEMJJT1KVJV4FNBE4` apply /Users/darinhaener/code/_genesis/orchestrator/m3plan/m2-evidence.json
