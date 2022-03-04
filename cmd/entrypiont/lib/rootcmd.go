package lib

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var rootCmd = &cobra.Command{
	Use:   "entrypoint",
	Short: "Entrypoint",
	Long:  "自定义Entrypoint",
	Run: func(cmd *cobra.Command, args []string) {
		//  core logic
		CheckFlags()
		CheckWaitFile()
		fmt.Println("core logic")
	},
}

func InitCommand() {
	rootCmd.Flags().StringVar(&waitFile, "wait", "", "entrypoint --wait")
	rootCmd.Flags().StringVar(&out, "out", "", "entrypoint --out /path/to/name")
	rootCmd.Flags().StringVar(&command, "cmd", "", "entrypoint --cmd")
	if err := rootCmd.Execute(); err != nil {
		klog.Exit(err)
	}
}
