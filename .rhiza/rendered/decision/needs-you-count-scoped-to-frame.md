# needs-you-count-scoped-to-frame

decision · `01M7B2JGSK8YSNVHGZDS2GVG0C`

## reasoning

The finding offers two fixes: add the dependency to 01M79NM3RJZFSZE1XHGHCRZSRM, or have this leaf drop the claim. The first edits another proposed change set's plan, which this change set (01M79NKERTK62YZM9C42BJ0JW3, the web frame) does not own and cannot enforce; that change set has its own leaves, its own acceptance (TestHomeNeedsYou) and its own planner. The second is fully within this change set and removes the contradiction outright. Nothing is lost for the frame: the header and picker still share one function, which is what web-project-picker ('a real project picker shared by every screen') and web-frame-header need. Both definitions already list the same five kinds and the same 20-minute stuck rule, so Home can adopt the shared function later at no cost; recording that as a follow-up for Home's plan keeps the intent without stating a dependency the model does not hold.

## statement

web-project-needs-you drops its claim about change set 01M79NM3RJZFSZE1XHGHCRZSRM: it defines the needs-you count and the one function that counts and lists the items for this change set's header and picker only, and promises nothing about Home. Whether web-home-needs-you reuses that function (and gains a depends-on link) is settled when that change set is planned or revised, not here.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7b13vsrpbp9gad9mdp7cpvd-r3-f1](../question/review-01m7b13vsrpbp9gad9mdp7cpvd-r3-f1.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
