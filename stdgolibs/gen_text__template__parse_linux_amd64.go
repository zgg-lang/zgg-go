package stdgolibs

import (
	pkg "text/template/parse"

	"reflect"
)

func init() {
	registerValues("text/template/parse", map[string]reflect.Value{
		// Functions
		"Parse":         reflect.ValueOf(pkg.Parse),
		"New":           reflect.ValueOf(pkg.New),
		"IsEmptyTree":   reflect.ValueOf(pkg.IsEmptyTree),
		"NewIdentifier": reflect.ValueOf(pkg.NewIdentifier),

		// Consts

		"ParseComments":  reflect.ValueOf(pkg.ParseComments),
		"NodeText":       reflect.ValueOf(pkg.NodeText),
		"NodeAction":     reflect.ValueOf(pkg.NodeAction),
		"NodeBool":       reflect.ValueOf(pkg.NodeBool),
		"NodeChain":      reflect.ValueOf(pkg.NodeChain),
		"NodeCommand":    reflect.ValueOf(pkg.NodeCommand),
		"NodeDot":        reflect.ValueOf(pkg.NodeDot),
		"NodeField":      reflect.ValueOf(pkg.NodeField),
		"NodeIdentifier": reflect.ValueOf(pkg.NodeIdentifier),
		"NodeIf":         reflect.ValueOf(pkg.NodeIf),
		"NodeList":       reflect.ValueOf(pkg.NodeList),
		"NodeNil":        reflect.ValueOf(pkg.NodeNil),
		"NodeNumber":     reflect.ValueOf(pkg.NodeNumber),
		"NodePipe":       reflect.ValueOf(pkg.NodePipe),
		"NodeRange":      reflect.ValueOf(pkg.NodeRange),
		"NodeString":     reflect.ValueOf(pkg.NodeString),
		"NodeTemplate":   reflect.ValueOf(pkg.NodeTemplate),
		"NodeVariable":   reflect.ValueOf(pkg.NodeVariable),
		"NodeWith":       reflect.ValueOf(pkg.NodeWith),
		"NodeComment":    reflect.ValueOf(pkg.NodeComment),

		// Variables

	})
	registerTypes("text/template/parse", map[string]reflect.Type{
		// Non interfaces

		"Tree":           reflect.TypeOf((*pkg.Tree)(nil)).Elem(),
		"Mode":           reflect.TypeOf((*pkg.Mode)(nil)).Elem(),
		"NodeType":       reflect.TypeOf((*pkg.NodeType)(nil)).Elem(),
		"Pos":            reflect.TypeOf((*pkg.Pos)(nil)).Elem(),
		"ListNode":       reflect.TypeOf((*pkg.ListNode)(nil)).Elem(),
		"TextNode":       reflect.TypeOf((*pkg.TextNode)(nil)).Elem(),
		"CommentNode":    reflect.TypeOf((*pkg.CommentNode)(nil)).Elem(),
		"PipeNode":       reflect.TypeOf((*pkg.PipeNode)(nil)).Elem(),
		"ActionNode":     reflect.TypeOf((*pkg.ActionNode)(nil)).Elem(),
		"CommandNode":    reflect.TypeOf((*pkg.CommandNode)(nil)).Elem(),
		"IdentifierNode": reflect.TypeOf((*pkg.IdentifierNode)(nil)).Elem(),
		"VariableNode":   reflect.TypeOf((*pkg.VariableNode)(nil)).Elem(),
		"DotNode":        reflect.TypeOf((*pkg.DotNode)(nil)).Elem(),
		"NilNode":        reflect.TypeOf((*pkg.NilNode)(nil)).Elem(),
		"FieldNode":      reflect.TypeOf((*pkg.FieldNode)(nil)).Elem(),
		"ChainNode":      reflect.TypeOf((*pkg.ChainNode)(nil)).Elem(),
		"BoolNode":       reflect.TypeOf((*pkg.BoolNode)(nil)).Elem(),
		"NumberNode":     reflect.TypeOf((*pkg.NumberNode)(nil)).Elem(),
		"StringNode":     reflect.TypeOf((*pkg.StringNode)(nil)).Elem(),
		"BranchNode":     reflect.TypeOf((*pkg.BranchNode)(nil)).Elem(),
		"IfNode":         reflect.TypeOf((*pkg.IfNode)(nil)).Elem(),
		"RangeNode":      reflect.TypeOf((*pkg.RangeNode)(nil)).Elem(),
		"WithNode":       reflect.TypeOf((*pkg.WithNode)(nil)).Elem(),
		"TemplateNode":   reflect.TypeOf((*pkg.TemplateNode)(nil)).Elem(),
	})
}
