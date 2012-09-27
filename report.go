package symexpr

import (
	"fmt"
	"sort"
	"container/list"
)

var MAX_BBQ_SIZE = 1 << 20

const (
	BBSORT_NULL = iota

	BBSORT_HITS
	BBSORT_ERRS
	BBSORT_SIZE

	BBSORT_SIZEHITS
	BBSORT_SIZEERRS
	BBSORT_HITSSIZE
	BBSORT_ERRSSIZE

	BBSORT_PERR
	BBSORT_PHIT
)

type ExprReport struct {
	expr  Expr
	coeff []float64

	// metrics
	PredErr, RealErr   float64
	PredHits, RealHits int

	// per data set metrics, if multiple data sets used
	PredErrz []float64
	PredHitz []int
	RealErrz []float64
	RealHitz []int

	// ids
	uniqID int // unique ID among all exprs ?
	procID int // ID of the search process
	iterID int // iteration of the search
	unitID int // ID with respect to the search

	// production information
	// p1,p2 int  // parent IDs
	// method int // method that produced this expression
}

type ExprReportArray []*ExprReport

func (p ExprReportArray) Len() int      { return len(p) }
func (p ExprReportArray) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ExprReportArray) Less(i, j int) bool {
	if p[i] == nil {
		return false
	}
	if p[j] == nil {
		return true
	}
	return p[i].expr.AmILess(p[j].expr)
}

func (r *ExprReport) String() string {
	format := "%v\n%d  %f\nuId: %d   pId: %d\niId: %d   tId: %d\n"
	return fmt.Sprintf(format, r.expr, r.score, r.error, r.uniqID, r.procID, r.iterID, r.unitID)
}

func (r *ExprReport) Expr() Expr     { return r.expr }
func (r *ExprReport) SetExpr(e Expr) { r.expr = e }

func (r *ExprReport) Coeff() []float64     { return r.coeff }
func (r *ExprReport) SetCoeff(c []float64) { r.coeff = c }

func (r *ExprReport) Score() int     { return r.score }
func (r *ExprReport) SetScore(s int) { r.score = s }

func (r *ExprReport) Error() float64     { return r.error }
func (r *ExprReport) SetError(e float64) { r.error = e }

func (r *ExprReport) UniqID() int     { return r.uniqID }
func (r *ExprReport) SetUniqID(i int) { r.uniqID = i }

func (r *ExprReport) ProcID() int     { return r.procID }
func (r *ExprReport) SetProcID(i int) { r.procID = i }

func (r *ExprReport) IterID() int     { return r.iterID }
func (r *ExprReport) SetIterID(i int) { r.iterID = i }

func (r *ExprReport) UnitID() int     { return r.unitID }
func (r *ExprReport) SetUnitID(i int) { r.unitID = i }

type ReportQueue struct {
	queue      []*ExprReport
	less       func(i, j *ExprReport) bool
	sortmethod int
}

func NewBBQ() *ReportQueue {
	B := new(ReportQueue)
	B.queue = make([]*ExprReport, 0, MAX_BBQ_SIZE)
	return B
}

func (bb ReportQueue) Len() int { return len(bb.queue) }
func (bb ReportQueue) Less(i, j int) bool {
	return bb.less(bb.queue[i], bb.queue[j])
}
func (bb ReportQueue) Swap(i, j int) {
	bb.queue[i], bb.queue[j] = bb.queue[j], bb.queue[i]
	bb.queue[i].index, bb.queue[j].index = i, j
}
func (bb *ReportQueue) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	// To simplify indexing expressions in these methods, we save a copy of the
	// slice object. We could instead write (*pq)[i].
	if len(bb.queue) == cap(bb.queue) {
		B := make([]*ExprReport, 0, len(bb.queue)*2)
		copy(B[:len(bb.queue)], (bb.queue)[:])
		bb.queue = B
	}
	a := bb.queue
	n := len(a)
	a = a[0 : n+1]
	item := x.(*ExprReport)
	item.index = n
	a[n] = item
	bb.queue = a
}

func (bb *ReportQueue) Pop() interface{} {
	a := bb.queue
	n := len(a)
	item := a[n-1]
	item.index = -1 // for safety
	bb.queue = a[0 : n-1]
	return item
}

func (bb *ReportQueue) Sort() {
	// all of the methods need to be reversed conceptually except for paretos
	// this is because the queue push/pops from the back
	switch bb.sortmethod {
	case BBSORT_HITS:
		bb.less = moreHits
		sort.Sort(bb)
	case BBSORT_ERRS:
		bb.less = moreError
		sort.Sort(bb)
	case BBSORT_SIZE:
		bb.less = moreSize
		sort.Sort(bb)
	case BBSORT_SIZEHITS:
		bb.less = moreSizeHits
		sort.Sort(bb)
	case BBSORT_SIZEERRS:
		bb.less = moreSizeError
		sort.Sort(bb)
	case BBSORT_HITSSIZE:
		bb.less = moreHitsSize
		sort.Sort(bb)
	case BBSORT_ERRSSIZE:
		bb.less = moreErrorSize
		sort.Sort(bb)
	case BBSORT_PHIT:
		bb.ParetoHits()
	case BBSORT_PERR:
		bb.ParetoError()
	default:

	}

}

func lessHits(l, r *ExprReport) bool {
	sc := l.score - r.score
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	return l.cand.AmILess(r.cand)
}
func lessError(l, r *ExprReport) bool {
	sc := l.error - r.error
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	return l.cand.AmILess(r.cand)
}
func lessSize(l, r *ExprReport) bool {
	sz := l.size - r.size
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	return l.cand.AmILess(r.cand)
}
func moreHits(l, r *ExprReport) bool {
	sc := l.score - r.score
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	return l.cand.AmILess(r.cand)
}
func moreError(l, r *ExprReport) bool {
	sc := l.error - r.error
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	return l.cand.AmILess(r.cand)
}
func moreSize(l, r *ExprReport) bool {
	sz := l.size - r.size
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	return l.cand.AmILess(r.cand)
}

func lessSizeHits(l, r *ExprReport) bool {
	sz := l.size - r.size
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	sc := l.score - r.score
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	return l.cand.AmILess(r.cand)
}
func lessSizeError(l, r *ExprReport) bool {
	sz := l.size - r.size
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	sc := l.error - r.error
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	return l.cand.AmILess(r.cand)
}
func lessHitsSize(l, r *ExprReport) bool {
	sc := l.score - r.score
	if sc > 0 {
		return true
	} else if sc < 0 {
		return false
	}
	sz := l.size - r.size
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	return l.cand.AmILess(r.cand)
}
func lessErrorSize(l, r *ExprReport) bool {
	sc := l.error - r.error
	if sc < 0.0 {
		return true
	} else if sc > 0.0 {
		return false
	}
	sz := l.size - r.size
	if sz < 0 {
		return true
	} else if sz > 0 {
		return false
	}
	return l.cand.AmILess(r.cand)
}
func moreSizeHits(l, r *ExprReport) bool {
	sz := l.size - r.size
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	sc := l.score - r.score
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	return l.cand.AmILess(r.cand)
}
func moreSizeError(l, r *ExprReport) bool {
	sz := l.size - r.size
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	sc := l.error - r.error
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	return l.cand.AmILess(r.cand)
}
func moreHitsSize(l, r *ExprReport) bool {
	sc := l.score - r.score
	if sc > 0 {
		return false
	} else if sc < 0 {
		return true
	}
	sz := l.size - r.size
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	return l.cand.AmILess(r.cand)
}
func moreErrorSize(l, r *ExprReport) bool {
	sc := l.error - r.error
	if sc < 0.0 {
		return false
	} else if sc > 0.0 {
		return true
	}
	sz := l.size - r.size
	if sz < 0 {
		return false
	} else if sz > 0 {
		return true
	}
	return l.cand.AmILess(r.cand)
}

func (bb *ReportQueue) ParetoHits() {
	bb.less = lessSizeHits
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.size
		cScore := pb.score
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.score > cScore {
				cScore = pb.score
				if pb.size > cSize {
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					cSize = pb.size
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}

func (bb *ReportQueue) ParetoError() {
	bb.less = lessSizeError
	sort.Sort(bb)

	var pareto list.List
	pareto.Init()
	for i, _ := range bb.queue {
		pareto.PushBack(bb.queue[i])
	}

	over := len(bb.queue) - 1
	for pareto.Len() > 0 && over >= 0 {
		pe := pareto.Front()
		eLast := pe
		pb := pe.Value.(*ExprReport)
		cSize := pb.size
		cError := pb.error
		pe = pe.Next()
		for pe != nil && over >= 0 {
			pb := pe.Value.(*ExprReport)
			if pb.error < cError {
				cError = pb.error
				if pb.size > cSize {
					bb.queue[over] = eLast.Value.(*ExprReport)
					over--
					pareto.Remove(eLast)
					cSize = pb.size
					eLast = pe
				}
			}
			pe = pe.Next()
		}
		if over < 0 {
			break
		}

		bb.queue[over] = eLast.Value.(*ExprReport)
		over--
		pareto.Remove(eLast)
	}
}

func RegressExpr(E SE.Expr, P *Problem) (R *BBfunctor) {

	c := make([]float64, 0)
	c, eqn := E.ConvertToConstants(c)

	var coeff []float64
	if len(c) > 0 {
		coeff = LevmarExpr(P, eqn, c)
	}

	R = new(BBfunctor)
	R.cand = eqn /*.ConvertToConstantFs(coeff)*/
	R.coeff = coeff
	s1, s2, serr := scoreExpr(E, P, coeff)
	R.score = s2
	R.score2 = s1
	R.error = serr
	R.size = E.CalcExprStats(0)

	return R
}

func scoreExpr(e SE.Expr, P *Problem, coeff []float64) (int, int, float64) {
	score := 0
	score2 := 0
	error := 0.0

	for _, PS := range P.test {
		for _, p := range PS.Points() {
			y := p.Depnd(P.searchDim)
			var out float64
			if P.task == "explicit" {
				out = e.Eval(0, p.Indeps(), coeff, PS.SysVals())
			} else if P.task == "diffeq" {
				out = e.Eval(p.Indep(0), p.Indeps()[1:], coeff, PS.SysVals())
			}

			diff := math.Abs(out - y)
			if diff < P.hitRatio {
				score++
			}
			err := math.Abs(diff / y)
			if math.IsNaN(err) || math.IsInf(err, 0) {
				err = diff
			}
			error += err
			if err < P.hitRatio {
				score2++
			}
		}
	}

	eAve := error / (float64(len(P.test)) * float64(P.test[0].NumPoints()))

	return score, score2, eAve
}
