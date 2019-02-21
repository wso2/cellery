package interact

import (
	"io"
)

type (
	Interview interface {
		//Next()
		//Prev()
		//Goto(interface{})
	}
	interview struct {
		index    int
		interact *Interact
		current  current
	}
	current struct {
		list  []*Question
		index int
	}
)

// Interact element
type Interact struct {
	settings
	Questions     []*Question
	After, Before ErrorFunc
}

// General settings and options
type settings struct {
	skip bool
	err  interface{}
	prefix
	def
	end
}

// Answer end
type end struct {
	status bool
	value  string
}

// Default answer value
type def struct {
	Value interface{}
	Text  interface{}
}

// Questions prefix
type prefix struct {
	Writer io.Writer
	Text   interface{}
}

// Run a questions list
func Run(i *Interact) error {
	if err := i.ask(); err != nil {
		return err
	}
	return nil
}

// Create a new interact configuration
func New(i *Interact) Interview {
	return &interview{interact: i, current: current{index: 0, list: i.Questions}}
}

// Ask one by one with interact after/before
func (i *Interact) ask() (err error) {
	context := &context{i: i}
	if err := context.method(i.Before); err != nil {
		return err
	}
	if !context.i.skip {
		for index := range i.Questions {
			i.Questions[index].interact = i
			if err = i.Questions[index].ask(); err != nil {
				return err
			}
		}
		if err := context.method(i.After); err != nil {
			return err
		}
	}
	return nil
}
