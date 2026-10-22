// Package depgraph is used to create dependency graphs. The specifics are
// abstract so it's easier to understand the motivation. In some applications
// (CAD models were the original motivation) there are a few top level objects
// (or just a single) that are being computed. These depend on lower level
// objects, down to configuration values. Changing a configuration value might
// drastically change the output. And it may be of value to inspect how changes
// effect each object.
//
// The dependency graph encapsulates this logic. An abstract graph is defined in
// code. Each abstract node is a function that returns a value. The node
// function doesn't have arguments passed in, instead it requests them from the
// graph. The graph then either returns cached values or computes new ones as
// needed. It also tracks these requests to build the graph of which objects
// have direct dependancies. This forms an a cyclic graph with the leaves being
// nodes which have no dependancies.
//
// A concrete instance of an abstract graph is a Context. The leaves are Values
// designated in it, and each Context holds its own Values and its own computed
// results. Updating a Value propigates changes to the nodes that depend on it.
//
// Nodes are always added to a Context by name, but can also be added by type.
// Only one node of a type can be found by type in a Context; adding another
// replaces it.
package depgraph
