Maestro Knowledge

Usage:
  maestro validate YAML_FILE [options]
  maestro validate SCHEMA_FILE YAML_FILE [options]

  maestro vectordb list [options]
  maestro vectordb create YAML_FILE [options]
  maestro vectordb delete NAME [options]

  maestro collection list --vdb=VDB_NAME [options]
  maestro collection create --name=COLLECTION_NAME --vdb=VDB_NAME [options]
  maestro collection delete COLLECTION_NAME --vdb=VDB_NAME [options]

  maestro embedding list --vdb=VDB_NAME [options]

  maestro document list --vdb=VDB_NAME --collection=COLLECTION_NAME [options]
  maestro document create --name=DOC_NAME --file=FILE_PATH --vdb=VDB_NAME --collection=COLLECTION_NAME [options]
  maestro document delete DOC_NAME --vdb=VDB_NAME --collection=COLLECTION_NAME [options]

  maestro query "QUERY_STRING" --vdb=VDB_NAME [options]

  maestro (-h | --help)
  maestro (-v | --version)

Options:
  --verbose              Show all output.
  --silent               Show no additional output on success, e.g., no OK or Success etc
  --dry-run              Mocks agents and other parts of workflow execution.

  -h --help              Show this screen.
  -v --version           Show version.