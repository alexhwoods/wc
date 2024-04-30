/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wc",
	Short: "word, line, character, and byte count",
	Long:  `A clone of the wc command in Unix. Do "man wc" for more information.`,
	RunE: func(cmd *cobra.Command, files []string) error {
		bytes, _ := cmd.Flags().GetBool("bytes")

		for _, file := range files {
			fileInfo, err := os.Lstat(file)

			if err != nil {
				fmt.Println(err)
				return err
			}

			line := ""

			if bytes {
				line = line + strconv.FormatInt(fileInfo.Size(), 10) + " "
			}

			line = line + file

			fmt.Println(line)
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("bytes", "c", false, "The number of bytes in each input file is written to the standard output.  This will cancel out any prior usage of the -m option.")
}
