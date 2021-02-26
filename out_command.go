package resource

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v28/github"
)

type OutCommand struct {
	github GitHub
	writer io.Writer
}

func NewOutCommand(github GitHub, writer io.Writer) *OutCommand {
	return &OutCommand{
		github: github,
		writer: writer,
	}
}

func (c *OutCommand) Run(sourceDir string, request OutRequest) (OutResponse, error) {
	if request.Params.ID == nil {
		return OutResponse{}, errors.New("id is a required parameter")
	}
	if request.Params.State == nil {
		return OutResponse{}, errors.New("state is a required parameter")
	}

	idInt, err := strconv.ParseInt(*request.Params.ID, 10, 64)
	if err != nil {
		return OutResponse{}, err
	}
	fmt.Fprintln(c.writer, "getting deployment")
	deployment, err := c.github.GetDeployment(idInt)
	if err != nil {
		return OutResponse{}, err
	}

	logURL := fmt.Sprintf("%s/teams/%s/pipelines/%s/jobs/%s/builds/%s",
		os.Getenv("ATC_EXTERNAL_URL"),
		os.Getenv("BUILD_TEAM_NAME"),
		os.Getenv("BUILD_PIPELINE_NAME"),
		os.Getenv("BUILD_JOB_NAME"),
		os.Getenv("BUILD_NAME"),
	)

	newStatus := &github.DeploymentStatusRequest{
		State:          request.Params.State,
		Description:    request.Params.Description,
		EnvironmentURL: request.Params.EnvironmentURL,
		LogURL:         github.String(logURL),
	}

	fmt.Fprintln(c.writer, "creating deployment status")
	_, err = c.github.CreateDeploymentStatus(*deployment.ID, newStatus)
	if err != nil {
		return OutResponse{}, err
	}

	fmt.Fprintln(c.writer, "getting deployment statuses list")
	statuses, err := c.github.ListDeploymentStatuses(*deployment.ID)
	if err != nil {
		return OutResponse{}, err
	}

	latestStatus := ""
	if len(statuses) > 0 {
		latestStatus = *statuses[0].State
	}

	return OutResponse{
		Version: Version{
			ID:         *request.Params.ID,
			LastStatus: latestStatus,
		},
		Metadata: metadataFromDeployment(deployment, statuses),
	}, nil
}

func (c *OutCommand) fileContents(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(contents)), nil
}
