package crowdflow

import (
	//	"encoding/xml"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/exp/mturk"
	"time"
)

type MTurkServe struct{}

func (m MTurkServe) Execute(jobs chan MetaJob, j MetaJob) {
	assignment := &MTurkAssignment{
		JobsChan: jobs,
		Job:      j,
	}

	assignment.Create()
	mturk_assignments = append(mturk_assignments, assignment)

	assignment.Poll()
}

type MTurkAssignment struct {
	Assignment
	JobsChan chan MetaJob
	Job      MetaJob
	Hit      *mturk.HIT
}

var (
	mturk_assignments MTurkAssignments
	mturk_auth        *mturk.MTurk
)

type MTurkAssignments []*MTurkAssignment

// Create a HIT and send it to MTurk.
func (assign *MTurkAssignment) Create() {
	html_question_template := `<![CDATA[
<!DOCTYPE html>
<html>
 <head>
  <meta http-equiv='Content-Type' content='text/html; charset=UTF-8'/>
  <script type='text/javascript' src='https://s3.amazonaws.com/mturk-public/externalHIT_v1.js'></script>
 </head>
 <body>
  <form name='mturk_form' method='post' id='mturk_form' action='https://www.mturk.com/mturk/externalSubmit'>
  <input type='hidden' value='' name='assignmentId' id='assignmentId'/>
%s

%s
  </form>
  <script language='Javascript'>turkSetAssignmentID();</script>
 </body>
</html>
]]>`

	input_html := ""
	for _, inp := range assign.Job.InputFields {
		input_html += "<div class=form-group>\n"
		input_html += fmt.Sprintf("<label>%s</label><br>\n", inp.Description)
		switch inp.Type {
		case ImageType:
			input_html += fmt.Sprintf("<img src=\"%s\" />\n", inp.Value)
		default:
			input_html += fmt.Sprintf("<p>%s</p>\n", inp.Value)
		}
		input_html += "</div>\n"
	}

	output_html := ""
	for _, out := range assign.Job.OutputFields {
		output_html += "<div class=form-group>\n"
		output_html += fmt.Sprintf("<label>%s</label>\n", out.Description)
		switch out.Type {
		case CheckBoxType:
			output_html += fmt.Sprintf("<input type=checkbox name=\"%s\" value=\"yes\"/>\n", out.Id)
		default:
			output_html += fmt.Sprintf("<br><input type=text name=\"%s\" value=\"%s\"/>\n", out.Id, out.Value)
		}
		output_html += "</div>\n"
	}

	output_html += `<input type=submit value=Submit class="btn btn-default" />`

	html_question := fmt.Sprintf(html_question_template, input_html, output_html)

	hit, err := mturk_auth.CreateHIT(
		assign.Job.TaskDesc.Title,
		assign.Job.TaskDesc.Description,
		mturk.HTMLQuestion{
			HTMLContent: mturk.HTMLContent{html_question},
			FrameHeight: 1000,
		},
		mturkPrice(assign.Job.TaskDesc.Price),
		100,
		200,
		assign.Job.TaskDesc.Tags,
		5,
		nil,
		"",
	)

	assign.Hit = hit

	fmt.Printf("%v", hit)

	if err != nil {
		panic(err)
	}
}

func mturkPrice(price uint) mturk.Price {
	return mturk.Price{
		fmt.Sprintf("%0.2f", float32(price)/100),
		"USD",
		"",
	}
}

// Connect to Mturk every so often to check if the work has
// been updated.
func (a *MTurkAssignment) Poll() {
	for {
		assignment_resp, err := mturk_auth.GetAssignmentsForHIT(a.Hit.HITId)

		if err != nil {
			panic(err)
		}

		fmt.Printf("No assignment yet\n")

		if assignment_resp.GetAssignmentsForHITResult.TotalNumResults > 0 {
			mturk_assignment := assignment_resp.GetAssignmentsForHITResult.Assignment

			fmt.Printf("Answer: %s", mturk_assignment.Answer)

			a.JobsChan <- a.Job
			break
		}

		// Wait for 1 minutes before checking again.
		time.Sleep(1 * time.Minute)
	}
}

func MTurkServer(access_key, secret_key string, sandbox bool) {
	auth := aws.Auth{
		AccessKey: access_key,
		SecretKey: secret_key,
	}

	mturk_auth = mturk.New(auth, sandbox)
}
