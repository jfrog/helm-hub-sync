package github

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

// HelmRepo ...
type HelmRepo struct {
	Sync Sync
}

// Sync ...
type Sync struct {
	Repos []Repo
}

// Repo is the construct that represents a Helm repository in the Helm Hub configuration
type Repo struct {
	Name string
	URL  string
}

const (
	chartURL = "https://raw.githubusercontent.com/helm/hub/master/config/repo-values.yaml"
)

func getChartConfig() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, chartURL, nil)
	if err != nil {
		return nil, err
	}

	// TODO: change this to not use the DefaultClient
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub returned non-OK statuscode [%d]: %s", res.StatusCode, string(body))
	}

	return body, nil
}

func marshalRepos(data []byte) ([]Repo, error) {
	repos := HelmRepo{}

	err := yaml.Unmarshal(data, &repos)
	if err != nil {
		return nil, err
	}

	return repos.Sync.Repos, nil
}

// GetHelmChartRepos gets the data from Helm Hub and creates an array of Helm Chart repos out of that
func GetHelmChartRepos() ([]Repo, error) {
	data, err := getChartConfig()
	if err != nil {
		return nil, err
	}

	return marshalRepos(data)
}

// RepositoryHashmap creates a hashmap of repositories that makes searching faster and easier
func RepositoryHashmap(repos []Repo) map[string]bool {
	m := make(map[string]bool)

	for _, item := range repos {
		m[item.Name] = true
	}

	return m
}
