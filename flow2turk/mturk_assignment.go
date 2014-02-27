package crowdflow

import (
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/exp/mturk"
	"log"
	"time"
)

func NewMTurkAssigner(assignDone chan bool, a *Assignment) {
	ma := MTurkAssignment{
		Assignment: a,
	}

	ma.Create()
	go ma.Poll()
}

type MTurkAssignment struct {
	SharedAssignment

	Assignment *Assignment
	Hit        *mturk.HIT
}

var (
	mturk_assignments MTurkAssignments
)

type MTurkAssignments []*MTurkAssignment

// Create a HIT and send it to MTurk.
func (assign *MTurkAssignment) Create() {
	html_question_template := `<![CDATA[
<!DOCTYPE html>
<html>
 <head>
  <meta http-equiv='Content-Type' content='text/html; charset=UTF-8'/>
  <link href="//netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css" rel="stylesheet">
  <script type='text/javascript' src='https://s3.amazonaws.com/mturk-public/externalHIT_v1.js'></script>
 </head>
 <body>
  <div class="container">
    <h1>%s</h1>

    <div>
      %s
    </div>
    <form name='mturk_form' method='post' id='mturk_form' action='https://www.mturk.com/mturk/externalSubmit'>
    <input type='hidden' value='' name='assignmentId' id='assignmentId'/>
%s

%s
    </form>
   </div>
   <script language='Javascript'>turkSetAssignmentID();</script>
 </body>
</html>
]]>`

	input_html := ""
	for _, inp := range assign.Assignment.Job.InputFields {
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
	for _, out := range assign.Assignment.Job.OutputFields {
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

	html_question := fmt.Sprintf(html_question_template, assign.Assignment.Job.TaskDesc.Title, assign.Assignment.Job.TaskDesc.Description, input_html, output_html)

	hit, err := MTurkAuth().CreateHIT(
		assign.Assignment.Job.TaskDesc.Title,
		assign.Assignment.Job.TaskDesc.Description,
		mturk.HTMLQuestion{
			HTMLContent: mturk.HTMLContent{html_question},
			FrameHeight: 1000,
		},
		mturkPrice(assign.Assignment.Job.TaskDesc.Price),
		100,
		200,
		assign.Assignment.Job.TaskDesc.Tags,
		5,
		nil,
		"",
	)

	assign.Hit = hit

	fmt.Printf("%v", assign.Hit)

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
		assignment_resp, err := MTurkAuth().GetAssignmentsForHIT(a.Hit.HITId)

		if err != nil {
			panic(err)
		}

		fmt.Printf("No assignment yet\n")

		if assignment_resp.AssignmentStatus == "Submitted" {
			answers := assignment_resp.Answers()
			log.Printf("Answer: %s", answers)

			for i, out := range a.Assignment.Job.OutputFields {
				if v, ok := answers[out.Id]; ok {
					a.Assignment.Job.OutputFields[i].Value = v
				}
			}

			a.Assignment.Finished = true
			a.Assignment.AssignDone <- true

			a.Assignment.Job.TaskDesc.Write(a.Assignment.Job)

			break
		}

		// Wait for 1 minutes before checking again.
		time.Sleep(1 * time.Minute)
	}
}

func MTurkAuth() *mturk.MTurk {
	if config.mturkAuth == nil {
		config.mturkAuth = mturk.New(aws.Auth{
			AccessKey: config.AwsAccessKey,
			SecretKey: config.AwsSecretKey,
		}, config.AwsSandbox)

	}

	return config.mturkAuth
}
