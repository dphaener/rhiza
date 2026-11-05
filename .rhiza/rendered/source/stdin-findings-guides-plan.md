# stdin-findings-guides-plan

source · `01M6MM9QTBY8HKPWN25PCBEKDE`

## quotes

```json
[
  {
    "by": "arbiter (assumed, would have escalated)",
    "text": "Change set 01M6MM3 serves milestone-25 ... Every place the build, review and plan-review guides show `rhiza run findings ... <file>` or `rhiza run verdict ... <file>` switches to the stdin heredoc form ... No CLI change: both commands already accept `-`.",
    "where": "ruling stdin-guides-ruling"
  }
]
```

## raw_sha256

a0c7fcdbfba3c29b67d3201c3c80252c7ff1452554c2a1cc67009bc81ea43a40

## session

9bf069b9-df44-4e1b-b29f-01a1a941fa85

## summary

Autonomous plan run, no owner present, so no owner words are quoted. The change set's intent: agents in runs hit Claude Code permission prompts when recording findings, because the build, review and plan-review guides show `rhiza run findings <file>`, so agents write mktemp-and-echo shell one-liners. The planner found both `rhiza run findings` and `rhiza run verdict` already read `-` as stdin (readReviewFile, internal/api/run_review.go). One combined question went to the arbiter (intent, scope, acceptance); its ruling stdin-guides-ruling: serve milestone-25; switch every findings and verdict `<file>` in the three guides to the stdin heredoc form `- <<'EOF'`, keeping flags and JSON shapes; no CLI change; `rhiza apply <file>` stays as is; test through internal/guide's build, review and plan tests. Planned: one leaf, guides-stdin-findings, under built-in-workflows. Nothing left open.

## title

Plan: guides show stdin form for run findings and verdicts (autonomous)

## Links

### In

- derived-from [stdin-guides-ruling](../decision/stdin-guides-ruling.md)

## History

- `01M6MM3MAXFA2ARM9GAF4RV1ZW` Fix Agents in runs hit Claude Code permission prompts despite bypassPermissions when recording findings: the build, review and plan-review guides show 'rhiza run findings <file>', so agents write mktemp-and-echo one-liners that Claude Code's shell parser flags ('Dangerous rm operation not ruled out'), stalling the session up to 2 minutes until the prompt auto-blocks (seen 2026-11-05 in run 01M6MJW2 verifier and model check). The guides should show the stdin form ('rhiza run findings --check X - <<EOF') so no temp file or compound shell is needed.
