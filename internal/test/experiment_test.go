package test

import "testing"

type tContextSlice struct {
	Counter int
	Funcs   []func(sl *tContextSlice) error
}

func (t *tContextSlice) addTrigger(f func(st *tContextSlice) error) {
	t.Funcs = append(t.Funcs, f)
}

func (t *tContextSlice) invokeTrigger() (err error) {
	for i := 0; i < len(t.Funcs); i++ {
		if f := t.Funcs[i]; f != nil {
			if err = f(t); err != nil {
				return
			}
		}
	}
	return
}

func step2(st *tContextSlice) error {
	st.Counter++
	return nil
}

func BenchmarkThunkCallSlice(b *testing.B) {
	ctx := &tContextSlice{}
	ctx.addTrigger(step2)
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = ctx.invokeTrigger()
	}
	b.StopTimer()
	_ = err
	println(ctx.Counter)
}
func BenchmarkThunkCallSliceWithClosure(b *testing.B) {
	ctx := &tContextSlice{}
	ctx.addTrigger(func(st *tContextSlice) error {
		st.Counter++
		return nil
	})
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = ctx.invokeTrigger()
	}
	b.StopTimer()
	_ = err
	println(ctx.Counter)
}

type tContext struct {
	Counter int
	trigger mnext
}

func (t *tContext) addTrigger(f mnext) {
	if t.trigger == nil {
		t.trigger = f
	} else {
		prev := t.trigger
		t.trigger = func(ctx *tContext, next func() error) error {
			nextFunc := func() error {
				return prev(ctx, next)
			}
			return f(ctx, nextFunc)
		}
	}
}

func (t *tContext) noop() error {
	return nil
}

func (t *tContext) invokeTrigger() error {
	return t.trigger(t, t.noop)
}

type mnext func(ctx *tContext, next func() error) error

func step1(t *tContext, next func() error) error {
	t.Counter++
	return next()
}

func BenchmarkThunkCall(b *testing.B) {
	if pprofX {
		defer profCpu()()
	}
	ctx := &tContext{}
	ctx.addTrigger(step1)
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = ctx.invokeTrigger()
	}
	b.StopTimer()
	_ = err
	println(ctx.Counter)
}

func BenchmarkThunkCallWithClosure(b *testing.B) {
	if pprofX {
		defer profCpu()()
	}
	ctx := &tContext{}
	next := func(ctx *tContext, next func() error) error {
		ctx.Counter++
		return next()
	}
	ctx.addTrigger(next)
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = ctx.invokeTrigger()
	}
	b.StopTimer()
	_ = err
	println(ctx.Counter)
}
