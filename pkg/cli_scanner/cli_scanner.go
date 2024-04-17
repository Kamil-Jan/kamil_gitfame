package cli_scanner

import (
	"github.com/spf13/cobra"
)

type Settings struct {
	Repository   string
	Revision     string
	OrderBy      string
	UseCommitter bool
	Format       string
	Extensions   string
	Languages    string
	Exclude      string
	RestrictTo   string
}

func Scan(args []string) *Settings {
	var settings Settings

	var rootCmd = &cobra.Command{
		Use: "gitfame",
		Run: func(cmd *cobra.Command, args []string) {
			settings.Repository, _ = cmd.Flags().GetString("repository")
			settings.Revision, _ = cmd.Flags().GetString("revision")
			settings.OrderBy, _ = cmd.Flags().GetString("order-by")
			settings.UseCommitter, _ = cmd.Flags().GetBool("use-committer")
			settings.Format, _ = cmd.Flags().GetString("format")
			settings.Extensions, _ = cmd.Flags().GetString("extensions")
			settings.Languages, _ = cmd.Flags().GetString("languages")
			settings.Exclude, _ = cmd.Flags().GetString("exclude")
			settings.RestrictTo, _ = cmd.Flags().GetString("restrict-to")
		},
	}

	rootCmd.Flags().StringP("repository", "r", ".", "Path to Git repository")
	rootCmd.Flags().StringP("revision", "", "HEAD", "Git revision")
	rootCmd.Flags().StringP("order-by", "", "lines", "Sort results by 'lines', 'commits', or 'files'")
	rootCmd.Flags().BoolP("use-committer", "", false, "Use committer instead of author in calculations")
	rootCmd.Flags().StringP("format", "", "tabular", "Output format: 'tabular', 'csv', 'json', 'json-lines'")
	rootCmd.Flags().StringP("extensions", "", "", "List of file extensions to include")
	rootCmd.Flags().StringP("languages", "", "", "List of programming languages to include")
	rootCmd.Flags().StringP("exclude", "", "", "Glob patterns to exclude files")
	rootCmd.Flags().StringP("restrict-to", "", "", "Glob patterns to include files")

	rootCmd.SetArgs(args)
	_ = rootCmd.Execute()

	return &settings
}
