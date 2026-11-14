# plan-skill-check

source · `01M7DW3G5E24NZVATBZGT5YQBQ`

## quotes

```json
[]
```

## raw_sha256

d24fec6cd3ef4ccfa0faf6662699e0233e6e2cc571c0aed7b920c7b9e7f72c31

## session

61c60493-d643-4ab8-b81d-01a1dbbdfd32

## summary

Autonomous plan of change set 01M79KHB in run 01M7DVV7GNJAWG29QZC5F2JWC3 (ship-rhiza, no owner). Research: rhiza workflow check (internal/api/workflow_check.go) reports problems and warnings (exposure, placeholder); every handback:<agent step> route is already capped by handback_cap (internal/run/engine.go), so the only loop it does not count is a red outcome naming an earlier step by plain name; session start already refuses a start command whose skill is not installed (checkStartSkill in internal/api/session.go, adapter.SkillChecker's StartSkill and SkillGlobs); the built-in ship-change merges with rhiza ship merge after a gate step. Two questions went to the orchestrator as arbiter: what counts as an uncapped hand-back loop (ruled: a plain red loop, the only loop handback_cap does not count, and it delivers no handback.md) and how skill check runs a skill (ruled: a scratch repo outside every git repository with a one-leaf fixture change set and one role, because the intent says the skill passes rhiza skill check before it runs for real, so it must not touch the user's repo). The orchestrator also set machine rules for every child: capped systemd-run scopes, scoped go test -p 2, one test command at a time.

## title

Plan of rhiza skill check and workflow check warnings

## Links

### In

- derived-from [uncapped-loop-is-plain-red-loop](../decision/uncapped-loop-is-plain-red-loop.md)
- derived-from [skill-check-scratch-one-role](../decision/skill-check-scratch-one-role.md)
- derived-from [skill-check-plan-shape](../decision/skill-check-plan-shape.md)
- derived-from [answer-61c60493-3](../decision/answer-61c60493-3.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
