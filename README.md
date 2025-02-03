# Terradendroflow

## Overview
Terradendroflow is a Terraform helper tool that generates a **prettified** Terraform plan output. It reads a Terraform plan (`stdout` or `JSON` format) and produces a structured, markdown-formatted summary of created, updated, replaced, deleted, and read resources.

## Features
- Parses Terraform plan (`stdout` or JSON output)
- Extracts resource identifiers and categorizes them
- Highlights changes (Created, Updated, Deleted, Replaced, Read)
- Outputs a structured Markdown report
- Includes CI integration for automated validation

## Installation
You can build the binary manually:

```bash
make build
```

Or install it via Go:

```bash
go install github.com/danny-molnar/terradendroflow@latest
```

## Usage
Run Terradendroflow on a Terraform plan output:

```bash
terradendroflow --input=tfplan.stdout --output=prettified_plan.md
```

### Example Output
```markdown
# Prettified Terraform Plan

## Created Resources
module.network.aws_s3_bucket.logs

## Updated Resources
module.compute.aws_instance.web

## Replaced Resources
module.storage.aws_ebs_volume.data
```

## Development
Run the tests locally before pushing:

```bash
make test
```

Lint the code:

```bash
make lint
```

## CI/CD
Terradendroflow uses GitHub Actions for:
- **Linting (`golangci-lint`)**
- **Running unit tests**
- **Building the binary**

CI will automatically validate code changes on each PR.

## Contributing
1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Open a pull request
