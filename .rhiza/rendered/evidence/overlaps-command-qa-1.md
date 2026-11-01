# overlaps-command-qa-1

evidence · `01M6AEBRWGM0Q17JPRVT26FX5V` · actor [claude-code-e44db580](../actor/claude-code-e44db580.md)

## result

PASS

## statement

At 04deb2ed I built rhiza with go build -o $TMPDIR/rhiza-qa ./cmd/rhiza and used it in a scratch git repo: landed A adding capabilities base and other, proposed B setting base.statement, drafts C (base.statement, other.acceptance) and D (other.acceptance, current). rhiza overlaps printed 'No overlaps' and --json {"overlaps": []} with D current, so the current draft is not counted unless named. --change C printed 'change sets B and C both write field:statement of base' and JSON with a, b, target base and aspects [field:statement]; rhiza api call overlaps --change C gave byte-identical output. --change D gave No overlaps (unnamed draft C ignored), --change B (proposed) gave No overlaps, a 10-char prefix and RHIZA_CHANGE both resolved C, an unknown ID exited 4 with 'no change set'. After proposing C, --change D listed B/C on base and C/D on field:acceptance of other. go test of internal/model Overlap and the overlaps and describe-golden scripts pass.

## Links

### Out

- verifies [overlaps-command](../capability/rhiza/views/overlaps-command.md)

## History

- `01M5SWTMR4G5EFRVRQHSEKVA15` Genesis reviews each plan against the model and the other runs in flight before work starts, as an agent check on the plan rather than the code.
