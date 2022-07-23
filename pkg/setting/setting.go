package setting

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Jedi struct {
	Keyword string `yaml:"keyword"`
}

func (j *Jedi) GetConfTerm() *Jedi {
	file, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal("Failed to load ./config/config.yaml:", err)
	}
	err = yaml.Unmarshal(file, j)
	if err != nil {
		log.Fatal("Failed to unmarshall:", err)
	}
	return j
}
