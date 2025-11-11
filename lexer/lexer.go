package lexer

import (
	"compiler/constants"
	h "compiler/helpers"
	"compiler/types/errors"
	tok "compiler/types/token"
)

type Lexer struct {
	Text             string
	CurrentTokenText string
	Line             int
	Column           int
	Pos              int
	Token_state      constants.AutomataState
	Read_state       constants.AutomataState
	Current_char     byte
}

func NewLexer(progText string) *Lexer {
	return &Lexer{Text: progText, Line: 1, Column: 1, Pos: 0, CurrentTokenText: "", Token_state: constants.Init, Read_state: constants.Reading}
}

func (lex *Lexer) Advance() {
	lex.Pos++
	lex.Column++
	lex.check_EOF()
	if lex.Read_state != constants.EOF {
		lex.Current_char = lex.Text[lex.Pos]
	}
}

func (lex *Lexer) NextToken() (tok.Token, error) {

	for lex.Read_state != constants.ReadingEnd && lex.Read_state != constants.EOF {
		err := lex.ReadChar()
		if err != nil {
			return tok.Token{"NULL", constants.ERROR}, err
		}
	}

	current_token := lex.CurrentTokenText

	if current_token == "" && lex.Read_state == constants.EOF {
		return tok.Token{current_token, constants.ENDOFSTREAM}, nil
	}

	switch current_token {
	case "+":
		lex.Reset()
		return tok.Token{current_token, constants.ADDITTIVE}, nil
	case "-":
		lex.Reset()
		return tok.Token{current_token, constants.ADDITTIVE}, nil
	case "*":
		lex.Reset()
		return tok.Token{current_token, constants.MULTIPLICATIVE}, nil
	case "/":
		lex.Reset()
		return tok.Token{current_token, constants.MULTIPLICATIVE}, nil
	case "VAR":
		lex.Reset()
		return tok.Token{current_token, constants.VAR}, nil
	case "BEGIN":
		lex.Reset()
		return tok.Token{current_token, constants.BEGIN}, nil
	case "END":
		lex.Reset()
		return tok.Token{current_token, constants.END}, nil
	case "FOR":
		lex.Reset()
		return tok.Token{current_token, constants.FOR}, nil
	case "READ":
		lex.Reset()
		return tok.Token{current_token, constants.READ}, nil
	case "WRITE":
		lex.Reset()
		return tok.Token{current_token, constants.WRITE}, nil
	case "TO":
		lex.Reset()
		return tok.Token{current_token, constants.TO}, nil
	case "IF":
		lex.Reset()
		return tok.Token{current_token, constants.IFT}, nil
	case "THEN":
		lex.Reset()
		return tok.Token{current_token, constants.THENT}, nil
	case "ELSE":
		lex.Reset()
		return tok.Token{current_token, constants.ELSET}, nil
	case ">=":
		lex.Reset()
		return tok.Token{current_token, constants.COMPARATIVE}, nil
	case "<=":
		lex.Reset()
		return tok.Token{current_token, constants.COMPARATIVE}, nil
	case "<":
		lex.Reset()
		return tok.Token{current_token, constants.COMPARATIVE}, nil
	case ">":
		lex.Reset()
		return tok.Token{current_token, constants.COMPARATIVE}, nil
	case "==":
		lex.Reset()
		return tok.Token{current_token, constants.COMPARATIVE}, nil
	case "(":
		lex.Reset()
		return tok.Token{current_token, constants.LPAREN}, nil
	case ")":
		lex.Reset()
		return tok.Token{current_token, constants.RPAREN}, nil
	case "=":
		lex.Reset()
		return tok.Token{current_token, constants.ASSIGN}, nil
	case ";":
		lex.Reset()
		return tok.Token{current_token, constants.SEMI}, nil
	case ",":
		lex.Reset()
		return tok.Token{current_token, constants.COMMA}, nil
	case "{":
		lex.Reset()
		return tok.Token{current_token, constants.LCURL}, nil
	case "}":
		lex.Reset()
		return tok.Token{current_token, constants.RCURL}, nil
	case "FUNC":
		lex.Reset()
		return tok.Token{current_token, constants.FUNCT}, nil
	case "RETURN":
		lex.Reset()
		return tok.Token{current_token, constants.RETURNT}, nil
	case "INT":
		lex.Reset()
		return tok.Token{current_token, constants.INT}, nil
	case "STRING":
		lex.Reset()
		return tok.Token{current_token, constants.STRING}, nil
	case "&&":
		lex.Reset()
		return tok.Token{current_token, constants.AND}, nil
	case "||":
		lex.Reset()
		return tok.Token{current_token, constants.OR}, nil

	default:
		switch lex.Token_state {
		case constants.ReadIdent:
			lex.Reset()
			return tok.Token{current_token, constants.IDENT}, nil
		case constants.ReadNum:
			lex.Reset()
			return tok.Token{current_token, constants.NUMBER}, nil
		}
	}

	lex.Read_state = constants.ERROR_STATE
	lex.Token_state = constants.ERROR_STATE
	return tok.Token{"NULL", constants.ERROR}, errors.UnknownTokenError{lex.Line, lex.Column, lex.CurrentTokenText, ""}
}

func (lex *Lexer) Reset() {
	if lex.Read_state != constants.EOF {
		lex.CurrentTokenText = ""
		lex.Read_state = constants.Reading
		lex.Token_state = constants.Init
	}

}

func (lex *Lexer) ReadChar() error {

	lex.skip_whitespace()
	if lex.Read_state == constants.EOF {
		return nil
	}

	lex.Current_char = lex.Text[lex.Pos]

	switch lex.Token_state {
	case constants.Init:
		lex.define_state(lex.Current_char)

	case constants.ReadIdent:
		lex.CurrentTokenText += string(lex.Current_char)
		lex.Advance()
		if lex.Read_state != constants.EOF {
			if lex.Current_char == ' ' || h.Is_special(lex.Current_char) {
				lex.Read_state = constants.ReadingEnd
			}
		} else {
			return nil
		}

	case constants.ReadNum:
		if h.Is_char(lex.Current_char) {
			lex.Token_state = constants.ERROR_STATE
			lex.Read_state = constants.ERROR_STATE
			return errors.WrongTokenError{lex.Line, lex.Column, lex.CurrentTokenText, "(Was reading Number, but encountered Text. Typo in number, or maybe you tried to name variable starting with digit?)"}
		}

		lex.CurrentTokenText += string(lex.Current_char)
		lex.Advance()

		if lex.Read_state != constants.EOF {
			if lex.Current_char == ' ' || h.Is_special(lex.Current_char) {
				lex.Read_state = constants.ReadingEnd
			}
		} else {
			return nil
		}

	case constants.ReadSpecial:
		lex.CurrentTokenText += string(lex.Current_char)
		lex.Advance()

		if lex.Read_state != constants.EOF {
			if (h.Is_comparative(lex.CurrentTokenText[0]) && lex.Current_char == '=') || h.Is_boolean(lex.CurrentTokenText[0]) && h.Is_boolean(lex.Current_char) {
				lex.CurrentTokenText += string(lex.Current_char)
				lex.Advance()
				lex.Read_state = constants.ReadingEnd
			} else {
				lex.Read_state = constants.ReadingEnd
			}
		} else {
			return nil
		}
	}

	return nil

}

func (lex *Lexer) define_state(char byte) {
	if h.Is_digit(char) {
		lex.Token_state = constants.ReadNum
	} else if h.Is_char(char) {
		lex.Token_state = constants.ReadIdent
	} else if h.Is_special(char) || h.Is_comparative(char) {
		lex.Token_state = constants.ReadSpecial
	}
}

func (lex *Lexer) skip_whitespace() {
	lex.Current_char = lex.Text[lex.Pos]

	for lex.Current_char == ' ' || lex.Current_char == '\n' || lex.Current_char == '\t' || lex.Current_char == '\r' {

		if lex.Read_state == constants.EOF {
			return
		}

		if lex.Current_char == '\n' {
			lex.Column = 1
			lex.Line++
		}
		lex.Advance()
	}

}

func (lex *Lexer) check_EOF() {
	if lex.Pos >= len(lex.Text) {
		lex.Read_state = constants.EOF
	}
}
