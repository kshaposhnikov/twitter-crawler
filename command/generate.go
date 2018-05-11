package command

import (
	"github.com/spf13/cobra"
	"time"
	"github.com/kshaposhnikov/twitter-crawler/generator"
	"regexp"
	"strings"
	"strconv"
	"github.com/sirupsen/logrus"
)

var generateCmd = &cobra.Command{
	Use: "generate-graph",
	Short: "Generate graph using Bollobas-Riordan method",
	Run: generate,
}

var graphConfig string
var threadCount int

func init() {
	generateCmd.Flags().StringVarP(&graphConfig, "config", "c", "", "Config has format \"VertexCount;M\"")
	generateCmd.Flags().IntVarP(&threadCount, "threads", "t", 1, "Number of threads")
}

func generate(cmd *cobra.Command, args []string) {
	template := regexp.MustCompile(`[0-9]+;[0-9]+`)
	if template.MatchString(graphConfig) {
		nm := strings.Split(template.FindString(graphConfig), ";")
		n, _ := strconv.Atoi(nm[0])
		m, _ := strconv.Atoi(nm[1])

		logrus.WithFields(logrus.Fields{
			"n": n,
			"m": m,
			"threads": threadCount,
		}).Debug("[Generate Command] Input parameters")

		start := time.Now()
		result := generator.SecondPhaseMultithreadGenerator{
			GeneralGenerator: generator.GeneralGenerator{
				VCount: n,
				ECount: m,
			},
			NumberThreads: threadCount,
		}.Generate()

		logrus.WithField("duration", time.Now().Sub(start)).Info("[Generate Command] Generation is done")
		logrus.Debug("[Generate Command] Final Graph", result)
	} else {
		logrus.Error("Need to specify format `n;m`")
	}
}
