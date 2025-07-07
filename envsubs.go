package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Updated regex: matches ${var} or ${var:default}
var pattern = regexp.MustCompile(`\$\{([^:{}]+)(?::([^{}]*))?\}`)

func convertKey(varName string) string {
	converted := strings.ToUpper(varName)
	converted = strings.ReplaceAll(converted, ".", "_")
	for i := 0; i < len(converted); i++ {
		c := converted[i]
		if !(c >= 'A' && c <= 'Z') && !(c >= '0' && c <= '9') && c != '_' {
			converted = strings.ReplaceAll(converted, string(c), "_")
		}
	}
	return converted
}

func resolveVariable(prefix, varName, defaultValue string) string {
	fmt.Fprintf(os.Stderr, "[DEBUG] prefix=%s, varName=%s, defaultValue=%s, result=??\n", prefix, varName, defaultValue)
	result := resolveVariable1(prefix, varName, defaultValue)
	fmt.Fprintf(os.Stderr, "[DEBUG] prefix=%s, varName=%s, defaultValue=%s, result=%s\n", prefix, varName, defaultValue, result)
	return result
}

func resolveVariable1(prefix, varName, defaultValue string) string {
	converted := convertKey(varName)
	fullKey := prefix + converted
	fmt.Fprintf(os.Stderr, "[DEBUG] Resolving: original='%s', converted='%s'\n", varName, converted)

	// Try direct env
	if val := os.Getenv(fullKey); val != "" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Found direct env var: %s=%s\n", fullKey, val)
		return val
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] Not found direct env var: %s\n", fullKey)

	// Try HEX fallback
	hexKey := strings.ToUpper(hex.EncodeToString([]byte(varName)))
	hexEnvKey := prefix + "HEX_" + hexKey
	if val := os.Getenv(hexEnvKey); val != "" {
		decoded, err := hex.DecodeString(val)
		if err == nil {
			fmt.Fprintf(os.Stderr, "[DEBUG] Found HEX env: %s (decoded: %s)\n", hexEnvKey, string(decoded))
			return string(decoded)
		} else {
			fmt.Fprintf(os.Stderr, "[DEBUG] Failed to base64 decode: %s: %v\n", hexEnvKey, err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "[DEBUG] Not found HEX env: %s\n", hexEnvKey)
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] Using default: %s\n", defaultValue)
	return defaultValue
}


func processLine(line, prefix string) string {
	return pattern.ReplaceAllStringFunc(line, func(m string) string {
		matches := pattern.FindStringSubmatch(m)
		if len(matches) < 2 {
			return m
		}
		varName := matches[1]
		defaultValue := ""
		if len(matches) == 3 && matches[2] != "" {
			defaultValue = matches[2]
		}
		return resolveVariable(prefix, varName, defaultValue)
	})
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input_file> <output_file> <env_prefix>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	envPrefix := os.Args[3]

	in, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)

	for scanner.Scan() {
		line := scanner.Text()
		processed := processLine(line, envPrefix)
		_, err := writer.WriteString(processed + "\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to output: %v\n", err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	writer.Flush()
}

