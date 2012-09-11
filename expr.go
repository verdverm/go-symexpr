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
  STARTVAR   // for serialization reduction of variables

)


type Expr interface {

  // types.go (this file)
  ExprType() exprType
  Clone() Expr


  // stats.go
  Size() int
  Depth() int
  Height() int
  NumChildren() int
  CalcExprStats( currDepth int ) (mySize int)

  // compare.go
  AmILess( rhs Expr ) bool
  AmIEqual( rhs Expr ) bool
  AmISame( rhs Expr ) bool        // equality without coefficient values/index
  AmIAlmostSame( rhs Expr ) bool  // adds flexibility to mul comparison to AmISame
  Sort()

  // has.go
  HasVar() bool
  HasVarI(i int) bool
  NumVar() int

  HasConst() bool
  HasConstI(i int) bool
  NumConstants() int // only Constant NOT ConstantF

  // convert.go
  ConvertToConstantFs( cs []float64 ) Expr
  ConvertToConstants( cs[]float64 ) ([]float64, Expr)
  //   IndexConstants( ci int ) int

  // getset.go   DFS retrieval
  GetExpr( pos *int ) Expr
  SetExpr( pos *int , e Expr ) (bool,bool)

  // print.go
  String() string
  Serial([]int) []int
  PrettyPrint( dnames,snames []string, cvals []float64 ) string
  // 	WriteString( buf *bytes.Buffer )

  // eval.go
  Eval(t float64, x, c, s []float64) float64

  // simp.go
  Simplify( rules SimpRules ) Expr

  // deriv.go
  DerivVar( i int ) Expr
  DerivConst( i int ) Expr

}

type ExprArray []Expr
func (p ExprArray) Len() int { return len(p) }
func (p ExprArray) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ExprArray) Less(i, j int) bool {
  if p[i] == nil && p[j] == nil { return false }
  if p[i] == nil { return false }
  if p[j] == nil { return true 	}
  return p[i].AmILess( p[j] )
}




// SAMPLE VIA NULL


// Null Leaf  (shouldn't appear ever... if code is correct... )
type Null struct{
  ExprStats
}
func NewNull() Expr { return new(Null) }
func (n *Null) ExprType() exprType   { return NULL }
func (n *Null) Clone() Expr     { return &Null{ExprStats{0,0,0,0}} }

func (n *Null) CalcExprStats( currDepth int ) (mySize int) {
  n.depth = currDepth+1
  n.height = 0
  n.size = 0
  n.numchld = 0
  return n.size
}

func (n *Null) AmILess( r Expr ) bool       { return NULL < r.ExprType() }
func (n *Null) AmIEqual( r Expr ) bool       { return r.ExprType() == NULL }
func (n *Null) AmISame( r Expr ) bool       { return r.ExprType() == NULL }
func (n *Null) AmIAlmostSame( r Expr ) bool { return r.ExprType() == NULL }
func (n *Null) Sort()                       { return }


func (n *Null) HasVar() bool            { return false }
func (n *Null) HasVarI( i int ) bool    { return false }
func (n *Null) NumVar() int             { return 0 }
func (n *Null) HasConst() bool          { return false }
func (n *Null) HasConstI( i int ) bool  { return false }
func (n *Null) NumConstants() int       { return 0 }

func (n *Null) ConvertToConstantFs( cs []float64 ) Expr     { return n }
func (n *Null) ConvertToConstants( cs []float64 ) ( []float64, Expr )     { return cs,n }

func (n *Null) GetExpr( pos *int ) Expr     { if (*pos) == 0 { return n }; (*pos)--; return nil }
func (n *Null) SetExpr( pos *int, e Expr ) ( replace_me, replaced bool )     { if (*pos) == 0 { return true,false }; (*pos)--; return false,false }

func (n *Null) String() string     { return "NULL" }
func (n *Null) Serial( sofar []int ) []int { return append(sofar,int(NULL)) }
func (n *Null) PrettyPrint( dnames,snames []string, cvals []float64 ) string { return "NULL" }

func (n *Null) Eval(t float64, x, c, s []float64) float64        { return math.NaN() }

func (n *Null) Simplify( rules SimpRules ) Expr      { return n }

func (n *Null) DerivConst( i int) Expr     { return &ConstantF{F: 0} }
func (n *Null) DerivVar( i int) Expr     { return &ConstantF{F: 0} }
