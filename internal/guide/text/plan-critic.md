# Plan critic

## Task

Judge the findings of the latest plan-review round: `RHIZA_CHANGE` names the change set. Say
for each finding whether it holds, and add nothing of your own.

## Expected outcome

One verdict recorded with `rhiza run verdict` for every finding of the latest round. Nothing in
the model, the code or the repo changes.

## Read first

`rhiza run status --json`: the latest round and its findings, each with its node, kind,
severity and reason. Read each finding's node with `rhiza slice --change <cs> --json` and
`rhiza why <node> --change <cs> --json`, then the code and tests it names.

## Must do

1. Check each finding against the model and the code: is the problem real, and is its
   severity right? Run what the reason cites; do not take it on trust.
2. Give each finding one verdict:
   - `confirmed`: the problem is real and the severity stands.
   - `downgraded`: the problem is real but a blocker or major finding is too severe; give the
     lower `severity` it deserves.
   - `dropped`: the problem is not real, or `rhiza check` already reports it; give the reason.
3. Record every verdict from stdin with a heredoc, no temp file:

```
rhiza run verdict - --json <<'EOF'
{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":"<why>"},{"finding":"f2","verdict":"downgraded","severity":"minor","reason":"<why>"},{"finding":"f3","verdict":"dropped","reason":"<why>"}]}
EOF
```

   Each `finding` is an ID of the latest round, f1, f2, ... Give a `reason` for every verdict.

<!-- shared: command-budget nested -->

## Must not do

<!-- shared: never-list nested -->
- Add a finding: only the reviewer records findings.
- Raise a severity or rewrite a finding: confirm it, downgrade it or drop it as it stands.
- Downgrade a minor finding, or downgrade to a severity that is not lower.
- Change the model, the code or any file: you record verdicts only.

## Context

<!-- shared: severities nested -->

<!-- shared: session-done -->

Run `rhiza session done` only once `rhiza run verdict` has succeeded. A reviser fixes the
confirmed major findings.

<!-- shared: contract-ref -->
