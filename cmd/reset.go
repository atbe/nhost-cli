/*
MIT License

Copyright (c) Nhost

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// reset approval flag
var approval bool

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Delete all saved Nhost configuration from this project",
	Long: `This is an irreversible action, and should always be avoided
	unless absolutely necesarry.
	
	It will permanently delete the 'nhost/' and '.nhost/' from this project
	root and you will lose all your saved configurations with Nhost.
	
	It will, however, not cause any changes or damage to the services already
	running on remote.`,
	Run: func(cmd *cobra.Command, args []string) {

		if !approval && !contains(args, "exit") {

			// warn the user of consequences
			log.Warn("This is irreversible and will remove all installed Nhost config from this project")

			// configure interative prompt
			prompt := promptui.Prompt{
				Label:     "Are you sure you want to proceed?",
				IsConfirm: true,
			}

			response, err := prompt.Run()
			if err != nil {
				log.Debug(err)
				log.Fatal("Input prompt aborted")
			}

			if strings.ToLower(response) == "y" || strings.ToLower(response) == "approval" {
				approval = true
			}
		}

		if approval {

			if err := deleteAllPaths(dotNhost); err != nil {
				log.Debug(err)
				log.Warnf("Please delete %s manually", path.Base(dotNhost))
				log.Fatal("Failed to delete " + path.Base(dotNhost))
			}

			if err := deleteAllPaths(path.Join(workingDir, "api")); err != nil {
				log.Debug(err)
				log.Warnf("Please delete %s manually", "/api")
				log.Fatal("Failed to delete " + "/api")
			}

			if err := deleteAllPaths(nhostDir); err != nil {
				log.Debug(err)
				log.Warnf("Please delete %s manually", path.Base(nhostDir))
				log.Fatal("Failed to delete " + path.Base(nhostDir))
			}

			if err := deleteAllPaths(path.Join(workingDir, "frontend")); err != nil {
				log.Debug(err)
				log.Warnf("Please delete %s manually", "/frontend")
				log.Fatal("Failed to delete " + "/frontend")
			}

			if err := deletePath(envFile); err != nil {
				log.Debug(err)
				log.Warnf("Please delete %s manually", path.Base(envFile))
				log.Fatal("Failed to delete " + path.Base(envFile))
			}
		}

		// signify reset completion
		log.Infof("Directories permanently removed from this project: %v, %v, %v, %v, %v",
			path.Base(nhostDir), path.Base(dotNhost), path.Base(envFile), "/api", "/frontend")

		// if an exit argument has been passed,
		// provide a graceful exit
		if contains(args, "exit") {
			log.Info("Reset complete. See you later, grasshopper!")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	// Here you will define your flags and configuration settings.
	resetCmd.PersistentFlags().BoolVarP(&approval, "approval", "y", false, "Bypass approval prompt and proceed ahead with reset")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
