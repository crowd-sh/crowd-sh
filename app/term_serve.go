package workmachine

import (
	"bufio"
	"fmt"
	"os"
)

type TermServe struct{}

func (ss TermServe) Execute(jobs chan Job, j Job) {
	for _, inp := range j.InputFields {
		fmt.Println(inp.Description)
		fmt.Println(inp.Value)
	}

	bio := bufio.NewReader(os.Stdin)
	for i, out := range j.OutputFields {
		fmt.Println(out.Description)
		fmt.Println(out.Value)
		line, _, _ := bio.ReadLine()
		j.OutputFields[i].Value = string(line)
	}

	fmt.Println(j)

	jobs <- j
}
