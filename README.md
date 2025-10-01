# LightHouse

Lighthouse is a self-hosted CI/CD tool designed for developers who want automated deployment pipelines without the complexity or overhead of enterprise-grade systems. It watches your GitHub repositories for new commits, pulls the changes, builds your code, and deploys it to your local environment — all without needing external services.

Ideal for homelab enthusiasts, solo developers, or internal tools where speed and control matter.


## File Structure

| File                              | Purpose |
|-----------------------------------|---------|
| `cmd/LightHouse/main.go`          | Entry point for the Lighthouse CLI |
| `config/repos.json`               | Stores the list of watched GitHub repositories |
| `internal/watcher/watchlist.go`   | Handles read/write and CRUD operations on the repo watch list |
| `internal/watcher/watcher.go`     | Handles GitHub API queries and update detection |
| `internal/watcher/credentials.go` | Handles any credentials, like .env and COVE connections |
| `internal/watcher/models.go`      | Contains data models for WatchedRepo and all stats sections |
| `internal/cli/cli.go`             | CLI command parsing (planned or WIP) |




## To Do
- Add a self updater