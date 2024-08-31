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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hvpaiva/goaoc-cli/pkg"
)

var (
	day, year int
	cookie    string
)

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a challenge to the application",
	Aliases: []string{"create", "challenge"},
	Long: `Add (goaoc-cli add) will create a new challenge for a day and a year param.
	
This command requires the day and year flags to be set, and a session cookie to be set or be present in the config file.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		cobra.CheckErr(err)

		err = validateDayAndYear(day, year)
		cobra.CheckErr(err)

		modName := getModImportPath()

		challenge := &pkg.Challenge{
			Project: &pkg.Project{
				AbsolutePath: wd,
				PkgName:      modName,
				AppName:      path.Base(modName),
			},
			Day:  day,
			Year: year,
		}

		cobra.CheckErr(challenge.Create())
		cobra.CheckErr(challenge.GetInput(cookie))

		cobra.CheckErr(exec.Command("go", "mod", "tidy").Run())

		fmt.Printf("Challenge %d-%02d created successfully!\n", challenge.Year, challenge.Day)
	},
}

func init() {
	addCmd.Flags().IntVarP(&day, "day", "d", 0, "day of the challenge (required)")
	addCmd.Flags().IntVarP(&year, "year", "y", 0, "year of the challenge (required)")
	addCmd.MarkFlagsRequiredTogether("day", "year")

	addCmd.Flags().StringVarP(&cookie, "cookie", "c", "", "AOC session cookie")

	cobra.CheckErr(viper.BindPFlag("cookie", addCmd.Flags().Lookup("cookie")))
}

func validateDayAndYear(day, year int) error {
	if day < 1 || day > 25 {
		return errInvalidDay
	}

	if year < 2015 {
		return errInvalidYear
	}

	return nil
}

var errInvalidDay = errors.New("invalid day. Day must be between 1 and 25")

var errInvalidYear = errors.New("invalid year. Year must be greater than or equal to 2015")
