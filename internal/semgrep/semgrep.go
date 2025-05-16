package semgrep

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Pos struct {
	Col    int `json:"col"`
	Line   int `json:"line"`
	Offset int `json:"offset"`
}

type Result struct {
	CheckID string `json:"check_id"`
	End     Pos    `json:"end"`
	Extra   Extra  `json:"extra"`
	Path    string `json:"path"`
	Start   Pos    `json:"start"`
}

type Extra struct {
	EngineKind  string         `json:"engine_kind"`
	Fingerprint string         `json:"fingerprint"`
	IsIgnored   bool           `json:"is_ignored"`
	Lines       string         `json:"lines"`
	Message     string         `json:"message"`
	Metadata    map[string]any `json:"metadata"`
	Severity    string         `json:"severity"`
}

type Paths struct {
	Comment string   `json:"_comment"`
	Scanned []string `json:"scanned"`
}

type OutputError struct {
	Code    int    `json:"code"`
	Level   string `json:"level"`
	Type    any    `json:"type"`
	Message string `json:"message"`
}

type Output struct {
	Errors  []OutputError `json:"errors"`
	Paths   Paths         `json:"paths"`
	Results []Result      `json:"results"`
	Version string        `json:"version"`
}

func ReadFile(name string) (*Output, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var output Output
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

// ExtendLines returns r with the Start and End extended to include
// the full line content
func ExtendLines(r Result, data []byte) Result {
	if len(data) == 0 {
		return r
	}
	isNL := func(b byte) bool { return b == '\n' || b == '\r' }
	if !isNL(data[r.Start.Offset]) {
		for r.Start.Offset > 0 && !isNL(data[r.Start.Offset-1]) {
			r.Start.Offset--
			r.Start.Col--
		}
	}
	if !isNL(data[r.End.Offset]) {
		for r.End.Offset < len(data) && !isNL(data[r.End.Offset]) {
			r.End.Offset++
			r.End.Col++
		}
	}
	return r
}

func FormatLines(data []byte, lineno, indent int) string {
	var b strings.Builder
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		num := strconv.Itoa(lineno)
		b.WriteString(strings.Repeat(" ", indent-len(num)))
		b.WriteString(num)
		b.WriteString("| ")
		b.Write(scanner.Bytes())
		b.WriteByte('\n')
		lineno++
	}
	if scanner.Err() != nil {
		panic("unreachable")
	}
	return b.String()
}

type ReadLinesOptions struct {
	Dir    string
	Extend bool
	Format bool
}

func ReadLines(r Result, opt ReadLinesOptions) (string, error) {
	filename := r.Path
	if opt.Dir != "" {
		filepath.Join(opt.Dir, r.Path)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if opt.Extend {
		r = ExtendLines(r, data)
	}
	match := data[r.Start.Offset:r.End.Offset]
	if opt.Format {
		return FormatLines(match, r.Start.Line, 5), nil
	}
	return string(match), nil
}

type ScanOptions struct {
	Dir        string
	ConfigPath string
	Bin        string
}

func Scan(opt ScanOptions) (*Output, error) {
	if opt.Bin == "" {
		opt.Bin = "semgrep"
	}
	cmd := exec.Command(opt.Bin, "scan",
		"--config", opt.ConfigPath,
		"--json",
		opt.Dir,
	)
	var output Output
	stdout, err := cmd.Output()
	if err != nil {
		if json.Unmarshal(stdout, &output) == nil && len(output.Errors) > 0 {
			return nil, fmt.Errorf("semgrep errors: %v", output.Errors)
		}
		return nil, err
	}
	if err := json.Unmarshal(stdout, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

type Rule struct {
	ID        string   `yaml:"id"`
	Message   string   `yaml:"message"`
	Languages []string `yaml:"languages"`
	Severity  string   `yaml:"severity"`
}

type Config struct {
	Name  string
	Rules []Rule
}

func ReadConfigs(dir string) ([]Config, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var configs []Config
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if ext := filepath.Ext(name); ext != ".yml" && ext != ".yaml" {
			continue
		}
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config %s: %w", name, err)
		}
		var aux struct {
			Rules []Rule `yaml:"rules"`
		}
		if err := yaml.Unmarshal(data, &aux); err != nil {
			return nil, fmt.Errorf("failed to read config %s: %w", name, err)
		}
		configs = append(configs, Config{
			Name:  name,
			Rules: aux.Rules,
		})
	}
	return configs, nil
}
