package graphly

type processor func(graph *Graph)

type processorList []processor

func (pl *processorList) push(p processor) {
	if pl == nil {
		n := make(processorList, 0)
		pl = &n
	}
	*pl = append(*pl, p)
}

type phase struct {
	main           processor
	preprocessors  processorList
	postprocessors processorList
}

func (p *phase) run(graph *Graph) {
	for _, p := range p.preprocessors {
		p(graph)
	}
	p.main(graph)
	for _, p := range p.postprocessors {
		p(graph)
	}
}

type assembler struct {
	phases []phase
}

func newAlgorithmAssembler(totalPhases uint8) *assembler {
	return &assembler{
		phases: make([]phase, totalPhases),
	}
}

func (asm *assembler) addPhase(phase uint8, proc processor) {
	asm.phases[phase].main = proc
}

func (asm *assembler) addPreProcessor(proc processor, beforePhase uint8) {
	asm.phases[beforePhase].preprocessors.push(proc)
}

func (asm *assembler) addPostProcessor(proc processor, afterPhase uint8) {
	asm.phases[afterPhase].postprocessors.push(proc)
}

func (asm *assembler) algorithm() []processor {
	procs := make([]processor, len(asm.phases))
	for i, phase := range asm.phases {
		procs[i] = phase.run
	}
	return procs
}
