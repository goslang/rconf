package rconf

import (
	//"fmt"
	"io"
	"io/ioutil"

	"github.com/mitchellh/go-mruby"
	"github.com/pkg/errors"
)

var (
	mrb      = mruby.NewMrb()
	dslClass = mrb.DefineClass("RConfDSL", nil)
)

type Interpreter func(io.ReadCloser) error

func NewInterpreter(fn func(bc BindContext)) (Interpreter, error) {
	dsl, err := dslClass.New()
	if err != nil {
		return nil, errors.Wrap(err, "Error instantiating Ruby Class")
	}

	// Declare the DSL
	klass := dsl.SingletonClass()
	bc := BindContext{class: klass}
	fn(bc)

	parser := mruby.NewParser(mrb)

	interpeter := func(reader io.ReadCloser) error {
		defer reader.Close()

		rawRuby, err := ioutil.ReadAll(reader)
		if err != nil {
			return errors.Wrap(err, "Failed to read input")
		}

		parser.Parse(string(rawRuby), nil)
		defer parser.Close()

		// Now we can use parser.GenerateCode to create a Ruby Block and
		// execute it in the context of `dsl`.
		proc := parser.GenerateCode()
		_, err = mrb.Run(proc, dsl)
		if err != nil {
			return errors.Wrap(err, "Failed to executed Ruby")
		}

		return nil
	}

	return interpeter, nil
}
