# GitHub Actions Workflows

This directory contains GitHub Actions workflows for the Maestro CLI project.

## Workflows

### 1. `ci.yml` - Main CI Pipeline
- **Triggers**: Push to main/develop, PRs to main/develop, manual dispatch
- **Jobs**: 
  - `lint`: Runs code quality checks
  - `build`: Builds the maestro CLI binary
  - `test`: Runs all tests
  - `release-build`: Creates release artifacts for main branch pushes
- **Matrix**: Tests on multiple Go versions and operating systems

### 2. `lint.yml` - Code Quality
- **Triggers**: Push to main/develop, PRs to main/develop, manual dispatch
- **Purpose**: Runs `./tools/lint.sh` for code quality checks
- **Artifacts**: Uploads lint results on failure

### 3. `build.yml` - Build Verification
- **Triggers**: Push to main/develop, PRs to main/develop, manual dispatch
- **Purpose**: Runs `./build.sh` and verifies binary functionality
- **Matrix**: Tests on Go 1.21, 1.22, and 1.23
- **Artifacts**: Uploads build artifacts and logs

### 4. `test.yml` - Test Suite
- **Triggers**: Push to main/develop, PRs to main/develop, manual dispatch
- **Purpose**: Runs `./test.sh` for comprehensive testing
- **Matrix**: Tests on Go 1.21, 1.22, and 1.23
- **Artifacts**: Uploads test results and coverage reports

## Workflow Features

- **Parallel Execution**: Individual workflows can run in parallel
- **Sequential Dependencies**: Main CI workflow runs jobs in sequence (lint → build → test)
- **Matrix Testing**: Tests across multiple Go versions
- **Artifact Upload**: Build artifacts and test results are preserved
- **Cross-Platform**: Release builds for Ubuntu, Windows, and macOS
- **Manual Triggers**: All workflows support manual dispatch

## Status Badges

Add these badges to your README.md:

```markdown
![CI](https://github.com/your-org/maestro-cli/workflows/CI/badge.svg)
![Lint](https://github.com/your-org/maestro-cli/workflows/Lint/badge.svg)
![Build](https://github.com/your-org/maestro-cli/workflows/Build/badge.svg)
![Test](https://github.com/your-org/maestro-cli/workflows/Test/badge.svg)
```

## Local Development

To run the same checks locally:

```bash
# Run all checks
./tools/lint.sh && ./build.sh && ./test.sh

# Run individual checks
./tools/lint.sh    # Code quality
./build.sh         # Build verification
./test.sh          # Test suite
```

## Workflow Configuration

- **Go Version**: Primary testing on Go 1.21, with matrix testing on 1.22 and 1.23
- **Operating Systems**: Ubuntu (primary), Windows and macOS (release builds)
- **Artifact Retention**: 7 days for logs, 30 days for test results, 90 days for releases
- **Cache**: Go module cache is enabled for faster builds
