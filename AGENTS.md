# AGENTS.md

Instructions for agents onboarding a new Flux CLI plugin into this repository.

## Repository layout

- `catalog.yaml` — top-level catalog listing every plugin (name, description,
  homepage, source, license). Schema: `schemas/plugincatalog-cli-v1beta1.json`.
- `plugins/<name>.yaml` — per-plugin manifest with release artifacts for every
  supported platform. Schema: `schemas/plugin-cli-v1beta1.json`.
- `schemas/` — JSON Schemas consumed by `make validate`.
- `cmd/validate/` — Go validator binary used by CI.
- `.github/workflows/update-plugin.yaml` — reusable workflow that fetches the
  latest upstream release, updates the manifest, verifies checksums and opens a PR.
- `.github/workflows/update-<name>.yaml` — per-plugin dispatcher that calls the
  reusable workflow on a schedule.

## Prerequisites before onboarding

Confirm with the maintainer or verify yourself:

1. The upstream project publishes GitHub releases tagged `v<semver>`.
2. Each release attaches a `<bin>_<version>_checksums.txt` file plus
   `tar.gz` archives for darwin/linux and `zip` archives for windows,
   covering `amd64` and `arm64`. The archive naming must match
   `<bin>_<version>_<os>_<arch>.<ext>`.
3. The plugin `name` is lowercase, matches `^[a-z][a-z0-9-]*$`, and does NOT
   collide with a built-in Flux command (see the reserved list in
   `.github/workflows/validate-catalog.yaml`).
4. The `bin` name (the executable inside the archive) is known. If the binary
   is not at the archive root, capture the relative path — it will go in
   `extractPath` on each platform entry.
5. The license is a valid SPDX identifier (e.g. `Apache-2.0`, `AGPL-3.0`).

If any prerequisite is missing, stop and ask — do not improvise archive
layouts or checksum formats.

## Onboarding steps

1. **Add the catalog entry.** Append a new item to `catalog.yaml` with
   `name`, `description`, `homepage`, `source`, `license`. The `name` must
   equal the filename you will create in `plugins/`.

2. **Create the plugin manifest.** Add `plugins/<name>.yaml` with
   `apiVersion: cli.fluxcd.io/v1beta1`, `kind: Plugin`, and the metadata
   fields. Seed `versions` with at least one release — newest first. Use
   `plugins/operator.yaml` as the reference shape. Version strings must be
   bare semver (no `v` prefix).

3. **Populate platform artifacts.** For the initial version, download the
   upstream `checksums.txt` and fill in one platform entry per
   os/arch combination the upstream actually ships. Set `checksum` to
   `sha256:<64 hex>`. Add `extractPath` only when the binary is not at the
   archive root.

4. **Add the update workflow.** Copy `.github/workflows/update-operator.yaml`
   to `.github/workflows/update-<name>.yaml` and change:
   - workflow `name`
   - `plugin:` input (matches the filename)
   - `repo:` input (`<owner>/<repo>` of the upstream)
   - `bin:` input (binary name used in archive filenames)

   Keep the `cron` schedule (`0 */6 * * *`) unless the maintainer requests
   otherwise. The reusable workflow handles checksum download, verification,
   manifest update and PR creation — do not duplicate that logic.

5. **Validate locally.** Run `make validate`. This builds the validator and
   checks every manifest against its JSON Schema. Fix any reported error
   before opening a PR.

6. **Commit and open a PR.** Include the catalog entry, the plugin manifest,
   and the update workflow in a single PR titled
   `Add <name> plugin`. CI runs `validate-catalog` which enforces:
   - schema conformance
   - every `plugins/*.yaml` has a catalog entry and vice versa
   - no name collision with built-in commands

## Conventions

- Keep `versions` ordered newest → oldest. The update workflow prepends new
  versions; manual edits should preserve that order.
- Do not edit historical version entries. If a checksum was wrong, fix it in
  a dedicated PR and explain why.
- Platform list order: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64,
  windows/amd64, windows/arm64. The update workflow emits this order.
- Never add a platform the upstream does not ship — `make validate` passes
  but installs will 404.

## Things not to do

- Do not add plugins whose upstream does not publish a `checksums.txt`
  alongside its release assets. The update workflow depends on it.
- Do not hand-edit `bin/` artifacts — they are gitignored build output.
- Do not bypass `make validate` with `--no-verify` or similar. If validation
  fails, fix the manifest.
- Do not change `apiVersion` or `kind` values; they are pinned by schema.

## Contribution workflow for AI agents

These rules come from [`fluxcd/flux2/CONTRIBUTING.md`](https://github.com/fluxcd/flux2/blob/main/CONTRIBUTING.md) and apply to every Flux repository.

- **Do not add `Signed-off-by` or `Co-authored-by` trailers with your agent name.** Only a human can legally certify the DCO.
- **Disclose AI assistance** with an `Assisted-by` trailer naming your agent and model:
  ```sh
  git commit -s -m "Add X plugin" --trailer "Assisted-by: <agent-name>/<model-id>"
  ```
  The `-s` flag adds the human's `Signed-off-by` from their git config — do not remove it.
- **Commit message format:** Subject in imperative mood ("Add X plugin" instead of "Adding X plugin"), capitalized, no trailing period, ≤50 characters.
- **Commit body:** Add a succinct explanation of what the plugin does and link to the upstream project.
