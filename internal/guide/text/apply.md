# apply

`rhiza apply <file> --json` writes many operations from one JSON file as one entry, or
writes nothing. `-` reads stdin. Use it instead of many single writes.

## The file

```json
{"writes":[ {"op":"add", ...}, {"op":"link", ...} ]}
```

At least one item. Each item is one write's request plus `"op"` naming it. Items apply in
order, and later items may refer to handles added by earlier ones. Unknown keys are errors.

| op | fields |
|---|---|
| `add` | `kind`, `handle`; `parent`, `before`, `after`, `fields`, `json_fields` |
| `set` | `ref`; `handle` (new), `fields`, `json_fields`, `unset` |
| `mv` | `ref`; `parent` or `root: true`; `before`, `after` |
| `retire` | `ref`, `reason` |
| `link` | `from`, `type`, `to`; `note`, `decision` |
| `unlink` | `ref` (the link, or its from node with `type` and `to`) |

- `fields` is a list of `"name=text"`; `json_fields` a list of `"name=<json>"`. Field names
  are `^[a-z][a-z0-9_]*$`. `statement` is the usual main field.
- `kind`: intent, capability, constraint, decision, question, evidence, source, actor, anchor, note.
  Only capabilities and intents take a `parent` (of the same kind).
- A ref is a handle, an ID, or a unique ID prefix of 4+ characters.
- Link types: serves, constrains, exempts, depends-on, relates-to, decided-by, answers,
  blocks, supersedes, derived-from, discovered-from, implements, verifies, observes, discusses.
  `relates-to` and `exempts` need a `note`; `exempts` needs a `decision`.
- Common links: anything `decided-by` a decision; a decision `answers` a question; a question
  `blocks` a capability; a capability `serves` an intent; a constraint `constrains` a
  capability; anything `derived-from` a source.

## Build order

A leaf `depends-on` another only when it must be built first: it uses the other's code or
output. Don't link leaves that are merely related.

```json
{"op":"link","from":"export-csv","type":"depends-on","to":"report-query"}
```

`rhiza change show --json` then lists `leaves` in build order, each with its `needs`. A
link that would close a loop is refused (exit 3, `depends-on from X to Y would close a
loop: ...`) and writes nothing.

## Running it

- `--dry-run` prints the operations it would write and writes nothing.
- `--record` lands the file in a new record change set in one step.
- Otherwise it goes into `--change`, `RHIZA_CHANGE` or the current change set.
- A bad item writes nothing; the error names it (`writes[2]: ...`) and `next` names
  `rhiza describe <op>`.
- Items that cancel out are an error: nothing to write.

`rhiza describe apply --json` gives the full schema, and `rhiza describe add --json` (and
so on) each item's.
