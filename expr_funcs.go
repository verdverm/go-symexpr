package symexpr

import (
  rand "math/rand"
  "fmt"
  "bufio"
)


var (
  roots  = [...]int{ ADD,MUL }
  leafs  = [...]int{ VAR,CONSTANTF }
  opers  = [...]int{ NEG,ABS,SQRT,SIN,COS,EXP,LOG,ADD,MUL,DIV,POWE }
  usable = [...]int{ VAR,CONSTANTF,NEG,ABS,SQRT,SIN,COS,EXP,LOG,ADD,MUL,DIV,POWE }
  intrig = [...]int{ VAR,NEG,ABS,SQRT,EXP,LOG,ADD,MUL,DIV,POWE }
)

type ExprGenParams struct {
  // bounds on tree
  MaxSize, MaxDepth,
  MinSize, MinDepth int
  // tpm bounds on tree (for subtree distributions)
  TmpMaxSize, TmpMaxDepth,
  TmpMinSize, TmpMinDepth int

  // Current values
  CurrSize, CurrDepth int
  InTrig bool

  // bounds on some operators
  NumDim, NumSys, NumCoeff int

  // usable terms at each location
  Roots, Nodes, Leafs, NonTrig, Var []int

}

func (egp *ExprGenParams) CheckExprTmp( e Expr ) bool {
    if e.Size() < egp.TmpMinSize || e.Size() > egp.TmpMaxSize ||
      e.Height() < egp.TmpMinDepth || e.Height() > egp.TmpMaxDepth {
          return false
        }
    return true
}
func (egp *ExprGenParams) CheckExpr( e Expr ) bool {
    if e.Size() < egp.MinSize {
      return false
    } else if e.Size() > egp.MaxSize {
      return false
    } else if  e.Height() < egp.MinDepth {
      return false
    } else if  e.Height() > egp.MaxDepth {
      return false
    }
    return true
}

func (egp *ExprGenParams) CheckExprLog( e Expr, log *bufio.Writer) bool {
//   if e.Size() < egp.TmpMinSize || e.Size() > egp.TmpMaxSize ||
//     e.Height() < egp.TmpMinDepth || e.Height() > egp.TmpMaxDepth {
//       return false
//     }
    if e.Size() < egp.MinSize {
         fmt.Fprintf( log, "Too SMALL:  e:%v  l:%v\n", e.Size(), egp.MinSize )
      return false
    } else if e.Size() > egp.MaxSize {
         fmt.Fprintf( log, "Too LARGE:  e:%v  l:%v\n", e.Size(), egp.MaxSize )
      return false
    } else if  e.Height() < egp.MinDepth {
         fmt.Fprintf( log, "Too SHORT:  e:%v  l:%v\n", e.Height(), egp.MinDepth )
      return false
    } else if  e.Height() > egp.MaxDepth {
         fmt.Fprintf( log, "Too TALL:  e:%v  l:%v\n", e.Height(), egp.MaxDepth )
      return false
    }
    return true
}

func (egp *ExprGenParams) ResetCurr() {
	egp.CurrSize,egp.CurrDepth,egp.InTrig = 0,0,false
}
func (egp *ExprGenParams) ResetTemp() {
  egp.TmpMaxSize, egp.TmpMaxDepth = egp.MaxSize, egp.MaxDepth
  egp.TmpMinSize, egp.TmpMinDepth = egp.MinSize, egp.MinDepth
}

func ExprGen(egp ExprGenParams, srules SimpRules, rng *rand.Rand) Expr {
  var ret Expr
  good := false
  cnt := 0

  for !good {
    egp.ResetCurr()
    egp.ResetTemp()

//     eqn := ExprGenStats( &egp, rng )
    eqn := exprGrow( -1, ExprGenDepth, &egp, rng )
    eqn.CalcExprStats(0)

//     fmt.Printf( "%v\n", eqn)

    ret = eqn.Simplify( srules )
    ret.CalcExprStats(0)

//     fmt.Printf( "%v\n\n", ret)

    // check eqn after simp
    good = egp.CheckExpr(ret)
    cnt++
  }
  return ret
}

type ExprGenFunc (func ( egp *ExprGenParams, rng *rand.Rand) Expr)

func exprGrow( e int, egfunc ExprGenFunc, egp *ExprGenParams, rng *rand.Rand) Expr {

	if e == -1 {
		return egfunc( egp, rng )
	}

	switch {

  case e == TIME:
    return &Time{ExprStats: ExprStats{0,0,0,0}}
  case e == VAR:
    return &Var{P: egp.Var[rng.Intn(len(egp.Var))]} /// this needs changing to parameters value
  case e == CONSTANT:
    return &Constant{P: rng.Int() % egp.NumCoeff}
  case e == CONSTANTF:
    return &ConstantF{F:  rng.NormFloat64()*2}
  case e == SYSTEM:
    return &System{P: rng.Int() % egp.NumSys}

  case e == NEG:
    egp.CurrDepth++
    return &Neg{C: egfunc(egp,rng)}
  case e == ABS:
    egp.CurrDepth++
    return &Abs{C: egfunc(egp,rng)}
  case e == SQRT:
    egp.CurrDepth++
    return &Sqrt{C: egfunc(egp,rng)}
  case e == SIN:
    egp.CurrDepth++
    egp.InTrig = true
    tmp := Sin{C: egfunc(egp,rng)}
    egp.InTrig = false
    return &tmp
  case e == COS:
    egp.CurrDepth++
    egp.InTrig = true
    tmp := Cos{C: egfunc(egp,rng)}
    egp.InTrig = false
    return &tmp
  case e == EXP:
    egp.CurrDepth++
    return &Exp{C: egfunc(egp,rng)}
  case e == LOG:
    egp.CurrDepth++
    return &Log{C: egfunc(egp,rng)}
  case e == POWI:
    egp.CurrDepth++
    return &PowI{Base: egfunc(egp,rng),Power: (rng.Int() % 5)-2}
  case e == POWF:
    egp.CurrDepth++
    return &PowF{Base: egfunc(egp,rng),Power: rng.Float64()*2}

  case e == ADD:
    egp.CurrDepth++
    var add Add
    nchld := 0
    for nchld < 2 || nchld > MaxAddGenChildren {
      nchld = (rng.Int() % MaxAddGenChildren-2)+2
    }
    if nchld < 2 { fmt.Printf( "AHHHH   %d\n\n\n", nchld ) }

    for i:=0; i<nchld; i++ { add.CS[i] = egfunc(egp,rng) }
    return &add
  case e == SUB:
    egp.CurrDepth++
    return &Sub{C1: egfunc(egp,rng),
    C2: egfunc(egp,rng)}
  case e == MUL:
    egp.CurrDepth++
    var mul Mul
    nchld := 0
    for nchld < 2 || nchld > MaxMulGenChildren {
      nchld = (rng.Int() % MaxMulGenChildren-2)+2
    }
    if nchld < 2 { fmt.Printf( "AHHHH   %d\n\n\n", nchld ) }

    for i:=0; i<nchld; i++ { mul.CS[i] = egfunc(egp,rng) }
    return &mul
  case e == DIV:
    egp.CurrDepth++
    return &Div{Numer: egfunc(egp,rng),
    Denom: egfunc(egp,rng)}
  case e == POWE:
    egp.CurrDepth++
    return &PowE{Base: egfunc(egp,rng),
    Power: egfunc(egp,rng)}
  }
  return &Null{}
}

func ExprGenDepth(egp *ExprGenParams, rng *rand.Rand) Expr {
	rnum := rng.Int()
  var e int

  if egp.CurrDepth == 0 {
    e = egp.Roots[ rnum % len(egp.Roots) ]
  } else if egp.CurrDepth >= egp.TmpMaxDepth {
    e = egp.Leafs[ rnum % len(egp.Leafs) ]
  } else {
    if egp.InTrig {
      e = egp.NonTrig[ rnum % len(egp.NonTrig) ]
    } else {
      e = egp.Nodes[ rnum % len(egp.Nodes) ] // this is to deal with NULL (ie so we dont get switch on 0)
    }
  }
  return exprGrow( e, ExprGenDepth, egp, rng )
}


func ExprGenSize(egp *ExprGenParams, rng *rand.Rand) Expr {
  rnum := rng.Int()
  var e int

  if egp.CurrDepth == 0 {
    e = egp.Roots[ rnum % len(egp.Roots) ]
  } else if egp.CurrDepth >= egp.MaxDepth {
    e = egp.Leafs[ rnum % len(egp.Leafs) ]
  } else {
    if egp.InTrig {
      e = egp.NonTrig[ rnum % len(egp.NonTrig) ]
    } else {
      e = egp.Nodes[ rnum % len(egp.Nodes) ] // this is to deal with NULL (ie so we dont get switch on 0)
    }
  }

	return exprGrow( e, ExprGenSize, egp, rng )
}



func CrossEqns_Vanilla( p1, p2 Expr, egp *ExprGenParams, rng *rand.Rand ) Expr {
  eqn := p1.Clone()
  eqn.CalcExprStats(0)

  s1,s2 := rng.Intn(eqn.Size()), rng.Intn(p2.Size())
//   eqn.SetExpr(&s1, p2.GetExpr(&s2).Clone())
  //   eqn.SetExpr(&s1, new_eqn )
  SwapExpr( eqn, p2.GetExpr(&s2).Clone(), s1 )

  eqn.CalcExprStats(0)
  return eqn

}

func InjectEqn_Vanilla( p1 Expr, egp *ExprGenParams, rng *rand.Rand ) Expr {
  eqn := p1.Clone()
  eqn.CalcExprStats(0)

  s1 := rng.Intn(eqn.Size())
  s2 := s1
  e2 := eqn.GetExpr(&s2)

  egp.CurrSize = eqn.Size() - e2.Size()
  egp.CurrDepth = e2.Depth()
  egp.ResetTemp()

  // not correct (should be size based)
  new_eqn := exprGrow( -1, ExprGenDepth, egp, rng )
//   eqn.SetExpr(&s1, new_eqn )
  SwapExpr( eqn, new_eqn, s1 )

  eqn.CalcExprStats(0)
  return eqn
}


func InjectEqn_50_150( p1 Expr, egp *ExprGenParams, rng *rand.Rand ) Expr {
  eqn := p1.Clone()
  eqn.CalcExprStats(0)


  // begin loop
  s1 := rng.Intn(eqn.Size())
  s1_tmp := s1
  e1 := eqn.GetExpr(&s1_tmp)

  egp.ResetCurr()
  egp.ResetTemp()
  egp.TmpMinSize = e1.Size() / 2
  egp.TmpMaxSize = (e1.Size()*3)/2
  // loop if min/max out of bounds
  // and select new subtree

  // not correct (should be size based)
  new_eqn := exprGrow( -1, ExprGenDepth, egp, rng )
//   eqn.SetExpr(&s1, new_eqn )
  SwapExpr( eqn, new_eqn, s1 )
  eqn.CalcExprStats(0)
  return eqn
}


func InjectEqn_SubtreeFair( p1 Expr, egp *ExprGenParams, rng *rand.Rand ) Expr {
  eqn := p1.Clone()
  eqn.CalcExprStats(0)

  // begin loop
  s1, s2:= rng.Intn(eqn.Size()), rng.Intn(eqn.Size())
  s2_tmp := s2
  e2 := eqn.GetExpr(&s2_tmp)

  egp.ResetCurr()
  egp.ResetTemp()
  egp.TmpMinSize = e2.Size() / 2
  egp.TmpMaxSize = (e2.Size()*3)/2
  // loop if min/max out of bounds
  // and select new subtree

  // not correct (should be size based)
  new_eqn := exprGrow( -1, ExprGenDepth, egp, rng )
  //   eqn.SetExpr(&s1, new_eqn )
  SwapExpr( eqn, new_eqn, s1 )

  eqn.CalcExprStats(0)
  return eqn
}

func MutateEqn_Vanilla( eqn Expr, egp *ExprGenParams, rng *rand.Rand, sysvals []float64 ) {
  mut := false
  for !mut {
    s2:= rng.Intn(eqn.Size())
    s2_tmp := s2
    e2 := eqn.GetExpr(&s2_tmp)

    t := e2.ExprType()


    switch t {
      case CONSTANTF:
        if egp.NumSys == 0 {
          e2.(*ConstantF).F += rng.NormFloat64()
        } else {
          // mod coeff
					if rng.Intn(2) == 0 {
						e2.(*ConstantF).F += rng.NormFloat64()
					} else {
						// coeff -> system
						var mul Mul
						s := &System{P: rng.Int() % egp.NumSys}
						e2.(*ConstantF).F /= sysvals[s.P]
						mul.CS[0] = s
						mul.CS[1] = e2.(*ConstantF)
						e2 = &mul
					}
        }
        mut = true
      case SYSTEM:
        e2.(*System).P = rng.Int() % egp.NumSys
        mut = true
      case VAR:
// 				if len(egp.Var) < 2 { continue }
				e2.(*Var).P = egp.Var[rng.Intn(len(egp.Var))]
        mut = true
      case ADD:

      case MUL:

    }
  }
}

/*
void mutateEqn( gsl_rng * rand_gen, Params* params, EqnBlock* eqn, int xn, double* sys_data ) {

  Time_ simper;
//   simper->simplify(eqn);
  eqn->calc_size();
  eqn->calc_depth();

  unsigned long int rnum = -1;
  int pos2 = -1;
  int sz1 = -1;

  EqnBlock *p1, *p2;

  p1 = eqn;
  sz1 = p1->size();

  rnum = gsl_rng_get(rand_gen);
  pos2 = rnum % sz1;

  p2 = p1->getNode(pos2);

  switch( p2->getType() ) {
    case ExprInfo::TIME:
      break;
    case ExprInfo::CONST: {
      double mutate = gsl_rng_uniform_pos(rand_gen);
      if( mutate > 0.23 ) {
        long ival = (gsl_rng_get(rand_gen) % 100001 ) - 50000;
        double val = double(ival) / 1000.0;
        Const_* c = static_cast<Const_*>(p2);
        c->setVal( c->getVal() * val );
      }
      else if(  params->num_sys > 0 ){
        int mut_type = gsl_rng_get(rand_gen) % 6;
        EqnBlock* exp = 0;
        Const_* c = static_cast<Const_*>(p2);
        Const_* nc = 0;
        Sys_* sys = new Sys_( static_cast<int>(gsl_rng_get(rand_gen) % params->num_sys) );
        switch( mut_type ) {
          case 0: //  nc + s = c -> nc = c - s
            exp = new Add_();
            nc = new Const_( c->getVal() - sys_data[ sys->getIdx() ] );
            exp->addOperand( nc );
            exp->addOperand( sys );
            break;
          case 1: // nc - s = c => nc = c + s
            exp = new Add_();
            nc = new Const_( c->getVal() + sys_data[ sys->getIdx() ] );
            exp->addOperand( nc );
            exp->addOperand( new Neg_(sys) );
            break;
          case 2: // s - nc = c => nc = s - c
            exp = new Add_();
            nc = new Const_( sys_data[ sys->getIdx() ] -  c->getVal());
            exp->addOperand( sys );
            exp->addOperand( nc );
            break;
          case 3: //  nc*s = c => nc = c/s
            exp = new Mul_();
            nc = new Const_( c->getVal() / sys_data[ sys->getIdx() ]);
            exp->addOperand( sys );
            exp->addOperand( nc );
            break;
          case 4: // s/nc = c => nc = s/c
            exp = new Div_();
            nc = new Const_( sys_data[ sys->getIdx() ] /  c->getVal());
            exp->addOperand( sys );
            exp->addOperand( nc );
            break;
          case 5: //  nc/s = c => c*s = nc
            exp = new Div_();
            nc = new Const_( sys_data[ sys->getIdx() ] * c->getVal());
            exp->addOperand( nc );
            exp->addOperand( sys );
            break;
          default:
            abort();
        }
        bool replaced = p1->replaceNode(pos2,exp);
//         if( replaced )
//           delete p2;
      }
    }
    break;
    case ExprInfo::SYS: {
      static_cast<Sys_*>(p2)->setIdx( static_cast<int>(gsl_rng_get(rand_gen) % params->num_sys) );
    }
    break;
    case ExprInfo::VAR: {
      int mut_type = gsl_rng_get(rand_gen)%2;
      Var_* v = static_cast<Var_*>(p2);
      if( mut_type == 0 ) {  // change which variable
        int x;
        x = static_cast<int>(gsl_rng_get(rand_gen) % params->usable_dim[xn].size());
        x = params->usable_dim[xn][x];
        v->setIdx(x);
      }
      else { // make into add or mul
        mut_type = gsl_rng_get(rand_gen)%2;
        // mul by random tree
        if( mut_type == 0 ) {
          EqnBlock* m = new Mul_, *tmp;
          m->addOperand( genTreeGrow( rand_gen, params, xn, params->max_eqn_size-sz1, 0 ) );
          if( p1 == p2 ) {
            m->addOperand(p1);
            eqn = m;
          }
          else if( p2->getParent() == NULL ) {
            m->addOperand(p2->clone(false));
            p1->replaceNode(pos2,m);
          }
          else {
            tmp = p2->getParent();
            m->addOperand(tmp->removeOperand(p2));
            tmp->addOperand(m);
          }
        }
        // add to random tree
        else {
          EqnBlock* m = new Add_, *tmp;
          m->addOperand( genTreeGrow( rand_gen, params, xn, params->max_eqn_size-sz1, 0 ) );
          if( p1 == p2 ) {
            m->addOperand(p1);
            eqn = m;
          }
          else if( p2->getParent() == NULL ) {
            m->addOperand(p2->clone(false));
            p1->replaceNode(pos2,m);
          }
          else {
            tmp = p2->getParent();
            m->addOperand(tmp->removeOperand(p2));
            tmp->addOperand(m);
          }
        }
        // mul by const
        // mul by sys
        // mul by variable
        // add to variable
      }
      break;
    }
    case ExprInfo::MUL: {
      int mut_type = gsl_rng_get(rand_gen)%3;
      if( mut_type == 0 ) {
        EqnBlock* b = genTreeGrow( rand_gen, params, xn, params->max_eqn_size-sz1, 0 );
        bool added = p2->addOperand( b );
        if( !added ) {
          delete b;
        }
        b = NULL;
      }
      else if( mut_type == 1 ) {
        EqnBlock* m = new Add_, *tmp;
        m->addOperand( genTreeGrow( rand_gen, params, xn, params->max_eqn_size-sz1, 0 ) );
        if( p1 == p2 ) {
          m->addOperand(p1);
          eqn = m;
        }
        else if( p2->getParent() == NULL ) {
          m->addOperand(p2->clone(false));
          p1->replaceNode(pos2,m);
        }
        else {
          tmp = p2->getParent();
          m->addOperand(tmp->removeOperand(p2));
          tmp->addOperand(m);
        }
      }
      else {  // change mul to add
        EqnBlock* a = new Add_, *tmp;
        Mul_* m = static_cast<Mul_*>(p2);
        for( int i = m->num_ops()-1; i >= 0; --i ) {
          a->addOperand( m->removeOperand(i) );
        }
        if( p1 == p2 ) {
          a->addOperand(p1);
          eqn = a;
        }
        else if( p2->getParent() == NULL ) {
          a->addOperand(p2->clone(false));
          p1->replaceNode(pos2,a);
        }
        else {
          tmp = p2->getParent();
          a->addOperand(tmp->removeOperand(p2));
          tmp->addOperand(a);
        }
      }
      break;
    }
    case ExprInfo::ADD: {
      int mut_type = gsl_rng_get(rand_gen)%3;
      if( mut_type == 0 ) {
        EqnBlock* b = genTreeGrow( rand_gen, params, xn, params->max_eqn_size-sz1, 0 );
        bool added = p2->addOperand( b );
        if( !added ) {
          delete b;
        }
        b = NULL;
      }
      else if( mut_type == 1 ) {
        EqnBlock* m = new Mul_, *tmp;
        m->addOperand( genTreeGrow( rand_gen, params, xn, params->max_eqn_size-sz1, 0 ) );
        if( p1 == p2 ) {
          m->addOperand(p1);
          eqn = m;
        }
        else if( p2->getParent() == NULL ) {
          m->addOperand(p2->clone(false));
          p1->replaceNode(pos2,m);
        }
        else {
          tmp = p2->getParent();
          m->addOperand(tmp->removeOperand(p2));
          tmp->addOperand(m);
        }
      }
      else {  // Change add to mul
        EqnBlock* m = new Mul_, *tmp;
        Add_* a = static_cast<Add_*>(p2);
        for( int i = a->num_ops()-1; i >= 0; --i ) {
          m->addOperand( a->removeOperand(i) );
        }
        if( p1 == p2 ) {
          m->addOperand(p1);
          eqn = m;
        }
        else if( p2->getParent() == NULL ) {
          m->addOperand(p2->clone(false));
          p1->replaceNode(pos2,m);
        }
        else {
          tmp = p2->getParent();
          m->addOperand(tmp->removeOperand(p2));
          tmp->addOperand(m);
        }
      }
      break;
    }
    case ExprInfo::POW1: {
      int e = static_cast<int>(gsl_rng_get(rand_gen) % 6 ) - 2;
      static_cast<Pow1_*>(p2)->setPower(e);
    }
    break;
    default:
      ;
  };
}
*/
