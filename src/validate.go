package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

var validateCmd = &cobra.Command{
	Use:   "validate [YAML_FILE] [SCHEMA_FILE]",
	Short: "Validate YAML files against schemas",
	Long: `Validate YAML files against schemas.
	
Usage:
  maestro validate YAML_FILE [options]                    # Uses default schema
  maestro validate SCHEMA_FILE YAML_FILE [options]        # Uses custom schema

Examples:
  maestro validate config.yaml                            # Validates against default schema
  maestro validate custom-schema.json config.yaml         # Validates against custom schema`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Suppress usage for all errors except usage errors
		cmd.SilenceUsage = true

		var yamlFile, schemaFile string

		if len(args) == 1 {
			// Single argument: validate YAML_FILE
			yamlFile = args[0]
			schemaFile = "" // Will use default schema
		} else {
			// Two arguments: validate SCHEMA_FILE YAML_FILE
			schemaFile = args[0]
			yamlFile = args[1]
		}

		return validateFiles(yamlFile, schemaFile)
	},
}

// downloadSchemaFromRepo downloads the schema file from the maestro-knowledge repository
func downloadSchemaFromRepo(schemaPath string) error {
	// Create the schemas directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(schemaPath), 0755); err != nil {
		return fmt.Errorf("failed to create schemas directory: %v", err)
	}

	// URL to the schema file in the maestro-knowledge repository
	schemaURL := "https://raw.githubusercontent.com/AI4quantum/maestro-knowledge/refs/heads/main/schemas/vector-database-schema.json"

	if verbose {
		fmt.Printf("Downloading schema from: %s\n", schemaURL)
	}

	// Download the schema file
	resp, err := http.Get(schemaURL)
	if err != nil {
		return fmt.Errorf("failed to download schema: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download schema: HTTP %d", resp.StatusCode)
	}

	// Create the local schema file
	file, err := os.Create(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to create schema file: %v", err)
	}
	defer file.Close()

	// Copy the content
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write schema file: %v", err)
	}

	if verbose {
		fmt.Printf("Schema downloaded successfully to: %s\n", schemaPath)
	}

	return nil
}

func validateFiles(yamlFile, schemaFile string) error {
	// If no schema file provided, use the default schema
	if schemaFile == "" {
		// Try multiple possible locations for the default schema
		possiblePaths := []string{
			"schemas/vector-database-schema.json",       // From project root
			"../schemas/vector-database-schema.json",    // From cli directory
			"../../schemas/vector-database-schema.json", // From cli/tests directory
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				schemaFile = path
				break
			}
		}

		// If no schema found, try to download it from the maestro-knowledge repo
		if schemaFile == "" {
			schemaFile = "schemas/vector-database-schema.json"

			// Try to download the schema if it doesn't exist
			if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
				if err := downloadSchemaFromRepo(schemaFile); err != nil {
					// If download fails, continue with the default path
					// The error will be caught later when trying to read the file
					if verbose {
						fmt.Printf("Warning: Failed to download schema from maestro-knowledge repository: %v\n", err)
						fmt.Printf("Please ensure the schema file exists at the expected location or provide a custom schema file.\n")
					}
				}
			}
		}
	}

	if verbose {
		fmt.Printf("Validating YAML file: %s\n", yamlFile)
		fmt.Printf("Using schema file: %s\n", schemaFile)
	}

	// Check if YAML file exists
	if _, err := os.Stat(yamlFile); os.IsNotExist(err) {
		return fmt.Errorf("yaml file not found: %s", yamlFile)
	}

	// Check if schema file exists
	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return fmt.Errorf("schema file not found: %s", schemaFile)
	}

	if dryRun {
		fmt.Println("[DRY RUN] Would validate files")
		return nil
	}

	// Perform validation
	if err := performValidation(yamlFile, schemaFile); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if !silent {
		fmt.Println("✅ Validation successful")
	}

	return nil
}

func performValidation(yamlFile, schemaFile string) error {
	if verbose {
		fmt.Printf("Performing validation of %s\n", filepath.Base(yamlFile))
	}

	// Step 1: Parse the YAML file
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Parse YAML to check if it's valid YAML
	var yamlContent interface{}
	if err := yaml.Unmarshal(yamlData, &yamlContent); err != nil {
		return fmt.Errorf("invalid yaml format: %w", err)
	}

	// Step 2: Always validate against schema (default or provided)

	// Step 3: Load and validate against JSON schema
	// Convert to absolute path for schema file
	absSchemaPath, err := filepath.Abs(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for schema: %w", err)
	}

	// Convert YAML to JSON for validation
	jsonData, err := yamlToJSON(yamlData)
	if err != nil {
		return fmt.Errorf("failed to convert yaml to JSON: %w", err)
	}

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + absSchemaPath)
	documentLoader := gojsonschema.NewStringLoader(string(jsonData))

	// Validate the document against the schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("schema validation error: %w", err)
	}

	// Step 4: Check validation results
	if !result.Valid() {
		if verbose {
			fmt.Println("Validation failed:")
		}

		for _, err := range result.Errors() {
			errorMsg := fmt.Sprintf("- %s: %s", err.Field(), err.Description())
			if verbose {
				fmt.Println(errorMsg)
			}
		}

		return fmt.Errorf("validation failed with %d errors", len(result.Errors()))
	}

	if verbose {
		fmt.Println("Schema validation passed")
	}

	return nil
}

// yamlToJSON converts YAML data to JSON format
func yamlToJSON(yamlData []byte) ([]byte, error) {
	var yamlContent interface{}

	// Parse YAML
	if err := yaml.Unmarshal(yamlData, &yamlContent); err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(yamlContent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return jsonData, nil
}
