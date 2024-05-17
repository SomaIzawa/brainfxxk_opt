package optimizer

import (
	"fmt"

	"github.com/rosylilly/brainfxxk/ast"
)

type Optimizer struct {
}

func NewOptimizer() *Optimizer {
	return &Optimizer{}
}

func (o *Optimizer) Optimize(p *ast.Program) (*ast.Program, error) {
	exprs, err := o.optimizeExpressions(p.Expressions)
	if err != nil {
		return nil, err
	}

	prog := &ast.Program{
		Expressions: exprs,
	}

	return prog, nil
}

func (o *Optimizer) optimizeExpressions(exprs []ast.Expression) ([]ast.Expression, error) {
	optimized := []ast.Expression{}
	for _, expr := range exprs {
		optExpr, err := o.optimizeExpression(expr)
		if err != nil {
			return nil, err
		}

		switch optExpr.(type) {
		case *ast.PointerIncrementExpression:
			if len(optimized) > 0 {
				// 一つ前が"MultiplePointerIncrementExpression"かどうか
				if last, ok := optimized[len(optimized)-1].(*ast.MultiplePointerIncrementExpression); ok {
					last.Count += 1
					last.Expressions = append(last.Expressions, optExpr)
					continue
				}
			}

			optExpr = &ast.MultiplePointerIncrementExpression{
				Count:       1,
				Expressions: []ast.Expression{optExpr},
			}
		case *ast.PointerDecrementExpression:
			if len(optimized) > 0 {
				if last, ok := optimized[len(optimized)-1].(*ast.MultiplePointerDecrementExpression); ok {
					last.Count += 1
					last.Expressions = append(last.Expressions, optExpr)
					continue
				}
			}

			optExpr = &ast.MultiplePointerDecrementExpression{
				Count:       1,
				Expressions: []ast.Expression{optExpr},
			}
		case *ast.ValueIncrementExpression: {
			if len(optimized) > 0 {
				if last, ok := optimized[len(optimized)-1].(*ast.MultipleValueIncrementExpression); ok {
					last.Count += 1
					last.Expressions = append(last.Expressions, optExpr)
					continue
				}
			}

			optExpr = &ast.MultipleValueIncrementExpression{
				Count: 1,
				Expressions: []ast.Expression{optExpr},
			}
		}
		case *ast.ValueDecrementExpression: {
			if len(optimized) > 0 {
				if last, ok := optimized[len(optimized)-1].(*ast.MultipleValueDecrementExpression); ok {
					last.Count += 1
					last.Expressions = append(last.Expressions, optExpr)
					continue
				}
			}

			optExpr = &ast.MultipleValueDecrementExpression{
				Count: 1,
				Expressions: []ast.Expression{optExpr},
			}
		}
	case *ast.WhileExpression: {
		we := optExpr.(*ast.WhileExpression)

		if we.String() == "[-]" {
			optExpr = &ast.LoadZeroExpression{
				Expressions: we.Body,
			}
		} else {
			ex, _ := o.optimizeExpressions(we.Body)
			if a, b, c := checkConditions(ex); a && b && c {
				copyOpt := &ast.CopyExpression{}
				phase := 0
				for i, exi := range ex {
					switch exi.(type) {
					case *ast.MultiplePointerDecrementExpression:
						phase -= exi.(*ast.MultiplePointerDecrementExpression).Count
					case *ast.MultiplePointerIncrementExpression:
						phase += exi.(*ast.MultiplePointerIncrementExpression).Count
					case *ast.MultipleValueDecrementExpression:
						if i != 0 || i != len(ex) {
							copyOpt.Copys = append(copyOpt.Copys, ast.Copy{
								Phase: phase,
								Count: 0 - exi.(*ast.MultipleValueDecrementExpression).Count,
							})
						}
					case *ast.MultipleValueIncrementExpression:
						copyOpt.Copys = append(copyOpt.Copys, ast.Copy{
							Phase: phase,
							Count: exi.(*ast.MultipleValueIncrementExpression).Count,
						})
					default:
						fmt.Println("error!!")
					}
				}
				optExpr = copyOpt
			} else {
				optExpr = &ast.WhileExpression{
					Body: ex,
				}
			}
		}
	}
	}

		optimized = append(optimized, optExpr)
	}

	return optimized, nil
}

func (o *Optimizer) optimizeExpression(expr ast.Expression) (ast.Expression, error) {
	return expr, nil
}

func checkConditions(exprs []ast.Expression) (bool, bool, bool) {
	// 1. 先頭か末尾に.MultipleValueDecrementExpressionがある
	firstOrLastMultipleValueDecrementExpression := false
	if len(exprs) > 0 {
		if _, ok := exprs[0].(*ast.MultipleValueDecrementExpression); ok {
			firstOrLastMultipleValueDecrementExpression = true
		}
		if _, ok := exprs[len(exprs)-1].(*ast.MultipleValueDecrementExpression); ok {
			firstOrLastMultipleValueDecrementExpression = true
		}
	}

	// 2. +,-,>,<のみで配列は構成されている
	validTarget := true
	for _, expr := range exprs {
		switch expr.(type) {
		case *ast.MultiplePointerDecrementExpression, *ast.MultiplePointerIncrementExpression, *ast.MultipleValueDecrementExpression, *ast.MultipleValueIncrementExpression:
			// Do nothing, valid animal
		default:
			validTarget = false
			break
		}
		if !validTarget {
			break
		}
	}

	// 3. DogまたはCatにはAgeというプロパティが存在する。このプロパティの合計が全体で0になる
	pSum := 0
	for _, expr := range exprs {
		switch e := expr.(type) {
		case *ast.MultiplePointerIncrementExpression:
			pSum += e.Count
		case *ast.MultiplePointerDecrementExpression:
			pSum -= e.Count
		}
	}
	pSumZero := pSum == 0

	return firstOrLastMultipleValueDecrementExpression, validTarget, pSumZero
}