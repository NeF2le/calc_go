package calculation

import "errors"

var (
	ErrInvalidExpression = errors.New("invalid expression")
	ErrDivisionByZero	 = errors.New("division by zero")
	ErrInvalidBrackets	 = errors.New("invalid brackets in expression")
)