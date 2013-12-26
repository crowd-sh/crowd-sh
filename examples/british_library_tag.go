package main

import (
	"encoding/csv"
	"fmt"
	. "github.com/abhiyerra/workmachine/app"
	"os"
)

type ImageTagging struct {
	ImageUrl             InputField  `work_desc:"Use this image to fill the information below." work_id:"image_url" work_type:"image"`
	Tags                 OutputField `work_desc:"List all the relavent tags for this image. Ex. trees, castle, person" work_id:"tags"`
	TextInImage          OutputField `work_desc:"Put any text that appears on the image here. One line per block of text." work_id:"text_in_image"`
	IsCorrectOrientation OutputField `work_desc:"Is the image in the correct orientation?" work_id:"is_correct_orientation" work_type:"checkbox"`
	IsLandscape          OutputField `work_desc:"Is the image of a landscape?" work_id:"is_landscape" work_type:"checkbox"`
	IsPattern            OutputField `work_desc:"Is the image of a pattern?" work_id:"is_pattern" work_type:"checkbox"`
	IsPerson             OutputField `work_desc:"Is the image of a person?" work_id:"is_person" work_type:"checkbox"`
	TraditionalClothing  OutputField `work_desc:"Are they wearing traditional costume?" work_id:"tradtional_clothing" work_type:"checkbox"`
}

func imageUrls() (images []ImageTagging) {
	file, err := os.Open("list_of_pictures.csv")
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	for _, i := range records {
		//		fmt.Printf("%s\n", i[1])
		images = append(images, ImageTagging{ImageUrl: InputField(i[1])})
	}

	if err != nil {
		panic(err)
	}

	return
}

func main() {
	results_file, err := os.Create("results.csv")
	if err != nil {
		panic(err)
	}
	defer results_file.Close()

	writer := csv.NewWriter(results_file)

	image_tasks := Task{
		Title:       "Tag the appropriate images",
		Description: "Look at the image and fill out the appropriate fields. We want to be able to tag all the images correctly. Fill out any appropriate tag that you see.",
		Write: func(j *Job) {
			fmt.Printf("%v\n", j)

			var output []string

			for _, i := range j.InputFields {
				output = append(output, i.Value)
				fmt.Println(i.Value)
			}

			for _, i := range j.OutputFields {
				output = append(output, i.Value)
				fmt.Println(i.Value)
			}

			if err := writer.Write(output); err != nil {
				panic(err)
			}
			writer.Flush()
		},
		Tasks: imageUrls(),
	}

	fmt.Println("Loaded and starting")
	serve := HtmlServe{}
	go Serve()

	fmt.Println("Serving")
	var backend Assigner = serve
	NewBatch(image_tasks).Run(backend)

}
