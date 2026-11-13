# m36-byo

intent · `01M79KKVMQRDB3EETMZHRCNATG` · parent [milestone-36](../milestone-36.md)

## statement

A team brings its own skills, agent and workflow: its already-installed skill runs through an invoke template, reads brief.md, slice.json and handback.md, meets the agent contract that the daemon's gates check, and passes rhiza skill check before it runs for real.

## Links

### In

- serves change set `01M79KGHYF3272SZ9H9Y12FWF4` (Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps)
- serves change set `01M79KHBKTKGWFD3ZGV1B5QPHA` (rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate)
- serves change set `01M79KHHVJK2RMMC25KY63QW13` (Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses)

## History

- `01M79KF4RGGHZCJ7PJWMP01XXM` Conductor workflow. Owner, 2026-11-13, after the overnight watch and the ship-change proposal (https://claude.ai/artifact/Vp1BQRYN5793H1N3GURe3h): an agent conducts each change set and spawns visible sessions through a dumb, durable daemon that only ships and guards; workflows are built from a few daemon primitives so teams bring their own skills and agents; the new default replaces the builtin ship-change at once.
