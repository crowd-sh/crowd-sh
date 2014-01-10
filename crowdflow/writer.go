package crowdflow

import (
	"encoding/csv"
	"fmt"
	"os"
)

func CsvJobWriter(file *os.File) func(j *MetaJob) {
	writer := csv.NewWriter(file)

	return func(j *MetaJob) {
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
