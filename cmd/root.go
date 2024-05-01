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

func getCounts(rd io.Reader) (FileParseResult, error) {
	// @note: cannot use scanner because new line characters
	//        are stripped, and \n vs. \n\r affects the char count
	reader := bufio.NewReader(rd)
	lines := 0
	words := 0
	chars := 0
	bytes := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return FileParseResult{}, err
		}

		// @note: will count an extra line if the file ends with a newline
		if err == io.EOF && len(line) == 0 {
			break
		}

		lines++
		words += len(strings.Fields(line))
		chars += utf8.RuneCountInString(line)
		bytes += len(line)

		if err == io.EOF {
			break
		}
	}

	return FileParseResult{
		lines: lines,
		words: words,
		chars: chars,
		bytes: bytes,
	}, nil
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
		fileParseResults := []FileParseResult{}

		if len(files) == 0 {
			reader := bufio.NewReader(os.Stdin)
			fileParseResult, err := getCounts(reader)
			check(err)

			fileParseResults = append(fileParseResults, fileParseResult)
		}

		for _, file := range files {
			fileReader, err := os.Open(file)
			check(err)

			fileParseResult, err := getCounts(fileReader)
			check(err)

			fileParseResults = append(fileParseResults, fileParseResult)
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
