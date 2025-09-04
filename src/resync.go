package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var resyncCmd = &cobra.Command{
	Use:   "resync-databases",
	Short: "Discover and register existing Milvus collections into MCP server memory",
	Long:  `Trigger an on-demand resynchronization so the MCP server registers any Milvus collections created outside this process.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURI, err := getMCPServerURI(mcpServerURI)
		if err != nil {
			return err
		}

		client, err := NewMCPClient(serverURI)
		if err != nil {
			return err
		}
		defer client.Close()

		result, err := client.ResyncDatabases()
		if err != nil {
			return fmt.Errorf("resync failed: %w", err)
		}

		fmt.Println(result)
		return nil
	},
}

// Note: the command is registered from main.go (where rootCmd is defined)
