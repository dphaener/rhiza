# byo-levels-skill-conductor-workflow

decision · `01M7E5H41JVA4Q7FXBRXE8B9DM`

## question

Change set 01M79KHH (Bring-your-own docs) asks for a 'Bring your own agent' page with one worked example per level. What are the levels? (Today: conductor spawns follow `rhiza guide <role>` with no per-role invoke; invokes exist on agent steps and on leaves/checks role tables; rhiza skill check checks builder, verifier and check roles.)

## reasoning

Conductor spawns take no invoke (run spawn starts each child with 'follow rhiza guide <role>'), so a team's own role skill is reachable today only through a leaves or checks role table, which ship-change-classic has; the conduct step's prompt is the place a team's own conductor goes. A level built on a second agent adapter would describe code that does not exist (claude-code is the only adapter), so the agent in 'bring your own agent' is the team's own agent session and skill. Example files kept valid by a test follow docs/examples/workflows and TestDocsExampleWorkflows.

## statement

The Bring your own agent page has three levels, each with one worked example: a skill (a team's installed skill plays one role: rhiza skill check it, then put it in a role table's invoke in a workflow forked from ship-change-classic), a conductor (the team's own conduct skill on ship-change's conduct step, such as prompt = "/myteam:conduct --prompt {brief}", reading {brief}, {slice} and {handback} while the daemon's gate still checks the result) and a workflow (the team's own workflow built from the primitives). Each example is a TOML file under docs/examples/byo/ that a test keeps passing rhiza workflow check.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [byo-docs-planning](../source/byo-docs-planning.md)

### In

- decided-by [docs-bring-your-own-agent](../capability/rhiza/documentation/docs-bring-your-own-agent.md)

## History

- `01M79KHHVJK2RMMC25KY63QW13` Bring-your-own docs: the primitives reference and a Bring your own agent page with one worked example per level, and rhiza init asks which agent the project uses
