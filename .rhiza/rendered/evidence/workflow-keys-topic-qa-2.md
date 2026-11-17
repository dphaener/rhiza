# workflow-keys-topic-qa-2

evidence · `01M7MMRHQZ6AFV31CE66KEBA6Q` · actor [claude-code-d305ef13](../actor/claude-code-d305ef13.md)

## result

PASS

## statement

QA cycle 2: with a fresh scratch build, rhiza guide workflow-keys lists every FileKeys, SharedKeys, kind (gate included) and role key in internal/workflow/spec.go with its default; rhiza guide workflows names the gate kind and workflow-keys (916 tokens). In a scratch repo, rhiza workflow new myflow writes start=build, agent step build, placeholder command step test failing on purpose with max_visits=3 and on.fail=handback:build; a second run refuses the existing file (exit 3); --from ship-change copies with name=forked. workflow check confirms the topic's claims: an unnamed file takes its base name and first step as start; handback_cap=0, RHIZA_X in [env], a step named done, a gate with no commands, handback to a command step and on.timeout on a gate are all refused. In a throwaway worktree, a new spec.go key and renaming gate in the topic fail TestWorkflowKeys naming the key/kind; dropping the pointer fails TestWorkflowsText.

## Links

### Out

- verifies [workflow-keys-topic](../capability/rhiza/guide/agent-guide/workflow-keys-topic.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
