# Git Workflow Guidelines

## Overview

This document outlines the Git workflow and best practices for project to ensure consistent collaboration and code quality.

## Branching Strategy

### Git Flow Model

We use a simplified Git Flow model with the following branches:

- **main**: Production-ready code
- **develop**: Integration branch for features
- **feature/***: Feature development branches
- **bugfix/***: Bug fix branches
- **hotfix/***: Critical fixes for production
- **release/***: Release preparation branches

### Branch Naming Conventions

#### Feature Branches
```
feature/user-authentication
feature/product-catalog
feature/order-management
feature/payment-integration
```

#### Bug Fix Branches
```
bugfix/user-login-error
bugfix/product-search-bug
bugfix/order-validation-fix
```

#### Hotfix Branches
```
hotfix/security-patch
hotfix/critical-bug-fix
hotfix/performance-issue
```

#### Release Branches
```
release/v1.2.0
release/v1.3.0
```

## Commit Message Format

### Conventional Commits

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Commit Types

- **feat**: New features
- **fix**: Bug fixes
- **docs**: Documentation changes
- **style**: Code style changes (formatting, missing semi-colons, etc.)
- **refactor**: Code refactoring without changing functionality
- **test**: Adding or modifying tests
- **chore**: Maintenance tasks (build, dependencies, etc.)
- **perf**: Performance improvements
- **ci**: CI/CD configuration changes
- **build**: Build system changes

### Examples

```
feat(auth): add JWT token authentication

Implement JWT-based authentication system with:
- Token generation and validation
- Middleware for protected routes
- User session management

Closes #123

fix(database): resolve connection pool exhaustion

Fixed issue where database connections were not being properly
released back to the pool, causing connection exhaustion under
high load.

Fixes #456

docs(api): update OpenAPI specification

Updated API documentation to reflect new endpoints and
authentication requirements.

chore(deps): update dependencies to latest versions

Updated all dependencies to their latest stable versions:
- gin-gonic/gin: v1.9.1
- jackc/pgx: v5.4.3
- go-redis/redis: v8.11.5
```

## Workflow Process

### 1. Starting New Work

```bash
# Switch to develop branch
git checkout develop

# Pull latest changes
git pull origin develop

# Create new feature branch
git checkout -b feature/user-authentication

# Push branch to remote
git push -u origin feature/user-authentication
```

### 2. Working on Features

```bash
# Make changes and commit frequently
git add .
git commit -m "feat(auth): implement user registration endpoint"

# Push changes regularly
git push origin feature/user-authentication

# Keep branch up to date with develop
git checkout develop
git pull origin develop
git checkout feature/user-authentication
git merge develop
```

### 3. Preparing for Review

```bash
# Ensure all tests pass
go test ./...

# Run linters
golangci-lint run

# Update documentation if needed
# Commit any final changes
git add .
git commit -m "docs(auth): update authentication documentation"

# Push final changes
git push origin feature/user-authentication
```

### 4. Creating Pull Request

#### PR Title Format
```
feat(auth): implement JWT token authentication
fix(database): resolve connection pool exhaustion
docs(api): update OpenAPI specification
```

#### PR Description Template
```markdown
## Summary
Brief description of the changes made.

## Changes
- [ ] Added new feature X
- [ ] Fixed bug Y
- [ ] Updated documentation Z

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing completed

## Breaking Changes
- [ ] No breaking changes
- [ ] Breaking changes (describe below)

## Screenshots/Logs
(If applicable)

## Checklist
- [ ] Code follows project coding standards
- [ ] Self-review completed
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No merge conflicts
```

### 5. Code Review Process

#### For Authors
- Ensure CI/CD passes
- Address all review comments
- Keep PR focused and small
- Provide clear description and context

#### For Reviewers
- Review within 24 hours
- Check for:
  - Code quality and standards
  - Security vulnerabilities
  - Performance implications
  - Test coverage
  - Documentation updates

#### Review Comments
```
# Constructive feedback
Consider using a more descriptive variable name here.

# Questions
Why did you choose this approach over X?

# Suggestions
LGTM! Small suggestion: we could extract this into a helper function.

# Approval
LGTM! Great work on the error handling.
```

### 6. Merging

```bash
# Squash and merge for feature branches
# This creates a clean history on develop

# After PR is approved and merged
git checkout develop
git pull origin develop
git branch -d feature/user-authentication
git push origin --delete feature/user-authentication
```

## Release Process

### 1. Creating Release Branch

```bash
# Create release branch from develop
git checkout develop
git pull origin develop
git checkout -b release/v1.2.0
git push -u origin release/v1.2.0
```

### 2. Release Preparation

```bash
# Update version numbers
# Update CHANGELOG.md
# Final testing
# Bug fixes only on release branch

# Example changelog entry
## [1.2.0] - 2024-01-15

### Added
- JWT authentication system
- User profile management
- Product search functionality

### Fixed
- Database connection pool issues
- Memory leak in cache implementation

### Changed
- Improved error handling across all endpoints
- Updated OpenAPI documentation

### Removed
- Deprecated legacy authentication endpoints
```

### 3. Release Deployment

```bash
# Merge release to main
git checkout main
git merge release/v1.2.0
git tag -a v1.2.0 -m "Release version 1.2.0"
git push origin main
git push origin v1.2.0

# Merge back to develop
git checkout develop
git merge release/v1.2.0
git push origin develop

# Delete release branch
git branch -d release/v1.2.0
git push origin --delete release/v1.2.0
```

## Hotfix Process

### 1. Creating Hotfix

```bash
# Create hotfix branch from main
git checkout main
git pull origin main
git checkout -b hotfix/security-patch
git push -u origin hotfix/security-patch
```

### 2. Implementing Fix

```bash
# Make necessary changes
git add .
git commit -m "fix(security): patch SQL injection vulnerability"

# Test thoroughly
go test ./...
# Run security scans
# Manual testing
```

### 3. Deploying Hotfix

```bash
# Merge to main
git checkout main
git merge hotfix/security-patch
git tag -a v1.2.1 -m "Hotfix version 1.2.1"
git push origin main
git push origin v1.2.1

# Merge to develop
git checkout develop
git merge hotfix/security-patch
git push origin develop

# Delete hotfix branch
git branch -d hotfix/security-patch
git push origin --delete hotfix/security-patch
```

## Best Practices

### Commit Guidelines

1. **Atomic Commits**: Each commit should represent a single logical change
2. **Clear Messages**: Write descriptive commit messages
3. **Frequent Commits**: Commit early and often
4. **No Secrets**: Never commit sensitive information

### Branch Management

1. **Short-lived Branches**: Keep feature branches short-lived
2. **Regular Updates**: Regularly merge develop into feature branches
3. **Clean History**: Use squash merges for feature branches
4. **Delete Merged Branches**: Clean up merged branches promptly

### Code Quality

1. **Pre-commit Hooks**: Use pre-commit hooks for basic checks
2. **CI/CD Integration**: Ensure all checks pass before merging
3. **Code Reviews**: All changes must be reviewed
4. **Testing**: Include appropriate tests with all changes

## Git Hooks

### Pre-commit Hook

```bash
#!/bin/sh
# .git/hooks/pre-commit

# Run tests
echo "Running tests..."
go test ./...
if [ $? -ne 0 ]; then
    echo "Tests failed, aborting commit"
    exit 1
fi

# Run linters
echo "Running linters..."
golangci-lint run
if [ $? -ne 0 ]; then
    echo "Linting failed, aborting commit"
    exit 1
fi

# Check for secrets
echo "Checking for secrets..."
if grep -r "password\|secret\|key" --include="*.go" --include="*.yaml" --include="*.json" .; then
    echo "Potential secrets found, please review"
    exit 1
fi

echo "Pre-commit checks passed"
```

### Commit Message Hook

```bash
#!/bin/sh
# .git/hooks/commit-msg

# Check commit message format
commit_regex='^(feat|fix|docs|style|refactor|test|chore|perf|ci|build)(\(.+\))?: .{1,50}'

if ! grep -qE "$commit_regex" "$1"; then
    echo "Invalid commit message format!"
    echo "Format: type(scope): description"
    echo "Example: feat(auth): add user authentication"
    exit 1
fi
```

## Troubleshooting

### Common Issues

1. **Merge Conflicts**
   ```bash
   git checkout develop
   git pull origin develop
   git checkout feature/branch
   git merge develop
   # Resolve conflicts
   git add .
   git commit -m "resolve merge conflicts"
   ```

2. **Accidental Commits**
   ```bash
   # Undo last commit (keep changes)
   git reset --soft HEAD~1
   
   # Undo last commit (discard changes)
   git reset --hard HEAD~1
   ```

3. **Wrong Branch**
   ```bash
   # Move commits to correct branch
   git checkout correct-branch
   git cherry-pick commit-hash
   git checkout wrong-branch
   git reset --hard HEAD~1
   ```

This workflow ensures consistent collaboration, maintains code quality, and provides a clear history of changes in project.