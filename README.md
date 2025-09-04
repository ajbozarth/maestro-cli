# Maestro CLI

A command-line interface for managing vector databases and their resources with support for YAML configuration and environment variable substitution.

## Features

- **List vector databases**: List all available vector database instances
- **List embeddings**: List supported embeddings for a specific vector database
- **List collections**: List all collections in a specific vector database
- **List documents**: List documents in a specific collection of a vector database
- **Query documents**: Query documents using natural language with semantic search
- **Pluggable document chunking**: Configure per-collection chunking (None, Fixed with size/overlap, Sentence, Semantic)
   - Discover supported strategies with `maestro chunking list`
- **Create vector databases**: Create vector databases from YAML configuration files
- **Delete vector databases**: Delete vector databases by name
- **Validate configurations**: Validate YAML configuration files
- **Environment variable substitution**: Replace `{{ENV_VAR_NAME}}` placeholders in YAML files
- **Environment variable support**: Configure MCP server URI via environment variables
- **Command-line flag override**: Override MCP server URI via command-line flags
- **Dry-run mode**: Test commands without making actual changes
- **Verbose output**: Get detailed information about operations
- **Silent mode**: Suppress success messages
- **Safety features**: Confirmation prompts for destructive operations with `--force` flag bypass

### Enhanced User Experience

- **Command suggestions**: Intelligent suggestions for mistyped commands with "Did you mean..." functionality
- **Command aliases**: Short aliases for common commands (e.g., `vdb` for `vectordb`, `coll` for `collection`, `doc` for `document`)
- **Contextual help**: Helpful tips and next steps shown after successful operations
- **Command examples**: Comprehensive examples for all commands and subcommands
- **Error guidance**: Helpful suggestions for common error scenarios when using `--verbose` mode
- **Interactive selection**: When resource names aren't provided, the CLI prompts for selection
- **Auto-completion**: Shell completion for commands, subcommands, flags, and resource names
- **Progress indicators**: Visual feedback for long-running operations
- **Status commands**: Quick overview of system state with `maestro status`

## Installation

### Prerequisites

- Go 1.21 or later
- Access to a running Maestro MCP server

### Building

```bash
# Use the build script (recommended)
./build.sh
```

## Documentation

- **[User Guide](docs/USER_GUIDE.md)** - Comprehensive getting started guide with examples
- **[TODOs.md](TODOs.md)** - Documentation of completed migration tasks
- **[USAGE.md](USAGE.md)** - Detailed command reference and examples

## Usage

### Basic Commands

```bash
# Show help
./maestro --help

# Show version
./maestro --version

# List vector databases
./maestro vectordb list

# List vector databases with verbose output
./maestro vectordb list --verbose

# List collections for a specific vector database
./maestro collection list --vdb=my-database

# List collections with verbose output
./maestro collection list --vdb=my-database --verbose

# List embeddings for a specific vector database
./maestro embedding list --vdb=my-database

# List embeddings with verbose output
./maestro embedding list --vdb=my-database --verbose

# List documents in a specific collection
./maestro document list --vdb=my-database --collection=my-collection

# List documents with verbose output
./maestro document list --vdb=my-database --collection=my-collection --verbose

# Query documents using natural language
# Tip: pass --collection to avoid an interactive prompt
./maestro query "What is the main topic of the documents?" --vdb=my-database
./maestro query "Find information about API endpoints" --vdb=my-database --collection=my-collection --doc-limit 10

# Create vector database from YAML file
./maestro vectordb create config.yaml

# Create vector database with environment variable substitution
./maestro vectordb create config.yaml --verbose

# Delete vector database (with confirmation prompt)
./maestro vectordb delete my-database

# Delete vector database without confirmation
./maestro vectordb delete my-database --force

# Validate YAML configuration
./maestro validate config.yaml

# Create collection with chunking
./maestro create collection my-database my-collection \
   --embedding=text-embedding-3-small \
   --chunking-strategy=Sentence \
   --chunk-size=512 \
   --chunk-overlap=32

# List supported chunking strategies
./maestro chunking list
```

### Enhanced UX Features

#### Interactive Selection Examples

```bash
# The CLI will prompt you to select a vector database if not specified
./maestro collection list

# The CLI will prompt you to select a collection if not specified
./maestro document list --vdb=my-database

# The CLI will prompt you to select a document if not specified
./maestro document delete --vdb=my-database --collection=my-collection

# Search and query also prompt for a collection when --collection is omitted
./maestro search "example" --vdb=my-database
./maestro query "example" --vdb=my-database

# Note: Prompts are skipped in non-interactive contexts and with --dry-run; pass --collection to avoid prompting
```

#### Auto-completion Setup

**Bash:**
```bash
# Generate completion script
./maestro completion bash > ~/.local/share/bash-completion/completions/maestro

# Or add to your .bashrc
echo 'source <(./maestro completion bash)' >> ~/.bashrc
```

**Zsh:**
```bash
# Generate completion script
./maestro completion zsh > ~/.zsh/completions/_maestro

# Or add to your .zshrc
echo 'source <(./maestro completion zsh)' >> ~/.zshrc
```

**Fish:**
```bash
# Generate completion script
./maestro completion fish > ~/.config/fish/completions/maestro.fish
```

**PowerShell:**
```powershell
# Generate completion script
./maestro completion powershell | Out-String | Invoke-Expression
```

#### Progress Indicators

Progress indicators are automatically shown for long-running operations:

```bash
# Document creation with progress indicator
./maestro document create --name=my-doc --file=document.txt --vdb=my-database --collection=my-collection

# Query execution with progress indicator
./maestro query "What is the main topic?" --vdb=my-database

# Status check with progress indicator
./maestro status
```

#### Status Command

```bash
# Show overall system status
./maestro status

# Show status for a specific vector database
./maestro status --vdb=my-database

# Show detailed status with verbose output
./maestro status --verbose
```

Example output:
```
🔍 Maestro Knowledge System Status
==================================
📊 Vector Database: test_remote_weaviate (weaviate)
   📁 Collection: test_collection
   📄 Documents: 3
   📂 Collections: test_collection, another_collection
   🧠 Supported Embeddings: text-embedding-3-small, text-embedding-3-large
   ✅ Status: Online

📈 Summary:
   • Total Vector Databases: 1
   • Total Documents: 3
   • MCP Server: http://localhost:8030/mcp
   • Connection: ✅ Active
```

### Command Aliases

For convenience, the CLI provides shorter aliases for all resource commands:

```bash
# Vector databases
./maestro vectordb list    # or: ./maestro vdb list

# Collections  
./maestro collection list  # or: ./maestro coll list

# Documents
./maestro document list    # or: ./maestro doc list

# Embeddings
./maestro embedding list   # or: ./maestro embed list
```

### Interactive Selection Examples

```bash
# Interactive selection when flags are missing
./maestro collection list             # Prompts you to select a vector database
./maestro document list               # Prompts you to select both VDB and collection
./maestro query "test"                # Prompts you to select a vector database

# Auto-completion for resource names and file paths
./maestro collection list --vdb=<TAB> # Completes vector database names
./maestro document list --collection=<TAB> # Completes collection names
./maestro document create --file=<TAB>     # Completes file paths
./maestro collection create --embedding=<TAB> # Completes embedding models

# Command suggestions for typos
./maestro vectord                     # Suggests: vectordb
./maestro docum                       # Suggests: document
./maestro embedd                      # Suggests: embedding

# Contextual help appears after operations
./maestro vectordb list               # Shows tip about creating new databases
./maestro collection create --vdb=my-db --name=my-coll  # Shows tip about adding documents
./maestro query "test" --vdb=my-db    # Shows tips about doc-limit and collection flags

# Error guidance with --verbose
./maestro collection list --verbose   # Shows helpful suggestions for common errors
```

### MCP Server Configuration

The CLI can connect to an MCP server using several methods:

#### 1. Environment Variable (Recommended)

Set the `MAESTRO_MCP_SERVER_URI` environment variable:

```bash
export MAESTRO_MCP_SERVER_URI="http://localhost:8030"
./maestro vectordb list
```

#### 2. .env File

Create a `.env` file in the current directory with your configuration:

```bash
# MCP Server configuration
MAESTRO_MCP_SERVER_URI=http://localhost:8030

# Weaviate configuration (for Weaviate backend)
WEAVIATE_API_KEY=your-weaviate-api-key
WEAVIATE_URL=https://your-weaviate-cluster.weaviate.network

# OpenAI configuration (for OpenAI embeddings)
OPENAI_API_KEY=your-openai-api-key
```

The CLI will automatically load the `.env` file if it exists in the current directory.

#### 3. Command-line Flag

Override the MCP server URI via command-line flag:

```bash
./maestro vectordb list --mcp-server-uri="http://localhost:8030"
```

**Priority order**: Command-line flag > Environment variable > .env file > Default (http://localhost:8030)

**Supported Environment Variables**:
- `MAESTRO_MCP_SERVER_URI`: MCP server URI
- `WEAVIATE_API_KEY`: Weaviate API key for Weaviate backend
- `WEAVIATE_URL`: Weaviate cluster URL
- `OPENAI_API_KEY`: OpenAI API key for embeddings

### Chunking support

- Configure chunking when defining collections via YAML. The CLI exposes discovery of supported strategies with:

```bash
./maestro chunking list
```

#### Chunking Strategies

**None**: No chunking is performed (default)
**Fixed**: Split documents into fixed-size chunks with optional overlap
**Sentence**: Split documents at sentence boundaries with size limits
**Semantic**: AI-powered chunking that identifies semantic boundaries using sentence embeddings

#### Semantic Chunking Example

Semantic chunking uses sentence transformers to identify natural break points in documents:

```bash
# Create collection with semantic chunking
./maestro create collection my-database my-collection \
  --chunking-strategy=Semantic \
  --chunk-size=768 \
  --chunk-overlap=0

# The semantic strategy will:
# - Split text into sentences
# - Use AI embeddings to find semantic boundaries
# - Respect the chunk_size limit while preserving meaning
# - Default to 768 characters (vs 512 for other strategies)
```

**Note**: Semantic chunking uses sentence-transformers for chunking decisions, but the resulting chunks are embedded using your collection's embedding model (e.g., nomic-embed-text) for search operations.

Additional notes:

- Advanced semantic parameters are fully supported via flags in the CLI in addition to the common ones:
  - `--semantic-model` (model identifier, e.g., all-MiniLM-L6-v2)
  - `--semantic-window-size` (integer context window)
  - `--semantic-threshold-percentile` (0–100 split sensitivity)
  - Plus common: `--chunk-size`, `--chunk-overlap`
- Completion: the CLI provides completion for `--chunking-strategy` (includes `Semantic`). The `--semantic-model` value is free-form (no static suggestions); numeric flags disable file completion.

Example with semantic-specific flags:

```bash
./maestro collection create --vdb=my-database --name=my-collection \
   --chunking-strategy=Semantic \
   --chunk-size=768 \
   --semantic-model=all-MiniLM-L6-v2 \
   --semantic-window-size=1 \
   --semantic-threshold-percentile=95
```

### Environment Variable Substitution in YAML Files

The CLI supports environment variable substitution in YAML files using the `{{ENV_VAR_NAME}}` syntax. This allows you to use environment variables directly in your configuration files:

```yaml
apiVersion: maestro/v1alpha1
kind: VectorDatabase
metadata:
  name: my-weaviate-db
spec:
  type: weaviate
  uri: {{WEAVIATE_URL}}
  collection_name: my_collection
  embedding: text-embedding-3-small
  mode: remote
```

When you run `./maestro create vector-db config.yaml`, the CLI will:
 
1. Load environment variables from `.env` file (if present)
2. Replace `{{WEAVIATE_URL}}` with the actual value from the environment
3. Process the YAML file with the substituted values

**Features**:

- **Automatic substitution**: All `{{ENV_VAR_NAME}}` placeholders are replaced before YAML parsing
- **Error handling**: Clear error messages if required environment variables are missing
- **Verbose logging**: Shows which environment variables are being substituted (when using `--verbose`)
- **Validation**: Ensures all required environment variables are set before processing

#### URL Format Flexibility

The CLI automatically normalizes URLs to ensure they have the correct protocol prefix:

- **Hostname only**: `localhost` → `http://localhost:8030`
- **Hostname with port**: `localhost:8030` → `http://localhost:8030`
- **Full URL**: `http://localhost:8030` → `http://localhost:8030` (unchanged)
- **HTTPS URL**: `https://example.com:9000` → `https://example.com:9000` (unchanged)

This makes it easy to specify server addresses in any format:

```bash
# All of these work the same way:
./maestro vectordb list --mcp-server-uri="localhost:8030"
./maestro vectordb list --mcp-server-uri="http://localhost:8030"
./maestro vectordb list --mcp-server-uri="https://example.com:9000"
```

### Global Flags

- `--verbose`: Show detailed output
- `--silent`: Suppress success messages
- `--dry-run`: Test commands without making changes
- `--force` / `-f`: Skip confirmation prompts for destructive operations
- `--mcp-server-uri`: Override MCP server URI
- `--help`: Show help information
- `--version`: Show version information

### List Commands

The CLI provides resource-based list commands for vector databases, collections, and documents:

```bash
# List all vector databases
./maestro vectordb list

# List with verbose output
./maestro vectordb list --verbose

# Test the command without connecting to server
./maestro vectordb list --dry-run

# List collections for a specific vector database
./maestro collection list --vdb=my-database

# List collections with verbose output
./maestro collection list --vdb=my-database --verbose

# Test collections command without connecting to server
./maestro collection list --vdb=my-database --dry-run

# List embeddings for a specific vector database
./maestro embedding list --vdb=my-database

# List embeddings with verbose output
./maestro embedding list --vdb=my-database --verbose

# Test embeddings command without connecting to server
./maestro embedding list --vdb=my-database --dry-run

# List documents in a specific collection
./maestro document list --vdb=my-database --collection=my-collection

# List documents with verbose output
./maestro document list --vdb=my-database --collection=my-collection --verbose

# Test documents command without connecting to server
./maestro document list --vdb=my-database --collection=my-collection --dry-run
```

#### Output Format

**Vector Databases**: When databases are found, the output shows:
- Database name and type
- Collection name
- Document count

Example:
```text
Found 2 vector database(s):

1. project_a_db (weaviate)
   Collection: ProjectADocuments
   Documents: 15

2. project_b_db (milvus)
   Collection: ProjectBDocuments
   Documents: 8
```

**Embeddings**: When listing embeddings for a vector database, the output shows:
- Supported embedding models for the specific database type

Example:
```text
Supported embeddings for weaviate vector database 'my-database': [
  "default",
  "text2vec-weaviate",
  "text2vec-openai",
  "text2vec-cohere",
  "text2vec-huggingface",
  "text-embedding-ada-002",
  "text-embedding-3-small",
  "text-embedding-3-large"
]
```
- `--collection`: Specific collection to search in (optional; if omitted you'll be prompted interactively unless in --dry-run or non-interactive mode)

**Collections**: When listing collections for a vector database, the output shows:
- All collections available in the vector database

Example:
```text
Collections in vector database 'my-database': [
  "Collection1",
  "Collection2",
  "MaestroDocs"
]
```

**Documents**: When listing documents in a collection, the output shows:
- All documents in the specified collection with their properties

Example:
```json
Found 3 documents in collection 'my-collection' of vector database 'my-database': [
  {
    "id": "doc1",
    "url": "https://example.com/doc1",
    "text": "Document content...",
    "metadata": {
      "source": "web",
      "timestamp": "2024-01-01T00:00:00Z"
    }
  },
  {
    "id": "doc2",
    "url": "https://example.com/doc2",
    "text": "Another document...",
    "metadata": {
      "source": "file",
      "timestamp": "2024-01-02T00:00:00Z"
    }
  }
]
```

### Create Commands

The CLI provides resource-based create commands for vector databases, collections, and documents:

#### Create Vector Database Command

```bash
# Create vector database from YAML file
./maestro vdb create config.yaml

# Create with verbose output
./maestro vdb create config.yaml --verbose

# Create with dry-run mode
./maestro vdb create config.yaml --dry-run

# Override configuration values
./maestro vdb create config.yaml --type=weaviate --uri=localhost:8080
```

**Supported Override Flags**:

- `--type`: Override database type (milvus, weaviate)
- `--uri`: Override connection URI
- `--collection-name`: Override collection name
- `--embedding`: Override embedding model
- `--mode`: Override deployment mode (local, remote)

#### Create Collection Command

```bash
# Create collection in vector database
./maestro collection create --name=my-collection --vdb=my-database

# Create collection with verbose output
./maestro collection create --name=my-collection --vdb=my-database --verbose

# Create collection with dry-run mode
./maestro collection create --name=my-collection --vdb=my-database --dry-run

# Create collection with chunking configuration
./maestro collection create --name=my-collection --vdb=my-database \
   --embedding=text-embedding-3-small \
   --chunking-strategy=Sentence \
   --chunk-size=512 \
   --chunk-overlap=32
```

#### Create Document Command

```bash
# Create document from file
./maestro document create --name=my-doc --file=document.txt --vdb=my-database --collection=my-collection

# Create document with dry-run mode
./maestro document create --name=my-doc --file=document.txt --vdb=my-database --collection=my-collection --dry-run
```

### Write Command

The `write` command is an alias for creating documents:

```bash
# Write document from file
./maestro write document my-database my-collection my-doc --file-name=document.txt

# Write document using short aliases
./maestro write doc my-database my-collection my-doc --file-name=document.txt
./maestro write vdb-doc my-database my-collection my-doc --file-name=document.txt

# Write document with dry-run mode
./maestro write document my-database my-collection my-doc --file-name=document.txt --dry-run

Note: --embed on write is deprecated and ignored; embedding is configured per collection when the collection is created.
```

### Confirmation Prompts for Destructive Operations

The CLI includes safety features to prevent accidental deletion of resources. All delete operations require user confirmation unless the `--force` flag is used.

#### Confirmation Behavior

- **Interactive Confirmation**: Delete commands prompt for confirmation before proceeding
- **Force Flag**: Use `--force` or `-f` to skip confirmation prompts
- **Dry-run Mode**: Confirmation is automatically skipped in dry-run mode
- **Silent Mode**: Confirmation is automatically skipped in silent mode

#### Confirmation Examples

```bash
# Delete vector database with confirmation prompt
./maestro vectordb delete my-database
# Output: ⚠️  Are you sure you want to delete 'vector database 'my-database''? This action cannot be undone. [y/N]:

# Skip confirmation with --force flag
./maestro vectordb delete my-database --force

# Skip confirmation with -f flag
./maestro vectordb delete my-database -f

# Confirmation automatically skipped in dry-run mode
./maestro vectordb delete my-database --dry-run

# Confirmation automatically skipped in silent mode
./maestro vectordb delete my-database --silent
```

### Delete Commands

The CLI provides resource-based delete commands for vector databases, collections, and documents:

#### Delete Vector Database Command

```bash
# Delete vector database (with confirmation prompt)
./maestro vdb delete my-database

# Delete with verbose output
./maestro vdb delete my-database --verbose

# Delete with dry-run mode
./maestro vdb delete my-database --dry-run

# Skip confirmation with force flag
./maestro vdb delete my-database --force
```

#### Delete Collection Command

```bash
# Delete collection from vector database (with confirmation prompt)
./maestro collection delete my-collection --vdb=my-database

# Delete collection with verbose output
./maestro collection delete my-collection --vdb=my-database --verbose

# Delete collection with dry-run mode
./maestro collection delete my-collection --vdb=my-database --dry-run

# Skip confirmation with force flag
./maestro collection delete my-collection --vdb=my-database --force
```

#### Delete Document Command

```bash
# Delete document from collection (with confirmation prompt)
./maestro document delete my-document --vdb=my-database --collection=my-collection

# Delete document with verbose output
./maestro document delete my-document --vdb=my-database --collection=my-collection --verbose

# Delete document with dry-run mode
./maestro document delete my-document --vdb=my-database --collection=my-collection --dry-run

# Skip confirmation with force flag
./maestro document delete my-document --vdb=my-database --collection=my-collection --force
```

### Search Command

The `search` command performs a vector search and returns JSON results suitable for programmatic use.

```bash
# Basic search (prompts for collection if omitted)
./maestro search "Find information about API endpoints" --vdb=my-database

# Search with specific document limit and collection
./maestro search "quantum circuits" --vdb=my-database --collection=my-collection --doc-limit 10
```

Search output schema (normalized across backends):

- id, url, text
- metadata:
   - doc_name
   - chunk_sequence_number
   - total_chunks
   - offset_start, offset_end
   - chunk_size
- similarity: canonical score in [0..1]
- distance: cosine distance (for reference)
- rank: 1-based rank
- _metric: e.g., "cosine"
- _search_mode: "vector" or "keyword"

Flags:

- `--doc-limit, -d`: Maximum number of documents to consider (default: 5)
- `--collection`: Specific collection to search in (optional; if omitted you'll be prompted interactively unless in --dry-run or non-interactive mode)

### Query Command

The `query` command allows you to search documents using natural language queries with semantic search:

```bash
# Query documents using natural language
./maestro query "What is the main topic of the documents?" --vdb=my-database

# Query with specific document limit
./maestro query "Find information about API endpoints" --vdb=my-database --doc-limit 10

# Query with collection name specification
./maestro query "Search for technical documentation" --vdb=my-database --collection=my-collection

# Query with dry-run mode
./maestro query "Test query" --vdb=my-database --dry-run

# Query with verbose output
./maestro query "Complex search query" --vdb=my-database --verbose
```

#### Query Command Features

- **Natural Language Queries**: Use plain English to search through your documents
- **Semantic Search**: Finds relevant documents based on meaning, not just keywords
- **Document Limit Control**: Control how many documents to consider with `--doc-limit`
- **Collection Targeting**: Optionally specify which collection to search in
- **Dry-run Mode**: Test queries without actually executing them
- **Verbose Output**: Get detailed information about the query process

#### Query Command Flags

- `--doc-limit, -d`: Maximum number of documents to consider (default: 5)
- `--collection`: Specific collection to search in (optional; if omitted you'll be prompted interactively unless in --dry-run or non-interactive mode)
- `--dry-run`: Test the command without making changes
- `--verbose`: Show detailed output
- `--silent`: Suppress success messages

#### Query Examples

```bash
# Basic query
./maestro query "What is machine learning?" --vdb=my-database

# Query with higher document limit
./maestro query "Find all API documentation" --vdb=my-database --doc-limit 20

# Query specific collection
./maestro query "Search for user guides" --vdb=my-database --collection=documentation

# Test query without execution
./maestro query "Test query" --vdb=my-database --dry-run
```

### Collection Info Command

Show collection information (embedding and chunking):

```bash
./maestro collection info --vdb=my-database --name=my-collection
```

**Collection Information Output**: The output shows:
- Collection name
- Document count
- Database type
- Embedding information
- Chunking configuration (strategy and parameters)
- Additional metadata

Example:

Collection information for 'my-collection' in vector database 'my-database':

```json
{
  "name": "my-collection",
  "document_count": 15,
  "db_type": "weaviate",
  "embedding": "text2vec-weaviate",
  "chunking": {
    "strategy": "Sentence",
    "parameters": { "chunk_size": 512, "overlap": 32 }
  },
  "metadata": {
    "description": "My test collection",
    "vectorizer": "text2vec-weaviate",
    "properties_count": 4,
    "module_config": null
  }
}
```

Note: The CLI does not currently provide a standalone "document get" command.

### Validate Command

The `validate` command validates YAML configuration files:

```bash
# Validate YAML configuration
./maestro validate config.yaml

# Validate with verbose output
./maestro validate config.yaml --verbose
```

## Examples

### Complete Workflow

1. **Start the MCP server**:
```bash
   cd /path/to/maestronowledge
   ./start.sh --http
   ```

2. **List databases**:
   ```bash
   cd cli
   ./maestro vdb list --mcp-server-uri="http://localhost:8030"
   ```

3. **List with verbose output**:
   ```bash
   ./maestro vdb list --mcp-server-uri="http://localhost:8030" --verbose
   ```

4. **List collections for a database**:
   ```bash
   ./maestro collection list --vdb=my-database --mcp-server-uri="http://localhost:8030"
   ```

5. **List documents in a collection**:
   ```bash
   ./maestro document list --vdb=my-database --collection=my-collection --mcp-server-uri="http://localhost:8030"
   ```

6. **Query documents using natural language**:
   ```bash
   ./maestro query "What is the main topic?" --vdb=my-database --mcp-server-uri="http://localhost:8030"
   ./maestro query "Find API documentation" --vdb=my-database --doc-limit 10 --mcp-server-uri="http://localhost:8030"
   ```

7. **Create a vector database from YAML**:
   ```bash
   ./maestro vdb create config.yaml --mcp-server-uri="http://localhost:8030"
   ```

8. **Delete a vector database**:
   ```bash
   ./maestro vdb delete my-database --mcp-server-uri="http://localhost:8030"
   ```

### Examples

See the [examples/](examples/) directory for usage examples:

- [example_usage.sh](examples/example_usage.sh) - Comprehensive CLI usage demonstration with MCP server

## Troubleshooting

### Connection Issues

If you get connection errors:

1. **Check if the MCP server is running**:
   ```bash
   cd /path/to/maestronowledge
   ./stop.sh status
   ```

2. **Verify the server URI**:

```bash
./maestro vectordb list --mcp-server-uri="http://localhost:8030" --verbose
```

3. **Check server logs**:

```bash
tail -f /path/to/maestronowledge/mcp_server.log
```

### Common Issues

- **"connection refused"**: MCP server is not running or wrong port
- **"HTTP error 404"**: Wrong endpoint or server not configured correctly
- **"failed to parse database list"**: Server response format issue
- **"missing required environment variables"**: Environment variable substitution failed
- **"vector database already exists"**: Database with that name already exists
- **"vector database does not exist"**: Database with that name doesn't exist

## Development

### Project Structure

```text
cli/
├── src/
│   ├── main.go          # Main CLI entry point
│   ├── list.go          # List command implementation
│   ├── create.go        # Create command implementation
│   ├── delete.go        # Delete command implementation
│   ├── query.go         # Query command implementation
│   ├── validate.go      # Validate command implementation
│   └── mcp_client.go    # MCP server client
├── examples/
│   ├── example_usage.sh # Comprehensive CLI usage examples
│   └── README.md        # Examples documentation
├── tests/
│   ├── list_test.go     # List command tests
│   ├── create_test.go   # Create command tests
│   ├── delete_test.go   # Delete command tests
│   ├── query_test.go    # Query command tests
│   ├── validate_test.go # Validate command tests
│   └── main_test.go     # Main CLI tests
├── go.mod               # Go module dependencies
├── go.sum               # Go module checksums
├── tools/
│   └── lint.sh          # Comprehensive linting script
├── test_integration.sh  # Integration test script
└── README.md           # This file
```

### Code Quality and Linting

The CLI includes comprehensive linting and code quality checks to ensure maintainable, high-quality Go code.

#### Available Linting Tools

- **staticcheck**: Detects unused code, unreachable code, and other code quality issues
- **golangci-lint**: Advanced Go linting with multiple analyzers
- **go fmt**: Code formatting
- **go vet**: Static analysis
- **go mod tidy/verify**: Dependency management
- **Race condition checks**: Thread safety validation

#### Running Linting

```bash
# Run all linting checks
./tools/lint.sh

# Run specific checks
go fmt ./src/...           # Format code
go vet ./src/...           # Static analysis
staticcheck ./src/...      # Unused code detection
golangci-lint run ./src/... # Advanced linting
```

#### Linting in CI/CD

The project includes automated linting in CI/CD pipelines:

- **Main CI**: Runs CLI linting for all changes
- **CLI CI**: Dedicated CLI linting job for CLI-specific changes
- **Quality Gate**: Linting failures block merges until resolved

#### Linting Features

- **Unused Code Detection**: Automatically identifies unused variables, functions, and imports
- **Code Formatting**: Ensures consistent code style across the project
- **Static Analysis**: Catches potential bugs and code smells
- **Dependency Management**: Verifies module dependencies are clean and secure
- **Thread Safety**: Detects race conditions in concurrent code

### Adding New Commands

1. Create a new command file (e.g., `src/new_command.go`)
2. Define the command using Cobra
3. Add the command to `main.go`
4. Update this README
5. **Run linting**: `./tools/lint.sh` to ensure code quality

### Testing

```bash
# Run all tests
go test ./tests/...

# Run integration tests
./test_integration.sh

# Build and test manually
go build -o maestro src/*.go
./maestro --help

# Run with verbose output
go test -v ./tests/...
```

### Development Workflow

1. **Make changes** to CLI code
2. **Run linting**: `./tools/lint.sh` to check code quality
3. **Run tests**: `go test ./tests/...` to verify functionality
4. **Run integration tests**: `./test_integration.sh` for end-to-end validation
5. **Commit changes** with descriptive commit messages

### Code Quality Standards

- **No unused code**: All variables, functions, and imports must be used
- **Consistent formatting**: Code follows `go fmt` standards
- **Static analysis clean**: No `go vet` warnings
- **Dependency hygiene**: Clean module dependencies
- **Thread safety**: No race conditions in concurrent code

## License

Apache 2.0 License - see the main project LICENSE file for details.

## Semantic Chunking Example

The CLI supports semantic chunking for intelligent document splitting:

```bash
# Create a collection with semantic chunking
cli/maestro collection create --vdb my-vdb --name my-collection

# Check collection information to see chunking strategy
cli/maestro collection info --vdb "Qiskit_studio_algo" --name "Qiskit_studio_algo"

# Search with semantic chunking to see results
./cli/maestro search "quantum circuit" --vdb qiskit_studio_algo --collection qiskit_studio_algo --doc-limit 1
```

**Note**: The semantic chunking strategy uses sentence-transformers for chunking decisions, while the collection's own embedding model is used for search operations.
