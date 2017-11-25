package sources

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// ExecSource is a way to get input and output from command line expressions.
// Essentially this needs to be a list of JSON values that will be used to send
// values to MTurk.
type ExecSource struct {
	Config map[string]string

	headers []string
	records []map[string]string
}

func (w *ExecSource) Init() {
	cmd := strings.Split(w.Config["Command"], " ")

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(out))
	err = json.Unmarshal(out, &w.records)
	log.Println(err)

	log.Println(w.records)

	var headers []string
	for k := range w.records[0] {
		headers = append(headers, k)
	}

	w.headers = headers
}

func (w *ExecSource) Headers() []string {
	return w.headers
}

func (w *ExecSource) Records() []map[string]string {
	return w.records
}

func (w *ExecSource) WriteAll(headers []string, rows []map[string]string) {
	cmd2 := strings.Split(w.Config["Command"], " ")
	cmd := exec.Command(cmd2[0], cmd2[1:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		b, _ := json.Marshal(rows)

		defer stdin.Close()
		stdin.Write(b)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
