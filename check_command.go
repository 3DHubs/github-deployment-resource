package resource

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

type CheckCommand struct {
	github GitHub
	writer io.Writer
}

func NewCheckCommand(github GitHub, writer io.Writer) *CheckCommand {
	return &CheckCommand{
		github: github,
		writer: writer,
	}
}

func (c *CheckCommand) Run(request CheckRequest) ([]Version, error) {
	fmt.Fprintln(c.writer, "getting deployments list")
	deployments, err := c.github.ListDeployments()

	if err != nil {
		return []Version{}, err
	}

	var latestVersions []Version

	for _, deployment := range deployments {
		if len(request.Source.Environments) > 0 {
			found := false
			for _, env := range request.Source.Environments {
				if env == *deployment.Environment {
					found = true
				}
			}

			if !found {
				continue
			}
		}

		id := *deployment.ID
		statuses, err := c.github.ListDeploymentStatuses(id)
		if err != nil {
			return []Version{}, err
		}

		// Assume first returned status is the latest one
		latestStatus := ""
		latestStatusID := ""
		if len(statuses) > 0 {
			latestStatus = *statuses[0].State
			latestStatusID = strconv.FormatInt(*statuses[0].ID, 10)
		}

		lastID, err := strconv.ParseInt(request.Version.ID, 10, 64)
		if err != nil || id >= lastID {
			latestVersions = append(latestVersions, Version{
				ID:         strconv.FormatInt(id, 10),
				LastStatus: latestStatus,
				StatusID:   latestStatusID,
			})
		}
	}

	if len(latestVersions) == 0 {
		return []Version{}, nil
	}

	sort.Slice(latestVersions[:], func(i, j int) bool {
		iID, _ := strconv.Atoi(latestVersions[i].ID)
		jID, _ := strconv.Atoi(latestVersions[j].ID)
		if iID == jID {
			iStatusID, _ := strconv.Atoi(latestVersions[i].StatusID)
			jStatusID, _ := strconv.Atoi(latestVersions[j].StatusID)
			return iStatusID < jStatusID
		}
		return iID < jID
	})

	latestVersion := latestVersions[len(latestVersions)-1]

	if request.Version.ID == "" {
		return []Version{
			latestVersion,
		}, nil
	}

	return latestVersions, nil
}
