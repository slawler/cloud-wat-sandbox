package jobs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kendru/darwin/go/depgraph"
	"gopkg.in/yaml.v3"
)

type Dag map[int][]Directive

type Directive struct {
	Plugin  string `json:"plugin"`
	Payload string `json:"payload"`
	Status  string `json:"status"`
}

func ParseJob(s string) (Job, error) {
	var job Job

	jobConfig, err := os.ReadFile(s)
	if err != nil {
		return job, err
	}

	err = yaml.Unmarshal(jobConfig, &job)
	if err != nil {
		return job, err
	}

	return job, nil
}

func (j Job) ParsePayload(pluginName string) (string, error) {
	var payloadFile string
	for _, plugin := range j.Job.Plugins {
		if pluginName == plugin.Plugin.Name {
			payloadFile = filepath.Join(j.Job.Config.Bucket,
				j.Job.Config.PrefixIn, plugin.Plugin.Payload)
		}
	}

	jsonFile, err := os.Open(payloadFile)
	if err != nil {
		return "", nil
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil

}

func (j Job) Plan() ([]byte, error) {

	// inputDir := filepath.Join(j.Job.Config.Bucket, j.Job.Config.PrefixIn)
	// outputDir := filepath.Join(j.Job.Config.Bucket, j.Job.Config.PrefixOut)

	programOrder := depgraph.New()
	dag := make(Dag)

	for _, plugin := range j.Job.Plugins {
		pluginName := plugin.Plugin.Name
		for _, pluginDep := range plugin.Plugin.DependsOn {
			programOrder.DependOn(pluginName, pluginDep)
		}
	}

	for i, taskList := range programOrder.TopoSortedLayers() {

		if _, ok := dag[i]; ok {
			continue
		} else {
			dag[i] = make([]Directive, 0)
		}

		for _, task := range taskList {
			for _, plugin := range j.Job.Plugins {
				if task == plugin.Plugin.Name {
					payloadFile := plugin.Plugin.Payload
					directive := Directive{Plugin: task, Payload: payloadFile, Status: "blocked"}

					dag[i] = append(dag[i], directive)
				}
			}
		}
	}

	fmt.Println("dag", dag)
	dagJson, err := json.Marshal(dag)
	if err != nil {
		return nil, err
	}

	return dagJson, nil
}
