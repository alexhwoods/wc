/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func processFile(file string) (lines int, err error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	lines = 0

	scanner := bufio.NewScanner(f)
	// scans one line at a time
	for scanner.Scan() {
		// _ := scanner.Text()

		lines++
	}

	if err := scanner.Err(); err != nil {
		return lines, err
	}

	return lines, nil
}

type FileParseResult struct {
	filename string
	lines    int
	bytes    int
}

func (f FileParseResult) String() string {
	return "lines: " + strconv.Itoa(f.lines) + " bytes: " + strconv.Itoa(f.bytes) + " " + f.filename
}

var rootCmd = &cobra.Command{
	Use:   "wc",
	Short: "word, line, character, and byte count",
	Long:  `A clone of the wc command in Unix. Do "man wc" for more information.`,
	RunE: func(cmd *cobra.Command, files []string) error {
		// bytesEnabled, _ := cmd.Flags().GetBool("bytes")
		// linesEnabled, _ := cmd.Flags().GetBool("lines")

		fileParseResults := make([]FileParseResult, len(files))

		for i, file := range files {
			fileInfo, err := os.Lstat(file)
			check(err)

			lines, err := processFile(file)
			check(err)

			fileParseResults[i] = FileParseResult{
				filename: file,
				lines:    lines,
				bytes:    int(fileInfo.Size()),
			}
		}

		for _, fileParseResult := range fileParseResults {
			fmt.Println(fileParseResult)
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

	rootCmd.Flags().BoolP("lines", "l", false, "The number of lines in each input file is written to the standard output.")
}
