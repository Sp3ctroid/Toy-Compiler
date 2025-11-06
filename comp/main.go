// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TODO ЛИСТ																																																	//
//
//																																																				//																																																	//
//																																																				//
//	 1. Отметил места где надо либо вынести код в метод соответствующего класса, либо надо отнять этот метод у другого класса и дать другому (часто Parser и Symbol Table)										//
//	    возможно придётся создать еще один класс "глобальной" таблицы, которая хранит таблицы символов (хранится щас в Parser которая)																			//
//	    ИЛИ хранить какой-то контекст в котором эта таблица тоже будет для всей проги																															//
//																																																				//
//	   																																																			//																																																		//
//																																																				//
//	 2. Добавить обработку аргументов запуска чтобы можно было из текста читать код моего языка и компилировать его в исполняемый с именем																		//
//	    																																																		//
//																																																				//
//	 3. Grapher пока сломан и не работает																																										//
//	    																																																		//
//																																																				//
//	 4. Есть огромная конструкция считывания IDENT токена. Она получается от того, что IDENT это может быть и вызов функции, и переменная. Ее либо надо упростить либо вынести в отдельную функцию.				//
//	    																																																		//
//																																																				//
//																																																				//
//
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	llir "github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	constants "compiler/constants"
	h "compiler/helpers"
	lex "compiler/lexer"
	nd "compiler/types/nodes"
	st "compiler/types/symbol_table"
	tok "compiler/types/token"
)

var SCANF_FUNC *llir.Func
var PRINT_FUNC *llir.Func

var SCANF_FORMAT_INT *llir.Global
var PRINTF_FORMAT_INT *llir.Global

var SCANF_FORMAT_TYPE *types.ArrayType
var PRINTF_FORMAT_TYPE *types.ArrayType

var CALCULATED_PTR_FOR_FORMAT_SCANF bool
var CALCULATED_PTR_FOR_FORMAT_PRINTF bool

func main() {

	// 	program_text :=

	// 		`INT x;

	// FUNC add(a, b){
	// 	INT sum;
	// 	sum = a + b;
	// 	WRITE sum;

	// 	RETURN sum

	// }

	// FUNC main(){

	// 	INT i;
	// 	READ i;
	// 	WRITE i;

	// 	INT j;
	// 	READ j;
	// 	WRITE j;

	// 	add(i, j);

	// 	RETURN i
	// }

	// `
	var filename string
	var outputname string
	flag.StringVar(&filename, "FN", "", "your file to compile")
	flag.StringVar(&outputname, "ON", "", "your output file name")

	flag.Parse()

	if filename == "" {
		panic("You need to specify file to compile (-h for help)")
	}

	if outputname == "" {
		panic("You need to specify output file name (-h for help)")
	}

	read_file, ok := os.ReadFile(filename)

	if ok != nil {
		panic("error reading file")
	}

	read_file, _ = os.ReadFile(filename)
	llvm_ir_output := start_compiling(string(read_file))

	os.WriteFile(outputname, []byte(llvm_ir_output), 0644)
}

func start_compiling(code string) string {
	p := NewParser(code + "\n")
	p.InitParser()

	p.Program()
	if _, ok := p.sts["main"]; !ok {
		fmt.Println("main function is not defined, no entry-point")
		os.Exit(1)
	}
	//ctx := context.Background()
	//g, err := graphviz.New(ctx)
	//if err != nil {
	//	fmt.Println("Graphviz creation failed")
	//	os.Exit(1)
	//}

	//graph, _ := g.Graph()
	//gr := graph_t.NewGrapher(p.prog_tree.Tree_root)

	//gr.Graph_tree(graph)

	//var buf bytes.Buffer
	//if err := g.Render(ctx, graph, "dot", &buf); err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println(buf.String())

	// fmt.Println(p.prog_tree.Tree_root.Var_n.Ids)
	// fmt.Println(p.prog_tree.Tree_root.Operator_n.Operators...)

	// fmt.Println()

	data := fmt.Sprint(p.gen.M)

	return data

}

type Generator struct {
	M                  *llir.Module
	F                  *llir.Func
	B                  *llir.Block
	strInd             string
	Printf_Format_calc *llir.InstGetElementPtr
	Scanf_Format_calc  *llir.InstGetElementPtr
}

func (g *Generator) GlobalFuncDeclares(M *llir.Module) {
	SCANF_FUNC = M.NewFunc("scanf", types.I32, llir.NewParam("format", types.I8Ptr))
	SCANF_FUNC.Sig.Variadic = true
	SCANF_FORMAT_INT = M.NewGlobalDef("format_read_int", constant.NewCharArrayFromString("%d\x00"))

	PRINT_FUNC = M.NewFunc("printf", types.I32, llir.NewParam("format", types.I8Ptr))
	PRINT_FUNC.Sig.Variadic = true
	PRINTF_FORMAT_INT = M.NewGlobalDef("format_write_int", constant.NewCharArrayFromString("%d\n\x00"))

	SCANF_FORMAT_TYPE = types.NewArray(3, types.I8)
	PRINTF_FORMAT_TYPE = types.NewArray(4, types.I8)

	CALCULATED_PTR_FOR_FORMAT_PRINTF = false
	CALCULATED_PTR_FOR_FORMAT_SCANF = false
}

type Parser struct {
	lex             lex.Lexer
	gen             Generator
	currentToken    tok.Token
	prog_tree       nd.Program
	sts             map[string][]*st.Symbol_table
	currentFunction string
	scope_level     int
}

func (p Parser) Exists(var_name string) bool {

	for i := 1; i <= p.scope_level; i++ {
		_, ok := p.sts[p.currentFunction][i].Items[var_name]

		if ok {
			return true
		}
	}

	_, ok := p.sts["global"][0].Items[var_name]

	return ok
}

func (p Parser) ExistsOnCurrentLevel(var_name string) bool {
	_, ok := p.sts[p.currentFunction][p.scope_level].Items[var_name]

	return ok
}

func NewParser(prog_text string) *Parser {

	return &Parser{*lex.NewLexer(prog_text), Generator{strInd: "0"}, tok.Token{}, nd.Program{}, make(map[string][]*st.Symbol_table), "global", 0}

}

func (p *Parser) InitParser() {
	new_token, err := p.lex.NextToken()

	if err != nil {
		return
	}

	global_table := st.NewSymbolTable()

	p.sts[p.currentFunction] = append(p.sts[p.currentFunction], global_table)
	p.currentToken = new_token

	p.gen.M = llir.NewModule()

	p.gen.GlobalFuncDeclares(p.gen.M)

}

func (p *Parser) Eat(t constants.TokenType) {
	if p.currentToken.TokType == t {
		new_token, err := p.lex.NextToken()

		if err != nil {
			fmt.Printf("ERROR GETTING NEXT TOKEN")
			os.Exit(1)
		}

		p.currentToken = new_token
	} else {
		fmt.Printf("%d:%d Unexpected Token: %s\n, Expected Token: %v\n", p.lex.Line, p.lex.Column, p.currentToken.Text, t)
		os.Exit(1)
	}
}

func (p *Parser) Program() {

	p.Eat(constants.INT)
	p.prog_tree.Tree_root.Var_n.Text = "GLOBAL VARS"
	p.prog_tree.Tree_root.Var_n.Ids = p.VarsDecl(constants.Integer)
	p.Eat(constants.SEMI)

	p.prog_tree.Tree_root.Operator_n.Text = "OPERATORS"
	for p.currentToken.TokType != constants.ENDOFSTREAM {
		p.prog_tree.Tree_root.Operator_n.Operators = append(p.prog_tree.Tree_root.Operator_n.Operators, p.Operators())
	}

}

func (p *Parser) VarsDecl(typ constants.ItemType) []nd.IdentNode {

	ids := []nd.IdentNode{}
	zero := constant.NewInt(types.I32, 0) // МОЖНО ВЫНЕСТИ В КОНСТАНТУ
	for p.currentToken.TokType != constants.SEMI {
		var llirval value.Value
		var_name := p.currentToken.Text
		if p.Exists(var_name) {
			fmt.Printf("%s redeclared at %d:%d\n", var_name, p.lex.Line, p.lex.Column)
			os.Exit(1)
		} else {
			if p.currentFunction == "global" {
				llirval = p.gen.M.NewGlobalDef(p.currentToken.Text, zero)
			} else {
				llirval = p.gen.B.NewAlloca(types.I32)
				p.gen.B.NewStore(zero, llirval)
			}

			n_ti := st.NewTableItem(p.scope_level, typ, llirval)
			p.sts[p.currentFunction][p.scope_level].Items[var_name] = n_ti
		}

		ids = append(ids, nd.NewIdentNode(p.currentToken, llirval, typ))
		p.Eat(constants.IDENT)
		if p.currentToken.TokType == constants.COMMA {
			p.Eat(constants.COMMA)
		} else if p.currentToken.TokType == constants.SEMI {
			break
		}
	}

	return ids

}

func (p *Parser) Operators() nd.Node {

	var operators_segment interface{}

	if p.currentFunction == "global" {
		if p.currentToken.TokType != constants.FUNCT {
			fmt.Printf("Expected function declaration.\n")
			p.Eat(constants.FUNCT)
		} else {
			p.Eat(constants.FUNCT)
			func_name := p.currentToken.Text
			if p.ExistsOnCurrentLevel(func_name) {
				fmt.Printf("Function %s is already declared.\n", func_name)
				os.Exit(1)
			}

			var llirParameteres []*llir.Param
			p.scope_level++
			p.sts[func_name] = st.NewSymbolTableArr()
			p.sts[func_name][p.scope_level] = st.NewSymbolTable()

			p.Eat(constants.IDENT)
			p.Eat(constants.LPAREN)
			args := []nd.IdentNode{}
			for p.currentToken.TokType != constants.RPAREN {
				if p.currentToken.TokType == constants.IDENT {
					LLIRParam := llir.NewParam(p.currentToken.Text, types.I32)
					id := nd.NewIdentNode(p.currentToken, LLIRParam, constants.Integer)
					args = append(args, id)

					llirParameteres = append(llirParameteres, LLIRParam)
					p.sts[func_name][p.scope_level].Items[id.Name] = st.NewTableItem(p.scope_level, constants.Integer, LLIRParam)
				}
				p.Eat(constants.IDENT)
				if p.currentToken.TokType != constants.RPAREN {
					p.Eat(constants.COMMA)
				}
			}
			p.Eat(constants.RPAREN)

			LLIRFunc := p.gen.M.NewFunc(func_name, types.I32, llirParameteres...)
			p.gen.F = LLIRFunc
			p.sts["global"][0].Items[func_name] = st.NewTableItem(0, constants.F, LLIRFunc)
			p.gen.B = p.gen.F.NewBlock("")

			p.currentFunction = func_name

			p.Eat(constants.LCURL)

			var func_body []interface{}
			for p.currentToken.TokType != constants.RCURL {
				func_body = append(func_body, p.Operators())
			}
			p.Eat(constants.RCURL)
			func_ops_n := nd.OperatorNode{"BODY", func_body}
			func_n := nd.NewFuncNode("FUNC_DECL", func_name, args, func_ops_n)
			p.currentFunction = "global"
			p.scope_level--
			CALCULATED_PTR_FOR_FORMAT_SCANF = false
			CALCULATED_PTR_FOR_FORMAT_PRINTF = false
			return func_n
		}
	} else {
		switch p.currentToken.TokType {
		case constants.READ:
			p.Eat(constants.READ)
			if !p.Exists(p.currentToken.Text) {
				fmt.Printf("%s undefined\n", p.currentToken.Text)
				os.Exit(1)
			}

			_, level := p.TableFirstOccurance(p.currentToken.Text)
			var item st.Table_item
			var llirVal value.Value
			if level == 0 {
				item = p.sts["global"][0].Items[p.currentToken.Text]
			} else {
				item = p.sts[p.currentFunction][level].Items[p.currentToken.Text]
			}
			llirVal = item.LLIRVal
			Id_n := nd.NewIdentNode(p.currentToken, llirVal, item.T)
			p.Eat(constants.IDENT)
			Read_n := nd.NewReadOpNode("READ", Id_n)
			p.Eat(constants.SEMI)
			if !CALCULATED_PTR_FOR_FORMAT_SCANF {
				p.gen.Scanf_Format_calc = p.gen.B.NewGetElementPtr(SCANF_FORMAT_TYPE, SCANF_FORMAT_INT, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
				CALCULATED_PTR_FOR_FORMAT_SCANF = true
			}

			p.gen.B.NewCall(SCANF_FUNC, p.gen.Scanf_Format_calc, llirVal)
			return Read_n

		case constants.WRITE:
			p.Eat(constants.WRITE)
			if p.currentToken.TokType == constants.IDENT {
				if !p.Exists(p.currentToken.Text) {
					fmt.Printf("%s undefined\n", p.currentToken.Text)
					os.Exit(1)
				}
			}

			_, level := p.TableFirstOccurance(p.currentToken.Text)
			var item st.Table_item
			var llirVal value.Value
			if level == 0 { //ОЧЕНЬ ЧАСТО ПОВТОРЯЮЩАЯСЯ КОНСТРУКЦИЯ. МОЖНО ВЫНЕСТИ В МЕТОД ТАБЛИЦЫ СКОРЕЕ ВСЕГО?
				item = p.sts["global"][0].Items[p.currentToken.Text]
			} else {
				item = p.sts[p.currentFunction][level].Items[p.currentToken.Text]
			}
			llirVal = item.LLIRVal

			if types.IsPointer(llirVal.Type()) {
				llirVal = p.gen.B.NewLoad(types.I32, llirVal)
			}

			Write_n, err := nd.NewWriteOpNode("WRITE", p.currentToken, llirVal)
			if err != nil {
				fmt.Errorf("ERROR READING WRITE ARGUMENT")
				os.Exit(1)
			}

			if p.currentToken.TokType == constants.IDENT { //ВОТ ЭТО ВООБЩЕ БРЕД КАКОЙ-ТО )))
				p.Eat(constants.IDENT)
			} else if p.currentToken.TokType == constants.NUMBER {
				p.Eat(constants.NUMBER)
			} else if p.currentToken.TokType == constants.STRING {
				p.Eat(constants.STRING)
			} else {
				fmt.Errorf("ERROR WHEN EATING WRITE ARGUMENT")
				os.Exit(1)
			}

			p.Eat(constants.SEMI)
			if !CALCULATED_PTR_FOR_FORMAT_PRINTF {
				p.gen.Printf_Format_calc = p.gen.B.NewGetElementPtr(PRINTF_FORMAT_TYPE, PRINTF_FORMAT_INT, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
				CALCULATED_PTR_FOR_FORMAT_PRINTF = true
			}

			p.gen.B.NewCall(PRINT_FUNC, p.gen.Printf_Format_calc, llirVal)
			return Write_n

		case constants.IFT:
			p.Eat(constants.IFT)

			var comp_n interface{}

			comp_n = p.ComparativeExpression()

			p.Eat(constants.BEGIN)
			var else_operators []interface{}
			var if_operators []interface{}
			p.scope_level++
			p.sts[p.currentFunction][p.scope_level] = st.NewSymbolTable()

			new_if_block := p.gen.F.NewBlock("")
			new_else_block := p.gen.F.NewBlock("")
			end_block := p.gen.F.NewBlock("")

			var condition_llirVal value.Value

			if value_node, ok := comp_n.(nd.ValueExtractor); ok {
				condition_llirVal = value_node.GetLLirValue()
			} else {

			}
			p.gen.B.NewCondBr(condition_llirVal, new_if_block, new_else_block)

			p.gen.B = new_if_block
			for p.currentToken.TokType != constants.END && p.currentToken.TokType != constants.ELSET {
				if_operators = append(if_operators, p.Operators())
			}
			p.gen.B.NewBr(end_block)
			p.gen.B = new_else_block
			switch p.currentToken.TokType {
			case constants.ELSET:

				p.Eat(constants.ELSET)

				for p.currentToken.TokType != constants.END {
					else_operators = append(else_operators, p.Operators())
				}
				p.Eat(constants.END)

			case constants.END:
				p.Eat(constants.END)
			}
			p.gen.B.NewBr(end_block)
			p.gen.B = end_block
			var If_n interface{}

			If_n = nd.NewIfNode("IF STATEMENT", comp_n, if_operators, else_operators)

			p.sts[p.currentFunction][p.scope_level] = nil
			p.scope_level--
			return If_n
		case constants.IDENT:
			if _, ok := p.sts[p.currentToken.Text]; ok {
				func_name := p.currentToken.Text
				p.Eat(constants.IDENT)
				p.Eat(constants.LPAREN)
				var args []interface{}
				var llirArgs []value.Value
				for p.currentToken.TokType != constants.RPAREN {

					switch p.currentToken.TokType {
					case constants.IDENT:
						if !p.Exists(p.currentToken.Text) {
							fmt.Printf("%s undefined\n", p.currentToken.Text)
							os.Exit(1)
						}

						_, level := p.TableFirstOccurance(p.currentToken.Text)
						var item st.Table_item
						var llirVal value.Value
						if level == 0 {
							item = p.sts["global"][0].Items[p.currentToken.Text]
						} else {
							item = p.sts[p.currentFunction][level].Items[p.currentToken.Text]

						}
						llirVal = item.LLIRVal
						id_n := nd.NewIdentNode(p.currentToken, llirVal, item.T)
						args = append(args, id_n)

						if types.IsPointer(id_n.LLIRVal.Type()) {
							llirVal = p.gen.B.NewLoad(types.I32, llirVal)
						}
						llirArgs = append(llirArgs, llirVal)
						p.Eat(constants.IDENT)

					case constants.NUMBER:

						integer_val, _ := strconv.Atoi(p.currentToken.Text)

						n_n := nd.NewNumNode(p.currentToken.Text, constant.NewInt(types.I32, int64(integer_val)))
						args = append(args, n_n)

						p.Eat(constants.NUMBER)

					case constants.STRING:
					}
					if p.currentToken.TokType != constants.RPAREN {
						p.Eat(constants.COMMA)
					}
				}

				p.Eat(constants.RPAREN)
				p.Eat(constants.SEMI)

				func_llir_call := p.gen.B.NewCall(p.sts["global"][0].Items[func_name].LLIRVal, llirArgs...)
				F_call := nd.NewFuncCallNode("FUNC CALL", func_name, args, func_llir_call)

				return F_call
			} else {

				if !p.Exists(p.currentToken.Text) {
					fmt.Printf("%s undefined\n", p.currentToken.Text)
					os.Exit(1)
				}
				_, level := p.TableFirstOccurance(p.currentToken.Text)
				var item st.Table_item
				var llirVal value.Value
				if level == 0 {
					item = p.sts["global"][0].Items[p.currentToken.Text]
					llirVal = item.LLIRVal
				} else {
					item = p.sts[p.currentFunction][level].Items[p.currentToken.Text]
					llirVal = item.LLIRVal
				}

				Id_n := nd.NewIdentNode(p.currentToken, llirVal, item.T)

				p.Eat(constants.IDENT)
				if p.currentToken.TokType == constants.LPAREN {
					fmt.Printf("%s function undefined\n", Id_n.Name)
					os.Exit(1)
				}
				p.Eat(constants.ASSIGN)
				n := p.ComparativeExpression()

				var typeofRnode constants.ItemType

				if type_node_r, ok := n.(nd.TypeExtractor); ok {
					typeofRnode = type_node_r.GetType()
				} else {

				}

				l_type, idx := p.TableFirstOccurance(Id_n.Name)
				if idx == -1 {
					fmt.Printf("%s undefined %d:%d", Id_n.Name, p.lex.Line, p.lex.Column)
				}

				if !IsTypeEqual(l_type, typeofRnode) {
					fmt.Printf("Assign to incompatible type %d:%d\n", p.lex.Line, p.lex.Column)
					os.Exit(1)
				}

				Assign_n := nd.NewAssignOpNode("=", Id_n, n, l_type)

				switch right := n.(type) {
				case nd.NumNode:
					p.gen.B.NewStore(right.LLIRVal, Id_n.LLIRVal)
				case nd.IdentNode:
					p.gen.B.NewStore(right.LLIRVal, Id_n.LLIRVal)
				case nd.AdditiveOpNode:
					p.gen.B.NewStore(right.LLIRVal, Id_n.LLIRVal) //ВИЗИТОР НА ВЫЗОВ НУЖНОГО NEWSTORE
				case nd.BooleanOpNode:
					Id_n.LLIRVal = p.gen.B.NewBitCast(Id_n.LLIRVal, types.I1Ptr)
					p.gen.B.NewStore(right.LLIRVal, Id_n.LLIRVal)

					if level == 0 {
						item = p.sts["global"][0].Items[Id_n.Name]
						item.LLIRVal = Id_n.LLIRVal
						p.sts["global"][0].Items[Id_n.Name] = item
					} else {
						item = p.sts[p.currentFunction][level].Items[Id_n.Name]
						item.LLIRVal = Id_n.LLIRVal
						p.sts[p.currentFunction][level].Items[Id_n.Name] = item
					}

				case nd.MultiplicativeOpNode:
					p.gen.B.NewStore(right.LLIRVal, Id_n.LLIRVal)
				case nd.FuncCallNode:
					p.gen.B.NewStore(right.LLIRVal, Id_n.LLIRVal)
				}

				p.Eat(constants.SEMI)
				return Assign_n
			}

		case constants.FOR:
			p.Eat(constants.FOR)
			current_token := p.currentToken
			p.Eat(constants.NUMBER)
			from := current_token.Text
			p.Eat(constants.TO)
			current_token = p.currentToken
			p.Eat(constants.NUMBER)
			to := current_token.Text
			range_cl_n := nd.NewRangeClause(from, to)
			p.Eat(constants.BEGIN)

			from_integer, _ := strconv.Atoi(from)
			to_integer, _ := strconv.Atoi(to)

			llir_from := constant.NewInt(types.I32, int64(from_integer)) //
			llir_to := constant.NewInt(types.I32, int64(to_integer))     //
			ONE := constant.NewInt(types.I32, 1)                         //
			//
			loop_cond := p.gen.F.NewBlock("loop_cond" + p.gen.strInd) //
			p.gen.strInd = h.String_plus(p.gen.strInd)                //
			loop_body := p.gen.F.NewBlock("loop_body" + p.gen.strInd) //
			p.gen.strInd = h.String_plus(p.gen.strInd)                //
			exit := p.gen.F.NewBlock("loop_exit" + p.gen.strInd)      //
			p.gen.strInd = h.String_plus(p.gen.strInd)                // ВЫНЕСТИ В ОТДЕЛЬНЫЙ МЕТОД ГЕНЕРАТОРА
			//
			p.gen.B.NewBr(loop_cond)                //
			counter := p.gen.B.NewAlloca(types.I32) //
			p.gen.B.NewStore(llir_from, counter)    //
			p.gen.B = loop_cond                     //
			//
			counter_load := p.gen.B.NewLoad(types.I32, counter)                     //
			condition_llir := p.gen.B.NewICmp(enum.IPredSLE, counter_load, llir_to) //
			p.gen.B.NewCondBr(condition_llir, loop_body, exit)                      //

			p.scope_level++
			p.sts[p.currentFunction][p.scope_level] = st.NewSymbolTable()

			p.gen.B = loop_body

			var for_operators []interface{}
			for p.currentToken.TokType != constants.END {
				for_operators = append(for_operators, p.Operators())
			}

			new_counter := p.gen.B.NewAdd(ONE, counter_load)
			p.gen.B.NewStore(new_counter, counter)
			p.gen.B.NewBr(loop_cond)

			p.gen.B = exit

			for_op_n := nd.OperatorNode{"BODY", for_operators}
			for_n := nd.NewForNode("FOR", range_cl_n, for_op_n)
			p.Eat(constants.END)

			p.sts[p.currentFunction][p.scope_level] = nil
			p.scope_level--
			return for_n

		case constants.SEMI:
			p.Eat(constants.SEMI)

		case constants.VAR:
			p.Eat(constants.VAR)
			ids := p.VarsDecl(constants.Integer)
			var_decl_n := nd.VarNode{"VAR DECL", ids}
			p.Eat(constants.SEMI)
			return var_decl_n

		case constants.INT:
			p.Eat(constants.INT)
			ids := p.VarsDecl(constants.Integer)
			int_decl_n := nd.VarNode{"INT", ids}
			p.Eat(constants.SEMI)
			return int_decl_n

		case constants.STRING:
			p.Eat(constants.STRING)
			ids := p.VarsDecl(constants.String)
			str_decl_n := nd.VarNode{"STRING", ids}
			p.Eat(constants.SEMI)
			return str_decl_n
		case constants.RETURNT:
			p.Eat(constants.RETURNT)

			node := p.Expression()
			switch n := node.(type) {
			case nd.IdentNode:
				if !p.Exists(n.Name) {
					fmt.Printf("%s undefined\n", n.Name)
					os.Exit(1)
				}

				_, level := p.TableFirstOccurance(n.Name)
				var item st.Table_item
				var llirVal value.Value
				if level == 0 {
					item = p.sts["global"][0].Items[n.Name]
				} else {
					item = p.sts[p.currentFunction][level].Items[n.Name]

				}
				llirVal = item.LLIRVal

				temp := p.gen.B.NewLoad(types.I32, llirVal)
				p.gen.B.NewRet(temp)

				p.Eat(constants.SEMI)
				return nd.NewReturnNode(n)

			}

		default:
			fmt.Printf("%d:%d Unsupported Syntax: %s\n", p.lex.Line, p.lex.Column, p.lex.CurrentTokenText)
			os.Exit(1)
		}
	}

	return operators_segment
}

func (p *Parser) Expression() nd.Node {
	node := p.Term()

	var typeofLnode constants.ItemType

	if typed_n_l, ok := node.(nd.TypeExtractor); ok {
		typeofLnode = typed_n_l.GetType()
	} else {

	}

	for p.currentToken.TokType == constants.ADDITTIVE || p.currentToken.TokType == constants.OR {
		for p.currentToken.TokType == constants.ADDITTIVE {
			op := p.currentToken.Text
			p.Eat(constants.ADDITTIVE)

			right_n := p.Term()
			var typeofRnode constants.ItemType

			if typed_r_n, ok := right_n.(nd.TypeExtractor); ok {
				typeofRnode = typed_r_n.GetType()
			} else {

			}

			if !IsTypeEqual(typeofLnode, typeofRnode) {
				fmt.Printf("Additive operation between incompatible types %d:%d\n", p.lex.Line, p.lex.Column)
				os.Exit(1)
			}

			var left_llirVal value.Value

			if valued_l_n, ok := node.(nd.ValueExtractor); ok {
				left_llirVal = valued_l_n.GetLLirValue()
			} else {

			}

			var right_llirVal value.Value

			if valued_r_n, ok := right_n.(nd.ValueExtractor); ok {
				right_llirVal = valued_r_n.GetLLirValue()
			} else {

			}

			if types.IsPointer(right_llirVal.Type()) {
				right_llirVal = p.gen.B.NewLoad(types.I32, right_llirVal)
			}

			if types.IsPointer(left_llirVal.Type()) {
				left_llirVal = p.gen.B.NewLoad(types.I32, left_llirVal)
			}

			var llirAdd value.Value
			if op == "+" {
				llirAdd = p.gen.B.NewAdd(left_llirVal, right_llirVal)
			} else {
				llirAdd = p.gen.B.NewSub(left_llirVal, right_llirVal) // ОТДЕЛЬНЫЙ МЕТОД ГЕНЕРАТОРА
			}

			node = nd.NewAdditiveOpNode(op, node, right_n, typeofRnode, llirAdd)

		}
		for p.currentToken.TokType == constants.OR {
			op := p.currentToken.Text
			p.Eat(constants.OR)

			right_n := p.Term()
			var typeofRnode constants.ItemType

			if typed_r_n, ok := right_n.(nd.TypeExtractor); ok {
				typeofRnode = typed_r_n.GetType()
			} else {

			}

			if !IsTypeEqual(typeofLnode, typeofRnode) {
				fmt.Printf("Boolean operation between incompatible types %d:%d\n", p.lex.Line, p.lex.Column)
				os.Exit(1)
			}

			var left_llirVal value.Value

			if valued_l_n, ok := node.(nd.ValueExtractor); ok {
				left_llirVal = valued_l_n.GetLLirValue()
			} else {

			}

			var right_llirVal value.Value

			if valued_r_n, ok := right_n.(nd.ValueExtractor); ok {
				right_llirVal = valued_r_n.GetLLirValue()
			} else {

			}

			if types.IsPointer(right_llirVal.Type()) {
				right_llirVal = p.gen.B.NewLoad(types.I32, right_llirVal)
			}

			if types.IsPointer(left_llirVal.Type()) {
				left_llirVal = p.gen.B.NewLoad(types.I32, left_llirVal)
			}

			llirOr := p.gen.B.NewOr(left_llirVal, right_llirVal)

			node = nd.NewBooleanOpNode(op, node, right_n, typeofRnode, llirOr)

		}
	}

	return node
}

func (p *Parser) Term() nd.Node {
	node := p.Factor()

	var typeofLnode constants.ItemType

	if typed_l_n, ok := node.(nd.TypeExtractor); ok {
		typeofLnode = typed_l_n.GetType()
	} else {
	}

	for p.currentToken.TokType == constants.MULTIPLICATIVE || p.currentToken.TokType == constants.AND {
		for p.currentToken.TokType == constants.MULTIPLICATIVE {
			op := p.currentToken.Text
			p.Eat(constants.MULTIPLICATIVE)
			right_n := p.Factor()
			var typeofRnode constants.ItemType

			if typed_r_n, ok := right_n.(nd.TypeExtractor); ok {
				typeofRnode = typed_r_n.GetType()
			} else {
			}

			if !IsTypeEqual(typeofLnode, typeofRnode) {
				fmt.Printf("Multiplicative operation between incompatible types %d:%d\n", p.lex.Line, p.lex.Column)
				os.Exit(1)
			}

			var left_llirVal value.Value

			if valued_l_n, ok := node.(nd.ValueExtractor); ok {
				left_llirVal = valued_l_n.GetLLirValue()
			} else {

			}

			var right_llirVal value.Value

			if valued_r_n, ok := right_n.(nd.ValueExtractor); ok {
				right_llirVal = valued_r_n.GetLLirValue()
			} else {

			}

			var llirMul value.Value

			if types.IsPointer(right_llirVal.Type()) {
				right_llirVal = p.gen.B.NewLoad(types.I32, right_llirVal)
			}

			if types.IsPointer(left_llirVal.Type()) {
				left_llirVal = p.gen.B.NewLoad(types.I32, left_llirVal)
			}

			if op == "*" {
				llirMul = p.gen.B.NewMul(left_llirVal, right_llirVal)
			} else {
				llirMul = p.gen.B.NewSDiv(left_llirVal, right_llirVal)
			}

			node = nd.NewMultiplicativeOpNode(op, node, right_n, typeofRnode, llirMul)
		}

		for p.currentToken.TokType == constants.AND {

			op := p.currentToken.Text
			p.Eat(constants.AND)

			right_n := p.Factor()
			var typeofRnode constants.ItemType

			if typed_r_n, ok := right_n.(nd.TypeExtractor); ok {
				typeofRnode = typed_r_n.GetType()
			} else {

			}

			if !IsTypeEqual(typeofLnode, typeofRnode) {
				fmt.Printf("Boolean operation between incompatible types %d:%d\n", p.lex.Line, p.lex.Column)
				os.Exit(1)
			}

			var left_llirVal value.Value

			if valued_l_n, ok := node.(nd.ValueExtractor); ok {
				left_llirVal = valued_l_n.GetLLirValue()
			} else {

			}

			var right_llirVal value.Value

			if valued_r_n, ok := right_n.(nd.ValueExtractor); ok {
				right_llirVal = valued_r_n.GetLLirValue()
			} else {

			}

			if types.IsPointer(right_llirVal.Type()) {
				right_llirVal = p.gen.B.NewLoad(types.I32, right_llirVal)
			}

			if types.IsPointer(left_llirVal.Type()) {
				left_llirVal = p.gen.B.NewLoad(types.I32, left_llirVal)
			}

			llirAnd := p.gen.B.NewAnd(left_llirVal, right_llirVal)

			node = nd.NewBooleanOpNode(op, node, right_n, typeofRnode, llirAnd)

		}
	}

	return node
}

func (p *Parser) Factor() nd.Node {
	switch p.currentToken.TokType {
	case constants.IDENT:

		if _, ok := p.sts[p.currentToken.Text]; ok { //
			func_name := p.currentToken.Text                 //
			p.Eat(constants.IDENT)                           //
			p.Eat(constants.LPAREN)                          //
			var args []interface{}                           //
			var llirArgs []value.Value                       //
			for p.currentToken.TokType != constants.RPAREN { //
				//
				switch p.currentToken.TokType { //
				case constants.IDENT: //
					if !p.Exists(p.currentToken.Text) { //
						fmt.Printf("%s undefined\n", p.currentToken.Text) //
						os.Exit(1)                                        //
					} //
					//
					_, level := p.TableFirstOccurance(p.currentToken.Text) //
					var item st.Table_item                                 //
					var llirVal value.Value                                //
					if level == 0 {                                        //
						item = p.sts["global"][0].Items[p.currentToken.Text] //
					} else { //
						item = p.sts[p.currentFunction][level].Items[p.currentToken.Text] //
						//
					} //
					//
					llirVal = item.LLIRVal               //
					if types.IsPointer(llirVal.Type()) { //
						llirVal = p.gen.B.NewLoad(types.I32, llirVal) //
					} //
					id_n := nd.NewIdentNode(p.currentToken, llirVal, item.T) //
					args = append(args, id_n)                                //
					llirArgs = append(llirArgs, llirVal)                     //
					p.Eat(constants.IDENT)                                   //
					//
				case constants.NUMBER: //
					//
					integer_val, _ := strconv.Atoi(p.currentToken.Text) //
					// ДАННАЯ КОНСТРУКЦИЯ УЖЕ ВСТРЕЧАЕТСЯ В КОДЕ
					n_n := nd.NewNumNode(p.currentToken.Text, constant.NewInt(types.I32, int64(integer_val))) // КАК-ТО ЭТО ОГРОМНОЕ ПОВТОРЕНИЕ НАДО УБРАТЬ
					args = append(args, n_n)                                                                  // ПОКА ХЗ КАК ЭТО ДЕЛАТЬ
					//
					p.Eat(constants.NUMBER) //
					//
				case constants.STRING: //
				} //
				if p.currentToken.TokType != constants.RPAREN { //
					p.Eat(constants.COMMA) //
				} //
			} //
			//
			p.Eat(constants.RPAREN) //
			//
			llir_val_call := p.gen.B.NewCall(p.sts["global"][0].Items[func_name].LLIRVal, llirArgs...) //
			F_call := nd.NewFuncCallNode("FUNC CALL", func_name, args, llir_val_call)                  //
			//
			return F_call //
		} else { //
			if !p.Exists(p.currentToken.Text) { //
				fmt.Printf("%s undefined", p.currentToken.Text) //
				os.Exit(1)                                      //
			} //
			//
			_, level := p.TableFirstOccurance(p.currentToken.Text) //
			var item st.Table_item                                 //
			var llirVal value.Value                                //
			if level == 0 {                                        //
				item = p.sts["global"][0].Items[p.currentToken.Text] //
				llirVal = p.gen.B.NewLoad(types.I32, item.LLIRVal)   //
				//
			} else { //
				item = p.sts[p.currentFunction][level].Items[p.currentToken.Text] //
				llirVal = item.LLIRVal                                            //
			} //
			//
			node := nd.NewIdentNode(p.currentToken, llirVal, item.T) //
			p.Eat(constants.IDENT)                                   //
			return node                                              //
		} //

	case constants.NUMBER:

		integer_num, _ := strconv.Atoi(p.currentToken.Text)

		node := nd.NewNumNode(p.currentToken.Text, constant.NewInt(types.I32, int64(integer_num)))
		p.Eat(constants.NUMBER)
		return node
	case constants.LPAREN:
		p.Eat(constants.LPAREN)
		node := p.Expression()
		p.Eat(constants.RPAREN)
		return node

	default:
		fmt.Printf("%d:%d Forbidden Factor\n", p.lex.Line, p.lex.Column)
		os.Exit(1)
	}

	return nil
}

func (p *Parser) TableFirstOccurance(var_name string) (constants.ItemType, int) {

	for i := p.scope_level; i >= 1; i-- {
		if item, ok := p.sts[p.currentFunction][i].Items[var_name]; ok {
			return item.T, i
		}
	}

	item, ok := p.sts["global"][0].Items[var_name]

	if ok {
		return item.T, 0
	}

	return -1, -1
}

func IsTypeEqual(l constants.ItemType, r constants.ItemType) bool {
	if l == r ||
		(l == constants.Dynamic && r == constants.Integer) ||
		(r == constants.Dynamic && l == constants.Integer) ||
		(r == constants.F && l == constants.Integer) ||
		(r == constants.Integer && l == constants.F) {
		return true
	}

	return false
}

func (p *Parser) ComparativeExpression() nd.Node {
	node := p.ComparativeTerm()

	for p.currentToken.TokType == constants.COMPARATIVE {

		op := p.currentToken.Text

		p.Eat(constants.COMPARATIVE)
		right_comp_n := p.ComparativeTerm()

		var left_llirVal value.Value

		if valued_l_n, ok := node.(nd.ValueExtractor); ok {
			left_llirVal = valued_l_n.GetLLirValue()
		} else {

		}

		var right_llirVal value.Value

		if valued_r_n, ok := right_comp_n.(nd.ValueExtractor); ok {
			right_llirVal = valued_r_n.GetLLirValue()
		} else {

		}

		if types.IsPointer(right_llirVal.Type()) {
			right_llirVal = p.gen.B.NewLoad(types.I32, right_llirVal)
		}

		if types.IsPointer(left_llirVal.Type()) {
			left_llirVal = p.gen.B.NewLoad(types.I32, left_llirVal)
		}
		var predicate enum.IPred

		switch op {
		case ">=":
			predicate = enum.IPredSGE
		case "<=":
			predicate = enum.IPredSLE
		case ">":
			predicate = enum.IPredSGT
		case "<":
			predicate = enum.IPredSLT
		case "==":
			predicate = enum.IPredEQ
		}
		comp_inst := p.gen.B.NewICmp(predicate, left_llirVal, right_llirVal)

		node = nd.NewBooleanOpNode(op, node, right_comp_n, constants.Integer, comp_inst)
	}

	return node
}

func (p *Parser) ComparativeTerm() nd.Node {

	switch p.currentToken.TokType {
	case constants.IDENT:
		node := p.Expression()
		return node
	case constants.NUMBER:
		node := p.Expression()
		return node
	case constants.LPAREN:
		p.Eat(constants.LPAREN)
		node := p.ComparativeExpression()
		p.Eat(constants.RPAREN)
		return node

	default:
		fmt.Printf("%d:%d Forbidden Factor in Comparative Expression", p.lex.Line, p.lex.Column)
		os.Exit(1)
	}

	return nil
}
