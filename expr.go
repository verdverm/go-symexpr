// Package symexpr implements symbolic equations as an AST.
// It aims to provide ease of dynamic manipulation of the tree.
//
// This work comes out of my Masters thesis at Binghamton University
// and is geared towards Symbolic Regression.
//
package symexpr

import (
	"math"
)

type exprType int

const (
	NULL exprType = iota
	CONSTANT
	CONSTANTF
	TIME
	SYSTEM
	VAR

	NEG
	ABS
	SQRT
	SIN
	COS
	TAN
	EXP
	LOG
	POWI
	POWF

	POWE
	DIV

	ADD
	MUL

	EXPR_MAX
	STARTVAR // for serialization reduction of variables
)

// Expr is the interface to all node types for the AST of mathematical expression
// 
type Expr interface {

	// types.go (this file)
	ExprType() exprType
	Clone() Expr

	// stats.go
	Size() int
	Depth() int
	Height() int
	NumChildren() int
	CalcExprStats(currDepth int) (mySize int)

	// compare.go
	AmILess(rhs Expr) bool
	AmIEqual(rhs Expr) bool
	AmISame(rhs Expr) bool       // equality without coefficient values/index
	AmIAlmostSame(rhs Expr) bool // adds flexibility to mul comparison to AmISame
	Sort()

	// has.go
	HasVar() bool
	HasVarI(i int) bool
	NumVar() int

	// DFS for a floating point valued ConstantF
	HasConst() bool
	// DFS for a indexed valued Constant
	HasConstI(i int) bool
	// Counts the number of indexed Constant nodes
	NumConstants() int

	// convert.go

	// Converts indexed Constant nodes to ConstantF nodes
	// using the input slice as the values for replacement 
	ConvertToConstantFs(cs []float64) Expr
	// DFS converting float valued constants to indexed constants
	// the input should be an empty slice
	// the output is an appended slice the size = |ConstantF| in the tree
	ConvertToConstants(cs []float64) ([]float64, Expr)
	//   IndexConstants( ci int ) int

	// getset.go   
	// DFS retrieval of a node by index
	GetExpr(pos *int) Expr
	// DFS replacement of a node and it's subtree
	// replaced is used to discontinue the DFS after replacement
	// replace_me gets triggered when pos == 0 and informs the parent node to replace the respective child node
	SetExpr(pos *int, e Expr) (replace_me, replaced bool)

	// print.go

	// prints the AST 
	String() string

	// creates an integer representation of the AST in ~prefix notation
	// The input is an empty slice, output is the representation.
	// The output is generally the ExprType integer value
	// Associative operators (+ & *) also include the number of children.
	// The terminal nodes include the index when appropriate.
	Serial([]int) []int

	// Pretty print acts like String, but replaces the internal indexed
	// formatting with user specified strings and values
	PrettyPrint(dnames, snames []string, cvals []float64) string
	// 	WriteString( buf *bytes.Buffer )

	// eval.go
	// Evaluates an expression at one point
	// t is a time value
	// x are the input Var values
	// c are the indexed Constant values
	// s are the indexed System values
	// the output is the result of DFS evaluation
	Eval(t float64, x, c, s []float64) float64

	// simp.go
	Simplify(rules SimpRules) Expr

	// deriv.go
	// Calculate the derivative w.r.t. Var_i
	DerivVar(i int) Expr
	// Calculate the derivative w.r.t. Constant_i
	DerivConst(i int) Expr
}

type ExprArray []Expr

func (p ExprArray) Len() int      { return len(p) }
func (p ExprArray) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ExprArray) Less(i, j int) bool {
	if p[i] == nil && p[j] == nil {
		return false
	}
	if p[i] == nil {
		return false
	}
	if p[j] == nil {
		return true
	}
	return p[i].AmILess(p[j])
}

// Null Leaf  (shouldn't really appear)
// This is a sample for the other node types
type Null struct {
	ExprStats
}

func NewNull() Expr                { return new(Null) }
func (n *Null) ExprType() exprType { return NULL }
func (n *Null) Clone() Expr        { return &Null{ExprStats{0, 0, 0, 0}} }

func (n *Null) CalcExprStats(currDepth int) (mySize int) {
	n.depth = currDepth + 1
	n.height = 0
	n.size = 0
	n.numchld = 0
	return n.size
}

func (n *Null) AmILess(r Expr) bool       { return NULL < r.ExprType() }
func (n *Null) AmIEqual(r Expr) bool      { return r.ExprType() == NULL }
func (n *Null) AmISame(r Expr) bool       { return r.ExprType() == NULL }
func (n *Null) AmIAlmostSame(r Expr) bool { return r.ExprType() == NULL }
func (n *Null) Sort()                     { return }

func (n *Null) HasVar() bool         { return false }
func (n *Null) HasVarI(i int) bool   { return false }
func (n *Null) NumVar() int          { return 0 }
func (n *Null) HasConst() bool       { return false }
func (n *Null) HasConstI(i int) bool { return false }
func (n *Null) NumConstants() int    { return 0 }

func (n *Null) ConvertToConstantFs(cs []float64) Expr             { return n }
func (n *Null) ConvertToConstants(cs []float64) ([]float64, Expr) { return cs, n }

func (n *Null) GetExpr(pos *int) Expr {
	if (*pos) == 0 {
		return n
	}
	(*pos)--
	return nil
}
func (n *Null) SetExpr(pos *int, e Expr) (replace_me, replaced bool) {
	if (*pos) == 0 {
		return true, false
	}
	(*pos)--
	return false, false
}

func (n *Null) String() string                                              { return "NULL" }
func (n *Null) Serial(sofar []int) []int                                    { return append(sofar, int(NULL)) }
func (n *Null) PrettyPrint(dnames, snames []string, cvals []float64) string { return "NULL" }

func (n *Null) Eval(t float64, x, c, s []float64) float64 { return math.NaN() }

func (n *Null) Simplify(rules SimpRules) Expr { return n }

func (n *Null) DerivConst(i int) Expr { return &ConstantF{F: 0} }
func (n *Null) DerivVar(i int) Expr   { return &ConstantF{F: 0} }
