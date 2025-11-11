package constants

type AutomataState int
type TokenType int

const (
	ReadIdent AutomataState = iota
	ReadNum
	ReadSpecial
	ReadingEnd
	Reading
	Init
	EOF
	ERROR_STATE
)
const (
	IDENT TokenType = iota
	NUMBER
	STRING
	INT
	VAR
	BEGIN
	END
	FOR
	READ
	WRITE
	TO
	LPAREN
	RPAREN
	ADDITTIVE
	MULTIPLICATIVE
	COMPARATIVE
	ASSIGN
	SEMI
	COMMA
	IFT
	THENT
	ELSET
	FUNCT
	RETURNT
	LCURL
	RCURL
	ERROR
	ENDOFSTREAM
	AND
	OR

	VOID
	INTEGERT
)

type ItemType int

const (
	F ItemType = iota
	Dynamic
	Integer
	String
)
