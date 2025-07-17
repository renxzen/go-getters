package test

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/renxzen/go-getters/pkg/generator"
	"github.com/renxzen/go-getters/pkg/parser"
)

var update = flag.Bool("update", false, "update golden files")

func TestGenerateGetters(t *testing.T) {
	tests := []struct {
		name       string
		structName string
		goldenFile string
	}{
		{
			name:       "basic_example",
			structName: "Example",
			goldenFile: "basic_struct.golden",
		},
		{
			name:       "multiple_structs",
			structName: "Container",
			goldenFile: "multiple_structs.golden",
		},
		{
			name:       "pointer_fields",
			structName: "Pointer",
			goldenFile: "pointer_fields.golden",
		},
		{
			name:       "slice_fields",
			structName: "Slices",
			goldenFile: "slice_fields.golden",
		},
		{
			name:       "dynamic_imports",
			structName: "DynamicImports",
			goldenFile: "dynamic_imports.golden",
		},
		{
			name:       "map_types",
			structName: "Maps",
			goldenFile: "map_types.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New()
			result, err := p.ParseDirectory("testdata")
			if err != nil {
				t.Fatalf("Failed to parse directory: %v", err)
			}

			// Generate getters from testdata directory
			gen := generator.New()
			structNames := []string{tt.structName}
			outBytes, err := gen.GenerateGetters(structNames, result)
			if err != nil {
				t.Fatalf("GenerateGetters failed: %v", err)
			}

			goldenPath := filepath.Join("testdata", tt.goldenFile)

			if *update {
				// Update golden file
				if err := os.WriteFile(goldenPath, outBytes, 0644); err != nil {
					t.Fatalf("Failed to update golden file: %v", err)
				}
				t.Logf("Updated golden file: %s", goldenPath)
				return
			}

			// Compare with golden file
			expected, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("Failed to read golden file %s: %v", goldenPath, err)
			}

			// TODO: implement diff viewing pkg
			if !bytes.Equal(outBytes, expected) {
				t.Errorf("Generated output doesn't match golden file %s", tt.goldenFile)
				t.Errorf("Expected:\n%s", string(expected))
				t.Errorf("Got:\n%s", string(outBytes))
			}
		})
	}
}
