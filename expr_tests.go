package symexpr

import (
  "fmt"
  "math"
)


func TestPrint() {

  c := &ConstantF{F: 1.57}
  s := &Sin{C:c}
  fmt.Printf( "Hello Me!   %v\n", s )

  c.F = 90.0
  fmt.Printf( "Hello Me!   %v\n", c )

}

func TestSimp() {
  sRules := SimpRules{true,true}

  testAddSimp(sRules)

}



func TestAmIAlmostSame() {
	var add1,add2,add3,add4 Add
	add1.CS[0] = &ConstantF{F: 4.20}
	add1.CS[1] = &Var{P:1}
	add2.CS[0] = &ConstantF{F: 4.20}
	add2.CS[1] = &Var{P:1}
	add3.CS[0] = &ConstantF{F: 4.20}
	add3.CS[1] = &Var{P:2}
	add4.CS[0] = &ConstantF{F: -4.20}
	add4.CS[1] = &Var{P:1}

	var a1,a2,a3,a4 Expr
	a1,a2,a3,a4 = &add1,&add2,&add3,&add4

	fmt.Printf( "Testing Almost\n\n" )

	fmt.Printf( "a1: %v\na2: %v\na3: %v\na4: %v\n", a1,a2,a3,a4 )

	fmt.Printf( "a1->a2 = %v (true)\n", a1.AmIAlmostSame(a2) )
	fmt.Printf( "a1->a3 = %v (false)\n", a1.AmIAlmostSame(a3) )
	fmt.Printf( "a1->a4 = %v (true)\n", a1.AmIAlmostSame(a4) )

}

func testAddSimp(srules SimpRules) {
  var add,bdd,cdd Add
  var mul Mul
  mul.CS[0] = &ConstantF{F: 4.20}
  mul.CS[1] = &Var{P:1}
  add.CS[0] = &ConstantF{F: 4.20}
  add.CS[1] = &cdd
  add.CS[2] = &bdd
  add.CS[3] = &mul
  bdd.CS[0] = &Var{P:1}
  bdd.CS[1] = &Var{P:0}
  cdd.CS[0] = &Var{P:1}
  cdd.CS[1] = &Var{P:1}

  var e Expr
  e = &add
//   f = &cdd

  fmt.Printf( "Orig:  %v\n", e )
  es := e.Simplify( srules )
  fmt.Printf( "Simp:  %v\n\n\n", es )
  es.CalcExprStats(0)

/*  for i:= 0; i < es.Size(); i++ {
    p,q := i,i
    fmt.Printf( "%v: %v\n", i, es.GetExpr(&p) )
    et := es.Clone()
    fmt.Printf( "Orig:  %v\n", et )
    f := cdd.Clone()
    eu := SwapExpr( et, f, q )
    if eu != nil {
      et = eu
    } else {
      fmt.Printf( "nil\n" )
    }
    fmt.Printf( "Swap:  %v\n", et )
    ev := et.Simplify(srules)
    fmt.Printf( "Simp:  %v\n\n", ev )
  }*/



}


func TestFlatten() {
  testSimpPendFlatten0()  // x0
  testSimpPendFlatten1()  // x1
  testSimpPendFlatten2()  // x1 on numerical derivatives

}

func testSimpPendFlatten0() {
  fmt.Printf( "Flattening the Pendulum x0\n" )

  var (
    origData DataSet
    flat_01  DataSet
    flat_02  DataSet
    flat_03  DataSet
    flat_04  DataSet
  )

  err := origData.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_01.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_02.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_03.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_04.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }


  var xTmp [2]float64
  sys := origData.SysVals()

  x0 := &Var{P:0}
  x1 := &Var{P:1}
  sin0 := &Sin{C: x0}
  sin1 := &Sin{C: x1}

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1,p2 := origData.Point(i),origData.Point(i+1)
    m1,m2 := flat_01.Point(i),flat_01.Point(i+1)
    t1,t2 := p1.Time(), p2.Time()
    v1,v2 := p1.Vals(), p2.Vals()
    out := PRK4( 0, x0, t1,t2, v1,v2, xTmp[:],nil,sys )
    mod := m2.Val(0) - (m1.Val(0) + out)
    m1.SetVal(0,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1,p2 := origData.Point(i),origData.Point(i+1)
    m1,m2 := flat_02.Point(i),flat_02.Point(i+1)
    t1,t2 := p1.Time(), p2.Time()
    v1,v2 := p1.Vals(), p2.Vals()
    out := PRK4( 0, x1, t1,t2, v1,v2, xTmp[:],nil,sys )
    mod := m2.Val(0) - (m1.Val(0) + out)
    m1.SetVal(0,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1,p2 := origData.Point(i),origData.Point(i+1)
    m1,m2 := flat_03.Point(i),flat_03.Point(i+1)
    t1,t2 := p1.Time(), p2.Time()
    v1,v2 := p1.Vals(), p2.Vals()
    out := PRK4( 0, sin0, t1,t2, v1,v2, xTmp[:],nil,sys )
    mod := m2.Val(0) - (m1.Val(0) + out)
    m1.SetVal(0,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1,p2 := origData.Point(i),origData.Point(i+1)
    m1,m2 := flat_04.Point(i),flat_04.Point(i+1)
    t1,t2 := p1.Time(), p2.Time()
    v1,v2 := p1.Vals(), p2.Vals()
    out := PRK4( 0, sin1, t1,t2, v1,v2, xTmp[:],nil,sys )
    mod := m2.Val(0) - (m1.Val(0) + out)
    m1.SetVal(0,mod)
  }



  m,v,s := calcMeanStdDev(origData,0)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m,v,s )

  m1,v1,s1 := calcMeanStdDev(flat_01,0)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m1,v1,s1 )

  m2,v2,s2 := calcMeanStdDev(flat_02,0)
  fmt.Printf( "%.4f  %.4f  %.4f   <= Real Ans\n", m2,v2,s2 )

  m3,v3,s3 := calcMeanStdDev(flat_03,0)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m3,v3,s3 )

  m4,v4,s4 := calcMeanStdDev(flat_04,0)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m4,v4,s4 )




}






func testSimpPendFlatten1() {
  fmt.Printf( "Flattening the Pendulum x1\n" )

  var (
    origData DataSet
    flat_01  DataSet
    flat_02  DataSet
    flat_03  DataSet
    flat_04  DataSet
    flat_05  DataSet
    flat_06  DataSet
  )

  err := origData.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_01.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_02.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_03.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_04.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_05.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_06.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }



  var xTmp [2]float64
  sys := origData.SysVals()

  x0 := &Var{P:0}
  sin := &Sin{C: x0}
  div1 := &Div{Numer: &ConstantF{F:-9.8}, Denom: &ConstantF{F:1.0}}
  div2 := &Div{Numer: &ConstantF{F:1.0},  Denom: &System{P:1}}
  divR := &Div{Numer: &ConstantF{F:-9.8}, Denom: &System{P:1}}

  var mul Mul
  mul.CS[0] = divR
  mul.CS[1] = sin

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1,p2 := origData.Point(i),origData.Point(i+1)
    m1,m2 := flat_01.Point(i),flat_01.Point(i+1)
    t1,t2 := p1.Time(), p2.Time()
    v1,v2 := p1.Vals(), p2.Vals()
    out := PRK4( 1, x0, t1,t2, v1,v2, xTmp[:],nil,sys )
    mod := m2.Val(1) - (m1.Val(1) + out)
    m1.SetVal(1,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1,p2 := origData.Point(i),origData.Point(i+1)
    m1,m2 := flat_02.Point(i),flat_02.Point(i+1)
    t1,t2 := p1.Time(), p2.Time()
    v1,v2 := p1.Vals(), p2.Vals()
    out := PRK4( 1, sin, t1,t2, v1,v2, xTmp[:],nil,sys )
    mod := m2.Val(1) - (m1.Val(1) + out)
    m1.SetVal(1,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1,p2 := origData.Point(i),origData.Point(i+1)
    m1,m2 := flat_03.Point(i),flat_03.Point(i+1)
    t1,t2 := p1.Time(), p2.Time()
    v1,v2 := p1.Vals(), p2.Vals()
    out := PRK4( 1, div1, t1,t2, v1,v2, xTmp[:],nil,sys )
    mod := m2.Val(1) - (m1.Val(1) + out)
    m1.SetVal(1,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1,p2 := origData.Point(i),origData.Point(i+1)
    m1,m2 := flat_04.Point(i),flat_04.Point(i+1)
    t1,t2 := p1.Time(), p2.Time()
    v1,v2 := p1.Vals(), p2.Vals()
    out := PRK4( 1, div2, t1,t2, v1,v2, xTmp[:],nil,sys )
    mod := m2.Val(1) - (m1.Val(1) + out)
    m1.SetVal(1,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1,p2 := origData.Point(i),origData.Point(i+1)
    m1,m2 := flat_05.Point(i),flat_05.Point(i+1)
    t1,t2 := p1.Time(), p2.Time()
    v1,v2 := p1.Vals(), p2.Vals()
    out := PRK4( 1, divR, t1,t2, v1,v2, xTmp[:],nil,sys )
    mod := m2.Val(1) - (m1.Val(1) + out)
    m1.SetVal(1,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1,p2 := origData.Point(i),origData.Point(i+1)
    m1,m2 := flat_06.Point(i),flat_06.Point(i+1)
    t1,t2 := p1.Time(), p2.Time()
    v1,v2 := p1.Vals(), p2.Vals()
    out := PRK4( 1, &mul, t1,t2, v1,v2, xTmp[:],nil,sys )
    mod := m2.Val(1) - (m1.Val(1) + out)
    m1.SetVal(1,mod)
  }


  m,v,s := calcMeanStdDev(origData,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m,v,s )

  m1,v1,s1 := calcMeanStdDev(flat_01,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m1,v1,s1 )

  m2,v2,s2 := calcMeanStdDev(flat_02,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m2,v2,s2 )

  m3,v3,s3 := calcMeanStdDev(flat_03,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m3,v3,s3 )

  m4,v4,s4 := calcMeanStdDev(flat_04,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m4,v4,s4 )

  m5,v5,s5 := calcMeanStdDev(flat_05,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m5,v5,s5 )

  m6,v6,s6 := calcMeanStdDev(flat_06,1)
  fmt.Printf( "%.4f  %.4f  %.4f  <= Real Ans\n", m6,v6,s6 )




}




func testSimpPendFlatten2() {
  fmt.Printf( "Flattening the Pendulum  dx1/dt\n" )

  var (
    origData DataSet
    flat_01  DataSet
    flat_02  DataSet
    flat_03  DataSet
    flat_04  DataSet
    flat_05  DataSet
    flat_06  DataSet
  )

  err := origData.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_c.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_01.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_cd.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_02.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_cd.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_03.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_cd.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_04.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_cd.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_05.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_cd.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }
  err = flat_06.ReadDiffeq("/home/aworm1/data/diffeq/PendulumSimp/PendulumSimp_6_0_cd.txt")
  if err != nil {
    fmt.Printf( "Error:  %v\n", err )
  }



  sys := origData.SysVals()

  x0 := &Var{P:0}
  sin := &Sin{C: x0}
  div1 := &Div{Numer: &ConstantF{F:-9.8}, Denom: &ConstantF{F:1.0}}
  div2 := &Div{Numer: &ConstantF{F:1.0},  Denom: &System{P:1}}
  divR := &Div{Numer: &ConstantF{F:-9.8}, Denom: &System{P:1}}

  var mul Mul
  mul.CS[0] = divR
  mul.CS[1] = sin

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1 := origData.Point(i)
    m1,m2 := flat_01.Point(i),flat_01.Point(i+1)
    t1 := p1.Time()
    v1 := p1.Vals()
    out := x0.Eval(t1,v1,nil,sys)
    mod := m2.Val(1) - (m1.Val(1) + out)
    m1.SetVal(1,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1 := origData.Point(i)
    m1,m2 := flat_02.Point(i),flat_02.Point(i+1)
    t1 := p1.Time()
    v1 := p1.Vals()
    out := sin.Eval(t1,v1,nil,sys)
    mod := m2.Val(1) - (m1.Val(1) + out)
    m1.SetVal(1,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1 := origData.Point(i)
    m1,m2 := flat_03.Point(i),flat_03.Point(i+1)
    t1 := p1.Time()
    v1 := p1.Vals()
    out := div1.Eval(t1,v1,nil,sys)
    mod := m2.Val(1) - (m1.Val(1) + out)
    m1.SetVal(1,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1 := origData.Point(i)
    m1,m2 := flat_04.Point(i),flat_04.Point(i+1)
    t1 := p1.Time()
    v1 := p1.Vals()
    out := div2.Eval(t1,v1,nil,sys)
    mod := m2.Val(1) - (m1.Val(1) + out)
    m1.SetVal(1,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1 := origData.Point(i)
    m1,m2 := flat_05.Point(i),flat_05.Point(i+1)
    t1 := p1.Time()
    v1 := p1.Vals()
    out := divR.Eval(t1,v1,nil,sys)
    mod := m2.Val(1) - (m1.Val(1) + out)
    m1.SetVal(1,mod)
  }

  for i := 0; i < flat_01.NumPoints()-1; i++ {
    p1 := origData.Point(i)
//     m1,m2 := flat_06.Point(i),flat_06.Point(i+1)
    m1 := flat_06.Point(i)
    t1 := p1.Time()
    v1 := p1.Vals()
    out := mul.Eval(t1,v1,nil,sys)
//     mod := m2.Val(1) - (m1.Val(1) + out)
    mod := (m1.Val(1) - out)
//     fmt.Printf( "%.2f  %.5f      %.5f   %.5f\n", t1, p1.Val(1), m1.Val(1), out )
    m1.SetVal(1,mod)
  }


  m,v,s := calcMeanStdDev(origData,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m,v,s )

  m1,v1,s1 := calcMeanStdDev(flat_01,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m1,v1,s1 )

  m2,v2,s2 := calcMeanStdDev(flat_02,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m2,v2,s2 )

  m3,v3,s3 := calcMeanStdDev(flat_03,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m3,v3,s3 )

  m4,v4,s4 := calcMeanStdDev(flat_04,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m4,v4,s4 )

  m5,v5,s5 := calcMeanStdDev(flat_05,1)
  fmt.Printf( "%.4f  %.4f  %.4f\n", m5,v5,s5 )

  m6,v6,s6 := calcMeanStdDev(flat_06,1)
  fmt.Printf( "%.4f  %.4f  %.4f  <= Real Ans\n", m6,v6,s6 )



}





func calcMean( data DataSet, dim int ) (mean float64) {
  sum := 0.
  good := 0
  for i := 0; i < data.NumPoints(); i++ {
    d := data.Point(i).Val(dim)
    if math.IsNaN(d) { continue }
    sum += d
    good++
  }
  mean = sum / float64(good)
  return
}

func calcMeanStdDev( data DataSet, dim int ) (mean,vari,stddev float64) {
  dif, sum, sumsqr := 0.,0.,0.
  good := 0
  for i := 3; i < data.NumPoints()-3; i++ {
    d := data.Point(i).Val(dim)
    if math.IsNaN(d) { continue }
    sum += d
    good++
  }
  mean = sum / float64(good)
  for i := 3; i < data.NumPoints()-3; i++ {
    d := data.Point(i).Val(dim)
    if math.IsNaN(d) { continue }
    dif = (d-mean)
    sumsqr += dif*dif
  }
  vari = sumsqr / float64(good-1)
  stddev = math.Sqrt(vari)
  return
}



