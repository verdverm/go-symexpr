package symexpr

import (
	// "fmt"
	"strings"
)

type Benchmark struct {
	Name string
	VarsText  string	   // variable names 
	InputText string   // range of input data

	FuncText string    // function as text
	FuncTree Expr   // function as tree
	// InputData data.Pointset
}

func ParseFunc( text,vars string ) Expr {
	// fmt.Printf( "Parsing %s\n", text)
	varary := strings.Split(vars,",")
	expr := parse(text,varary)
	// fmt.Printf( "Result: %s\n\n",expr)
	return expr
}

var	benchmarks = []Benchmark{
	Benchmark{"Koza_1","x", "U[-1,1,20]", "x^4 + x^3 + x^2 + x",nil},
	Benchmark{"Koza_2","x", "U[-1,1,20]", "x^5 - 2x^3 + x",nil},
	Benchmark{"Koza_3","x", "U[-1,1,20]", "x^6 - 2x^4 + x^2",nil},

	Benchmark{"Nguyen_1","x",    "U[-1,1,20]", "x^3 + x^2 + x",nil},
	Benchmark{"Nguyen_2","x",    "U[-1,1,20]", "x^4 + x^3 + x^2 + x",nil},
	Benchmark{"Nguyen_3","x",    "U[-1,1,20]", "x^5 + x^4 + x^3 + x^2 + x",nil},
	Benchmark{"Nguyen_4","x",    "U[-1,1,20]", "x^6 + x^5 + x^4 + x^3 + x^2 + x",nil},
	Benchmark{"Nguyen_5","x",    "U[-1,1,20]", "sin(x^2)cos(x) - 1",nil},
	Benchmark{"Nguyen_6","x",    "U[-1,1,20]", "sin(x) + sin(x + x^2)",nil},
	Benchmark{"Nguyen_7","x",    "U[0,2,20]",  "ln(x+1) + ln(x^2 + 1)",nil},
	Benchmark{"Nguyen_8","x",    "U[0,4,20]",  "sqrt(x)",nil},
	Benchmark{"Nguyen_9","x,y",  "U[0,1,20]",  "sin(x) + sin(y^2)",nil},
	Benchmark{"Nguyen_10","x,y", "U[0,1,20]",  "2sin(x)cos(y)",nil},
	Benchmark{"Nguyen_11","x,y", "U[0,1,20]",  "x^y",nil},
	Benchmark{"Nguyen_12","x,y", "U[0,1,20]",  "x^4 - x^3 + 0.5y^2 - y",nil},

	Benchmark{"Pagie_1","x,y",   "E[-5,5,0.4]","1 / (1 + x^-4) + 1 / (1 + y^-4)",nil},

	// 5 inputs: x,y,z,v,w
	Benchmark{"Korns_1", "x,y,z,v,w", "U[-50,50,10000]", "1.57 + 24.3v",nil},
	Benchmark{"Korns_2", "x,y,z,v,w", "U[-50,50,10000]", "0.23 + 14.2(v+y)/3w",nil},
	Benchmark{"Korns_3", "x,y,z,v,w", "U[-50,50,10000]", "-5.41 + 4.9(v-x+y/w)/3w",nil},
	Benchmark{"Korns_4", "x,y,z,v,w", "U[-50,50,10000]", "-2.3 + 0.13sin(z)",nil},
	Benchmark{"Korns_5", "x,y,z,v,w", "U[-50,50,10000]", "3 + 2.13 ln(w)",nil},
	Benchmark{"Korns_6", "x,y,z,v,w", "U[-50,50,10000]", "1.3 + 0.13 sqrt(x)",nil},
	Benchmark{"Korns_7", "x,y,z,v,w", "U[-50,50,10000]", "213.80940889(1 - e^(-0.54723748542x))",nil},
	Benchmark{"Korns_8", "x,y,z,v,w", "U[-50,50,10000]", "6.87 + 11 sqrt(7.23xvw)",nil},
	Benchmark{"Korns_9", "x,y,z,v,w", "U[-50,50,10000]", "sqrt(x)/ln(y) * e^z / v^2",nil},
	Benchmark{"Korns_10","x,y,z,v,w", "U[-50,50,10000]", "0.81 + 24.3 (2y+3z^2)/(4(v)^3+5(w)^4)",nil},
	Benchmark{"Korns_11","x,y,z,v,w", "U[-50,50,10000]", "6.87 + 11 cos(7.23 x^3)",nil},
	Benchmark{"Korns_12","x,y,z,v,w", "U[-50,50,10000]", "2 - 2.1cos(9.8x)sin(1.3w)",nil},
	Benchmark{"Korns_13","x,y,z,v,w", "U[-50,50,10000]", "32 - 3 (tan(x)tan(z))/(tan(y)tan(v))",nil},
	Benchmark{"Korns_14","x,y,z,v,w", "U[-50,50,10000]", "22 - 4.2(cos(x)-tan(y))(tanh(z)/sin(v))",nil},
	Benchmark{"Korns_15","x,y,z,v,w", "U[-50,50,10000]", "12 - 6(tan(x)/e^y)(ln(z)-tan(v))",nil},

	Benchmark{"Keijzer_1","x",    "E[-1,1,0.1]",  "0.3 x sin(2PIx)",nil},
	Benchmark{"Keijzer_2","x",    "E[-2,2,0.1]",  "0.3 x sin(2PIx)",nil},
	Benchmark{"Keijzer_3","x",    "E[-3,3,0.1]",  "0.3 x sin(2PIx)",nil},
	Benchmark{"Keijzer_4","x",    "E[0,10,0.05]", "x^3 e^-x cos(x) sin(x) (sin^2(x)cos(x) - 1)",nil},
	Benchmark{"Keijzer_5","x,y,z","x,z: U[-1,1,1000] y: U[1,2,1000]",   "30xz / (x-10)y^2",nil},
	Benchmark{"Keijzer_6","x",    "E[1,50,1]",    "\\SUM_i^x (1/i)",nil},
	Benchmark{"Keijzer_7","x",    "E[1,100,1]",   "ln(x)",nil},
	Benchmark{"Keijzer_8","x",    "E[0,100,1]",   "sqrt(x)",nil},
	// arcsinh(x) == ln(x+sqrt(x^2+1))
	Benchmark{"Keijzer_9","x",    "E[0,100,1]",   "ln(x+sqrt(x^2+1))",nil},
	Benchmark{"Keijzer_10","x,y", "U[0,1,100]",   "x^y",nil},
	// xy ? sin((x-1)(y-1))
	Benchmark{"Keijzer_11","x,y", "U[-3,3,20]",   "xy ? sin((x-1)(y-1))",nil},
	Benchmark{"Keijzer_12","x,y", "U[-3,3,20]",   "x^4 - x^3 + 0.5y^2 - y ",nil},
	Benchmark{"Keijzer_13","x,y", "U[-3,3,20]",   "6 sin(x) cos(y)",nil},
	Benchmark{"Keijzer_14","x,y", "U[-3,3,20]",   "8 / (2 + x^2 + y^2) ",nil},
	Benchmark{"Keijzer_15","x,y", "U[-3,3,20]",   "0.2x^3 + 0.5y^2 - y - x ",nil},

	// (e^{-(x-1)^2}) / (1.2 + (y-2.5)^2)
	Benchmark{"Vladislavleva_1","x,y", "?U[0.3,4,10]?",                        "(e^{-(x-1)^2}) / (1.2 + (y-2.5)^2)",nil},
	Benchmark{"Vladislavleva_2","x",   "E[0.05,10,0.1]",                       "e^-x x^3 (cos(x)sin(x)) (cos(x)sin^2(x)-1)",nil},
	Benchmark{"Vladislavleva_3","x,y", "x: E[0.05,10,0.1]  y: E[0.05,10.05,2]","e^-x x^3 (cos(x)sin(x)) (cos(x)sin^2(x)-1)",nil},
	Benchmark{"Vladislavleva_4","x_i", " U[0.05, 6.05, 1024]",                 "10 / (5 + \\SUM_1^5 (x_i - 3)^2)",nil},
	Benchmark{"Vladislavleva_5","x,y", "x,z: U[0.05,2,300]  y: U[1,2,300]",    "30(x-1)(z-1) / y^2(x-10)  ",nil},
	Benchmark{"Vladislavleva_6","x,y", "U[0.1,5.9,30]",                        "6sin(x)cos(y)",nil},
	Benchmark{"Vladislavleva_7","x,y", "U[0.05,6.05,300]",                     "(x-3)(y-3) + 2sin((x-4)(y-4))",nil},
	Benchmark{"Vladislavleva_8","x,y", "U[0.05,6.05,50]",                      "((x-3)^4 + (y-3)^3) - (y-3)) / ((y-2)^4 + 10)",nil},
}
