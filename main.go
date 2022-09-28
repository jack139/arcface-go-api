package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"

	"github.com/jack139/go-infer/cli"
	"github.com/jack139/go-infer/types"
	"github.com/jack139/go-infer/helper"

	// 推理模型
	"arcface-go-api/models/arcface"
)


var (
	rootCmd = &cobra.Command{
		Use:   "arcface-api",
		Short: "Arcface-go API",
	}
)


func init() {
	// 重载 PersistentPreRunE
	cli.HttpCmd.PersistentPreRunE = preRun
	cli.ServerCmd.PersistentPreRunE = preRun

	// 命令行设置
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(cli.HttpCmd)
	rootCmd.AddCommand(cli.ServerCmd)
}


func preRun(cmd *cobra.Command, args []string) error {
	yaml, _ := cmd.Flags().GetString("yaml")
	helper.InitSettings(yaml)

	// 初始化时根据配置文件加载模型
	types.ModelList = append(types.ModelList, &arcface.FaceLocate{})
	types.ModelList = append(types.ModelList, &arcface.FaceVerify{})
	types.ModelList = append(types.ModelList, &arcface.FaceFeatures{})

	return nil
}


func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

