package resource_test

import (
	"github.com/google/go-github/v28/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGithubDeploymentResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GithubDeploymentResource Suite")
}

func newDeployment(id int64) *github.Deployment {
	return &github.Deployment{
		ID: github.Int64(id),
	}
}

func newDeploymentStatus(id int64, state string) *github.DeploymentStatus {
	return &github.DeploymentStatus{
		ID:    github.Int64(id),
		State: github.String(state),
	}
}

func newDeploymentWithEnvironment(id int64, env string) *github.Deployment {
	return &github.Deployment{
		ID:          github.Int64(id),
		Environment: &env,
	}
}
