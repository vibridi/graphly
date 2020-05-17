package graphly

const (
	phase1_CYCLE_BREAKING = iota
	phase2_LAYERING
	phase3_ORDERING
	phase4_POSITIONING
	phase5_ROUTING
)

type CycleBreakingStrategy uint8

const (
	CycleBreakingStrategy_GREEDY CycleBreakingStrategy = iota
	CycleBreakingStrategy_DEPTH_FIRST
)

// Maybe unexported. Build with all the defaults
// Could also be zero values
type LayeredOptions struct {
	totalPhases           uint8
	CycleBreakingStrategy CycleBreakingStrategy
}

func NewLayeredOptions() LayeredOptions {
	return LayeredOptions{
		// set the defaults
		totalPhases: 5,
	}
}

func Layered(graph *Graph, opts LayeredOptions) {
	layout := layeredLayout{opts}

	// algorithm assembler
	asm := newAlgorithmAssembler(opts.totalPhases)
	asm.addPhase(phase1_CYCLE_BREAKING, layout.phase1Factory())

	// for p in asm.algorithm [list of processors]
	// p(graph)

	for _, f := range asm.algorithm() {
		f(graph)
	}
}

type layeredLayout struct {
	opts LayeredOptions
}

func (l *layeredLayout) algorithmStrategy() {

}
