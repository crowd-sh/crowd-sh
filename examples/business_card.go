package main

type BusinessCard struct {
	ImageUrl    InputField  `crowd_desc:"Use this information for the data below." crowd_id:"image_url" crowd_type:"image"`
	Name        OutputField `crowd_desc:"Find the name from the business card" crowd_id:"name"`
	Company     OutputField `crowd_desc:"Find the company from the business card" crowd_id:"company"`
	Email       OutputField `crowd_desc:"Find the email from the business card" crowd_id:"email"`
	PhoneNumber OutputField `crowd_desc:"Find the phone number from the business card" crowd_id:"phone_number"`
}

type BusinessCardVerify struct {
	ImageUrl InputField  `crowd_desc:"Use this information for the data below." crowd_id:"image_url" crowd_type:"image"`
	Name     InputField  `crowd_desc:"Find the name from the business card" crowd_id:"name"`
	Email    InputField  `crowd_desc:"Find the email from the business card" crowd_id:"email"`
	IsValid  OutputField `crowd_desc:"Is the information above valid?" crowd_id:"is_valid"`
}

func main() {
	business_cards := Task{
		Title:       "Business Card Fields",
		Description: "Enter the fields.",
		Price:       1,
		Tasks: []BusinessCard{
			BusinessCard{
				ImageUrl: "google.com",
			},
			// BusinessCard{
			// 	ImageUrl: "yahoo.com",
			// },
		},
	}

	serve := TermServe{}
	var backend Backender = serve
	backend_cards_results := NewBatch(business_cards).Run(backend)

}
