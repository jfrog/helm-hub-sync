package artifactory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	allRemoteRepositoriesURL = "/api/repositories?type=remote&packageType=helm"
	repoConfigURL            = "/api/repositories/%s"
	usageUrl                 = "/api/system/usage"
)

// Repository is a collection of metadata surrounding the artifact storage in JFrog Artifactory
type Repository struct {
	Key                                           string                 `json:"key,omitempty"`
	Type                                          string                 `json:"type,omitempty"`
	URL                                           string                 `json:"url,omitempty"`
	PackageType                                   string                 `json:"packageType,omitempty"`
	Description                                   string                 `json:"description,omitempty"`
	Notes                                         string                 `json:"notes,omitempty"`
	IncludesPattern                               string                 `json:"includesPattern,omitempty"`
	ExcludesPattern                               string                 `json:"excludesPattern,omitempty"`
	RepoLayoutRef                                 string                 `json:"repoLayoutRef,omitempty"`
	EnableComposerSupport                         bool                   `json:"enableComposerSupport,omitempty"`
	EnableNuGetSupport                            bool                   `json:"enableNuGetSupport,omitempty"`
	EnableGemsSupport                             bool                   `json:"enableGemsSupport,omitempty"`
	EnableNpmSupport                              bool                   `json:"enableNpmSupport,omitempty"`
	EnableBowerSupport                            bool                   `json:"enableBowerSupport,omitempty"`
	EnableCocoaPodsSupport                        bool                   `json:"enableCocoaPodsSupport,omitempty"`
	EnableConanSupport                            bool                   `json:"enableConanSupport,omitempty"`
	EnableDebianSupport                           bool                   `json:"enableDebianSupport,omitempty"`
	DebianTrivialLayout                           bool                   `json:"debianTrivialLayout,omitempty"`
	EnablePypiSupport                             bool                   `json:"enablePypiSupport,omitempty"`
	EnablePuppetSupport                           bool                   `json:"enablePuppetSupport,omitempty"`
	EnableDockerSupport                           bool                   `json:"enableDockerSupport,omitempty"`
	DockerAPIVersion                              string                 `json:"dockerApiVersion,omitempty"`
	ForceNugetAuthentication                      bool                   `json:"forceNugetAuthentication,omitempty"`
	EnableVagrantSupport                          bool                   `json:"enableVagrantSupport,omitempty"`
	EnableGitLFSSupport                           bool                   `json:"enableGitLfsSupport,omitempty"`
	EnableDistRepoSupport                         bool                   `json:"enableDistRepoSupport,omitempty"`
	Repositories                                  []string               `json:"repositories,omitempty"`
	ArtifactoryRequestsCanRetrieveRemoteArtifacts bool                   `json:"artifactoryRequestsCanRetrieveRemoteArtifacts,omitempty"`
	KeyPair                                       string                 `json:"keyPair,omitempty"`
	PomRepositoryReferencesCleanupPolicy          string                 `json:"pomRepositoryReferencesCleanupPolicy,omitempty"`
	DefaultDeploymentRepo                         string                 `json:"defaultDeploymentRepo,omitempty"`
	ExternalDependenciesEnabled                   bool                   `json:"externalDependenciesEnabled,omitempty"`
	VirtualRetrievalCachePeriodSecs               int64                  `json:"virtualRetrievalCachePeriodSecs,omitempty"`
	ForceMavenAuthentication                      bool                   `json:"forceMavenAuthentication,omitempty"`
	DebianDefaultArchitectures                    string                 `json:"debianDefaultArchitectures,omitempty"`
	EnabledChefSupport                            bool                   `json:"enabledChefSupport,omitempty"`
	Rclass                                        string                 `json:"rclass,omitempty"`
	Username                                      string                 `json:"username,omitempty"`
	Password                                      string                 `json:"password,omitempty"`
	HandleReleases                                bool                   `json:"handleReleases,omitempty"`
	HandleSnapshots                               bool                   `json:"handleSnapshots,omitempty"`
	SuppressPomConsistencyChecks                  bool                   `json:"suppressPomConsistencyChecks,omitempty"`
	RemoteRepoChecksumPolicyType                  string                 `json:"remoteRepoChecksumPolicyType,omitempty"`
	HardFail                                      bool                   `json:"hardFail,omitempty"`
	Offline                                       bool                   `json:"offline,omitempty"`
	BlackedOut                                    bool                   `json:"blackedOut,omitempty"`
	StoreArtifactsLocally                         bool                   `json:"storeArtifactsLocally,omitempty"`
	SocketTimeoutMillis                           int64                  `json:"socketTimeoutMillis,omitempty"`
	LocalAddress                                  string                 `json:"localAddress,omitempty"`
	RetrievalCachePeriodSecs                      int64                  `json:"retrievalCachePeriodSecs,omitempty"`
	AssumedOfflinePeriodSecs                      int64                  `json:"assumedOfflinePeriodSecs,omitempty"`
	MissedRetrievalCachePeriodSecs                int64                  `json:"missedRetrievalCachePeriodSecs,omitempty"`
	UnusedArtifactsCleanupPeriodHours             int64                  `json:"unusedArtifactsCleanupPeriodHours,omitempty"`
	FetchJarsEagerly                              bool                   `json:"fetchJarsEagerly,omitempty"`
	FetchSourcesEagerly                           bool                   `json:"fetchSourcesEagerly,omitempty"`
	ShareConfiguration                            bool                   `json:"shareConfiguration,omitempty"`
	SynchronizeProperties                         bool                   `json:"synchronizeProperties,omitempty"`
	MaxUniqueSnapshots                            int64                  `json:"maxUniqueSnapshots,omitempty"`
	MaxUniqueTags                                 int64                  `json:"maxUniqueTags,omitempty"`
	PropertySets                                  []interface{}          `json:"propertySets,omitempty"`
	ArchiveBrowsingEnabled                        bool                   `json:"archiveBrowsingEnabled,omitempty"`
	ListRemoteFolderItems                         bool                   `json:"listRemoteFolderItems,omitempty"`
	RejectInvalidJars                             bool                   `json:"rejectInvalidJars,omitempty"`
	AllowAnyHostAuth                              bool                   `json:"allowAnyHostAuth,omitempty"`
	EnableCookieManagement                        bool                   `json:"enableCookieManagement,omitempty"`
	EnableTokenAuthentication                     bool                   `json:"enableTokenAuthentication,omitempty"`
	PropagateQueryParams                          bool                   `json:"propagateQueryParams,omitempty"`
	BlockMismatchingMIMETypes                     bool                   `json:"blockMismatchingMimeTypes,omitempty"`
	MismatchingMIMETypesOverrideList              string                 `json:"mismatchingMimeTypesOverrideList,omitempty"`
	BypassHeadRequests                            bool                   `json:"bypassHeadRequests,omitempty"`
	ContentSynchronisation                        ContentSynchronisation `json:"contentSynchronisation,omitempty"`
	XrayIndex                                     bool                   `json:"xrayIndex,omitempty"`
}

// Marshal takes a repository and turns it into a byte slice
func (r *Repository) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// ContentSynchronisation ...
type ContentSynchronisation struct {
	Enabled    bool       `json:"enabled,omitempty"`
	Statistics Properties `json:"statistics,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Source     Source     `json:"source,omitempty"`
}

// Properties ...
type Properties struct {
	Enabled bool `json:"enabled,omitempty"`
}

// Source ...
type Source struct {
	OriginAbsenceDetection bool `json:"originAbsenceDetection,omitempty"`
}

func unmarshalRepositories(data []byte) ([]Repository, error) {
	var r []Repository
	err := json.Unmarshal(data, &r)
	return r, err
}

func unmarshalRepository(data []byte) (Repository, error) {
	var r Repository
	err := json.Unmarshal(data, &r)
	return r, err
}

// GetRepositories retrieves all remote Helm repositories configured in JFrog Artifactory
func GetRepositories(baseURL string, authHeaderName string, authHeaderValue string) ([]Repository, error) {
	url := fmt.Sprintf("%s%s", baseURL, allRemoteRepositoriesURL)
	response, err := callArtifactory(url, authHeaderName, authHeaderValue, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return unmarshalRepositories(response)
}

// RepositoryHashmap creates a hashmap of repositories that makes searching faster and easier
func RepositoryHashmap(repos []Repository) map[string]bool {
	m := make(map[string]bool)

	for _, item := range repos {
		m[item.Key] = true
	}

	return m
}

// GetRepository retrieves configuration details of a single repo in JFrog Artifactory
func GetRepository(baseURL string, authHeaderName string, authHeaderValue string, repoName string) (Repository, error) {
	urlSuffix := fmt.Sprintf(repoConfigURL, repoName)
	url := fmt.Sprintf("%s%s", baseURL, urlSuffix)
	response, err := callArtifactory(url, authHeaderName, authHeaderValue, http.MethodGet, nil)
	if err != nil {
		return Repository{}, err
	}
	return unmarshalRepository(response)
}

// CreateRepository creates a new remote repository with default settings in JFrog Artifactory
func CreateRepository(baseURL string, authHeaderName string, authHeaderValue string, repoName string, repoURL string) error {
	urlSuffix := fmt.Sprintf(repoConfigURL, repoName)
	url := fmt.Sprintf("%s%s", baseURL, urlSuffix)

	repo := Repository{
		Key:                               repoName,
		PackageType:                       "helm",
		Description:                       fmt.Sprintf("Remote repository for %s", repoName),
		Notes:                             "",
		IncludesPattern:                   "**/*",
		ExcludesPattern:                   "",
		RepoLayoutRef:                     "simple-default",
		EnableComposerSupport:             false,
		EnableNuGetSupport:                false,
		EnableGemsSupport:                 false,
		EnableNpmSupport:                  false,
		EnableBowerSupport:                false,
		EnableCocoaPodsSupport:            false,
		EnableConanSupport:                false,
		EnableDebianSupport:               false,
		DebianTrivialLayout:               false,
		EnablePypiSupport:                 false,
		EnablePuppetSupport:               false,
		EnableDockerSupport:               false,
		DockerAPIVersion:                  "V2",
		ForceNugetAuthentication:          false,
		EnableVagrantSupport:              false,
		EnableGitLFSSupport:               false,
		EnableDistRepoSupport:             false,
		URL:                               repoURL,
		Username:                          "",
		Password:                          "",
		HandleReleases:                    true,
		HandleSnapshots:                   true,
		SuppressPomConsistencyChecks:      true,
		RemoteRepoChecksumPolicyType:      "generate-if-absent",
		HardFail:                          false,
		Offline:                           false,
		BlackedOut:                        false,
		StoreArtifactsLocally:             true,
		SocketTimeoutMillis:               15000,
		LocalAddress:                      "",
		RetrievalCachePeriodSecs:          600,
		AssumedOfflinePeriodSecs:          300,
		MissedRetrievalCachePeriodSecs:    1800,
		UnusedArtifactsCleanupPeriodHours: 0,
		FetchJarsEagerly:                  false,
		FetchSourcesEagerly:               false,
		ShareConfiguration:                false,
		SynchronizeProperties:             false,
		MaxUniqueSnapshots:                0,
		MaxUniqueTags:                     0,
		PropertySets:                      nil,
		ArchiveBrowsingEnabled:            false,
		ListRemoteFolderItems:             true,
		RejectInvalidJars:                 false,
		AllowAnyHostAuth:                  false,
		EnableCookieManagement:            false,
		EnableTokenAuthentication:         false,
		PropagateQueryParams:              false,
		BlockMismatchingMIMETypes:         true,
		BypassHeadRequests:                false,
		ContentSynchronisation: ContentSynchronisation{
			Enabled: false,
			Statistics: Properties{
				Enabled: false,
			},
			Properties: Properties{
				Enabled: false,
			},
			Source: Source{
				OriginAbsenceDetection: false,
			},
		},
		ExternalDependenciesEnabled: false,
		XrayIndex:                   false,
		EnabledChefSupport:          false,
		Rclass:                      "remote",
	}

	payload, err := repo.Marshal()
	if err != nil {
		return err
	}

	_, err = callArtifactory(url, authHeaderName, authHeaderValue, http.MethodPut, payload)
	if err != nil {
		return err
	}

	return nil
}

// UpdateRepository updates the virtual repository in JFrog Artifactory with all newly created remote Helm repos
func UpdateRepository(baseURL string, authHeaderName string, authHeaderValue string, repoName string, repo Repository) error {
	urlSuffix := fmt.Sprintf(repoConfigURL, repoName)
	url := fmt.Sprintf("%s%s", baseURL, urlSuffix)

	payload, err := repo.Marshal()
	if err != nil {
		return err
	}

	_, err = callArtifactory(url, authHeaderName, authHeaderValue, http.MethodPost, payload)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRepository deletes a remote repository in JFrog Artifactory
func DeleteRepository(baseURL string, authHeaderName string, authHeaderValue string, repoName string) error {
	urlSuffix := fmt.Sprintf(repoConfigURL, repoName)
	url := fmt.Sprintf("%s%s", baseURL, urlSuffix)

	_, err := callArtifactory(url, authHeaderName, authHeaderValue, http.MethodDelete, nil)
	if err != nil {
		return err
	}

	return nil
}

func SendUsage(baseUrl string, authHeaderName string, authHeaderValue string, version string) error {
	url := fmt.Sprintf("%s%s", baseUrl, usageUrl)
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(map[string]interface{}{
		"productId": fmt.Sprintf("JFrogHelmHubSync/%s", version),
		"features": []string{},
	})
	if err != nil {
		return err
	}
	payload := buf.Bytes()
	_, err = callArtifactory(url, authHeaderName, authHeaderValue, http.MethodPost, payload)
	return err
}

func callArtifactory(url string, authHeaderName string, authHeaderValue string, httpMethod string, payload []byte) ([]byte, error) {
	var req *http.Request
	var err error

	if len(payload) > 0 {
		req, err = http.NewRequest(httpMethod, url, bytes.NewReader(payload))
		req.Header.Add("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(httpMethod, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Add(authHeaderName, authHeaderValue)

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
		return nil, fmt.Errorf("Artifactory returned non-OK statuscode [%d]: %s", res.StatusCode, string(body))
	}

	return body, nil
}
