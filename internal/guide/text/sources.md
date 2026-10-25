# Sources

A Source node is where an idea came from: a conversation, a meeting, an interview. Decisions
and capabilities link to it, so `genesis why` can show the owner's own words behind them.

## Two different things

- **The Source node** (committed): a scrubbed summary, the owner's verbatim quotes, and the
  sha256 of the raw material.
- **The raw material** (never committed, never stored): the transcript itself. Keep it on your
  machine or throw it away. Its hash proves which material the summary was made from.

## Stored fields

| field | what |
|---|---|
| `title` | a short name for the conversation |
| `summary` | what was discussed, decided and left open, scrubbed |
| `quotes` | `[{"text", "by", "where"}]`: exact words, mostly the owner's decisions |
| `raw_sha256` | the hex sha256 of the raw material |

## Adding one

1. Scrub the conversation with the genesis-scrub skill (`genesis guide scrub`). It writes a
   JSON file `{"title", "summary", "quotes":[{"text", "by", "where"}]}`; `by` and `where` are
   optional, and `quotes` may be empty.
2. `genesis source add <handle> <file.json> --raw <transcript> --json` inside a change set, or
   with `--record` to land it in one step. `--raw-sha256 <hex>` gives the hash when the
   transcript is elsewhere. `--dry-run` shows what it would write.
3. Link what came from it: `derived-from` from a decision or capability to the source, and
   `--source <handle>` on `genesis change new` for a change set that comes from it.

## Rules

- A decision's quote is copied, never paraphrased.
- Secrets and personal data are removed before `source add`; the log is append-only, so what
  goes in stays in the history.
