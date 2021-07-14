// +build native

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"path"
)

type GenerationContext struct {
	NeededImports []Import
}

type Import struct {
	Name string
	Path string
}

// ImportName returns the shorthand name we'll use in the generated file to refer to a specific import path
func (gc *GenerationContext) ImportName(impPath string) string {
	for _, imRec := range gc.NeededImports {
		if imRec.Path == impPath {
			return imRec.Name
		}
	}
	base := path.Base(impPath)
	idx := 1
	for {
		storeName := base
		if idx != 1 {
			storeName = fmt.Sprintf("%s%d", base, idx)
		}
		free := true
		for _, imRec := range gc.NeededImports {
			if imRec.Name == base {
				free = false
				break
			}
		}
		if free {
			gc.NeededImports = append(gc.NeededImports, Import{
				Name: storeName,
				Path: impPath,
			})
			return storeName
		}
	}
}

func (gc *GenerationContext) WriteImports(out io.Writer) error {
	_, err := fmt.Fprintln(out, "import (")
	if err != nil {
		return err
	}

	for _, entry := range gc.NeededImports {
		if path.Base(entry.Path) == entry.Name {
			_, err = fmt.Fprintf(out, "    \"%s\"\n", entry.Path)
			if err != nil {
				return err
			}
		} else {
			_, err = fmt.Fprintf(out, "    %s \"%s\"\n", entry.Name, entry.Path)
			if err != nil {
				return err
			}
		}
	}
	_, err = fmt.Fprintln(out, ")")
	return err
}

type ComprehendedSource struct {
	fileSet    *token.FileSet
	parsedFile *ast.File
	importMap  map[string]string
	interfaces map[string]Interface
}

func ReadGoSource(fileName string) (*ComprehendedSource, error) {
	result := &ComprehendedSource{
		fileSet:    token.NewFileSet(),
		importMap:  map[string]string{},
		interfaces: map[string]Interface{},
	}
	var err error
	result.parsedFile, err = parser.ParseFile(result.fileSet, fileName, nil, 0)
	if err != nil {
		return nil, err
	}

	// parse imports
	for _, ispec := range result.parsedFile.Imports {
		imPath := ispec.Path.Value[1 : len(ispec.Path.Value)-1]
		imName := path.Base(imPath)
		if ispec.Name == nil {
			imName = ispec.Name.String()
		}
		result.importMap[imName] = imPath
	}

	// parse interfaces
	for _, decl := range result.parsedFile.Decls {
		genDec, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if len(genDec.Specs) < 1 {
			continue
		}
		if _, ok = genDec.Specs[0].(*ast.ImportSpec); ok {
			continue
		}
		if _, ok = genDec.Specs[0].(*ast.ValueSpec); ok {
			continue
		}
		if len(genDec.Specs) > 1 {
			return nil, fmt.Errorf("multiple spec in declaration at %v", genDec.Pos())
		}
		typeSpec, ok := genDec.Specs[0].(*ast.TypeSpec)
		if !ok {
			return nil, fmt.Errorf("expected TypeSpec at %v", genDec.Pos())
		}
		ifType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}
		ifName := typeSpec.Name.Name

		methods := Interface{}
		for _, method := range ifType.Methods.List {
			ifMethod := InterfaceMethod{}
			if len(method.Names) != 1 {
				return nil, fmt.Errorf("expected a single name per method at %v", method.Pos())
			}
			methodName := method.Names[0].String()
			fun, ok := method.Type.(*ast.FuncType)
			if !ok {
				return nil, fmt.Errorf("expected method to be a function type at %v", method.Pos())
			}
			for idx, arg := range fun.Params.List {
				fArg := FuncParameter{
					Name:       fmt.Sprintf("arg%d", idx),
					Type:       arg.Type,
					SourceFile: result,
				}
				if len(arg.Names) == 1 {
					fArg.Name = arg.Names[0].Name
				}
				ifMethod.ArgTypes = append(ifMethod.ArgTypes, fArg)
			}
			for idx, ret := range fun.Results.List {
				rArg := FuncParameter{
					Name:       fmt.Sprintf("ret%d", idx),
					Type:       ret.Type,
					SourceFile: result,
				}
				if len(ret.Names) == 1 {
					rArg.Name = ret.Names[0].Name
				}
				ifMethod.RetTypes = append(ifMethod.RetTypes, rArg)
			}
			ifMethod.RetTypes = ifMethod.RetTypes[:len(ifMethod.RetTypes)-1]
			methods[methodName] = ifMethod
		}
		result.interfaces[ifName] = methods
	}
	return result, nil
}

func (cs *ComprehendedSource) PackageName() string {
	return cs.parsedFile.Name.Name
}

type FuncParameter struct {
	Name       string
	Type       ast.Expr
	SourceFile *ComprehendedSource
}

func formatType(gc *GenerationContext, cs *ComprehendedSource, expr ast.Expr) string {
	if idExp, ok := expr.(*ast.Ident); ok {
		return idExp.String()
	}
	if starExp, ok := expr.(*ast.StarExpr); ok {
		return "*" + formatType(gc, cs, starExp.X)
	}
	if arrExp, ok := expr.(*ast.ArrayType); ok {
		// FIXME: handle arrays in addition to slices
		return "[]" + formatType(gc, cs, arrExp.Elt)
	}
	if mapExp, ok := expr.(*ast.MapType); ok {
		return "map[" + formatType(gc, cs, mapExp.Key) + "]" + formatType(gc, cs, mapExp.Value)
	}
	if selExp, ok := expr.(*ast.SelectorExpr); ok {
		pkg := selExp.X.(*ast.Ident).Name
		pkgPath := cs.importMap[pkg]
		impName := gc.ImportName(pkgPath)
		return impName + "." + selExp.Sel.Name
	}
	// FIXME: handle function types (*ast.FuncType), even though gob and net/rpc can't transport them?
	panic(fmt.Errorf("unknown expression type at %v", expr.Pos()))
}

func (fp FuncParameter) TypeText(gc *GenerationContext) string {
	return formatType(gc, fp.SourceFile, fp.Type)
}

type InterfaceMethod struct {
	ArgTypes []FuncParameter
	RetTypes []FuncParameter
}

type Interface map[string]InterfaceMethod

type ImportSet []Import

func (is *ImportSet) Add(ispec *ast.ImportSpec) {
	iport := Import{}
	iport.Path = ispec.Path.Value[1 : len(ispec.Path.Value)-1]
	if ispec.Name == nil {
		iport.Name = path.Base(iport.Path)
	} else {
		iport.Name = ispec.Name.String()
	}
	*is = append(*is, iport)
}

func (is *ImportSet) Lookup(abbrev string) string {
	for _, importEntry := range *is {
		if importEntry.Name == abbrev {
			return importEntry.Path
		}
	}
	return abbrev
}
