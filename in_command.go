package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type InCommand struct {
	github GitHub
	writer io.Writer
}

func NewInCommand(github GitHub, writer io.Writer) *InCommand {
	return &InCommand{
		github: github,
		writer: writer,
	}
}

func (c *InCommand) Run(destDir string, request InRequest) (InResponse, error) {
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return InResponse{}, err
	}

	id, _ := strconv.ParseInt(request.Version.ID, 10, 64)
	fmt.Fprintln(c.writer, "getting deployment")
	deployment, err := c.github.GetDeployment(id)
	if err != nil {
		return InResponse{}, err
	}

	if deployment == nil {
		return InResponse{}, errors.New("no deployment")
	}

	idPath := filepath.Join(destDir, "id")
	err = ioutil.WriteFile(idPath, []byte(request.Version.ID), 0644)
	if err != nil {
		return InResponse{}, err
	}

	statusPath := filepath.Join(destDir, "status")
	err = ioutil.WriteFile(statusPath, []byte(request.Version.LastStatus), 0644)
	if err != nil {
		return InResponse{}, err
	}

	refPath := filepath.Join(destDir, "ref")
	err = ioutil.WriteFile(refPath, []byte(*deployment.Ref), 0644)
	if err != nil {
		return InResponse{}, err
	}

	shaPath := filepath.Join(destDir, "sha")
	err = ioutil.WriteFile(shaPath, []byte(*deployment.SHA), 0644)
	if err != nil {
		return InResponse{}, err
	}

	if deployment.Task != nil {
		taskPath := filepath.Join(destDir, "task")
		err = ioutil.WriteFile(taskPath, []byte(*deployment.Task), 0644)
		if err != nil {
			return InResponse{}, err
		}
	}

	if deployment.Environment != nil {
		envPath := filepath.Join(destDir, "environment")
		err = ioutil.WriteFile(envPath, []byte(*deployment.Environment), 0644)
		if err != nil {
			return InResponse{}, err
		}
	}

	if deployment.Description != nil {
		descPath := filepath.Join(destDir, "description")
		err = ioutil.WriteFile(descPath, []byte(*deployment.Description), 0644)
		if err != nil {
			return InResponse{}, err
		}
	}

	// Save the whole deployment too I guess.
	deploymentPath := filepath.Join(destDir, "deployment.json")
	deploymentJSON, _ := json.Marshal(deployment)
	err = ioutil.WriteFile(deploymentPath, deploymentJSON, 0644)
	if err != nil {
		return InResponse{}, err
	}

	statuses, err := c.github.ListDeploymentStatuses(id)
	if err != nil {
		return InResponse{}, err
	}

	// Save all the metadata why not.
	metadata := metadataFromDeployment(deployment, statuses)
	metadataPath := filepath.Join(destDir, "metadata.json")
	metadataJSON, _ := json.Marshal(metadata)
	err = ioutil.WriteFile(metadataPath, metadataJSON, 0644)
	if err != nil {
		return InResponse{}, err
	}

	return InResponse{
		Version:  request.Version,
		Metadata: metadata,
	}, nil
}
