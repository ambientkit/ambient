package lib

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

var (
	ErrVarMustString = errors.New("variable type must be a string")
	ErrVarNotFound   = errors.New("variable not found")
)

type Tree struct {
	FileSet    *token.FileSet
	File       *ast.File
	CommentMap ast.CommentMap
}

type ConstVisitor struct {
	Name  string
	Value string
	Err   error
}

type StructChanger struct {
	Name     string
	Fields   []*ast.Field
	Comments *ast.CommentGroup
	Ok       bool
}

type StructCopier struct {
	Name     string
	Fields   []*ast.Field
	Comments *ast.CommentGroup
	Ok       bool
}

type CommentVisitor struct {
	Tree *Tree
}

// Visit walks the tree for CommentVisitor
func (v *CommentVisitor) Visit(n ast.Node) (w ast.Visitor) {
	var s string
	switch x := n.(type) {
	case *ast.BasicLit:
		//s = x.Value
	case *ast.Ident:
		//s = x.Name
	case *ast.StructType:
		//s = string()
		//fmt.Printf("%s %s\n", x)
	case *ast.Comment:
		s = x.Text
	default:
		//fmt.Printf("Missed: %v\n", n)
	}
	if s != "" {
		fmt.Printf("%s=\t%s\n", v.Tree.FileSet.Position(n.Pos()), s)
	}

	return v
}

// PrintComments prints all the comments and locations
func (gt *Tree) PrintComments() error {
	cv := &CommentVisitor{
		Tree: gt,
	}

	ast.Walk(cv, gt.File)

	return nil
}

// Visit walks the tree for StructCopier
func (v *StructCopier) Visit(n ast.Node) (w ast.Visitor) {
	switch spec := n.(type) {
	case *ast.TypeSpec:
		if spec.Name.String() == v.Name {
			switch structType := spec.Type.(type) {
			case *ast.StructType:
				v.Fields = structType.Fields.List
				v.Comments = spec.Comment
				v.Ok = true
				return nil
			}
		}
	}

	return v
}

// StructFields returns the struct fields
func (gt *Tree) StructFields(varName string) ([]*ast.Field, *ast.CommentGroup, error) {
	sc := &StructCopier{
		Name: varName,
	}

	ast.Walk(sc, gt.File)

	if !sc.Ok {
		return nil, nil, ErrVarNotFound
	}

	return sc.Fields, sc.Comments, nil
}

// DiscoverType builds clean fields for ast
func DiscoverType(express ast.Expr) ast.Expr {
	switch t := express.(type) {
	case *ast.Ident:
		return &ast.Ident{
			Name: t.Name,
		}
	case *ast.StarExpr:
		return &ast.StarExpr{
			X: DiscoverType(t.X),
		}
	case *ast.ArrayType:
		return &ast.ArrayType{
			Elt: DiscoverType(t.Elt),
		}
	case *ast.MapType:
		return &ast.MapType{
			Key:   t.Key,
			Value: t.Value,
		}
	case *ast.InterfaceType:
		return &ast.InterfaceType{
			Methods: t.Methods,
		}
	case *ast.SelectorExpr:
		return &ast.SelectorExpr{
			X: DiscoverType(t.X),
			Sel: &ast.Ident{
				Name: t.Sel.Name,
			},
		}
	case *ast.StructType:
		return &ast.StructType{
			Fields: &ast.FieldList{
				List: DiscoverList(t.Fields.List),
			},
		}
	case *ast.ChanType:
		return &ast.ChanType{
			Arrow: t.Arrow,
			Dir:   t.Dir,
			Value: DiscoverType(t.Value),
		}
	case *ast.FuncType:
		return &ast.FuncType{
			Params: &ast.FieldList{
				List: DiscoverList(t.Params.List),
			},
			Results: &ast.FieldList{
				List: DiscoverList(t.Results.List),
			},
		}
	default:
		fmt.Printf("Missed Type: %#v\n", express)
	}
	return nil
}

// DiscoverList builds clean lists for ast
func DiscoverList(list []*ast.Field) []*ast.Field {
	var l []*ast.Field

	for i := 0; i < len(list); i++ {
		field := &ast.Field{}

		// Comments
		if list[i].Comment != nil {

			// Create a comment group
			field.Comment = &ast.CommentGroup{}

			// Loop through the comments
			for c := 0; c < len(list[i].Comment.List); c++ {
				field.Comment.List = append(field.Comment.List, &ast.Comment{
					Text: list[i].Comment.List[c].Text,
				})
			}

		}

		// Loop through the names
		for n := 0; n < len(list[i].Names); n++ {
			field.Names = append(field.Names, &ast.Ident{
				Name: list[i].Names[n].Name,
			})
		}

		// Tags
		if list[i].Tag != nil {
			field.Tag = &ast.BasicLit{
				Kind:  list[i].Tag.Kind,
				Value: list[i].Tag.Value,
			}
		}

		// Types
		field.Type = DiscoverType(list[i].Type)

		l = append(l, field)
	}

	return l
}

// Visit walks the tree for StructChanger
func (v *StructChanger) Visit(n ast.Node) (w ast.Visitor) {
	switch spec := n.(type) {
	case *ast.TypeSpec:
		if spec.Name.String() == v.Name {
			switch structType := spec.Type.(type) {
			case *ast.StructType:
				structType.Fields.List = DiscoverList(v.Fields)
				spec.Comment = v.Comments
				v.Ok = true
				return nil
			}
		}
	}
	return v
}

// ChangeStruct changes the struct field list
func (gt *Tree) ChangeStruct(varName string, varValue []*ast.Field, varComments *ast.CommentGroup) error {
	sc := &StructChanger{
		Name:     varName,
		Fields:   varValue,
		Comments: varComments,
	}

	ast.Walk(sc, gt.File)

	if !sc.Ok {
		return ErrVarNotFound
	}

	return nil
}

// Visit walks the tree for ConstVisitor
func (v *ConstVisitor) Visit(n ast.Node) (w ast.Visitor) {
	switch spec := n.(type) {
	case *ast.ValueSpec:
		for i := 0; i < len(spec.Names); i++ {
			if spec.Names[i].Name == v.Name {
				switch val := spec.Values[i].(type) {
				case *ast.BasicLit:
					if val.Kind == token.STRING {
						spec.Values[i] = &ast.BasicLit{
							Value: strconv.Quote(v.Value),
						}
					} else {
						v.Err = ErrVarMustString
					}
				}
				return nil
			}
		}

	}
	return v
}

// ChangeConstString changes the value of a named const
func (gt *Tree) ChangeConstString(varName string, varValue string) error {
	cv := &ConstVisitor{
		Name:  varName,
		Value: varValue,
	}

	ast.Walk(cv, gt.File)

	if cv.Err != nil {
		return cv.Err
	}

	return nil
}

// New creates a new Go package tree
func New(name string) *Tree {
	gt := &Tree{
		FileSet: token.NewFileSet(),
		File: &ast.File{
			Name: &ast.Ident{
				Name: name,
			},
		},
	}

	return gt
}

// Load creates a Go package tree from a file
func Load(filepath string, mode parser.Mode) (*Tree, error) {
	var err error

	gt := &Tree{
		FileSet: token.NewFileSet(),
	}

	gt.File, err = parser.ParseFile(gt.FileSet, filepath, nil, mode)
	if err != nil {
		return nil, err
	}
	return gt, nil
}

// SetPackageName sets the package name
func (gt *Tree) SetPackageName(name string) {
	gt.File.Name = &ast.Ident{
		Name: name,
	}
}

// AddImport adds an import section with import paths
func (gt *Tree) AddImportSection(imports []string) {
	var specs []ast.Spec

	// Add all the imports
	for i := 0; i < len(imports); i++ {
		specs = append(specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Value: strconv.Quote(imports[i]),
			},
		})
	}

	// Generate the import declaration
	decl := &ast.GenDecl{
		Lparen: 1,
		Rparen: 1,
		Tok:    token.IMPORT,
		Specs:  specs,
	}
	gt.File.Decls = append(gt.File.Decls, decl)
}

// AddHelloMainFunc adds a main func the outputs: hello world
func (gt *Tree) AddHelloMainFunc() {
	fd := &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "main",
		},
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.Ident{
							Name: "fmt.Println",
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: strconv.Quote("hello world"),
							},
						},
					},
				},
			},
		},
	}

	gt.File.Decls = append(gt.File.Decls, fd)
}

// AddImport adds an import
func (gt *Tree) AddImport(imp string) {

	// Add the import
	for i := 0; i < len(gt.File.Decls); i++ {
		switch decl := gt.File.Decls[i].(type) {
		case *ast.GenDecl:
			if decl.Tok == token.IMPORT {
				iSpec := &ast.ImportSpec{
					Path: &ast.BasicLit{
						Value: strconv.Quote(imp),
					}}
				decl.Specs = append(decl.Specs, iSpec)
			}
		}
	}

	// Sort the imports (removes any lines between imports)
	ast.SortImports(gt.FileSet, gt.File)
}

// Bytes returns the code as a byte array
func (gt *Tree) Bytes(doFormat bool) ([]byte, error) {
	var output []byte
	var err error

	buffer := bytes.NewBuffer(output)
	if err = printer.Fprint(buffer, gt.FileSet, gt.File); err != nil {
		return nil, err
	}

	output = buffer.Bytes()

	if doFormat {
		// Format the source
		output, err = format.Source(output)
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}

// WriteFile writes the code to a file and create the folder structure
func (gt *Tree) WriteFile(filepath string, doFormat bool, dirPerm os.FileMode, filePerm os.FileMode) error {
	var err error

	// Create folders
	err = os.MkdirAll(path.Dir(filepath), dirPerm)
	if err != nil {
		return err
	}

	// Generate the code
	data, err := gt.Bytes(doFormat)
	if err != nil {
		return err
	}

	// Write file
	err = ioutil.WriteFile(filepath, data, filePerm)
	if err != nil {
		return err
	}

	return nil
}
