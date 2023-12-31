/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Can't read config:", err)
		}

		ecName := viper.GetString("ECNAME")
		//docker compose -f "$DATA_DIR/compose.yaml" up -d $ECNAME
		command := exec.Command("docker", "compose", "--file", filepath.Join(dataDir, "compose.yaml"), "up", "-d", ecName)
		fmt.Println(command.Args)
		for k, v := range viper.AllSettings() {
			env := fmt.Sprintf("%s=%s", strings.ToUpper(k), v)
			fmt.Println(env)
			command.Env = append(command.Env, env)
		}

		stdout, err := command.StdoutPipe()
		command.Stderr = command.Stdout
		if err != nil {
			log.Error(err)
		}
		if err = command.Start(); err != nil {
			log.Fatal(err)
		}
		for {
			tmp := make([]byte, 1024)
			_, err := stdout.Read(tmp)
			fmt.Print(string(tmp))
			if err != nil {
				break
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
