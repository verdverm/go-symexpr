package symexpr

import (
	"testing"

	"fmt"
)

func Test_Parser( TEST *testing.T) {
	fmt.Printf( "Testing: Parser\n\n")

	s,e := 0,len(benchmarks)
	for i := s; i<e; i++ {
		fmt.Printf( "Benchmark: %d\n",i)
		b := benchmarks[i]
		fmt.Printf( "Input:     %s\n",b.FuncText)

		expr := ParseFunc(b.FuncText,b.VarsText)

		fmt.Printf( "Result:    %v\n\n\n",expr)

	}
}

