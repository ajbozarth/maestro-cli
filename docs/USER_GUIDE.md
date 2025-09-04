# Maestro CLI User Guide

Welcome to the Maestro CLI! This guide will help you get started with managing vector databases and their resources using the Maestro command-line interface.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Basic Commands](#basic-commands)
- [Vector Database Management](#vector-database-management)
- [Collection Management](#collection-management)
- [Document Management](#document-management)
- [Validation](#validation)
- [Environment Variables](#environment-variables)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)

## Installation

### Prerequisites

- Go 1.21 or later
- Access to a vector database (Milvus or Weaviate)

### Building from Source

1. Clone the repository:
```bash
git clone <repository-url>
cd maestro-cli
```

2. Build the CLI:
```bash
./build.sh
```

3. Verify installation:
```bash
./maestro --version
```

## Quick Start

1. **Set up your MCP server connection** (optional):
```bash
export MAESTRO_MCP_SERVER_URI="http://localhost:8030/mcp"
```

2. **Validate a configuration file**:
```bash
./maestro validate config.yaml
```

3. **List available vector databases**:
```bash
./maestro vectordb list
```

4. **Create a new vector database**:
```bash
./maestro vectordb create config.yaml
```

## Configuration

### Vector Database Configuration Schema

The Maestro CLI uses YAML configuration files that follow a specific schema. The schema is automatically downloaded from the maestro-knowledge repository when needed.

Example configuration file (`config.yaml`):
```yaml
apiVersion: maestro/v1alpha1
kind: VectorDatabase
metadata:
  name: my-vector-db
  labels:
    app: my-app
spec:
  type: milvus  # or weaviate
  uri: localhost:19530
  collection_name: my_collection
  embedding: text-embedding-3-small
  mode: local  # or remote
```

### Configuration Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `apiVersion` | string | Yes | Must be `maestro/v1alpha1` |
| `kind` | string | Yes | Must be `VectorDatabase` |
| `metadata.name` | string | Yes | Unique name for the vector database |
| `metadata.labels` | object | No | Optional labels for the configuration |
| `spec.type` | string | Yes | Type of vector database (`milvus` or `weaviate`) |
| `spec.uri` | string | Yes | Connection URI (host:port for local, full URL for remote) |
| `spec.collection_name` | string | Yes | Name of the collection to use |
| `spec.embedding` | string | Yes | Embedding model to use |
| `spec.mode` | string | Yes | Deployment mode (`local` or `remote`) |

## Basic Commands

### Help and Version

```bash
# Show help
./maestro --help

# Show version
./maestro --version

# Show help for a specific command
./maestro vectordb --help
```

### Global Flags

- `--mcp-server-uri string`: MCP server URI (overrides MAESTRO_MCP_SERVER_URI environment variable)
- `--verbose`: Enable verbose output
- `--silent`: Suppress output (except errors)
- `--dry-run`: Show what would be done without executing

## Vector Database Management

### List Vector Databases

```bash
# List all vector databases
./maestro vectordb list

# List with verbose output
./maestro vectordb list --verbose

# Dry run (show what would be listed)
./maestro vectordb list --dry-run
```

### Create Vector Database

```bash
# Create from configuration file
./maestro vectordb create config.yaml

# Create with verbose output
./maestro vectordb create config.yaml --verbose

# Dry run (show what would be created)
./maestro vectordb create config.yaml --dry-run
```

### Delete Vector Database

```bash
# Delete a vector database
./maestro vectordb delete my-vector-db

# Delete with verbose output
./maestro vectordb delete my-vector-db --verbose

# Dry run (show what would be deleted)
./maestro vectordb delete my-vector-db --dry-run
```

## Collection Management

### List Collections

```bash
# List collections in a vector database
./maestro collection list my-vector-db

# List with verbose output
./maestro collection list my-vector-db --verbose
```

### Create Collection

```bash
# Create a collection
./maestro collection create my-vector-db my-collection

# Create with verbose output
./maestro collection create my-vector-db my-collection --verbose
```

### Delete Collection

```bash
# Delete a collection
./maestro collection delete my-vector-db my-collection

# Delete with verbose output
./maestro collection delete my-vector-db my-collection --verbose
```

## Document Management

### List Documents

```bash
# List documents in a collection
./maestro document list my-vector-db my-collection

# List with verbose output
./maestro document list my-vector-db my-collection --verbose
```

### Write Documents

```bash
# Write documents to a collection
./maestro document write my-vector-db my-collection data.json

# Write with verbose output
./maestro document write my-vector-db my-collection data.json --verbose
```

### Delete Documents

```bash
# Delete a document
./maestro document delete my-vector-db my-collection doc-id

# Delete with verbose output
./maestro document delete my-vector-db my-collection doc-id --verbose
```

## Validation

### Validate Configuration Files

```bash
# Validate a configuration file
./maestro validate config.yaml

# Validate with verbose output
./maestro validate config.yaml --verbose

# Validate with custom schema
./maestro validate config.yaml schema.json

# Dry run validation
./maestro validate config.yaml --dry-run
```

The validation command automatically downloads the latest schema from the maestro-knowledge repository if no local schema is found.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MAESTRO_MCP_SERVER_URI` | MCP server URI for communication | `http://localhost:8030/mcp` |
| `MAESTRO_TEST_MODE` | Enable test mode (for testing) | `false` |

## Examples

### Example 1: Setting up a Milvus Vector Database

1. Create a configuration file:
```yaml
# milvus-config.yaml
apiVersion: maestro/v1alpha1
kind: VectorDatabase
metadata:
  name: my-milvus-db
spec:
  type: milvus
  uri: localhost:19530
  collection_name: documents
  embedding: text-embedding-3-small
  mode: local
```

2. Validate the configuration:
```bash
./maestro validate milvus-config.yaml
```

3. Create the vector database:
```bash
./maestro vectordb create milvus-config.yaml
```

4. List to verify:
```bash
./maestro vectordb list
```

### Example 2: Setting up a Weaviate Vector Database

1. Create a configuration file:
```yaml
# weaviate-config.yaml
apiVersion: maestro/v1alpha1
kind: VectorDatabase
metadata:
  name: my-weaviate-db
spec:
  type: weaviate
  uri: http://localhost:8080
  collection_name: documents
  embedding: text-embedding-3-small
  mode: local
```

2. Validate and create:
```bash
./maestro validate weaviate-config.yaml
./maestro vectordb create weaviate-config.yaml
```

### Example 3: Working with Collections and Documents

```bash
# List collections
./maestro collection list my-vector-db

# Create a collection
./maestro collection create my-vector-db my-documents

# List documents
./maestro document list my-vector-db my-documents

# Write documents (assuming you have a data.json file)
./maestro document write my-vector-db my-documents data.json
```

## Troubleshooting

### Common Issues

1. **Schema download fails**:
   - The CLI automatically tries to download the schema from the maestro-knowledge repository
   - If download fails, ensure you have internet connectivity
   - You can provide a custom schema file: `./maestro validate config.yaml custom-schema.json`

2. **MCP server connection issues**:
   - Check that your MCP server is running
   - Verify the URI: `./maestro --mcp-server-uri http://your-server:port/mcp`

3. **Vector database connection issues**:
   - Ensure your vector database (Milvus/Weaviate) is running
   - Check the URI in your configuration file
   - Verify network connectivity

4. **Permission issues**:
   - Ensure the binary has execute permissions: `chmod +x maestro`
   - Check file permissions for configuration files

### Getting Help

- Use `./maestro --help` for general help
- Use `./maestro <command> --help` for command-specific help
- Check the logs with `--verbose` flag for detailed information
- Use `--dry-run` to see what would be executed without making changes

### Debug Mode

Enable verbose output to see detailed information about what the CLI is doing:

```bash
./maestro <command> --verbose
```

This will show:
- Schema download attempts
- MCP server communication
- Detailed error messages
- Step-by-step execution information

---

For more information, see the main [README.md](../README.md) or run `./maestro --help`.
