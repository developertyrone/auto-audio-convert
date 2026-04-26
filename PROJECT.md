# PROJECT.md

## Project Name
auto-audio-convert

## Description
A lightweight, cross-platform batch audio converter that automatically scans directories and converts audio files using FFmpeg. Outputs a single self-contained binary with zero runtime dependencies.

## Purpose
Simplify bulk audio conversion across platforms (Windows, macOS, Linux) with smart skipping of already-converted files, resource limiting, and optional forced re-conversion via --overwrite flag.

## Tech Stack
- Language: Go
- Framework: None (CLI binary)
- Key Libraries: Go stdlib + FFmpeg (auto-downloaded if missing)
- Build: build.sh / GitHub Actions

## Coding Mode
- Mode: direct
- Reason: Small focused Go project, direct file writing is fastest
- Set Date: 2026-02-25

## Git Remote
- URL: https://github.com/developertyrone/auto-audio-convert.git

## Key Commands
- Build: `./build.sh` or `go build -o auto-audio-convert .`
- Test: `./test.sh`
- Run: `./auto-audio-convert [options] <directory>`
- Install: `./install.sh`

## Last Agent Used
- Tool: GitHub Copilot CLI
- Task: docs: add --overwrite flag documentation for v1.2.0 release
- Date: 2026-02-25

## Notes
- v1.2.0 adds --overwrite flag to force re-conversion of existing files
- Worker pool with configurable concurrency
- 512MB memory cap, CPU throttling built-in
- .gitattributes enforces LF line endings
- Release notes: RELEASE_v1.2.0.md, DEPLOYMENT_v1.2.0.md
