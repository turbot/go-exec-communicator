package shared

import "fmt"

type Outputter struct{}

func (*Outputter) Output(s string) {
	fmt.Println(s)
}
