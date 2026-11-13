# review-01m7b13vsrpbp9gad9mdp7cpvd-r3-f1

question · `01M7B2GQSBCTN0TD11H6832ZQJ`

## statement

Plan review round 3 left finding f1 (contradiction on web-project-needs-you): The statement says web-home-needs-you (change set 01M79NM3RJZFSZE1XHGHCRZSRM) 'depends on this change set and builds after it' and counts Home's items with this leaf's function, but the model bears none of it out: 01M79NM3RJZFSZE1XHGHCRZSRM's leaves list no needs, it has no depends-on link to web-project-needs-you, and web-home-needs-you's statement defines its own list and its own stuck rule. Either that change set gets the dependency (its plan) or this leaf drops the claim; as written the shared-count guarantee (header, picker and Home showing the same number) has nothing enforcing it. The critic confirmed it: web-project-needs-you's statement says change set 01M79NM3RJZFSZE1XHGHCRZSRM depends on this change set and counts with this leaf's function. rhiza change show 01M79NM3RJZFSZE1XHGHCRZSRM lists only its own three node.create and three link.create operations, with no depends-on link. Its web-home-needs-you statement and acceptance (TestHomeNeedsYou) define their own list and stuck rule and never name a shared function. Nothing in the model enforces the claimed shared count, so the leaf states a dependency the model contradicts. Major holds for a contradiction. How should it be settled?

## Links

### Out

- blocks [web-project-needs-you](../capability/rhiza/web-monitor/web-project-needs-you.md)

### In

- answers [needs-you-count-scoped-to-frame](../decision/needs-you-count-scoped-to-frame.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
