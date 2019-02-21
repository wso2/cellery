package interact

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Question params
type Quest struct {
	Choices
	parent            *Question
	Options, Msg, Tag string
	Resolve           BoolFunc
}

// Question entity
type Question struct {
	Quest
	settings
	reload
	choices       bool
	response      string
	value         interface{}
	interact      *Interact
	parent        *Question
	Action        InterfaceFunc
	Subs          []*Question
	After, Before ErrorFunc
}

// Reload entity
type reload struct {
	status bool
	loop   int
}

// Choice option
type Choice struct {
	Text     string
	Response interface{}
}

// Choices list and prefix color
type Choices struct {
	Alternatives []Choice
	Color        func(...interface{}) string
}

func (q *Question) ask() (err error) {
	context := &context{i: q.interact, q: q}
	if q.checkEnd() {
		return nil
	}
	if !context.i.skip {
		if err := context.method(q.Before); err != nil {
			return err
		}
		if !context.i.skip {
			if q.prefix.Text != nil {
				q.print(q.prefix.Text, " ")
			} else if q.parent != nil && q.parent.prefix.Text != nil {
				q.print(q.parent.prefix.Text, " ")
			} else if q.interact.prefix.Text != nil {
				fmt.Print(q.interact.prefix.Text, " ")
			}
			if q.Msg != "" {
				q.print(q.Msg, " ")
			}
			if q.Options != "" {
				q.print(q.Options, " ")
			}
			if q.def.Value != nil && q.def.Text != nil {
				q.print(q.def.Text, " ")
			}
			if q.Alternatives != nil && len(q.Alternatives) > 0 {
				q.multiple()
			}
			if err = q.wait(); err != nil {
				return q.loop(err)
			}
			if q.Subs != nil && len(q.Subs) > 0 {
				if q.Resolve != nil {
					if q.Resolve(context) {
						for _, s := range q.Subs {
							s.interact = q.interact
							s.parent = q
							s.ask()
						}
					}
				} else {
					for _, s := range q.Subs {
						s.interact = q.interact
						s.parent = q.parent
						s.ask()
					}
				}
			}
			if q.checkEnd() {
				return nil
			}
			if q.Action != nil {
				if err := q.Action(context); err != nil {
					q.print(err, " ")
					return q.ask()
				}
			}
			if q.reload.status {
				q.reload.status = false
				q.value = nil
				return q.ask()
			}
			if err := context.method(q.After); err != nil {
				return err
			}
		} else {
			context.i.skip = false
		}
	} else {
		context.i.skip = false
	}
	q.value = nil
	return nil
}

func (q *Question) wait() error {
	reader := bufio.NewReader(os.Stdin)
	if q.choices {
		q.print(q.color("?"), " ", "Answer", " ")
	}
	r, _, err := reader.ReadLine()
	if err != nil {
		return err
	}
	q.response = string(r)
	if abort := q.abort(q.response); abort.value != "" {
		abort.status = true
		return nil
	}
	if len(q.response) == 0 && q.def.Value != nil {
		q.value = q.def.Value
		return nil
	} else if len(q.response) == 0 {
		return errors.New("Answer invalid")
	}
	// multiple choice
	if q.choices {
		choice, err := strconv.ParseInt(q.response, 10, 64)
		if err != nil || int(choice) > len(q.Alternatives) || int(choice) < 1 {
			return errors.New("out of range")
		}
		q.value = q.Alternatives[choice-1].Response
	}
	return nil
}

func (q *Question) print(v ...interface{}) {
	if q.prefix.Writer != nil {
		fmt.Fprint(q.prefix.Writer, v...)
	} else if q.parent != nil && q.parent.prefix.Writer != nil {
		fmt.Fprint(q.parent.prefix.Writer, v...)
	} else if q.interact != nil && q.interact.prefix.Writer != nil {
		fmt.Fprint(q.interact.prefix.Writer, v...)
	} else {
		fmt.Print(v...)
	}

}

func (q *Question) color(v ...interface{}) string {
	if q.Color != nil {
		return q.Color(v...)
	}
	return fmt.Sprint(v...)
}

func (q *Question) loop(err error) error {
	if q.err != nil {
		q.print(q.err, " ")
	} else if q.interact.err != nil {
		q.print(q.interact.err, " ")
	}
	return q.ask()
}

func (q *Question) multiple() error {
	for index, i := range q.Alternatives {
		q.print("\n\t", q.color(index+1, ") "), i.Text, " ")
	}
	q.choices = true
	q.print("\n")
	return nil
}

func (q *Question) abort(s string) *end {
	if q.end.value == "" || q.end.value != s {
		if q.parent != nil {
			return q.parent.abort(s)
		} else if q.interact.end.value == s {
			return &q.interact.end
		}
		return &end{}
	}
	return &q.end
}

func (q *Question) checkEnd() bool {
	if !q.end.status {
		if q.parent != nil {
			return q.parent.checkEnd()
		} else if q.interact.end.status {
			return true
		}
		return false
	}
	return true
}
