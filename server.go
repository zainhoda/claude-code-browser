package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type ProjectInfo struct {
	Name     string
	Path     string
	ModTime  time.Time
	Sessions []SessionInfo
}

type SessionInfo struct {
	UUID        string
	Filename    string
	ModTime     time.Time
	Size        int64
	LatestTodos *TodoWriteInput
}

func startServer(port string) {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/project/", projectHandler)
	http.HandleFunc("/session/", sessionHandler)
	
	fmt.Printf("Starting Claude Code Parser server on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := getClaudeProjects()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading projects: %v", err), http.StatusInternalServerError)
		return
	}
	
	component := ProjectsIndex(projects)
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
	}
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract project name from URL: /project/my-project-name
	projectName := strings.TrimPrefix(r.URL.Path, "/project/")
	if projectName == "" {
		http.Error(w, "Project name required", http.StatusBadRequest)
		return
	}
	
	claudeDir := os.ExpandEnv("$HOME/.claude/projects")
	projectPath := filepath.Join(claudeDir, projectName)
	
	sessions, err := getProjectSessions(projectPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading project sessions: %v", err), http.StatusInternalServerError)
		return
	}
	
	component := ProjectDetail(projectName, sessions)
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
	}
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	// Extract project and session from URL: /session/project-name/session-uuid
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/session/"), "/")
	if len(pathParts) != 2 {
		http.Error(w, "Invalid session URL format", http.StatusBadRequest)
		return
	}
	
	projectName, sessionUUID := pathParts[0], pathParts[1]
	
	claudeDir := os.ExpandEnv("$HOME/.claude/projects")
	sessionPath := filepath.Join(claudeDir, projectName, sessionUUID+".jsonl")
	
	// Check if file exists
	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	
	// Parse the JSONL file
	entries, err := parseJSONL(sessionPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing session: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Render the conversation log
	component := ConversationLog(entries, sessionPath)
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
	}
}

func getClaudeProjects() ([]ProjectInfo, error) {
	claudeDir := os.ExpandEnv("$HOME/.claude/projects")
	
	entries, err := os.ReadDir(claudeDir)
	if err != nil {
		return nil, fmt.Errorf("error reading Claude projects directory: %v", err)
	}
	
	var projects []ProjectInfo
	
	for _, entry := range entries {
		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			
			projectPath := filepath.Join(claudeDir, entry.Name())
			sessions, err := getProjectSessions(projectPath)
			if err != nil {
				// Skip projects we can't read
				continue
			}
			
			projects = append(projects, ProjectInfo{
				Name:     entry.Name(),
				Path:     projectPath,
				ModTime:  info.ModTime(),
				Sessions: sessions,
			})
		}
	}
	
	// Sort by modification time (most recent first)
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].ModTime.After(projects[j].ModTime)
	})
	
	return projects, nil
}

func getProjectSessions(projectPath string) ([]SessionInfo, error) {
	entries, err := os.ReadDir(projectPath)
	if err != nil {
		return nil, err
	}
	
	var sessions []SessionInfo
	
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jsonl") {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			
			// Extract UUID from filename
			uuid := strings.TrimSuffix(entry.Name(), ".jsonl")
			if len(uuid) == 36 && isValidUUID(uuid) {
				sessions = append(sessions, SessionInfo{
					UUID:     uuid,
					Filename: entry.Name(),
					ModTime:  info.ModTime(),
					Size:     info.Size(),
				})
			}
		}
	}
	
	// Sort by modification time (most recent first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].ModTime.After(sessions[j].ModTime)
	})
	
	// For the first 5 sessions (most recent), get the latest TodoWrite
	for i := 0; i < len(sessions) && i < 5; i++ {
		sessionPath := filepath.Join(projectPath, sessions[i].Filename)
		sessions[i].LatestTodos = getLatestTodoWrite(sessionPath)
	}
	
	return sessions, nil
}

func isValidUUID(uuid string) bool {
	return len(uuid) == 36 &&
		uuid[8] == '-' && uuid[13] == '-' &&
		uuid[18] == '-' && uuid[23] == '-'
}

func getLatestTodoWrite(sessionPath string) *TodoWriteInput {
	entries, err := parseJSONL(sessionPath)
	if err != nil {
		return nil
	}
	
	// Look through entries in reverse order (most recent first)
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		
		if blocks, ok := entry.Message.Content.([]ContentBlock); ok {
			for _, block := range blocks {
				if toolUse, ok := block.(*ToolUseBlock); ok {
					if toolUse.Name == "TodoWrite" {
						if todoInput, ok := toolUse.Input.(TodoWriteInput); ok {
							return &todoInput
						}
					}
				}
			}
		}
	}
	
	return nil
}