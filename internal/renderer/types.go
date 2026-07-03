package renderer

type SegmentKind int

const (
	SegmentWord SegmentKind = iota
	SegmentSpace
)

type RenderedSegment struct {
	Kind SegmentKind
	Text string
	Rows []string
}

type RenderedLine struct {
	Segments []RenderedSegment
}

type RenderedASCII []RenderedLine
