package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pcornish/depevent/eventmodel"
	"log"
	"os"
	"path"
	"regexp"
	"time"
)

type deployment struct {
	RepoName    string    `json:"repoName"`
	Commit      string    `json:"commit"`
	Environment string    `json:"environment"`
	EventTime   time.Time `json:"eventTime"`
}

var rootDir string
var environment string
var format string

func main() {
	flag.StringVar(&rootDir, "dir", ".", "root directory")
	flag.StringVar(&environment, "env", "", "filter by environment (optional)")
	flag.StringVar(&format, "output-format", "json", "output format (json or text)")
	flag.Parse()
	if rootDir == "." {
		rootDir, _ = os.Getwd()
	}
	log.Println("fetching files in", rootDir)

	files, err := fetchAllFiles(rootDir)
	if err != nil {
		panic(err)
	}
	events, err := discoverDeployments(files)
	if err != nil {
		panic(err)
	}
	for _, event := range events {
		if environment != "" && event.Environment != environment {
			continue
		}
		printDeployment(event, format)
	}
}

// fetchAllFiles returns all files in the given directory recursively.
func fetchAllFiles(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			subFiles, err := fetchAllFiles(path.Join(dir, entry.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, path.Join(dir, entry.Name()))
		}
	}

	return files, nil
}

func discoverDeployments(files []string) ([]deployment, error) {
	var deployments []deployment
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %s: %w", file, err)
		}
		events, err := getEvents(string(content))
		if err != nil {
			return nil, fmt.Errorf("error getting events from file: %s: %w", file, err)
		}
		if len(events) > 1 {
			log.Println("multiple events in file", file)
		}
		for _, event := range events {
			dep, err := parseDeployment(event)
			if err != nil {
				return nil, fmt.Errorf("error parsing deployment event in file: %s: %w", file, err)
			}
			if dep != nil {
				deployments = append(deployments, *dep)
			}
		}
	}
	return deployments, nil
}

// getEvents returns all events in the given content.
// The content is _sometimes_ a single JSON object, but _sometimes_ a
// concatenation of JSON objects (i.e. not a valid JSON array).
func getEvents(content string) ([]string, error) {
	r := regexp.MustCompile(`}(\s)*{`)
	parts := r.Split(content, -1)
	if len(parts) > 1 {
		for i, part := range parts {
			if i > 0 {
				part = "{" + part
			}
			if i < len(parts)-1 {
				part = part + "}"
			}
			parts[i] = part
		}
	}
	return parts, nil
}

func parseDeployment(event string) (*deployment, error) {
	var d eventmodel.DeploymentEvent
	err := json.Unmarshal([]byte(event), &d)
	if err != nil {
		return nil, err
	}
	if d.ResponsePayload == nil || d.ResponsePayload.Message == nil {
		return nil, nil
	}
	msg := d.ResponsePayload.Message
	if msg.GitRepository == "" || msg.GitCommitSha == "" {
		return nil, nil
	}

	eventTime, err := time.Parse(time.RFC3339, d.RequestPayload.Time)
	if err != nil {
		return nil, fmt.Errorf("error parsing event time: %s: %w", d.RequestPayload.Time, err)
	}

	env := "none"
	for _, param := range msg.StackTemplateParameters {
		if param.ParameterKey == "Environment" {
			env = param.ParameterValue
			break
		}
	}
	dep := deployment{
		EventTime:   eventTime,
		Environment: env,
		RepoName:    msg.GitRepository,
		Commit:      msg.GitCommitSha,
	}
	return &dep, nil
}

func printDeployment(event deployment, format string) {
	switch format {
	case "json":
		bytes, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bytes))

	case "text":
		fmt.Printf("%s %s %s %s\n", event.EventTime.Format(time.RFC3339), event.Environment, event.RepoName, event.Commit)
		return
	}
}
