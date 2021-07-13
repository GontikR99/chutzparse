// +build native

// Package rpcgen reads a go file containing a single interface, and generates stubs and helper functions to
// allow golang's net/rpc package to carry the RPC
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func upperCamlCase(name string) string {
	if len(name)>0 {
		return strings.ToUpper(name[:1])+name[1:]
	} else {
		return name
	}
}

func main() {
	inputFilename := os.Args[1]
	inputBasename := path.Base(inputFilename)
	if strings.HasSuffix(inputBasename, ".go") {
		inputBasename = inputBasename[:len(inputBasename)-3]
	}
	outputBasename := inputBasename+"_genrpc.go"
	outputFilename := path.Join(path.Dir(inputFilename), outputBasename)

	compFile, err := ReadGoSource(inputFilename)
	if err!=nil {
		panic(err)
	}

	mainCode := &bytes.Buffer{}
	gc := &GenerationContext{}

	for ifName, methods := range compFile.interfaces {
		// Write stub declarations
		fmt.Fprintf(mainCode, "type Stub_serverSide_%s struct {impl %s}\n", ifName, ifName)
		fmt.Fprintf(mainCode, "type Stub_clientSide_%s struct {client *%s.Client}\n", ifName, gc.ImportName("net/rpc"))

		// Create stub implementations
		for methodName, method := range methods {
			fmt.Fprintf(mainCode, "type RequestMessage_%s_%s struct {\n", ifName, methodName)
			for _, field := range method.ArgTypes {
				fmt.Fprintf(mainCode, "  %s %s\n", upperCamlCase(field.Name), field.TypeText(gc))
			}
			fmt.Fprintln(mainCode, "}")

			fmt.Fprintf(mainCode, "type ResponseMessage_%s_%s struct {\n", ifName, methodName)
			for _, field := range method.RetTypes {
				fmt.Fprintf(mainCode, "  %s %s\n", upperCamlCase(field.Name), field.TypeText(gc))
			}
			fmt.Fprintln(mainCode, "}")

			fmt.Fprintf(mainCode, "func (ss Stub_serverSide_%s) %s(req *RequestMessage_%s_%s, res *ResponseMessage_%s_%s) error {\n",
				ifName,
				methodName,
				ifName, methodName,
				ifName, methodName,
			)
			fmt.Fprintln(mainCode, "  var err error")
			mainCode.WriteString("  ")
			for _, field := range method.RetTypes {
				fmt.Fprintf(mainCode, "res.%s, ", upperCamlCase(field.Name))
			}
			fmt.Fprintf(mainCode, "err = ss.impl.%s(", methodName)
			needSep := false
			for _, field := range method.ArgTypes {
				if needSep {
					mainCode.WriteString(", ")
				} else {
					needSep = true
				}
				fmt.Fprintf(mainCode, "req.%s", upperCamlCase(field.Name))
			}
			fmt.Fprintln(mainCode, ")")
			fmt.Fprintln(mainCode, "  return err")
			fmt.Fprintln(mainCode, "}")

			fmt.Fprintf(mainCode, "func (cs Stub_clientSide_%s) %s(", ifName, methodName)
			needSep = false
			for idx, field := range method.ArgTypes {
				if needSep {
					mainCode.WriteString(", ")
				} else {
					needSep = true
				}
				fmt.Fprintf(mainCode, "arg%d %s", idx, field.TypeText(gc))
			}
			fmt.Fprintf(mainCode, ") (")
			for _, rField := range method.RetTypes {
				fmt.Fprintf(mainCode, "%s, ", rField.TypeText(gc))
			}
			fmt.Fprintln(mainCode, "error) {")
			fmt.Fprintf(mainCode, "  req := &RequestMessage_%s_%s{\n", ifName, methodName)
			for idx, field := range method.ArgTypes {
				fmt.Fprintf(mainCode, "    %s: arg%d,\n", upperCamlCase(field.Name), idx)
			}
			fmt.Fprintln(mainCode, "  }")
			fmt.Fprintf(mainCode, "  res := new(ResponseMessage_%s_%s)\n", ifName, methodName)
			fmt.Fprintf(mainCode, "  err := cs.client.Call(\"Stub_serverSide_%s.%s\", req, res)\n", ifName, methodName)
			fmt.Fprintf(mainCode, "  return ")
			for _, field := range method.RetTypes {
				fmt.Fprintf(mainCode, "res.%s, ", upperCamlCase(field.Name))
			}
			fmt.Fprintln(mainCode, "err")
			fmt.Fprintln(mainCode, "}")
		}
		fmt.Fprintf(mainCode, `
func New%sClient(client *%s.Client) %s {
    return Stub_clientSide_%s{client}
}

func handle%s(server %s) func(*%s.Server) {
    ss:=Stub_serverSide_%s{server}
    return func(srv *%s.Server) {srv.Register(ss)}
}
`, ifName, gc.ImportName("net/rpc"), ifName, ifName,
ifName, ifName, gc.ImportName("net/rpc"), ifName, gc.ImportName("net/rpc"))
	}

	fullBuf := &bytes.Buffer{}

	// Write package
	fmt.Fprintf(fullBuf, "package %s\n", compFile.PackageName())
	fmt.Fprintln(fullBuf)

	// Write warning
	fmt.Fprintln(fullBuf, "// This code is generated, do not edit!")

	// Write imports
	gc.WriteImports(fullBuf)

	// Write code
	fullBuf.Write(mainCode.Bytes())

	err = ioutil.WriteFile(outputFilename, fullBuf.Bytes(), 0644)
	if err!=nil {
		panic(err)
	}
}