package crowdflow

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
