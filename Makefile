# Makefile for AWS Research Wizard
# Provides convenient commands for development, testing, and deployment

.PHONY: help install install-dev test test-unit test-integration test-coverage lint format type-check security clean docs build deploy-docs pre-commit setup-dev go-build go-test go-install build-all

# Default target
help: ## Show this help message
	@echo "AWS Research Wizard - Development Commands"
	@echo "=========================================="
	@echo ""
	@echo "Python Implementation:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {if ($$1 !~ /^go-/) printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Go Implementation:"
	@awk 'BEGIN {FS = ":.*?## "} /^go-[a-zA-Z_-]+:.*?## / {printf "\033[32m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Combined Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^build-all:.*?## / {printf "\033[33m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Installation targets
install: ## Install production dependencies
	pip install -e .

install-dev: ## Install development dependencies
	pip install -e ".[dev]"
	pip install -r requirements-dev.txt
	pre-commit install

# Testing targets
test: ## Run all tests
	python -m pytest

test-unit: ## Run unit tests only
	python -m pytest tests/unit/ -v

test-integration: ## Run integration tests (requires AWS credentials)
	python -m pytest tests/integration/ -v --tb=short

test-coverage: ## Run tests with coverage report
	python -m pytest --cov=. --cov-report=html --cov-report=term --cov-fail-under=85

test-performance: ## Run performance benchmarks
	python -m pytest tests/performance/ --benchmark-only --benchmark-json=benchmark-results.json

# Code quality targets
lint: ## Run all linting checks
	black --check --diff .
	isort --check-only --diff .
	flake8 .
	pylint --rcfile=pyproject.toml $(shell find . -name "*.py" -not -path "./tests/*" -not -path "./venv/*" -not -path "./.venv/*")

format: ## Auto-format code
	autoflake --in-place --recursive --remove-all-unused-imports --remove-unused-variables .
	black .
	isort .

type-check: ## Run static type checking
	mypy .

security: ## Run security checks
	bandit -r . -f json -o bandit-report.json
	bandit -r . -ll
	safety check

# Documentation targets
docs: ## Build documentation
	sphinx-build -W -b html docs docs/_build/html

docs-live: ## Build documentation with live reload
	sphinx-autobuild docs docs/_build/html --host 0.0.0.0 --port 8000

docs-clean: ## Clean documentation build
	rm -rf docs/_build/

# Development targets
clean: ## Clean up build artifacts and cache
	rm -rf build/
	rm -rf dist/
	rm -rf *.egg-info/
	rm -rf .pytest_cache/
	rm -rf .mypy_cache/
	rm -rf .coverage
	rm -rf htmlcov/
	rm -rf .tox/
	find . -type d -name __pycache__ -delete
	find . -type f -name "*.pyc" -delete
	find . -type f -name "*.pyo" -delete
	rm -f bandit-report.json safety-report.json benchmark-results.json coverage.xml

setup-dev: ## Set up development environment
	python -m venv venv
	@echo "Virtual environment created. Activate with:"
	@echo "  source venv/bin/activate  # Linux/macOS"
	@echo "  venv\\Scripts\\activate     # Windows"
	@echo "Then run: make install-dev"

# Pre-commit targets
pre-commit: ## Run pre-commit hooks on all files
	pre-commit run --all-files

pre-commit-install: ## Install pre-commit hooks
	pre-commit install

pre-commit-update: ## Update pre-commit hooks
	pre-commit autoupdate

# Build and distribution targets
build: ## Build distribution packages
	python -m build

build-clean: ## Clean build and rebuild
	make clean
	make build

# AWS and deployment targets
aws-check: ## Run AWS environment checker
	python aws_environment_checker.py

aws-check-export: ## Run AWS environment checker and export results
	python aws_environment_checker.py --export aws_environment_check.json

gui: ## Start the Streamlit GUI
	streamlit run gui_research_wizard.py

gui-dev: ## Start the Streamlit GUI in development mode
	streamlit run gui_research_wizard.py --server.runOnSave true --server.port 8501

# Docker targets (if Docker support is added)
docker-build: ## Build Docker image
	docker build -t aws-research-wizard .

docker-run: ## Run Docker container
	docker run -p 8501:8501 aws-research-wizard

docker-clean: ## Clean Docker images
	docker system prune -f

# Tox targets for multi-environment testing
tox: ## Run tests in multiple Python environments
	tox

tox-parallel: ## Run tox tests in parallel
	tox -p auto

# Release targets
version-patch: ## Bump patch version
	@echo "Current version: $(shell python -c 'import toml; print(toml.load(\"pyproject.toml\")[\"project\"][\"version\"])')"
	@read -p "Enter new patch version (x.y.Z): " version; \
	sed -i.bak "s/version = \"[^\"]*\"/version = \"$$version\"/" pyproject.toml && \
	rm pyproject.toml.bak
	@echo "Version updated to: $(shell python -c 'import toml; print(toml.load(\"pyproject.toml\")[\"project\"][\"version\"])')"

version-minor: ## Bump minor version
	@echo "Current version: $(shell python -c 'import toml; print(toml.load(\"pyproject.toml\")[\"project\"][\"version\"])')"
	@read -p "Enter new minor version (x.Y.z): " version; \
	sed -i.bak "s/version = \"[^\"]*\"/version = \"$$version\"/" pyproject.toml && \
	rm pyproject.toml.bak
	@echo "Version updated to: $(shell python -c 'import toml; print(toml.load(\"pyproject.toml\")[\"project\"][\"version\"])')"

# CI/CD simulation
ci: ## Simulate CI/CD pipeline locally
	make lint
	make type-check
	make security
	make test-coverage
	@echo "✅ All CI checks passed!"

ci-quick: ## Quick CI checks (faster subset)
	make format
	make test-unit
	@echo "✅ Quick CI checks passed!"

# Development workflow targets
dev-setup: ## Complete development setup
	make setup-dev
	@echo "🔧 Development environment setup complete!"
	@echo "Next steps:"
	@echo "1. Activate virtual environment"
	@echo "2. Run: make install-dev"
	@echo "3. Configure AWS credentials"
	@echo "4. Run: make aws-check"
	@echo "5. Start developing!"

dev-test: ## Development testing workflow
	make format
	make test-unit
	make type-check
	@echo "🧪 Development testing complete!"

dev-check: ## Full development quality check
	make lint
	make type-check
	make security
	make test-coverage
	@echo "✅ Full development check complete!"

# Research pack development
new-pack: ## Create a new research pack template
	@read -p "Enter research pack name (e.g., neuroscience): " name; \
	python scripts/create_research_pack.py "$$name"

test-pack: ## Test a specific research pack
	@read -p "Enter research pack module name: " pack; \
	python -m pytest tests/unit/test_"$$pack"_pack.py -v

# Documentation helpers
docs-api: ## Generate API documentation
	sphinx-apidoc -o docs/api . --separate --force

docs-check: ## Check documentation for issues
	sphinx-build -W -b linkcheck docs docs/_build/linkcheck

# Performance and profiling
profile: ## Run performance profiling
	python -m cProfile -o profile_output.prof -m pytest tests/performance/
	@echo "Profile saved to profile_output.prof"
	@echo "View with: python -c 'import pstats; pstats.Stats(\"profile_output.prof\").sort_stats(\"cumulative\").print_stats(20)'"

benchmark: ## Run benchmarks and save results
	python -m pytest tests/performance/ --benchmark-only --benchmark-json=benchmark.json
	@echo "Benchmark results saved to benchmark.json"

# Monitoring and analysis
coverage-html: ## Generate HTML coverage report
	python -m pytest --cov=. --cov-report=html
	@echo "Coverage report generated in htmlcov/"
	@echo "Open htmlcov/index.html in your browser"

dependency-check: ## Check for dependency vulnerabilities
	safety check
	pip-audit

dependency-update: ## Update dependencies
	pip-compile --upgrade requirements-dev.in
	pip-sync requirements-dev.txt

# Environment validation
validate-env: ## Validate development environment
	@echo "🔍 Validating development environment..."
	@python --version
	@pip --version
	@git --version
	@echo "📦 Checking Python packages..."
	@python -c "import boto3, streamlit, pytest, black, mypy; print('✅ Core packages available')"
	@echo "🔑 Checking AWS credentials..."
	@python -c "import boto3; boto3.Session().get_credentials() and print('✅ AWS credentials found') or print('⚠️  AWS credentials not configured')"
	@echo "✅ Environment validation complete!"

# Quick start for new developers
quickstart: ## Quick start for new developers
	@echo "🚀 AWS Research Wizard - Quick Start"
	@echo "==================================="
	@echo ""
	@echo "Setting up your development environment..."
	make dev-setup
	@echo ""
	@echo "Run the following commands to get started:"
	@echo "1. source venv/bin/activate  # Activate virtual environment"
	@echo "2. make install-dev          # Install dependencies"
	@echo "3. make aws-check           # Validate AWS setup"
	@echo "4. make gui                 # Start the GUI"
	@echo ""
	@echo "For more commands, run: make help"

# Advanced targets for maintainers
release-check: ## Pre-release validation
	make clean
	make ci
	make docs
	make build
	@echo "🚀 Release validation complete!"

release-test: ## Test release in clean environment
	make clean
	python -m venv test_env
	test_env/bin/pip install dist/*.whl
	test_env/bin/python -c "import aws_research_wizard; print('✅ Package installs correctly')"
	rm -rf test_env

# Debug and troubleshooting
debug-imports: ## Debug import issues
	python -c "import sys; print('Python path:'); [print(f'  {p}') for p in sys.path]"
	python -c "import pkg_resources; [print(f'{pkg.key}: {pkg.version}') for pkg in pkg_resources.working_set]"

debug-aws: ## Debug AWS configuration
	python -c "import boto3; session = boto3.Session(); print(f'Region: {session.region_name}'); print(f'Profile: {session.profile_name}'); print('Credentials:', '✅' if session.get_credentials() else '❌')"

# Utility targets
line-count: ## Count lines of code
	@echo "Lines of code by file type:"
	@find . -name "*.py" -not -path "./venv/*" -not -path "./.venv/*" -not -path "./build/*" | xargs wc -l | tail -1
	@echo ""
	@echo "Test coverage:"
	@find ./tests -name "*.py" | xargs wc -l | tail -1

loc: line-count ## Alias for line-count

todo: ## Find TODO comments in code
	@echo "TODO items found:"
	@grep -r "TODO\|FIXME\|XXX\|HACK" --include="*.py" . || echo "No TODO items found"

# Default Python and pip commands with virtual environment awareness
python-version: ## Show Python version info
	python --version
	python -c "import sys; print(f'Python executable: {sys.executable}')"

pip-list: ## List installed packages
	pip list

pip-outdated: ## Show outdated packages
	pip list --outdated

# Go Implementation Targets
go-build: ## Build Go binary
	@echo "🔨 Building Go implementation..."
	cd go && make build
	@echo "✅ Go binary built: go/build/aws-research-wizard-config"

go-test: ## Run Go tests
	@echo "🧪 Running Go tests..."
	cd go && make test

go-install: ## Install Go binary to ~/bin
	@echo "📦 Installing Go binary..."
	cd go && make install
	@echo "✅ Go binary installed to ~/bin/aws-research-wizard"

go-build-all: ## Build Go binaries for all platforms
	@echo "🏗️ Building Go binaries for all platforms..."
	cd go && make build-all
	@echo "✅ Multi-platform binaries built in go/build/"

go-clean: ## Clean Go build artifacts
	@echo "🧹 Cleaning Go build artifacts..."
	cd go && make clean

go-dev: ## Quick Go development build and run
	@echo "⚡ Quick Go development cycle..."
	cd go && make dev

go-fmt: ## Format Go code
	@echo "🎨 Formatting Go code..."
	cd go && make fmt

go-lint: ## Run Go linter
	@echo "🔍 Running Go linter..."
	cd go && make lint

go-security: ## Run Go security scanner
	@echo "🔒 Running Go security scanner..."
	cd go && make security

# Combined Build Targets
build-all: ## Build both Python and Go implementations
	@echo "🏗️ Building both Python and Go implementations..."
	make build
	make go-build
	@echo "✅ Both implementations built successfully!"

install-all: ## Install both Python and Go implementations
	@echo "📦 Installing both implementations..."
	make install
	make go-install
	@echo "✅ Both implementations installed!"

test-all: ## Run tests for both implementations
	@echo "🧪 Running tests for both implementations..."
	make test
	make go-test
	@echo "✅ All tests completed!"

clean-all: ## Clean artifacts for both implementations
	@echo "🧹 Cleaning all build artifacts..."
	make clean
	make go-clean
	@echo "✅ All artifacts cleaned!"

# Development workflow for both implementations
dev-all: ## Development workflow for both implementations
	@echo "⚡ Running development workflow for both implementations..."
	make format
	make go-fmt
	make test-unit
	make go-test
	@echo "✅ Development workflow complete for both implementations!"