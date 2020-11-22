package ci

import (
	"context"
	"fmt"
	"github.com/tahirali-csc/task-executor-engine/engine"
	"github.com/tahirali-csc/task-executor-engine/engine/kube"
	"log"
	"os"
)

type runner struct {
}

func NewRunner() *runner {
	return nil
}

func isKubernetes() bool {
	return os.Getenv("KUBERNETES_SERVICE_HOST") != ""
}

func (runner *runner) Run(step *Step, stepId int64) {

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

	spec := &engine.Spec{
		Image:   step.Image,
		Command: step.Cmd,
		Args:    step.Args,
		Metadata: engine.Metadata{
			Namespace: "default",
			//TODO: Can add more randomization
			UID: fmt.Sprintf("te-step-%d", stepId),
		},
	}

	err := kubeEngine.Start(context.Background(), spec)
	if err != nil {
		log.Println(err)
		return
	}

	kubeEngine.Wait(context.Background(), spec)
	log.Println("I am done")
}
