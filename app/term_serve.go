package workmachine

import (
	"bufio"
	"fmt"
	. "github.com/abhiyerra/workmachine/app"
	"os"
)

type TermServe struct {
	Batch *Batch
}

func (ss TermServe) Execute(batch *Batch) {
	for _, j := range ss.Batch.Jobs {
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
	}
}

func (ss TermServe) Aggregate() {
	fmt.Println(ss.Batch)

	// for _, j := range ss.Batch.Jobs {
	// 	// for _, out := range j.OutputFields {
	// 	// 	fmt.Printf("%s\t%s", out.Id, out.Description)
	// 	// }
	// 	fmt.Println(j)
	// }
}
