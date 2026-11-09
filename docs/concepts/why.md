# Decisions, sources and why

Most of the model says what the product is. This page is about the other half: why it is that way. Code can't tell you that. A year from now somebody will ask why reports stream instead of loading into memory, and the answer is in a hallway conversation nobody wrote down. Rhiza keeps it next to the thing it explains.

Four kinds of node carry this: decisions, questions, sources and evidence. They're all [nodes like any other](model.md), written through [change sets](change-sets.md).

## Decisions and the owner's words

A decision is a choice somebody made. Its `statement` says what was chosen, and its `quote` keeps the words of whoever made the call.

```bash
cat > decision.json <<'EOF'
{"writes":[
  {"op":"add","kind":"decision","handle":"csv-streaming",
   "fields":["statement=Stream CSV rows instead of building the file in memory",
             "quote=\"Stream it, reports get huge\" (owner, 2026-10-20)"]},
  {"op":"link","from":"export-csv","type":"decided-by","to":"csv-streaming"}
]}
EOF
rhiza apply decision.json
```

The quote is verbatim. If the owner said something, write down what they said, not your summary of it. A summary has already been through your reading of it, and that's the part somebody will want to check later. It matters most when an agent is the one writing: the quote is how you tell a decision the owner made from one an agent made up.

The `decided-by` link is what ties the decision to what it explains. Any node can point at one, so a capability, a constraint or a whole change set can say "this is the way it is because of that." A decision that replaces an older one gets a `supersedes` link, and the old one stays in the model so you can still see what it said.

Accepting a change set records a decision too. `rhiza accept` and `rhiza reject` take a `--reason` and a `--rule` (the default is owner), so "the owner said yes" ends up in the log without anyone typing it twice.

## Questions

A question is something nobody knows yet. Writing it down is better than guessing and moving on, because a guess looks exactly like knowledge by the time it's in the code.

```bash
rhiza add question csv-delimiter --field "statement=Which delimiter do exports use for locales that write decimals with a comma?"
rhiza link csv-delimiter blocks export-csv
```

`blocks` means the question stops that capability or change set until somebody answers. Agents that hit one stop and wait instead of choosing for you. When you have the answer, `rhiza question answer` writes it as a decision quoting your words and links it `answers` to the question. It works on a question that blocks the current change set, which is how a run's questions get answered:

```bash
rhiza question answer csv-delimiter "Semicolon there, comma everywhere else."
```

The answer goes in as given. Rhiza doesn't tidy it.

## Sources

Decisions often come out of a long conversation: a planning chat, a call, an email thread. A source holds that conversation in a form the repo can keep. It has a title, a short summary and the exact quotes that mattered.

```bash
cat > source.json <<'EOF'
{"title":"CSV export planning",
 "summary":"Reports can reach millions of rows, so the owner wants export streamed. Locale handling is left open.",
 "quotes":[{"text":"Stream it, reports get huge","by":"owner","where":"planning chat, 2026-10-20"}]}
EOF
rhiza source add csv-planning source.json --raw-sha256 a34bb6c0bbacc3adca0cb8cf388cd5955381bdb6c896103669217955bd86010a
```

The raw material itself doesn't go in. You give `rhiza source add` one of `--raw` (a file, which it hashes), `--raw-sha256` (a hash you already have) or `--raw-session` (an agent session, hashed the same way). Only the hash is stored, so if the original turns up later you can show it's the one the summary came from, and nothing private has to live in the repo.

That's also why you scrub first. A transcript has secrets, personal details and a lot of "um" in it. The summary and the quotes are what you keep, so go through them before you add the source and take out anything that shouldn't be in a public repo. Quotes stay exact, and everything around them gets shorter.

Nodes that came out of a source get a `derived-from` link to it, which is how a decision points back at the conversation it was made in.

## Evidence

Evidence is a measurement that backs a claim: a test result, a benchmark, a check somebody ran. It has a `statement` saying what was measured and a `result`.

```bash
rhiza add evidence export-csv-memory --field "statement=Exporting 2M rows peaks at 40 MB resident" --field "result=PASS"
rhiza link export-csv-memory verifies export-csv
```

`verifies` says the measurement shows the capability works. `observes` points evidence at an anchor instead, when it's about a piece of code and not a behaviour. In a build run the verifier writes this, so "done" rests on somebody checking and not on the builder saying so.

## Reading it back

You don't have to go find any of this. Two commands pull it out.

`rhiza why` walks backwards from a node: the decisions behind it, the sources those came from, and the change sets that touched it.

```bash
rhiza why export-csv
```

It's what you'd run on code you don't understand, to find out whether the odd part is a mistake or a decision. If it's a decision, the quote is right there. If there's nothing behind it, that's worth knowing too.

`rhiza slice` goes wider. It composes the context around a node in one read: the node, its parents, the decisions and constraints that apply to it, trimmed to a token budget. It's the first thing an agent runs before touching code, so the decisions you recorded are the ones it sees.

```bash
rhiza slice export-csv
```

Both read the model as it'll be once the current change set lands, by default, so a draft shows its own decisions. `--view as-is` reads only what has landed.

## Where next

- [The product model](model.md) covers nodes, links and handles.
- [Change sets](change-sets.md) covers how all of this gets written and landed.
- [The quickstart](../quickstart.md) has a first decision in it.
