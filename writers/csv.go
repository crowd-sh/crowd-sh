package writers

import (
	"encoding/csv"
	"fmt"
	"github.com/abhiyerra/workmachine/crowdflow"
	"os"
)

func Csv(file *os.File) func(j *crowdflow.MetaJob) {
	writer := csv.NewWriter(file)

	return func(j *crowdflow.MetaJob) {
		fmt.Printf("%v\n", j)

		var output []string

		for _, i := range j.InputFields {
			output = append(output, i.Value)
			fmt.Println(i.Value)
		}

		for _, i := range j.OutputFields {
			output = append(output, i.Value)
			fmt.Println(i.Value)
		}

		if err := writer.Write(output); err != nil {
			panic(err)
		}
		writer.Flush()
	}
}
