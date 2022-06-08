package docs

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "generate various documentations on how to use the hcloud client",
	Run:   genDocs,
}

var root *cobra.Command

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(docsCmd)

	root = rootCmd
}

func genDocs(cmd *cobra.Command, args []string) {
	doc.GenManTree(root, &doc.GenManHeader{Title: "hcloud"}, "/Users/sev/Desktop/hcloud")
}
