package symexpr

import (
// "fmt"
// "strings"
)

type RangeType int

const (
	Uniform RangeType = iota
	Equal
)

type BenchmarkVar struct {
	Name    string
	Index   int
	Rtype   RangeType
	L, H, S float64 // low,high,step of range
}

type Benchmark struct {
	Name         string
	TrainVars    []BenchmarkVar
	TrainSamples int
	TestVars     []BenchmarkVar
	TestSamples  int

	FuncText string // function as text
	FuncTree Expr   // function as tree
	// InputData data.Pointset
}

func ParseFunc(text string, varNames []string) Expr {
	expr := parse(text, varNames)
	return expr
}

var xU11 = BenchmarkVar{"x", 0, Uniform, -1.0, 1.0, 0.0}
var xU22 = BenchmarkVar{"x", 0, Uniform, -2.0, 2.0, 0.0}
var xU33 = BenchmarkVar{"x", 0, Uniform, -3.0, 3.0, 0.0}
var xU01 = BenchmarkVar{"x", 0, Uniform, 0.0, 1.0, 0.0}
var xU02 = BenchmarkVar{"x", 0, Uniform, 0.0, 2.0, 0.0}
var xU04 = BenchmarkVar{"x", 0, Uniform, 0.0, 4.0, 0.0}
var yU01 = BenchmarkVar{"y", 1, Uniform, 0.0, 1.0, 0.0}

var xU5050 = BenchmarkVar{"x", 0, Uniform, -50.0, 50.0, 0.0}
var yU5050 = BenchmarkVar{"y", 1, Uniform, -50.0, 50.0, 0.0}
var zU5050 = BenchmarkVar{"z", 2, Uniform, -50.0, 50.0, 0.0}
var vU5050 = BenchmarkVar{"v", 3, Uniform, -50.0, 50.0, 0.0}
var wU5050 = BenchmarkVar{"w", 4, Uniform, -50.0, 50.0, 0.0}

var Korns5 = []BenchmarkVar{xU5050, yU5050, zU5050, vU5050, wU5050}
var xyU01 = []BenchmarkVar{xU01, yU01}

var xE11_01 = BenchmarkVar{"x", 0, Equal, -1.0, 1.0, 0.01}
var xE22_01 = BenchmarkVar{"x", 0, Equal, -2.0, 2.0, 0.01}
var xE33_01 = BenchmarkVar{"x", 0, Equal, -3.0, 3.0, 0.01}
var xE11_001 = BenchmarkVar{"x", 0, Equal, -1.0, 1.0, 0.001}
var xE22_001 = BenchmarkVar{"x", 0, Equal, -2.0, 2.0, 0.001}
var xE33_001 = BenchmarkVar{"x", 0, Equal, -3.0, 3.0, 0.001}

var xE55_4 = BenchmarkVar{"x", 0, Equal, -5.0, 5.0, 0.4}
var yE55_4 = BenchmarkVar{"y", 1, Equal, -5.0, 5.0, 0.4}
var xyE55_4 = []BenchmarkVar{xE55_4, yE55_4}

var benchmarks = []Benchmark{
	Benchmark{"Koza_1", []BenchmarkVar{xU11}, 20, nil, 0, "x^4 + x^3 + x^2 + x", nil},
	Benchmark{"Koza_2", []BenchmarkVar{xU11}, 20, nil, 0, "x^5 - 2x^3 + x", nil},
	Benchmark{"Koza_3", []BenchmarkVar{xU11}, 20, nil, 0, "x^6 - 2x^4 + x^2", nil},

	Benchmark{"Nguyen_01", []BenchmarkVar{xU11}, 20, nil, 0, "x^3 + x^2 + x", nil},
	Benchmark{"Nguyen_02", []BenchmarkVar{xU11}, 20, nil, 0, "x^4 + x^3 + x^2 + x", nil},
	Benchmark{"Nguyen_03", []BenchmarkVar{xU11}, 20, nil, 0, "x^5 + x^4 + x^3 + x^2 + x", nil},
	Benchmark{"Nguyen_04", []BenchmarkVar{xU11}, 20, nil, 0, "x^6 + x^5 + x^4 + x^3 + x^2 + x", nil},
	Benchmark{"Nguyen_05", []BenchmarkVar{xU11}, 20, nil, 0, "sin(x^2)*cos(x) - 1", nil},
	Benchmark{"Nguyen_06", []BenchmarkVar{xU11}, 20, nil, 0, "sin(x) + sin(x + x^2)", nil},
	Benchmark{"Nguyen_07", []BenchmarkVar{xU02}, 20, nil, 0, "ln(x+1) + ln(x^2 + 1)", nil},
	Benchmark{"Nguyen_08", []BenchmarkVar{xU04}, 20, nil, 0, "sqrt(x)", nil},
	Benchmark{"Nguyen_09", xyU01, 20, nil, 0, "sin(x) + sin(y^2)", nil},
	Benchmark{"Nguyen_10", xyU01, 20, nil, 0, "2*sin(x)*cos(y)", nil},
	Benchmark{"Nguyen_11", xyU01, 20, nil, 0, "x^y", nil},
	Benchmark{"Nguyen_12", xyU01, 20, nil, 0, "x^4 - x^3 + 0.5*y^2 - y", nil},

	Benchmark{"Pagie_1", xyE55_4, 0, nil, 0, "1 / (1 + x^-4) + 1 / (1 + y^-4)", nil},

	// 5 inputs: x,y,z,v,w
	Benchmark{"Korns_01", Korns5, 10000, Korns5, 10000, "1.57 + 24.3*v", nil},
	Benchmark{"Korns_02", Korns5, 10000, Korns5, 10000, "0.23 + 14.2*(v+y)/3w", nil},
	Benchmark{"Korns_03", Korns5, 10000, Korns5, 10000, "-5.41 + 4.9*(v-x+y/w)/3w", nil},
	Benchmark{"Korns_04", Korns5, 10000, Korns5, 10000, "-2.3 + 0.13sin(z)", nil},
	Benchmark{"Korns_05", Korns5, 10000, Korns5, 10000, "3 + 2.13*ln(w)", nil},
	Benchmark{"Korns_06", Korns5, 10000, Korns5, 10000, "1.3 + 0.13*sqrt(x)", nil},
	Benchmark{"Korns_07", Korns5, 10000, Korns5, 10000, "213.80940889*(1 - e^(-0.54723748542x))", nil},
	Benchmark{"Korns_08", Korns5, 10000, Korns5, 10000, "6.87 + 11*sqrt(7.23*x*v*w)", nil},
	Benchmark{"Korns_09", Korns5, 10000, Korns5, 10000, "sqrt(x)/ln(y) * e^z / v^2", nil},
	Benchmark{"Korns_10", Korns5, 10000, Korns5, 10000, "0.81 + 24.3*(2y+3*z^2)/(4*(v)^3+5*(w)^4)", nil},
	Benchmark{"Korns_11", Korns5, 10000, Korns5, 10000, "6.87 + 11*cos(7.23*x^3)", nil},
	Benchmark{"Korns_12", Korns5, 10000, Korns5, 10000, "2 - 2.1*cos(9.8*x)*sin(1.3*w)", nil},
	Benchmark{"Korns_13", Korns5, 10000, Korns5, 10000, "32 - 3*(tan(x)*tan(z))/(tan(y)*tan(v))", nil},
	Benchmark{"Korns_14", Korns5, 10000, Korns5, 10000, "22 - 4.2*(cos(x)-tan(y))*(tanh(z)/sin(v))", nil},
	Benchmark{"Korns_15", Korns5, 10000, Korns5, 10000, "12 - 6*(tan(x)/e^y)(ln(z)-tan(v))", nil},

	Benchmark{"Keijzer_01", []BenchmarkVar{xE11_01}, 0, []BenchmarkVar{xE11_001}, 0, "0.3*x*sin(2*PI*x)", nil},
	Benchmark{"Keijzer_02", []BenchmarkVar{xE22_01}, 0, []BenchmarkVar{xE22_001}, 0, "0.3*x*sin(2*PI*x)", nil},
	Benchmark{"Keijzer_03", []BenchmarkVar{xE33_01}, 0, []BenchmarkVar{xE33_001}, 0, "0.3*x*sin(2*PI*x)", nil},
	// Benchmark{"Keijzer_04", "x", "E[0,10,0.05]", "x^3*e^-x*cos(x)*sin(x)*((sin(x))^2*cos(x) - 1)", nil},
	// Benchmark{"Keijzer_05", "x,y,z", "x,z: U[-1,1,1000] y: U[1,2,1000]", "(30*x*z) / ((x-10)*y^2)", nil},
	// // Benchmark{"Keijzer_06", "x", "E[1,50,1]", "\\SUM_i^x (1/i)", nil},
	// Benchmark{"Keijzer_07", "x", "E[1,100,1]", "ln(x)", nil},
	// Benchmark{"Keijzer_08", "x", "E[0,100,1]", "sqrt(x)", nil},
	// // arcsinh(x) == ln(x+sqrt(x^2+1))
	// Benchmark{"Keijzer_09", "x", "E[0,100,1]", "ln(x+sqrt(x^2+1))", nil},
	// Benchmark{"Keijzer_10", "x,y", "U[0,1,100]", "x^y", nil},
	// // xy ? sin((x-1)(y-1))
	// // Benchmark{"Keijzer_11", "x,y", "U[-3,3,20]", "xy ? sin((x-1)*(y-1))", nil},
	// Benchmark{"Keijzer_12", "x,y", "U[-3,3,20]", "x^4 - x^3 + 0.5*y^2 - y ", nil},
	// Benchmark{"Keijzer_13", "x,y", "U[-3,3,20]", "6*sin(x)*cos(y)", nil},
	// Benchmark{"Keijzer_14", "x,y", "U[-3,3,20]", "8 / (2 + x^2 + y^2) ", nil},
	// Benchmark{"Keijzer_15", "x,y", "U[-3,3,20]", "0.2*x^3 + 0.5*y^2 - y - x ", nil},

	// // (e^{-(x-1)^2}) / (1.2 + (y-2.5)^2)
	// Benchmark{"Vladislavleva_1", "x,y", "?U[0.3,4,10]?", "(e^{-(x-1)^2}) / (1.2 + (y-2.5)^2)", nil},
	// Benchmark{"Vladislavleva_2", "x", "E[0.05,10,0.1]", "e^-x* x^3 * (cos(x)*sin(x)) * (cos(x)*(sin(x))^2-1)", nil},
	// Benchmark{"Vladislavleva_3", "x,y", "x: E[0.05,10,0.1]  y: E[0.05,10.05,2]", "e^-y * x^3 * (cos(x)*sin(x)) * (cos(x)*(sin(x))^2-1)", nil},
	// // Benchmark{"Vladislavleva_4", "x_i", " U[0.05, 6.05, 1024]", "10 / (5 + \\SUM_1^5 (x_i - 3)^2)", nil},
	// Benchmark{"Vladislavleva_5", "x,y,z", "x,z: U[0.05,2,300]  y: U[1,2,300]", "(30*(x-1)*(z-1)) / (y^2*(x-10))  ", nil},
	// Benchmark{"Vladislavleva_6", "x,y", "U[0.1,5.9,30]", "6*sin(x)*cos(y)", nil},
	// Benchmark{"Vladislavleva_7", "x,y", "U[0.05,6.05,300]", "(x-3)*(y-3) + 2sin((x-4)*(y-4))", nil},
	// Benchmark{"Vladislavleva_8", "x,y", "U[0.05,6.05,50]", "((x-3)^4 + (y-3)^3 - (y-3)) / ((y-2)^4 + 10)", nil},
}
