package symboltable

import (
	"compiler/constants"

	"github.com/llir/llvm/ir/value"
)

type Table_item struct {
	Scope_level int
	T           constants.ItemType
	LLIRVal     value.Value
}

func NewTableItem(sl int, t constants.ItemType, LLIRVal value.Value) Table_item {
	return Table_item{sl, t, LLIRVal}
}

type Symbol_table struct {
	Items map[string]Table_item
}

func NewSymbolTable() *Symbol_table {
	return &Symbol_table{make(map[string]Table_item)}
}

func NewSymbolTableArr() []*Symbol_table {
	return make([]*Symbol_table, 256)
}
