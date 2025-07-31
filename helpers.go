package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func countUserMessages(entries []LogEntry) int {
	count := 0
	for _, entry := range entries {
		if entry.Type == "user" {
			count++
		}
	}
	return count
}

func countAssistantMessages(entries []LogEntry) int {
	count := 0
	for _, entry := range entries {
		if entry.Type == "assistant" {
			count++
		}
	}
	return count
}

func countToolUses(entries []LogEntry) int {
	count := 0
	for _, entry := range entries {
		if blocks, ok := entry.Message.Content.([]ContentBlock); ok {
			for _, block := range blocks {
				if _, ok := block.(*ToolUseBlock); ok {
					count++
				}
			}
		}
	}
	return count
}

func getToolCounts(entries []LogEntry) map[string]int {
	toolCounts := make(map[string]int)
	for _, entry := range entries {
		if blocks, ok := entry.Message.Content.([]ContentBlock); ok {
			for _, block := range blocks {
				if toolUse, ok := block.(*ToolUseBlock); ok {
					toolCounts[toolUse.Name]++
				}
			}
		}
	}
	return toolCounts
}

func countLines(text string) int {
	if text == "" {
		return 0
	}
	return len(strings.Split(strings.TrimSpace(text), "\n"))
}

func getFirstLines(text string, maxLines int) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) <= maxLines {
		return text
	}
	return strings.Join(lines[:maxLines], "\n")
}

func formatToolUseResult(result interface{}) string {
	if result == nil {
		return "null"
	}
	
	// Try to format as pretty JSON first
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		// Fall back to Go's %+v format if JSON marshaling fails
		return fmt.Sprintf("%+v", result)
	}
	
	return string(jsonBytes)
}

func extractSessionUUID(filename string) string {
	// Extract UUID from filename like "91cc2e2a-2d04-46ba-a5cf-5fcadf00f1da.jsonl"
	if len(filename) < 36 {
		return ""
	}
	
	// Remove path and extension
	basename := filename
	if lastSlash := strings.LastIndex(filename, "/"); lastSlash != -1 {
		basename = filename[lastSlash+1:]
	}
	if lastDot := strings.LastIndex(basename, "."); lastDot != -1 {
		basename = basename[:lastDot]
	}
	
	// Check if it looks like a UUID (36 characters with hyphens in right places)
	if len(basename) == 36 && 
		basename[8] == '-' && basename[13] == '-' && 
		basename[18] == '-' && basename[23] == '-' {
		return basename
	}
	
	return ""
}

func getSessionCwd(entries []LogEntry) string {
	// Get the cwd from the first entry (all should be the same for a session)
	if len(entries) > 0 {
		return entries[0].Cwd
	}
	return ""
}