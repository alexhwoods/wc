/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func processFile(file string) (lines int, words int, chars int, err error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, 0, 0, err
	}
	defer f.Close()

	// @note: cannot use scanner because new line characters
	//        are stripped, and \n vs. \n\r affects the char count
	reader := bufio.NewReader(f)
	lines = 0
	words = 0
	chars = 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return lines, words, chars, err
		}

		// @note: will count an extra line if the file ends with a newline
		if err == io.EOF && len(line) == 0 {
			break
		}

		lines++
		words += len(strings.Fields(line))
		chars += utf8.RuneCountInString(line)

		if err == io.EOF {
			break
		}
	}

	return lines, words, chars, nil
}

type FileParseResult struct {
	filename string
	lines    int
	words    int
	chars    int
	bytes    int
}

func (f FileParseResult) String() string {
	return "lines: " + strconv.Itoa(f.lines) + " bytes: " + strconv.Itoa(f.bytes) + " chars: " + strconv.Itoa(f.chars) +
		" words: " + strconv.Itoa(f.words) +
		" " + f.filename
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

			lines, words, chars, err := processFile(file)
			check(err)

			fileParseResults[i] = FileParseResult{
				filename: file,
				lines:    lines,
				words:    words,
				chars:    chars,
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
