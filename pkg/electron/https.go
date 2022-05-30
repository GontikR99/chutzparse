//go:build wasm && electron
// +build wasm,electron

package electron

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gontikr99/chutzparse/pkg/jsbinding"
	"github.com/gontikr99/chutzparse/pkg/nodejs"
	"syscall/js"
)

var http = nodejs.Require("http")
var https = nodejs.Require("https")

func HttpCall(scheme string, method string, hostname string, port int16, path string, headers map[string]string, reqText []byte) (resBody []byte, statCode int, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("Error fetching from site: %v", r)
			}
		}
	}()
	var schemelib js.Value
	if scheme == "http" {
		schemelib = http
	} else if scheme == "https" {
		schemelib = https
	} else {
		return nil, 500, errors.New("Unsupported scheme " + scheme)
	}
	intHeaders := map[string]interface{}{
		//"Content-Type":   "application/json",
		//"Accept":         "application/json",
		"Content-Length": len(reqText),
	}
	if headers != nil {
		for k, v := range headers {
			intHeaders[k] = v
		}
	}
	options := map[string]interface{}{
		"hostname": hostname,
		"port":     port,
		"path":     path,
		"method":   method,
		"headers":  intHeaders,
	}

	doneChan := make(chan struct{})
	buffer := new(bytes.Buffer)
	errHolder := new(error)
	statusCodeHolder := new(int)

	responseFunc := new(js.Func)
	*responseFunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		responseFunc.Release()
		res := args[0]
		*statusCodeHolder = res.Get("statusCode").Int()

		onDataFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			chunk := args[0]
			buffer.Write(jsbinding.ReadArrayBuffer(chunk))
			return nil
		})
		onEndFunc := new(js.Func)
		onErrorFunc := new(js.Func)
		*onEndFunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			onDataFunc.Release()
			onEndFunc.Release()
			onErrorFunc.Release()
			doneChan <- struct{}{}
			return nil
		})
		*onErrorFunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			onDataFunc.Release()
			onEndFunc.Release()
			onErrorFunc.Release()
			*errHolder = errors.New(args[0].String())
			return nil
		})
		res.Call("on", "data", onDataFunc)
		res.Call("on", "error", *onErrorFunc)
		res.Call("on", "end", *onEndFunc)

		return nil
	})
	// FIXME: how to catch the error?
	ucHandle := RegisterUncaughtException(func(err error) {
		*errHolder = err
		doneChan <- struct{}{}
	})
	req := schemelib.Call("request", options, *responseFunc)
	req.Call("write", jsbinding.BufferOf(reqText))
	req.Call("end")

	<-doneChan
	ucHandle.Release()
	if *errHolder != nil {
		return nil, 0, *errHolder
	} else {
		return buffer.Bytes(), *statusCodeHolder, nil
	}
}
