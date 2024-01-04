/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nodeset-org/hyperdrive-stakewise/config"
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

		c, err := config.LoadConfig()
		if err != nil {
			log.Fatal("Can't read config:", err)
		}

		command := exec.Command("docker", "compose", "--file", filepath.Join(c.DataDir, "compose.yaml"), "up", "-d", c.ExceutionClientName)

		command.Dir = c.DataDir
		command.Env = append(command.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
		command.Env = append(command.Env, fmt.Sprintf("DATA_DIR=%s", c.DataDir))
		for k, v := range viper.AllSettings() {
			env := fmt.Sprintf("%s=%s", strings.ToUpper(k), v)
			command.Env = append(command.Env, env)
		}

		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		// Start the command
		if err := command.Start(); err != nil {
			log.Error(err)
		}

		command.Wait()

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
