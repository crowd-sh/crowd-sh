package backend

type SelfServe struct {
}

func (ss SelfServe) Publish(batch Batch) {
	// func (jf JobField) Html() (html string) {
	// 	html += "    <div>\n"
	// 	html += "        <label>" + jf.Description + "</label>\n"
	// 	switch jf.Type {
	// 	case "image":
	// 		html += "        <img src=\"" + jf.Value + "\"/>\n"
	// 	default:
	// 		html += "        <input type=\"text\" name=\"" + jf.Id + "\" value=\"" + jf.Description + "\"/>\n"
	// 	}
	// 	html += "    </div>"
	// 	return
	// }
}

func (ss SelfServe) Execute(j Job) {
}
