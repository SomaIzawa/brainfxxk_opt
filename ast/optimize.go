package ast

type MultiplePointerIncrementExpression struct {
	Count       int
	Expressions []Expression
}

func (e *MultiplePointerIncrementExpression) StartPos() int {
	return e.Expressions[0].StartPos()
}

func (e *MultiplePointerIncrementExpression) EndPos() int {
	return e.Expressions[len(e.Expressions)-1].EndPos()
}

func (e *MultiplePointerIncrementExpression) Bytes() []byte {
	b := []byte{}
	for _, expr := range e.Expressions {
		b = append(b, expr.Bytes()...)
	}
	return b
}

func (e *MultiplePointerIncrementExpression) String() string {
	return string(e.Bytes())
}

type MultiplePointerDecrementExpression struct {
	Count       int
	Expressions []Expression
}

func (e *MultiplePointerDecrementExpression) StartPos() int {
	return e.Expressions[0].StartPos()
}

func (e *MultiplePointerDecrementExpression) EndPos() int {
	return e.Expressions[len(e.Expressions)-1].EndPos()
}

func (e *MultiplePointerDecrementExpression) Bytes() []byte {
	b := []byte{}
	for _, expr := range e.Expressions {
		b = append(b, expr.Bytes()...)
	}
	return b
}

func (e *MultiplePointerDecrementExpression) String() string {
	return string(e.Bytes())
}

type MultipleValueIncrementExpression struct {
	Count       int
	Expressions []Expression
}

func (e *MultipleValueIncrementExpression) StartPos() int {
	return e.Expressions[0].StartPos()
}

func (e *MultipleValueIncrementExpression) EndPos() int {
	return e.Expressions[len(e.Expressions)-1].EndPos()
}

func (e *MultipleValueIncrementExpression) Bytes() []byte {
	b := []byte{}
	for _, expr := range e.Expressions {
		b = append(b, expr.Bytes()...)
	}
	return b
}

func (e *MultipleValueIncrementExpression) String() string {
	return string(e.Bytes())
}

type MultipleValueDecrementExpression struct {
	Count       int
	Expressions []Expression
}

func (e *MultipleValueDecrementExpression) StartPos() int {
	return e.Expressions[0].StartPos()
}

func (e *MultipleValueDecrementExpression) EndPos() int {
	return e.Expressions[len(e.Expressions)-1].EndPos()
}

func (e *MultipleValueDecrementExpression) Bytes() []byte {
	b := []byte{}
	for _, expr := range e.Expressions {
		b = append(b, expr.Bytes()...)
	}
	return b
}

func (e *MultipleValueDecrementExpression) String() string {
	return string(e.Bytes())
}

type LoadZeroExpression struct {
	Expressions []Expression
}

func (e *LoadZeroExpression) StartPos() int {
	return e.Expressions[0].StartPos()
}

func (e *LoadZeroExpression) EndPos() int {
	return e.Expressions[len(e.Expressions)-1].EndPos()
}

func (e *LoadZeroExpression) Bytes() []byte {
	b := []byte{}
	for _, expr := range e.Expressions {
		b = append(b, expr.Bytes()...)
	}
	return b
}

func (e *LoadZeroExpression) String() string {
	return string(e.Bytes())
}