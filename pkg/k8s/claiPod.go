package k8s

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ContainerName string = "python"

func GetPod(name string) v1.Pod {
	job := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-job", name),
		},
		Spec: v1.PodSpec{
			RestartPolicy: "Never",
			Containers: []v1.Container{
				{
					Name:    ContainerName,
					Image:   "python:3.12",
					Command: []string{"sleep", "3600"}, // Sleep for One Hour
				},
			},
		},
	}

	return job
}
