package tpl

func AddMainTemplate() []byte {
	return []byte(`package main
	
import (
	_ "embed"
	"log"

	"{{ .PkgName }}/pkg/parser"

	"github.com/hvpaiva/goaoc"
)
	
//go:embed input.txt
var input string
	
func main() {
	err := goaoc.Run(parser.Normalize(input), partOne, partTwo)
	if err != nil {
		log.Fatalf("Error running AOC: %v", err)
	}
}
	
func partOne(input string) int {
	return len(input)
}
	
func partTwo(input string) int {
	return len(input) * 2
}
`)
}

func AddMainTestTemplate() []byte {
	return []byte(`package main
	
import (
	"testing"
)
	
func TestPartOne(t *testing.T) {
	tests := []ChallengeTest{
		{
			Name:     "Case 1",
			Input:    "0",
			Expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := partOne(tt.Input)
			if got != tt.Expected {
				t.Errorf("partOne() = %d; expected %d", got, tt.Expected)
			}
		})
	}
}

func TestPartTwo(t *testing.T) {
	tests := []ChallengeTest{
		{
			Name:     "Case 1",
			Input:    "0",
			Expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := partTwo(tt.Input)
			if got != tt.Expected {
				t.Errorf("partTwo() = %d; expected %d", got, tt.Expected)
			}
		})
	}
}

type ChallengeTest struct {
	Name     string
	Input    string
	Expected int
}
`)
}

func AddParserTemplate() []byte {
	return []byte(`package parser

import (
	"log"
	"strings"
)

func Normalize(input string) string {
	input = strings.TrimRight(input, "\n")
	if len(input) == 0 {
		log.Fatalf("input.txt is empty")
	}

	return input
}
`)
}
