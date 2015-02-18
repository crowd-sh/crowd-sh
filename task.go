package main

import (
	"github.com/crowdmob/goamz/exp/mturk"
	"time"
)

const (
	Working  = "Working"
	Finished = "Finished"
)

type Task struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Status    string
	CreatedAt time.Time `json:"-"`

	TaskFields []TaskField
}

type TaskField struct {
	TaskId int64

	Type  string
	Key   string
	Value string
}

func (t *Task) PublishToMTurk() {
	question := mturk.HTMLQuestion{
		HTMLContent: mturk.HTMLContent{`<![CDATA[
<!DOCTYPE html>
<html>
 <head>
  <meta http-equiv='Content-Type' content='text/html; charset=UTF-8'/>
  <script type='text/javascript' src='https://s3.amazonaws.com/mturk-public/externalHIT_v1.js'></script>
 </head>
 <body>
  <form name='mturk_form' method='post' id='mturk_form' action='https://www.mturk.com/mturk/externalSubmit'>
  <input type='hidden' value='' name='assignmentId' id='assignmentId'/>
  <h1>What's up?</h1>
  <p><textarea name='comment' cols='80' rows='3'></textarea></p>
  <p><input type='submit' id='submitButton' value='Submit' /></p></form>
  <script language='Javascript'>turkSetAssignmentID();</script>
 </body>
</html>
]]>`},
		FrameHeight: 200,
	}
	reward := mturk.Price{
		Amount:       "0.01",
		CurrencyCode: "USD",
	}

	hit, err := s.mturk.CreateHIT("title", "description", question, reward, 1, 2, "key1,key2", 3, nil, "annotation")

}

func (t *Task) PollMTurk() {

}

// func (t *Task) Status() {
// 	// All the assignments are done
// }

// func NewTask(w *Workflow) {
// 	assignDone := make(chan bool)

// 	// Create the assignments
// 	log.Println("Creating assignments of jobs", len(p.Jobs))

// 	for i := range p.Jobs {
// 		for j := range p.Jobs[i].OutputFields {
// 			go NewAssignment(assignDone, &p.Jobs[i], &p.Jobs[i].OutputFields[j])
// 		}
// 	}

// 	// Wait for responses
// 	log.Println("Created assignments")
// 	for _ = range assignDone {
// 		log.Println("Waiting for assignments:", len(AvailableAssignments))
// 	}
// }
