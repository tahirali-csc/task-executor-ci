package ci

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Step struct {
	Name        string
	Image       string
	Cmd         []string
	CpuLimit    int
	MemoryLimit int
	//BuildID???
	//Host??
}

type stepExec struct {
	Name     string
	Image    string
	Cmd      []string
	CpuLimit int
	Memory   int
	BuildId  int64
}

const hostURL string = "TE_HOST_URL"
const buildID string = "TE_BUILD_ID"

func RunStep(step *Step) error {
	if len(step.Name) == 0 {
		return errors.New("name is missing")
	}

	if len(step.Image) == 0 {
		return errors.New("image is missing")
	}

	if len(step.Cmd) == 0 {
		return errors.New("commands are missing")
	}

	hostURL := os.Getenv(hostURL)
	if len(hostURL) == 0 {
		return errors.New("host URL is missing")
	}

	buildIdEnv := os.Getenv(buildID)
	if len(buildIdEnv) == 0 {
		return errors.New("build Id is missing")
	}

	buildId, err := strconv.ParseInt(buildIdEnv, 10, 64)
	if err != nil {
		return errors.New("invalid build Id")
	}

	stepExec := &stepExec{
		Name:     step.Name,
		BuildId:  buildId,
		Image:    step.Image,
		Cmd:      step.Cmd,
		CpuLimit: step.CpuLimit,
		Memory:   step.MemoryLimit,
	}

	data, err := json.Marshal(stepExec)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s:/api/steps", hostURL)
	client := http.Client{}
	_, err = client.Post(url, "application/json", bytes.NewReader(data))
	return err
}
