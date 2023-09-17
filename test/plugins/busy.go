package main

import (
	"fmt"
	. "github.com/compliance-framework/assessment-runtime/internal/provider"
	"google.golang.org/protobuf/types/known/structpb"
	"strconv"
)

type BusyPlugin struct {
	message string
}

func (p *BusyPlugin) EvaluateSelector(_ *SubjectSelector) (*SubjectList, error) {
	subjects := make([]*Subject, 0)
	for i := 0; i < 3; i++ {
		subjects = append(subjects, &Subject{Id: strconv.Itoa(i)})
	}
	list := &SubjectList{
		Subjects: subjects,
	}
	return list, nil
}

func (p *BusyPlugin) Execute(in *ActionInput) (*ActionOutput, error) {
	data := map[string]interface{}{
		"message": fmt.Sprintf("busy provider completed for subject: %s %s", in.Subject.Id, p.message),
	}
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}
	return &ActionOutput{
		ResultData: s,
	}, nil
}

func main() {
	Register(&BusyPlugin{
		message: "busy provider completed",
	})
}
