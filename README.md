# artisync-hub

[![Go Report Card](https://goreportcard.com/badge/github.com/retgits/artisync-hub?style=flat-square)](https://goreportcard.com/report/github.com/retgits/artisync-hub)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/retgits/artisync-hub)
[![Release](https://img.shields.io/github/release/retgits/artisync-hub.svg?style=flat-square)](https://github.com/retgits/artisync-hub/releases/latest)

![logo](./logo.png)

A command line tool to synchronize [Helm Hub](https://github.com/helm/hub) repositories with [JFrog Artifactory](https://jfrog.com/artifactory/)

## Environment variables

To run the app, you'll need to set a few command line variables

* **ARTIFACTORY_HOST**: The hostname of JFrog Artifactory to connect to (like `http://jfrog.local/artifactory`)
* **ARTIFACTORY_REPO**: The Helm Virtual Repository to use (like `helm`)
* **ARTIFACTORY_AUTHTYPE**: The authentication type to use (either `basic` or `apikey`)
* **ARTIFACTORY_AUTH**: The authentication data to use (see below for details)

## Authentication

**artisync-hub** uses JFrog Artifactory's REST API to update the Helm repositories. The JFrog Artifactory REST API supports a few forms of authentication:

* Basic authentication using your username and password
  * Set `ARTIFACTORY_AUTHTYPE` to `basic` and `ARTIFACTORY_AUTH` to `<username>:<password>`
* Basic authentication using your username and API Key.
  * Set `ARTIFACTORY_AUTHTYPE` to `basic` and `ARTIFACTORY_AUTH` to `<username>:<apikey>`
* Using an access token instead of a password for basic authentication.
  * Set `ARTIFACTORY_AUTHTYPE` to `basic` and `ARTIFACTORY_AUTH` to `<username>:<token>`
* Using a dedicated header (X-JFrog-Art-Api) with your API Key.
  * Set `ARTIFACTORY_AUTHTYPE` to `apikey` and `ARTIFACTORY_AUTH` to `your api key>`
