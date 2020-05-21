package graphly

// NOTES
// not extensible

// opts := graphly.NewLayeredOptions()
// has all defaults set
// opts.SomeOpt = ABC

// graphly.Layered(graph, opts)
// graphly.Force(graph, opts)

// graph class
// processors and phases mutate the graph object in place

// * A layout processor processes a graph. Layout processors are the secondary components of layout algorithms, the
// * primary being layout phases. Layout processors are inserted before or after phases to do further
// * processing on a graph.
// *
// * The AlgorithmAssembler class can be used to build algorithms by specifying phases and letting the assembler
// * worry about instantiating all required processors.

// public interface ILayoutProcessor<G> {
//    void process(G graph, IElkProgressMonitor progressMonitor);
// }

//     * Rebuilds the configuration to include all processors required to layout the given graph. The list
//     * of processors is attached to the graph in the {@link InternalProperties#PROCESSORS} property.
//     *
//    public void prepareGraphForLayout(final LGraph lgraph) {
//        // Make sure the graph properties are sensible
//        configureGraphProperties(lgraph);
//
//        // Setup the algorithm assembler
//        algorithmAssembler.reset();
//
// Main algorithm phases

//        algorithmAssembler.setPhase(LayeredPhases.P1_CYCLE_BREAKING,
//                lgraph.getProperty(LayeredOptions.CYCLE_BREAKING_STRATEGY));
//        algorithmAssembler.setPhase(LayeredPhases.P2_LAYERING,
//                lgraph.getProperty(LayeredOptions.LAYERING_STRATEGY));
//        algorithmAssembler.setPhase(LayeredPhases.P3_NODE_ORDERING,
//                lgraph.getProperty(LayeredOptions.CROSSING_MINIMIZATION_STRATEGY));
//        algorithmAssembler.setPhase(LayeredPhases.P4_NODE_PLACEMENT,
//                lgraph.getProperty(LayeredOptions.NODE_PLACEMENT_STRATEGY));
//        algorithmAssembler.setPhase(LayeredPhases.P5_EDGE_ROUTING,
//                EdgeRouterFactory.factoryFor(lgraph.getProperty(LayeredOptions.EDGE_ROUTING)));
//
//        algorithmAssembler.addProcessorConfiguration(getPhaseIndependentLayoutProcessorConfiguration(lgraph));
//
//        lgraph.setProperty(InternalProperties.PROCESSORS, algorithmAssembler.build(lgraph));
//    }
