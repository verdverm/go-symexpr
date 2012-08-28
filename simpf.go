package symexpr

// import (
//   "math"
//   "sort"
// )
// 
// 
type SimpRulesF struct {
    EvalUnary bool
    GroupCoeff bool

    MulToPow bool
    MulInCoeff int // add must be >= to size, 0 == No
}

// func (l *Leaf) Simplify( rules SimpRules ) Expr      { return nil }
// func (u *Unary) Simplify( rules SimpRules ) Expr     { return nil }
// func (n *N_ary) Simplify( rules SimpRules ) Expr     { return nil }
// 
// 
// func (n *Time) Simplify( rules SimpRules ) Expr      { return n }
// 
// func (v *Var) Simplify( rules SimpRules ) Expr       { return v }
// 
// func (c *Constant) Simplify( rules SimpRules ) Expr  { return c }
// 
// func (c *ConstantF) Simplify( rules SimpRules ) Expr {
//   if math.IsNaN(c.F) || math.IsInf(c.F,0) {
//     return &Null{}
//   }
//   // "close" to system value ??? TODO
// 
//   return c
// }
// 
// func (s *System) Simplify( rules SimpRules ) Expr    { return s }
// 
// 
// 
// 
// func (u *Neg) Simplify( rules SimpRules ) Expr {
//   var (
//     ret Expr = u
//     t int = NULL
//   )
//   if u.C != nil {
//     u.C = u.C.Simplify(rules)
//     t = u.C.ExprType()
//   }
//   switch t {
//     case NULL, NEG:
//       ret = u.C
//       u.C = nil
//     case CONSTANTF:
//       ret = u.C
//       u.C = nil
//       ret.(*ConstantF).F *= -1.0
//     case MUL:
//       m := u.C.(*Mul)
//       if m.CS[0].ExprType() == CONSTANTF {
//         m.CS[0].(*ConstantF).F *= -1.0
//         ret = u.C
//         u.C = nil
//       }
//   }
//   return ret
// }
// 
// 
// func (u *Abs) Simplify( rules SimpRules ) Expr {
//   var (
//     ret Expr = u
//     t int = NULL
//   )
//   if u.C != nil {
//     u.C = u.C.Simplify(rules)
//     t = u.C.ExprType()
//   }
//   switch t {
//     case NULL, ABS:
//       ret = u.C
//       u.C = nil
//     case CONSTANTF:
//       ret = u.C
//       u.C = nil
//       ret.(*ConstantF).F = math.Abs( ret.(*ConstantF).F )
//   }
//   return ret
// }
// 
// 
// func (u *Sqrt) Simplify( rules SimpRules ) Expr {
//   var (
//     ret Expr = u
//     t int = NULL
//   )
//   if u.C != nil {
//     u.C = u.C.Simplify(rules)
//     t = u.C.ExprType()
//   }
//   switch t {
//     case NULL:
//       ret = u.C
//       u.C = nil
//     case CONSTANTF:
//       ret = u.C
//       u.C = nil
//       ret.(*ConstantF).F = math.Sqrt( ret.(*ConstantF).F )
//   }
//   return ret
// }
// 
// 
// func (u *Sin) Simplify( rules SimpRules ) Expr    {
//   var (
//     ret Expr = u
//     t int = NULL
//   )
//   if u.C != nil {
//     u.C = u.C.Simplify(rules)
//     t = u.C.ExprType()
//   }
//   switch t {
//     case NULL,SIN,COS:
//       ret = u.C
//       u.C = nil
//     case CONSTANTF:
//       ret = u.C
//       u.C = nil
//       ret.(*ConstantF).F = math.Sin( ret.(*ConstantF).F )
//   }
//   return ret
// }
// 
// 
// func (u *Cos) Simplify( rules SimpRules ) Expr    {
//   var (
//     ret Expr = u
//     t int = NULL
//   )
//   if u.C != nil {
//     u.C = u.C.Simplify(rules)
//     t = u.C.ExprType()
//   }
//   switch t {
//     case NULL,SIN,COS:
//       ret = u.C
//       u.C = nil
//     case CONSTANTF:
//       ret = u.C
//       u.C = nil
//       ret.(*ConstantF).F = math.Cos( ret.(*ConstantF).F )
//   }
//   return ret
// }
// 
// 
// func (u *Exp) Simplify( rules SimpRules ) Expr    {
//   var (
//     ret Expr = u
//     t int = NULL
//   )
//   if u.C != nil {
//     u.C = u.C.Simplify(rules)
//     t = u.C.ExprType()
//   }
//   switch t {
//     case NULL:
//       ret = u.C
//       u.C = nil
//     case CONSTANTF:
//       ret = u.C
//       u.C = nil
//       ret.(*ConstantF).F = math.Exp( ret.(*ConstantF).F )
//   }
//   return ret
// }
// 
// 
// func (u *Log) Simplify( rules SimpRules ) Expr    {
//   var (
//     ret Expr = u
//     t int = NULL
//   )
//   if u.C != nil {
//     u.C = u.C.Simplify(rules)
//     t = u.C.ExprType()
//   }
//   switch t {
//     case NULL:
//       ret = u.C
//       u.C = nil
//     case CONSTANTF:
//       ret = u.C
//       u.C = nil
//       ret.(*ConstantF).F = math.Log( ret.(*ConstantF).F )
//   }
//   return ret
// }
// 
// 
// func (u *PowI) Simplify( rules SimpRules ) Expr     {
//   var (
//     ret Expr = u
//     t int = NULL
//   )
//   if u.Base != nil {
//     u.Base = u.Base.Simplify(rules)
//     t = u.Base.ExprType()
//   }
//   if u.Power == 0 {
//     ret = &ConstantF{F:1}
//   } else if u.Power == 1 {
//     ret = u.Base
//     u.Base = nil
//   } else {
//     switch t {
//       case NULL:
//         ret = u.Base
//         u.Base = nil
//       case CONSTANTF:
//         ret = u.Base
//         u.Base = nil
//         ret.(*ConstantF).F = math.Pow( ret.(*ConstantF).F, float64(u.Power) )
//     }
//   }
//   return ret
// }
// 
// 
// func (u *PowF) Simplify( rules SimpRules ) Expr     {
//   var (
//     ret Expr = u
//     t int = NULL
//   )
//   if u.Base != nil {
//     u.Base = u.Base.Simplify(rules)
//     t = u.Base.ExprType()
//   }
//   switch t {
//     case NULL:
//       ret = u.Base
//       u.Base = nil
//     case CONSTANTF:
//       ret = u.Base
//       u.Base = nil
//       ret.(*ConstantF).F = math.Pow( ret.(*ConstantF).F, float64(u.Power) )
//   }
//   return ret
// }
// 
// 
// 
// 
// func (n *PowE) Simplify( rules SimpRules ) Expr     {
//   var (
//     ret Expr = n
//     t1, t2 int = NULL,NULL
//   )
//   if n.Base != nil {
//     n.Base = n.Base.Simplify(rules)
//     t1 = n.Base.ExprType()
//   }
//   if n.Power != nil {
//     n.Power = n.Power.Simplify(rules)
//     t2 = n.Power.ExprType()
//   }
//   if t1 == NULL && t2 == NULL {
//     return &Null{}
//   } else if t1 == NULL {
//      ret = n.Power
//      n.Base = nil
//      n.Power = nil
//   } else if t2 == NULL {
//      ret = n.Base
//      n.Base = nil
//      n.Power = nil
//   } else if n.Base.ExprType() == n.Power.ExprType() &&
//      n.Base.ExprType() == CONSTANTF {
//      ret = n.Base
//      ret.(*ConstantF).F = math.Pow( ret.(*ConstantF).F, n.Power.(*ConstantF).F )
//      n.Base = nil
//      n.Power = nil
//   }
//   return ret
// }
// 
// func (n *Div) Simplify( rules SimpRules ) Expr      {
//   var (
//     ret Expr = n
//     t1, t2 int = NULL,NULL
//   )
//   if n.Numer != nil {
//     n.Numer = n.Numer.Simplify(rules)
//     t1 = n.Numer.ExprType()
//   }
//   if n.Denom != nil {
//     n.Denom = n.Denom.Simplify(rules)
//     t2 = n.Denom.ExprType()
//   }
// 
//   if t1 == NULL && t2 == NULL {
//     return &Null{}
//   } else if t1 == NULL {
//      ret = n.Denom
//      n.Numer = nil
//      n.Denom = nil
//   } else if t2 == NULL {
//      ret = n.Numer
//      n.Numer = nil
//      n.Denom = nil
//   } else if n.Numer.ExprType() == n.Denom.ExprType() &&
//      n.Numer.ExprType() == CONSTANTF {
//      ret = n.Numer
//      ret.(*ConstantF).F /= n.Denom.(*ConstantF).F
//      n.Numer = nil
//      n.Denom = nil
//   }
//   return ret
// }
// 
// 
// 
// 
// func (n *Add) Simplify( rules SimpRules ) Expr      {
//   var ret Expr = n
//   for i,C := range n.CS {
//     if C != nil {
//       n.CS[i] = C.Simplify(rules)
//       if n.CS[i] == nil { continue }
//       if n.CS[i].ExprType() == NULL { n.CS[i] = nil }
//     }
//   }
// 
//   groupCoeff( n.CS[:], plus )
//   cnt := leftAlignTerms(n.CS[:])
//   n.CS = n.CS[:cnt]
//   
//   groupAddTerms(n.CS[:])
//   sort.Sort( n )
// 
//   gatherAdds(n)
//   sort.Sort( n )
// 
//   groupCoeff( n.CS[:], plus )
//   cnt = leftAlignTerms(n.CS[:])
//   n.CS = n.CS[:cnt]
//   
//   groupAddTerms(n.CS[:])
//   sort.Sort( n )
// 
// 
//   cnt = countTerms(n.CS[:])
//   n.CS = n.CS[:cnt]
//   if cnt > 0 {
//     if n.CS[0].ExprType() == CONSTANTF {
//       f := math.Abs(n.CS[0].(*ConstantF).F)
//       if f < 0.00000001 {
//         n.CS[0] = nil
//         n.CS = n.CS[1:]
//         cnt--
//       }
//     }
//   }
// 
//   if cnt == 0 {
//     ret = &Null{}
//   } else if cnt == 1 {
//     ret = n.CS[0]
//     n.CS[0] = nil
//   }
// 
//   return ret
// }
// 
// func (n *Mul) Simplify( rules SimpRules ) Expr      {
//   var ret Expr = n
//   for i,C := range n.CS {
//     if C != nil {
//       n.CS[i] = C.Simplify(rules)
//       if n.CS[i].ExprType() == NULL { n.CS[i] = nil }
//     }
//   }
//   groupCoeff( n.CS[:], mult )
//   cnt := leftAlignTerms(n.CS[:])
//   n.CS = n.CS[:cnt]
//   
//   groupMulTerms(n.CS[:])
//   sort.Sort( n )
// 
//   gatherMuls(n)
//   sort.Sort( n )
// 
//   groupCoeff( n.CS[:], mult )
//   cnt = leftAlignTerms(n.CS[:])
//   n.CS = n.CS[:cnt]
//   
//   groupMulTerms(n.CS[:])
// 
//   sort.Sort( n )
//   cnt = countTerms(n.CS[:])
//   n.CS = n.CS[:cnt]
//   if cnt == 0 { ret = nil }
//   if cnt == 1 { ret = n.CS[0]; n.CS[0] = nil }
//   if cnt > 1 {
//     if n.CS[0].ExprType() == CONSTANTF {
//       f := math.Abs(n.CS[0].(*ConstantF).F)
//       if f < 0.00000001 {
//         ret = nil
//       }
//     }
//   }
//   return ret
// }
// 
// 
// 
// 
// func countTerms( terms []Expr ) int {
//   cnt := 0
//   for _,e := range terms {
//     if e != nil {
//       cnt++
//     }
//   }
//   return cnt
// }
// 
// // this function left aligns children terms ( ie move nils to end of terms[] )
// // and returns the number of children
// func leftAlignTerms( terms []Expr ) int {
//   cnt, nilp := 0, -1
//   for i,e := range terms {
//     if e != nil {
//       //           fmt.Printf( "TERMS(%d/%d): %v\n", i,nilp, terms )
//       cnt++
//       if nilp >= 0 {
//         terms[nilp], terms[i] = terms[i], nil
//         nilp++
//         // find next nil spot
//         for nilp < len(terms) {
//           if terms[nilp] == nil { break } else { nilp++ }
//         }
//         if nilp >= len(terms) {break} // the outer loop
//       }
//     } else if nilp < 0 {
//       nilp = i
//     }
//   }
//   return cnt
// }
// 
// func plus(lhs,rhs float64) float64{return lhs+rhs}
// func mult(lhs,rhs float64) float64{return lhs*rhs}
// 
// 
// func groupCoeff( terms []Expr, fold(func (lhs,rhs float64) float64) ) {
//   var (
//     fC *Constant = nil // first coeff (nil until found)
//   )
//   for i,t := range terms {
//     if t == nil { continue }
//     if t.ExprType() == CONSTANT {
//       if fC == nil {
//         fC = t.(*Constant)
//       } else {
//         terms[i] = nil
//       }
//     }
//   }
// }
// 
// // func groupCoeff( terms []Expr, fold(func (lhs,rhs float64) float64) ) {
// //   var (
// //     fC *ConstantF = nil // first coeff (nil until found)
// //   )
// //   for i,t := range terms {
// //     if t == nil { continue }
// //     if t.ExprType() == CONSTANTF {
// //       if fC == nil {
// //         fC = t.(*ConstantF)
// //       } else {
// //         fC.F = fold(fC.F,t.(*ConstantF).F)
// //         terms[i] = nil
// //       }
// //     }
// //   }
// // }
// 
// 
// 
// func gatherAddsF( n* Add ) {
//   terms := make([]Expr,0)
//   for i,e := range n.CS {
//     if e == nil { continue }
//     if e.ExprType() == ADD {
//       a := e.(*Add)
//       leftAlignTerms(a.CS[:])
//       for j,E := range a.CS {
//         if E == nil { continue }
//         terms = append(terms,E)
//         a.CS[j] = nil
//       }
//       rem := leftAlignTerms(a.CS[:])
//       if rem == 0 {
//         n.CS[i] = nil
//       }
//       leftAlignTerms(terms)
//     } else {
//       terms = append(terms,e)
//     }
//   }
//   n.CS = terms
// }
// 
// func groupAddTermsF( terms []Expr ) {
//   l := len(terms)
//   for i,t := range terms {
//     if t == nil { continue }
//     ty := t.ExprType()
//     switch ty {
//       case VAR,SYSTEM: 
//         sum := 1.0
//         for j := i+1; j < l; j++ {
//           if terms[j] == nil {continue}
//           // Xi + Xi
//           if t.AmISame(terms[j]) {
//             sum += 1.0
//             terms[j] = nil
//             continue
//           }
//           // Xi + -Xi
//           if terms[j].ExprType() == NEG && t.AmISame( terms[j].(*Neg).C ) {
//             sum -= 1.0
//             terms[j] = nil
//             continue
//           }
//           // Xi + cXi
//           if terms[j].ExprType() == MUL {
//             m := terms[j].(*Mul)
//             nc := leftAlignTerms(m.CS[:])
//             if nc == 2 && m.CS[0].ExprType() == CONSTANT && t.AmISame(m.CS[1]) {
//               terms[j] = nil
//               continue
//             }
//           }
//         }
// 
//         if sum == 0.0 {
//           terms[i] = nil
//         } else if sum != 1.0 {
//           var mul Mul
//           mul.Insert( NewConstant(-1) )
//           mul.Insert( t )
//           var e Expr = &mul
//           terms[i] = e
//         }
// 	  case MUL:
// 		// cXi + dXi
// 		tm := t.(*Mul)
// 		if tm.CS[0].ExprType() == CONSTANT { continue }
// 		for j := i+1; j < l; j++ {
// 		  if terms[j] == nil {continue}
// 		  if terms[j].ExprType() == MUL {
// 			m := terms[j].(*Mul)
// 			nc := leftAlignTerms(m.CS[:])
// 			if nc == 2 && m.CS[0].ExprType() == CONSTANT && tm.CS[1].AmISame(m.CS[1]) {
// 			  terms[j] = nil
// 			  continue
// 			}
// 		  }
// 		}
//     }
//   }
// }
// 
// func groupMulTermsF( terms []Expr ) {
//   l := len(terms)
//   for i,t := range terms {
//     if t == nil { continue }
//     ty := t.ExprType()
//     switch ty {
//       case VAR,SYSTEM:
//         sum := 1
//         is_neg := false
//         for j := i+1; j < l; j++ {
//           if terms[j] == nil {continue}
//           // Xi * Xi
//           if t.AmISame(terms[j]) {
//             sum += 1
//             terms[j] = nil
//             continue
//           }
//           // Xi * -Xi
//           if terms[j].ExprType() == NEG && t.AmISame( terms[j].(*Neg).C ) {
//             sum += 1
//             if !is_neg { is_neg = true } else { is_neg = false }
//             terms[j] = nil
//             continue
//           }
//           // Xi * cXi  [NOT NEEDED FOR MUL (won't be seen like this)]
//           // Xi * Xi^m
//           if terms[j].ExprType() == POWI {
//             pow := terms[j].(*PowI)
//             // Xi^n
//             if t.AmISame( pow.Base ) {
//               sum += pow.Power
//               if !is_neg { is_neg = true } else { is_neg = false }
//               terms[j] = nil
//               continue
//             }
//             // (-Xi)^n
//             if pow.Base.ExprType() == NEG && t.AmISame(pow.Base.(*Neg).C ) {
//               sum += pow.Power
//               // swap is_neg if power is odd
//               if pow.Power % 2 == 1 { if !is_neg { is_neg = true } else { is_neg = false } }
//               terms[j] = nil
//               continue
//             }
//           }
//         }
// 
//         if sum == 0 {
//           var e Expr = NewConstantF(1.0)
//           terms[i] = e
//         } else if sum != 1 {
//           var pow PowI
//           if is_neg {
//             pow.Base = NewNeg(t)
//           } else {
//             pow.Base = t
//           }
//           pow.Power = sum
//           var e Expr = &pow
//           terms[i] = e
//         }
// 
//     }
//   }
// }
// 
// func gatherMuls( n* Mul ) {
//   terms := make([]Expr,0)
//   for i,e := range n.CS {
//     if e == nil { continue }
//     if e.ExprType() == MUL {
//       a := e.(*Mul)
//       leftAlignTerms(a.CS[:])
//       for j,E := range a.CS {
//         if E == nil { continue }
//         terms = append(terms,E)
//         a.CS[j] = nil
//       }
//       rem := leftAlignTerms(a.CS[:])
//       if rem == 0 {
//         n.CS[i] = nil
//       }
//       leftAlignTerms(terms)
//     } else {
//       terms = append(terms,e)
//     }
//   }
//   n.CS = terms
// }
// 
