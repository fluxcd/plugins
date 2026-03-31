/*
Copyright 2026 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Command validate checks plugin YAML manifests against their JSON schemas.
//
// It resolves the schema from apiVersion/kind and validates
// using JSON Schema (draft 2020-12).
//
// Usage:
//
//	./bin/validate schemas/ catalog.yaml plugins/*.yaml
//	kustomize build . | ./bin/validate schemas/ /dev/stdin
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"sigs.k8s.io/yaml"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <schemas-dir> <files...>\n", os.Args[0])
		os.Exit(2)
	}

	schemasDir := os.Args[1]
	files := os.Args[2:]
	compiler := jsonschema.NewCompiler()

	failed := 0
	for _, file := range files {
		errs := validateFile(compiler, schemasDir, file)
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "FAIL %s: %v\n", file, err)
			failed++
		}
		if len(errs) == 0 {
			fmt.Printf("OK   %s\n", file)
		}
	}

	if failed > 0 {
		fmt.Fprintf(os.Stderr, "\n%d error(s) found\n", failed)
		os.Exit(1)
	}
}

// schemaFileName returns the schema filename for a given apiVersion and kind,
// following the kubeconform/kubeval naming convention: {kind}-{group}-{version}.json.
// For example, cli.fluxcd.io/v1beta1 Plugin -> plugin-cli-v1beta1.json.
func schemaFileName(apiVersion, kind string) string {
	group, version, _ := strings.Cut(apiVersion, "/")
	group, _, _ = strings.Cut(group, ".")
	return strings.ToLower(kind) + "-" + group + "-" + version + ".json"
}

// validateFile reads a YAML file, splits it into individual documents
// (separated by ---), and validates each against the matching JSON schema.
func validateFile(compiler *jsonschema.Compiler, schemasDir, path string) []error {
	data, err := os.ReadFile(path)
	if err != nil {
		return []error{err}
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 4*1024), 5*1024*1024)
	scanner.Split(splitYAMLDocument)

	var errs []error
	docIndex := 0
	for scanner.Scan() {
		raw := bytes.TrimSpace(scanner.Bytes())
		if len(raw) == 0 {
			continue
		}
		docIndex++
		if err := validateDocument(compiler, schemasDir, raw); err != nil {
			errs = append(errs, fmt.Errorf("document %d: %w", docIndex, err))
		}
	}
	if err := scanner.Err(); err != nil {
		errs = append(errs, err)
	}
	return errs
}

// validateDocument validates a single YAML document against its schema.
func validateDocument(compiler *jsonschema.Compiler, schemasDir string, raw []byte) error {
	var doc any
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		return fmt.Errorf("YAML parse error: %w", err)
	}

	m, ok := doc.(map[string]any)
	if !ok {
		return fmt.Errorf("document is not a YAML mapping")
	}

	apiVersion, _ := m["apiVersion"].(string)
	kind, _ := m["kind"].(string)
	if apiVersion == "" || kind == "" {
		return fmt.Errorf("missing 'apiVersion' or 'kind' field")
	}

	schemaFile := schemaFileName(apiVersion, kind)
	schemaPath := filepath.Join(schemasDir, schemaFile)
	if _, err := os.Stat(schemaPath); err != nil {
		return fmt.Errorf("no schema found for %s %s (expected %s)", apiVersion, kind, schemaFile)
	}

	sch, err := compiler.Compile(schemaPath)
	if err != nil {
		return fmt.Errorf("schema %s: %w", schemaFile, err)
	}

	if err := sch.Validate(doc); err != nil {
		return fmt.Errorf("schema validation failed:\n%v", err)
	}

	return nil
}

// splitYAMLDocument is a bufio.SplitFunc for splitting multi-doc YAML on "\n---" markers.
func splitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	sep := []byte("\n---")
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, sep); i >= 0 {
		// Found a potential separator — consume past the next newline.
		after := data[i+len(sep):]
		if len(after) == 0 {
			if atEOF {
				return len(data), data[:i], nil
			}
			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + len(sep) + j + 1, data[:i], nil
		}
		return 0, nil, nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
