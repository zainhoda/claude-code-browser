# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Claude Code Parser is a Go application that parses JSONL files containing Claude Code session logs and provides both CLI and web-based visualization. The application uses strongly-typed Go structs to model conversation data without relying on `interface{}` types, providing type safety throughout the parsing pipeline.

**Dual Mode Operation:**
- **CLI Mode**: Process individual JSONL files and generate standalone HTML reports
- **Web Server Mode**: Browse all Claude Code projects and sessions through a web interface

## Core Architecture

### Data Flow Pipeline
1. **JSONL Parsing** (`main.go:parseJSONL`) - Reads large JSONL files with 10MB buffer support
2. **Struct Mapping** (`main.go:parseMessageContent`, `parseToolInput`) - Converts raw JSON to typed structs based on message/tool types  
3. **HTML Generation** (`templates.templ`) - Uses templ library to render structured HTML output

### Key Components

**structs.go** - Contains the complete data model hierarchy:
- `LogEntry` - Root structure with session metadata (timestamps, UUIDs, git info)
- `Message` - User/assistant message container with role and content
- `ContentBlock` interface - Polymorphic content types (text, tool_use, tool_result)
- Tool-specific input structs - Strongly typed inputs for all 7 tool types (TodoWrite, Bash, Edit, Read, Glob, Grep, LS)

**Template System** - Uses Go 1.24's local tool installation pattern:
- `templates.templ` - templ template definitions with embedded CSS and tool-specific formatters
- `templates_templ.go` - Generated Go code (do not edit manually)
- `helpers.go` - Statistical functions for template data

**Tool-Specific Templates** - Each tool type has custom HTML formatting:
- `TodoWrite` - Visual todo list with status/priority badges and color coding
- `Edit` - Side-by-side old/new text comparison with syntax highlighting
- `Bash` - Terminal-style command display with descriptions
- `Read/Glob/Grep/LS` - Structured parameter display with relevant options

### Critical Design Decisions

**Type Safety**: The parser intelligently routes JSON data to specific structs based on:
- `message.role` for user vs assistant messages
- `content[].type` for content block types (text/tool_use/tool_result) 
- `tool_use.name` for tool-specific input parsing

**Content Polymorphism**: `MessageContentType` can be either:
- `string` - Simple user messages (7 instances in typical data)
- `[]ContentBlock` - Complex messages with multiple blocks (240 instances)

## Development Commands

### Setup and Build
```bash
make setup          # Initialize dependencies and install templ tool
make generate       # Generate Go code from .templ files  
make build          # Build standalone binary
```

### CLI Mode (Single File Processing)
```bash
make run            # Quick run with default JSONL file
make run-file FILE=data.jsonl OUTPUT=report.html  # Custom input/output
```

### Web Server Mode (Browse All Sessions)
```bash
make server         # Start web server on port 8080
make server-port PORT=3000  # Start on custom port
```

### Development Workflow
```bash
make dev            # Complete setup -> generate -> run pipeline
make check          # Full validation (fmt + vet + test + build)
```

### Template Development
Always run `make generate` or `go tool templ generate` after editing `.templ` files. The generated `templates_templ.go` file must not be edited manually.

## Tool Architecture

The application models 7 different tool types with dedicated input structs:
- **TodoWrite**: Manages todo lists with status tracking
- **Edit**: File modifications with optional replace_all flag
- **Bash**: Command execution with descriptions
- **Read**: File reading with optional offset/limit
- **Glob**: Pattern matching with optional path filtering
- **Grep**: Text search with multiple output modes and filters
- **LS**: Directory listing

Tool inputs are parsed dynamically based on the `tool_use.name` field, providing compile-time type safety while supporting runtime polymorphism.

## Web Server Architecture

**URL Structure:**
- `/` - Project index (lists all projects from `~/.claude/projects`)
- `/project/{name}` - Project detail (shows all JSONL sessions)
- `/session/{project}/{uuid}` - Session viewer (full conversation display)

**Key Components:**
- `server.go` - HTTP handlers and routing logic
- `web_templates.templ` - Web-specific templates for project/session browsing
- `ProjectInfo/SessionInfo` structs - Metadata for directory listings

## Data Characteristics

Typical JSONL files contain:
- ~60% assistant messages, ~40% user messages
- ~35% of messages involve tool usage
- Large line sizes requiring 10MB scanner buffer
- Mixed `toolUseResult` types (string or object)
- Files stored in `~/.claude/projects/{project-name}/{session-uuid}.jsonl`