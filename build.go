package ci

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"encoding/json"
	"errors"
)

type Build struct {
}

func NewBuild() *Build {
	return &Build{}
}

const hostURL string = "TE_HOST_URL"
const buildID string = "TE_BUILD_ID"

func uploadLog(baseURL string, file string, stepId int64) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	const chunkSize = 100

	br := bufio.NewReader(f)
	pr, pw := io.Pipe()

	url := fmt.Sprintf("%s:/api/steps/%d/logs", baseURL, stepId)
	req, err := http.NewRequest(http.MethodPost, url, pr)
	if err != nil {
		log.Println(err)
		return
	}

	req.TransferEncoding = []string{"chunked"}

	go func() {
		defer pw.Close()
		for {
			_, err := io.CopyN(pw, br, chunkSize)
			if err != nil {
				return
			}
		}
	}()

	client := http.Client{}
	_, err = client.Do(req)

	if err != nil {
		log.Println(err)
	}

	log.Println("Done loading.")
}

func (b *Build) Exec(step *Step) error {
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
		Args:     step.Args,
		CpuLimit: step.CpuLimit,
		Memory:   step.MemoryLimit,
	}

	data, err := json.Marshal(stepExec)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/steps", hostURL)
	client := http.Client{}
	res, err := client.Post(url, "application/json", bytes.NewReader(data))

	if err != nil {
		return err
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	stepInfo := make(map[string]interface{})
	err = json.Unmarshal(data, &stepInfo)
	if err != nil {
		return err
	}

	logsChann := make(chan []byte)
	doneChann := make(chan bool)
	logFilePath := make(chan string)

	go func() {
		mountPath := os.Getenv("MOUNT_PATH")

		currTime := time.Now().Unix()
		filePath := path.Join(mountPath, fmt.Sprintf("%s-%d.log", step.Name, currTime))
		f, err := os.Create(filePath)
		// defer os.Remove(filePath)

		if err != nil {
			log.Println("Can not create file::::", err)
			return
		}

		f.Chmod(0700)
		bw := bufio.NewWriter(f)
	loop:
		for {
			select {
			case line := <-logsChann:
				bw.WriteString(string(line) + "\n")
				bw.Flush()

			case <-doneChann:
				break loop
			}
		}

		logFilePath <- filePath
	}()

	stepId := int64(stepInfo["Id"].(float64))
	runner := NewRunner()
	runner.Run(step, stepId, logsChann, doneChann)

	file := <-logFilePath
	uploadLog(hostURL, file, stepId)
	os.Remove(file)

	//TODO:
	url = fmt.Sprintf("%s:/api/steps/%d/status/%s", hostURL, stepId, "Finished")
	res, err = client.Post(url, "application/json", bytes.NewReader(data))

	if err != nil {
		return err
	}

	return nil
}

func (b *Build) Done() error {
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

	client := http.Client{}
	url := fmt.Sprintf("%s/api/builds/%d/status/%s", hostURL, buildId, "Finished")
	_, err = client.Post(url, "application/json", nil)

	if err != nil {
		return err
	}

	return nil
}
