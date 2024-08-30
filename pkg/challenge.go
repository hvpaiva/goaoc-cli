package pkg

import (
	"fmt"
	"html/template"
	"os"

	"github.com/spf13/cobra"

	"github.com/hvpaiva/goaoc-cli/tpl"
)

type Challenge struct {
	*Project
	Day  int
	Year int
}

func (c *Challenge) Create() error {
	if _, err := os.Stat(fmt.Sprintf("%s/internal/%d", c.AbsolutePath, c.Year)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/internal/%d", c.AbsolutePath, c.Year), 0o751))
	}

	if _, err := os.Stat(fmt.Sprintf("%s/internal/%d/day%02d", c.AbsolutePath, c.Year, c.Day)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/internal/%d/day%02d", c.AbsolutePath, c.Year, c.Day), 0o751))
	}

	mainFile, err := os.Create(fmt.Sprintf("%s/internal/%d/day%02d/main.go", c.AbsolutePath, c.Year, c.Day))
	if err != nil {
		return err
	}
	defer func(mainFile *os.File) {
		err := mainFile.Close()
		if err != nil {
			panic(err)
		}
	}(mainFile)

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.AddMainTemplate())))

	err = mainTemplate.Execute(mainFile, c)
	if err != nil {
		return err
	}

	mainTestFile, err := os.Create(fmt.Sprintf("%s/internal/%d/day%02d/main_test.go", c.AbsolutePath, c.Year, c.Day))
	if err != nil {
		return err
	}
	defer func(mainTestFile *os.File) {
		err := mainTestFile.Close()
		if err != nil {
			panic(err)
		}
	}(mainTestFile)

	mainTestTemplate := template.Must(template.New("main_test").Parse(string(tpl.AddMainTestTemplate())))

	err = mainTestTemplate.Execute(mainTestFile, c)
	if err != nil {
		return err
	}

	return nil
}
