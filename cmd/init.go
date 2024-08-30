/*
Copyright © 2024 Highlander Paiva

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hvpaiva/goaoc-cli/pkg"
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"initialize", "initialise"},
	Short:   "Initialize a GoAOC Application",
	Long: `Initialize (goaoc-cli init) will create a new application, with 
the appropriate structure for a GoAOC application.

GoAOC init must be run inside of a go module (please run "go mod init <MODNAME>" first)`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var comps []string
		var directive cobra.ShellCompDirective
		if len(args) == 0 {
			comps = cobra.AppendActiveHelp(comps, "Optionally specify the path of the go module to initialize")
			directive = cobra.ShellCompDirectiveDefault
		} else if len(args) == 1 {
			comps = cobra.AppendActiveHelp(comps, "This command does not take any more arguments (but may accept flags)")
			directive = cobra.ShellCompDirectiveNoFileComp
		} else {
			comps = cobra.AppendActiveHelp(comps, "ERROR: Too many arguments specified")
			directive = cobra.ShellCompDirectiveNoFileComp
		}

		return comps, directive
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, err := initializeProject(args)
		cobra.CheckErr(err)
		cobra.CheckErr(goGet("github.com/hvpaiva/goaoc"))
		fmt.Printf("Your GoAOC application is ready at\n%s\n", projectPath)
	},
}

func initializeProject(args []string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if len(args) > 0 {
		if args[0] != "." {
			wd = fmt.Sprintf("%s/%s", wd, args[0])
		}
	}

	modName := getModImportPath()

	project := &pkg.Project{
		AbsolutePath: wd,
		PkgName:      modName,
		AppName:      path.Base(modName),
	}

	if err := project.Create(); err != nil {
		return "", err
	}

	return project.AbsolutePath, nil
}

func getModImportPath() string {
	mod, cd := parseModInfo()

	return path.Join(mod.Path, fileToURL(strings.TrimPrefix(cd.Dir, mod.Dir)))
}

func fileToURL(in string) string {
	i := strings.Split(in, string(filepath.Separator))

	return path.Join(i...)
}

func parseModInfo() (Mod, CurDir) {
	var (
		mod Mod
		dir CurDir
	)

	m := modInfoJSON("-m")
	cobra.CheckErr(json.Unmarshal(m, &mod))

	if mod.Path == "command-line-arguments" {
		cobra.CheckErr("Please run `go mod init <MODNAME>` before `cobra-cli init`")
	}

	e := modInfoJSON("-e")
	cobra.CheckErr(json.Unmarshal(e, &dir))

	return mod, dir
}

type Mod struct {
	Path, Dir, GoMod string
}

type CurDir struct {
	Dir string
}

func goGet(mod string) error {
	return exec.Command("go", "get", mod).Run()
}

func modInfoJSON(args ...string) []byte {
	cmdArgs := append([]string{"list", "-json"}, args...)
	out, err := exec.Command("go", cmdArgs...).Output()
	cobra.CheckErr(err)

	return out
}
