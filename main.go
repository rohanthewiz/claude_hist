package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rohanthewiz/logger"
	"github.com/rohanthewiz/rweb"
	"github.com/rohanthewiz/serr"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var dirArg string
	if len(args) >= 2 && args[0] == "--dir" {
		dirArg = args[1]
		args = args[2:]
	}

	projectDir, err := resolveProjectDir(dirArg)
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
	case "--web":
		addr := ":7070"
		if len(args) >= 2 {
			addr = args[1]
		}
		displayDir := dirArg
		if displayDir == "" {
			displayDir, _ = os.Getwd()
		}
		return serveWeb(projectDir, displayDir, addr)
	case "--help", "-h":
		printUsage()
		return nil
	default:
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
  (none)              Show responses from the latest session
  <1-999>             Show responses for session at that index (see --list)
  --dir <path>        Use the given project path instead of the current directory
  --list              List available sessions with numeric indices
  --all               Show responses from all sessions
  --web [addr]        Start web UI (default addr: :7070)
  --help, -h          Show this help

Output goes to stdout — pipe to a file as needed:
  claude_hist > responses.md`)
}

func serveWeb(projectDir, displayDir, addr string) error {
	logger.InitLog(logger.LogConfig{
		Formatter: "text",
		LogLevel:  "info",
	})
	defer logger.CloseLog()

	logger.Info("Starting Claude History server", "addr", addr, "project_dir", projectDir)

	s := rweb.NewServer(rweb.ServerOptions{
		Address: addr,
		Verbose: true,
	})
	s.Use(rweb.RequestInfo)

	s.Get("/", func(ctx rweb.Context) error {
		sessions, err := getSessions(projectDir)
		if err != nil {
			logger.Err(serr.Wrap(err, "msg", "listing sessions"))
			return ctx.SetStatus(500).WriteHTML(renderErrorPage("Failed to load sessions: " + err.Error()))
		}
		return ctx.WriteHTML(renderSessionsPage(sessions, displayDir))
	})

	s.Get("/session/:id", func(ctx rweb.Context) error {
		id := ctx.Request().PathParam("id")

		sessions, err := getSessions(projectDir)
		if err != nil {
			logger.Err(serr.Wrap(err, "msg", "loading sessions for detail", "id", id))
			return ctx.SetStatus(500).WriteHTML(renderErrorPage("Failed to load sessions"))
		}

		var found *sessionInfo
		for i := range sessions {
			if sessions[i].id == id {
				found = &sessions[i]
				break
			}
		}
		if found == nil {
			return ctx.SetStatus(404).WriteHTML(renderErrorPage("Session not found: " + id))
		}

		responses, err := getResponses(found.path)
		if err != nil {
			logger.Err(serr.Wrap(err, "msg", "loading responses", "id", id))
			return ctx.SetStatus(500).WriteHTML(renderErrorPage("Failed to load session"))
		}

		return ctx.WriteHTML(renderSessionPage(id, found.timestamp, responses))
	})

	log.Fatal(s.Run())
	return nil
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