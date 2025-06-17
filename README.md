# Provider Spotify <!-- omit in toc -->

<div align="center">

![CI](https://github.com/crossplane-contrib/provider-spotify/actions/workflows/ci.yml/badge.svg?branch=main)
[![GitHub release](https://img.shields.io/github/release/crossplane-contrib/provider-spotify/all.svg?style=flat-square)](https://github.com/crossplane-contrib/provider-spotify/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/crossplane-contrib/provider-spotify)](https://goreportcard.com/report/github.com/crossplane-contrib/provider-spotify)

</div>

`provider-spotify` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
Spotify API. It currently only supports managing Spotify playlists.

This crossplane provider is generated from [conradludgate/terraform-provider-spotify](https://github.com/conradludgate/terraform-provider-spotify).

# Table of contents
- [Getting Started](#getting-started)
  - [Requirements](#requirements)
  - [Configuration](#configuration)
  - [Installation](#installation)
- [Usage](#usage)
- [Developing](#developing)
- [Report a Bug](#report-a-bug)

## Getting Started

### Requirements

You need to create a Spotify developer app and run Spotify's authorization
proxy server `spotify-auth-proxy`. It is recommended to install the
`spotify-auth-proxy` in the same Kubernetes cluster as the `provider-spotify`
via [this Helm
chart](https://github.com/tampakrap/helm-charts/tree/main/charts/spotify-auth-proxy).
Check its README and the comments in the values.yaml to set it up.

Additional documentation:
- https://developer.hashicorp.com/terraform/tutorials/community-providers/spotify-playlist
- https://github.com/conradludgate/terraform-provider-spotify/tree/main/spotify_auth_proxy

### Configuration

Assuming that `spotify-auth-proxy` is running, and that the Auhentication
against Spotify has been successful, you need to create a Kubernetes Secret
that contains the API Key and the URL of the `spotify-auth-proxy` Kubernetes
Service.

- If you have not set a custom API Key in the Helm chart (default):
  ```bash
  export SPOTIFY_API_KEY=$(kubectl -n spotify-auth-proxy logs spotify-auth-proxy-0 | grep APIKey | cut -d':' -f2 | xargs)
  ```

- If you have set a custom API Key in the Helm chart:
  ```bash
  export SPOTIFY_API_KEY=$(kubectl -n spotify-auth-proxy exec spotify-auth-proxy-0 -- env | grep API_KEY | cut -d'=' -f2)
  ```

Next, create the Kubernetes Secret with the API Key and the URL of the
`spotify-auth-proxy` Kubernetes Service:

```bash
sed -e "s/YOUR_API_KEY/$SPOTIFY_API_KEY/" examples/providerconfig/secret.yaml.tmpl > examples/providerconfig/secret.yaml
```

### Installation

Install the provider by using the following command after changing the image tag
to the [latest release](https://marketplace.upbound.io/providers/crossplane-contrib/provider-spotify)
using either of the following methods:

- Using [up](https://docs.upbound.io/reference/cli/):
  ```bash
  up ctp provider install crossplane-contrib/provider-spotify:v0.2.5
  ```

- Using [crossplane](https://docs.crossplane.io/latest/cli/):
  ```bash
  crossplane xpkg install provider crossplane-contrib/provider-spotify:v0.2.5
  ```

- Using declarative installation:
  ```bash
  kubectl apply -f examples/install.yaml
  ```

You can see the API reference [here](https://doc.crds.dev/github.com/crossplane-contrib/provider-spotify).

Finally, you can install the Secret and the ProviderConfig:

```bash
kubectl apply -f examples/providerconfig/
```

You should get outputs similar to the following:
```
➜ kubeclt get providers
NAME               INSTALLED   HEALTHY   PACKAGE                                      AGE
provider-spotify   True        True      crossplane-contrib/provider-spotify:v0.2.0   12m
➜ kubectl get spotify
NAME                                           AGE
providerconfig.spotify.crossplane.io/default   4m9s
➜ kubectl get secrets provider-spotify-example-creds
NAME                             TYPE     DATA   AGE
provider-spotify-example-creds   Opaque   1      7m7s
```

## Usage

See [this example playlist](examples/playlist/playlist.yaml). Example outputs:

```
➜ kubectl apply -f examples/playlist/playlist.yaml
playlist.playlist.spotify.crossplane.io/crossplane-can-play-music created
➜ kubectl get spotify
NAME                                                                READY   SYNCED   EXTERNAL-NAME            AGE
playlist.playlist.spotify.crossplane.io/crossplane-can-play-music   True    True     3HXwBJSvBPHnWHQZ3z0o3b   4m44s

NAME                                           AGE
providerconfig.spotify.crossplane.io/default   13m

NAME                                                                             AGE     CONFIG-NAME   RESOURCE-KIND   RESOURCE-NAME
providerconfigusage.spotify.crossplane.io/46502e43-db94-4ba1-85bc-6f7df2352459   4m44s   default       Playlist        crossplane-can-play-music
```

📣 🔊 Side note, this is [a real
playlist](https://open.spotify.com/playlist/3HXwBJSvBPHnWHQZ3z0o3b), which is
tracked in [this Git
repository](https://github.com/tampakrap/spotify-playlists-as-code), PRs are
always welcome!

## Developing

Run code-generation pipeline:
```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/crossplane-contrib/provider-spotify/issues).
