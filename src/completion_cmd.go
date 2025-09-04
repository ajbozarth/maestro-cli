package main

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
  $ source <(maestro completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ maestro completion bash > /etc/bash_completion.d/maestro
  # macOS:
  $ maestro completion bash > /usr/local/etc/bash_completion.d/maestro

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ maestro completion zsh > "${fpath[1]}/_maestro"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ maestro completion fish | source

  # To load completions for each session, execute once:
  $ maestro completion fish > ~/.config/fish/completions/maestro.fish

PowerShell:
  PS> maestro completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> maestro completion powershell > maestro.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

// SetupCustomCompletions sets up custom completion functions for commands
func SetupCustomCompletions() {
	// Vector database completion for --vdb flag
	if vdbListCmd != nil {
		vdbListCmd.Flags().String("vdb", "", "Vector database name")
		vdbListCmd.RegisterFlagCompletionFunc("vdb", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			provider := NewCompletionProvider()
			completions, err := provider.CompleteVectorDatabases(toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			var results []string
			for _, completion := range completions {
				results = append(results, completion.Text)
			}
			return results, cobra.ShellCompDirectiveDefault
		})
	}

	// Collection completion for --collection flag
	if documentListCmd != nil {
		documentListCmd.RegisterFlagCompletionFunc("collection", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// Get VDB name from --vdb flag
			vdbName, _ := cmd.Flags().GetString("vdb")
			if vdbName == "" {
				return nil, cobra.ShellCompDirectiveError
			}

			provider := NewCompletionProvider()
			completions, err := provider.CompleteCollections(vdbName, toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			var results []string
			for _, completion := range completions {
				results = append(results, completion.Text)
			}
			return results, cobra.ShellCompDirectiveDefault
		})
	}

	// File completion for --file flag
	if documentCreateCmd != nil {
		documentCreateCmd.RegisterFlagCompletionFunc("file", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			provider := NewCompletionProvider()
			completions, err := provider.CompleteFiles(toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			var results []string
			for _, completion := range completions {
				results = append(results, completion.Text)
			}
			return results, cobra.ShellCompDirectiveDefault
		})
	}

	// Embedding completion for --embedding flag
	if collectionCreateCmd != nil {
		collectionCreateCmd.RegisterFlagCompletionFunc("embedding", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			provider := NewCompletionProvider()
			completions, err := provider.CompleteEmbeddings(toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			var results []string
			for _, completion := range completions {
				results = append(results, completion.Text)
			}
			return results, cobra.ShellCompDirectiveDefault
		})

		// Chunking strategy completion
		collectionCreateCmd.RegisterFlagCompletionFunc("chunking-strategy", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			provider := NewCompletionProvider()
			completions, err := provider.CompleteChunkingStrategies(toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			var results []string
			for _, completion := range completions {
				results = append(results, completion.Text)
			}
			return results, cobra.ShellCompDirectiveDefault
		})

		// Chunk size and overlap are numeric; we don't compute values but we can hint directive to file completion disabled
		collectionCreateCmd.RegisterFlagCompletionFunc("chunk-size", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		})
		collectionCreateCmd.RegisterFlagCompletionFunc("chunk-overlap", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		})

		// Semantic-specific flags
		collectionCreateCmd.RegisterFlagCompletionFunc("semantic-model", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		})
		collectionCreateCmd.RegisterFlagCompletionFunc("semantic-window-size", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		})
		collectionCreateCmd.RegisterFlagCompletionFunc("semantic-threshold-percentile", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		})
	}

	// Status --vdb completion using vector database names
	if statusCmd != nil {
		statusCmd.RegisterFlagCompletionFunc("vdb", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			provider := NewCompletionProvider()
			completions, err := provider.CompleteVectorDatabases(toComplete)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			var results []string
			for _, completion := range completions {
				results = append(results, completion.Text)
			}
			return results, cobra.ShellCompDirectiveDefault
		})
	}
}

// AddCompletionCommand adds the completion command to the root command
func AddCompletionCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(completionCmd)
}
