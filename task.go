package ci

import "log"

type Step struct {
	Image string
	Cmd   []string
}

func RunStep(step *Step) {
	log.Println("Running::", step)
}
