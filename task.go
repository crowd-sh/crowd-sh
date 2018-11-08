package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mturk"

	"golang.org/x/net/html/charset"
)

type questionFormAnswers struct {
	Answer []struct {
		QuestionIdentifier string
		FreeText           string
	}
}

type Task struct {
	workflow *Workflow

	AirtableID string `json:"id"`
	Fields     map[string]string

	MTurk struct {
		HIT                 *mturk.HIT
		QuestionFormAnswers questionFormAnswers
		Assignments         []*mturk.Assignment
	}
}

func (t *Task) Question() string {
	var fieldsHTML string
	for k, v := range t.Fields {
		if k == MTurkDataField || k == MTurkHitIDField {
			continue
		}

		for _, ft := range t.workflow.FieldTypes {
			if ft.Name == k {
				ft.Value = v
				fieldsHTML += ft.HTML()
			}
		}
	}

	return fmt.Sprintf(`
<HTMLQuestion xmlns="http://mechanicalturk.amazonaws.com/AWSMechanicalTurkDataSchemas/2011-11-11/HTMLQuestion.xsd">
  <HTMLContent><![CDATA[
<!DOCTYPE html>
<html>
 <head>
  <meta http-equiv='Content-Type' content='text/html; charset=UTF-8'/>
  <script type='text/javascript' src='https://s3.amazonaws.com/mturk-public/externalHIT_v1.js'></script>

<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>

 </head>
 <body>
  <div class="container">
    <form name='mturk_form' method='post' id='mturk_form' action='https://www.mturk.com/mturk/externalSubmit'>
    <h1>%s</h1>

    <p>
    %s
    </p>

    %s

    <p>
	<input type='hidden' value='' name='assignmentId' id='assignmentId'/>
	<input type='submit' id='submitButton' value='Submit' class='btn btn-success' />
    </p>
    </form>
  </div>

  <script language='Javascript'>turkSetAssignmentID();</script>
 </body>
</html>
]]>
  </HTMLContent>
  <FrameHeight>1000</FrameHeight>
</HTMLQuestion>
`, html.EscapeString(t.workflow.Title), html.EscapeString(t.workflow.Description), fieldsHTML)
}

func (t *Task) Sync(w *Workflow) {
	t.workflow = w

	if t.Fields[MTurkHitIDField] == "" {
		log.Println("New")
		resp, err := w.client.CreateHITWithHITType(&mturk.CreateHITWithHITTypeInput{
			HITTypeId:         aws.String(w.MTurk.HitTypeId),
			MaxAssignments:    aws.Int64(1),
			Question:          aws.String(t.Question()),
			LifetimeInSeconds: aws.Int64(86400), // 1 day
		})

		if err == nil {
			t.Fields[MTurkHitIDField] = *resp.HIT.HITId
			t.Save()
		} else {
			fmt.Println(err)
			if r := recover(); r != nil {
				fmt.Println("Recovered", r)
			}
			t.Save()
		}
	} else {
		log.Println("Update")

		resp, err := w.client.ListAssignmentsForHIT(&mturk.ListAssignmentsForHITInput{
			HITId: aws.String(t.Fields[MTurkHitIDField]),
		})

		fmt.Println(err)
		fmt.Println(resp)

		if len(resp.Assignments) > 0 {
			var q questionFormAnswers

			b := bytes.NewBufferString(*resp.Assignments[0].Answer)
			decoder := xml.NewDecoder(bufio.NewReader(b))
			decoder.CharsetReader = charset.NewReaderLabel
			err = decoder.Decode(&q)

			// err := xml.Unmarshal(, &q)
			if err != nil {
				log.Println(err)
			}
			t.MTurk.QuestionFormAnswers = q

			log.Println("Hi")
			for _, answer := range t.MTurk.QuestionFormAnswers.Answer {
				log.Println("Answer", answer)
				t.Fields[answer.QuestionIdentifier] = strings.TrimSpace(answer.FreeText)
			}

			log.Println("ID", t.AirtableID)
		}

		t.Save()

	}
}

func (t *Task) Save() {
	updatedFields := make(map[string]interface{})

	for k, v := range t.Fields {
		updatedFields[k] = v
	}

	log.Println(updatedFields)

	if err := t.workflow.AirTable.client.UpdateRecord(t.workflow.AirTable.Table, t.AirtableID, updatedFields, t); err != nil {
		panic(err)
	}

}

// func (t *Task) Approve(w *Workflow) {
// 	log.Println("Approve")

// 	w.client.ApproveAssignment(&mturk.ApproveAssignmentInput{
// 		AssignmentId:      t.MTurk.Assignments[0].AssignmentId,
// 		RequesterFeedback: aws.String("Great Job"),
// 	})
// }

// func (t *Task) Reject(w *Workflow) {
// 	log.Println("Rejected")

// 	w.client.RejectAssignment(&mturk.RejectAssignmentInput{
// 		AssignmentId:      t.MTurk.Assignments[0].AssignmentId,
// 		RequesterFeedback: aws.String("Incorrect Input / Spam"),
// 	})

// 	// So we can remake it.
// 	t.HitID = ""
// }

// func (t *Task) Expire(w *Workflow) {
// 	if len(t.MTurk.Assignments) == 0 || (len(t.MTurk.Assignments) > 0 && (*t.MTurk.Assignments[0].AssignmentStatus != "Approved" && *t.MTurk.Assignments[0].AssignmentStatus != "Submitted")) {
// 		log.Println(t.HitID)
// 		w.client.UpdateExpirationForHIT(&mturk.UpdateExpirationForHITInput{
// 			ExpireAt: aws.Time(time.Now().AddDate(0, 0, -1)),
// 			HITId:    aws.String(t.HitID),
// 		})
// 	}

// 	t.HitID = ""
// }
