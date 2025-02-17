package cmd

import (
	"os"

	"github.com/bndrmrtn/zxl/internal/lexer"
	"github.com/bndrmrtn/zxl/pkg/prettycode"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// rootCmd is the root command for the CLI
var highlightCmd = &cobra.Command{
	Use:     "highligh file.zx",
	Aliases: []string{"hl"},
	Short:   "Highlight Zx (.zx) files",
	Run:     execHighlight,
}

func init() {
	// Add the run command to the root command
	rootCmd.AddCommand(highlightCmd)
	highlightCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
}

func execHighlight(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	colors := cmd.Flag("nocolor").Value.String() == "false"
	if !colors {
		color.NoColor = true
	}

	if len(args) == 0 {
		cmd.PrintErr("No file specified")
		return
	}

	if len(args) > 1 {
		cmd.PrintErr("Only one file can be run at a time")
		return
	}

	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		cmd.PrintErr("File does not exist")
		return
	}

	file, err := os.Open(args[0])
	if err != nil {
		cmd.PrintErr(err)
		return
	}
	defer file.Close()

	lx := lexer.New(args[0])
	tokens, err := lx.Parse(file)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	pretty := prettycode.New(tokens)
	html := pretty.HighlightHtml()
	if err := os.WriteFile(args[0]+".html", []byte(html), os.ModePerm); err != nil {
		cmd.PrintErr(err)
		return
	}
}
