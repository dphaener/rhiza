# overlap-retire-link-verified

evidence · `01M6A89T72D5HYQ3879E3C8JQW` · actor [claude-code-56bde82d](../actor/claude-code-56bde82d.md)

## result

PASS

## statement

At bf373189, go test ./internal/model -run TestOverlaps -count=1 -v passes: TestOverlaps has new cases 'retire vs link create to the node' and 'retire vs link create from the node' (each wants one Overlap on the retired node with aspects [linked retired]), plus negative cases 'retire and link in one set', 'two link creates to one node' and 'link create vs update' (no overlap); TestOverlapsInputOrder (order-independence over every permutation) still passes. go test ./... -count=1 passes. Diff read: Overlaps indexes each link create under both endpoints with aspect 'linked', which overlapAspects never reports as a shared aspect, so it only overlaps a retirement.

## Links

### Out

- verifies [overlap-retire-link](../capability/rhiza/change-sets/conflicts/overlap-retire-link.md)

## History

- `01M5SWTNWB9GEJ5AA8R7K867KB` Genesis detects conflicting decisions across change sets in flight and surfaces them to the owner.
