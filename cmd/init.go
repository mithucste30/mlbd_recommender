package cmd

func init() {
	serverCmd.Flags().IntVarP(&Port, "port", "p",8000, "Server running port")
	serverCmd.Flags().StringVarP(&RedisHost, "redisHost", "r", "redis://redis:6379", "Redis host address with port and protocol(e.g. redis://redis:6379)")
	serverCmd.Flags().BoolVarP(&Doc, "doc", "d", true, "Specify whether to mount documentation or not(true/false). Default is true")

	batchUpdateCmd.Flags().IntVarP(&Port, "port", "p",8000, "Server running port")
	batchUpdateCmd.Flags().StringVarP(&RedisHost, "redisHost", "r", "redis://redis:6379", "Redis host address with port and protocol(e.g. redis://redis:6379)")
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(batchUpdateCmd)
}
