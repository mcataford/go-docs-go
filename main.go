package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
)

type ArgConfiguration struct {
	argName     string
	shortFlag   string
	longFlag    string
	expectValue bool
}

type Args struct {
	source    []string
	outputDir string
}

var rIsFlag = regexp.MustCompile(`^(-[a-zA-Z])|(--[a-zA-Z]+)$`)

var argumentConfiguration = []ArgConfiguration{
	ArgConfiguration{"outputDir", "-o", "--out", true},
}

// Parses command-line arguments into an easy to use Args struct.
// This also contains validation for the provided arguments and
// can terminate the program if invalid arguments are provided.
func parseArgs(args []string) Args {
	parsedArgs := map[string]string{}

	unmatched := []string{}

	expectValue := false
	expectedArg := ""
	expectNoFlag := false

	for _, argument := range args {
		isFlag := rIsFlag.MatchString(argument)

		if isFlag && expectNoFlag {
			log.Fatal("Unexpected flag")
		}

		if isFlag && expectValue {
			panic(fmt.Sprintf("Unexpected argument %s", argument))
		}

		if expectValue {
			parsedArgs[expectedArg] = argument
			expectValue = false
			expectedArg = ""
			continue
		}

		matched := false
		for _, argConfig := range argumentConfiguration {
			if argument == argConfig.shortFlag || argument == argConfig.longFlag {
				expectValue = argConfig.expectValue
				expectedArg = argConfig.argName
				matched = true
				break
			}
		}

		if !isFlag && !matched {
			expectNoFlag = true
			unmatched = append(unmatched, argument)
		}

		if matched && !expectValue {
			parsedArgs[expectedArg] = "true"
		}
	}

	return Args{unmatched, parsedArgs["outputDir"]}
}

func main() {
	args := parseArgs(os.Args[1:])

	log.Println(fmt.Sprintf("Processing %d files...", len(args.source)))

	outputDirectory := "docs"

	if args.outputDir != "" {
		outputDirectory = args.outputDir
	}

	err := os.Mkdir(outputDirectory, 0750)

	if err != nil {
		log.Fatal(err)
	}

	for _, sourceFile := range args.source {
		log.Println(fmt.Sprintf("Processing %s", sourceFile))
		readBytes, _ := ioutil.ReadFile(sourceFile)
		fileContent := string(readBytes)

		source := Parse(fileContent, sourceFile)

		generated_markup := GenerateMarkdown(source)

		baseName := path.Base(sourceFile)
		targetPath := fmt.Sprintf("%s/%s.md", outputDirectory, baseName)
		file, err := os.Create(targetPath)
		defer file.Close()

		if err != nil {
			return
		}

		file.WriteString(generated_markup)
		log.Println(fmt.Sprintf("Wrote %s", targetPath))
	}
}
