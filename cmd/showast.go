package main

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/samber/lo"
	"github.com/zgg-lang/zgg-go/ast"

	"github.com/zgg-lang/zgg-go/parser"
)

const INDENT = "   "

type (
	astFieldValue interface {
	}
	astAttr struct {
		Field string
		Value astFieldValue
	}
	astTreeNode struct {
		Name  string
		Attrs []astAttr
		Child any
	}
)

func runShowAstExpr(isDebug bool, args []string) {
	code := args[0]
	node, errs := parser.ParseReplFromString(code, !isDebug)
	if errs != nil {
		for i, e := range errs {
			fmt.Printf("err[%d]: %v\n", i, e)
		}
		return
	}
	tree := getAstNodeTree(node)
	fmt.Println(tree)
}

func getAstNodeTree(root any) (tnode *astTreeNode) {
	tnode = &astTreeNode{}
	switch rootVal := root.(type) {
	case ast.Expr:
		switch vv := rootVal.(type) {
		case *ast.ExprInt:
			tnode.Name = "INT"
			tnode.Child = vv.Value.Value()
			return
		case *ast.ExprFloat:
			tnode.Name = "FLOAT"
			tnode.Child = vv.Value.Value()
			return
		case *ast.ExprBool:
			tnode.Name = "BOOL"
			tnode.Child = vv.Value.Value()
			return
		case *ast.ExprStr:
			tnode.Name = "STR"
			tnode.Child = vv.Value.Value()
			return
		case *ast.ExprArray:
			tnode.Name = "ARRAY"
			tnode.Child = lo.Map(vv.Items, func(item *ast.ArrayItem, _ int) *astTreeNode {
				expr := getAstNodeTree(item.Expr)
				return expr
			})
			return
		case ast.Expr:
		default:
			tnode.Name = "OTHER"
			tnode.Child = vv
			return
		}
	case ast.CallArgument:
		if rootVal.Arg != nil {
			tnode = getAstNodeTree(rootVal.Arg)
			return
		}
	}
	typ := reflect.TypeOf(root)
	val := reflect.ValueOf(root)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		panic(fmt.Sprintf("showAstTree: unexpected type %+v", typ))
	}
	tnode.Name = typ.String()
	for i := 0; i < typ.NumField(); i++ {
		ft := typ.Field(i)
		fv := val.Field(i)
		switch ft.Name {
		case "BinOp":
			tnode.Attrs = append(tnode.Attrs,
				getAstFieldNode("Left", fv.FieldByName("Left")),
				getAstFieldNode("Right", fv.FieldByName("Right")),
			)
			continue
		case "Pos":
			continue
		}
		tnode.Attrs = append(tnode.Attrs, getAstFieldNode(ft.Name, fv))
	}
	return
}

func getAstFieldNode(label string, value reflect.Value) (fnode astAttr) {
	fnode.Field = label
	vi := value.Interface()
	switch vv := vi.(type) {
	case ast.Expr:
		fnode.Value = getAstNodeTree(vv)
	default:
		vvv := reflect.ValueOf(vv)
		if vvv.Kind() == reflect.Slice {
			values := make([]*astTreeNode, 0, vvv.Len())
			for i := 0; i < vvv.Len(); i++ {
				item := vvv.Index(i).Interface()
				values = append(values, getAstNodeTree(item))
			}
			fnode.Value = values
		} else {
			fnode.Value = &astTreeNode{Name: "OTHER", Child: vv}
		}
	}
	return
}

func (n *astTreeNode) String() string {
	var buf strings.Builder
	n.textTo("", &buf)
	return buf.String()
}

func (n *astTreeNode) textTo(indent string, w io.Writer) {
	if n.Name == "OTHER" {
		fmt.Fprint(w, n.Child)
		return
	}
	if len(n.Attrs) == 0 {
		switch child := n.Child.(type) {
		case []*astTreeNode:
			fmt.Fprintf(w, "%s(", n.Name)
			nindent := indent + INDENT
			for _, item := range child {
				fmt.Fprint(w, "\n"+nindent)
				item.textTo(nindent, w)
			}
			fmt.Fprintf(w, "\n%s)", indent)
		default:
			s, _ := json.Marshal(child)
			fmt.Fprintf(w, "%s(%s)", n.Name, string(s))
		}
		return
	}
	nextIndent := indent + INDENT
	fmt.Fprintf(w, "%s(\n", n.Name)
	for i := range n.Attrs {
		attr := &n.Attrs[i]
		fmt.Fprintf(w, "%s%s = ", nextIndent, attr.Field)
		switch vv := attr.Value.(type) {
		case *astTreeNode:
			vv.textTo(nextIndent, w)
		case []*astTreeNode:
			fmt.Fprintln(w, "[")
			nnindent := nextIndent + INDENT
			for _, nn := range vv {
				fmt.Fprint(w, nnindent)
				nn.textTo(nnindent, w)
				fmt.Fprintln(w)
			}
			fmt.Fprintf(w, "%s]", nextIndent)
		default:
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintf(w, "%s)", indent)
}
