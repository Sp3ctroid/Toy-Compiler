package token

import (
	"compiler/constants"
	"fmt"
)

type Token struct {
	Text    string
	TokType constants.TokenType
}

func (token Token) String() string {
	return fmt.Sprintf("VAL: %s  TYPE: %v", token.Text, token.TokType)
}
