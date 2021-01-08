package ci

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestBuild(t *testing.T) {
	//Temp:::
	//client := http.Client{}
	//res, err := client.Post("http://localhost:8080/api/builds?namespace=tahirali-csc&repoName=hello-app",
	//	"application/json", nil)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//dat, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//vv := make(map[string]interface{})
	//json.Unmarshal(dat, &vv)

	log.Println("Executing task1")
	step := &Step{
		Name:  "task1",
		Image: "alpine:latest",
		Cmd:   []string{"ls"},
	}

	//v := vv["Id"].(float64)
	v := 212
	os.Setenv("TE_HOST_URL", "http://localhost:8080")
	os.Setenv("TE_BUILD_ID", fmt.Sprintf("%d", int64(v)))

	os.Setenv("MOUNT_PATH","/Users/tahir/workspace/build-workspace/app1/")
	os.Setenv("CLAIM_NAME","task-pvc-volume")
	build := NewBuild()
	err := build.Exec(step)
	if err != nil {
		log.Println(err)
	}

	time.Sleep(time.Second * 6)
	log.Println("Executing task2")
	step = &Step{
		Name:  "task2",
		Image: "alpine:latest",
		Cmd:   []string{"/bin/sh", "-c", "ls -al / && date"},
	}
	err = build.Exec(step)
	if err != nil {
		log.Println(err)
	}

	build.Done()

	os.Unsetenv("TE_HOST_URL")
	os.Unsetenv("TE_BUILD_ID")
}
