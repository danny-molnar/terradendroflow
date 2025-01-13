// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	filePath := "tfplan.stdout"      // Replace with your actual stdout file path
	outputPath := "prettified_plan.md" // Output markdown file

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

	stats := map[string]int{
		"Created": 0,
		"Updated": 0,
		"Deleted": 0,
		"Replaced": 0,
		"Read": 0,
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
			outputFile.WriteString(fmt.Sprintf("## %s Resources\n\n", action))
			for _, res := range resList {
				outputFile.WriteString(fmt.Sprintf("%s\n", res))
			}
			outputFile.WriteString("\n")
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
