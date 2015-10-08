package gtreap

type Treap struct {
	compare Compare
	root    *node
}

// Compare returns an integer comparing the two items
// lexicographically. The result will be 0 if a==b, -1 if a < b, and
// +1 if a > b.
type Compare func(a, b interface{}) int

// Item can be anything.
type Item interface{}

type node struct {
	item     Item
	priority int
	left     *node
	right    *node
	count    int
}

func (n *node) getCount() int {
	if n == nil {
		return 0
	}
	return n.count
}

func NewTreap(c Compare) *Treap {
	return &Treap{compare: c, root: nil}
}

func (t *Treap) Len() int {
	if t.root == nil {
		return 0
	}
	return t.root.count
}

func (t *Treap) Min() Item {
	n := t.root
	if n == nil {
		return nil
	}
	for n.left != nil {
		n = n.left
	}
	return n.item
}

func (t *Treap) Max() Item {
	n := t.root
	if n == nil {
		return nil
	}
	for n.right != nil {
		n = n.right
	}
	return n.item
}

func (t *Treap) Get(target Item) Item {
	n := t.root
	for n != nil {
		c := t.compare(target, n.item)
		if c < 0 {
			n = n.left
		} else if c > 0 {
			n = n.right
		} else {
			return n.item
		}
	}
	return nil
}

func (t *Treap) Upsert(item Item, itemPriority int) *Treap {
	r := t.union(t.root, &node{item: item, priority: itemPriority, count: 1})
	return &Treap{compare: t.compare, root: r}
}

func (t *Treap) union(this *node, that *node) *node {
	if this == nil {
		return that
	}
	if that == nil {
		return this
	}
	if this.priority > that.priority {
		left, middle, right := t.split(that, this.item)
		if middle == nil {
			uLeft, uRight := t.union(this.left, left), t.union(this.right, right)
			return &node{
				item:     this.item,
				priority: this.priority,
				left:     uLeft,
				right:    uRight,
				count:    1 + uLeft.getCount() + uRight.getCount(),
			}
		}
		uLeft, uRight := t.union(this.left, left), t.union(this.right, right)
		return &node{
			item:     middle.item,
			priority: middle.priority,
			left:     uLeft,
			right:    uRight,
			count:    1 + uLeft.getCount() + uRight.getCount(),
		}
	}
	// We don't use middle because the "that" has precendence.
	left, _, right := t.split(this, that.item)
	uLeft, uRight := t.union(left, that.left), t.union(right, that.right)
	return &node{
		item:     that.item,
		priority: that.priority,
		left:     uLeft,
		right:    uRight,
		count:    1 + uLeft.getCount() + uRight.getCount(),
	}
}

// Splits a treap into three treaps based on a split item "s".
// The result tuple-3 means (left, X, right), where X is either...
// nil - meaning the item s was not in the original treap.
// non-nil - returning a treap with a single node having item s.
// The tuple-3's left result has items < s,
// and the tuple-3's right result has items > s.
func (t *Treap) Split(s Item) (*Treap, *Treap, *Treap) {
	nleft, nmiddle, nright := t.split(t.root, s)

	left := &Treap{compare: t.compare, root: nleft}
	var middle *Treap
	if nmiddle != nil {
		middle = &Treap{compare: t.compare, root: &node{
			item:     nmiddle.item,
			priority: nmiddle.priority,
			left:     nil,
			right:    nil,
			count:    1,
		}}
	}
	right := &Treap{compare: t.compare, root: nright}
	return left, middle, right
}

func (t *Treap) split(n *node, s Item) (*node, *node, *node) {
	if n == nil {
		return nil, nil, nil
	}
	c := t.compare(s, n.item)
	if c == 0 {
		return n.left, n, n.right
	}
	if c < 0 {
		left, middle, right := t.split(n.left, s)
		return left, middle, &node{
			item:     n.item,
			priority: n.priority,
			left:     right,
			right:    n.right,
			count:    1 + right.getCount() + n.right.getCount(),
		}
	}
	left, middle, right := t.split(n.right, s)
	return &node{
		item:     n.item,
		priority: n.priority,
		left:     n.left,
		right:    left,
		count:    1 + n.left.getCount() + left.getCount(),
	}, middle, right
}

func (t *Treap) Delete(target Item) *Treap {
	left, _, right := t.split(t.root, target)
	return &Treap{compare: t.compare, root: t.join(left, right)}
}

// All the items from this are < items from that.
func (t *Treap) join(this *node, that *node) *node {
	if this == nil {
		return that
	}
	if that == nil {
		return this
	}
	if this.priority > that.priority {
		right := t.join(this.right, that)
		return &node{
			item:     this.item,
			priority: this.priority,
			left:     this.left,
			right:    right,
			count:    1 + this.left.getCount() + right.getCount(),
		}
	}
	left := t.join(this, that.left)
	return &node{
		item:     that.item,
		priority: that.priority,
		left:     left,
		right:    that.right,
		count:    1 + left.getCount() + that.right.getCount(),
	}
}

type ItemVisitor func(i Item) bool

// Visit items greater-than-or-equal to the pivot.  If the pivot is null, this visits all items in ascending order.
func (t *Treap) VisitAscend(pivot Item, visitor ItemVisitor) {
	t.visitAscend(t.root, pivot, visitor)
}

func (t *Treap) visitAscend(n *node, pivot Item, visitor ItemVisitor) bool {
	if n == nil {
		return true
	}
	if pivot == nil || t.compare(pivot, n.item) <= 0 {
		if !t.visitAscend(n.left, pivot, visitor) {
			return false
		}
		if !visitor(n.item) {
			return false
		}
		//since n.right is > n.item by the comparison, we can speed up by not comparing
		return t.visitAscend(n.right, nil, visitor)
	}
	return t.visitAscend(n.right, pivot, visitor)
}

// Visit items less-than-or-equal to the pivot.  If the pivot is null, this visits all items in descending order.
func (t *Treap) VisitDescend(pivot Item, visitor ItemVisitor) {
	t.visitDescend(t.root, pivot, visitor)
}

func (t *Treap) visitDescend(n *node, pivot Item, visitor ItemVisitor) bool {
	if n == nil {
		return true
	}
	if pivot == nil || t.compare(pivot, n.item) >= 0 {
		if !t.visitDescend(n.right, pivot, visitor) {
			return false
		}
		if !visitor(n.item) {
			return false
		}
		//since n.left is < n.item by the comparison, we can speed up by not comparing
		return t.visitDescend(n.left, nil, visitor)
	}
	return t.visitDescend(n.left, pivot, visitor)
}
