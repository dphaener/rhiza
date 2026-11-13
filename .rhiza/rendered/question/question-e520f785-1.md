# question-e520f785-1

question · `01M7AHJAA6T4TC41CBJ52D8V2N`

## session

e520f785-7bad-4af3-ae27-01a1d514808d

## statement

Change set 'Every agent session a run' and its sibling 'The daemon offers spawn, gate and hand-back primitives' (leaf run-handback) both say they write handback.md, and nothing hands back yet. Who owns handback.md? / The milestone-36 interview left brief.md's exact format open. What goes in a session's brief.md? / Where do a session's brief.md, slice.json and handback.md live? There is no per-run directory today: only the journal <commondir>/rhiza/runs/<id>.jsonl and the run's temp folder, which is removed when the run ends. / Where does the new 'invoke' key apply in a workflow file?

## Links

### In

- answers [brief-files-and-invoke-placement](../decision/brief-files-and-invoke-placement.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
