package command

import (
	"github.com/spf13/cobra"
	"github.com/kshaposhnikov/twitter-crawler/analyzer/graph"
	"github.com/kshaposhnikov/twitter-crawler/analyzer"
)

var analyzeCmd = &cobra.Command{
	Use: "analyze-graph",
	Short: "Analyze graph which is stored in a database",
	Run: analyze,
}

var contextDir string

func init() {
	analyzeCmd.Flags().StringVarP(&contextDir, "dir", "d", ".",
		"Directory for storing results of graph analysing")
}

func analyze(cmd *cobra.Command, args []string) {
	if DB == nil {
		openDBConnection()
	}

	graphIterator := graph.LoadGraphByIter(DB)
	var context analyzer.AnalyzerContext
	analyzer.CalcluatePowerByIter(graphIterator, &context)
	context.SaveContext(contextDir)
}