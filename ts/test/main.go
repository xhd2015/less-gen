package main

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/xhd2015/less-gen/go/project"

	go2ts_types "github.com/xhd2015/less-gen/go/go2ts/types"
	"github.com/xhd2015/less-gen/go/gofile"
	"github.com/xhd2015/less-gen/go/gofmt"
	load "github.com/xhd2015/less-gen/go/load/legacy"
	"github.com/xhd2015/less-gen/template"
	"github.com/xhd2015/less-gen/ts/format"
)

const help = `
PROG help to parse

Usage: Prog x [OPTIONS]
Options:
  --help   show help message
`

func main() {
	err := handle(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
func handle(args []string) error {
	var some string

	var remainArgs []string
	n := len(args)
	for i := 0; i < n; i++ {
		if args[i] == "--some" {
			if i+1 >= n {
				return fmt.Errorf("%v requires arg", args[i])
			}
			some = args[i+1]
			i++
			continue
		}
		if args[i] == "--help" {
			fmt.Println(strings.TrimSpace(help))
			return nil
		}
		if args[i] == "--" {
			remainArgs = append(remainArgs, args[i+1:]...)
			break
		}
		if strings.HasPrefix(args[i], "-") {
			return fmt.Errorf("unrecognized flag: %v", args[i])
		}
		remainArgs = append(remainArgs, args[i])
	}
	// TODO handle
	_ = some

	pkgPath := "github.com/xhd2015/lifelog/handle/learning/sentence"
	dir := filepath.Join(os.Getenv("X"), "lifelog/server_go/handle/learning/sentence")
	fnName := "Generate"
	api := "/api/TODO"

	return generateCode(dir, pkgPath, fnName, api)
}

// every specific go file has a naturely parseable
// or at least part of it is parseable.
func generateCode(dir string, pkgPath string, fnName string, api string) error {
	project, err := project.Load([]string{pkgPath}, &load.LoadOptions{
		Dir: dir,
	})
	if err != nil {
		return err
	}

	pkg := project.GetPkg(pkgPath)
	if pkg == nil {
		return fmt.Errorf("pkg not found: %s", pkgPath)
	}

	gen := pkg.Lookup(fnName)
	if gen == nil {
		return fmt.Errorf("not found: %s", fnName)
	}

	fmt.Printf("gen: %T %v\n", gen.GoObject(), gen.GoObject())

	fn := gen.GoObject().(*types.Func)

	sig := fn.Type().(*types.Signature)

	req := sig.Params().At(1)
	res := sig.Results()
	err = go2ts_types.ValidateRespErr(res)
	if err != nil {
		return err
	}

	var __RESULT__ string
	var __RESP__ string
	if res.Len() > 1 {
		__RESULT__ = "data, "
		__RESP__ = "\"data\": data, "
	}

	reqTypeNamed := req.Type().(*types.Pointer).Elem().(*types.Named)
	reqType := reqTypeNamed.Underlying().(*types.Struct)

	typesInfo := pkg.TypesInfo()

	reqTypeName := reqTypeNamed.Obj().Name()
	pkgName := pkg.GoPkg().Name

	_ = typesInfo

	var reqDef string

	flattenFields := go2ts_types.GetStructFields(reqType)
	var hasInt64 bool
	var hasUserID bool
	for _, flattenField := range flattenFields {
		if len(flattenFields) > 0 {
			continue
		}
		if flattenField.Field.Name == "UserID" {
			hasUserID = true
			continue
		}
		hasInt64 = hasInt64 || gofile.Int64.Equals(flattenField.Field.Type)
	}

	var argToReq string
	var reqName string = "req"

	fctx := gofile.NewContext()
	if hasInt64 {
		var argToReqInitLines []string
		var argToReqParseI64Lines []string
		reqName = "freq"
		argToReqInitLines = append(argToReqInitLines, "var err error")
		argToReqInitLines = append(argToReqInitLines, fmt.Sprintf("freq := %s.%s{", pkgName, reqTypeName))
		var fields []*gofile.Field
		for _, flattenField := range flattenFields {
			if len(flattenField.Nested) > 0 {
				continue
			}
			field := flattenField.Field
			if gofile.Int64.Equals(flattenField.Field.Type) {
				v := *field
				v.Type = &gofile.Named{
					PkgPath: "github.com/xhd2015/lifelog/model",
					Name:    "OptionalNumber",
				}
				field = &v
				argToReqParseI64Lines = append(argToReqParseI64Lines,
					fmt.Sprintf("freq.%s, err = req.%s.Int64()", field.Name, field.Name),
					"if err!=nil{",
					"routehelp.AbortWithErr(ctx, err)",
					"return",
					"}",
				)
			} else {
				argToReqInitLines = append(argToReqInitLines, fmt.Sprintf("%s: %s,", field.Name, field.Name))
			}
			fields = append(fields, field)
		}
		argToReqInitLines = append(argToReqInitLines, "}")
		strct := gofile.Struct{Fields: fields}

		reqDef = strct.Format(fctx)
		argToReq = strings.Join(argToReqInitLines, "\n") + "\n" + strings.Join(argToReqParseI64Lines, "\n")
	} else {
		reqDef = pkgName + "." + reqTypeName
	}

	_ = req

	goFile := &gofile.File{PkgName: "sentence"}

	goFile.Import(&gofile.Import{Path: "net/http"})
	goFile.Import(&gofile.Import{Path: "github.com/gin-gonic/gin"})
	goFile.Import(&gofile.Import{Path: "github.com/xhd2015/lifelog/route/routehelp"})
	goFile.Import(&gofile.Import{Path: "github.com/xhd2015/lifelog/service/session"})

	goFile.Decls = append(goFile.Decls)

	// some user_id conversion
	tpl := `package __ROUTE_PKG_NAME__

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"__PKG__"
	"github.com/xhd2015/lifelog/route/routehelp"
	"github.com/xhd2015/lifelog/service/session"
)

func RouteGenerate(r *gin.Engine) {
	r.Any("__API__", func(ctx *gin.Context) {
		var req __REQ_DEF__
		if !routehelp.ParseRequest(ctx, &req) {
			return
		}
		__ROUTE_ARG_TO_REQ_ARG__
		__USER_ID__
		__RESULT__err := __PKG_NAME__.__FUNC_NAME__(ctx, &__REQ_NAME__)
		if err != nil {
			routehelp.AbortWithErr(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			__RESP__
		})
	})
}
`
	var userIDLine string
	if hasUserID {
		userIDLine = fmt.Sprintf("session := session.Get(ctx)\n%s.UserID = session.UserID", reqName)
	}

	//
	// go
	//   model
	//   route
	// ts
	//   model.ts
	//   route.ts
	//
	code := template.Format(tpl, map[string]string{
		"__API__":                  api,
		"__ROUTE_PKG_NAME__":       "sentence",
		"__PKG__":                  pkgPath,
		"__PKG_NAME__":             pkgName,
		"__REQ_TYPE_NAME__":        reqTypeName,
		"__REQ_NAME__":             reqName,
		"__FUNC_NAME__":            "Generate",
		"__REQ_DEF__":              reqDef,
		"__USER_ID__":              userIDLine,
		"__RESULT__":               __RESULT__,
		"__RESP__":                 __RESP__,
		"__ROUTE_ARG_TO_REQ_ARG__": argToReq,
	})
	fmtCode := gofmt.TryFormatCode(code)
	fmt.Println(fmtCode)

	ctx := gofile.NewContext()
	tsCode := genTSTypes(ctx, flattenFields)
	tsCodePretty, err := format.Pretty(tsCode)
	if err != nil {
		return err
	}

	fmt.Println(tsCodePretty)
	return nil
}

func genTSTypes(ctx gofile.Context, fields []*gofile.StructField) string {
	var extends []string
	var tsFields []string
	for _, field := range fields {
		tsType := toTSType(field.Field.Type)
		if field.IsNested() {
			extends = append(extends, tsType)
			continue
		}
		jsonName := getJSONName(field.Field.Name, reflect.StructTag(field.Field.Tag).Get("json"))
		if jsonName == "" {
			continue
		}
		tsFields = append(tsFields, fmt.Sprintf("%s:%s", jsonName, tsType))
	}
	var ext string
	if len(extends) > 0 {
		ext = " extends " + strings.Join(extends, ",")
	}
	return fmt.Sprintf("export interface X%s{\n%s\n}", ext, strings.Join(tsFields, "\n"))
}

func toTSType(t gofile.Type) string {
	switch t := t.(type) {
	case gofile.BuiltinType:
		switch t {
		case gofile.Bool:
			return "bool"
		case gofile.Int, gofile.Int64:
			return "number"
		case gofile.String:
			return "string"
		default:
			return "any"
		}
	case *gofile.Named:
		return t.Name
	default:
		return "any"
	}
}

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
