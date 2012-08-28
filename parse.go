package symexpr

import (
	"fmt"
	"strconv"
	"strings"
)



var itemPrec = map[itemType]int{
	itemAdd:		  3,
	itemNeg:		  4,
	itemMul:		  5,
	itemDiv:		  5,
	itemKeyword:	  6,
	itemCarrot:		  7,
}



func parse(input string, variables []string) Expr {
	L := lex("test", input, variables)
	go L.run()

	// var rules = DefaultRules()

	// fmt.Printf( "Input:     %s\n",input)
	expr := parseExpr("",L,0)
	// fmt.Printf( "Result:    %v\n\n",expr)
	// simp := expr.Simplify(rules)
	// fmt.Printf( "Output:    %v\n\n",simp)

	return expr
}

func parseExpr( prefix string, L *lexer, p int) (Expr) {
	// fmt.Printf( "%sin-E(%d)\n", prefix,p)
	e := parsePiece(prefix+"  ",L)
	// fmt.Printf( "%se: %v\n",prefix,e)
	for i:= 0; ; i++{
		var next item
		if L.nextToken.typ == itemNil {
			next = <- L.items
		} else {
			next,L.nextToken.typ = L.nextToken,itemNil
		}
		// fmt.Printf( "%snext(%d-%d): %v\n",prefix,p,i,next)
		typ := next.typ
		
		if isBinary(typ) && itemPrec[typ] >= p {
			q := itemPrec[typ]
			e2 := parseExpr(prefix+"  ",L,q)
			switch typ {
			case itemAdd:
				add := NewAdd()
				add.Insert(e)
				add.Insert(e2)
				e = add
			case itemNeg:
				add := NewAdd()
				add.Insert(e)
				add.Insert(NewNeg(e2))
				e = add
			case itemCarrot:
				pow := NewPowE(e,e2)
				e = pow
			}
		} else if typ == itemIdentifier {
			L.nextToken = next
			e2 := parseExpr(prefix+"  ",L,itemPrec[itemMul])
			mul := NewMul()
			mul.Insert(e)
			mul.Insert(e2)
			e = mul
		} else if typ > itemKeyword && p < itemPrec[itemKeyword]  {
			L.nextToken = next
		 	e2 := parseExpr(prefix+"  ",L,itemPrec[itemMul])
			mul := NewMul()
			mul.Insert(e)
			mul.Insert(e2)
			e = mul
			// fmt.Printf( "e: %v\ne2: %v\n",e,e2)
		} else if typ == itemLParen || typ == itemLBrack {
			// consumed '(' or '{'
			e2 := parseExpr(prefix+"  ",L,0)
			// now consume ')','}'
			var next2 item
			if L.nextToken.typ == itemNil {
			next2 = <- L.items
			} else {
				next2,L.nextToken.typ = L.nextToken,itemNil
			}
			typ2 := next2.typ
			if (typ == itemLParen && typ2 != itemRParen) || (typ == itemLBrack && typ2 != itemRBrack) {
				fmt.Printf( "error: expected rhs of %v\n", next2)
			}
			mul := NewMul()
			mul.Insert(e)
			mul.Insert(e2)
			e = mul
		} else {
			L.nextToken = next
			break
		}
	}

	// fmt.Printf( "%sout-E(%d): %v\n", prefix,p,e)
	return e
}

func parsePiece( prefix string, L *lexer) (Expr) {
	// fmt.Printf( "%sin-P()\n",prefix)
	var next item
	if L.nextToken.typ == itemNil {
		next = <- L.items
	} else {
		next,L.nextToken.typ = L.nextToken,itemNil
	}
	typ := next.typ
	var e Expr

	if isUnary(typ) {
		e1 := parseExpr(prefix+"  ",L,0)
		e = NewNeg(e1)
	} else if typ == itemLParen || typ == itemLBrack {
		// consumed '(' or '{'
		e = parseExpr(prefix+"  ",L,0)
		// now consume ')','}'
		var next2 item
		if L.nextToken.typ == itemNil {
		next2 = <- L.items
		} else {
			next2,L.nextToken.typ = L.nextToken,itemNil
		}
		typ2 := next2.typ
		if (typ == itemLParen && typ2 != itemRParen) || (typ == itemLBrack && typ2 != itemRBrack) {
			fmt.Printf( "error: expected rhs of %v\n", next2)
		}
		return e
	} else if  typ == itemIdentifier { // leaf
		if next.val[0] == 'c' { 
			// coefficient
			ipos := strings.Index(next.val, "_")+1
			index,err := strconv.ParseInt(next.val[ipos:],0,64)
			if err != nil {
				fmt.Printf( "Error (%v) parsing index in '%s'\n", err, next.val)
			}
			e = NewConstant(int(index))
		} else if next.val[0] == 'X' {
			// variable
			ipos := strings.Index(next.val, "_")+1
			index,err := strconv.ParseInt(next.val[ipos:],0,64)
			if err != nil {
				fmt.Printf( "Error (%v) parsing index in '%s'\n", err, next.val)
			}
			return NewVar(int(index))
		} else {
			// is it a named variable
			for p,v := range L.vars {
				if next.val == v {
					e = NewVar(p)
					break
				}
			}

			if e == nil {
				fmt.Printf("UNIMPLEMENTED IDENTIFIER:  %v\n",next)
			}
		}
	} else if typ == itemNumber {
		flt,err := strconv.ParseFloat(next.val,64)
		if err != nil {
			fmt.Printf( "Error (%v) parsing number in '%s'\n", err, next.val)
		}
		e = NewConstantF(flt)
	
	// } elsi if { // func name...
	} else if typ > itemKeyword {
		// is it a named function
		switch typ {
		case itemSin:
			e2 := parseExpr(prefix+"  ",L,itemPrec[itemKeyword])
			e = NewSin(e2)
		case itemCos:
			e2 := parseExpr(prefix+"  ",L,itemPrec[itemKeyword])
			e = NewCos(e2)
		}

		if e == nil {
			fmt.Printf("Unimplemented Function:  %v\n",next)
		}



	} else { // error
		L.nextToken = next
		e = nil
	}
	// } else if { is it a user defined function

	// fmt.Printf( "%sout-P(): %v\n",prefix,e)
	return e
}
