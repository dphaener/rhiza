# genesis-scrub

Turn a raw conversation (a chat transcript, meeting notes, an interview) into a Source node.
The Source is committed to the repo; the raw transcript never is.

## Steps

1. Read the whole transcript. Find who owns the decisions (the product owner, the person
   whose call it is) and every place something was decided, rejected or left open.
2. Summarise. Write a short `summary`: what the conversation was about, what was decided,
   what is still open. Plain sentences, no filler.
3. Quote every decision. Wherever the owner decided something, copy their exact words into
   `quotes`. Never paraphrase a decision, never tidy its grammar, never merge two quotes.
   "Yes, but custom kinds are down the roadmap" is worth more verbatim than summarised.
   Quote only what the owner typed: their messages, typed answers and notes. An option an agent
   wrote (an AskUserQuestion choice, a "(Recommended)" pick) is never the owner's words, even
   when the owner picked it. Quote it with `by` naming the agent and saying the owner accepted
   it, e.g. `"by": "agent (accepted by the owner)"`.
4. Remove what must not be committed:
   - secrets: API keys, tokens, passwords, connection strings, private URLs;
   - personal data: emails, phone numbers, addresses, names of people outside the project
     (replace with a role: "a customer", "the designer");
   - noise: greetings, tool output, retries, tangents that decided nothing.
   If a decision quote contains a secret or personal data, cut that part and mark the cut
   with `[...]`; keep the rest verbatim.
5. Write the source file, JSON:

```json
{
  "title": "Export formats, with the owner",
  "summary": "Agreed CSV export streams rows. Encoding is still open.",
  "quotes": [
    {"text": "Stream it, reports get huge.", "by": "owner", "where": "2026-10-20 call"}
  ]
}
```

   `title`, `summary` and `quotes` are required; `quotes` may be empty. In a quote, `text` is
   required; `by` (who said it, a role or handle) and `where` (when or where in the material)
   are optional.
6. Store it:
   - inside a change set: `genesis source add <handle> <file.json> --raw <transcript> --json`
   - on its own: `genesis source add <handle> <file.json> --raw <transcript> --record --json`
   `--raw` hashes the transcript (sha256) and stores only the hash. If the transcript is not on
   this machine, pass `--raw-sha256 <hex>` instead. Check first with `--dry-run`.
7. Link what came from it: decisions and capabilities derived from the conversation get a
   `derived-from` link to the source (`genesis link <node> derived-from <handle> --json`, or in
   an `apply` file). Pass `--source <handle>` to `genesis change new` when a change set comes
   from it.

## Rules

- Never commit the raw transcript, and never paste it into a field. Keep it and the scratch
  `<file.json>` out of git; `source add` stores what is needed.
- Keep the owner's words exactly where they decided something; summarise everything else.
- A picked option is the agent's wording, accepted by the owner, never the owner's own words.
- When unsure whether something is sensitive, leave it out.
- Handles are lower-case words joined by dashes, such as `export-call-2026-10-20`.

`genesis guide sources` covers the stored fields.
