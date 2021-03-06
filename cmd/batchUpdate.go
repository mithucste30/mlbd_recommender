package cmd

import (
	"fmt"
	"github.com/mithucste30/mlbd_recommender/app"
	"github.com/spf13/cobra"
)

var (
	batchUpdateCmd = &cobra.Command{
		Use: "batch_update",
		Short: "Run recommender batch update operation.",
		Long: `This operation updates all the user's recommendation.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting batch update.")
			app.BatchUpdate(Port, RedisHost)
			fmt.Println("Batch update completed successfully.")
		},
	}
)
