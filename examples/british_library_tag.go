package main

import (
	. "github.com/abhiyerra/workmachine/app"
)

type ImageTagging struct {
	ImageUrl             InputField  `crowd_desc:"Use this information for the data below." crowd_id:"image_url" crowd_type:"image"`
	Tags                 OutputField `crowd_desc:"List all the relavent tags for this image." crowd_id:"tags"`
	TextInImage          OutputField `crowd_desc:"Put any text that appears on the image here. One line per block of text." crowd_id:"text_in_image"`
	IsCorrectOrientation OutputField `crowd_desc:"Is the image in the correct orientation?" crowd_id:"is_correct_orientation"`
}

func main() {
	image_tasks := Task{
		Title:       "Tag the appropriate images",
		Description: "Look at the image and fill out the appropriate fields. We want to be able to tag all the images correctly. Fill out any appropriate tag that you see.",
		Tasks: []ImageTagging{
			ImageTagging{
				ImageUrl: "http://www.flickr.com/photos/britishlibrary/11115401504",
			},
			ImageTagging{
				ImageUrl: "wwwphotos/britishlibrary/11115401504",
			},
		},
	}

	serve := HtmlServe{}
	go Serve()

	var backend Assigner = serve
	NewBatch(image_tasks).Run(backend)

}
