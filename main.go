package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <jsonl-file> [output.html]")
	}

	filename := os.Args[1]
	outputFile := "output.html"
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}
	
	entries, err := parseJSONL(filename)
	if err != nil {
		log.Fatalf("Error parsing JSONL: %v", err)
	}

	fmt.Printf("Successfully parsed %d log entries\n", len(entries))
	
	// Print some stats
	userCount := countUserMessages(entries)
	assistantCount := countAssistantMessages(entries)
	toolCounts := getToolCounts(entries)
	
	fmt.Printf("User messages: %d\n", userCount)
	fmt.Printf("Assistant messages: %d\n", assistantCount)
	fmt.Println("\nTool usage:")
	for tool, count := range toolCounts {
		fmt.Printf("  - %s: %d\n", tool, count)
	}
	
	// Generate HTML output
	if err := generateHTML(entries, outputFile, filename); err != nil {
		log.Fatalf("Error generating HTML: %v", err)
	}
	
	fmt.Printf("\nHTML output written to: %s\n", outputFile)
}

func generateHTML(entries []LogEntry, outputFile, inputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	// Render the templ component to the file
	component := ConversationLog(entries, inputFile)
	return component.Render(context.Background(), file)
}

func parseJSONL(filename string) ([]LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	
	// Increase buffer size to handle large lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024) // 10MB max line size
	
	lineNum := 0
	
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		
		if line == "" {
			continue
		}
		
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			log.Printf("Error parsing JSON on line %d: %v", lineNum, err)
			continue
		}
		
		// Parse the message content based on its type
		if err := parseMessageContent(&entry.Message); err != nil {
			log.Printf("Error parsing message content on line %d: %v", lineNum, err)
			continue
		}
		
		entries = append(entries, entry)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	
	return entries, nil
}

func parseMessageContent(message *Message) error {
	// First, unmarshal the raw content to determine its type
	var rawContent interface{}
	contentBytes, err := json.Marshal(message.Content)
	if err != nil {
		return err
	}
	
	if err := json.Unmarshal(contentBytes, &rawContent); err != nil {
		return err
	}
	
	// Check if it's a string (simple user messages)
	if str, ok := rawContent.(string); ok {
		message.Content = str
		return nil
	}
	
	// Otherwise, it should be an array of content blocks
	if arr, ok := rawContent.([]interface{}); ok {
		var blocks []ContentBlock
		
		for _, item := range arr {
			blockBytes, err := json.Marshal(item)
			if err != nil {
				return err
			}
			
			// Determine block type
			var typeCheck struct {
				Type string `json:"type"`
			}
			if err := json.Unmarshal(blockBytes, &typeCheck); err != nil {
				return err
			}
			
			var block ContentBlock
			switch typeCheck.Type {
			case "text":
				var textBlock TextBlock
				if err := json.Unmarshal(blockBytes, &textBlock); err != nil {
					return err
				}
				block = &textBlock
				
			case "tool_use":
				var toolUseBlock ToolUseBlock
				if err := json.Unmarshal(blockBytes, &toolUseBlock); err != nil {
					return err
				}
				
				// Parse the tool input based on tool name
				if err := parseToolInput(&toolUseBlock); err != nil {
					return err
				}
				block = &toolUseBlock
				
			case "tool_result":
				var toolResultBlock ToolResultBlock
				if err := json.Unmarshal(blockBytes, &toolResultBlock); err != nil {
					return err
				}
				block = &toolResultBlock
				
			default:
				return fmt.Errorf("unknown content block type: %s", typeCheck.Type)
			}
			
			blocks = append(blocks, block)
		}
		
		message.Content = blocks
		return nil
	}
	
	return fmt.Errorf("unexpected content type: %v", reflect.TypeOf(rawContent))
}

func parseToolInput(toolUse *ToolUseBlock) error {
	inputBytes, err := json.Marshal(toolUse.Input)
	if err != nil {
		return err
	}
	
	switch toolUse.Name {
	case "TodoWrite":
		var input TodoWriteInput
		if err := json.Unmarshal(inputBytes, &input); err != nil {
			return err
		}
		toolUse.Input = input
		
	case "Bash":
		var input BashInput
		if err := json.Unmarshal(inputBytes, &input); err != nil {
			return err
		}
		toolUse.Input = input
		
	case "Edit":
		var input EditInput
		if err := json.Unmarshal(inputBytes, &input); err != nil {
			return err
		}
		toolUse.Input = input
		
	case "Read":
		var input ReadInput
		if err := json.Unmarshal(inputBytes, &input); err != nil {
			return err
		}
		toolUse.Input = input
		
	case "Glob":
		var input GlobInput
		if err := json.Unmarshal(inputBytes, &input); err != nil {
			return err
		}
		toolUse.Input = input
		
	case "Grep":
		var input GrepInput
		if err := json.Unmarshal(inputBytes, &input); err != nil {
			return err
		}
		toolUse.Input = input
		
	case "LS":
		var input LSInput
		if err := json.Unmarshal(inputBytes, &input); err != nil {
			return err
		}
		toolUse.Input = input
		
	default:
		// Keep as interface{} for unknown tools
		log.Printf("Warning: unknown tool type %s, keeping input as interface{}", toolUse.Name)
	}
	
	return nil
}