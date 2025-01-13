// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	filePath := "tfplan.stdout"      
	outputPath := "prettified_plan.md"

	err := parseAndPrettifyStdout(filePath, outputPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Prettified plan saved to", outputPath)
}

// parseAndPrettifyStdout parses the Terraform stdout plan and writes a prettified output to a markdown file
func parseAndPrettifyStdout(filePath, outputPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	outputFile.WriteString("# Prettified Terraform Plan\n\n")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip lines that are likely attributes or JSON fragments
		if isAttributeLine(line) || isInvalidLine(line) {
			continue
		}

		if isResourceLine(line) {
			resource := extractResourceName(line)
			action := determineAction(line)
			output := fmt.Sprintf("- **%s**: `%s`\n", action, resource)
			fmt.Print(output)
			outputFile.WriteString(output)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// isAttributeLine checks if a line is likely an attribute, JSON fragment, or brace
func isAttributeLine(line string) bool {
	// Skip lines that are single braces, empty braces, or JSON-like fragments
	return line == "{" || line == "}" || strings.HasPrefix(line, "\"") || strings.HasSuffix(line, ")")
}

// isInvalidLine checks if a line contains irrelevant or invalid content
func isInvalidLine(line string) bool {
	// Skip lines that are irrelevant, such as single characters or invalid fragments
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

// extractResourceName extracts the resource identifier from a plan line
func extractResourceName(line string) string {
	parts := strings.Fields(line)

	// Handle resource lines starting with +, ~, -, or +/-
	if len(parts) > 2 && (strings.HasPrefix(parts[0], "+") || strings.HasPrefix(parts[0], "~") || strings.HasPrefix(parts[0], "-") || strings.HasPrefix(parts[0], "+/-")) {
		return parts[1] // The second part typically contains the full resource identifier
	}

	// Handle lines starting with <= data
	if len(parts) > 2 && strings.HasPrefix(parts[0], "<=") {
		return parts[1] // The second part contains the data source identifier
	}

	// Fallback for unknown lines
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
