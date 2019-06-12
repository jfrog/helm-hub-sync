package sync

import (
	"testing"

	"github.com/jfrog/helm-hub-sync/artifactory"
	"github.com/jfrog/helm-hub-sync/github"
	"github.com/stretchr/testify/assert"
)

func TestRemoveFromSlice(t *testing.T) {
	slice := []string{"hello", "world", "jfrog", "artifactory"}
	newSlice := RemoveFromSlice(slice, "bla")
	assert.ElementsMatch(t, newSlice, slice)
	newSlice = RemoveFromSlice(slice, "world")
	assert.NotContains(t, newSlice, "world")
	assert.NotEmpty(t, newSlice)
}

func TestNotInArtifactory(t *testing.T) {
	artiRepos := make(map[string]bool)
	artiRepos["myAwesomeRepo"] = true

	ghRepos := []github.Repo{
		github.Repo{
			Name: "myAwesomeRepo",
			URL:  "http://example.com",
		},
		github.Repo{
			Name: "notExistingAwesomeRepo",
			URL:  "http://example.com",
		},
	}

	repos := NotInArtifactory(artiRepos, ghRepos)
	assert.NotEmpty(t, repos)
	assert.Equal(t, len(repos), 1)
}

func TestNotInGitHub(t *testing.T) {
	ghRepos := make(map[string]bool)
	ghRepos["myAwesomeRepo"] = true

	artiRepos := []artifactory.Repository{
		artifactory.Repository{
			Key: "myAwesomeRepo",
		},
		artifactory.Repository{
			Key: "notExistingAwesomeRepo",
		},
	}

	repos := NotInGitHub(ghRepos, artiRepos)
	assert.NotEmpty(t, repos)
	assert.Equal(t, len(repos), 1)
}
