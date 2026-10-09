package main

import (
	"os"
	"time"

	"github.com/lehatrutenb/course-go-python-2026/homeworks/hw1/dice/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	zapCfg := zap.NewProductionConfig()
	zapCfg.Level.SetLevel(zapcore.InfoLevel)
	lg, _ := zapCfg.Build()

	defer lg.Sync()
	sugar := lg.Sugar()

	var verbose bool
	var to_sum bool
	var count int
	var seed int64

	var rootCmd = &cobra.Command{
		Use:   "dice",
		Short: "Dice rolling",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := config.NewConfig(sugar, count, seed, to_sum, verbose, args[0])
			if err != nil {
				sugar.Errorw("failed to create config", "error", err)
				os.Exit(1)
			}
			sugar.Infow("Run command", "config", config)
		},
		// RunE: func(cmd *cobra.Command, args []string) error { TODO do we really need this?
		// 	return nil
		// },
	}

	rootCmd.Flags().IntVarP(&count, "count", "c", 1, "Number of iterations (default: 1)")
	rootCmd.Flags().Int64Var(&seed, "seed", time.Now().UnixNano(), "Seed for the random number generator")
	rootCmd.Flags().BoolVarP(&to_sum, "sum", "s", false, "To calculate the sum of the rolls (default: false)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "To print the rolls (default: false)")

	sugar.Info("Starting dice rolling")
	if err := rootCmd.Execute(); err != nil {
		sugar.Error(err)
		os.Exit(1)
	}
}
