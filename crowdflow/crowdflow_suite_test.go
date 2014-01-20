package crowdflow_test

import (
	. "github.com/abhiyerra/workmachine/crowdflow"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	//	"net/http/httptest"
	//"fmt"
	"testing"
)

type Research struct {
	Name      InputField  `work_desc:"The name of the person to research." work_id:"name"`
	Born      OutputField `work_desc:"Find out when the person was born" work_id:"born"`
	IsPainter OutputField `work_desc:"Was the person a painter?" work_id:"is_painter" work_type:"checkbox"`
}

func TestCrowdflow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crowdflow Suite")
}

/*
func TestCrowdflow2(t *testing.T) {

	// serve := HtmlServe{}
	// //	go HtmlServer()

	// var backend Assigner = serve
	// NewBatch(research).Run(backend)

	// fmt.Printf("Length of assignments %d", len(assignments))
	// for _, assign := range assignments {
	// 	assign.JobsChan <- Job{}
	// }

	/* CrowdFlow test one

	   - Make sure that the jobs are created
	   - Make sure that the we can get the assignment.
	   - Make sure that we can post the assignment
	* /
}
*/
