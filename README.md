# actionspanel
[![Actions Panel](https://img.shields.io/badge/actionspanel-enabled-brightgreen)](https://www.actionspanel.app/app/phunki/actionspanel)
[![Go Report Card](https://goreportcard.com/badge/github.com/phunki/actionspanel)](https://goreportcard.com/report/github.com/phunki/actionspanel)
[![Docs](https://godoc.org/github.com/phunki/actionspanel?status.svg)](https://pkg.go.dev/github.com/phunki/actionspanel?tab=doc)
[![codecov](https://codecov.io/gh/phunki/actionspanel/branch/master/graph/badge.svg)](https://codecov.io/gh/phunki/actionspanel)
[![Test Actions Panel](https://github.com/phunki/actionspanel/workflows/Test%20Actions%20Panel/badge.svg)](https://github.com/phunki/actionspanel/actions?query=workflow%3A%22Release+Actions+Panel%22)
[![release](https://img.shields.io/github/release/phunki/actionspanel.svg)](https://github.com/phunki/actionspanel/releases/latest)
[![GitHub release date](https://img.shields.io/github/release-date/phunki/actionspanel.svg)](https://github.com/phunki/actionspanel/releases)
![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/phunki/actionspanel)
![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/phunki/actionspanel)
![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/phunki/actionspanel?sort=semver)
![Docker Pulls](https://img.shields.io/docker/pulls/phunki/actionspanel)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=phunki/actionspanel)](https://dependabot.com)
[![license](https://img.shields.io/github/license/phunki/actionspanel.svg)](https://github.com/phunki/actionspanel/blob/master/LICENSE)

Manually trigger your GitHub Actions

Visit https://www.actionspanel.app to learn more.

## Running locally
actionspanel is built with [Go 1.14+](https://golang.org/dl/) so you'll need to install that.

You will need a Kubernetes cluster to deploy to for local development. We
recommend using the [kind](https://github.com/kubernetes-sigs/kind) project
which makes running a local Kubernetes cluster extremely easy. We currently
develop using `kind v0.7.0`.

Lastly, for doing the development and deployment to Kubernetes, we use
[skaffold
v1.8.0](https://github.com/GoogleContainerTools/skaffold/releases/tag/v1.8.0).

Once you have all of these tools installed, you can run and deploy the
application locally using the following command:
`skaffold dev --no-prune=false --cache-artifacts=false --port-forward`

This will build the Docker image locally, deploy it to your Kubernertes cluster
and make the application available at port `8080`.

### Creating your GitHub App
Actions Panel runs as a GitHub App. If you want to run Actions Panel locally,
you'll need to [create a GitHub
App](https://developer.github.com/apps/building-github-apps/creating-a-github-app/).

Make sure that you set your GitHub App's Webhook url to `/webhook` as that's
what the application registers for its path for receiving GitHub events.

To trigger `repository_dispatch` events, you need to configure your GitHub App
to have `Read & write` permissions for `Repository` `Contents`.
`repository_dispatch` needs write permissions in order to create these events.
Unfortunately, there's not a better or more granular way to gain access to
creating `repository_dispatch`.

We use [ngrok](https://ngrok.com/) to tunnel the GitHub webhook requests to run
locally.

### Secrets used for local development
When you create your GitHub App, you'll get a few secret values which you need
to pass to this application in order to get it to run. Create a `ConfigMap` in
your Kubernetes cluster with the following values:

```
kubectl apply -f - << EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: actionspanel
data:
  AP_INTEGRATION_ID: ""
  AP_OAUTH_CLIENT_ID: ""
  AP_OAUTH_CLIENT_SECRET: ""
  AP_PRIVATE_KEY: ""
  AP_WEBHOOK_SECRET: ""
EOF
```

*NOTE:* The `AP_PRIVATE_KEY` is the generated private key you create as part of the GitHub App, but base64 encoded.
