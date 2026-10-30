# build-review-render-plan

source · `01M650620KYZAB9S44N5MJVFCK`

## quotes

```json
[]
```

## raw_sha256

0eed233534f114ea32cba94d50757fe3532c4c2426b88b1f2a3979af924b582a

## session

5a87ddc7-48a8-4fe9-b35c-01a18a012b40

## summary

Change set 01M61GW7K1PDG7HMFXMN3YNA7B (Build and review render) intends that plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit. It serves milestone-4's done point m4-runs-hold-up. The run had no owner (autonomous), so nothing was asked of anyone and no owner words exist. The plan follows the ship-change workflow, which already renders and commits after build and after review. Decided by assumption: the render lives in the build and review workflow files, not in the run engine, and it is the last step of each. Open: whether build renders after a failed build too.

## title

Build and review render, planned autonomously

## Links

### In

- derived-from [render-at-workflow-end](../decision/render-at-workflow-end.md)

## History

- `01M61GW7K1PDG7HMFXMN3YNA7B` Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.
