package queue

// Keeping below as var so it is possible to run the slice size bench tests with no coding changes.
var (
	// firstSliceSize holds the size of the first slice.
	firstSliceSize = 1

	// maxFirstSliceSize holds the maximum size of the first slice.
	maxFirstSliceSize = 16

	// maxInternalSliceSize holds the maximum size of each internal slice.
	maxInternalSliceSize = 128
)

// Queue represents an unbounded, dynamically growing FIFO queue.
// The zero value for queue is an empty queue ready to use.
type Queue struct {
	// Head points to the first node of the linked list.
	head *Node

	// Tail points to the last node of the linked list.
	// In an empty queue, head and tail points to the same node.
	tail *Node

	// Hp is the index pointing to the current first element in the queue
	// (i.e. first element added in the current queue values).
	hp int

	// Len holds the current queue values length.
	len int

	// lastSliceSize holds the size of the last created internal slice.
	lastSliceSize int
}

// Node represents a queue node.
// Each node holds a slice of user managed values.
type Node struct {
	// v holds the list of user added values in this node.
	v []interface{}

	// n points to the next node in the linked list.
	n *Node
}

// New returns an initialized queue.
func New() *Queue {
	return new(Queue).Init()
}

// Init initializes or clears queue q.
func (q *Queue) Init() *Queue {
	q.head = nil
	q.tail = nil
	q.hp = 0
	q.len = 0
	return q
}

// Len returns the number of elements of queue q.
// The complexity is O(1).
func (q *Queue) Len() int { return q.len }

// Front returns the first element of queue q or nil if the queue is empty.
// The second, bool result indicates whether a valid value was returned;
//   if the queue is empty, false will be returned.
// The complexity is O(1).
func (q *Queue) Front() (interface{}, bool) {
	if q.head == nil {
		return nil, false
	}
	return q.head.v[q.hp], true
}

// Push adds a value to the queue.
// The complexity is O(1).
func (q *Queue) Push(v interface{}) {
	if q.head == nil {
		h := newNode(firstSliceSize)
		q.head = h
		q.tail = h
		q.lastSliceSize = maxFirstSliceSize
	} else if len(q.tail.v) >= q.lastSliceSize {
		n := newNode(maxInternalSliceSize)
		q.tail.n = n
		q.tail = n
		q.lastSliceSize = maxInternalSliceSize
	}

	q.tail.v = append(q.tail.v, v)
	q.len++
}

// Pop retrieves and removes the current element from the queue.
// The second, bool result indicates whether a valid value was returned;
// 	if the queue is empty, false will be returned.
// The complexity is O(1).
func (q *Queue) Pop() (interface{}, bool) {
	if q.head == nil {
		return nil, false
	}

	v := q.head.v[q.hp]
	q.head.v[q.hp] = nil // Avoid memory leaks
	q.len--
	q.hp++
	if q.hp >= len(q.head.v) {
		n := q.head.n
		q.head.n = nil // Avoid memory leaks
		q.head = n
		q.hp = 0
	}
	return v, true
}

// newNode returns an initialized node.
func newNode(capacity int) *Node {
	return &Node{
		v: make([]interface{}, 0, capacity),
	}
}
