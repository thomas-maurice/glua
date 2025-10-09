# Contributing to glua

Thanks for your interest in contributing to glua! This document outlines the requirements and best practices for contributing.

## Before Opening a Pull Request

**CRITICAL**: Test your changes thoroughly before submitting a PR. All PRs must pass the following checks:

### 1. Run All Tests

```bash
make
```

This command runs:

- Unit tests with race detection and coverage
- k8sclient integration tests with Kind cluster
- Builds all binaries (stubgen, example)

**Your PR will be rejected if tests fail.**

### 2. Run Pre-commit Hooks

```bash
pre-commit run -a
```

This checks:

- Code formatting (gofmt with `-s`)
- go mod tidy
- golangci-lint
- go vet
- Markdown linting
- Shell script linting

**Install pre-commit first:**

```bash
pip install pre-commit
pre-commit install
```

### 3. Test GitHub Actions Locally (Recommended)

```bash
make act-test-unit
```

This runs the same tests that GitHub Actions will run, catching CI failures before you push.

**Install act:**

```bash
# See https://github.com/nektos/act for installation instructions
```

## Code Quality Standards

### Go Code

- **Format**: Use `gofmt -s -w .` (simplify mode)
- **Linting**: Fix all `golangci-lint` warnings
- **Comments**: Use standard Go comment style: `// FuncName: description`
- **Complexity**: Keep cyclomatic complexity low (< 15 per function)
- **Tests**: Write unit tests for all new functions

### Commit Messages

- Be descriptive about WHAT changed and WHY
- Don't just describe the code changes (that's what diffs are for)
- Keep first line under 72 characters

Example:

```
Good: Add k8sclient module for dynamic Kubernetes operations

Bad: update files
```

### Pull Request Description

Include:

- **Summary**: What does this PR do?
- **Motivation**: Why is this change needed?
- **Testing**: How did you test this?
- **Breaking Changes**: Any API changes?

## Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. **Run `make` to ensure all tests pass**
5. **Run `pre-commit run -a` to check code quality**
6. Commit your changes
7. Push to your fork
8. Open a Pull Request

## Testing Requirements

### Unit Tests

- Must pass with race detection: `make test-unit`
- Add tests for new functionality in the same package
- Test files go in `pkg/*/testdata/` for Lua scripts

### Integration Tests

- k8sclient example must pass: `make test-k8sclient`
- Requires Kind and kubectl installed
- Tests run in isolated Kind cluster

### Example Updates

If you add new modules or features:

- Update `example/` directory with usage examples
- Run the example to verify it works: `cd example && go run .`
- Update README.md if adding user-facing features

## What Gets Rejected

Pull requests will be rejected if:

- Tests fail (unit or integration)
- Pre-commit hooks fail
- Code is not formatted with `gofmt -s`
- No tests added for new functionality
- Breaking changes without discussion
- Unclear or missing PR description

## Getting Help

- **Questions**: Open a GitHub Discussion
- **Bugs**: Open an issue with reproduction steps
- **Features**: Open an issue to discuss before implementing

## Final Checklist

Before submitting your PR, confirm:

- [ ] `make` passes (all tests + build)
- [ ] `pre-commit run -a` passes
- [ ] New tests added for new code
- [ ] Documentation updated if needed
- [ ] Commit messages are descriptive
- [ ] PR description explains the change

**Remember**: It's YOUR responsibility to ensure your code works before asking for review. Test your shit before opening a PR.
