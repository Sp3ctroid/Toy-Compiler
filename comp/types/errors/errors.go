package errors

import "fmt"

type WrongTokenError struct {
	Line   int
	Column int
	Val    string
	Detail string
}

type UnknownTokenError struct {
	Line   int
	Column int
	Val    string
	Detail string
}

func (err UnknownTokenError) Error() string {
	return fmt.Sprintf("%d:%d Unknown identifier or keyword: %s %s", err.Line, err.Column, err.Val, err.Detail)
}

func (err WrongTokenError) Error() string {
	return fmt.Sprintf("%d:%d Forbidden language: %s %s", err.Line, err.Column, err.Val, err.Detail)
}
