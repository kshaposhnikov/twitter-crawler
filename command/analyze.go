package command

import (
	"github.com/kshaposhnikov/twitter-crawler/analyzer"
	"github.com/kshaposhnikov/twitter-crawler/graph"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze-graph",
	Short: "Analyze graph which is stored in a database",
	Run:   analyze,
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

	gateway := graph.NewGateway(DB)
	graphIterator := gateway.LoadGraphByIter()
	var context analyzer.AnalyzerContext
	analyzer.CalcluatePowerByIter(graphIterator, &context)
	context.SaveContext(contextDir)
}
