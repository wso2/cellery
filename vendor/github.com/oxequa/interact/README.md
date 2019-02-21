### Interact

An easy and fast Go library, without external imports, to handle questions and answers by command line

##### Features

- [Single question](#single-question)
- [Questions list](#questions-list)
- [Multiple choice](#multiple-choice)
- [Sub questions](#sub-questions)
- [Question prefix](#question-prefix)
- [Default values](#default-values)
- [Custom errors](#custom-errors)
- [After/Before listeners](#after-before)
- [Skip a Question](#skip-a-question)
- [Reload a Question](#reload-a-question)
- [End signal](#end-signal)
- [Colors support (fatih/color)](#color-support)

##### Installation

To install interact:
```
$ go get github.com/tockins/interact
```

##### Single question

Run a simple question and manage the response. 
The response field is used to get the answer as a specific type.
``` go
package main

import (
	i "github.com/tockins/interact"
)

func main() {
	i.Run(&i.Interact{
		Questions: []*i.Question{
			{
				Quest: i.Quest{
					Msg:      "Would you like some coffee?",
				},
				Action: func(c i.Context) interface{} {
					val, err := c.Ans().Bool()
					if err != nil{
					    return err
					}
					fmt.Println(val)
					return nil
				},
			},
		},
	})
}
``` 

##### Questions list

Define a list of questions to be run in sequence.
The Action func can be used for validate the answer and can return a custom error.

Question struct is only for single question whereas **Interact struct** supports multiple questions
``` go
package main

import (
	i "github.com/tockins/interact"
)

func main() {
	i.Run(&i.Interact{
		Questions: []*i.Question{
			{
				Quest: i.Quest{
					Msg:     "Would you like some coffee?",
				},
				Action: func(c i.Context) interface{} {
					val, err := c.Ans().Bool()
					if err != nil{
					    return err
					}
					fmt.Println(val)
					return nil
				},
			},
			{
				Quest: i.Quest{
					Msg:     "What's 2+2?",
				},
				Action: func(c i.Context) interface{} {
				    val, _ := c.Ans().Int()
					// get the answer as integer
					if val < 4 {
						// return a custom error and rerun the question
						return "INCREASE"
					}else if val > 4 {
						return "DECREASE"
					}
					return nil
				},
			},
		},
	})
}
```

##### Multiple choice

Define a multiple choice question

``` go
package main

import (
	i "github.com/tockins/interact"
)

func main() {
	i.Run(&i.Interact{
		Questions: []*i.Question{
			{
				Quest: i.Quest{
                    Msg:     "how much for a teacup?",
                    Choices: i.Choices{
                        Alternatives: []i.Choice{
                            {
                                Text: "Gyokuro teapcup",
                                Response: "20",
                            },
                            {
                                Text: "Sencha teacup",
                                Response: -10,
                            },
                            {
                                Text: "Matcha teacup",
                                Response: 15.50,
                            },
                        },
                    },
                },
                Action: func(c i.Context) interface{} {
                    val, _ := c.Ans().Int()
                    fmt.Println(val)
                    return nil
                },
			},
		},
	})
}
```

##### Sub questions

The sub questions list is managed by the **"Resolve"** func.
Each sub question can access to the parent answer by the **"Parent"** method

``` go
package main

import (
	i "github.com/tockins/interact"
)

func main() {
    i.Run(&i.Interact{
        Questions: []*i.Question{
            {
                Quest: i.Quest{
                    Msg:     "Would you like some coffee?",
                    Resolve: func(c i.Context) bool {     
                        val, _ := c.Ans().Bool()
                        return val
                    },
                },
                Subs: []*i.Question{
                    {
                        Quest: i.Quest{
                            Msg:     "What Kind of Coffee?",
                            Choices: i.Choices{
                                Alternatives: []i.Choice{
                                    {
                                        Text: "Black coffee",
                                        Response: "black",
                                    },
                                    {
                                        Text: "With milk",
                                        Response: "milk",
                                    },
                                },
                            },
                        },
                        Action: func(c i.Context) interface{} {
                            // question (sub) answer
                            val, _ := c.Ans().String()
                            fmt.Println(val)
                            // parent answer
                            val, _ = c.Parent().Ans().String()
                            fmt.Println(val)
                            return nil
                        },
                    },
                },
                Action: func(c i.Context) interface{} {
                    // question answer   
                    val, _ := c.Ans().String()
                    fmt.Println(val)
                    // sub question answer
                    fmt.Println(c.Qns().Get(0).Ans().Raw())
                    return nil
                },
            },
        },
    })
}
```

##### Question prefix

Interact support a custom prefix for each question

You can define a **global prefix** for all questions but you can **overwrite it** in each question with ease

As the first param you can pass a custom **io.writer** instance

``` go
package main

import (
	i "github.com/tockins/interact"
)

func main() {
    i.Run(&i.Interact{
        Before: func(c i.Context) error{
            c.SetPrfx(nil,"GLOBAL PREFIX")
            return nil
        },
        Questions: []*i.Question{
            {
                Before: func(c i.Context) error{
                    c.SetPrfx(nil,"OVERWRITTEN PREFIX")
                    // print current prefix
                    fmt.Println(c.Prfx())
                    return nil
                },
                Quest: i.Quest{
                    Msg:     "Would you like some coffee?",
                },
                Action: func(c i.Context) interface{} {
                    return nil
                },
            },
            {
                Before: func(c i.Context) error{
                    fmt.Println(c.Prfx())
                    return nil
                },
                Quest: i.Quest{
                    Msg:     "What's 2+2?",
                },
                Action: func(c i.Context) interface{} {
                    return nil
                },
            },
        },
    })
}
```

##### Default values

You can define a default value for each question and get it in the action func as an answer

``` go
package main

import (
	i "github.com/tockins/interact"
)

func main() {
    i.Run(&i.Interact{
        Questions: []*i.Question{
            {
                Before: func(c i.Context) error{
                    c.SetDef("test",default val")
                    return nil
                },
                Quest: i.Quest{
                    Msg:     "Would you like some coffee?",
                },
                Action: func(c i.Context) interface{} {
                    val, _ := c.Ans().String()
                    fmt.Println(val)
                    return nil
                },
            },
            {
                Before: func(c i.Context) error{
                    c.SetDef("default","default val")
                    return nil
                },
                Quest: i.Quest{
                    Msg:     "Would you like some coffee?",
                },
                Action: func(c i.Context) interface{} {
                    val, _ := c.Ans().Bool()
                    fmt.Println(val)
                    return nil
                },
            },
        },
    })
}
``` 

##### Custom errors

You can define a default error for every question or you can set a default error message


``` go
package main

import (
	i "github.com/tockins/interact"
)

func main() {

	i.Run(&i.Interact{
		Before: func(c i.Context) error{
			c.SetErr("Default error")
			return nil
		},
		Questions: []*i.Question{
			{
				Quest: i.Quest{
					Msg: "Would you like some coffee?",
					Err: "Custom error fot this question",
				},
				Action: func(c i.Context) interface{} {
                    val, err := c.Ans().Bool()
					if err {
						return "Invalid answer"
					}
					return nil
				},
			},
			{
				Quest: i.Quest{
					Msg: "Would you like some coffee?",
				},
				Action: func(c i.Context) interface{} {
					val, err := c.Ans().Bool()
                    if err {
                        return c.Err()
                    }
                    return nil
				},
			},
		},
	})
}
```

##### After Before

For every question and for each list of questions you can define custom commands to be run before or after the relative instance


``` go
package main

import (
	i "github.com/tockins/interact"
)

function main(){
    i.Run(&i.Interact{
        Before: func(c i.Context) error{
            c.SetPrfx(nil, "TEST")
            return nil
        },
        Questions: []*i.Question{
            {
                Before: func(c i.Context) error{
                    c.SetPrfx(nil, "TEST A")
                    return nil
                },
                Quest: i.Quest{
                    Msg:     "How much coffee?",
                },
                Action: func(c i.Context) interface{} {
                    return nil
                },
                After: func(c i.Context) error{
                    val, _ := c.Ans().Int()
                    fmt.Println(val)
                    return nil
                },
            },
            {
                Quest: i.Quest{
                    Msg:     "How much coffee?",
                },
                Action: func(c i.Context) interface{} {
                    return nil
                },
            },
        },
        After: func(c i.Context) error{
            for _, v := range c.Qns().List(){
                fmt.Println(v.Quest(),v.Ans().Raw())
            }
            return nil
        },
    })
}
```

##### Skip a Question

With the skip func you can stop the execution of the current question or you can skip the next.

``` go
package main

import (
	i "github.com/tockins/interact"
)

function main(){
    i.Run(&i.Interact{
        Before: func(c i.Context) error{
            // skip all questions
            //c.Skip()
            return nil
        },
        Questions: []*i.Question{
            {
                Before: func(c i.Context) error{
                    // skip the current question
                    //c.Skip()
                    return nil
                },
                Quest: i.Quest{
                    Msg:     "How much coffee?",
                },
                Action: func(c i.Context) interface{} {
                    return nil
                },
                After: func(c i.Context) error{
                    // skip the next question
                    c.Skip()
                    return nil
                },
            },
            {
                Before: func(c i.Context) error{
                    return nil
                },
                Quest: i.Quest{
                    Msg:     "How much tea?",
                },
                Action: func(c i.Context) interface{} {
                    return nil
                },
                After: func(c i.Context) error{
                    return nil
                },
            },
        },
    })
}
```

##### Reload a Question

You can reload a question how many times as you want

``` go
package main

import (
	i "github.com/tockins/interact"
)

function main(){
    i.Run(&i.Interact{
        Questions: []*i.Question{
            {
                Quest: i.Quest{
                    Msg:     "Would you like Interact?",
                },
                Action: func(c i.Context) interface{}{
                    val, err := c.Ans().Bool()
                    if (err != nil || !val){
                        c.Reload()
                    }
                    return nil
                },
            },
        },
    })
}
```

##### End signal

End a group of questions or sub-questions with a specific character or string

```
package main

import (
	i "github.com/tockins/interact"
)

func main() {
	i.Run(&i.Interact{
		Before: func(c i.Context) error {
			c.SetEnd("!*")
			return nil
		},
		Questions: []*i.Question{
			{
				Before: func(c i.Context) error {
					c.SetEnd("*")
					return nil
				},
				Quest: i.Quest{
					Msg: "Would you like some coffee? (insert '*' to stop this question or the sub questions)",
					Resolve: func(c i.Context) bool {
						val, _ := c.Ans().Bool()
						return val
					},
				},
				Subs: []*i.Question{
					{
						Before: func(c i.Context) error {
							c.SetEnd("!")
							return nil
						},
						Quest: i.Quest{
							Msg: "What kind of Coffee? (insert '!' to stop)",
						},
						Action: func(c i.Context) interface{} {
							c.Reload()
							return nil
						},
					},
					{
						Quest: i.Quest{
							Msg: "What type of Coffee?",
						},
					},
				},
			},
			{
				Quest: i.Quest{
					Msg: "Would you like some tea?",
				},
				Action: func(c i.Context) interface{} {
					c.Reload()
					return nil
				},
			},
		},
	})
}
```

##### Colors support

Interact supports the color scheme defined by the package "fatih/color"

``` go
package main

import (
	"github.com/fatih/color"
	i "github.com/tockins/interact"
	"fmt"
)

func main() {

	b := color.New(color.FgHiWhite).Add(color.BgRed).SprintfFunc()
	y := color.New(color.FgYellow).SprintFunc()
	r := color.New(color.FgRed).SprintFunc()
	g := color.New(color.FgGreen).SprintFunc()
	prefix := y("[") + "INTERACT" + y("]")

	i.Run(&i.Interact{
		Before: func(c i.Context) error{
			c.SetPrfx(color.Output, prefix)
			return nil
		},
		Questions: []*i.Question{
			{
				Before: func(c i.Context) error{
					c.SetPrfx(nil,y("[") + "INTERACT QUEST" + y("]"))
					return nil
				},
				Quest: i.Quest{
					Msg:     "Would you like some coffee?",
					Options:  g("[yes/no]"),
					Err:      b("INVALID"),
					Resolve: func(c i.Context) bool{
						val, _ := c.Ans().Bool()
						return val
					},
				},
				Subs: []*i.Question{
					{
						Quest: i.Quest{
							Msg:     "What Kind of Coffee?",
							Err:      b("INVALID"),
							Choices: i.Choices{
								Color: g,
								Alternatives: []i.Choice{
									{
										Text: "Black coffee",
										Response: "black",
									},
									{
										Text: "With milk",
										Response: "milk",
									},
								},
							},
						},
						Action: func(c i.Context) interface{}{
						    val, _ := c.Ans().String()
							fmt.Println(val)
							val, _ := c.Parent().Ans().String()
							fmt.Println(val)
							return nil
						},
					},
				},
				Action: func(c i.Context) interface{} {
				    val, _ := c.Ans().Bool()
					if !val{
						return r("INVALID INPUT")
					}
					fmt.Println(c.Quest(), val)
					return nil
				},
			},
		},
	})
}
```