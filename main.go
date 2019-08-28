package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jfrog/helm-hub-sync/artifactory"
	"github.com/jfrog/helm-hub-sync/github"
	"github.com/jfrog/helm-hub-sync/sync"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	version             = "devbuild"
	buildTime           = time.Now().Format("20060102150405")
	modifiedVirtualRepo = false
	authHeaderName      = ""
	authHeaderValue     = ""
	githubDigest        = ""
)

// Config keeps the configuration of the app
type Config struct {
	LogLevel                    string `required:"true"`
	ConsoleLog                  bool
	ArtifactoryHost             string   `required:"true" split_words:"true"`
	ArtifactoryHelmRepo         string   `required:"true" split_words:"true"`
	ArtifactoryAuthType         string   `required:"true" split_words:"true"`
	ArtifactoryAuthData         string   `required:"true" split_words:"true"`
	ArtifactoryKeepDeletedRepos bool     `split_words:"true"`
	ArtifactoryKeepList         []string `split_words:"true"`
	GithubIgnoreList            []string `split_words:"true"`
	TimeInterval                int
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		panic(fmt.Errorf("Error parsing configuration: %s", err.Error()))
	}

	log.Debug().Msgf("%+v", config)

	loglevel, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.SetGlobalLevel(loglevel)

	if config.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msgf("Running helm-hub-sync version [%s-%s]", version, buildTime)

	log.Debug().Msg("Validating environment variables")

	switch strings.ToLower(config.ArtifactoryAuthType) {
	case "basic":
		authHeaderName = "authorization"
		authHeaderValue = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(config.ArtifactoryAuthData)))
	case "apikey":
		authHeaderName = "X-JFrog-Art-Api"
		authHeaderValue = config.ArtifactoryAuthData
	default:
		log.Panic().Msg("No authentication method set. Set environment variable ARTIFACTORY_AUTH_TYPE to continue...")
	}

	log.Debug().Msgf("Using %s authentication", strings.ToLower(config.ArtifactoryAuthType))

	if config.TimeInterval < 1 {
		log.Debug().Msg("No time interval set, performing only one run")
	}

	log.Debug().Msg("All configuration checked")

	err = artifactory.SendUsage(config.ArtifactoryHost, authHeaderName, authHeaderValue, version)
	if err != nil {
		log.Debug().Msgf("Error sending usage data to Artifactory: %s", err.Error())
	}

	log.Info().Msg("Started successfully and waiting for runs (use CTRL+c to stop)...")

	for {
		log.Debug().Msg("Starting run")

		githubRepos, err := github.GetHelmChartRepos()
		if err != nil {
			log.Panic().Msgf("Error getting Helm Chart data from GitHub: %s", err.Error())
		}

		newDigest := sync.GetMD5Hash(fmt.Sprintf("%+v", githubRepos))

		if githubDigest != newDigest {
			githubDigest = newDigest
			log.Debug().Msg("Get data from Artifactory")

			artiRepos, err := artifactory.GetRepositories(config.ArtifactoryHost, authHeaderName, authHeaderValue)
			if err != nil {
				log.Panic().Msgf("Error getting repositories from Artifactory: %s", err.Error())
			}

			helmRepo, err := artifactory.GetRepository(config.ArtifactoryHost, authHeaderName, authHeaderValue, config.ArtifactoryHelmRepo)
			if err != nil {
				log.Panic().Msgf("Error getting Helm virtual repository from Artifactory: %s", err.Error())
			}

			log.Debug().Msg("Get data from GitHub")

			artiRepoMap := artifactory.RepositoryHashmap(artiRepos)
			githubRepoMap := github.RepositoryHashmap(githubRepos)

			log.Debug().Msg("Checking which charts are no longer in GitHub")

			for _, repo := range sync.NotInGitHub(githubRepoMap, artiRepos) {
				log.Info().Msgf("Helm Chart repository [%s] no longer present in GitHub data", repo.Key)
				if !sync.Contains(repo.Key, config.ArtifactoryKeepList) || !config.ArtifactoryKeepDeletedRepos {
					err := artifactory.DeleteRepository(config.ArtifactoryHost, authHeaderName, authHeaderValue, repo.Key)
					if err != nil {
						log.Error().Msgf("Error removing %s from Artifactory: %s", repo.Key, err.Error())
					}
					helmRepo.Repositories = sync.RemoveFromSlice(helmRepo.Repositories, repo.Key)
					modifiedVirtualRepo = true
				} else {
					log.Debug().Msgf("Helm Chart repository [%s] is in Artifactory Keep List or Keep Deleted Repos is set to true, repo will not be removed", repo.Key)
				}
			}

			log.Debug().Msg("Checking which chart repos are not in Artifactory")

			for _, repo := range sync.NotInArtifactory(artiRepoMap, githubRepos) {
				log.Info().Msgf("Adding Helm Chart repository [%s] to Artifactory", repo.Name)
				if !sync.Contains(repo.Name, config.GithubIgnoreList) {
					err := artifactory.CreateRepository(config.ArtifactoryHost, authHeaderName, authHeaderValue, repo.Name, repo.URL)
					if err != nil {
						log.Error().Msgf("Error adding %s to Artifactory: %s", repo.Name, err.Error())
					}
					helmRepo.Repositories = append(helmRepo.Repositories, repo.Name)
					modifiedVirtualRepo = true
				} else {
					log.Debug().Msgf("Helm Chart repository [%s] is in GitHub Ignore List, repo will not be added", repo.Name)
				}
			}

			if modifiedVirtualRepo {
				log.Info().Msgf("Made changes to Artifactory, updating virtual repository %s", config.ArtifactoryHelmRepo)
				newRepoContent := artifactory.Repository{
					Key:          config.ArtifactoryHelmRepo,
					PackageType:  "helm",
					Repositories: helmRepo.Repositories,
					Rclass:       helmRepo.Rclass,
				}
				err := artifactory.UpdateRepository(config.ArtifactoryHost, authHeaderName, authHeaderValue, config.ArtifactoryHelmRepo, newRepoContent)
				if err != nil {
					log.Error().Msgf("Error updating %s: %s", config.ArtifactoryHelmRepo, err.Error())
				}
			} else {
				log.Info().Msg("Artifactory and Helm Hub are in sync...")
			}
		} else {
			log.Debug().Msg("No changes detected in GitHub content...")
		}

		log.Debug().Msg("Completed run")

		if config.TimeInterval < 1 {
			os.Exit(0)
		} else {
			log.Debug().Msgf("Sleeping for %d seconds", config.TimeInterval)
			time.Sleep(time.Duration(config.TimeInterval) * time.Second)
		}
	}
}
