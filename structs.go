package main

import "time"

// Root message structure
type LogEntry struct {
	IsSidechain  bool        `json:"isSidechain"`
	UserType     string      `json:"userType"`
	ParentUuid   string      `json:"parentUuid"`
	RequestId    *string     `json:"requestId,omitempty"`
	ToolUseResult interface{} `json:"toolUseResult,omitempty"`
	SessionId    string      `json:"sessionId"`
	Version      string      `json:"version"`
	GitBranch    string      `json:"gitBranch"`
	Timestamp    time.Time   `json:"timestamp"`
	Message      Message     `json:"message"`
	Uuid         string      `json:"uuid"`
	Cwd          string      `json:"cwd"`
	Type         string      `json:"type"`
}

// Message can be either user or assistant message
type Message struct {
	Role         string           `json:"role"`
	Content      MessageContentType   `json:"content"`
	
	// Assistant-only fields
	Id           *string          `json:"id,omitempty"`
	Model        *string          `json:"model,omitempty"`
	StopReason   *string          `json:"stop_reason,omitempty"`
	StopSequence *string          `json:"stop_sequence,omitempty"`
	Type         *string          `json:"type,omitempty"`
	Usage        *Usage           `json:"usage,omitempty"`
}

// Usage information for assistant messages
type Usage struct {
	InputTokens               int    `json:"input_tokens"`
	CacheCreationInputTokens  int    `json:"cache_creation_input_tokens"`
	CacheReadInputTokens      int    `json:"cache_read_input_tokens"`
	OutputTokens              int    `json:"output_tokens"`
	ServiceTier               string `json:"service_tier"`
}

// MessageContentType can be either a string (for simple user messages) or an array of content blocks
type MessageContentType interface{}

// ContentBlock represents different types of content blocks
type ContentBlock interface{}

// TextBlock represents a text content block
type TextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ToolUseBlock represents a tool use request
type ToolUseBlock struct {
	Type  string    `json:"type"`
	Id    string    `json:"id"`
	Name  string    `json:"name"`
	Input ToolInput `json:"input"`
}

// ToolInput represents the input for different tools
type ToolInput interface{}

// TodoWriteInput represents input for TodoWrite tool
type TodoWriteInput struct {
	Todos []TodoItem `json:"todos"`
}

// TodoItem represents a single todo item
type TodoItem struct {
	Id       string `json:"id"`
	Content  string `json:"content"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
}

// BashInput represents input for Bash tool
type BashInput struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// EditInput represents input for Edit tool
type EditInput struct {
	FilePath   string `json:"file_path"`
	OldString  string `json:"old_string"`
	NewString  string `json:"new_string"`
	ReplaceAll *bool  `json:"replace_all,omitempty"`
}

// ReadInput represents input for Read tool
type ReadInput struct {
	FilePath string `json:"file_path"`
	Limit    *int   `json:"limit,omitempty"`
	Offset   *int   `json:"offset,omitempty"`
}

// GlobInput represents input for Glob tool
type GlobInput struct {
	Pattern string  `json:"pattern"`
	Path    *string `json:"path,omitempty"`
}

// GrepInput represents input for Grep tool
type GrepInput struct {
	Pattern    string  `json:"pattern"`
	Glob       *string `json:"glob,omitempty"`
	Path       *string `json:"path,omitempty"`
	OutputMode *string `json:"output_mode,omitempty"`
	LineNumbers *bool  `json:"-n,omitempty"`
}

// LSInput represents input for LS tool
type LSInput struct {
	Path string `json:"path"`
}

// ToolResultBlock represents a tool execution result
type ToolResultBlock struct {
	Type       string  `json:"type"`
	ToolUseId  string  `json:"tool_use_id"`
	Content    string  `json:"content"`
	IsError    *bool   `json:"is_error,omitempty"`
}