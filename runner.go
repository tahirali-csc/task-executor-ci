package ci

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tahirali-csc/task-executor-engine/engine"
	"github.com/tahirali-csc/task-executor-engine/engine/kube"
)

type runner struct {
}

func NewRunner() *runner {
	return nil
}

func isKubernetes() bool {
	return os.Getenv("KUBERNETES_SERVICE_HOST") != ""
}

func (runner *runner) Run(step *Step, stepId int64, logsChann chan []byte, doneChann chan bool) {

	var (
		kubeEngine engine.Engine
		initError  error
	)

	if isKubernetes() {
		kubeEngine, initError = kube.NewFile("", "", "")
	} else {
		userHome, _ := os.UserHomeDir()
		kubeConfig := userHome + "/.kube/config"
		kubeEngine, initError = kube.NewFile("", kubeConfig, "")
	}

	if initError != nil {
		log.Println(initError)
		return
	}

	mountPath := os.Getenv("MOUNT_PATH")
	claimName := os.Getenv("CLAIM_NAME")

	spec := &engine.Spec{
		Image:   step.Image,
		Command: step.Cmd,
		Args:    step.Args,
		Metadata: engine.Metadata{
			Namespace: "default",
			//TODO: Can add more randomization
			UID: fmt.Sprintf("te-step-%d", stepId),
		},
		Volumes: []engine.VolumeMount{
			{
				Name:      "logs-drive",
				ClaimName: claimName,
				MountPath: mountPath,
			},
		},
	}

	err := kubeEngine.Start(context.Background(), spec)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		r, err := kubeEngine.Tail(context.Background(), spec)
		if err != nil {
			return
		}

		br := bufio.NewReader(r)
		for {
			line, _, err := br.ReadLine()
			if err != nil {
				doneChann <- true
				return
			}

			logsChann <- line
		}
	}()

	kubeEngine.Wait(context.Background(), spec)

	log.Println("Finished running step ", step.Name)
}
