# plan-role-topics-shape

decision · `01M79BD2TMFK8RCJJHZS5N5M8G`

## reasoning

Followed Decision build-role-topics-shape (change set 01M76XR5), which already settled the same four choices for the build roles: new topics installed as role skills, builtin prompts cut to one line naming both routes, the old topic kept as an index so custom workflows and the rhiza-plan skill's triggers keep working (AGENTS.md treats removing something users rely on as breaking), and shared blocks expanded at load from internal/guide/text/shared/ with the review and plan topics named as later users. Read internal/guide/guide.go (plan is a skill, plan-review a topic not installed as a skill), internal/api/guide_install.go, the builtin plan.toml and plan-milestone.toml prompts (the reviewer, critic and reviser all point at `rhiza guide plan-review`), internal/run/journal.go (Kinds, KindSeverity) and plan-review.md. Leaf contract 01M76XT3 carries a note that this change set owns uncovered-intent and oversized-leaf, so it adds them here. Both are major because plan-review.md says minor and nit are never fixed: a minor oversized-leaf would never be split, defeating the milestone's build-retry target, and an uncovered intent clause is as serious as outside-intent, already major. The research pass and self-check stay in the guide because the intent asks for prompt content, prompt-gap's CS1 row keeps these change sets to guide and toml, and a self-check against the reviewer's own kinds catches the findings before a review round. The severity block is already required to be byte-identical between review and plan-review by TestReviewSeverityBlockShared, the exact duplication the shared-block mechanism exists to remove, so it becomes one shared file.

## statement

Delivery: four new topics, `planner`, `plan-reviewer`, `plan-critic` and `plan-reviser`, each on the six-section skeleton (Task, Expected outcome, Read first, Must do, Must not do, Context), installed by guide install as the skills rhiza-planner, rhiza-plan-reviewer, rhiza-plan-critic and rhiza-plan-reviser. The builtin plan.toml and plan-milestone.toml prompts shrink to one line per role that carries the TASK variables and names both routes, e.g. `/rhiza-planner Plan change set {change} ({title}): follow rhiza guide planner.`, `/rhiza-plan-reviewer Review the plan of change set {change}: follow rhiza guide plan-reviewer.`, and the same for the critic and reviser; the session-done and command rules now in those prompts move into the topics through the shared blocks. The planner part of today's plan topic (7985 bytes, 1997 tokens, against TopicCap 1000) cannot fit one role topic, which Build role topics measures, shared blocks expanded, against TopicCap; so it splits three ways, each text within TopicCap: the `planner` role topic holds the steps, the research pass, the self-check, autonomous runs and the shared blocks, for a plan run and for planning with the owner in a chat, and its Read first names two plain reference topics (not installed as skills) the planner loads with `rhiza guide`: `plan-recipes` (scope by kind, the apply recipes, who records a pick, Decisions) and `plan-interview` (the interview with the owner, the Source it keeps, and milestone planning). The `plan` topic and rhiza-plan skill stay, with their description (so the skill still triggers), cut to a short index naming the four role topics and which to load; the `plan-review` topic stays as a short index the same way. Nothing is removed, so there is no breaking note. Kinds: this change set adds `uncovered-intent` (a clause of the change set's intent that no leaf serves) and `oversized-leaf` (a leaf with more than one concern, or too large for one fresh session) to run.Kinds, documented after the five existing kinds, and KindSeverity makes both major; Leaf contract (01M76XT3) adds only no-failure-case and the move of mechanical checks to review-check, as its scope note says. The plan reviewer topic tells the reviewer not to record what `rhiza check` already reports. Research pass and self-check: guide prose in the planner topic, no new command or Go. Before writing any Decision or leaf, the planner reads the slice and `rhiza why` of the change set, then the code and tests the intent touches (searching for the names it uses) and the patterns a builder would follow, and names in each leaf the files it touches; it raises a question only for what that reading cannot settle. After writing, before session done, it re-reads `rhiza slice <change> --json` and checks the plan against the reviewer's kinds: every intent clause served by a leaf, every leaf one concern with a runnable acceptance scoped to its tests, every Decision with provenance, no unreal depends-on, then runs `rhiza check`, and fixes what it finds. Shared blocks: this change set depends on Build role topics (01M76XR5) and its leaves that use shared blocks depend-on its guide-shared-blocks capability, so it builds after that lands. The four plan topics name the existing shared blocks that apply (session done rule, command budget, never-list, outside work) with `<!-- shared: <block> -->` lines. The severity block moves out of review.md and plan-review.md into a new shared file internal/guide/text/shared/severities.md, named by the review topic(s) and by plan-reviewer and plan-critic; its plan-review line becomes `missing-acceptance`, `contradiction`, `outside-intent`, `uncovered-intent` and `oversized-leaf` major, the rest minor; review_test's TestReviewSeverityBlockShared becomes a test that the block lives in one shared file and is expanded in each topic that names it.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-3c4cbc82-1](../question/question-3c4cbc82-1.md)
- derived-from [plan-role-topics-interview](../source/plan-role-topics-interview.md)

### In

- decided-by [plan-finding-kinds](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-finding-kinds.md)
- decided-by [shared-severities-block](../capability/rhiza/workflows/built-in-workflows/plan-workflow/shared-severities-block.md)
- decided-by [planner-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/planner-topic.md)
- decided-by [plan-reviewer-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviewer-topic.md)
- decided-by [plan-critic-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-critic-topic.md)
- decided-by [plan-reviser-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-topic.md)
- decided-by [plan-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-role-skills-installed.md)
- decided-by [plan-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-prompts-one-line.md)
- decided-by [plan-guide](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-guide.md)
- decided-by [plan-review-guide](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-guide.md)

## History

- `01M76XRPN6WYMZSV2GWA05Y1FM` Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).
