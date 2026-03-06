# LLM Evaluation Web Tool

A web-based tool for evaluating Large Language Models on medical datasets.

## Features

- Real-time evaluation monitoring with SSE streaming
- Visual dashboard with metrics visualization
- Config file upload and visual form builder
- Support for multiple LLM providers
- Dataset: MMLU, CMMLU, MedQA, PubMedQA, MedMCQA

## Quick Start

```bash
# Install dependencies
go mod download
cd web && npm install

# Run dev servers
make dev

# Build production binary
make build
./bin/llm-eval
```

## Configuration

Copy `configs/models.yaml.example` to `configs/models.yaml` and configure your models.

## API Documentation

See [API.md](docs/API.md)

## License

MIT
