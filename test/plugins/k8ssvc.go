package main

import (
	"context"
	"fmt"
	"time"

	. "github.com/compliance-framework/assessment-runtime/provider"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sSvcPlugin struct {
	message string
}

func (p *K8sSvcPlugin) Evaluate(ei *EvaluateInput) (*EvaluateResult, error) {
	subjects := make([]*Subject, 0)
	clientset, err := prepareClient(ei.Selector.Labels["host"], ei.Selector.Labels["token"])
	if err != nil {
		return nil, err
	}
	d := time.Now().Add(time.Second * 10)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	services, err := clientset.CoreV1().Services(ei.Selector.Labels["namespace"]).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, svc := range services.Items {
		subjects = append(subjects, &Subject{Id: svc.Name})
	}
	er := &EvaluateResult{
		Subjects: subjects,
	}
	er.Props = make(map[string]string)
	er.Props["host"] = ei.Selector.Labels["host"]
	er.Props["token"] = ei.Selector.Labels["token"]
	return er, nil
}

func (p *K8sSvcPlugin) Execute(in *ExecuteInput) (*ExecuteResult, error) {
	clientset, err := prepareClient(in.Props["host"], in.Props["token"])
	if err != nil {
		return nil, err
	}
	d := time.Now().Add(time.Second * 10)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	service, err := clientset.CoreV1().Services(in.Props["namespace"]).Get(ctx, in.Subject.Id, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	found := false
	logEntries := make([]*LogEntry, 0)
	for _, port := range service.Spec.Ports {
		le := &LogEntry{
				Timestamp: time.Now().Format("2006-02-01T15:04:05Z"),
				Type:      LogType_DEBUG,
				Details:   fmt.Sprintf("Service %s exposes port %d", service.Name, port.Port),
			}
		logEntries = append(logEntries, le)
		if port.Port == 80 {
			found = true
			break
		}
	}
	er := &ExecuteResult{}

	if !found {
		er.Status = ExecutionStatus_SUCCESS
	} else {
		er.Status = ExecutionStatus_FAILURE
		observations := make([]*Observation, 0)

		obs := &Observation{
			SubjectId:   in.Subject.Id,
			Title:       "Service exposes port 80",
			Description: "The automated assessment tool detected that the the service exposes the default port 80.",
			Collected:   time.Now().Format("2006-02-01T15:04:05Z"),
			Expires:     time.Now().Add(time.Hour * 24).Format("2006-02-01T15:04:05Z"),
			Links: []*Link{
				{
					Rel:  "related",
					Href: "https://kubernetes.io/docs/concepts/services-networking/service/",
				},
			},
			Props: []*Property{
				{
					Name:  "Risk Level",
					Value: "High",
				},
				{
					Name:  "Recommendation",
					Value: "Expose ports other than default http port.",
				},
			},
			RelevantEvidence: []*Evidence{
				{
					Description: "Automated tool log indicating the exposure of the default http port.",
				},
			},
			Remarks: "Immediate action required to mitigate potential data breaches.",
		}
		er.Observations = append(observations, obs)
	}
	er.Logs = logEntries
	return er, nil
}

func main() {
	Register(&K8sSvcPlugin{
		message: "K8sSvcPlugin provider completed",
	})
}

func prepareClient(host, token string) (*kubernetes.Clientset, error) {
	config := &rest.Config{
		Host:        host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
