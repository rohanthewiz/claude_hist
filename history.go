package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rohanthewiz/serr"
)

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

type sessionInfo struct {
	id        string
	path      string
	timestamp string
	modTime   int64
}

type Response struct {
	Number    int
	Timestamp string
	Texts     []string
}

func resolveProjectDir(path string) (string, error) {
	if path == "" {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return "", serr.Wrap(err, "msg", "getting working directory")
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", serr.Wrap(err, "msg", "resolving path")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", serr.Wrap(err, "msg", "getting home directory")
	}

	key := strings.ReplaceAll(absPath, "/", "-")
	key = strings.TrimPrefix(key, "-")

	dir := filepath.Join(home, ".claude", "projects", "-"+key)
	if _, err := os.Stat(dir); err == nil {
		return dir, nil
	}

	altKey := strings.ReplaceAll(key, "_", "-")
	if altKey != key {
		altDir := filepath.Join(home, ".claude", "projects", "-"+altKey)
		if _, err := os.Stat(altDir); err == nil {
			return altDir, nil
		}
	}

	return "", serr.New("no Claude Code project found", "path", absPath, "dir", dir)
}

func getSessions(projectDir string) ([]sessionInfo, error) {
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return nil, serr.Wrap(err, "msg", "reading project directory", "dir", projectDir)
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

func getResponses(path string) ([]Response, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, serr.Wrap(err, "msg", "opening session file", "path", path)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 2*1024*1024), 2*1024*1024)

	var responses []Response
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
		responses = append(responses, Response{
			Number:    msgNum,
			Timestamp: formatTimestamp(rec.Timestamp),
			Texts:     texts,
		})
	}
	if err := scanner.Err(); err != nil {
		return responses, serr.Wrap(err, "msg", "scanning session file", "path", path)
	}
	return responses, nil
}

func formatTimestamp(ts string) string {
	if len(ts) >= 19 {
		return strings.Replace(ts[:19], "T", " ", 1)
	}
	return ts
}