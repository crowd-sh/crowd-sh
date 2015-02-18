package main

import (
	"encoding/json"
	"github.com/crowdmob/goamz/exp/mturk"
	"log"
	"time"
)

const (
	Working  = "Working"
	Finished = "Finished"
)

type Task struct {
	Id           int64       `json:"id"`
	WorkflowId   int64       `json:"-"`
	Status       string      `json:"status"`
	InputFields  []TaskField `json:"input_fields"`
	OutputFields []TaskField `json:"output_fields"`
	Token        string      `json:"-"` // MTurkId
	RawData      string      `json:"-"`
	CreatedAt    time.Time   `json:"-"`
}

type TaskField struct {
	Key   string
	Value string
}

func (t *Task) Parse() {
	err := json.Unmarshal([]byte(t.RawData), t)
	if err != nil {
		log.Println(err)
	}
}

func (t *Task) VerifyWithWorkflow(w *Workflow) bool {
	// TODO: Verify that the workflow inputs are included here.

	return true
}

func (t *Task) PublishToMTurk() {
	// question := mturk.HTMLQuestion{
	// 	HTMLContent: mturk.HTMLContent{`<![CDATA[
	// <!DOCTYPE html>
	// <html>
	//  <head>
	//   <meta http-equiv='Content-Type' content='text/html; charset=UTF-8'/>
	//   <script type='text/javascript' src='https://s3.amazonaws.com/mturk-public/externalHIT_v1.js'></script>
	//  </head>
	//  <body>
	//   <form name='mturk_form' method='post' id='mturk_form' action='https://www.mturk.com/mturk/externalSubmit'>
	//   <input type='hidden' value='' name='assignmentId' id='assignmentId'/>
	//   <h1>What's up?</h1>
	//   <p><textarea name='comment' cols='80' rows='3'></textarea></p>
	//   <p><input type='submit' id='submitButton' value='Submit' /></p></form>
	//   <script language='Javascript'>turkSetAssignmentID();</script>
	//  </body>
	// </html>
	// ]]>`},
	// 	FrameHeight: 200,
	// }

	// reward := mturk.Price{
	// 	Amount:       "0.01",
	// 	CurrencyCode: "USD",
	// }

	// hit, err := s.mturk.CreateHIT("title", "description", question, reward, 1, 2, "key1,key2", 3, nil, "annotation")
}

func (t *Task) PollMTurk() {

}
