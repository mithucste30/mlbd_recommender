package cmd

import "github.com/spf13/cobra"

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.recommender.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().IntVarP(&Port, "port", "p",8080, "Server running port")
	serverCmd.Flags().StringVarP(&RedisHost, "redisHost", "r", "redis://redis:6379", "Redis host address with port and protocol(e.g. redis://redis:6379)")

	batchUpdateCmd.Flags().IntVarP(&Port, "port", "p",8080, "Server running port")
	batchUpdateCmd.Flags().StringVarP(&RedisHost, "redisHost", "r", "redis://redis:6379", "Redis host address with port and protocol(e.g. redis://redis:6379)")
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(batchUpdateCmd)
}
