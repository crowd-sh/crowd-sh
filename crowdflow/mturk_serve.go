package crowdflow

// import (
// 	"fmt"
// 	"github.com/crowdmob/goamz/aws"
// 	"github.com/crowdmob/goamz/exp/mturk"
// 	"time"
// )

// type MTurkServe struct{}

// func (m MTurkServe) Execute(jobs chan Job, j Job) {
// 	assignment := &MTurkAssignment{
// 		JobsChan: jobs,
// 		Job:      j,
// 	}

// 	assignment.Create()
// 	go assignment.Poll()

// 	mturk_assignments = append(mturk_assignments, assignment)
// }

// type MTurkAssignment struct {
// 	JobsChan chan Job
// 	Job      Job
// 	Hit      mturk.HIT
// }

// // Create a HIT and send it to MTurk.
// func (assign *MTurkAssignment) Create() {
// 	html_question_template := `
// <HTMLQuestion xmlns="http://mechanicalturk.amazonaws.com/AWSMechanicalTurkDataSchemas/2011-11-11/HTMLQuestion.xsd">
//   <HTMLContent><![CDATA[
// <!DOCTYPE html>
// <html>
//  <head>
//   <meta http-equiv='Content-Type' content='text/html; charset=UTF-8'/>
//   <script type='text/javascript' src='https://s3.amazonaws.com/mturk-public/externalHIT_v1.js'></script>
//  </head>
//  <body>
//   <form name='mturk_form' method='post' id='mturk_form' action='https://www.mturk.com/mturk/externalSubmit'>
//   <input type='hidden' value='' name='assignmentId' id='assignmentId'/>
// %s

// %s
//   </form>
//   <script language='Javascript'>turkSetAssignmentID();</script>
//  </body>
// </html>
// ]]>
//   </HTMLContent>
//   <FrameHeight>450</FrameHeight>
// </HTMLQuestion>`

// 	//	question_form := mturk.QuestionForm{}

// 	input_html := ""
// 	for _, inp := range assign.Job.InputFields {
// 		input_html += "<div class=form-group>\n"
// 		input_html += fmt.Sprintf("<label>%s</label><br>\n", inp.Description)
// 		switch inp.Type {
// 		case ImageType:
// 			input_html += fmt.Sprintf("<img src=\"%s\" />\n", inp.Value)
// 		default:
// 			input_html += fmt.Sprintf("<p>%s</p>\n", inp.Value)
// 		}
// 		input_html += "</div>\n"
// 	}

// 	output_html := ""
// 	for _, out := range assign.Job.OutputFields {
// 		output_html += "<div class=form-group>\n"
// 		output_html += fmt.Sprintf("<label>%s</label>\n", out.Description)
// 		switch out.Type {
// 		case CheckBoxType:
// 			output_html += fmt.Sprintf("<input type=checkbox name=\"%s\" value=\"yes\"/>\n", out.Id)
// 		default:
// 			output_html += fmt.Sprintf("<br><input type=text name=\"%s\" value=\"%s\"/>\n", out.Id, out.Value)
// 		}
// 		output_html += "</div>\n"
// 	}

// 	output_html += `<input type=submit value=Submit class="btn btn-default" />\n`

// 	html_question := fmt.Sprintf(html_question_template, input_html, output_html)

// 	fmt.Println(html_question)

// 	//	mturk.CreateHit()
// 	// Create the Hit on MTurk and send it out to be uploaded and
// 	// used.
// }

// // Connect to Mturk every so often to check if the work has
// // been updated.
// func (a *MTurkAssignment) Poll() {
// 	for {
// 		hit_assignment, err := mturk_auth.GetAssignmentsForHIT(a.Hit.HITId)
// 		if err == nil {
// 			// TODO: Add the assignment values to the Job
// 			a.JobsChan <- a.Job
// 			break
// 		}

// 		time.Sleep(2 * time.Minute)

// 		_ = hit_assignment
// 	}
// }

// type MTurkAssignments []*MTurkAssignment

// var (
// 	mturk_assignments MTurkAssignments
// 	mturk_auth        *mturk.MTurk
// )

// func MTurkServer(access_key, secret_key string) {
// 	auth := aws.Auth{
// 		AccessKey: access_key,
// 		SecretKey: secret_key,
// 	}
// 	fmt.Println(auth)

// 	mturk_auth = mturk.New(auth, true)
// 	hit, err := mturk_auth.CreateHIT(
// 		"hello2",
// 		"hi",
// 		mturk.ExternalQuestion{
// 			ExternalURL: "http://britishlibrary.crowdflow.us",
// 			// 		mturk.HTMLQuestion{
// 			// 			HTMLContent: mturk.HTMLContent{`<![CDATA[
// 			// <!DOCTYPE html>
// 			// <html>
// 			//  <head>
// 			//   <meta http-equiv='Content-Type' content='text/html; charset=UTF-8'/>
// 			//   <script type='text/javascript' src='https://s3.amazonaws.com/mturk-public/externalHIT_v1.js'></script>
// 			//  </head>
// 			//  <body>
// 			//   <form name='mturk_form' method='post' id='mturk_form' action='https://www.mturk.com/mturk/externalSubmit'>
// 			//   <input type='hidden' value='' name='assignmentId' id='assignmentId'/>
// 			//   <h1>What's up?</h1>
// 			//   <p><textarea name='comment' cols='80' rows='3'></textarea></p>
// 			//   <p><input type='submit' id='submitButton' value='Submit' /></p></form>
// 			//   <script language='Javascript'>turkSetAssignmentID();</script>
// 			//  </body>
// 			// </html>
// 			// ]]>`},
// 			FrameHeight: 1000,
// 		},
// 		mturk.Price{"0.01", "USD", ""},
// 		100,
// 		200,
// 		"test, image",
// 		5,
// 		nil,
// 		"",
// 	)

// 	fmt.Println(hit)

// 	h, err := mturk_auth.GetAssignmentsForHIT("2C8VJ8EBDPGFDNAPG6Q9UNDWVSI0FI")
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(h)
// }
