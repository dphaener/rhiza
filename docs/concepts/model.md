# The product model

The model is what Rhiza keeps about your product: what it does, why it's that way, and what's planned. It's a set of nodes joined by typed links, and it lives in your repo next to the code. This page covers the pieces. The change sets page explains how the model changes, and the why page explains why a node is there at all.

It's easier to follow with a few nodes of your own, so if you haven't yet, run through the [quickstart](../quickstart.md) first.

## Nodes

Everything in the model is a node, and every node has a kind. There are ten of them.

| Kind | What it holds |
|---|---|
| `intent` | A goal: something you want to be true. A milestone is an intent. |
| `capability` | Something the product does, or should. This is your feature, requirement or spec. |
| `constraint` | A rule the product has to keep, like a limit or a non-functional requirement. |
| `decision` | A choice that was made, and who made it. |
| `question` | Something nobody knows yet. |
| `evidence` | A test result, a benchmark or any other measurement that backs a claim. |
| `source` | A conversation, a meeting or a transcript, summarised, with the exact words that mattered. |
| `actor` | A person or an agent that does something in the repo. |
| `anchor` | A pointer to code: a path, and a symbol in it. |
| `note` | A remark about any other node or change set. |

Most nodes carry a `statement`, one sentence saying what the node is. Past that, the fields depend on the kind. A decision keeps a `quote`, evidence keeps a `result`, a capability can keep an `acceptance` check. `rhiza describe` lists the fields each kind accepts, and `rhiza add` rejects one it doesn't know.

You add a node with `rhiza add`, giving the kind, a handle and the fields you want. `rhiza set` changes them later.

```bash
rhiza add capability export-csv --field "statement=Reports can be exported as CSV"
rhiza set export-csv --field "acceptance=A 10,000 row report exports in under a second"
```

Only capabilities and intents sit in a tree. The other kinds are flat, and they get their place from the links they have.

## Links

A link joins two nodes and says how they relate. There are fifteen types, and each only allows certain endpoints, so you can't link a decision to a benchmark by accident. `rhiza link` checks the pair before it writes anything.

| Type | From | To | Reads as |
|---|---|---|---|
| `serves` | capability, change set | intent | this work is for that goal |
| `constrains` | constraint | capability, intent | this rule applies there |
| `exempts` | capability | constraint | this capability is excused from that rule |
| `depends-on` | capability, or change set | capability, or change set | this can't be done before that |
| `relates-to` | anything | anything | they're connected, and the note says how |
| `decided-by` | anything | decision | this is the way it is because of that |
| `answers` | decision | question | this settles that |
| `blocks` | question | capability, change set | this can't go ahead until that's answered |
| `supersedes` | a node | a node of the same kind | this replaces that |
| `derived-from` | anything | source | this came out of that conversation |
| `discovered-from` | change set | change set | this was noticed while doing that |
| `implements` | anchor | capability | this code is that capability |
| `verifies` | evidence | capability | this measurement shows it works |
| `observes` | evidence | anchor | this measurement is about that code |
| `discusses` | note | anything | this remark is about that |

"Anything" means any node or a change set. Two of the types need a bit more from you. `relates-to` and `exempts` take a `--note` saying why, and `exempts` also takes a `--decision`, because an exemption with no recorded reason is just a hole in the rule.

```bash
rhiza add decision csv-streaming --field "statement=Export CSV as a stream, not a buffered file"
rhiza link export-csv decided-by csv-streaming
```

## Handles, IDs and references

A node can be named three ways, and any command that takes a node takes any of them.

- **An ID.** Every node gets one when it's created: 26 characters of upper-case letters and digits, sortable by time. It never changes.
- **An ID prefix.** The first four or more characters are enough, as long as they're unique in the repo.
- **A handle.** A short name you choose, lower-case words joined by dashes, like `export-csv`. It's only for people, so you can rename it later. The ID stays the same.

Handles are unique among live nodes across the whole repo, whatever their kind or parent. Two capabilities under different parents can't both be `export-csv`, and a decision can't take a capability's handle. Rhiza refuses a taken handle instead of quietly picking another. A retired node frees its handle for reuse.

Free text can point at other nodes too. Write `@export-csv` in a statement or a note and Rhiza stores the reference by ID, so it still works after the handle changes. A reference starts with `@` at the beginning of the text or after a space or an opening bracket, which is why an email address like `someone@example.com` is left alone. To write a literal `@` where a reference would start, double it: `@@`.

## The capability tree

Capabilities nest. A capability's parent is the bigger thing it's part of, so `reports` can have `export-csv` under it, and `export-csv` can have its own children when it gets big enough to break up.

```bash
rhiza tree
```

The tree is how you find where something belongs before you add it. Before describing a new feature, look for the capability it would sit under and read that capability's decisions and constraints first. They probably apply to yours too. Constraints and decisions on a parent reach its children, and `rhiza slice` pulls them together for you.

Intents have a tree of their own. A milestone is an intent at the top, with child intents under it for its parts, and capabilities link up to them with `serves`.

## Where next

- [The quickstart](../quickstart.md) walks all of this by hand.
- `rhiza describe` prints every kind, field and link type from the code itself.
