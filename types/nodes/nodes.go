package nodes

import (
	"compiler/constants"
	tok "compiler/types/token"
	"fmt"

	"github.com/llir/llvm/ir/value"
)

type Node interface {
}

type ValueExtractor interface {
	Node
	GetLLirValue() value.Value
}

type TypeExtractor interface {
	Node
	GetType() constants.ItemType
}

type Program struct {
	Tree_root RootNode
}

type RootNode struct {
	Prog_name  string
	Var_n      VarNode
	Operator_n OperatorNode
}

type OperatorNode struct {
	Text      string
	Operators []interface{}
}

type VarNode struct {
	Text string
	Ids  []IdentNode
}

type IdentNode struct {
	Name    string
	LLIRVal value.Value
	Type    constants.ItemType
}

func (n IdentNode) GetLLirValue() value.Value {
	return n.LLIRVal
}

func (n IdentNode) GetType() constants.ItemType {
	return n.Type
}

func NewIdentNode(t tok.Token, LLIRVal value.Value, tp constants.ItemType) IdentNode {
	return IdentNode{t.Text, LLIRVal, tp}
}

type ReadOpNode struct {
	Text string
	Id   IdentNode
}

func NewReadOpNode(t string, id IdentNode) ReadOpNode {
	return ReadOpNode{t, id}
}

type NumNode struct {
	Number  string
	LLIRVal value.Value
}

func (n NumNode) GetLLirValue() value.Value {
	return n.LLIRVal
}

func NewNumNode(n string, llirval value.Value) NumNode {
	return NumNode{n, llirval}
}

type WriteOpNode struct {
	Text     string
	To_write interface{}
	LLIRVal  value.Value
}

func (n WriteOpNode) GetLLirValue() value.Value {
	return n.LLIRVal
}

func NewWriteOpNode(t string, something interface{}, llirval value.Value) (WriteOpNode, error) {

	if tok, ok := something.(tok.Token); ok {
		switch tok.TokType {
		case constants.STRING:
			return WriteOpNode{t, tok.Text, llirval}, nil
		case constants.IDENT:
			return WriteOpNode{t, tok.Text, llirval}, nil
		case constants.NUMBER:
			return WriteOpNode{t, tok.Text, llirval}, nil
		default:
			return WriteOpNode{}, fmt.Errorf("Write Killed")
		}
	} else {
		return WriteOpNode{}, fmt.Errorf("Write Killed")
	}
}

type AdditiveOpNode struct {
	Op      string
	Left    interface{}
	Right   interface{}
	TypeOf  constants.ItemType
	LLIRVal value.Value
}

func (n AdditiveOpNode) GetLLirValue() value.Value {
	return n.LLIRVal
}

func (n AdditiveOpNode) GetType() constants.ItemType {
	return n.TypeOf
}

func NewAdditiveOpNode(op string, left, right interface{}, TypeOf constants.ItemType, llirVal value.Value) AdditiveOpNode {
	return AdditiveOpNode{op, left, right, TypeOf, llirVal}
}

type MultiplicativeOpNode struct {
	Op      string
	Left    interface{}
	Right   interface{}
	TypeOf  constants.ItemType
	LLIRVal value.Value
}

func (n MultiplicativeOpNode) GetLLirValue() value.Value {
	return n.LLIRVal
}

func (n MultiplicativeOpNode) GetType() constants.ItemType {
	return n.TypeOf
}

func NewMultiplicativeOpNode(op string, left, right interface{}, TypeOf constants.ItemType, LLIRVal value.Value) MultiplicativeOpNode {
	return MultiplicativeOpNode{op, left, right, TypeOf, LLIRVal}
}

type AssignOpNode struct {
	Op     string
	Left   IdentNode
	Right  interface{}
	TypeOf constants.ItemType
}

func (n AssignOpNode) GetType() constants.ItemType {
	return n.TypeOf
}

func NewAssignOpNode(op string, left IdentNode, right interface{}, TypeOf constants.ItemType) AssignOpNode {
	return AssignOpNode{op, left, right, TypeOf}
}

type ForNode struct {
	Text       string
	R          RangeClause
	Operator_n OperatorNode
}

func NewForNode(t string, r RangeClause, o OperatorNode) ForNode {
	return ForNode{t, r, o}
}

type ComparisonNode struct {
	Text  string
	Left  interface{}
	Op    string
	Right interface{}
}

func NewComparisonNode(t string, l interface{}, op string, r interface{}) ComparisonNode {
	return ComparisonNode{t, l, op, r}
}

type IfStatementNode struct {
	Text      string
	Stmnt     interface{}
	Body      []interface{}
	Else_body []interface{}
}

func NewIfNode(t string, s interface{}, b []interface{}, eb []interface{}) IfStatementNode {
	return IfStatementNode{t, s, b, eb}
}

type RangeClause struct {
	From string
	To   string
}

func NewRangeClause(f, t string) RangeClause {
	return RangeClause{f, t}
}

type FuncNode struct {
	Text string
	Name string
	Args []IdentNode
	Body OperatorNode
}

func NewFuncNode(t string, name string, a []IdentNode, body OperatorNode) FuncNode {
	return FuncNode{t, name, a, body}
}

type FuncCallNode struct {
	Text     string
	FuncName string
	Args     []interface{}
	LLIRVal  value.Value
}

func (n FuncCallNode) GetLLirValue() value.Value {
	return n.LLIRVal
}

func NewFuncCallNode(t string, name string, a []interface{}, llirval value.Value) FuncCallNode {
	return FuncCallNode{t, name, a, llirval}
}

type BooleanOpNode struct {
	Op      string
	Left    interface{}
	Right   interface{}
	TypeOf  constants.ItemType
	LLIRVal value.Value
}

func (n BooleanOpNode) GetLLirValue() value.Value {
	return n.LLIRVal
}

func (n BooleanOpNode) GetType() constants.ItemType {
	return n.TypeOf
}

func NewBooleanOpNode(op string, left, right interface{}, TypeOf constants.ItemType, LLIRVal value.Value) BooleanOpNode {
	return BooleanOpNode{op, left, right, TypeOf, LLIRVal}
}

type ReturnNode struct {
	ToReturn interface{}
}

func NewReturnNode(TR interface{}) ReturnNode {
	return ReturnNode{ToReturn: TR}
}
