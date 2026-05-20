# Contributing to TrueKalender

Thanks for your interest in contributing. This document covers how to get set up, the conventions used, and the process for submitting changes.

---

## Prerequisites

- Go 1.26+
- GCC (required for `go-sqlite3` CGO compilation)
- Git

---

## Getting Started

```bash
git clone https://github.com/Shivam583-hue/TrueKalender.git
cd TrueKalender
go mod download
go build .
```

Run directly:

```bash
go run .
```

---

## Project Layout

| Path | Purpose |
|---|---|
| `main.go` | Entry point |
| `db/db.go` | All SQLite interaction |
| `model/model.go` | Calendar model, Update loop, View |
| `model/form.go` | Task form model and shared `models` slice |
| `model/renderPanel.go` | Side panel renderer |
| `model/sync.go` | Kanban sync Bubble Tea command |
| `styles/styles.go` | All Lipgloss style definitions |

---

## Architecture Notes

- This project follows the [Bubble Tea](https://github.com/charmbracelet/bubbletea) Model-Update-View pattern.
- The `models` slice in `model/form.go` is shared state used to switch between the calendar and form screens. When adding a new screen, add a new constant and extend that slice.
- All database access goes through `db/db.go`. Do not open SQLite connections outside that package (the sync to external DBs in `sync.go` is an exception and should ideally be refactored into its own package).
- Styles are centralised in `styles/styles.go`. Do not define inline `lipgloss` styles in model or render files.

---

## Making Changes

1. Fork the repository and create a branch from `main`:
   ```bash
   git checkout -b your-feature-name
   ```

2. Make your changes. Keep commits focused — one logical change per commit.

3. Ensure the project builds and runs without errors:
   ```bash
   go build ./...
   ```

4. Open a pull request against `main` with a clear description of what changed and why.

---

## Pull Request Guidelines

- Keep PRs small and focused. Large, unscoped PRs are harder to review.
- Describe the motivation, not just the change.
- If your PR fixes a bug, include a short description of how to reproduce it.
- If your PR adds a feature, update `README.md` (keybindings table, features list, etc.) as part of the same PR.

---

## Reporting Issues

Open an issue on GitHub with:
- A clear title
- Steps to reproduce
- What you expected vs. what happened
- Your OS, terminal emulator, and Go version
