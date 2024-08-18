package claiagent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lukasmetzner/clai/pkg/database"
	"github.com/lukasmetzner/clai/pkg/k8s"
	"github.com/lukasmetzner/clai/pkg/models"
	"github.com/lukasmetzner/clai/pkg/mq"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/remotecommand"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Start() {
	log.Println("Starting clai-agent...")

	// Init Kubernetes
	k8s.Init()

	// Init RabbitMQ
	mq.InitMQ()

	// Init Postgres
	database.InitDB()

	msgs, err := mq.Channel.Consume(
		mq.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("%s", err)
	}

	for msg := range msgs {
		log.Printf("Received a new job %s", msg.Body)
		job := models.Job{}
		if err := json.Unmarshal(msg.Body, &job); err != nil {
			log.Printf("%s", err)
			continue
		}

		go executeJob(job)
	}
}

func waitForPodRunning(podName string) error {
	return wait.PollUntilContextTimeout(context.TODO(), time.Duration(5), wait.ForeverTestTimeout, false, func(ctx context.Context) (bool, error) {
		pod, err := k8s.PodClient.Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return pod.Status.Phase == v1.PodRunning, nil
	})
}

// Function to execute a command inside the Pod
func execCommandInPod(podName string, command []string, stdinBuffer *bytes.Buffer, stdoutBuffer *bytes.Buffer, stderrBuffer *bytes.Buffer) error {
	req := k8s.ClientSet.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(k8s.Namespace).
		SubResource("exec").
		Param("container", k8s.ContainerName).
		Param("stdout", "true").
		Param("stderr", "true").
		Param("stdin", "true").
		Param("tty", "false")

	for _, cmd := range command {
		req = req.Param("command", cmd)
	}

	exec, err := remotecommand.NewSPDYExecutor(k8s.Config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("could not initialize executor: %v", err)
	}

	streamOptions := remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	}

	// TODO: Create buffer config struct
	if stdinBuffer != nil {
		streamOptions.Stdin = stdinBuffer
	}
	if stdoutBuffer != nil {
		streamOptions.Stdout = stdoutBuffer
	}
	if stderrBuffer != nil {
		streamOptions.Stderr = stderrBuffer
	}

	return exec.StreamWithContext(context.TODO(), streamOptions)
}

func executeJob(job models.Job) {
	idStr := strconv.FormatUint(uint64(job.ID), 10)
	k8sPod := k8s.GetPod(idStr)

	result, err := k8s.PodClient.Create(context.TODO(), &k8sPod, metav1.CreateOptions{})
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	podName := result.GetObjectMeta().GetName()

	if err := waitForPodRunning(podName); err != nil {
		log.Printf("%s\n", err)
	}

	log.Println("Pod is running")

	// Load requirements.txt into bytes buffer
	var buffer bytes.Buffer
	buffer.WriteString(job.Requirements)

	// Copy requirements.txt to remote pod
	command := []string{"/bin/sh", "-c", "cat - > /tmp/requirements.txt"}
	if err := execCommandInPod(podName, command, &buffer, nil, nil); err != nil {
		log.Printf("%s\n", err)
	}
	log.Println("Uploaded requirements.txt to remote pod")

	// Pip install requirements.txt
	command = []string{"/bin/sh", "-c", "pip3 install -r /tmp/requirements.txt"}
	if err := execCommandInPod(podName, command, nil, nil, nil); err != nil {
		log.Printf("%s\n", err)
	}
	log.Println("Installed requirements.txt into remote pod")

	// Load script.py into bytes buffer
	var scriptBuffer bytes.Buffer
	scriptBuffer.WriteString(job.Script)

	// Copy script.py to remote pod
	command = []string{"/bin/sh", "-c", "cat - > /tmp/script.py"}
	if err := execCommandInPod(podName, command, &scriptBuffer, nil, nil); err != nil {
		log.Printf("%s\n", err)
	}
	log.Println("Uploaded script.py to remote pod")

	// Execute python script

	var stdoutBuffer, stderrBuffer bytes.Buffer

	command = []string{"/bin/sh", "-c", "python3 /tmp/script.py"}
	if err := execCommandInPod(podName, command, nil, &stdoutBuffer, &stderrBuffer); err != nil {
		log.Printf("%s\n", err)
	}
	log.Println("Executed script.py in remote pod")

	database.AppendJobOutput(job.ID, &stdoutBuffer, &stderrBuffer)

	// Delete pod after script execution
	if err := k8s.PodClient.Delete(context.TODO(), podName, metav1.DeleteOptions{}); err != nil {
		log.Printf("%s\n", err)
	}
	log.Println("Deleted pod")
}
