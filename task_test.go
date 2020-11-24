package ci

import (
	"log"
	"os"
	"testing"
)

func TestTask(t *testing.T) {
	step := &Step{
		Name:  "task1",
		Image: "alpine:latest",
		Cmd:   []string{"ls"},
	}
	os.Setenv("TE_HOST_URL","http://localhost:8080")
	os.Setenv("TE_BUILD_ID","100")
	err := RunStep(step)
	if err != nil {
		log.Println(err)
	}

	os.Unsetenv("TE_HOST_URL")
	os.Unsetenv("TE_BUILD_ID")
}
