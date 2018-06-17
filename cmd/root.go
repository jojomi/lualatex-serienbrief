package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/template"
	script "github.com/jojomi/go-script"
	"github.com/jojomi/go-script/print"
	"github.com/jojomi/strtpl"
	"github.com/recursionpharma/go-csv-map"
	"github.com/spf13/cobra"
)

var (
	flagDataFile           string
	flagTemplateDir        string
	flagOutputDir          string
	flagTexFile            string
	flagOutputFileTemplate string

	flagVerbose bool
)

var RootCmd = &cobra.Command{
	Use: "lualatex-serienbrief",
	Run: rootCmd,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&flagDataFile, "data-file", "d", "data.csv", "data file (default is data.csv)")
	RootCmd.PersistentFlags().StringVarP(&flagTexFile, "tex-file", "l", "main.tex", "tex file (default is main.tex)")
	RootCmd.PersistentFlags().StringVarP(&flagTemplateDir, "template-dir", "t", "template", "template directory (default is template)")
	RootCmd.PersistentFlags().StringVarP(&flagOutputDir, "output-folder", "o", "output", "output directory (default is output)")
	RootCmd.PersistentFlags().StringVarP(&flagOutputFileTemplate, "output-file-template", "f", "{{ .Name }}", "template for output PDF (without extension, default is {{ .Name }})")

	RootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "show more output (e.g. while compiling document)")
}

func rootCmd(cmd *cobra.Command, args []string) {
	var (
		err error
	)
	file, err := os.Open(flagDataFile)
	if err != nil {
		log.Fatal(err)
	}
	reader := csvmap.NewReader(file)
	reader.Columns, err = reader.ReadHeader()
	if err != nil {
		log.Fatal(err)
	}
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	sc := script.NewContext()
	workingDir := sc.WorkingDir()
	buildDir := sc.AbsPath("_build_dir")
	compiledFile := strings.Replace(flagTexFile, ".tex", ".pdf", -1)
	fmt.Printf("Working dir %s...\n", workingDir)
	fmt.Printf("Build dir %s...\n", buildDir)
	for _, record := range records {
		// filter empty csv lines
		if record["Name"] == "" {
			continue
		}

		print.Boldf("Processing %s...\n", record["Name"])

		// copy template folder
		if sc.DirExists(buildDir) {
			os.RemoveAll(buildDir)
		}
		err = sc.CopyDir(flagTemplateDir, buildDir)
		if err != nil {
			log.Print(err)
			continue
		}

		// evaluate templates
		err = evaluateTemplates(buildDir, record)
		if err != nil {
			log.Print(err)
			continue
		}

		// execute lualatex
		fmt.Printf("Building document %s...\n", flagTexFile)
		sc.SetWorkingDir(buildDir)
		var execFunc func(name string, args ...string) (*script.ProcessResult, error)
		if flagVerbose {
			execFunc = sc.ExecuteDebug
		} else {
			execFunc = sc.ExecuteSilent
		}
		pr, err := execFunc("lualatex", "-interaction=nonstopmode", flagTexFile)
		if err != nil {
			log.Print(err)
			continue
		}
		if !pr.Successful() {
			log.Print("lualatex failed")
			continue
		}

		// copy output file over
		sc.SetWorkingDir(workingDir)
		outputFile := strtpl.MustEval(flagOutputFileTemplate, record) + ".pdf"
		err = sc.CopyFile(filepath.Join(buildDir, compiledFile), filepath.Join(flagOutputDir, outputFile))
		if err != nil {
			log.Print(err)
			continue
		}
	}
	// cleanup build dir
	if sc.DirExists(buildDir) {
		os.RemoveAll(buildDir)
	}
}

func evaluateTemplates(path string, data interface{}) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(info.Name())
		if ext != ".tex" && ext != ".lco" {
			return nil
		}

		// parse template
		fmt.Printf("Evaluating template %s...\n", info.Name())
		err = evalTemplate(path, data)
		if err != nil {
			return err
		}

		return nil
	})
}

func evalTemplate(path string, data interface{}) error {
	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	return nil
}
