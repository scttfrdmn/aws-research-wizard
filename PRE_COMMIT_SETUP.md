# Pre-commit Hooks Setup for AWS Research Wizard

## Overview

The AWS Research Wizard now has comprehensive pre-commit hooks configured to ensure code quality, security, and consistency across the codebase.

## Installed Hooks

### Code Quality & Formatting
- **Black**: Python code formatting with 100-character line length
- **isort**: Import statement organization and sorting
- **autoflake**: Removes unused imports and variables
- **pyupgrade**: Upgrades Python syntax to Python 3.9+

### Linting & Analysis
- **flake8**: Python linting with additional plugins:
  - flake8-docstrings: Docstring style checking
  - flake8-bugbear: Bug and anti-pattern detection
  - flake8-comprehensions: List/dict comprehension optimization
  - flake8-simplify: Code simplification suggestions
  - flake8-pytest-style: Pytest best practices
- **pylint**: Advanced Python static analysis
- **mypy**: Static type checking with boto3 type stubs

### Security
- **bandit**: Security vulnerability scanning
- Configured to scan all Python files except tests

### Documentation
- **docformatter**: Docstring formatting and standardization
- Ensures consistent documentation style across the codebase

### General Checks
- **trailing-whitespace**: Removes trailing whitespace
- **end-of-file-fixer**: Ensures files end with newlines
- **check-yaml**: YAML syntax validation
- **check-json**: JSON syntax validation
- **check-merge-conflict**: Detects merge conflict markers
- **check-case-conflict**: Prevents case-sensitive filename conflicts
- **check-added-large-files**: Prevents commits of large files (>10MB)
- **debug-statements**: Detects debug print statements

### Shell & Configuration
- **shellcheck**: Shell script linting
- **prettier**: YAML formatting
- **hadolint**: Dockerfile linting

## Custom Hooks

### Test Coverage Check (Pre-push)
- Runs pytest with coverage analysis
- Requires minimum 85% test coverage
- Only runs on pre-push to avoid slowing down commits
- Timeout: 5 minutes

### Code Comment Density Check
- Ensures minimum 15% comment density in Python files
- Counts both inline comments and docstrings
- Excludes test files and __init__.py files
- Promotes well-documented code

### TODO/FIXME Detection
- Scans for TODO, FIXME, XXX, and HACK comments
- Provides warnings but doesn't fail the commit
- Helps track technical debt

## Usage

### Automatic Execution
Pre-commit hooks run automatically on:
- **git commit**: Runs formatting, linting, and basic checks
- **git push**: Runs additional expensive checks like test coverage

### Manual Execution
```bash
# Run all hooks on all files
pre-commit run --all-files

# Run specific hook
pre-commit run black --all-files
pre-commit run flake8 --all-files

# Run hooks on specific files
pre-commit run --files config_loader.py dataset_manager.py
```

### Skip Hooks (Emergency Use)
```bash
# Skip all pre-commit hooks
git commit --no-verify -m "Emergency commit"

# Skip specific hooks
SKIP=mypy,pylint git commit -m "Skip type checking"
```

## Configuration Files

### Pre-commit Configuration
- **File**: `.pre-commit-config.yaml`
- **Auto-update**: Weekly via pre-commit.ci
- **CI Integration**: Configured for GitHub Actions

### Tool Configurations
All tools are configured via `pyproject.toml`:
- Black: 100-character line length, Python 3.9+ target
- isort: Black-compatible import sorting
- flake8: Comprehensive rule set with reasonable ignores
- mypy: Strict type checking with boto3 stubs
- bandit: Security scanning with TOML configuration
- pylint: Advanced analysis with documentation requirements

## Benefits

1. **Consistent Code Style**: Black and isort ensure uniform formatting
2. **Early Bug Detection**: Flake8, pylint, and mypy catch issues before they reach production
3. **Security**: Bandit scans for common security vulnerabilities
4. **Documentation Quality**: Docstring checks and comment density requirements
5. **Test Coverage**: Ensures new code is properly tested
6. **Automation**: Reduces manual code review overhead

## Troubleshooting

### Hook Failures
If a hook fails:
1. Review the error message and fix the issue
2. Stage your changes: `git add .`
3. Retry the commit

### Performance Issues
If hooks are too slow:
- Consider using `SKIP` environment variable for specific hooks
- Use `--no-verify` for emergency commits (not recommended)

### Installation Issues
```bash
# Reinstall pre-commit
pip install --upgrade pre-commit
pre-commit install --install-hooks

# Clear cache if hooks are misbehaving
pre-commit clean
pre-commit install --install-hooks
```

## Integration with Development Workflow

1. **Development**: Write code with documentation
2. **Pre-commit**: Hooks automatically format and check code
3. **Pre-push**: Coverage and expensive checks run
4. **CI/CD**: Additional checks in continuous integration
5. **Code Review**: Focus on logic and design rather than style

This setup ensures that all code committed to the AWS Research Wizard repository maintains high quality, security, and documentation standards.
