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
		totalPhases:           5,
		CycleBreakingStrategy: CycleBreakingStrategy_GREEDY,
	}
}

func Layered(graph *Graph, opts LayeredOptions) {
	// algorithm assembler
	asm := newAlgorithmAssembler(opts.totalPhases)
	asm.addPhase(phase1_CYCLE_BREAKING, layeredPhase1Factory(opts.CycleBreakingStrategy))
	// todo add post-processing step restore reversed edges after p5

	// for p in asm.algorithm [list of processors]
	// p(graph)

	for _, processor := range asm.algorithm() {
		processor.process(graph)
	}
}
