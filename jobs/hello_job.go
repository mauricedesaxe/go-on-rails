package jobs

import "log"

type HelloJob struct {
	Greeting string
	Name     string
}

func (e HelloJob) Execute() error {
	log.Printf("%s, %s!", e.Greeting, e.Name)
	return nil
}
