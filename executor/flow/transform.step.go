package flow

import (
	"fmt"

	v8 "rogchap.com/v8go"
)

type TransformStep struct {
	transformFunction string // Js function to transform the data
}

func NewTransformStep(transformFunction string) *TransformStep {
	return &TransformStep{
		transformFunction: transformFunction,
	}
}

func (s *TransformStep) Run(context *FlowContext) error {
	iso := v8.NewIsolate() // create a new VM
	// a template that represents a JS function
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		fmt.Printf("%v", info.Args())
		return nil
	})
	global := v8.NewObjectTemplate(iso)
	global.Set("print", printfn)
	ctx := v8.NewContext(iso, global)
	_, err := ctx.RunScript(s.transformFunction, "transform.js")
	if err != nil {
		return err
	}

	// get the function
	fn, err := ctx.RunScript("transform", "transform.js")
	if err != nil {
		return err
	}

	fu, err := fn.AsFunction()
	if err != nil {
		return err
	}

	// create Object
	obj := v8.NewObjectTemplate(iso)
	for k, v := range context.Data {
		if v == nil {
			continue
		}

		if s, ok := v.(string); ok {
			obj.Set(k, s)
		}

		if i, ok := v.(int); ok {
			obj.Set(k, i)
		}

		if b, ok := v.(bool); ok {
			obj.Set(k, b)
		}

		if f, ok := v.(float64); ok {
			obj.Set(k, f)
		}

		if f, ok := v.(float32); ok {
			obj.Set(k, f)
		}

		if f, ok := v.(map[string]interface{}); ok {
			obj.Set(k, f)
			for k, v := range f {
				obj.Set(k, v)
			}
		}
	}

	out, err := fu.Call(ctx.Global())
	if err != nil {
		return err
	}

	json, err := out.MarshalJSON()
	fmt.Println(string(json))
	return nil
}

// Write a function that takes a Go Object and converts it to a JS Object
