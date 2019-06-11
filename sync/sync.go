package sync

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/retgits/artisync-hub/artifactory"
	"github.com/retgits/artisync-hub/github"
)

// NotInGitHub checks which JFrog Artifactory repositories do not exist in GitHub (anymore), which means they
// can be deleted from JFrog Artifactory
func NotInGitHub(ghRepos map[string]bool, artiRepos []artifactory.Repository) (result []artifactory.Repository) {
	for _, item := range artiRepos {
		if _, found := ghRepos[item.Key]; !found {
			result = append(result, item)
		}
	}
	return
}

// NotInArtifactory checks which Helm Hub repositories do not exist in JFrog Artifactory, which means they
// can be added to JFrog Artifactory
func NotInArtifactory(artiRepos map[string]bool, ghRepos []github.Repo) (result []github.Repo) {
	for _, item := range ghRepos {
		if _, found := artiRepos[item.Name]; !found {
			result = append(result, item)
		}
	}
	return
}

// RemoveFromSlice removes a string from a slice
func RemoveFromSlice(slice []string, item string) []string {
	for idx := range slice {
		if slice[idx] == item {
			return append(slice[:idx], slice[idx+1:]...)
		}
	}

	return slice
}

// GetMD5Hash generates an MD5 hash of a string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
