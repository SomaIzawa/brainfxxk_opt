package interpreter

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/rosylilly/brainfxxk/ast"
	"github.com/rosylilly/brainfxxk/optimizer"
	"github.com/rosylilly/brainfxxk/parser"
)

var (
	ErrInputFinished  = fmt.Errorf("input finished")
	ErrMemoryOverflow = fmt.Errorf("memory overflow")
)

type Interpreter struct {
	Program *ast.Program
	Config  *Config
	Memory  []byte
	Pointer int
	Step Step
}

type Step struct {
	PointerIncrementExpression int
	PointerDecrementExpression int
	ValueIncrementExpression int
	ValueDecrementExpression int
	OutputExpression int
	InputExpression int
	WhileExpression int
	MultiplePointerIncrementExpression int
	MultiplePointerDecrementExpression int
	MultipleValueIncrementExpression int
	MultipleValueDecrementExpression int
	LoadZeroExpression int
	CopyExpression int
}

func (s *Step) Sum() int {
	return s.PointerIncrementExpression + s.PointerDecrementExpression + s.ValueIncrementExpression + s.ValueDecrementExpression + s.OutputExpression + s.InputExpression + s.WhileExpression + s.MultiplePointerIncrementExpression + s.MultiplePointerDecrementExpression + s.MultipleValueIncrementExpression + s.MultipleValueDecrementExpression + s.LoadZeroExpression + s.CopyExpression
}

func (s *Step) ShowStep() {
	fmt.Println("======================== Step ========================")
	fmt.Printf("PointerIncrementExpression:         %v\n", s.PointerIncrementExpression)
	fmt.Printf("PointerDecrementExpression:         %v\n", s.PointerDecrementExpression)
	fmt.Printf("ValueIncrementExpression:           %v\n", s.ValueIncrementExpression)
	fmt.Printf("ValueDecrementExpression:           %v\n", s.ValueDecrementExpression)
	fmt.Printf("OutputExpression:                   %v\n", s.OutputExpression)
	fmt.Printf("InputExpression:                    %v\n", s.InputExpression)
	fmt.Printf("WhileExpression:                    %v\n", s.WhileExpression)
	fmt.Printf("MultiplePointerIncrementExpression: %v\n", s.MultiplePointerIncrementExpression)
	fmt.Printf("MultiplePointerDecrementExpression: %v\n", s.MultiplePointerDecrementExpression)
	fmt.Printf("MultipleValueIncrementExpression:   %v\n", s.MultipleValueIncrementExpression)
	fmt.Printf("MultipleValueDecrementExpression:   %v\n", s.MultipleValueDecrementExpression)
	fmt.Printf("LoadZeroExpression:                 %v\n", s.LoadZeroExpression)
	fmt.Printf("CopyExpression:                     %v\n", s.CopyExpression)
	fmt.Println("------------------------------------------------------")
	fmt.Printf("Total:                              %v\n", s.Sum())
	fmt.Println("======================== Step ========================")
}

func Run(ctx context.Context, s io.Reader, c *Config) error {
	p, err := parser.Parse(s)
	if err != nil {
		return err
	}

	return NewInterpreter(p, c).Run(ctx)
}

func NewInterpreter(p *ast.Program, c *Config) *Interpreter {
	return &Interpreter{
		Program: p,
		Config:  c,
		Memory:  make([]byte, c.MemorySize),
		Pointer: 0,
	}
}

func (i *Interpreter) Run(ctx context.Context) error {
	p, err := optimizer.NewOptimizer().Optimize(i.Program)
	if err != nil {
		return err
	}

	err = i.runExpressions(ctx, p.Expressions)
	i.Step.ShowStep()
	if errors.Is(err, ErrInputFinished) && !i.Config.RaiseErrorOnEOF {
		return nil
	}
	return err

}

func (i *Interpreter) runExpressions(ctx context.Context, exprs []ast.Expression) error {
	for _, expr := range exprs {
		if err := i.runExpression(ctx, expr); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) runExpression(ctx context.Context, expr ast.Expression) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	switch e := expr.(type) {
	case *ast.PointerIncrementExpression:
		if i.Pointer == len(i.Memory)-1 && i.Config.RaiseErrorOnOverflow {
			return fmt.Errorf("%w: %d to pointer overflow, on %d:%d", ErrMemoryOverflow, i.Pointer, e.StartPos(), e.EndPos())
		}
		i.Pointer += 1
		i.Step.PointerIncrementExpression++
	case *ast.PointerDecrementExpression:
		if i.Pointer == 0 && i.Config.RaiseErrorOnOverflow {
			return fmt.Errorf("%w: %d to pointer underflow, on %d:%d", ErrMemoryOverflow, i.Pointer, e.StartPos(), e.EndPos())
		}
		i.Pointer -= 1
		i.Step.PointerDecrementExpression++
	case *ast.ValueIncrementExpression:
		if i.Memory[i.Pointer] == 255 && i.Config.RaiseErrorOnOverflow {
			return fmt.Errorf("%w: %d to memory overflow, on %d:%d", ErrMemoryOverflow, i.Pointer, e.StartPos(), e.EndPos())
		}
		i.Memory[i.Pointer] += 1
		i.Step.PointerIncrementExpression++
	case *ast.ValueDecrementExpression:
		if i.Memory[i.Pointer] == 0 && i.Config.RaiseErrorOnOverflow {
			return fmt.Errorf("%w: %d to memory underflow, on %d:%d", ErrMemoryOverflow, i.Pointer, e.StartPos(), e.EndPos())
		}
		i.Memory[i.Pointer] -= 1
		i.Step.PointerDecrementExpression++
	case *ast.OutputExpression:
		if _, err := i.Config.Writer.Write([]byte{i.Memory[i.Pointer]}); err != nil {
			return err
		}
		i.Step.OutputExpression++
	case *ast.InputExpression:
		b := make([]byte, 1)
		if _, err := i.Config.Reader.Read(b); err != nil {
			if errors.Is(err, io.EOF) {
				return ErrInputFinished
			}
			return err
		}
		i.Memory[i.Pointer] = b[0]
		i.Step.InputExpression++
	case *ast.WhileExpression:
		for i.Memory[i.Pointer] != 0 {
			if err := i.runExpressions(ctx, e.Body); err != nil {
				return err
			}
		}
		fmt.Println(e.Body)
		i.Step.WhileExpression++
	case *ast.MultiplePointerIncrementExpression:
		if i.Pointer + 1 + e.Count > len(i.Memory)  && !i.Config.RaiseErrorOnOverflow {
			return fmt.Errorf("%w: %d to pointer overflow, on %d:%d", ErrMemoryOverflow, i.Pointer, e.StartPos(), e.EndPos())
		}
		i.Pointer += e.Count
		i.Step.MultiplePointerIncrementExpression++
	case *ast.MultiplePointerDecrementExpression:
		if i.Pointer - e.Count < 0 && !i.Config.RaiseErrorOnOverflow {
			return fmt.Errorf("%w: %d to pointer overflow, on %d:%d", ErrMemoryOverflow, i.Pointer, e.StartPos(), e.EndPos())
		}
		i.Pointer -= e.Count
		i.Step.MultiplePointerDecrementExpression++
	case *ast.MultipleValueIncrementExpression:
		if i.Memory[i.Pointer] + byte(e.Count) > 255 && !i.Config.RaiseErrorOnOverflow {
			return fmt.Errorf("%w: %d to memory overflow, on %d:%d", ErrMemoryOverflow, i.Pointer, e.StartPos(), e.EndPos())
		}
		i.Memory[i.Pointer] += byte(e.Count)
		i.Step.MultipleValueIncrementExpression++
	case *ast.MultipleValueDecrementExpression:
		if i.Memory[i.Pointer] - byte(e.Count) < 0 && !i.Config.RaiseErrorOnOverflow {
			return fmt.Errorf("%w: %d to memory overflow, on %d:%d", ErrMemoryOverflow, i.Pointer, e.StartPos(), e.EndPos())
		}
		i.Memory[i.Pointer] -= byte(e.Count)
		i.Step.MultiplePointerDecrementExpression++
	case *ast.LoadZeroExpression:
		i.Memory[i.Pointer] = 0
		i.Step.LoadZeroExpression++
	case *ast.CopyExpression:
		for _, copy := range e.Copys {
			i.Memory[i.Pointer + copy.Phase] = byte(i.Pointer) * byte(copy.Count)
		}
		i.Step.CopyExpression++
	}
	return nil
}
