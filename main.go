package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var projectDir string
	var err error

	// Check for --dir flag
	var dirArg string
	if len(args) >= 2 && args[0] == "--dir" {
		dirArg = args[1]
		args = args[2:]
	}
	projectDir, err = resolveProjectDir(dirArg)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return extractLatest(projectDir)
	}

	switch args[0] {
	case "--list":
		return listSessions(projectDir)
	case "--all":
		return extractAll(projectDir)
	case "--help", "-h":
		printUsage()
		return nil
	default:
		// Check if arg is a number referring to a session index
		if n, err := strconv.Atoi(args[0]); err == nil && n >= 1 && n <= 999 {
			return extractByIndex(projectDir, n)
		}
		return extractSession(projectDir, args[0])
	}
}

func printUsage() {
	fmt.Println(`Usage: claude_hist [OPTIONS] [SESSION_ID | NUMBER]

Extract Claude Code's text responses from conversation history.

Options:
  (none)          Show responses from the latest session
  <1-999>         Show responses for session at that index (see --list)
  --dir <path>    Use the given project path instead of the current directory
  --list          List available sessions with numeric indices
  --all           Show responses from all sessions
  --help, -h      Show this help

Output goes to stdout — pipe to a file as needed:
  claude_hist > responses.md`)
}

// resolveProjectDir converts a project path to its Claude Code project directory.
// If path is empty, the current working directory is used.
func resolveProjectDir(path string) (string, error) {
	if path == "" {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("getting working directory: %w", err)
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolving path: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	// Claude Code uses the absolute path with "/" replaced by "-", minus the leading "-"
	key := strings.ReplaceAll(absPath, "/", "-")
	key = strings.TrimPrefix(key, "-")

	dir := filepath.Join(home, ".claude", "projects", "-"+key)
	if _, err := os.Stat(dir); err == nil {
		return dir, nil
	}

	// Fallback: also try with "_" converted to "-"
	altKey := strings.ReplaceAll(key, "_", "-")
	if altKey != key {
		altDir := filepath.Join(home, ".claude", "projects", "-"+altKey)
		if _, err := os.Stat(altDir); err == nil {
			return altDir, nil
		}
	}

	return "", fmt.Errorf("no Claude Code project found for %s, dir: %s", absPath, dir)
}

type sessionInfo struct {
	id        string
	path      string
	timestamp string
	modTime   int64
}

func getSessions(projectDir string) ([]sessionInfo, error) {
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return nil, fmt.Errorf("reading project directory: %w", err)
	}

	var sessions []sessionInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		s := sessionInfo{
			id:      strings.TrimSuffix(e.Name(), ".jsonl"),
			path:    filepath.Join(projectDir, e.Name()),
			modTime: info.ModTime().UnixNano(),
		}
		s.timestamp = firstTimestamp(s.path)
		sessions = append(sessions, s)
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].modTime > sessions[j].modTime
	})
	return sessions, nil
}

func firstTimestamp(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return "unknown"
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		var rec record
		if json.Unmarshal(scanner.Bytes(), &rec) == nil && rec.Timestamp != "" {
			return formatTimestamp(rec.Timestamp)
		}
	}
	return "unknown"
}

func listSessions(projectDir string) error {
	fmt.Printf("Sessions for %s:\n", projectDir)

	sessions, err := getSessions(projectDir)
	if err != nil {
		return err
	}
	for i, s := range sessions {
		fmt.Printf("  %2d. %s  (%s)\n", i+1, s.id, s.timestamp)
	}
	return nil
}

func extractByIndex(projectDir string, index int) error {
	sessions, err := getSessions(projectDir)
	if err != nil {
		return err
	}
	if index < 1 || index > len(sessions) {
		return fmt.Errorf("session index %d out of range (1-%d)", index, len(sessions))
	}
	s := sessions[index-1]
	fmt.Printf("# Session: %s\n\n", s.id)
	return extractResponses(s.path)
}

func extractLatest(projectDir string) error {
	sessions, err := getSessions(projectDir)
	if err != nil {
		return err
	}
	if len(sessions) == 0 {
		return fmt.Errorf("no conversations found")
	}
	s := sessions[0]
	fmt.Printf("# Session: %s\n\n", s.id)
	return extractResponses(s.path)
}

func extractAll(projectDir string) error {
	sessions, err := getSessions(projectDir)
	if err != nil {
		return err
	}
	for _, s := range sessions {
		fmt.Printf("# Session: %s\n\n", s.id)
		if err := extractResponses(s.path); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
		fmt.Println()
	}
	return nil
}

func extractSession(projectDir, sessionID string) error {
	path := filepath.Join(projectDir, sessionID+".jsonl")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Session '%s' not found.\nAvailable sessions:\n", sessionID)
		_ = listSessions(projectDir)
		return fmt.Errorf("session not found")
	}
	fmt.Printf("# Session: %s\n\n", sessionID)
	return extractResponses(path)
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type message struct {
	Content []contentBlock `json:"content"`
}

type record struct {
	Type      string  `json:"type"`
	Timestamp string  `json:"timestamp"`
	Message   message `json:"message"`
}

func extractResponses(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 2*1024*1024), 2*1024*1024)

	msgNum := 0
	for scanner.Scan() {
		var rec record
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			continue
		}
		if rec.Type != "assistant" {
			continue
		}

		var texts []string
		for _, block := range rec.Message.Content {
			if block.Type == "text" {
				t := strings.TrimSpace(block.Text)
				if t != "" {
					texts = append(texts, t)
				}
			}
		}
		if len(texts) == 0 {
			continue
		}

		msgNum++
		ts := formatTimestamp(rec.Timestamp)
		fmt.Println("---")
		fmt.Printf("### Response %d (%s)\n\n", msgNum, ts)
		for _, t := range texts {
			fmt.Println(t)
			fmt.Println()
		}
	}
	return scanner.Err()
}

func formatTimestamp(ts string) string {
	if len(ts) >= 19 {
		return strings.Replace(ts[:19], "T", " ", 1)
	}
	return ts
}
