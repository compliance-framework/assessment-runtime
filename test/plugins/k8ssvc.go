package main

import (
	"context"
	. "github.com/compliance-framework/assessment-runtime/internal/provider"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

type K8sSvcPlugin struct {
	message string
}

func (p *K8sSvcPlugin) EvaluateSelector(ss *SubjectSelector) (*SubjectList, error) {
	subjects := make([]*Subject, 0)
	clientset, err := prepareClient(ss.Labels["host"], ss.Labels["token"])
	if err != nil {
		return nil, err
	}
	d := time.Now().Add(time.Second * 10)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	services, err := clientset.CoreV1().Services(ss.Labels["namespace"]).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, svc := range services.Items {
		subjects = append(subjects, &Subject{Id: svc.Name})
	}
	list := &SubjectList{
		Subjects: subjects,
	}
	return list, nil
}

func (p *K8sSvcPlugin) Execute(in *JobInput) (*JobResult, error) {
	clientset, err := prepareClient(in.Parameters["host"], in.Parameters["token"])
	if err != nil {
		return nil, err
	}
	d := time.Now().Add(time.Second * 10)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	service, err := clientset.CoreV1().Services(in.Parameters["namespace"]).Get(ctx, in.Id, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	found := false
	for _, port := range service.Spec.Ports {
		if port.Port == 80 {
			found = true
			break
		}
	}
	jr := &JobResult{}

	if !found {
		jr.State = "Pass"
	} else {
		jr.State = "Fail"
		observations := make([]*Observation, 0)

		obs := &Observation{
			SubjectId:   in.SubjectId,
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
			Uuid:    uuid.New().String(),
		}
		jr.Observations = append(observations, obs)
	}

	return jr, nil
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
