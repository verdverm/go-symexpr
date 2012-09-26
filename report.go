package symexpr

import (
	"fmt"
)

type ExprReport struct {
	expr  Expr
	coeff []float64

	// metrics
	score int
	error float64

	// ids
	uniqID int // unique ID among all exprs
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
