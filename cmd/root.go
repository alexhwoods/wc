package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getCounts(rd io.Reader, name string) (FileParseResult, error) {
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
		lines:    lines,
		words:    words,
		chars:    chars,
		bytes:    bytes,
		filename: name,
	}, nil
}

type FileParseResult struct {
	filename string
	lines    int
	words    int
	chars    int
	bytes    int
}

func (f FileParseResult) Println(bytesFlag bool, linesFlag bool, wordsFlag bool, charsEnabled bool, allFlagsDisabled bool) {
	s := ""

	if linesFlag || allFlagsDisabled {
		s += fmt.Sprintf("%8d", f.lines)
	}
	if wordsFlag || allFlagsDisabled {
		s += fmt.Sprintf("%8d", f.words)
	}
	if charsEnabled {
		s += fmt.Sprintf("%8d", f.chars)
	}
	if bytesFlag || allFlagsDisabled {
		s += fmt.Sprintf("%8d", f.bytes)
	}

	fmt.Printf(s + " " + f.filename + "\n")
}

var rootCmd = &cobra.Command{
	Use:   "wc",
	Short: "word, line, character, and byte count",
	Long:  `A clone of the wc command in Unix. Do "man wc" for more information.`,
	RunE: func(cmd *cobra.Command, files []string) error {
		bytesFlag, _ := cmd.Flags().GetBool("bytes")
		linesFlag, _ := cmd.Flags().GetBool("lines")
		wordsFlag, _ := cmd.Flags().GetBool("words")
		charsFlag, _ := cmd.Flags().GetBool("chars")

		// @note: I'm varying from official wc behavior here.
		//        They will take the last of -c and -m if both are used.
		//        I'm simply going to use -m if both are used.
		//        Cobra does not have a simple way of getting the order.
		//        Also, I really dislike -cm and -mc giving different behavior.
		charsEnabled := charsFlag && !bytesFlag
		allFlagsDisabled := !bytesFlag && !linesFlag && !wordsFlag && !charsFlag

		totalLines := 0
		totalWords := 0
		totalChars := 0
		totalBytes := 0

		if len(files) == 0 {
			reader := bufio.NewReader(os.Stdin)
			fileParseResult, err := getCounts(reader, "")
			check(err)

			totalLines += fileParseResult.lines
			totalWords += fileParseResult.words
			totalChars += fileParseResult.chars
			totalBytes += fileParseResult.bytes

		}

		for _, file := range files {
			fileReader, err := os.Open(file)
			check(err)

			fileParseResult, err := getCounts(fileReader, file)
			check(err)

			fileParseResult.Println(bytesFlag, linesFlag, wordsFlag, charsEnabled, allFlagsDisabled)

			totalLines += fileParseResult.lines
			totalWords += fileParseResult.words
			totalChars += fileParseResult.chars
			totalBytes += fileParseResult.bytes
		}

		if len(files) > 1 {
			totalResult := FileParseResult{
				lines:    totalLines,
				words:    totalWords,
				chars:    totalChars,
				bytes:    totalBytes,
				filename: "total",
			}

			totalResult.Println(bytesFlag, linesFlag, wordsFlag, charsEnabled, allFlagsDisabled)
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

	rootCmd.Flags().BoolP("words", "w", false, "The number of words in each input file is written to the standard output.")

	rootCmd.Flags().BoolP("chars", "m", false, "The number of characters in each input file is written to the standard output.  If the current locale does not support multibyte characters, this is equivalent to the -c option.  This will cancel out any prior usage of the -c option.")
}
