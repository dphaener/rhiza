# tui-planning-conversation

source · `01M62JM7AH8GAVS9M2EYQD0TP6`

## quotes

```json
[
  {
    "by": "the owner",
    "text": "It should launch from a rex quick actions layer. This will require us to have an install command that installs a lua script in the rex config. We should be able to do something like ctrl-a>r to open a rex specific quick launch rex layer and then select from there which rex \"mini\" app we want to open. This isn't just a single one size fits all tui, it's many mini apps that sort of wrap around Rex and give you small bits of functionality rather than one large TUI that is always open.",
    "where": "Rhiza TUI planning, topic 7: how the TUI starts"
  }
]
```

## raw_sha256

097f856474f2bb324c1e5af1e5e5823dbb2b5a3d5867d3f02d2a33865af66c9f

## session

6bddbe2d-bba9-42bc-96a2-01a185209478

## summary

The owner and the planner planned the Rhiza TUI record on 2026-10-28. Picks, each the planner's option accepted by the owner: the milestone covers command apps and a model view, with writes in this milestone; the command groups are reads (roadmap, ready, show, why, slice), the change-set flow, record writes, and runs and sessions; the apps run in-process through the internal/api registry, so an app and the matching rhiza command --json print the same result; writes from an app show a what-if preview and write after one confirm; the model view starts as the as-is tree and roadmap; the change sets planned now are the shell with the read apps, the model view, the write apps with preview, and the runs and sessions apps. The owner's typed answer on how the TUI starts changed the shape: the apps launch from a Rex quick-launch layer, an install command writes a Lua script into the Rex config, and there is no single large TUI that stays open. Parked: whether that Lua script lives in its own file or in the owner's init.lua, which blocks the install leaf.

## title

Rhiza TUI planning conversation, 2026-10-28

## Links

### In

- derived-from [tui-many-mini-apps](../decision/tui-many-mini-apps.md)

## History

- `01M62J0ED188JTTR7R9709Y989` Rhiza has a terminal UI built on Rex, with rexuikit, in which common rhiza commands run as tiny apps in modals or pop-ups, alongside a view of the model.
