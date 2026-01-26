//go:build conformance

package conformance

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	simpleTextResponse        = "This is a simple text response for testing."
	staticTextResourceContent = "This is the content of the static text resource."
	staticTextResourceName    = "test_static_text"
	staticTextResourceURI     = "test://static-text"
)

// Register adds conformance-only MCP fixtures (tools, prompts, resources).
func Register(mcpServer *mcp.Server) {
	if mcpServer == nil {
		return
	}

	mcp.AddTool(mcpServer, simpleTextTool(), simpleTextHandler())
	mcpServer.AddPrompt(simplePrompt(), simplePromptHandler())
	mcpServer.AddResource(staticTextResource(), staticTextResourceHandler())
}

// simpleTextTool defines the MCP conformance tool schema for simple text output.
func simpleTextTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "test_simple_text",
		Description: "Conformance tool that returns a simple text response.",
	}
}

// simpleTextHandler returns a fixed text payload for conformance validation.
// TODO: Provide project-aware completion/tool examples once conformance fixtures map to Duality features.
func simpleTextHandler() mcp.ToolHandlerFor[struct{}, any] {
	return func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: simpleTextResponse},
			},
		}, nil, nil
	}
}

func simplePrompt() *mcp.Prompt {
	return &mcp.Prompt{
		Name:        "test_simple_prompt",
		Description: "Conformance prompt that returns a simple text message.",
	}
}

func simplePromptHandler() mcp.PromptHandler {
	return func(_ context.Context, _ *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{
				{
					Role: "user",
					Content: &mcp.TextContent{
						Text: "This is a simple prompt for testing.",
					},
				},
			},
		}, nil
	}
}

// staticTextResource defines the MCP conformance resource schema for static text content.
func staticTextResource() *mcp.Resource {
	return &mcp.Resource{
		Name:        staticTextResourceName,
		Description: "Conformance resource that returns fixed text content.",
		MIMEType:    "text/plain",
		URI:         staticTextResourceURI,
	}
}

// staticTextResourceHandler returns fixed text content for conformance validation.
func staticTextResourceHandler() mcp.ResourceHandler {
	return func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      staticTextResourceURI,
					MIMEType: "text/plain",
					Text:     staticTextResourceContent,
				},
			},
		}, nil
	}
}
