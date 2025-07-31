package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
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

func extractProjectName(filename string) string {
	// Extract project name from path like "/Users/user/.claude/projects/my-project/session.jsonl"
	// Look for the pattern "/.claude/projects/{project-name}/"
	
	claudeIndex := strings.Index(filename, "/.claude/projects/")
	if claudeIndex == -1 {
		return ""
	}
	
	// Start after "/.claude/projects/"
	start := claudeIndex + len("/.claude/projects/")
	remaining := filename[start:]
	
	// Find the next slash to get the project name
	if slashIndex := strings.Index(remaining, "/"); slashIndex != -1 {
		return remaining[:slashIndex]
	}
	
	// If no slash found, the remaining part might be the project name
	// (though this shouldn't happen in normal usage)
	return remaining
}

func getSessionCwd(entries []LogEntry) string {
	// Get the cwd from the first entry (all should be the same for a session)
	if len(entries) > 0 {
		return entries[0].Cwd
	}
	return ""
}

func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 0 {
			return "just now"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	} else {
		return t.Format("Jan 2, 2006")
	}
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen]
}