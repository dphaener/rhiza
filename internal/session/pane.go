package session

// The directions a pane splits in: the new pane to the anchor's right, or
// below it.
const (
	SplitRight = "right"
	SplitBelow = "below"
)

// SplitFor is the direction to split a pane of widthPx by heightPx pixels
// in, so the new pane takes half of its longer side: right when it is at
// least as wide as it is tall, below when it is taller. An override of
// right or below wins over the shape; any other override (auto, "") lets
// the shape decide. A size that is not known (zero or less on either side)
// splits right.
func SplitFor(widthPx, heightPx int, override string) string {
	switch override {
	case SplitRight, SplitBelow:
		return override
	}
	if widthPx <= 0 || heightPx <= 0 || widthPx >= heightPx {
		return SplitRight
	}
	return SplitBelow
}
