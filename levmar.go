package symexpr

/*
#cgo LDFLAGS: -L/usr/lib  -L/home/tony/src/levmar-2.6 -llevmar -llapack -lblas -lf2c  -lm

void func(double *p, double *x, int m, int n, void *data);
void jacfunc(double *p, double *x, int m, int n, void *data);
void levmar(double* ygiven, double* p, const int n, const int m, void* e );

*/
import "C"

import (
	"unsafe"
	"reflect"
)

type callback_data struct {
	Train []*PointSet
	Test  []*PointSet
	E     Expr
	J     []Expr
	Coeff []float64
	Task  string // "explicit" or "diffeq"
}

func LevmarExpr(e Expr, searchDim int, task string, guess []float64, train, test []*PointSet) []float64 {

	ps := train[0].NumPoints()
	PS := len(train) * ps

	c := make([]float64, len(guess))
	var cd callback_data
	cd.Train = train
	cd.Test = test
	cd.E = e
	cd.Coeff = c
	cd.Task = task
	cd.J = make([]Expr, len(guess))
	for i, _ := range cd.J {
		cd.J[i] = e.DerivConst(i) /*.Simplify(SimpRules{true,true})*/
	}

	// c/levmar inputs
	coeff := make([]C.double, len(guess))
	for i, g := range guess {
		coeff[i] = C.double(g)
	}

	y := make([]C.double, PS)
	for i1, T := range train {
		for i2, p := range T.Points() {
			i := i1*ps + i2
			y[i] = C.double(p.Depnd(searchDim))
		}
	}
	ya := (*C.double)(unsafe.Pointer(&y[0]))
	ca := (*C.double)(unsafe.Pointer(&coeff[0]))
	ni := C.int(PS)
	mi := C.int(len(c))

	C.levmar(ya, ca, ni, mi, unsafe.Pointer(&cd))

	for i, _ := range coeff {
		c[i] = float64(coeff[i])
	}
	return c
}

//export callback_func
func callback_func(p, x *C.double, e unsafe.Pointer) {

	cd := *(*callback_data)(e)
	eqn := cd.E
	coeff := cd.Coeff

	M1 := len(cd.Train)
	M2 := cd.Train[0].NumPoints()
	M3 := len(coeff)
	M23 := M2 * M3
	MA := M1 * M23

	var p_go []C.double
	p_head := (*reflect.SliceHeader)((unsafe.Pointer(&p_go)))
	p_head.Cap = M3
	p_head.Len = M3
	p_head.Data = uintptr(unsafe.Pointer(p))
	for i, _ := range p_go {
		coeff[i] = float64(p_go[i])
	}

	var x_go []C.double
	x_head := (*reflect.SliceHeader)((unsafe.Pointer(&x_go)))
	x_head.Cap = MA
	x_head.Len = MA
	x_head.Data = uintptr(unsafe.Pointer(x))

	N := len(cd.Train)
	var out float64
	for i1, PS := range cd.Train {
		for i2, pnt := range PS.Points() {
			i := i1*N + i2
			if cd.Task == "explicit" {
				out = eqn.Eval(0, pnt.Indeps(), coeff, PS.SysVals())
			} else if cd.Task == "diffeq" {
				out = eqn.Eval(pnt.Indep(0), pnt.Indeps()[1:], coeff, PS.SysVals())
			}

			x_go[i] = C.double(out)

		}
	}
}

//export callback_jacfunc
func callback_jacfunc(p, x *C.double, e unsafe.Pointer) {

	cd := *(*callback_data)(e)
	coeff := cd.Coeff

	M1 := len(cd.Train)
	M2 := cd.Train[0].NumPoints()
	M3 := len(coeff)
	M23 := M2 * M3
	MA := M1 * M23

	var p_go []C.double
	p_head := (*reflect.SliceHeader)((unsafe.Pointer(&p_go)))
	p_head.Cap = M3
	p_head.Len = M3
	p_head.Data = uintptr(unsafe.Pointer(p))
	for i, _ := range p_go {
		coeff[i] = float64(p_go[i])
	}

	var x_go []C.double
	x_head := (*reflect.SliceHeader)((unsafe.Pointer(&x_go)))
	x_head.Cap = MA
	x_head.Len = MA
	x_head.Data = uintptr(unsafe.Pointer(x))

	var out float64
	for i1, PS := range cd.Train {
		for i2, pnt := range PS.Points() {
			i := i1*M23 + i2*M3
			for ji, eqn := range cd.J {
				if cd.Task == "explicit" {
					out = eqn.Eval(0, pnt.Indeps(), coeff, PS.SysVals())
				} else if cd.Task == "diffeq" {
					out = eqn.Eval(pnt.Indep(0), pnt.Indeps()[1:], coeff, PS.SysVals())
				}

				x_go[i+ji] = C.double(out)
			}
		}
	}

}
