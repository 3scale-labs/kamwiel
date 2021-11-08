# Kamwiel
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)

A [Kuadrant](https://github.com/Kuadrant) service that facilitates its resources to build documentation sites.

## Overview
Kamwiel makes it possible to consume [Kuadrant](https://github.com/Kuadrant) resources, extend its features connecting
3rd party services (external APIs, Kafka streams, etc...) digest and expose them making it possible for other services
and  documentation tools to employ them. Its main consumer is [Kamrad](https://github.com/3scale-labs/kamrad),
a _Developer Portal_ builder.

![Kamwiel overview](docs/images/kamwiel-overview.png?raw=true)

## Usage

Given Kamwiel is meant to work within [Kuadrant](https://github.com/Kuadrant), the setup of your cluster should at least
include its resources. At this PoC level, the minimal secured way of running it includes Kuadrant CRDs, an
[Authorino](https://github.com/kuadrant/authorino) instance to handle AuthN/AuthZ and an [Envoy](https://www.envoyproxy.io/)
proxy as the cluster Ingress object.

![Kamwiel minimal setup](docs/images/kamwiel-cluster.png?raw=true)

### Local Setup

Kamwiel comes with some useful scripts that make it easier to try it locally. You could inspect them running `make help`.
However, it needs a couple of dependencies to make it possible:

* [Go](https://golang.org/doc/install)
* [Docker](https://www.docker.com/)
* [Kind](https://kind.sigs.k8s.io/)

The easiest way to try it out is to run the script in charge of setting everything up:

```bash
make local-setup
```

This will run a local Kubernetes server loaded up with a freshly built Kamwiel image plus Authorino, Envoy and Kuadrant
CRDs with some sample data. It will also configure Authorino as the protection layer, issuing an API key that'll make it
possible to consume the cluster resources like so:

```bash
curl -H 'X-API-KEY: YOUR_AUTO_GENERATED_API_KEY' http://kamwiel-authorino.127.0.0.1.nip.io:8000/ping
```
