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
