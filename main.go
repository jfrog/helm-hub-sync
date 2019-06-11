package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/retgits/artisync-hub/artifactory"
	"github.com/retgits/artisync-hub/github"
	"github.com/retgits/artisync-hub/sync"
)

var (
	version             = "dev"
	buildTime           = ""
	modifiedVirtualRepo = false
	jfrogHost           = os.Getenv("ARTIFACTORY_HOST")
	helmVirtualRepo     = os.Getenv("ARTIFACTORY_REPO")
	authType            = os.Getenv("ARTIFACTORY_AUTHTYPE")
	userPass            = os.Getenv("ARTIFACTORY_AUTH")
	authHeaderName      = ""
	authHeaderValue     = ""
)

func main() {
	fmt.Printf("Running artisync-hub version [%s-%s]", version, buildTime)

	if len(jfrogHost) == 0 {
		panic("No Artifactory host set. Set environment variable ARTIFACTORY_HOST to continue...")
	}

	if len(helmVirtualRepo) == 0 {
		panic("No Artifactory virtual Helm repository set. Set environment variable ARTIFACTORY_REPO to continue...")
	}

	switch strings.ToLower(authType) {
	case "basic":
		authHeaderName = "authorization"
		authHeaderValue = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(userPass)))
	case "apikey":
		authHeaderName = "X-JFrog-Art-Api"
		authHeaderValue = userPass
	}

	if len(authHeaderName) == 0 {
		panic("No authentication method set. Set environment variable ARTIFACTORY_AUTHTYPE to continue...")
	}

	artiRepos, err := artifactory.GetRepositories(jfrogHost, authHeaderName, authHeaderValue)
	if err != nil {
		panic(err)
	}

	helmRepo, err := artifactory.GetRepository(jfrogHost, authHeaderName, authHeaderValue, helmVirtualRepo)
	if err != nil {
		panic(err)
	}

	githubRepos, err := github.GetHelmChartRepos()
	if err != nil {
		panic(err)
	}

	artiRepoMap := artifactory.RepositoryHashmap(artiRepos)
	githubRepoMap := github.RepositoryHashmap(githubRepos)

	for _, repo := range sync.NotInGitHub(githubRepoMap, artiRepos) {
		fmt.Printf("Helm Chart repository [%s] no longer present in GitHub data\n", repo.Key)
		err := artifactory.DeleteRepository(jfrogHost, authHeaderName, authHeaderValue, repo.Key)
		if err != nil {
			fmt.Printf("Error removing %s from Artifactory: %s", repo.Key, err.Error())
		}
		helmRepo.Repositories = sync.RemoveFromSlice(helmRepo.Repositories, repo.Key)
		modifiedVirtualRepo = true
	}

	for _, repo := range sync.NotInArtifactory(artiRepoMap, githubRepos) {
		fmt.Printf("Adding Helm Chart repository [%s] to Artifactory\n", repo.Name)
		err := artifactory.CreateRepository(jfrogHost, authHeaderName, authHeaderValue, repo.Name, repo.URL)
		if err != nil {
			fmt.Printf("Error adding %s to Artifactory: %s", repo.Name, err.Error())
		}
		helmRepo.Repositories = append(helmRepo.Repositories, repo.Name)
		modifiedVirtualRepo = true
	}

	if modifiedVirtualRepo {
		fmt.Printf("Made changes to Artifactory, updating virtual repository %s", helmVirtualRepo)
		newRepoContent := artifactory.Repository{
			Key:          helmVirtualRepo,
			PackageType:  "helm",
			Repositories: helmRepo.Repositories,
			Rclass:       helmRepo.Rclass,
		}
		err := artifactory.UpdateRepository(jfrogHost, authHeaderName, authHeaderValue, helmVirtualRepo, newRepoContent)
		if err != nil {
			fmt.Printf("Error updating %s: %s", helmVirtualRepo, err.Error())
		}
	} else {
		fmt.Printf("Artifactory and Helm Hub are in sync...\n")
	}
}
