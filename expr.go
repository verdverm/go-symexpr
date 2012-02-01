package symexpr

import (
  "fmt"
  "math"
  "sort"
)

type CoeffSet struct {
  Vals []float64
}


type Expr interface {


  Size() int
  Depth() int
  Height() int
  NumChildren() int

  ExprType() int
  AmILess( rhs Expr ) bool
  AmISame( rhs Expr ) bool
  AmIAlmostSame( rhs Expr ) bool
  HasVar() bool

  Clone() Expr

  // DFS retrieval
  GetExpr( pos *int ) Expr
  setExpr( pos *int , e Expr ) (bool,bool)

  String() string
// 	WriteString( buf *bytes.Buffer )

  Eval(t float64, x, c, s []float64) float64
  // RK4
  // pRK4

  CalcExprStats( currDepth int ) (mySize int)
  Simplify( rules SimpRules ) Expr

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
type ExprArraySort struct {
  ExprReportArray
}
func (p ExprArraySort) Less(i, j int) bool {
  if p.ExprReportArray[i] == nil && p.ExprReportArray[j] == nil { return false }
  if p.ExprReportArray[i] == nil { return false }
  if p.ExprReportArray[j] == nil { return true 	}
  return p.ExprReportArray[i].AmILess( p.ExprReportArray[j] )
}

const (
  MaxAddChildren = 8
  MaxMulChildren = 8
  MaxAddGenChildren = 4
  MaxMulGenChildren = 4
)

const (
  NULL = iota
  TIME
  CONSTANT
  CONSTANTF
  SYSTEM
  VAR

  NEG
  ABS
  SQRT
  SIN
  COS
  EXP
  LOG
  POWI
  POWF

  DIV
  ADD
  SUB
  MUL
  POWE

  EXPR_MAX
)

type ExprStats struct {
  depth int // layers from top (root == 1)
  height int // layers of subtree (leaf == 1)
  size int
  numchld int
  // for tree based output (only with depth < D)?
  // output []float64
}

func max( l,r int ) int {
  if l > r { return l }
  return r
}

func (es *ExprStats) Depth() int { return es.depth }
func (es *ExprStats) Height()  int { return es.height }
func (es *ExprStats) Size()  int { return es.size }
func (es *ExprStats) NumChildren()  int { return es.numchld }

// Null Leaf  (shouldn't appear ever)
type Null struct{
  ExprStats
}

// Leaf Nodes
type Time struct {
  ExprStats
}

type Var struct {
  ExprStats
  P int
}
type Constant struct {
  ExprStats
  P int
}
type ConstantF struct {
  ExprStats
  F float64
}
type System struct {
  ExprStats
  P int
}

// Unary Operators
type Neg struct {
  ExprStats
  C  Expr
}
type Abs struct {
  ExprStats
  C Expr
}
type Sqrt struct {
  ExprStats
  C Expr
}
type Sin struct {
  ExprStats
  C Expr
}
type Cos struct {
  ExprStats
  C Expr
}
type Exp struct {
  ExprStats
  C Expr
}
type Log struct {
  ExprStats
  C Expr
}

// Half-ary Operators
type PowI struct {
  ExprStats
  Base Expr
  Power int
}
type PowF struct {
  ExprStats
  Base Expr
  Power float64
}



// Binary Operators
type Add struct {
  ExprStats
  CS [MaxAddChildren]Expr
}
func (n *Add) Len() int { return len(n.CS) }
func (n *Add) Swap(i, j int) { n.CS[i], n.CS[j] = n.CS[j], n.CS[i] }
func (n *Add) Less(i, j int) bool {
  if n.CS[i] == nil && n.CS[j] == nil { return false }
  if n.CS[i] == nil { return false }
  if n.CS[j] == nil { return true }
  return n.CS[i].AmILess( n.CS[j] )
}

type Sub struct {
  ExprStats
  C1 Expr
  C2 Expr
}
type Mul struct {
  ExprStats
  CS [MaxMulChildren]Expr
}
func (n *Mul) Len() int { return len(n.CS) }
func (n *Mul) Swap(i, j int) { n.CS[i], n.CS[j] = n.CS[j], n.CS[i] }
func (n *Mul) Less(i, j int) bool {
  if n.CS[i] == nil && n.CS[j] == nil { return false }
  if n.CS[i] == nil { return false }
  if n.CS[j] == nil { return true }
  return n.CS[i].AmILess( n.CS[j] )
}

type Div struct {
  ExprStats
  Numer Expr
  Denom Expr
}
type PowE struct {
  ExprStats
  Base Expr
  Power Expr
}

func (n *Null) ExprType() int    { return NULL }
func (t *Time) ExprType() int    { return TIME }
func (v *Var) ExprType() int     { return VAR }
func (c *Constant) ExprType() int{ return CONSTANT }
func (c *ConstantF) ExprType() int{ return CONSTANTF }
func (s *System) ExprType() int  { return SYSTEM }

func (u *Neg) ExprType() int 	   { return NEG }
func (u *Abs) ExprType() int 	   { return ABS }
func (u *Sqrt) ExprType() int 	{ return SQRT }
func (u *Sin) ExprType() int 	   { return SIN }
func (u *Cos) ExprType() int 	   { return COS }
func (u *Exp) ExprType() int 	   { return EXP }
func (u *Log) ExprType() int 	   { return LOG }
func (u *PowI) ExprType() int 	{ return POWI }
func (u *PowF) ExprType() int 	{ return POWF }

func (n *Add) ExprType() int     { return ADD }
func (n *Sub) ExprType() int     { return SUB }
func (n *Mul) ExprType() int     { return MUL }
func (n *Div) ExprType() int     { return DIV }
func (n *PowE) ExprType() int    { return POWE }




func (n *Null) AmILess( r Expr ) bool     { return NULL < r.ExprType() }
func (n *Time) AmILess( r Expr ) bool     { return TIME < r.ExprType() }
func (v *Var) AmILess( r Expr ) bool      {
  if VAR < r.ExprType() { return true }
  if VAR > r.ExprType() { return false }
  return v.P < r.(*Var).P
}
func (c *Constant) AmILess( r Expr ) bool {
  if CONSTANT < r.ExprType() { return true }
  if CONSTANT > r.ExprType() { return false }
  return c.P < r.(*Constant).P
}
func (c *ConstantF) AmILess( r Expr ) bool {
  if CONSTANTF < r.ExprType() { return true }
  if CONSTANTF > r.ExprType() { return false }
  return c.F < r.(*ConstantF).F
}
func (s *System) AmILess( r Expr ) bool   {
  if SYSTEM < r.ExprType() { return true }
  if SYSTEM > r.ExprType() { return false }
  return s.P < r.(*System).P
}

func (u *Neg) AmILess( r Expr ) bool 	  {
  if NEG < r.ExprType() { return true }
  if NEG > r.ExprType() { return false }
  return u.C.AmILess( r.(*Neg).C )
}
func (u *Abs) AmILess( r Expr ) bool 	  {
  if ABS < r.ExprType() { return true }
  if ABS > r.ExprType() { return false }
  return u.C.AmILess( r.(*Abs).C )
}
func (u *Sqrt) AmILess( r Expr ) bool 	  {
  if SQRT < r.ExprType() { return true }
  if SQRT > r.ExprType() { return false }
  return u.C.AmILess( r.(*Sqrt).C )
}
func (u *Sin) AmILess( r Expr ) bool 	  {
  if SIN < r.ExprType() { return true }
  if SIN > r.ExprType() { return false }
  return u.C.AmILess( r.(*Sin).C )
}
func (u *Cos) AmILess( r Expr ) bool 	  {
  if COS < r.ExprType() { return true }
  if COS > r.ExprType() { return false }
  return u.C.AmILess( r.(*Cos).C )
}
func (u *Exp) AmILess( r Expr ) bool 	  {
  if EXP < r.ExprType() { return true }
  if EXP > r.ExprType() { return false }
  return u.C.AmILess( r.(*Exp).C )
}
func (u *Log) AmILess( r Expr ) bool 	  {
  if LOG < r.ExprType() { return true }
  if LOG > r.ExprType() { return false }
  return u.C.AmILess( r.(*Log).C )
}
func (u *PowI) AmILess( r Expr ) bool     {
  if POWI < r.ExprType() { return true }
  if POWI > r.ExprType() { return false }
  if u.Base.AmILess( r.(*PowI).Base ) { return true }
  if r.(*PowI).Base.AmILess( u.Base ) { return false }
  return u.Power < r.(*PowI).Power
}
func (u *PowF) AmILess( r Expr ) bool 	  {
  if POWF < r.ExprType() { return true }
  if POWF > r.ExprType() { return false }
  if u.Base.AmILess( r.(*PowF).Base ) { return true }
  if r.(*PowF).Base.AmILess( u.Base ) { return false }
  return u.Power < r.(*PowF).Power
}

func (n *Add) AmILess( r Expr ) bool      {
  if ADD < r.ExprType() { return true }
  if ADD > r.ExprType() { return false }
  m := r.(*Add)
  for i,C := range n.CS {
    if C == nil { return false }
    if m.CS[i] == nil { return true }
    if C.AmILess( m.CS[i] ) { return true }
    if m.CS[i].AmILess( C ) { return false }
  }
  return false
}
func (n *Sub) AmILess( r Expr ) bool      {
  if SUB < r.ExprType() { return true }
  if SUB > r.ExprType() { return false }
  if n.C1.AmILess( r.(*Sub).C1 ) { return true }
  if r.(*Sub).C1.AmILess( n.C1 ) { return false }
  return n.C2.AmILess( r.(*Sub).C2 )
}
func (n *Mul) AmILess( r Expr ) bool      {
  if MUL < r.ExprType() { return true }
  if MUL > r.ExprType() { return false }
  m := r.(*Mul)
  for i,C := range n.CS {
    if C == nil { return false }
    if m.CS[i] == nil { return true }
    if C.AmILess( m.CS[i] ) { return true }
    if m.CS[i].AmILess( C ) { return false }
  }
  return false
}
func (n *Div) AmILess( r Expr ) bool      {
  if DIV < r.ExprType() { return true }
  if DIV > r.ExprType() { return false }
  if n.Numer.AmILess( r.(*Div).Numer ) { return true }
  if r.(*Div).Numer.AmILess( n.Numer ) { return false }
  return n.Denom.AmILess( r.(*Div).Denom )
}
func (n *PowE) AmILess( r Expr ) bool     {
  if POWE < r.ExprType() { return true }
  if POWE > r.ExprType() { return false }
  if n.Base.AmILess( r.(*PowE).Base ) { return true }
  if r.(*PowE).Base.AmILess( n.Base ) { return false }
  return n.Power.AmILess( r.(*PowE).Power )
}




func (n *Null) AmISame( r Expr ) bool     { return r.ExprType() == NULL }
func (n *Time) AmISame( r Expr ) bool     { return r.ExprType() == TIME }
func (v *Var) AmISame( r Expr ) bool      { return r.ExprType() == VAR && r.(*Var).P == v.P }
func (c *Constant) AmISame( r Expr ) bool { return r.ExprType() == CONSTANT && r.(*Constant).P == c.P }
func (c *ConstantF) AmISame( r Expr ) bool { return r.ExprType() == CONSTANTF && r.(*ConstantF).F == c.F }
func (s *System) AmISame( r Expr ) bool   { return r.ExprType() == SYSTEM && r.(*System).P == s.P }

func (u *Neg) AmISame( r Expr ) bool 	  { return r.ExprType() == NEG && u.C.AmISame(r.(*Neg).C) }
func (u *Abs) AmISame( r Expr ) bool 	  { return r.ExprType() == ABS && u.C.AmISame(r.(*Abs).C) }
func (u *Sqrt) AmISame( r Expr ) bool 	  { return r.ExprType() == SQRT && u.C.AmISame(r.(*Sqrt).C) }
func (u *Sin) AmISame( r Expr ) bool 	  { return r.ExprType() == SIN && u.C.AmISame(r.(*Sin).C) }
func (u *Cos) AmISame( r Expr ) bool  	  { return r.ExprType() == COS && u.C.AmISame(r.(*Cos).C) }
func (u *Exp) AmISame( r Expr ) bool 	  { return r.ExprType() == EXP && u.C.AmISame(r.(*Exp).C) }
func (u *Log) AmISame( r Expr ) bool 	  { return r.ExprType() == LOG && u.C.AmISame(r.(*Log).C) }
func (u *PowI) AmISame( r Expr ) bool 	  { return r.ExprType() == POWI && r.(*PowI).Power == u.Power && u.Base.AmISame(r.(*PowI).Base) }
func (u *PowF) AmISame( r Expr ) bool 	  { return r.ExprType() == POWF && r.(*PowF).Power == u.Power && u.Base.AmISame(r.(*PowF).Base) }

func (n *Add) AmISame( r Expr ) bool      {
  if r.ExprType() != ADD { return false }
  m := r.(*Add)
  for i,C := range n.CS {
    if C == nil && m.CS[i] == nil { continue }
    if C == nil && m.CS[i] != nil { return false }
    if C != nil && m.CS[i] == nil { return false }
    if !C.AmISame( m.CS[i] ) { return false }
    if m.CS[i].AmILess( C ) { return false }
  }
  return true
}

func (n *Sub) AmISame( r Expr ) bool      { return r.ExprType() == SUB  && n.C1.AmISame(r.(*Sub).C1) && n.C2.AmISame(r.(*Sub).C2)  }

func (n *Mul) AmISame( r Expr ) bool      {
  if r.ExprType() != MUL { return false }
  m := r.(*Mul)
  for i,C := range n.CS {
    if C == nil && m.CS[i] == nil { continue }
    if C == nil && m.CS[i] != nil { return false }
    if C != nil && m.CS[i] == nil { return false }
    if !C.AmISame( m.CS[i] ) { return false }
    if m.CS[i].AmILess( C ) { return false }
  }
  return true
}

func (n *Div) AmISame( r Expr ) bool      { return r.ExprType() == DIV  && n.Numer.AmISame(r.(*Div).Numer) && n.Denom.AmISame(r.(*Div).Denom)  }
func (n *PowE) AmISame( r Expr ) bool     { return r.ExprType() == POWE && n.Base.AmISame(r.(*PowE).Base) && n.Power.AmISame(r.(*PowE).Power) }




func (n *Null) AmIAlmostSame( r Expr ) bool     { return r.ExprType() == NULL }
func (n *Time) AmIAlmostSame( r Expr ) bool     { return r.ExprType() == TIME }
func (v *Var) AmIAlmostSame( r Expr ) bool      { return r.ExprType() == VAR && r.(*Var).P == v.P }
func (c *Constant) AmIAlmostSame( r Expr ) bool { return r.ExprType() == CONSTANT && r.(*Constant).P == c.P }
func (c *ConstantF) AmIAlmostSame( r Expr ) bool { return r.ExprType() == CONSTANTF }
func (s *System) AmIAlmostSame( r Expr ) bool   { return r.ExprType() == SYSTEM && r.(*System).P == s.P }

func (u *Neg) AmIAlmostSame( r Expr ) bool 	  { return r.ExprType() == NEG && u.C.AmIAlmostSame(r.(*Neg).C) }
func (u *Abs) AmIAlmostSame( r Expr ) bool 	  { return r.ExprType() == ABS && u.C.AmIAlmostSame(r.(*Abs).C) }
func (u *Sqrt) AmIAlmostSame( r Expr ) bool 	  { return r.ExprType() == SQRT && u.C.AmIAlmostSame(r.(*Sqrt).C) }
func (u *Sin) AmIAlmostSame( r Expr ) bool 	  { return r.ExprType() == SIN && u.C.AmIAlmostSame(r.(*Sin).C) }
func (u *Cos) AmIAlmostSame( r Expr ) bool  	  { return r.ExprType() == COS && u.C.AmIAlmostSame(r.(*Cos).C) }
func (u *Exp) AmIAlmostSame( r Expr ) bool 	  { return r.ExprType() == EXP && u.C.AmIAlmostSame(r.(*Exp).C) }
func (u *Log) AmIAlmostSame( r Expr ) bool 	  { return r.ExprType() == LOG && u.C.AmIAlmostSame(r.(*Log).C) }
func (u *PowI) AmIAlmostSame( r Expr ) bool 	  { return r.ExprType() == POWI && r.(*PowI).Power == u.Power && u.Base.AmIAlmostSame(r.(*PowI).Base) }
func (u *PowF) AmIAlmostSame( r Expr ) bool 	  { return r.ExprType() == POWF && r.(*PowF).Power == u.Power && u.Base.AmIAlmostSame(r.(*PowF).Base) }

func (n *Add) AmIAlmostSame( r Expr ) bool      {
  if r.ExprType() != ADD { return false }
  m := r.(*Add)
  for i,C := range n.CS {
    if C == nil && m.CS[i] == nil { continue }
    if C == nil && m.CS[i] != nil { return false }
    if C != nil && m.CS[i] == nil { return false }
    if !C.AmIAlmostSame( m.CS[i] ) && m.CS[i].AmILess( C ) { return false }
  }
  return true
}

func (n *Sub) AmIAlmostSame( r Expr ) bool      { return r.ExprType() == SUB  && n.C1.AmIAlmostSame(r.(*Sub).C1) && n.C2.AmIAlmostSame(r.(*Sub).C2)  }

func (n *Mul) AmIAlmostSame( r Expr ) bool      {
  if r.ExprType() != MUL { return false }
  m := r.(*Mul)
  for i,C := range n.CS {
    if C == nil && m.CS[i] == nil { continue }
    if C == nil && m.CS[i] != nil { return false }
    if C != nil && m.CS[i] == nil { return false }
    if !C.AmIAlmostSame( m.CS[i] ) && m.CS[i].AmILess( C ) { return false }
  }
  return true
}

func (n *Div) AmIAlmostSame( r Expr ) bool      { return r.ExprType() == DIV  && n.Numer.AmIAlmostSame(r.(*Div).Numer) && n.Denom.AmIAlmostSame(r.(*Div).Denom)  }
func (n *PowE) AmIAlmostSame( r Expr ) bool     { return r.ExprType() == POWE && n.Base.AmIAlmostSame(r.(*PowE).Base) && n.Power.AmIAlmostSame(r.(*PowE).Power) }




func (n *Null) HasVar() bool     { return false }
func (n *Time) HasVar() bool     { return false }
func (v *Var) HasVar() bool      { return true }
func (c *Constant) HasVar() bool { return false }
func (c *ConstantF) HasVar() bool { return false }
func (s *System) HasVar() bool   { return false }

func (u *Neg) HasVar() bool 	  { return u.C.HasVar() }
func (u *Abs) HasVar() bool 	  { return u.C.HasVar() }
func (u *Sqrt) HasVar() bool 	  { return u.C.HasVar() }
func (u *Sin) HasVar() bool 	  { return u.C.HasVar() }
func (u *Cos) HasVar() bool  	  { return u.C.HasVar() }
func (u *Exp) HasVar() bool 	  { return u.C.HasVar() }
func (u *Log) HasVar() bool 	  { return u.C.HasVar() }
func (u *PowI) HasVar() bool 	  { return u.Base.HasVar() }
func (u *PowF) HasVar() bool 	  { return u.Base.HasVar() }

func (n *Add) HasVar() bool      {
  for _,C := range n.CS {
    if C != nil && C.HasVar() { return true }
  }
  return false
}

func (n *Sub) HasVar() bool      { return n.C1.HasVar() || n.C2.HasVar()  }

func (n *Mul) HasVar() bool      {
  for _,C := range n.CS {
    if C != nil &&  C.HasVar() { return true }
  }
  return false
}

func (n *Div) HasVar() bool      { return n.Numer.HasVar() || n.Denom.HasVar()  }
func (n *PowE) HasVar() bool     { return n.Base.HasVar()  || n.Power.HasVar() }





func (n *Null) Clone() Expr     { return &Null{ExprStats{0,0,0,0}} }
func (n *Time) Clone() Expr     { return &Time{ExprStats{0,0,0,0}} }
func (v *Var) Clone() Expr      { return &Var{P: v.P} }
func (c *Constant) Clone() Expr { return &Constant{P: c.P} }
func (c *ConstantF) Clone() Expr { return &ConstantF{F: c.F} }
func (s *System) Clone() Expr   { return &System{P: s.P} }

func (u *Neg) Clone() Expr 	  { return &Neg{C: u.C.Clone()} }
func (u *Abs) Clone() Expr 	  { return &Abs{C: u.C.Clone()} }
func (u *Sqrt) Clone() Expr 	  { return &Sqrt{C: u.C.Clone()} }
func (u *Sin) Clone() Expr 	  { return &Sin{C: u.C.Clone()} }
func (u *Cos) Clone() Expr 	  { return &Cos{C: u.C.Clone()} }
func (u *Exp) Clone() Expr 	  { return &Exp{C: u.C.Clone()} }
func (u *Log) Clone() Expr 	  { return &Log{C: u.C.Clone()} }
func (u *PowI) Clone() Expr 	  { return &PowI{Base: u.Base.Clone(), Power: u.Power} }
func (u *PowF) Clone() Expr 	  { return &PowF{Base: u.Base.Clone(), Power: u.Power} }

func (n *Add) Clone() Expr      {
  var a Add
  for i,C := range n.CS {
    if C != nil {
      a.CS[i] = C.Clone()
    }
  }
  return &a
}
func (n *Sub) Clone() Expr      { return &Sub{C1: n.C1.Clone(), C2: n.C2.Clone()} }
func (n *Mul) Clone() Expr      {
  var m Mul
  for i,C := range n.CS {
    if C != nil {
      m.CS[i] = C.Clone()
    }
  }
  return &m
}
func (n *Div) Clone() Expr      { return &Div{Numer: n.Numer.Clone(), Denom: n.Denom.Clone()} }
func (n *PowE) Clone() Expr     { return &PowE{Base: n.Base.Clone(), Power: n.Power.Clone()} }



func (n *Null) GetExpr( pos *int ) Expr     { if (*pos) == 0 { return n }; (*pos)--; return nil }
func (n *Time) GetExpr( pos *int ) Expr     { if (*pos) == 0 { return n }; (*pos)--; return nil }
func (v *Var) GetExpr( pos *int ) Expr      { if (*pos) == 0 { return v }; (*pos)--; return nil }
func (c *Constant) GetExpr( pos *int ) Expr { if (*pos) == 0 { return c }; (*pos)--; return nil }
func (c *ConstantF) GetExpr( pos *int ) Expr { if (*pos)== 0 { return c }; (*pos)--; return nil }
func (s *System) GetExpr( pos *int ) Expr   { if (*pos) == 0 { return s }; (*pos)--; return nil }

func (u *Neg) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return u }
  (*pos)--
  return u.C.GetExpr(pos)
}
func (u *Abs) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return u }
  (*pos)--
  return u.C.GetExpr(pos)
}
func (u *Sqrt) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return u }
  (*pos)--
  return u.C.GetExpr(pos)
}
func (u *Sin) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return u }
  (*pos)--
  return u.C.GetExpr(pos)
}
func (u *Cos) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return u }
  (*pos)--
  return u.C.GetExpr(pos)
}
func (u *Exp) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return u }
  (*pos)--
  return u.C.GetExpr(pos)
}
func (u *Log) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return u }
  (*pos)--
  return u.C.GetExpr(pos)
}
func (u *PowI) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return u }
  (*pos)--
  return u.Base.GetExpr(pos)
}
func (u *PowF) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return u }
  (*pos)--
  return u.Base.GetExpr(pos)
}

func (n *Add) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return n }
  (*pos)--
  for _,C := range n.CS {
    if C == nil { continue }
    tmp := C.GetExpr(pos)
    if tmp != nil { return tmp }
    if *pos < 0 { return nil }
  }
  return nil
}
func (n *Sub) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return n }
  (*pos)--
  tmp := n.C1.GetExpr(pos)
  if tmp != nil { return tmp }
  return n.C2.GetExpr(pos)
}
func (n *Mul) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return n }
  (*pos)--
  for _,C := range n.CS {
    if C == nil { continue }
    tmp := C.GetExpr(pos)
    if tmp != nil { return tmp }
    if *pos < 0 { return nil }
  }
  return nil
}
func (n *Div) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return n }
  (*pos)--
  tmp := n.Numer.GetExpr(pos)
  if tmp != nil { return tmp }
  return n.Denom.GetExpr(pos)
}
func (n *PowE) GetExpr( pos *int ) Expr {
  if (*pos) == 0 { return n }
  (*pos)--
  tmp := n.Base.GetExpr(pos)
  if tmp != nil { return tmp }
  return n.Power.GetExpr(pos)
}

func SwapExpr( orig, newt Expr, pos int ) (ret Expr) {
//   fmt.Printf( "SWAP orig  %v\n", orig )
  p := pos
//   oldt := orig.GetExpr(&p)
//   fmt.Printf( "SWAP (%d)\n%v\n%v\n", pos, oldt, newt )
  rme,_ := orig.setExpr(&p,newt)
  if rme {
    ret = newt
  }

//   fmt.Printf( "SWAP ret  %v\n", ret )
  return
}

func (n *Null) setExpr( pos *int, e Expr ) ( replace_me, replaced bool )     { if (*pos) == 0 { return true,false }; (*pos)--; return false,false }
func (n *Time) setExpr( pos *int, e Expr ) ( replace_me, replaced bool )     { if (*pos) == 0 { return true,false }; (*pos)--; return false,false }
func (v *Var) setExpr( pos *int, e Expr ) ( replace_me, replaced bool )      { if (*pos) == 0 { return true,false }; (*pos)--; return false,false }
func (c *Constant) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) { if (*pos) == 0 { return true,false }; (*pos)--; return false,false }
func (c *ConstantF) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) { if (*pos)== 0 { return true,false }; (*pos)--; return false,false }
func (s *System) setExpr( pos *int, e Expr ) ( replace_me, replaced bool )   { if (*pos) == 0 { return true,false }; (*pos)--; return false,false }

func (u *Neg) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := u.C.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    u.C = e
    return false,true
  }
  return false,repd
}
func (u *Abs) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := u.C.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    u.C = e
    return false,true
  }
  return false,repd
}
func (u *Sqrt) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := u.C.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    u.C = e
    return false,true
  }
  return false,repd
}
func (u *Sin) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := u.C.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    u.C = e
    return false,true
  }
  return false,repd
}
func (u *Cos) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := u.C.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    u.C = e
    return false,true
  }
  return false,repd
}
func (u *Exp) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := u.C.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    u.C = e
    return false,true
  }
  return false,repd
}
func (u *Log) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := u.C.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    u.C = e
    return false,true
  }
  return false,repd
}
func (u *PowI) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := u.Base.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    u.Base = e
    return false,true
  }
  return false,repd
}
func (u *PowF) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := u.Base.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    u.Base = e
    return false,true
  }
  return false,repd
}

func (n *Add) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  var rme,repd bool
  for i,C := range n.CS {
    if C == nil { continue }
    rme,repd = C.setExpr(pos,e)
    if repd { return false,true }
    if rme {
      n.CS[i] = e
      return false,true
    }
//     if *pos < 0 { return nil }
  }
  return false,repd
}
func (n *Sub) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := n.C1.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    n.C1 = e
    return false,true
  }
//   if *pos < 0 { return nil }
  rme,repd = n.C2.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    n.C2 = e
    return false,true
  }
  return false,repd
}
func (n *Mul) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  var rme,repd bool
  for i,C := range n.CS {
    if C == nil { continue }
    rme,repd = C.setExpr(pos,e)
    if repd { return false,true }
    if rme {
      n.CS[i] = e
      return false,true
    }
//     if *pos < 0 { return nil }
  }
  return false,repd
}
func (n *Div) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := n.Numer.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    n.Numer = e
    return false,true
  }
//   if *pos < 0 { return nil }
  rme,repd = n.Denom.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    n.Denom = e
    return false,true
  }
  return false,repd
}
func (n *PowE) setExpr( pos *int, e Expr ) ( replace_me, replaced bool ) {
  if (*pos) == 0 { return true,false }
  (*pos)--
  rme,repd := n.Base.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    n.Base = e
    return false,true
  }
  rme,repd = n.Power.setExpr(pos,e)
  if repd { return false,true }
  if rme {
    n.Power = e
    return false,true
  }
  return false,repd
}


func (n *Null) String() string     { return "NULL" }
func (n *Time) String() string     { return "T" }
func (v *Var) String() string      { return "X_" + fmt.Sprint(v.P) }
func (c *Constant) String() string { return "C_" + fmt.Sprint(c.P) }
func (c *ConstantF) String() string { return fmt.Sprintf("%.4f",c.F) }
func (s *System) String() string   { return "S_" + fmt.Sprint(s.P) }

func (u *Neg) String() string 	  { return "-(" + u.C.String() + ")" }
func (u *Abs) String() string 	  { return "abs(" + u.C.String() + ")" }
func (u *Sqrt) String() string 	  { return "sqrt(" + u.C.String() + ")" }
func (u *Sin) String() string 	  { return "sin(" + u.C.String() + ")" }
func (u *Cos) String() string 	  { return "cos(" + u.C.String() + ")" }
func (u *Exp) String() string 	  { return "e^(" + u.C.String() + ")" }
func (u *Log) String() string 	  { return "ln(" + u.C.String() + ")" }
func (u *PowI) String() string 	  { return "(" + u.Base.String() + ")^" + fmt.Sprint(u.Power) }
func (u *PowF) String() string 	  { return "(" + u.Base.String() + ")^" + fmt.Sprint(u.Power) }

func (n *Add) String() string      {
  str := "( " + n.CS[0].String()
  for i := 1; i < len(n.CS); i++ {
    if n.CS[i] == nil { continue }
    str += " + " + n.CS[i].String()
  }
  str += " )"
  return str
}

func (n *Sub) String() string      { return n.C1.String() + " - " + n.C2.String() }
func (n *Mul) String() string      {
//   str := "|" + n.CS[0].String()
//   for i := 1; i < len(n.CS); i++ {
//     if n.CS[i] == nil { continue }
//     str += "*" + n.CS[i].String()
//   }
//   str += "|"
//   return str
  str := n.CS[0].String()
  for i := 1; i < len(n.CS); i++ {
    if n.CS[i] == nil { continue }
    str += "*" + n.CS[i].String()
  }
  return str
}
func (n *Div) String() string      { return "{ " + n.Numer.String() + " }//{ " + n.Denom.String() + " }" }
func (n *PowE) String() string      { return "{" + n.Base.String() + "}^(" + n.Power.String() + ")" }





func (n *Null) Eval(t float64, x, c, s []float64) float64        { return math.NaN() }
func (*Time) Eval(t float64, x, c, s []float64) float64          { return t }
func (v *Var) Eval(t float64, x, c, s []float64) float64         { return x[v.P] }
func (cnst *Constant) Eval(t float64, x, c, s []float64) float64 { return c[cnst.P] }
func (cnst *ConstantF) Eval(t float64, x, c, s []float64) float64 { return cnst.F }
func (sys *System) Eval(t float64, x, c, s []float64) float64    { return s[sys.P] }

func (u *Neg) Eval(t float64, x, c, s []float64) float64 {
  return -1. * u.C.Eval(t, x, c, s)
}
func (u *Abs) Eval(t float64, x, c, s []float64) float64 {
  return math.Abs( u.C.Eval(t, x, c, s) )
}
func (u *Sqrt) Eval(t float64, x, c, s []float64) float64 {
  return math.Sqrt( u.C.Eval(t, x, c, s) )
}
func (u *Sin) Eval(t float64, x, c, s []float64) float64 {
  return math.Sin( u.C.Eval(t, x, c, s) )
}
func (u *Cos) Eval(t float64, x, c, s []float64) float64 {
  return math.Cos( u.C.Eval(t, x, c, s) )
}
func (u *Exp) Eval(t float64, x, c, s []float64) float64 {
  return math.Exp( u.C.Eval(t, x, c, s) )
}
func (u *Log) Eval(t float64, x, c, s []float64) float64 {
  return math.Log( u.C.Eval(t, x, c, s) )
}
func (u *PowI) Eval(t float64, x, c, s []float64) float64 {
  return math.Pow( u.Base.Eval(t, x, c, s), float64(u.Power) )
}
func (u *PowF) Eval(t float64, x, c, s []float64) float64 {
  return math.Pow( u.Base.Eval(t, x, c, s), u.Power )
}

func (n *Add) Eval(t float64, x, c, s []float64) float64 {
  ret := 0.0
  for _,C := range n.CS {
    if C == nil { continue }
    ret += C.Eval(t, x, c, s)
  }
  return ret
}
func (n *Sub) Eval(t float64, x, c, s []float64) float64 {
  return n.C1.Eval(t, x, c, s) - n.C2.Eval(t, x, c, s)
}
func (n *Mul) Eval(t float64, x, c, s []float64) float64 {
  ret := 1.0
  for _,C := range n.CS {
    if C == nil { continue }
    ret *= C.Eval(t, x, c, s)
  }
  return ret}
func (n *Div) Eval(t float64, x, c, s []float64) float64 {
  return n.Numer.Eval(t, x, c, s) / n.Denom.Eval(t, x, c, s)
}
func (n *PowE) Eval(t float64, x, c, s []float64) float64 {
  return math.Pow(n.Base.Eval(t, x, c, s), n.Power.Eval(t, x, c, s))
}


func RK4(eqn []Expr, ti, tj float64, x_in, x_tmp, c, s []float64) (x_out []float64) {
  var k [32][4]float64
  L := len(x_in)
  h := tj - ti
  for i := 0; i < L; i++ { k[i][0] = eqn[i].Eval(ti, x_in, c, s)      }
  for i := 0; i < L; i++ { x_tmp[i] = x_in[i] + (h * k[i][0] / 2.0)  }
  for i := 0; i < L; i++ { k[i][1] = eqn[i].Eval(ti, x_tmp, c, s)    }
  for i := 0; i < L; i++ { x_tmp[i] = x_in[i] + (h * k[i][1] / 2.0)   }
  for i := 0; i < L; i++ { k[i][2] = eqn[i].Eval(ti, x_tmp, c, s)}
  for i := 0; i < L; i++ { x_tmp[i] = x_in[i] + (h * k[i][2])}
  for i := 0; i < L; i++ { k[i][3] = eqn[i].Eval(ti, x_tmp, c, s)}
  for i := 0; i < L; i++ { x_out[i] =  ((k[i][0] + 2.0*k[i][1] + 2.0*k[i][2] + k[i][3]) * (h / 6.0)) }

  return
}
/*
func eval_err() {
  if err := recover(); err != nil {
    fmt.Printf("Error: %v", err)
  }
}*/
func PRK4(xn int, eqn Expr, ti, tj float64, x_in, x_out, x_tmp, c, s []float64) float64 {
  //  defer eval_err()
  var k [4]float64
  L := len(x_in)
  h := tj - ti
  for i := 0; i < L; i++ {
    x_tmp[i] = x_in[i] + (0.5 * (x_out[i] - x_in[i]))
  }
  k[0] = eqn.Eval(ti, x_in, c, s)
  x_tmp[xn] = x_in[xn] + (h * k[0] / 2.0)
  k[1] = eqn.Eval(ti, x_tmp, c, s)
  x_tmp[xn] = x_in[xn] + (h * k[1] / 2.0)
  k[2] = eqn.Eval(ti, x_tmp, c, s)
  x_tmp[xn] = x_in[xn] + (h * k[2])
  k[3] = eqn.Eval(ti, x_tmp, c, s)
  return ((k[0] + 2.0*k[1] + 2.0*k[2] + k[3]) * (h / 6.0))
}



func PrintPRK4(xn int, eqn Expr, ti, to float64, x_in, x_out, x_tmp, c, s []float64) float64 {
  //  defer eval_err()
  var k [4]float64
  L := len(x_in)
  h := to - ti
//   fmt.Printf( "t: %.4f\n", h )
  for i := 0; i < L; i++ {
    x_tmp[i] = x_in[i] + (0.5 * (x_out[i] - x_in[i]))
  }
  fmt.Printf( "in:   %v\n", x_in )
  fmt.Printf( "out:  %v\n", x_out )

  fmt.Printf( "tmp:  %v\n", x_tmp )
  k[0] = eqn.Eval(ti, x_in, c, s)
  x_tmp[xn] = x_in[xn] + (h * k[0] / 2.0)
  fmt.Printf( "tmp:  %v\n", x_tmp )
  k[1] = eqn.Eval(ti, x_tmp, c, s)
  x_tmp[xn] = x_in[xn] + (h * k[1] / 2.0)
  fmt.Printf( "tmp:  %v\n", x_tmp )
  k[2] = eqn.Eval(ti, x_tmp, c, s)
  x_tmp[xn] = x_in[xn] + (h * k[2])
  fmt.Printf( "tmp:  %v\n", x_tmp )
  k[3] = eqn.Eval(ti, x_tmp, c, s)
  fmt.Printf( "k:    %v\n", k )
  ans := ((k[0] + 2.0*k[1] + 2.0*k[2] + k[3]) * (h / 6.0))
  fmt.Printf( "ans:  %.4f   =>   %.4f\n\n", ans, x_out[xn]-x_in[xn] )
  return ans
}





func (n *Null) CalcExprStats( currDepth int ) (mySize int) {
  n.depth = currDepth+1
  n.height = 0
  n.size = 0
  n.numchld = 0
  return n.size
}
func (t *Time) CalcExprStats( currDepth int ) (mySize int) {
  t.depth = currDepth+1
  t.height = 1
  t.size = 1
  t.numchld = 0
  return t.size
}
func (v *Var) CalcExprStats( currDepth int ) (mySize int) {
  v.depth = currDepth+1
  v.height = 1
  v.size = 1
  v.numchld = 0
  return v.size
}
func (c *Constant) CalcExprStats( currDepth int ) (mySize int) {
  c.depth = currDepth+1
  c.height = 1
  c.size = 1
  c.numchld = 0
  return c.size
}
func (c *ConstantF) CalcExprStats( currDepth int ) (mySize int) {
  c.depth = currDepth+1
  c.height = 1
  c.size = 1
  c.numchld = 0
  return c.size
}
func (s *System) CalcExprStats( currDepth int ) (mySize int) {
  s.depth = currDepth+1
  s.height = 1
  s.size = 1
  s.numchld = 0
  return s.size
}

func (u *Neg) CalcExprStats( currDepth int ) (mySize int) {
  u.depth = currDepth+1
  u.size = 1 + u.C.CalcExprStats(currDepth+1)
  u.height = 1 + u.C.Height()
  u.numchld = 1
  return u.size
}
func (u *Abs) CalcExprStats( currDepth int ) (mySize int) {
  u.depth = currDepth+1
  u.size = 1 + u.C.CalcExprStats(currDepth+1)
  u.height = 1 + u.C.Height()
  u.numchld = 1
  return u.size
}
func (u *Sqrt) CalcExprStats( currDepth int ) (mySize int) {
  u.depth = currDepth+1
  u.size = 1 + u.C.CalcExprStats(currDepth+1)
  u.height = 1 + u.C.Height()
  u.numchld = 1
  return u.size
}
func (u *Sin) CalcExprStats( currDepth int ) (mySize int) {
  u.depth = currDepth+1
  u.size = 1 + u.C.CalcExprStats(currDepth+1)
  u.height = 1 + u.C.Height()
  u.numchld = 1
  return u.size
}
func (u *Cos) CalcExprStats( currDepth int ) (mySize int) {
  u.depth = currDepth+1
  u.size = 1 + u.C.CalcExprStats(currDepth+1)
  u.height = 1 + u.C.Height()
  u.numchld = 1
  return u.size
}
func (u *Exp) CalcExprStats( currDepth int ) (mySize int) {
  u.depth = currDepth+1
  u.size = 1 + u.C.CalcExprStats(currDepth+1)
  u.height = 1 + u.C.Height()
  u.numchld = 1
  return u.size
}
func (u *Log) CalcExprStats( currDepth int ) (mySize int) {
  u.depth = currDepth+1
  u.size = 1 + u.C.CalcExprStats(currDepth+1)
  u.height = 1 + u.C.Height()
  u.numchld = 1
  return u.size
}
func (u *PowI) CalcExprStats( currDepth int ) (mySize int) {
  u.depth = currDepth+1
  u.size = 1 + u.Base.CalcExprStats(currDepth+1)
  u.height = 1 + u.Base.Height()
  u.numchld = 1
  return u.size
}
func (u *PowF) CalcExprStats( currDepth int ) (mySize int) {
  u.depth = currDepth+1
  u.size = 1 + u.Base.CalcExprStats(currDepth+1)
  u.height = 1 + u.Base.Height()
  u.numchld = 1
  return u.size
}

func (n *Add) CalcExprStats( currDepth int ) (mySize int) {
  n.depth = currDepth+1
  n.size = 1
  n.numchld = 0
  h := 0
  for _,C := range n.CS {
    if C == nil { continue } else { n.numchld++ }
    n.size +=  C.CalcExprStats(currDepth+1)
    if h < C.Height() { h = C.Height() }
  }
  n.height = 1 + h
  return n.size
}
func (n *Sub) CalcExprStats( currDepth int ) (mySize int) {
  n.depth = currDepth+1
  n.size = 1 + n.C1.CalcExprStats(currDepth+1) + n.C2.CalcExprStats(currDepth+1)
  n.height = 1 + max(n.C1.Height(),n.C2.Height())
  n.numchld = 2
  return n.size
}
func (n *Mul) CalcExprStats( currDepth int ) (mySize int) {
  n.depth = currDepth+1
  n.size = 1
  n.numchld = 0
  h := 0
  for _,C := range n.CS {
    if C == nil { continue } else { n.numchld++ }
    n.size +=  C.CalcExprStats(currDepth+1)
    if h < C.Height() { h = C.Height() }
  }
  n.height = 1 + h
  return n.size

}
func (n *Div) CalcExprStats( currDepth int ) (mySize int) {
  n.depth = currDepth+1
  n.size = 1 + n.Numer.CalcExprStats(currDepth+1) + n.Denom.CalcExprStats(currDepth+1)
  n.height = 1 + max(n.Numer.Height(),n.Denom.Height())
  n.numchld = 2
  return n.size
}
func (n *PowE) CalcExprStats( currDepth int ) (mySize int) {
  n.depth = currDepth+1
  n.size = 1 + n.Base.CalcExprStats(currDepth+1) + n.Power.CalcExprStats(currDepth+1)
  n.height = 1 + max(n.Base.Height(),n.Power.Height())
  n.numchld = 2
  return n.size
}



type SimpRules struct {
    EvalUnary bool
    GroupCoeff bool
}


func (n *Null) Simplify( rules SimpRules ) Expr      { return n }
func (n *Time) Simplify( rules SimpRules ) Expr      { return n }
func (v *Var) Simplify( rules SimpRules ) Expr       { return v }
func (c *Constant) Simplify( rules SimpRules ) Expr  { return c }
func (c *ConstantF) Simplify( rules SimpRules ) Expr {
	if math.IsNaN(c.F) || math.IsInf(c.F,0) {
		return &Null{}
	}
	// "close" to system value

	return c
}
func (s *System) Simplify( rules SimpRules ) Expr    { return s }

func (u *Neg) Simplify( rules SimpRules ) Expr {
	var (
	  ret Expr = u
	  t int = NULL
	)
	if u.C != nil {
		u.C = u.C.Simplify(rules)
		t = u.C.ExprType()
	}
	switch t {
		case NULL, NEG:
			ret = u.C
			u.C = nil
		case CONSTANTF:
			ret = u.C
			u.C = nil
			ret.(*ConstantF).F *= -1.0
    case MUL:
      m := u.C.(*Mul)
      if m.CS[0].ExprType() == CONSTANTF {
        m.CS[0].(*ConstantF).F *= -1.0
        ret = u.C
        u.C = nil
      }
	}
	return ret
}
func (u *Abs) Simplify( rules SimpRules ) Expr {
	var (
	  ret Expr = u
	  t int = NULL
	)
	if u.C != nil {
		u.C = u.C.Simplify(rules)
		t = u.C.ExprType()
	}
	switch t {
		case NULL, ABS:
			ret = u.C
			u.C = nil
		case CONSTANTF:
			ret = u.C
			u.C = nil
			ret.(*ConstantF).F = math.Abs( ret.(*ConstantF).F )
	}
	return ret
}
func (u *Sqrt) Simplify( rules SimpRules ) Expr {
	var (
	  ret Expr = u
	  t int = NULL
	)
	if u.C != nil {
		u.C = u.C.Simplify(rules)
		t = u.C.ExprType()
	}
	switch t {
		case NULL:
			ret = u.C
			u.C = nil
		case CONSTANTF:
			ret = u.C
			u.C = nil
			ret.(*ConstantF).F = math.Sqrt( ret.(*ConstantF).F )
	}
	return ret
}
func (u *Sin) Simplify( rules SimpRules ) Expr 	  {
	var (
	  ret Expr = u
	  t int = NULL
	)
	if u.C != nil {
		u.C = u.C.Simplify(rules)
		t = u.C.ExprType()
	}
	switch t {
		case NULL,SIN,COS:
			ret = u.C
			u.C = nil
		case CONSTANTF:
			ret = u.C
			u.C = nil
			ret.(*ConstantF).F = math.Sin( ret.(*ConstantF).F )
	}
	return ret
}
func (u *Cos) Simplify( rules SimpRules ) Expr 	  {
	var (
	  ret Expr = u
	  t int = NULL
	)
	if u.C != nil {
		u.C = u.C.Simplify(rules)
		t = u.C.ExprType()
	}
	switch t {
		case NULL,SIN,COS:
			ret = u.C
			u.C = nil
		case CONSTANTF:
			ret = u.C
			u.C = nil
			ret.(*ConstantF).F = math.Cos( ret.(*ConstantF).F )
	}
	return ret
}
func (u *Exp) Simplify( rules SimpRules ) Expr 	  {
	var (
	  ret Expr = u
	  t int = NULL
	)
	if u.C != nil {
		u.C = u.C.Simplify(rules)
		t = u.C.ExprType()
	}
	switch t {
		case NULL:
			ret = u.C
			u.C = nil
		case CONSTANTF:
			ret = u.C
			u.C = nil
			ret.(*ConstantF).F = math.Exp( ret.(*ConstantF).F )
	}
	return ret
}
func (u *Log) Simplify( rules SimpRules ) Expr 	  {
	var (
	  ret Expr = u
	  t int = NULL
	)
	if u.C != nil {
		u.C = u.C.Simplify(rules)
		t = u.C.ExprType()
	}
	switch t {
		case NULL:
			ret = u.C
			u.C = nil
		case CONSTANTF:
			ret = u.C
			u.C = nil
			ret.(*ConstantF).F = math.Log( ret.(*ConstantF).F )
	}
	return ret
}
func (u *PowI) Simplify( rules SimpRules ) Expr 	  {
	var (
	  ret Expr = u
	  t int = NULL
	)
	if u.Base != nil {
		u.Base = u.Base.Simplify(rules)
		t = u.Base.ExprType()
	}
	if u.Power == 0 {
	  ret = &ConstantF{F:1}
	} else if u.Power == 1 {
		ret = u.Base
		u.Base = nil
	} else {
		switch t {
			case NULL:
				ret = u.Base
				u.Base = nil
			case CONSTANTF:
				ret = u.Base
				u.Base = nil
				ret.(*ConstantF).F = math.Pow( ret.(*ConstantF).F, float64(u.Power) )
		}
	}
	return ret
}
func (u *PowF) Simplify( rules SimpRules ) Expr 	  {
	var (
	  ret Expr = u
	  t int = NULL
	)
	if u.Base != nil {
		u.Base = u.Base.Simplify(rules)
		t = u.Base.ExprType()
	}
	switch t {
		case NULL:
			ret = u.Base
			u.Base = nil
		case CONSTANTF:
			ret = u.Base
			u.Base = nil
			ret.(*ConstantF).F = math.Pow( ret.(*ConstantF).F, float64(u.Power) )
	}
	return ret
}
func (n *Sub) Simplify( rules SimpRules ) Expr      {
	var (
	  ret Expr = n
	  t1, t2 int = NULL,NULL
	)
	if n.C1 != nil {
		n.C1 = n.C1.Simplify(rules)
		t1 = n.C1.ExprType()
	}
	if n.C2 != nil {
		n.C2 = n.C2.Simplify(rules)
		t2 = n.C2.ExprType()
	}

	if t1 == NULL && t2 == NULL {
		return &Null{}
	} else if t1 == NULL {
	   ret = n.C2
	   n.C1 = nil
	   n.C2 = nil
	} else if t2 == NULL {
	   ret = n.C1
	   n.C1 = nil
	   n.C2 = nil
	} else if n.C1.ExprType() == n.C2.ExprType() &&
	   n.C1.ExprType() == CONSTANTF {
	   ret = n.C1
	   ret.(*ConstantF).F -= n.C2.(*ConstantF).F
	   n.C1 = nil
	   n.C2 = nil
	}
	return ret
}

func (n *Div) Simplify( rules SimpRules ) Expr      {
	var (
	  ret Expr = n
	  t1, t2 int = NULL,NULL
	)
	if n.Numer != nil {
		n.Numer = n.Numer.Simplify(rules)
		t1 = n.Numer.ExprType()
	}
	if n.Denom != nil {
		n.Denom = n.Denom.Simplify(rules)
		t2 = n.Denom.ExprType()
	}

	if t1 == NULL && t2 == NULL {
		return &Null{}
	} else if t1 == NULL {
	   ret = n.Denom
	   n.Numer = nil
	   n.Denom = nil
	} else if t2 == NULL {
	   ret = n.Numer
	   n.Numer = nil
	   n.Denom = nil
	} else if n.Numer.ExprType() == n.Denom.ExprType() &&
	   n.Numer.ExprType() == CONSTANTF {
	   ret = n.Numer
	   ret.(*ConstantF).F /= n.Denom.(*ConstantF).F
	   n.Numer = nil
	   n.Denom = nil
	}
	return ret
}
func (n *PowE) Simplify( rules SimpRules ) Expr     {
	var (
	  ret Expr = n
	  t1, t2 int = NULL,NULL
	)
	if n.Base != nil {
		n.Base = n.Base.Simplify(rules)
		t1 = n.Base.ExprType()
	}
	if n.Power != nil {
		n.Power = n.Power.Simplify(rules)
		t2 = n.Power.ExprType()
	}
	if t1 == NULL && t2 == NULL {
		return &Null{}
	} else if t1 == NULL {
	   ret = n.Power
	   n.Base = nil
	   n.Power = nil
	} else if t2 == NULL {
	   ret = n.Base
	   n.Base = nil
	   n.Power = nil
	} else if n.Base.ExprType() == n.Power.ExprType() &&
	   n.Base.ExprType() == CONSTANTF {
	   ret = n.Base
	   ret.(*ConstantF).F = math.Pow( ret.(*ConstantF).F, n.Power.(*ConstantF).F )
	   n.Base = nil
	   n.Power = nil
	}
	return ret
}



// Simplify an addition node
func (n *Add) Simplify( rules SimpRules ) Expr      {
  var ret Expr = n
  for i,C := range n.CS {
    if C != nil {
      n.CS[i] = C.Simplify(rules)
      if n.CS[i].ExprType() == NULL { n.CS[i] = nil }
    }
  }
  groupCoeff( n.CS[:], plus )
  groupAddTerms(n.CS[:])
  sort.Sort( n )

  gatherAdds(n)
  sort.Sort( n )

  groupCoeff( n.CS[:], plus )
  groupAddTerms(n.CS[:])

  sort.Sort( n )
  cnt := countTerms(n.CS[:])
  if cnt == 0 { ret = &Null{} }
  if cnt == 1 { ret = n.CS[0]; n.CS[0] = nil }
  return ret
}




// Simplify a multiplication node
func (n *Mul) Simplify( rules SimpRules ) Expr      {
  var ret Expr = n
  for i,C := range n.CS {
    if C != nil {
      n.CS[i] = C.Simplify(rules)
      if n.CS[i].ExprType() == NULL { n.CS[i] = nil }
    }
  }
  groupCoeff( n.CS[:], plus )
  groupMulTerms(n.CS[:])
  sort.Sort( n )

  gatherMuls(n)
  sort.Sort( n )

  groupCoeff( n.CS[:], plus )
  groupMulTerms(n.CS[:])

  sort.Sort( n )
  cnt := countTerms(n.CS[:])
  if cnt == 0 { ret = &Null{} }
  if cnt == 1 { ret = n.CS[0]; n.CS[0] = nil }
  return ret
}




func countTerms( terms []Expr ) int {
  cnt := 0
  for _,e := range terms {
    if e != nil {
      cnt++
    }
  }
  return cnt
}

// this function left aligns children terms ( ie move nils to end of terms[] )
// and returns the number of children
func leftAlignTerms( terms []Expr ) int {
  cnt, nilp := 0, -1
  for i,e := range terms {
    if e != nil {
      //           fmt.Printf( "TERMS(%d/%d): %v\n", i,nilp, terms )
      cnt++
      if nilp >= 0 {
        terms[nilp], terms[i] = terms[i], nil
        nilp++
        // find next nil spot
        for nilp < len(terms) {
          if terms[nilp] == nil { break } else { nilp++ }
        }
        if nilp >= len(terms) {break} // the outer loop
      }
    } else if nilp < 0 {
      nilp = i
    }
  }
  return cnt
}

func plus(lhs,rhs float64) float64{return lhs+rhs}
func mult(lhs,rhs float64) float64{return lhs*rhs}


func groupCoeff( terms []Expr, fold(func (lhs,rhs float64) float64) ) {
  var (
    fC *ConstantF = nil // first coeff (nil until found)
  )
  for i,t := range terms {
    if t == nil { continue }
    if t.ExprType() == CONSTANTF {
      if fC == nil {
        fC = t.(*ConstantF)
      } else {
        fC.F = fold(fC.F,t.(*ConstantF).F)
        terms[i] = nil
      }
    }
  }
}



func gatherAdds( n* Add ) {
  terms := n.CS[:]
  nops := countTerms(terms)
  for i,e := range terms {
    if e == nil { continue }
    if e.ExprType() == ADD {
      a := e.(*Add)
      cnt := leftAlignTerms(a.CS[:])
      for j:=0; j < cnt && nops < MaxAddChildren; j++ {
        terms[nops] = a.CS[j]
        a.CS[j] = nil
        nops++
      }
      rem := leftAlignTerms(a.CS[:])
      if rem == 0 {
        n.CS[i] = nil
      }
      leftAlignTerms(terms)
    }
  }
}

func groupAddTerms( terms []Expr ) {
  l := len(terms)
  for i,t := range terms {
    if t == nil { continue }
    ty := t.ExprType()
    switch ty {
      case VAR,SYSTEM: // ,SYSTEM
        sum := 1.0
        for j := i+1; j < l; j++ {
          if terms[j] == nil {continue}
          // Xi + Xi
          if t.AmISame(terms[j]) {
            sum += 1.0
            terms[j] = nil
            continue
          }
          // Xi + -Xi
          if terms[j].ExprType() == NEG && t.AmISame( terms[j].(*Neg).C ) {
            sum -= 1.0
            terms[j] = nil
            continue
          }
          // Xi + cXi
          if terms[j].ExprType() == MUL {
            m := terms[j].(*Mul)
            nc := leftAlignTerms(m.CS[:])
            if nc == 2 && m.CS[0].ExprType() == CONSTANTF && t.AmISame(m.CS[1]) {
              sum += m.CS[0].(*ConstantF).F
              terms[j] = nil
              continue
            }
          }
        }

        if sum == 0.0 {
          terms[i] = nil
        } else if sum != 1.0 {
          var mul Mul
          mul.CS[0] = &ConstantF{F:sum}
          mul.CS[1] = t
          var e Expr = &mul
          terms[i] = e
        }

    }
  }
}

func groupMulTerms( terms []Expr ) {
  l := len(terms)
  for i,t := range terms {
    if t == nil { continue }
    ty := t.ExprType()
    switch ty {
      case VAR,SYSTEM:
        sum := 1
        is_neg := false
        for j := i+1; j < l; j++ {
          if terms[j] == nil {continue}
          // Xi * Xi
          if t.AmISame(terms[j]) {
            sum += 1
            terms[j] = nil
            continue
          }
          // Xi * -Xi
          if terms[j].ExprType() == NEG && t.AmISame( terms[j].(*Neg).C ) {
            sum += 1
            if !is_neg { is_neg = true } else { is_neg = false }
            terms[j] = nil
            continue
          }
          // Xi * cXi  [NOT NEEDED FOR MUL (won't be seen like this)]
          // Xi * Xi^m
          if terms[j].ExprType() == POWI {
            pow := terms[j].(*PowI)
            // Xi^n
            if t.AmISame( pow.Base ) {
              sum += pow.Power
              if !is_neg { is_neg = true } else { is_neg = false }
              terms[j] = nil
              continue
            }
            // (-Xi)^n
            if pow.Base.ExprType() == NEG && t.AmISame(pow.Base.(*Neg).C ) {
              sum += pow.Power
              // swap is_neg if power is odd
              if pow.Power % 2 == 1 { if !is_neg { is_neg = true } else { is_neg = false } }
              terms[j] = nil
              continue
            }
          }
        }

        if sum == 0 {
          var e Expr = &ConstantF{F:1}
          terms[i] = e
        } else if sum != 1 {
          var pow PowI
          if is_neg {
            pow.Base = &Neg{C:t}
          } else {
            pow.Base = t
          }
          pow.Power = sum
          var e Expr = &pow
          terms[i] = e
        }

    }
  }
}

func gatherMuls( n* Mul ) {
  terms := n.CS[:]
  nops := countTerms(terms)
  for i,e := range terms {
    if e == nil { continue }
    if e.ExprType() == MUL {
      a := e.(*Mul)
      cnt := leftAlignTerms(a.CS[:])
      for j:=0; j < cnt && nops < MaxMulChildren; j++ {
        terms[nops] = a.CS[j]
        a.CS[j] = nil
        nops++
      }
      rem := leftAlignTerms(a.CS[:])
      if rem == 0 {
        n.CS[i] = nil
      }
      leftAlignTerms(terms)
    }
  }
}
