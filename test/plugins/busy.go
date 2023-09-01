package main

import (
	"context"
	plugins2 "github.com/compliance-framework/assessment-runtime/internal/plugins"
	structpb "google.golang.org/protobuf/types/known/structpb"
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

func (p *BusyPlugin) Execute(_ *plugins2.ActionInput) (*plugins2.ActionOutput, error) {
	time.Sleep(p.duration)
	data := map[string]interface{}{
		"message": p.message,
	}
	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, err
	}
	return &plugins2.ActionOutput{
		ResultData: s,
	}, nil
}

func (p *BusyPlugin) Shutdown(context.Context) error {
	return nil
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pluginsMap := make(map[string]plugins2.Plugin)
	pluginsMap["busy-plugin"] = &BusyPlugin{
		duration: time.Duration(r.Intn(10)) * time.Second,
		message:  "Busy Plugin completed",
	}
	plugins2.Register(pluginsMap)
}
