package layered

type processor interface {
	process(graph *Graph)
}

// Represents a major phase of an autolayout algorithm,
// including pre- and post-processing steps. Implements process interface
type phase struct {
	main           processor
	preprocessors  []processor
	postprocessors []processor
}

func (p *phase) process(graph *Graph) {
	for _, p := range p.preprocessors {
		p.process(graph)
	}
	p.main.process(graph)
	for _, p := range p.postprocessors {
		p.process(graph)
	}
}

type assembler struct {
	phases []*phase
}

func newAlgorithmAssembler(totalPhases uint8) *assembler {
	return &assembler{
		phases: make([]*phase, totalPhases),
	}
}

func (asm *assembler) addPhase(phase uint8, proc processor) {
	asm.phases[phase].main = proc
}

func (asm *assembler) addPreProcessor(proc processor, beforePhase uint8) {
	pre := &asm.phases[beforePhase].preprocessors
	*pre = append(*pre, proc)
}

func (asm *assembler) addPostProcessor(proc processor, afterPhase uint8) {
	post := &asm.phases[afterPhase].postprocessors
	*post = append(*post, proc)
}

func (asm *assembler) algorithm() []processor {
	procs := make([]processor, len(asm.phases))
	for i, phase := range asm.phases {
		procs[i] = phase
	}
	return procs
}
