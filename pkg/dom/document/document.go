// +build wasm,web

package document

import (
	"github.com/gontikr99/chutzparse/pkg/dom"
	"syscall/js"
)

var doc = js.Global().Get("document")

func QuerySelector(selector string) dom.Element {
	jsv := doc.Call("querySelector", selector)
	if jsv.IsNull() || jsv.IsUndefined() {
		return nil
	} else {
		return dom.WrapElement(jsv)
	}
}

func GetElementById(id string) dom.Element {
	jsv := doc.Call("getElementById", id)
	if jsv.IsNull() || jsv.IsUndefined() {
		return nil
	} else {
		return dom.WrapElement(jsv)
	}
}

func GetElementsByTagName(tag string) []dom.Element {
	jsv := doc.Call("getElementsByTagName", tag)
	var result []dom.Element
	for i := 0; i < jsv.Length(); i++ {
		result = append(result, dom.WrapElement(jsv.Index(i)))
	}
	return result
}

func CreateElement(tag string) dom.Element {
	return dom.WrapElement(doc.Call("createElement", tag))
}

func CreateElementNS(namespace string, name string) dom.Element {
	return dom.WrapElement(doc.Call("createElementNS", namespace, name))
}

func CreateTextNode(content string) dom.Element {
	return dom.WrapElement(doc.Call("createTextNode", content))
}

func ExecCommand(command string) {
	doc.Call("execCommand", command)
}
