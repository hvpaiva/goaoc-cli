package pkg

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/hvpaiva/goaoc-cli/tpl"
)

type Project struct {
	PkgName      string
	AppName      string
	AbsolutePath string
}

func (p *Project) Create() error {
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		if err = os.Mkdir(p.AbsolutePath, 0o754); err != nil {
			return err
		}
	}

	if _, err := os.Stat(fmt.Sprintf("%s/internal", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/internal", p.AbsolutePath), 0o751))
	}

	if _, err := os.Stat(fmt.Sprintf("%s/pkg", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/pkg", p.AbsolutePath), 0o751))
	}

	if _, err := os.Stat(fmt.Sprintf("%s/pkg/parser", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/pkg/parser", p.AbsolutePath), 0o751))
	}

	parserFile, err := os.Create(fmt.Sprintf("%s/pkg/parser/parser.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer func(parserFile *os.File) {
		err := parserFile.Close()
		if err != nil {
			panic(err)
		}
	}(parserFile)

	parserTemplate := template.Must(template.New("parser").Parse(string(tpl.AddMainTestTemplate())))

	err = parserTemplate.Execute(parserFile, p)
	if err != nil {
		return err
	}

	return nil
}
