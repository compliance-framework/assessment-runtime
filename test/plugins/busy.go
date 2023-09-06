package main

import (
	"context"
	. "github.com/compliance-framework/assessment-runtime/internal/plugin"
	"google.golang.org/protobuf/types/known/structpb"
	"math/rand"
	"time"
)

type BusyPlugin struct {
	duration time.Duration
	message  string
}

func (p *BusyPlugin) Init() error {
	return nil
}

func (p *BusyPlugin) Execute(_ *ActionInput) (*ActionOutput, error) {
	time.Sleep(p.duration)
	data := map[string]interface{}{
		"message": p.message,
	}
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}
	return &ActionOutput{
		ResultData: s,
	}, nil
}

func (p *BusyPlugin) Shutdown(context.Context) error {
	return nil
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pluginsMap := make(map[string]Plugin)
	pluginsMap["busy-plugin"] = &BusyPlugin{
		duration: time.Duration(r.Intn(10)) * time.Second,
		message:  "Busy Plugin completed",
	}
	Register(pluginsMap)
}
