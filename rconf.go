package rconf

import (
	"github.com/mitchellh/go-mruby"
)

// BindContext binds ruby methods in the config to variables in Go.
type BindContext struct {
	class *mruby.Class
}

// block binds a config block to a Go function. The function will be called and passed
// a new BindContext which builds a new, anonymous class that will in turn be used as the
// context for the Ruby config block.
func (bc BindContext) Block(attr string, f func(BindContext)) {
	bc.BlockWith(attr, mruby.ArgsNone(), f)
}

func (bc BindContext) BlockWithArg(attr string, f func(BindContext)) {
	bc.BlockWith(attr, mruby.ArgsReq(1), f)
}

func (bc BindContext) BlockWith(attr string, args mruby.ArgSpec, f func(BindContext)) {
	rbMethod := func(mrb *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
		dsl, _ := dslClass.New()
		f(BindContext{class: dsl.SingletonClass()})

		args := mrb.GetArgs()
		block := args[len(args)-1]

		// Run the block in the context of the anonymous class.
		// TODO: Discarding err here
		_, err := dsl.CallBlock("instance_eval", block)
		if err != nil {
			println(err.Error())
			return nil, nil
		}

		return nil, nil
	}
	bc.class.DefineMethod(attr, rbMethod, args|mruby.ArgsBlock())
}

// BindString binds a string value in the config to a go string var.
func (bc BindContext) BindString(attr string, str *string) {
	rbMethod := func(mrb *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
		val := mrb.GetArgs()[0]
		*str = val.String()
		return nil, nil
	}
	bc.class.DefineMethod(attr, rbMethod, mruby.ArgsReq(1))
}

//func (bc BindContext) BindStringFn(attr string, f func(string)) {
//	rbMethod := func(mrb *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
//		val := mrb.GetArgs()[0]
//		f(val.String())
//		return nil, nil
//	}
//	bc.class.DefineMethod(attr, rbMethod, mruby.ArgsReq(1))
//}

func (bc BindContext) BindMapAttr(attr string, m map[string]interface{}, key string) {
	rbMethod := func(mrb *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
		val := mrb.GetArgs()[0]
		m[key] = val.String()
		return nil, nil
	}
	bc.class.DefineMethod(attr, rbMethod, mruby.ArgsReq(1))
}
