package basic

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"strings"

	"github.com/xhd2015/less-gen/strcase"
	"golang.org/x/tools/go/packages"
)

func processPkg(fset *token.FileSet, pkg *packages.Package) string {
	c := &Conv{
		fset:      fset,
		pkg:       pkg,
		typePkg:   pkg.Types,
		typesInfo: pkg.TypesInfo,
	}

	var codes []string
	for _, file := range pkg.Syntax {
		code := c.file(file)
		codes = append(codes, code)
	}

	jointCode := strings.Join(codes, "\n")

	if pkg.Name == "main" {
		jointCode = jointCode + "\nmain()"
	}
	return jointCode
}

type Conv struct {
	fset      *token.FileSet
	pkg       *packages.Package
	typePkg   *types.Package
	typesInfo *types.Info
}

func (c *Conv) file(f *ast.File) string {
	return c.decls(f.Decls)
}
func (c *Conv) decls(decls []ast.Decl) string {
	var declCode []string
	for _, decl := range decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			declCode = append(declCode, c.genDecl(d))
		case *ast.FuncDecl:
			declCode = append(declCode, c.funcDecl(d))
		default:
			fmt.Printf("decl %T %v\n", decl, decl)
		}
	}
	return strings.Join(declCode, "\n")
}

func (c *Conv) genDecl(g *ast.GenDecl) string {
	switch g.Tok {
	case token.IMPORT:
		for _, spec := range g.Specs {
			c.spec(spec)
		}
	default:
		fmt.Printf("token %s\n", g.Tok)
	}
	return ""
}

func (c *Conv) spec(s ast.Spec) string {
	switch s := s.(type) {
	case *ast.ImportSpec:
		return c.importSpec(s)

	default:
		fmt.Printf("spec: %T %v\n", s, s)
	}
	return ""
}

func (c *Conv) importSpec(s *ast.ImportSpec) string {
	fmt.Printf("import %s\n", s.Path.Value)
	return ""
}

func (c *Conv) funcDecl(f *ast.FuncDecl) string {
	var name string
	if f.Name != nil {
		name = f.Name.Name
	}

	exported := isExported(name)

	var modifier string
	if exported {
		modifier = "export "
	}
	var body string
	if f.Body != nil {
		body = c.blockStmt(f.Body)
	} else {
		body = "throw \"implementation not found\""
	}
	return fmt.Sprintf("%sfunction %s(){\n%s\n}", modifier, name, body)
}

func (c *Conv) blockStmt(f *ast.BlockStmt) string {
	stmts := make([]string, 0, len(f.List))
	for _, stmt := range f.List {
		stmts = append(stmts, c.stmt(stmt))
	}
	return strings.Join(stmts, "\n")
}

func (c *Conv) stmt(f ast.Stmt) string {
	switch f := f.(type) {
	case *ast.ExprStmt:
		return c.exprStmt(f)
	default:
		fmt.Printf("stmt: %T %v\n", f, f)
	}
	return ""
}

func (c *Conv) exprStmt(f *ast.ExprStmt) string {
	return c.expr(f.X)
}

func (c *Conv) expr(f ast.Expr) string {
	switch f := f.(type) {
	case *ast.CallExpr:
		return c.callExpr(f)
	case *ast.BasicLit:
		return c.basicLit(f)
	default:
		fmt.Printf("expr: %T %v\n", f, f)
	}
	return ""
}

func (c *Conv) basicLit(f *ast.BasicLit) string {
	switch f.Kind {
	case token.STRING:
		return f.Value
	default:
		return f.Value
	}
}
func (c *Conv) callExpr(f *ast.CallExpr) string {
	var fn string
	if f.Fun != nil {
		var def types.Object
		switch t := f.Fun.(type) {
		case *ast.Ident:
			def = c.typesInfo.Uses[t]
		case *ast.SelectorExpr:
			def = c.typesInfo.Uses[t.Sel]
		}
		if def == nil {
			fmt.Printf("def not found\n")
		} else {
			fmt.Printf("def found: pkg=%s, name=%s\n", def.Pkg().Path(), def.Name())
		}

		if def.Pkg().Path() == "fmt" && (def.Name() == "Printf" || def.Name() == "Println") {
			fn = "console.log"
		} else {
			fn = c.expr(f.Fun)
		}
	}

	var exprs []string
	for _, arg := range f.Args {
		expr := c.expr(arg)
		exprs = append(exprs, expr)
	}
	return fmt.Sprintf("%s(%s)", fn, strings.Join(exprs, ","))
}

func isExported(name string) bool {
	return name != "" && strings.ToUpper(name[0:1]) == name[0:1]
}

func (c *Conv) Translate(object types.Object) string {
	fmt.Printf("obj: %s %T\n", object.Name(), object)
	objType := object.Type()
	switch o := object.(type) {
	case *types.TypeName:
		t := o.Type()
		switch t := t.(type) {
		case *types.Named:
			ut := t.Underlying()
			if st, ok := ut.(*types.Struct); ok {
				stLit := c.translateStruct(st)
				var modifier string
				if t.Obj().Exported() {
					modifier = "export "
				}
				return fmt.Sprintf("%sinterface %s %s", modifier, t.Obj().Name(), stLit)
			}
		}
	case *types.Func:
		return c.translateFunc(o)
		// case *types.Signature:
		// return fmt.Sprintf()
	}
	return fmt.Sprintf("TODO %T %v", objType, objType)
}

func (c *Conv) translateStruct(st *types.Struct) string {
	var fields []string
	n := st.NumFields()
	for i := 0; i < n; i++ {
		stf := st.Field(i)

		if stf.Anonymous() {
			// ut := stf.Type().Underlying()
			panic("TODO ut")
		}
		tag := st.Tag(i)

		stTag := reflect.StructTag(tag)
		name := stf.Name()

		jsonName := getJSONName(name, stTag.Get("json"))
		typStr := fmt.Sprint(stf.Type())

		fields = append(fields, fmt.Sprintf("    %s: %s", jsonName, typStr))
	}
	return fmt.Sprintf("{\n%s\n}", strings.Join(fields, "\n"))
}
func (c *Conv) translateFunc(o *types.Func) string {
	name := o.Name()
	var modifier string
	if o.Exported() {
		modifier = "export "
	}
	t := c.transalteType(o.Type())

	scope := c.transalteScope(o.Scope())

	jsName := strcase.Decapitalize(name)
	return fmt.Sprintf("%sfunction %s(%s){%s}", modifier, jsName, scope, t)
}

func (c *Conv) transalteType(t types.Type) string {
	switch t := t.(type) {
	case *types.Signature:
		return c.translateSignature(t)
	}
	return fmt.Sprintf("TODO type %T %s", t, t)
}

func (c *Conv) translateSignature(t *types.Signature) string {
	var list []string
	params := t.Params()
	n := params.Len()
	for i := 0; i < n; i++ {
		v := params.At(i)

		jsType := v.Type().String()
		list = append(list, fmt.Sprintf("%s: %s", v.Name(), jsType))
	}
	return strings.Join(list, ",")
}

func (c *Conv) transalteScope(t *types.Scope) string {
	n := t.NumChildren()
	for i := 0; i < n; i++ {
		child := t.Child(i)
		_ = child
	}
	return "TODO"
}

// xx.go

// xx.go
func getJSONName(fieldName string, jsonTag string) string {
	jsonName := parseJSONTagName(jsonTag)
	if jsonName != "" {
		if jsonName == "-" {
			return ""
		}
		return jsonName
	}
	return fieldName
}

// parseJSONTagName get json name that will appear in marshaled json
func parseJSONTagName(jsonTag string) string {
	before, _, _ := strings.Cut(jsonTag, ",")

	return before
}
