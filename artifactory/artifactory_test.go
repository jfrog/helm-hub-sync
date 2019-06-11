package artifactory

import (
	"fmt"
	"os"
	"testing"
	"net/http"

	"github.com/stretchr/testify/assert"
)

var (
	artifactoryFailHost = os.Getenv("ARTIFACTORY_DUMMY_SERVER")
	artifactoryTestHost = os.Getenv("ARTIFACTORY_LIVE_SERVER")
	artifactoryTestAuth = os.Getenv("ARTIFACTORY_BASIC_CREDS")
)

func TestUnmarshalRepositories(t *testing.T) {
	// unmarshalRepositories fails with improper input
	repos, err := unmarshalRepositories([]byte(""))
	assert.Error(t, err)
	assert.Nil(t, repos)

	// unmarshalRepositories succeeds with proper inputs
	repos, err = unmarshalRepositories([]byte(`[ { "key": "helm", "type": "VIRTUAL", "url": "http://jfrog.local:80/artifactory/helm", "packageType": "Helm" } ]`))
	assert.NoError(t, err)
	assert.Equal(t, len(repos), 1)
	assert.Equal(t, repos[0].PackageType, "Helm")
}

func TestCallArtifactory(t *testing.T) {
	// callArtifactory fails with an incorrect host
	_, err := callArtifactory(artifactoryFailHost, "authorization", "", http.MethodGet, nil)
	assert.Error(t, err)

	// callArtifactory succeeds with proper host
	_, err = callArtifactory(fmt.Sprintf("%s/artifactory%s", artifactoryTestHost, allRemoteRepositoriesURL), "authorization", "", http.MethodGet, nil)
	assert.NoError(t, err)
}

func TestGetAllRepositories(t *testing.T) {
	// GetRepositories fails with improper host
	_, err := GetRepositories(artifactoryTestHost, "authorization", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Artifactory returned non-OK statuscode [404]")

	// GetRepositories succeeds with proper host
	_, err = GetRepositories(fmt.Sprintf("%s/artifactory", artifactoryTestHost), "authorization", "")
	assert.NoError(t, err)
}

func TestGetRepository(t *testing.T) {
	// GetRepository fails with improper host
	_, err := GetRepository(fmt.Sprintf("%s/artifactory", artifactoryTestHost), "authorization", "", "bla")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Artifactory returned non-OK statuscode [400]")

	// GetRepository succeeds with proper host
	_, err = GetRepository(fmt.Sprintf("%s/artifactory", artifactoryTestHost), "authorization", "", "helm")
	assert.NoError(t, err)
}

func TestCreateRepository(t *testing.T) {
	// CreateRepository fails with improper authentication
	err := CreateRepository(fmt.Sprintf("%s/artifactory", artifactoryTestHost), "authorization", "bla", "myRemoteRepo","https://myrepo.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Artifactory returned non-OK statuscode [401]")

	// GetRepository succeeds with proper host
	err = CreateRepository(fmt.Sprintf("%s/artifactory", artifactoryTestHost), "authorization", artifactoryTestAuth, "myRemoteRepo","https://myrepo.com")
	assert.NoError(t, err)
}

func TestUpdateRepository(t *testing.T) {
	// UpdateRepository fails with improper authentication
	err := UpdateRepository(fmt.Sprintf("%s/artifactory", artifactoryTestHost), "authorization", "bla", "myRemoteRepo",Repository{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Artifactory returned non-OK statuscode [401]")

	// UpdateRepository succeeds with proper host
	repo := Repository{
		Key: "helm",
		PackageType: "helm",
		Repositories: []string{"helm-remote","myRemoteRepo"},
		Rclass: "virtual",
	}
	err = UpdateRepository(fmt.Sprintf("%s/artifactory", artifactoryTestHost), "authorization", artifactoryTestAuth, "helm",repo)
	assert.NoError(t, err)
}

func TestDeleteRepository(t *testing.T) {
	// DeleteRepository fails with improper authentication
	err := DeleteRepository(fmt.Sprintf("%s/artifactory", artifactoryTestHost), "authorization", "bla", "myRemoteRepo")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Artifactory returned non-OK statuscode [401]")

	// DeleteRepository succeeds with proper host
	err = DeleteRepository(fmt.Sprintf("%s/artifactory", artifactoryTestHost), "authorization", artifactoryTestAuth, "myRemoteRepo")
	assert.NoError(t, err)
}