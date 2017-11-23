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
)

type Field struct {
	Name        string
	Type        string
	Description string
	Value       string
	Options     []string
}

func (t *Field) TextArea() string {
	return fmt.Sprintf(`
		<div class="row">
		<div class="col-md-12">
		  <h2>%s</h2>
		  <p>%s</p>
		  <p>
			<textarea name='%s' class="form-control" cols='80' rows='3'>%s</textarea>
		  </p>
		</div>
		</div>`, html.EscapeString(t.Name), html.EscapeString(t.Description), t.Name, t.Value)
}

func (t *Field) TextField() string {
	return fmt.Sprintf(`
		<div class="row">
		<div class="col-md-12">
		  <h2>%s</h2>
		  <p>%s</p>
		  <p>
			<input type="text" class="form-control" name="%s" value="%s" />
		  </p>
		</div>
		</div>`, html.EscapeString(t.Name), html.EscapeString(t.Description), t.Name, t.Value)
}

func (t *Field) SelectField() string {
	options := ""

	for _, i := range t.Options {
		options += fmt.Sprintf(`<option value="%s">%s</option>`, i, i)
	}

	return fmt.Sprintf(`
		<div class="row">
		<div class="col-md-12">
		  <h2>%s</h2>
		  <p>%s</p>
		  <p>
		  	<select name="%s" class="form-control">
				%s
			</select>
		  </p>
		</div>
		</div>`, html.EscapeString(t.Name), html.EscapeString(t.Description), t.Name, options)
}

func (t *Field) RadioFields() string {
	options := ""

	for _, i := range t.Options {
		options += fmt.Sprintf(`<input type="radio" class="form-control" name="%s" value="%s" /> %s<br/>`, t.Name, i, i)
	}

	return fmt.Sprintf(`
		<div class="row">
		<div class="col-md-12">
		  <h2>%s</h2>
		  <p>%s</p>
		  <p>
			%s
		  </p>
		</div>
		</div>`, html.EscapeString(t.Name), html.EscapeString(t.Description), t.Name, options)
}

func (t *Field) HTML() string {
	switch t.Type {
	case "LongText":
		return t.TextArea()
	case "ShortText":
		return t.TextField()
	case "Select":
		return t.SelectField()
	case "Radio":
		return t.RadioFields()
	}

	return t.TextArea()
}
