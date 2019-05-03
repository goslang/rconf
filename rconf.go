package rconf

import (
	"github.com/mitchellh/go-mruby"
)

// BindContext binds ruby methods in the config to variables in Go.
type BindContext struct {
	class *mruby.Class
}

// Block binds a config block to a Go function. The function will be called
// and passed a new BindContext which builds a new, anonymous class that will
// in turn be used as the context for the Ruby config block.
func (bc BindContext) Block(attr string, f func(BindContext)) {
	bc.BlockWith(attr, mruby.ArgsNone(), f)
}

func (bc BindContext) BlockWithArg(attr string, f func(BindContext)) {
	bc.BlockWith(attr, mruby.ArgsReq(1), f)
}

func (
	bc BindContext,
) BlockWith(
	attr string,
	args mruby.ArgSpec,
	f func(BindContext),
) {
	rbMethod := func(mrb *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
		// Create a new instance of the DSL class and attach any new DSL methods
		// as Singletons.
		dsl, err := dslClass.New()
		if err != nil {
			return nil, err.(*mruby.Exception).MrbValue
		}

		f(BindContext{class: dsl.SingletonClass()})

		// Get the passed block, it will be the last argument in the list.
		// NOTE: We should do some more sanity checking here to make sure that the
		// block actually exists.
		args := mrb.GetArgs()
		block := args[len(args)-1]

		// Now that the DSL methods are defined, execute the block in the context
		// of the DSL instance.
		// TODO: Effectively discarding err here
		_, err = dsl.CallBlock("instance_eval", block)
		if err != nil {
			// NOTE: This leaves open the possibility of a panic and should be
			// fixed.
			return nil, err.(*mruby.Exception).MrbValue
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

// BindInt binds a Ruby method named `attr` to a target integer.
func (bc BindContext) BindInt(attr string, target *int) {
	rbMethod := func(mrb *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
		val := mrb.GetArgs()[0]
		*target = val.Fixnum()
		return nil, nil
	}
	bc.class.DefineMethod(attr, rbMethod, mruby.ArgsReq(1))
}

// BindFloat binds a Ruby method named `attr` to a target float64.
func (bc BindContext) BindFloat(attr string, target *float64) {
	rbMethod := func(mrb *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
		val := mrb.GetArgs()[0]
		*target = val.Float()
		return nil, nil
	}
	bc.class.DefineMethod(attr, rbMethod, mruby.ArgsReq(1))
}

// BindMapAttr binds a Ruby method named `attr` to the key, `key`, in the
// provided map.
func (
	bc BindContext,
) BindMapAttr(
	attr string,
	m map[string]interface{},
	key string,
) {
	rbMethod := func(mrb *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
		val := mrb.GetArgs()[0]
		m[key] = val.String()
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

func (bc BindContext) Arg(idx int) (*mruby.MrbValue, bool) {
	args := mrb.GetArgs()

	if len(args) <= idx {
		return nil, false
	}
	return args[0], true
}

func (bc BindContext) StringArg(idx int) string {
	arg, ok := bc.Arg(idx)

	if !ok {
		return ""
	}
	return arg.String()
}
