package ci

import (
	"log"
	"os"
	"testing"
)

func TestBuild(t *testing.T) {
	step := &Step{
		Name:  "task1",
		Image: "alpine:latest",
		Cmd:   []string{"ls"},
	}
	os.Setenv("TE_HOST_URL", "http://localhost:8080")
	os.Setenv("TE_BUILD_ID", "104")
	build := NewBuild()
	err := build.Exec(step)
	if err != nil {
		log.Println(err)
	}

	build.Done()

	os.Unsetenv("TE_HOST_URL")
	os.Unsetenv("TE_BUILD_ID")
}
