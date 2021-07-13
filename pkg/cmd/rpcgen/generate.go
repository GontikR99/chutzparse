package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type FuncArg struct {
	Name string
	Type string
}

type InterfaceMethod struct {
	Name string
	ArgTypes []FuncArg
	RetTypes []FuncArg
}

type Import struct {
	Name string
	Path string
}

type ImportSet []Import

func (is *ImportSet) Add(ispec *ast.ImportSpec) {
	iport := Import{}
	iport.Path = ispec.Path.Value[1:len(ispec.Path.Value)-1]
	if ispec.Name==nil {
		iport.Name=path.Base(iport.Path)
	} else {
		iport.Name=ispec.Name.String()
	}
	*is=append(*is, iport)
}

func (is *ImportSet) Lookup(abbrev string) string {
	for _, importEntry := range *is {
		if importEntry.Name == abbrev {
			return importEntry.Path
		}
	}
	return abbrev
}

func formatType(reqimports map[string]struct{}, expr ast.Expr) string {
	if selExp, ok := expr.(*ast.SelectorExpr); ok {
		pkg := selExp.X.(*ast.Ident).String()
		reqimports[pkg]=struct{}{}
		return pkg+"."+selExp.Sel.String()
	}
	if idExp, ok := expr.(*ast.Ident); ok {
		return idExp.String()
	}
	if starExp, ok := expr.(*ast.StarExpr); ok {
		return "*"+formatType(reqimports, starExp.X)
	}
	if arrExp, ok := expr.(*ast.ArrayType); ok {
		// FIXME: handle fixed length arrays
		return "[]"+formatType(reqimports, arrExp.Elt)
	}
	panic(fmt.Errorf("unknown expression type at %v", expr.Pos()))
}

func upperCamlCase(name string) string {
	if len(name)>0 {
		return strings.ToUpper(name[:1])+name[1:]
	} else {
		return name
	}
}

func main() {
	fset := token.NewFileSet()
	inputFilename := os.Args[1]
	inputBasename := path.Base(inputFilename)
	if strings.HasSuffix(inputBasename, ".go") {
		inputBasename = inputBasename[:len(inputBasename)-3]
	}
	outputBasename := inputBasename+"_genrpc.go"
	outputFilename := path.Join(path.Dir(inputFilename), outputBasename)
	parsedFile, err := parser.ParseFile(fset, inputFilename, nil, 0)
	if err!=nil {
		panic(err)
	}

	imports := new(ImportSet)
	for _, ispec := range parsedFile.Imports {
		imports.Add(ispec)
	}
	neededImports := map[string]struct{}{}

	if err!=nil {
		panic(err)
	}

	ifName:=""
	var ifType *ast.InterfaceType

	for _, decl := range parsedFile.Decls {
		genDec, ok := decl.(*ast.GenDecl)
		if !ok {continue}
		if len(genDec.Specs)<1 {
			continue
		}
		if _, ok = genDec.Specs[0].(*ast.ImportSpec); ok {
			continue
		}
		if _, ok = genDec.Specs[0].(*ast.ValueSpec); ok {
			continue
		}
		if len(genDec.Specs)>1 {
			panic(fmt.Errorf("multiple spec in declaration at %v", genDec.Pos()))
		}
		typeSpec, ok := genDec.Specs[0].(*ast.TypeSpec)
		if !ok {
			panic(fmt.Errorf("expected TypeSpec at %v", genDec.Pos()))
		}
		ifVal, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}
		if ifName=="" {
			ifName = typeSpec.Name.String()
			ifType = ifVal
		} else {
			panic(fmt.Errorf("more than one interface found at %v", genDec.Pos()))
		}
	}

	methods := []*InterfaceMethod{}
	for _, method := range ifType.Methods.List {
		ifMethod := &InterfaceMethod{}
		if len(method.Names)!=1 {
			panic(fmt.Errorf("expected a single name per method at %v", method.Pos()))
		}
		ifMethod.Name = method.Names[0].String()
		fun, ok := method.Type.(*ast.FuncType)
		if !ok {
			panic(fmt.Errorf("expected method to be a function type at %v", method.Pos()))
		}
		for idx, arg := range fun.Params.List {
			fArg := FuncArg{
				Name: fmt.Sprintf("Arg%d", idx),
				Type: formatType(neededImports, arg.Type),
			}
			if len(arg.Names)==1 {
				fArg.Name=upperCamlCase(arg.Names[0].Name)
			}
			ifMethod.ArgTypes = append(ifMethod.ArgTypes, fArg)
		}
		for idx, ret := range fun.Results.List {
			rArg := FuncArg{
				Name: fmt.Sprintf("Ret%d", idx),
				Type: formatType(neededImports, ret.Type),
			}
			if len(ret.Names)==1 {
				rArg.Name=upperCamlCase(ret.Names[0].Name)
			}
			ifMethod.RetTypes = append(ifMethod.RetTypes, rArg)
		}
		if len(ifMethod.RetTypes)==0 || ifMethod.RetTypes[len(ifMethod.RetTypes)-1].Type!="error" {
			panic(fmt.Errorf("Method %s must return 'error' as its last result type at %v", ifMethod.Name, method.Pos()))
		}
		ifMethod.RetTypes = ifMethod.RetTypes[:len(ifMethod.RetTypes)-1]
		methods=append(methods, ifMethod)
	}

	output := &bytes.Buffer{}

	// Write package
	fmt.Fprintf(output, "package %s\n", parsedFile.Name.String())
	fmt.Fprintln(output)

	// Write warning
	fmt.Fprintln(output, "// This code is generated, do not edit!")

	// Write imports
	if len(neededImports)!=0 {
		fmt.Fprintln(output, "import (")
		for impName, _ := range neededImports {
			fmt.Fprintf(output, "  %s \"%s\"\n", impName, imports.Lookup(impName))
		}
		fmt.Fprintln(output, ")")
	}
	fmt.Fprintln(output, "import (golang_net_rpc \"net/rpc\")")
	fmt.Fprintln(output)

	// Write stub declarations
	fmt.Fprintf(output, "type Stub_serverSide_%s struct {impl %s}\n", ifName, ifName)
	fmt.Fprintf(output, "type Stub_clientSide_%s struct {client *golang_net_rpc.Client}\n", ifName)

	// Create stub implementations
	for _, method := range methods {
		fmt.Fprintf(output, "type RequestMessage_%s_%s struct {\n", ifName, method.Name)
		for _, field := range method.ArgTypes {
			fmt.Fprintf(output, "  %s %s\n", field.Name, field.Type)
		}
		fmt.Fprintln(output, "}")

		fmt.Fprintf(output, "type ResponseMessage_%s_%s struct {\n", ifName, method.Name)
		for _, field := range method.RetTypes {
			fmt.Fprintf(output, "  %s %s\n", field.Name, field.Type)
		}
		fmt.Fprintln(output, "}")

		fmt.Fprintf(output, "func (ss Stub_serverSide_%s) %s(req *RequestMessage_%s_%s, res *ResponseMessage_%s_%s) error {\n",
			ifName,
			method.Name,
			ifName, method.Name,
			ifName, method.Name,
		)
		fmt.Fprintln(output, "  var err error")
		output.WriteString("  ")
		for _, field := range method.RetTypes {
			fmt.Fprintf(output, "res.%s, ", field.Name)
		}
		fmt.Fprintf(output, "err = ss.impl.%s(", method.Name)
		needSep := false
		for _, field := range method.ArgTypes {
			if needSep {
				output.WriteString(", ")
			} else {
				needSep = true
			}
			fmt.Fprintf(output, "req.%s", field.Name)
		}
		fmt.Fprintln(output, ")")
		fmt.Fprintln(output, "  return err")
		fmt.Fprintln(output, "}")

		fmt.Fprintf(output, "func (cs Stub_clientSide_%s) %s(", ifName, method.Name)
		needSep = false
		for idx, field := range method.ArgTypes {
			if needSep {
				output.WriteString(", ")
			} else {
				needSep = true
			}
			fmt.Fprintf(output, "arg%d %s", idx, field.Type)
		}
		fmt.Fprintf(output, ") (")
		for _, rField := range method.RetTypes {
			fmt.Fprintf(output, "%s, ", rField.Type)
		}
		fmt.Fprintln(output, "error) {")
		fmt.Fprintf(output, "  req := &RequestMessage_%s_%s{\n", ifName, method.Name)
		for idx, field := range method.ArgTypes {
			fmt.Fprintf(output, "    %s: arg%d,\n", field.Name, idx)
		}
		fmt.Fprintln(output, "  }")
		fmt.Fprintf(output, "  res := new(ResponseMessage_%s_%s)\n", ifName, method.Name)
		fmt.Fprintf(output, "  err := cs.client.Call(\"Stub_serverSide_%s.%s\", req, res)\n", ifName, method.Name)
		fmt.Fprintf(output, "  return ")
		for _, field := range method.RetTypes {
			fmt.Fprintf(output, "res.%s, ", field.Name)
		}
		fmt.Fprintln(output, "err")
		fmt.Fprintln(output, "}")
	}
	fmt.Fprintf(output, `
func New%sClient(client *golang_net_rpc.Client) %s {
    return Stub_clientSide_%s{client}
}

func Handle%s(server %s) func(*golang_net_rpc.Server) {
    ss:=Stub_serverSide_%s{server}
    return func(srv *golang_net_rpc.Server) {srv.Register(ss)}
}
`, ifName, ifName, ifName, ifName, ifName, ifName)
	err = ioutil.WriteFile(outputFilename, output.Bytes(), 0644)
	if err!=nil {
		panic(err)
	}
}