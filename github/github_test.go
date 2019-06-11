package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChartConfig(t *testing.T) {
	data, err := getChartConfig()
	assert.NoError(t, err)
	assert.NotNil(t, data)
}

func TestMarshalRepos(t *testing.T) {
	data := `sync:
  repos:
    - name: stable
      url: https://kubernetes-charts.storage.googleapis.com
    - name: jfrog
      url: https://charts.jfrog.io
    - name: rimusz
      url: https://charts.rimusz.net
`
	repos, err := marshalRepos([]byte(data))
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, len(repos), 3)
}
