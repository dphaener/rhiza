# plan-task-10

source · `01M5SWRGXHJD80QG3XVV3FXV2B`

## quotes

```json
[]
```

## raw_sha256

fbef4a00d82eec0490436d3bf1349b276aa6d119496ce3b1e6dc4577b6fed80b

## summary

The shipit work plan for task 10: the agent guide shipped in the binary (the genesis skill, the genesis-scrub skill, three topics), its install into Claude Code as a skills-directory plugin with a session-start hook, a size-cap test, a minimal status, and genesis source add. The plan file was lost when the task's worktree was removed and was recovered by the orchestrator from the task's transcripts with D1-D13 intact. Decided (D1-D13, agent wording accepted by the owner): a minimal status showing the repo, current change set, change sets needing attention and the installed guide version against the binary; source add takes a handle, a JSON file of title, summary and quotes, and --raw (or --raw-sha256), storing only the raw material's sha256, through a change set or --record, with --dry-run; the size cap counts tokens as ceil(bytes / 4), with each skill at most 2,000, each topic at most 800 and the hook output at most 100; topics change-sets, apply and sources; guide prints the skill, a topic or the scrub skill, guide ls lists them, and the text is embedded markdown in a pure internal/guide; the guide names the full milestone-1 loop and a test checks every command it names exists, with a pending allowlist for reads not yet built; status --hook is brief, silent outside a Genesis repo and always exits 0; guardrails exclude AGENTS.md, the reconcile skill and evals; install writes a skills-directory plugin at user level by default (or in the repo with --project), with no settings.json edits; a re-install rewrites our own folder and refuses a foreign one without --force; the hook runs genesis status --hook through PATH. Revised after sign-off by the orchestrator under the owner's authorization (not in the plan, applied in PR #8): the scrub skill says an option an agent wrote is never the owner's words, the core loop states the work contract and ends with genesis land, and link accepts change-set refs.

## title

Shipit plan, milestone 1 task 10: the agent guide

## Links

### In

- derived-from [source-add](../capability/genesis/lineage/source-add.md)
- derived-from [scrubbing](../capability/genesis/lineage/scrubbing.md)
- derived-from [agent-guide](../capability/genesis/guide/agent-guide.md)
- derived-from [skills](../capability/genesis/guide/skills.md)
- derived-from [guide-install](../capability/genesis/guide/guide-install.md)
- derived-from [session-start](../capability/genesis/guide/session-start.md)
- derived-from [guide-size-cap](../capability/genesis/guide/guide-size-cap.md)
- derived-from [t10-d1-minimal-status](../decision/t10-d1-minimal-status.md)
- derived-from [t10-d2-source-add-shape](../decision/t10-d2-source-add-shape.md)
- derived-from [t10-d3-size-cap](../decision/t10-d3-size-cap.md)
- derived-from [t10-d4-tests](../decision/t10-d4-tests.md)
- derived-from [t10-d5-topics](../decision/t10-d5-topics.md)
- derived-from [t10-d6-guide-command](../decision/t10-d6-guide-command.md)
- derived-from [t10-d7-full-loop-allowlist](../decision/t10-d7-full-loop-allowlist.md)
- derived-from [t10-d8-status-hook](../decision/t10-d8-status-hook.md)
- derived-from [t10-d9-guardrails](../decision/t10-d9-guardrails.md)
- derived-from [t10-d10-done](../decision/t10-d10-done.md)
- derived-from [t10-d11-skills-dir-plugin](../decision/t10-d11-skills-dir-plugin.md)
- derived-from [t10-d12-reinstall-rule](../decision/t10-d12-reinstall-rule.md)
- derived-from [t10-d13-hook-exec-form](../decision/t10-d13-hook-exec-form.md)
- derived-from [t10-fix1-only-typed-words](../decision/t10-fix1-only-typed-words.md)
- derived-from [t10-fix2-work-contract](../decision/t10-fix2-work-contract.md)
- derived-from [t10-fix3-minor-review-points](../decision/t10-fix3-minor-review-points.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
