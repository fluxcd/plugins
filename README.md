# plugins

[![license](https://img.shields.io/github/license/fluxcd/plugins.svg)](https://github.com/fluxcd/plugins/blob/main/LICENSE)

This repository contains the catalog of Flux CLI plugins as specified
in the [Flux CLI Plugin System RFC](https://github.com/fluxcd/flux2/blob/main/rfcs/0013-cli-plugin-system/README.md).

## Available plugins

### mirror

[Flux Mirror CLI](https://github.com/fluxcd/flux-mirror) — a command-line tool for mirroring Helm charts, OCI artifacts
and container images between registries using a declarative configuration. It populates internal mirror registries so
that clusters don't depend on external registries during reconciliation. Features include byte-for-byte mirroring with
multi-architecture support, conversion of HTTP/S Helm charts to OCI, optional signature verification, regex and semver
filtering, and a wide range of registry authentication options including cloud workload identity and OIDC.

### operator

[Flux Operator CLI](https://github.com/controlplaneio-fluxcd/flux-operator) — a command-line tool for managing
Flux Operator resources in Kubernetes clusters. It can build and validate Flux Operator resources, retrieve cluster
information, and reconcile, suspend or delete Flux components. It also provides commands for exporting resources,
creating secrets, tracing objects through the GitOps pipeline, and install AI agent skills.

### schema

[Flux Schema CLI](https://github.com/fluxcd/flux-schema) — a Kubernetes YAML validation tool that checks manifests
against JSON schemas and CEL rules using the same evaluation logic as the Kubernetes API server. It ships with built-in
schemas for Kubernetes, OpenShift, Gateway API and the Flux ecosystem CRDs, and supports custom catalogs extracted from
CRDs and OpenAPI specs. It integrates into CI/CD pipelines via GitHub Actions or Docker, and can discover and catalog
GitOps repositories into structured inventories for downstream automation.

## Managing plugins

Plugins are managed with the `flux plugin` sub-commands, which download binaries from this catalog
into `~/fluxcd/plugins` and register them as Flux CLI sub-commands.

### Search the catalog

List the plugins available in the catalog, optionally filtered by a query:

```shell
flux plugin search
flux plugin search schema
```

### Install a plugin

Download and install a plugin from the catalog. You can install the latest version, pin to a
specific version, or pin to a specific digest:

```shell
# Install the latest version
flux plugin install schema

# Install a specific version
flux plugin install schema@0.5.0

# Install pinned to a specific digest
flux plugin install schema@sha256:06e0a38db4fa6bc9f705a577c7e58dc020bfe2618e45488599e5ef7bb62e3a8a
```

Once installed, a plugin is invoked as a Flux sub-command, e.g. `flux schema --help`.

### List installed plugins

Show all installed plugins with their versions and paths (alias: `ls`):

```shell
flux plugin list
```

### Update plugins

Update installed plugins to their latest versions (alias: `upgrade`):

```shell
# Update a single plugin
flux plugin update schema

# Update all installed plugins
flux plugin update
```

### Uninstall a plugin

Remove a plugin binary and its receipt from the plugin directory (alias: `delete`):

```shell
flux plugin uninstall schema
```
