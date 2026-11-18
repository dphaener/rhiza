# m36-conductor

intent · `01M79KKVJZBSHJ5RXCG4JKBSP4` · parent [milestone-36](../milestone-36.md)

## statement

One conductor session conducts a change set end to end: it plans the leaves when there are none, spawns a builder and then a verifier per leaf, spawns the Plan, Code and QA review checks with a critic only on a blocker or major finding, spawns a fixer for confirmed findings, and lands the change set on the branch.

## Links

### Out

- decided-by [conductor-plans-leaves](../../decision/conductor-plans-leaves.md)
- decided-by [conductor-three-checks](../../decision/conductor-three-checks.md)

### In

- serves change set `01M79KGZ51V2F5YYWMMQTGR6NW` (A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands)
- serves change set `01M79KH59D0HEME6B85PXXB056` (The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic)
- serves change set `01M7EHHAG3CTRQCJCWSG2W7BTZ` (Fix the gate step: the conductor default's gate runs rhiza check with the installed rhiza, not the PR head's build (SHIP_HEAD_RHIZA), so a change set that adds a CLI command, flag or setting fails the gate on docs/reference files its own build passes and hands back to conduct in a loop (runs 01M7DVV7GNJAWG29QZC5F2JWC3 and 01M7E56MGTY4M60XZNEG2A6537; request 01M7EAYSYK7Z32QRP0FB4XPKJN))
- serves change set `01M7M51W2S30H06K0BFZ6MK7PF` (Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand)
- serves change set `01M7MD2K2BE981QPPF5HE8M551` (Fix A rejoined conductor session keeps its first visit's done mark, so after a gate hand-back the session reads as done before it has worked the hand-back (anything closing done sessions closes it mid-visit))
- serves change set `01M7MJN83Y6EF727C2T0AF3Y5D` (The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide)
- serves change set `01M7MN70WCHDYN7H7T01AVZ2VK` (Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
