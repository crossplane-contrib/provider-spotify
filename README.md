# Provider Spotify

`provider-spotify` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
Spotify API. It currently only supports managing Spotify playlists.

This crossplane provider is generated from [conradludgate/terraform-provider-spotify](https://github.com/conradludgate/terraform-provider-spotify).

## Installation

### Requirements

You need to create a Spotify developer app and run Spotify's authorization
proxy server `spotify-auth-proxy`. In oder to set them up, take a look at the
following documents:
- https://developer.hashicorp.com/terraform/tutorials/community-providers/spotify-playlist
- https://github.com/conradludgate/terraform-provider-spotify/tree/main/spotify_auth_proxy

It is recommended to install the `spotify-auth-proxy` in the same Kubernetes
cluster as the `provider-spotify` via [this Helm
chart](https://github.com/tampakrap/helm-charts/tree/main/charts/spotify-auth-proxy).
Check its README and the comments in the values.yaml to set it up.

### Configuration

Assuming that `spotify-auth-proxy` is running, and that the Auhentication
against Spotify has been successful, you need to create a Kubernetes Secret
that contains the API Key and the URL of the `spotify-auth-proxy` Kubernetes
Service.

- If you have not set a custom API Key in the Helm chart (default):

```
export SPOTIFY_API_KEY=$(kubectl -n spotify-auth-proxy logs spotify-auth-proxy-0 | grep APIKey | cut -d':' -f2 | xargs)
```

- If you have set a custom API Key in the Helm chart:

```
export SPOTIFY_API_KEY=$(kubectl -n spotify-auth-proxy exec spotify-auth-proxy-0 -- env | grep API_KEY | cut -d'=' -f2)
```

Next, create the Kubernetes Secret with the API Key and the URL of the
`spotify-auth-proxy` Kubernetes Service:

```
sed -e "s/YOUR_API_KEY/$SPOTIFY_API_KEY/" examples/providerconfig/secret.yaml.tmpl > examples/providerconfig/secret.yaml
```

### ProviderConfig and Provider installation

Install the provider by using the following command after changing the image tag
to the [latest release](https://marketplace.upbound.io/providers/tampakrap/provider-spotify)
using either of the following methods (replace `$LATEST_VERSION` accordingly):

- Using [up](https://docs.upbound.io/reference/cli/):
```
up ctp provider install tampakrap/provider-spotify:v$LATEST_VERSION
```

- Using [crossplane](https://docs.crossplane.io/latest/cli/):
```
crossplane xpkg install provider tampakrap/provider-spotify:v$LATEST_VERSION
```

- Using declarative installation:
```
cat <<EOF | kubectl apply -f -
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-spotify
spec:
  package: tampakrap/provider-spotify:v$LATEST_VERSION
EOF
```

You can see the API reference [here](https://doc.crds.dev/github.com/tampakrap/provider-spotify).

Finally, you can install the Secret and the ProviderConfig:

```
kubectl apply -f examples/providerconfig/
```

You should get outputs similar to the following:
```
➜ kubeclt get providers
NAME               INSTALLED   HEALTHY   PACKAGE                                       AGE
provider-spotify   True        True      tampakrap/provider-spotify:v$LATEST_VERSION   12m
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
open an [issue](https://github.com/tampakrap/provider-spotify/issues).
