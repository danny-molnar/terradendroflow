// main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Define flags for input and output file paths
	inputFilePath := flag.String("input", "", "Path to the input Terraform plan stdout file")
	outputFilePath := flag.String("output", "prettified_plan.md", "Path to the output markdown file")
	flag.Parse()

	if *inputFilePath == "" {
		fmt.Println("Error: Input file path is required")
		os.Exit(1)
	}

	err := parseAndPrettifyStdout(*inputFilePath, *outputFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Prettified plan saved to", *outputFilePath)
}

// parseAndPrettifyStdout parses the Terraform stdout plan and writes a prettified output to a markdown file
func parseAndPrettifyStdout(filePath, outputPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Failed to close input file: %v\n", err)
		}
	}()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			fmt.Printf("Failed to close output file: %v\n", err)
		}
	}()

	if _, err := outputFile.WriteString("# Prettified Terraform Plan\n\n"); err != nil {
		fmt.Printf("Failed to write header to output file: %v\n", err)
		return err
	}

	stats := map[string]int{
		"Created":  0,
		"Updated":  0,
		"Deleted":  0,
		"Replaced": 0,
		"Read":     0,
	}

	scanner := bufio.NewScanner(file)
	resources := map[string][]string{
		"Created":  {},
		"Updated":  {},
		"Deleted":  {},
		"Replaced": {},
		"Read":     {},
	}

	currentIdentifier := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check for full identifier line
		if strings.HasPrefix(line, "# module.") {
			currentIdentifier = strings.TrimPrefix(line, "# ")
			continue
		}

		// Skip lines that are likely attributes or JSON fragments
		if isAttributeLine(line) || isInvalidLine(line) {
			continue
		}

		if isResourceLine(line) {
			resource := extractResourceName(line)
			if currentIdentifier != "" {
				resource = currentIdentifier
				currentIdentifier = ""
			}
			action := determineAction(line)
			stats[action]++
			resources[action] = append(resources[action], resource)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Write grouped resources with full identifiers to the markdown file
	for action, resList := range resources {
		if len(resList) > 0 {
			if _, err := outputFile.WriteString(fmt.Sprintf("## %s Resources\n\n", action)); err != nil {
				fmt.Printf("Failed to write resource group to output file: %v\n", err)
				return err
			}
			for _, res := range resList {
				if _, err := outputFile.WriteString(fmt.Sprintf("%s\n", res)); err != nil {
					fmt.Printf("Failed to write resource to output file: %v\n", err)
					return err
				}
			}
			if _, err := outputFile.WriteString("\n"); err != nil {
				fmt.Printf("Failed to write newline to output file: %v\n", err)
				return err
			}
		}
	}

	// Print summary statistics to stdout
	fmt.Println("Summary of changes:")
	for action, count := range stats {
		fmt.Printf("- %s: %d\n", action, count)
	}

	return nil
}

// isAttributeLine checks if a line is likely an attribute, JSON fragment, or brace
func isAttributeLine(line string) bool {
	return line == "{" || line == "}" || strings.HasPrefix(line, "\"") || strings.HasSuffix(line, ")")
}

// isInvalidLine checks if a line contains irrelevant or invalid content
func isInvalidLine(line string) bool {
	return line == "[" || line == "]" || len(line) < 3
}

// isResourceLine checks if a line represents a resource action
func isResourceLine(line string) bool {
	return strings.HasPrefix(line, "+ resource") || strings.HasPrefix(line, "~ resource") ||
		strings.HasPrefix(line, "- resource") || strings.HasPrefix(line, "+/- resource") ||
		strings.HasPrefix(line, "+ module.") || strings.HasPrefix(line, "~ module.") ||
		strings.HasPrefix(line, "- module.") || strings.HasPrefix(line, "+ aws_") ||
		strings.HasPrefix(line, "~ aws_") || strings.HasPrefix(line, "- aws_") ||
		strings.HasPrefix(line, "<= data")
}

// extractResourceName extracts the full resource identifier from a resource line
func extractResourceName(line string) string {
	re := regexp.MustCompile(`"(.*?)"\s+"(.*?)"`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		return fmt.Sprintf("%s.%s", matches[1], matches[2])
	}

	parts := strings.Fields(line)
	if len(parts) > 2 {
		return fmt.Sprintf("%s.%s", parts[1], parts[2]) // Ensure full identifier with type and name
	}
	return "unknown resource"
}

// determineAction determines the action type based on the line prefix
func determineAction(line string) string {
	switch {
	case strings.HasPrefix(line, "+"):
		return "Created"
	case strings.HasPrefix(line, "~"):
		return "Updated"
	case strings.HasPrefix(line, "-"):
		return "Deleted"
	case strings.HasPrefix(line, "+/-"):
		return "Replaced"
	case strings.HasPrefix(line, "<="):
		return "Read"
	default:
		return "Unknown action"
	}
}
