package grapher

import (
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"

	h "compiler/helpers"
	nd "compiler/types/nodes"
)

func (g *Grapher) Graph_tree(graph *graphviz.Graph) {
	root := g.AST_Tree
	g.recursive_traverse(graph, root)
}

type Grapher struct {
	node_id  string
	edge_id  string
	AST_Tree nd.RootNode
}

func NewGrapher(AST nd.RootNode) *Grapher {
	return &Grapher{"0", "0", AST}
}

func (g *Grapher) recursive_traverse(graph *graphviz.Graph, tree_node interface{}) *cgraph.Node {
	switch n := tree_node.(type) {
	case nd.IdentNode:
		node, _ := graph.CreateNodeByName(n.Name + g.node_id)
		node.SetLabel(n.Name)
		g.node_id = h.String_plus(g.node_id)
		return node

	case nd.NumNode:
		node, _ := graph.CreateNodeByName(n.Number + g.node_id)
		node.SetLabel(n.Number)
		g.node_id = h.String_plus(g.node_id)
		return node

	case nd.RootNode:
		node, _ := graph.CreateNodeByName(n.Prog_name + g.node_id)
		node.SetLabel(n.Prog_name)
		g.node_id = h.String_plus(g.node_id)
		left := g.recursive_traverse(graph, n.Var_n)
		graph.CreateEdgeByName(g.edge_id, node, left)
		g.edge_id = h.String_plus(g.edge_id)
		right := g.recursive_traverse(graph, n.Operator_n)
		graph.CreateEdgeByName(g.edge_id, node, right)
		g.edge_id = h.String_plus(g.edge_id)
		return node

	case nd.BooleanOpNode:

		node, _ := graph.CreateNodeByName(n.Op + g.node_id)
		node.SetLabel(n.Op)
		g.node_id = h.String_plus(g.node_id)

		l := g.recursive_traverse(graph, n.Left)
		graph.CreateEdgeByName(g.edge_id, node, l)
		g.edge_id = h.String_plus(g.edge_id)
		r := g.recursive_traverse(graph, n.Right)
		graph.CreateEdgeByName(g.edge_id, node, r)
		g.edge_id = h.String_plus(g.edge_id)
		return node
	case nd.AdditiveOpNode:
		node, _ := graph.CreateNodeByName(n.Op + g.node_id)
		node.SetLabel(n.Op)
		g.node_id = h.String_plus(g.node_id)
		l := g.recursive_traverse(graph, n.Left)
		graph.CreateEdgeByName(g.edge_id, node, l)
		g.edge_id = h.String_plus(g.edge_id)
		r := g.recursive_traverse(graph, n.Right)
		graph.CreateEdgeByName(g.edge_id, node, r)
		g.edge_id = h.String_plus(g.edge_id)
		return node
	case nd.MultiplicativeOpNode:
		node, _ := graph.CreateNodeByName(n.Op + g.node_id)
		node.SetLabel(n.Op)
		g.node_id = h.String_plus(g.node_id)
		l := g.recursive_traverse(graph, n.Left)
		graph.CreateEdgeByName(g.edge_id, node, l)
		g.edge_id = h.String_plus(g.edge_id)
		r := g.recursive_traverse(graph, n.Right)
		graph.CreateEdgeByName(g.edge_id, node, r)
		g.edge_id = h.String_plus(g.edge_id)
		return node
	case nd.AssignOpNode:
		node, _ := graph.CreateNodeByName(n.Op + g.node_id)
		node.SetLabel(n.Op)
		g.node_id = h.String_plus(g.node_id)
		l := g.recursive_traverse(graph, n.Left)
		graph.CreateEdgeByName(g.edge_id, node, l)
		g.edge_id = h.String_plus(g.edge_id)
		r := g.recursive_traverse(graph, n.Right)
		graph.CreateEdgeByName(g.edge_id, node, r)
		g.edge_id = h.String_plus(g.edge_id)
		return node
	case nd.RangeClause:
		node, _ := graph.CreateNodeByName("RANGE" + g.node_id)
		node.SetLabel("RANGE")
		g.node_id = h.String_plus(g.node_id)
		left := g.recursive_traverse(graph, n.From)
		right := g.recursive_traverse(graph, n.To)

		graph.CreateEdgeByName(g.edge_id, node, left)
		g.edge_id = h.String_plus(g.edge_id)

		graph.CreateEdgeByName(g.edge_id, node, right)
		g.edge_id = h.String_plus(g.edge_id)
		return node

	case nd.ComparisonNode:

		node, _ := graph.CreateNodeByName("CMP" + g.node_id)
		node.SetLabel("CMP")
		g.node_id = h.String_plus(g.node_id)

		left := g.recursive_traverse(graph, n.Left)
		operator := g.recursive_traverse(graph, n.Op)
		right := g.recursive_traverse(graph, n.Right)

		graph.CreateEdgeByName(g.edge_id, node, left)
		g.edge_id = h.String_plus(g.edge_id)

		graph.CreateEdgeByName(g.edge_id, node, operator)
		g.edge_id = h.String_plus(g.edge_id)

		graph.CreateEdgeByName(g.edge_id, node, right)
		g.edge_id = h.String_plus(g.edge_id)

		return node

	case nd.IfStatementNode:

		node, _ := graph.CreateNodeByName("IF" + g.node_id)
		node.SetLabel("IF")
		g.node_id = h.String_plus(g.node_id)

		condition := g.recursive_traverse(graph, n.Stmnt)
		graph.CreateEdgeByName(g.edge_id, node, condition)
		g.edge_id = h.String_plus(g.edge_id)

		if_b, _ := graph.CreateNodeByName("IF BODY" + g.node_id)
		if_b.SetLabel("IF BODY")
		g.node_id = h.String_plus(g.node_id)
		graph.CreateEdgeByName(g.edge_id, node, if_b)
		g.edge_id = h.String_plus(g.edge_id)

		else_b, _ := graph.CreateNodeByName("ELSE BODY" + g.node_id)
		else_b.SetLabel("ELSE BODY")
		g.node_id = h.String_plus(g.node_id)
		graph.CreateEdgeByName(g.edge_id, node, else_b)
		g.edge_id = h.String_plus(g.edge_id)

		for i, _ := range n.Body {
			returned_node := g.recursive_traverse(graph, n.Body[i])
			graph.CreateEdgeByName(g.edge_id, if_b, returned_node)
			g.edge_id = h.String_plus(g.edge_id)
		}

		for i, _ := range n.Else_body {
			returned_node := g.recursive_traverse(graph, n.Else_body[i])
			graph.CreateEdgeByName(g.edge_id, else_b, returned_node)
			g.edge_id = h.String_plus(g.edge_id)
		}

		return node

	case nd.ForNode:
		node, _ := graph.CreateNodeByName("FOR LOOP" + g.node_id)
		node.SetLabel("FOR LOOP")
		g.node_id = h.String_plus(g.node_id)

		ran := g.recursive_traverse(graph, n.R)
		graph.CreateEdgeByName(g.edge_id, node, ran)
		g.edge_id = h.String_plus(g.edge_id)

		l := g.recursive_traverse(graph, n.Operator_n)
		graph.CreateEdgeByName(g.edge_id, node, l)
		g.edge_id = h.String_plus(g.edge_id)

		return node

	case nd.OperatorNode:
		node, _ := graph.CreateNodeByName(n.Text + g.node_id)
		node.SetLabel(n.Text)
		g.node_id = h.String_plus(g.node_id)

		for i, _ := range n.Operators {
			returned_node := g.recursive_traverse(graph, n.Operators[i])
			graph.CreateEdgeByName(g.edge_id, node, returned_node)
			g.edge_id = h.String_plus(g.edge_id)
		}
		return node

	case nd.FuncNode:

		node, _ := graph.CreateNodeByName(n.Text + g.node_id)
		node.SetLabel(n.Text)
		g.node_id = h.String_plus(g.node_id)

		name_n, _ := graph.CreateNodeByName(n.Name + g.node_id)
		name_n.SetLabel(n.Name)
		g.node_id = h.String_plus(g.node_id)

		name_flat, _ := graph.CreateNodeByName("NAME" + g.node_id)
		name_flat.SetLabel("NAME")
		g.node_id = h.String_plus(g.node_id)

		graph.CreateEdgeByName(g.edge_id, node, name_flat)
		g.edge_id = h.String_plus(g.edge_id)
		graph.CreateEdgeByName(g.edge_id, name_flat, name_n)
		g.edge_id = h.String_plus(g.edge_id)

		args_flat, _ := graph.CreateNodeByName("ARGS" + g.node_id)
		args_flat.SetLabel("ARGS")
		g.node_id = h.String_plus(g.node_id)

		for i, _ := range n.Args {
			returned_node := g.recursive_traverse(graph, n.Args[i])
			graph.CreateEdgeByName(g.edge_id, args_flat, returned_node)
			g.edge_id = h.String_plus(g.edge_id)
		}

		graph.CreateEdgeByName(g.edge_id, node, args_flat)
		g.edge_id = h.String_plus(g.edge_id)

		body_n, _ := graph.CreateNodeByName(n.Body.Text + g.node_id)
		body_n.SetLabel(n.Body.Text)
		g.node_id = h.String_plus(g.node_id)

		for i, _ := range n.Body.Operators {
			returned_node := g.recursive_traverse(graph, n.Body.Operators[i])
			graph.CreateEdgeByName(g.edge_id, body_n, returned_node)
			g.edge_id = h.String_plus(g.edge_id)
		}

		graph.CreateEdgeByName(g.edge_id, node, body_n)
		g.edge_id = h.String_plus(g.edge_id)

		return node

	case nd.VarNode:
		node, _ := graph.CreateNodeByName(n.Text + g.node_id)
		node.SetLabel(n.Text)
		g.node_id = h.String_plus(g.node_id)

		for i, _ := range n.Ids {
			returned_id := g.recursive_traverse(graph, n.Ids[i])
			graph.CreateEdgeByName(g.edge_id, node, returned_id)
			g.edge_id = h.String_plus(g.edge_id)
		}

		return node

	case nd.ReadOpNode:
		node, _ := graph.CreateNodeByName(n.Text + g.node_id)
		node.SetLabel(n.Text)
		g.node_id = h.String_plus(g.node_id)
		left := g.recursive_traverse(graph, n.Id)
		graph.CreateEdgeByName(g.edge_id, node, left)
		g.edge_id = h.String_plus(g.edge_id)
		return node

	case nd.WriteOpNode:
		node, _ := graph.CreateNodeByName(n.Text + g.node_id)
		node.SetLabel(n.Text)
		g.node_id = h.String_plus(g.node_id)
		left := g.recursive_traverse(graph, n.To_write)
		graph.CreateEdgeByName(g.edge_id, node, left)
		g.edge_id = h.String_plus(g.edge_id)
		return node

	case nd.FuncCallNode:
		node, _ := graph.CreateNodeByName(n.Text + g.node_id)
		node.SetLabel(n.Text)
		g.node_id = h.String_plus(g.node_id)

		name, _ := graph.CreateNodeByName("NAME" + g.node_id)
		name.SetLabel("NAME")
		g.node_id = h.String_plus(g.node_id)

		graph.CreateEdgeByName(g.edge_id, node, name)
		g.edge_id = h.String_plus(g.edge_id)

		func_name_n, _ := graph.CreateNodeByName(n.FuncName + g.node_id)
		func_name_n.SetLabel(n.FuncName)
		g.node_id = h.String_plus(g.node_id)

		graph.CreateEdgeByName(g.edge_id, name, func_name_n)
		g.edge_id = h.String_plus(g.edge_id)

		args, _ := graph.CreateNodeByName("ARGS" + g.node_id)
		args.SetLabel("ARGS")
		g.node_id = h.String_plus(g.node_id)

		graph.CreateEdgeByName(g.edge_id, node, args)
		g.edge_id = h.String_plus(g.edge_id)

		for i, _ := range n.Args {
			returned_id := g.recursive_traverse(graph, n.Args[i])
			graph.CreateEdgeByName(g.edge_id, args, returned_id)
			g.edge_id = h.String_plus(g.edge_id)
		}

		return node
	case string:
		node, _ := graph.CreateNodeByName("string" + g.node_id)
		node.SetLabel(n)
		g.node_id = h.String_plus(g.node_id)
		return node

	case nd.ReturnNode:
		node, _ := graph.CreateNodeByName("RETURN" + g.node_id)
		node.SetLabel("RETURN")
		g.node_id = h.String_plus(g.node_id)

		to_return := g.recursive_traverse(graph, n.ToReturn)
		graph.CreateEdgeByName(g.edge_id, node, to_return)
		g.edge_id = h.String_plus(g.edge_id)
		return node
	// case []interface{}:

	// 	for _, item := range n {
	// 		g.recursive_traverse(graph, item)
	// 	}

	default:
		return nil
	}

}
