package main

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
	"io/fs"
	"errors"
	"path"
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
	love by spf13 and friends in Go.
	Complete documentation is available at https://gohugo.io/documentation/`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("Hello World! Here are the args: %v", args)
	},
}

var syncCmd = &cobra.Command{
	Use: "sync",
	Args: cobra.ExactArgs(1),
	RunE: func (cmd *cobra.Command, args []string) error {
		currDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("Could not read current working directory: %s", err)
		}
		fmt.Printf("Current dir: %s\n", currDir)
		workDir := args[0]
		fmt.Println(workDir)
		entries, err := os.ReadDir(workDir)
		if err != nil {
			return fmt.Errorf("Could not read current workbench: %s", err)
		}

		var repos []fs.DirEntry
		for i, entry := range entries {
			if i > 5 {
				break
			}
			if entry.IsDir() {
				gitDir := fmt.Sprintf("%s/%s/.git", workDir, entry.Name()) 
				_, err := os.Stat(gitDir) 
				if err != nil && errors.Is(err, fs.ErrNotExist) {
					// .git dont exist, not a git repo
					continue
				}
				repos = append(repos, entry)
			}
		}
		fmt.Printf("%d entries. %d repos found\n", len(entries), len(repos))
		for _, repo := range repos {
			fullPath := path.Join(workDir, repo.Name())
			fmt.Printf("fetching %s...\n", repo.Name())
			err := os.Chdir(fullPath)
			if err != nil {
				return fmt.Errorf("Could not proceed sync: %s", err)
			}
			fetchCmd := exec.Command("git", "fetch")
			out, err := fetchCmd.Output()
			if err != nil {
				return fmt.Errorf("Could not execute git fetch: %s", err)
			}
			fmt.Printf("output: %s", string(out))

			// back to the initial dir
			err = os.Chdir(currDir)
			if err != nil {
				return fmt.Errorf("Could not proceed sync: %s", err)
			}
		}


		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main ()  {
	rootCmd.AddCommand(syncCmd)
	Execute()
}

