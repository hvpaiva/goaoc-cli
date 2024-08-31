package pkg

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hvpaiva/goaoc-cli/tpl"
)

const baseURL = "https://adventofcode.com"

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

func (c *Challenge) GetInput(cookie string) error {
	if cookie == "" && viper.IsSet("cookie") {
		cookie = viper.GetString("cookie")
	}

	if cookie == "" {
		return errCookieRequired
	}

	url := fmt.Sprintf("%s/%d/day/%d/input", baseURL, c.Year, c.Day)

	body, err := getInputWithAOCCookie(url, cookie)
	if err != nil {
		return fmt.Errorf("error getting input for day %d, year %d: %w", c.Day, c.Year, err)
	}

	if strings.HasPrefix(string(body), "Puzzle inputs differ by user") {
		return PuzzleInputsDifferError{Day: c.Day, Year: c.Year}
	}

	inputFile, err := os.Create(fmt.Sprintf("%s/internal/%d/day%02d/input.txt", c.AbsolutePath, c.Year, c.Day))
	if err != nil {
		return err
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			panic(err)
		}
	}(inputFile)

	err = os.WriteFile(inputFile.Name(), body, os.FileMode(0o644))
	if err != nil {
		return CreatingFileError{Err: err}
	}

	fmt.Printf("Input written to %s\n", inputFile.Name())

	return nil
}

const timout = 5 * time.Second

func getInputWithAOCCookie(url, cookie string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: cookie,
	})

	resp, err := http.DefaultClient.Do(req) //nolint:bodyclose // Body is closed in the defer statement, but the linter is buggy
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Printf("warning: closing response body: %v\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if strings.HasPrefix(string(body), "Please don't repeatedly") {
		return nil, RateLimitError{URL: url}
	}

	return body, nil
}

type RateLimitError struct {
	URL string
}

func (e RateLimitError) Error() string {
	return fmt.Sprintf("rate limited for %s", e.URL)
}

type PuzzleInputsDifferError struct {
	Day, Year int
}

func (e PuzzleInputsDifferError) Error() string {
	return fmt.Sprintf("puzzle inputs differ by user for day %d, year %d", e.Day, e.Year)
}

type CreatingFileError struct {
	Err error
}

func (e CreatingFileError) Error() string {
	return fmt.Sprintf("creating file: %v", e.Err)
}

var errCookieRequired = errors.New("cookie is required")
