package symexpr

type ExprParams struct {

	// bounds on tree
	MaxS, MaxD,
	MinS, MinD int

	// bounds on some operators
	NumDim, NumSys, NumCoeff int

	// usable terms at each location
	Roots, Nodes, Leafs, NonTrig []Expr
}

type ExprStats struct {
	depth   int // layers from top (root == 1)
	height  int // layers of subtree (leaf == 1)
	size    int
	numchld int
}

func (es *ExprStats) Depth() int       { return es.depth }
func (es *ExprStats) Height() int      { return es.height }
func (es *ExprStats) Size() int        { return es.size }
func (es *ExprStats) NumChildren() int { return es.numchld }

func (t *Leaf) CalcExprStats(currDepth int) (mySize int) {
	t.depth = currDepth + 1
	t.height = 1
	t.size = 1
	t.numchld = 0
	return t.size
}

func (u *Unary) CalcExprStats(currDepth int) (mySize int) {
	u.depth = currDepth + 1
	u.size = 1 + u.C.CalcExprStats(currDepth+1)
	u.height = 1 + u.C.Height()
	u.numchld = 1
	return u.size
}

func (n *N_ary) CalcExprStats(currDepth int) (mySize int) {
	n.depth = currDepth + 1
	n.size = 1
	n.numchld = 0
	h := 0
	for _, C := range n.CS {
		if C == nil {
			continue
		} else {
			n.numchld++
		}
		n.size += C.CalcExprStats(currDepth + 1)
		if h < C.Height() {
			h = C.Height()
		}
	}
	n.height = 1 + h
	return n.size
}

func (u *PowI) CalcExprStats(currDepth int) (mySize int) {
	u.depth = currDepth + 1
	u.size = 1 + u.Base.CalcExprStats(currDepth+1)
	u.height = 1 + u.Base.Height()
	u.numchld = 1
	return u.size
}

func (u *PowF) CalcExprStats(currDepth int) (mySize int) {
	u.depth = currDepth + 1
	u.size = 1 + u.Base.CalcExprStats(currDepth+1)
	u.height = 1 + u.Base.Height()
	u.numchld = 1
	return u.size
}

func max(l, r int) int {
	if l > r {
		return l
	}
	return r
}

func (n *PowE) CalcExprStats(currDepth int) (mySize int) {
	n.depth = currDepth + 1
	n.size = 1 + n.Base.CalcExprStats(currDepth+1) + n.Power.CalcExprStats(currDepth+1)
	n.height = 1 + max(n.Base.Height(), n.Power.Height())
	n.numchld = 2
	return n.size
}

func (n *Div) CalcExprStats(currDepth int) (mySize int) {
	n.depth = currDepth + 1
	n.size = 1 + n.Numer.CalcExprStats(currDepth+1) + n.Denom.CalcExprStats(currDepth+1)
	n.height = 1 + max(n.Numer.Height(), n.Denom.Height())
	n.numchld = 2
	return n.size
}
