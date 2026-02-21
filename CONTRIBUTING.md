# Contributing Guidelines

## Conventional Commits
We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.
Example: `feat(vlan): add support for isolation mode`

## Workflow
1. Use feature branches for all changes.
2. Use `git worktree` for isolated agent development (see `agent.md`).
3. Ensure `make pre-commit` passes before committing.
4. Merge to `main` via Pull Request to trigger automated releases.

## Adding a New Resource

To maintain consistency, follow these steps when adding a new entity:

1.  **Client**: Add the generic CRUD methods to `internal/client/client.go`.
2.  **Resource**: Create `provider/resource_<name>.go`, embedding `BaseResource`. Use `internal/provider/utils` for type conversions.
3.  **Registration**: Register the resource in `provider/provider.go`.
4.  **Tests**: Create `provider/resource_<name>_test.go` with full lifecycle acceptance tests.
5.  **Docs**: Run `go generate ./...` to update documentation.
