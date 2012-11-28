package kd

import (
	"../geometry"
	"container/list"
	"sort"
)

const (
	DIM_X = 0
	DIM_Y = 1
	DIM_Z = 2
)

// Everything stored in the KD-tree, is a KDValue
type KDValue interface {
	Position() geometry.Vec3
}

// The KDNodes are the nodes in the tree
// It has a value, a splitting dimension and left and right childs.
type KDNode struct {
	Value       KDValue
	Split       int
	Left, Right *KDNode
}

// Make KDNode implement the KDValue interface
func (n KDNode) Position() geometry.Vec3 {
    return n.Value.Position()
}

// Convinience distance function
func (me *KDNode) Distance(other KDValue) float64{
    return me.Position().Distance(other.Position())
}

// Convinience distance^2 function
func (me *KDNode) Distance2(other KDValue) float64{
    return me.Position().Distance2(other.Position())
}

// Extract the correct value from the KDValue to compare on
func comparingValue(item KDValue, dimension int) float64 {
    switch dimension {
	case DIM_X:
		return item.Position().X
	case DIM_Y:
		return item.Position().Y
	case DIM_Z:
		return item.Position().Z
	}
    panic("Trying to get higher dimensional value")
}

///////////////////////////////
// Implement sort.Interface
// Needed to use the sorting library
// in the Go standard library
///////////////////////////////
type valueList struct {
	values    []KDValue
	dimension int
}

func (l valueList) Len() int {
	return len(l.values)
}

func (l valueList) Less(i, j int) bool {
	return comparingValue(l.values[i], l.dimension) < comparingValue(l.values[j], l.dimension)
}

func (l valueList) Swap(i, j int) {
	l.values[i], l.values[j] = l.values[j], l.values[i]
}


// Debugging functions calculating the KD tree depth
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func (k *KDNode) Depth() int {
    if k == nil {
        return 0
    }
    depth := max(k.Left.Depth(), k.Right.Depth())
    
    return 1 + depth;
}

// Function that calls New asynconosly
// returning a channel from which you can read the value later
func AsyncNew(items *list.List, maxDimension int) (<-chan *KDNode) {
    channel := make(chan *KDNode);

    go func() {
        channel <- New(items, maxDimension)
    }()

    return channel
}

// Helper function to conditionally branch with go
func condGo(condition bool, f func()) {
    if condition {
        go f()
    } else {
        f()
    }
}

// Creates a new KD-tree by taking a *list.List of KDValues
// Works by finding the median in every dimension and
// recursivly creating KD-trees as children untill the list is empty.
//
// Uses Go routines and channels to acheive concurrency.
// Every level creates one new Go routine and processes one sub-tree
// on it's own. 
func New(items *list.List, maxDimension int) *KDNode {
	buffer := make([]KDValue, items.Len())

	index := 0
	for e := items.Front(); e != nil; e = e.Next() {
		buffer[index] = e.Value.(KDValue)
		index++
	}

	var create func(l []KDValue, result chan *KDNode, depth int)
	create = func(l []KDValue, result chan *KDNode, depth int) {
		if len(l) == 0 {
			result <- nil
			return
		}

        // Sort the array 
		sort.Sort(valueList{l, depth % maxDimension})
		median := len(l) / 2
        // Adjust the median to make sure it's the FIRST of any
        // identical values
        dimension := depth % maxDimension
        forbiddenValue := comparingValue(l[median], dimension)
        for comparingValue(l[median], dimension) == forbiddenValue && median > 0 {
            median--
        }
		value := l[median]

		left := make(chan *KDNode, 1)
		right := make(chan *KDNode, 1)

        // Branch if high enough in the tree
        condGo(depth < 4, func() { create(l[:median], left, depth+1) })
		create(l[median+1:], right, depth+1)

		result <- &KDNode{value, depth % maxDimension, <-left, <-right}
	}
	node := make(chan *KDNode, 1)
	create(buffer, node, 0)
	return <-node
}

// Searches the tree for any nodes within radius r
// from the target point. This is currently rather slow
// but accurate. By comparing every point to the leftmost
// and rightmost point to the resulting sphere
// irrelevant subtrees are cut of.
func (tree *KDNode) Neighbors(point KDValue, r float64) *list.List {
    if tree == nil {
        return nil
    }
    result := list.New()
    
    // Am I part of the sphere?
    // Compare Distance² to r² to avoid calling sqrt
    if tree.Distance2(point) < r*r {
        result.PushBack(tree)
    }
    
    split := tree.Split
    // Is the leftmost point to the left of us?
    if comparingValue(tree, split) > comparingValue(point, split) - r {
        if left := tree.Left.Neighbors(point, r); left != nil {
            result.PushBackList(left)
        }
    }

    // Is the rightmost point to the right of us?
    if comparingValue(tree, split) < comparingValue(point, split) + r {
        if right := tree.Right.Neighbors(point, r); right != nil {
            result.PushBackList(right)
        }
    }
    // Return all the found nodes
    return result;
}
