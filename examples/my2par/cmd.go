package main

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var (
	outputDir string
	dsn       string
	rowLimit  int
	pageSize  int
)

var rootCmd = &cobra.Command{
	Use:   path.Base(os.Args[0]),
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&outputDir, "output", "", "./output", "")
	flags.StringVarP(&dsn, "dsn", "", os.Getenv("BIKIT_DSN"), "example: root:@tcp(127.0.0.1:3306)/db_name")
	flags.IntVarP(&rowLimit, "limit", "", -1, "")
	flags.IntVarP(&pageSize, "page-size", "", 100000, "")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
