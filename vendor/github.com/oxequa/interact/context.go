package interact

import (
	"errors"
	"io"
)

type (
	Context interface {
		Skip()
		Reload()
		GetReload() int
		SetPrfx(io.Writer, interface{})
		SetDef(interface{}, interface{})
		SetErr(interface{})
		SetEnd(string)
		Ans() Cast
		Def() Cast
		Err() error
		Parent() Context
		Prfx() Cast
		Qns() Qns
		Tag() string
		Quest() string
	}
	context struct {
		i *Interact
		q *Question
	}
)

type ErrorFunc func(Context) error

type BoolFunc func(Context) bool

type InterfaceFunc func(Context) interface{}

func (c *context) Parent() Context {
	if c.q.parent != nil {
		return &context{i: c.i, q: c.q.parent}
	}
	return &context{i: c.i}
}

func (c *context) Ans() Cast {
	if c.q != nil {
		cast := &cast{answer: c.q.response, value: c.q.value}
		return cast
	}
	return &cast{}
}

func (c *context) Skip() {
	if c.q != nil {
		c.q.skip = true
	}
	c.i.skip = true
}

func (c *context) Reload() {
	if c.q != nil {
		c.q.reload.status = true
		c.q.reload.loop++
	}
}

func (c *context) Err() error {
	if c.q.err != nil {
		return errors.New(c.q.err.(string))
	} else if c.i.err != nil {
		return errors.New(c.i.err.(string))
	}
	return nil
}

func (c *context) Def() Cast {
	if c.q != nil {
		return &cast{value: c.q.def.Value}
	}
	return &cast{value: c.i.def.Value}
}

func (c *context) Prfx() Cast {
	if c.q != nil && c.q.prefix.Text != nil {
		return &cast{value: c.q.prefix.Text}
	}
	return &cast{value: c.i.prefix.Text}
}

func (c *context) Qns() Qns {
	list := []Context{}
	if c.q != nil {
		for _, q := range c.q.Subs {
			list = append(list, &context{i: c.i, q: q})
		}
	} else {
		for _, q := range c.i.Questions {
			list = append(list, &context{i: c.i, q: q})
		}
	}
	return &qns{list: list}
}

func (c *context) Tag() string {
	if c.q != nil {
		return c.q.Tag
	}
	return ""
}

func (c *context) Quest() string {
	if c.q != nil {
		return c.q.Quest.Msg
	}
	return ""
}

func (c *context) GetReload() int {
	if c.q != nil {
		return c.q.reload.loop
	}
	return 0
}

func (c *context) SetPrfx(w io.Writer, t interface{}) {
	if c.q != nil {
		c.q.prefix = prefix{w, t}
		return
	}
	c.i.prefix = prefix{w, t}
	return
}

func (c *context) SetDef(v interface{}, t interface{}) {
	if c.q != nil {
		c.q.def = def{v, t}
		return
	}
	c.i.def = def{v, t}
	return
}

func (c *context) SetErr(e interface{}) {
	if c.q != nil {
		c.q.err = e
		return
	}
	c.i.err = e
	return
}

func (c *context) SetEnd(e string) {
	if c.q != nil {
		c.q.end.value = e
		return
	}
	c.i.end.value = e
	return
}

func (c *context) method(f ErrorFunc) error {
	if f != nil {
		if err := f(c); err != nil {
			return err
		}
	}
	return nil
}
