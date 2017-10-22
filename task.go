// Copyright © 2017 Abhi Yerra <abhi@opszero.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"html"

	"github.com/aws/aws-sdk-go/service/mturk"
)

type Task struct {
	// Copied from Workflow
	Title       string
	Description string

	HitID    string
	SourceID string
	Fields   []Field

	MTurk struct {
		QuestionFormAnswers struct {
			Answer []struct {
				QuestionIdentifier string
				FreeText           string
			}
		}

		Assignments []*mturk.Assignment
	}
}

func (t *Task) Question() string {
	var fieldsHTML string
	for i := range t.Fields {
		fieldsHTML += t.Fields[i].HTML()
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
`, html.EscapeString(t.Title), html.EscapeString(t.Description), fieldsHTML)
}
