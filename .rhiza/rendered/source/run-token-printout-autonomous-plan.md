# run-token-printout-autonomous-plan

source · `01M65055ZY7G2C54ETJZGACA4G`

## quotes

```json
[]
```

## raw_sha256

1e33afcc9117a64cdd7e0ef5821d04afc683e930148c11d4c4fedec8a657eb5d

## session

4ba1eee9-94f7-4b51-8950-01a18a012c9c

## summary

No owner answered this plan run. The change set's intent is that a run's final printout and run status show each session's token use and never say 'no token use recorded' for a session that used tokens. Today the printout (internal/api/run.go, the Sessions section) prints that line for any session with no journaled token reading. Readings are journaled only when a session says done, its step ends or it exits, so a live session, or one whose transcript could not be read, shows the line though it has used tokens. Open: whether the printout and run status read the transcript at print time when no reading is journaled. Raised as a question; no owner words are recorded.

## title

Run token printout: autonomous plan plus arbiter-raised question

## Links

### In

- derived-from [token-gap-read-at-print](../question/token-gap-read-at-print.md)
- derived-from [printout-reads-transcript-when-unjournaled](../decision/printout-reads-transcript-when-unjournaled.md)

## History

- `01M61GW8RKSBM4BT1B8B3FR5JH` A run's final printout and run status show each session's token use, never 'no token use recorded' for sessions that used tokens.
