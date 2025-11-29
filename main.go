package main

import (
	"fmt"
	"os"

	"github.com/mrjosh/icloud-mailbox/cmd"
	"github.com/spf13/cobra"
)

func main() {

	defer cmd.PanicHandler()
	rootCmd := &cobra.Command{
		Use: "mailbox",
		Long: `
              _ _   ___           
  /\/\   __ _(_) | / __\ _____  __
 /    \ / _' | | |/__\/// _ \ \/ /
/ /\/\ \ (_| | | / \/  \ (_) >  < 
\/    \/\__,_|_|_\_____/\___/_/\_\`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.SetArgs(os.Args[1:])
	if err := cmd.Start(rootCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
