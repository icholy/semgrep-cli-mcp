package semgrep

import (
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func TestReadConfigs(t *testing.T) {
	configs, err := ReadConfigs(filepath.FromSlash("testdata/configs"))
	assert.NilError(t, err)
	assert.DeepEqual(t, configs, []Config{
		{
			Name: "error.no.cause.yml",
			Rules: []Rule{
				{
					ID:        "err.no.cause",
					Message:   "The caught error should be passed as the CDLError cause",
					Languages: []string{"javascript"},
					Severity:  "WARNING",
				},
			},
		},
		{
			Name: "nullish.coalescing.yaml",
			Rules: []Rule{
				{
					ID:        "nulling.coalescing.0",
					Message:   "Use nullish coalescing",
					Languages: []string{"typescript", "javascript"},
					Severity:  "ERROR",
				},
			},
		},
	})
}

func TestScan(t *testing.T) {
	output, err := Scan(ScanOptions{
		Dir:        filepath.FromSlash("testdata/code"),
		ConfigPath: filepath.FromSlash("testdata/configs/nullish.coalescing.yaml"),
	})
	assert.NilError(t, err)
	assert.DeepEqual(t, output, &Output{
		Paths: Paths{
			Comment: "",
			Scanned: []string{"testdata/code/file.js"},
		},
		Results: []Result{
			{
				CheckID: "testdata.configs.nulling.coalescing.0",
				End:     Pos{Col: 2, Line: 6, Offset: 48},
				Extra: Extra{
					EngineKind:  "OSS",
					Fingerprint: "requires login",
					Lines:       "requires login",
					Message:     "Use nullish coalescing",
					Severity:    "ERROR",
					Metadata:    map[string]any{},
				},
				Path:  "testdata/code/file.js",
				Start: Pos{Col: 1, Line: 4, Offset: 16},
			},
		},
		Errors:  []OutputError{},
		Version: "1.122.0",
	})
}
