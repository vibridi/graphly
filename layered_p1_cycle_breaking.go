package graphly

import (
	"errors"
	"fmt"
)

func (l *layeredLayout) phase1Factory() processor {
	strat := l.opts.CycleBreakingStrategy
	switch strat {
	case CycleBreakingStrategy_GREEDY:
		return l.greedyCycleBreaker

	case CycleBreakingStrategy_DEPTH_FIRST:
		return l.depthFirstCycleBreaker
	}
	panic(errors.New(fmt.Sprintf("unknown cycle breaking strategy: %d", strat)))
}

func (l *layeredLayout) greedyCycleBreaker(graph *Graph) {

}

func (l *layeredLayout) depthFirstCycleBreaker(graph *Graph) {

}
