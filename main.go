package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/retgits/artisync-hub/artifactory"
	"github.com/retgits/artisync-hub/github"
	"github.com/retgits/artisync-hub/sync"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	version             = "devbuild"
	buildTime           = time.Now().Format("20060102150405")
	modifiedVirtualRepo = false
	jfrogHost           = os.Getenv("ARTIFACTORY_HOST")
	helmVirtualRepo     = os.Getenv("ARTIFACTORY_REPO")
	authType            = os.Getenv("ARTIFACTORY_AUTHTYPE")
	userPass            = os.Getenv("ARTIFACTORY_AUTH")
	authHeaderName      = ""
	authHeaderValue     = ""
)

func main() {
	loglvl := os.Getenv("LOGLEVEL")
	if len(loglvl) == 0 {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		loglevel, err := zerolog.ParseLevel(loglvl)
		if err != nil {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		zerolog.SetGlobalLevel(loglevel)
	}

	if strings.EqualFold(os.Getenv("CONSOLELOG"), "true") {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msgf("Running artisync-hub version [%s-%s]", version, buildTime)

	if len(jfrogHost) == 0 {
		log.Panic().Msg("No Artifactory host set. Set environment variable ARTIFACTORY_HOST to continue...")
	}

	if len(helmVirtualRepo) == 0 {
		log.Panic().Msg("No Artifactory virtual Helm repository set. Set environment variable ARTIFACTORY_REPO to continue...")
	}

	switch strings.ToLower(authType) {
	case "basic":
		authHeaderName = "authorization"
		authHeaderValue = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(userPass)))
	case "apikey":
		authHeaderName = "X-JFrog-Art-Api"
		authHeaderValue = userPass
	default:
		log.Panic().Msg("No authentication method set. Set environment variable ARTIFACTORY_AUTHTYPE to continue...")
	}

	log.Debug().Msgf("Using %s authentication", strings.ToLower(authType))

	artiRepos, err := artifactory.GetRepositories(jfrogHost, authHeaderName, authHeaderValue)
	if err != nil {
		log.Panic().Msgf("Error getting repositories from Artifactory: %s", err.Error())
	}

	helmRepo, err := artifactory.GetRepository(jfrogHost, authHeaderName, authHeaderValue, helmVirtualRepo)
	if err != nil {
		log.Panic().Msgf("Error getting Helm virtual repository from Artifactory: %s", err.Error())
	}

	githubRepos, err := github.GetHelmChartRepos()
	if err != nil {
		log.Panic().Msgf("Error getting Helm Chart data from GitHub: %s", err.Error())
	}

	artiRepoMap := artifactory.RepositoryHashmap(artiRepos)
	githubRepoMap := github.RepositoryHashmap(githubRepos)

	for _, repo := range sync.NotInGitHub(githubRepoMap, artiRepos) {
		log.Info().Msgf("Helm Chart repository [%s] no longer present in GitHub data", repo.Key)
		err := artifactory.DeleteRepository(jfrogHost, authHeaderName, authHeaderValue, repo.Key)
		if err != nil {
			log.Error().Msgf("Error removing %s from Artifactory: %s", repo.Key, err.Error())
		}
		helmRepo.Repositories = sync.RemoveFromSlice(helmRepo.Repositories, repo.Key)
		modifiedVirtualRepo = true
	}

	for _, repo := range sync.NotInArtifactory(artiRepoMap, githubRepos) {
		log.Info().Msgf("Adding Helm Chart repository [%s] to Artifactory", repo.Name)
		err := artifactory.CreateRepository(jfrogHost, authHeaderName, authHeaderValue, repo.Name, repo.URL)
		if err != nil {
			log.Error().Msgf("Error adding %s to Artifactory: %s", repo.Name, err.Error())
		}
		helmRepo.Repositories = append(helmRepo.Repositories, repo.Name)
		modifiedVirtualRepo = true
	}

	if modifiedVirtualRepo {
		log.Info().Msgf("Made changes to Artifactory, updating virtual repository %s", helmVirtualRepo)
		newRepoContent := artifactory.Repository{
			Key:          helmVirtualRepo,
			PackageType:  "helm",
			Repositories: helmRepo.Repositories,
			Rclass:       helmRepo.Rclass,
		}
		err := artifactory.UpdateRepository(jfrogHost, authHeaderName, authHeaderValue, helmVirtualRepo, newRepoContent)
		if err != nil {
			log.Error().Msgf("Error updating %s: %s", helmVirtualRepo, err.Error())
		}
	} else {
		log.Info().Msg("Artifactory and Helm Hub are in sync...")
	}
}
