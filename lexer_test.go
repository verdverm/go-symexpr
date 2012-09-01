package symexpr

import (
	"testing"

	"fmt"
)

/* still getting errors on:
 * 'strings of identifiers' :: 23,31,32,33,35
 * 'unimplemented functions(tan,tanh)' :: 28,29,30
 * '\\SUM...' :: 36,
 * 'precedence erro ' :: 47,50
 * 'other' :: 24,25,41
 */

func Test_Parser(TEST *testing.T) {
	fmt.Printf("Testing: Parser\n\n")

	s, e := 0, len(benchmarks)
	for i := s; i < e; i++ {

		// switch i {
		// case 47:
		// if i < 21 {
		fmt.Printf("Benchmark: %d\n", i)
		b := benchmarks[i]
		fmt.Printf("Input:     %s\n", b.FuncText)

		varNames := make([]string, 0)
		for _, v := range b.TrainVars {
			fmt.Printf("  %v\n", v)
			varNames = append(varNames, v.Name)
		}

		expr := ParseFunc(b.FuncText, varNames)

		// fmt.Printf("Input:     %s\n", b.FuncText)
		fmt.Printf("Result:    %v\n\n\n", expr)
		// }
	}
}

/* latest errors: 
tan,tanh,pi

SUM
41
*/
