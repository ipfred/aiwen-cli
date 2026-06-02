<!-- TRELLIS:START -->
# Trellis Instructions

These instructions are for AI assistants working in this project.

This project is managed by Trellis. The working knowledge you need lives under `.trellis/`:

- `.trellis/workflow.md` — development phases, when to create tasks, skill routing
- `.trellis/spec/` — package- and layer-scoped coding guidelines (read before writing code in a given layer)
- `.trellis/workspace/` — per-developer journals and session traces
- `.trellis/tasks/` — active and archived tasks (PRDs, research, jsonl context)

If a Trellis command is available on your platform (e.g. `/trellis:finish-work`, `/trellis:continue`), prefer it over manual steps. Not every platform exposes every command.

If you're using Codex or another agent-capable tool, additional project-scoped helpers may live in:
- `.agents/skills/` — reusable Trellis skills
- `.codex/agents/` — optional custom subagents

Managed by Trellis. Edits outside this block are preserved; edits inside may be overwritten by a future `trellis update`.

<!-- TRELLIS:END -->

# Project Architecture Reference

This project must model its CLI architecture, build pipeline, npm installer, and release distribution flow after `larksuite/cli`.

Reference sources:

- Upstream repository: https://github.com/larksuite/cli
- Local checkout for implementation reference: `E:\my_work\github_pro\cli`

Before changing CLI architecture, npm installation, binary packaging, release workflow, or skills installation, inspect the matching files in the local `larksuite/cli` checkout and keep this project structurally aligned unless there is a documented product-specific reason to diverge.

Required reference points:

- `package.json` should follow the `larksuite/cli` npm wrapper pattern: scoped package, `bin` entry pointing to `scripts/run.js`, `postinstall` invoking `scripts/install.js`, supported `os`/`cpu`, and packaged installer files.
- `scripts/run.js` should behave like the reference wrapper: intercept the `install` subcommand for the setup wizard, recover/auto-install the native binary when needed, and delegate normal commands to the Go binary.
- `scripts/install.js` should follow the reference binary installer: resolve platform/arch archive names, download from GitHub release first, fall back to npm mirror URLs, verify `checksums.txt`, extract the archive, and install the native binary into `bin/`.
- `scripts/install-wizard.js` should follow the reference one-command setup flow: `npx aiwen-geoip-cli@latest install` installs or upgrades the global CLI package and installs the project skills. Interactive configuration may be product-specific, but the install + skills flow must stay aligned.
- Release distribution should follow the `larksuite/cli` GoReleaser-style contract: build `darwin/linux/windows` for `amd64/arm64`, publish archives named `aw-cli-<version>-<os>-<arch>.(tar.gz|zip)`, publish `checksums.txt`, then publish the npm package after the release assets are available.
- Archive contents must match the installer expectation: extracted archives must contain the executable as `aw-cli` on Unix-like systems and `aw-cli.exe` on Windows, unless `scripts/install.js` explicitly documents and tests a compatibility fallback.
- Skills distribution should follow the reference `npx skills add ... -y -g` pattern. Any change to skill names, repository source, or fallback behavior must be verified against the reference flow and documented in the installer.

Do not replace this architecture with a separate install/build/distribution approach without first making the deviation explicit in this file and in the changed implementation.
