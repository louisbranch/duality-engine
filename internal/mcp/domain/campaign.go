package domain

import (
	"context"
	"fmt"
	"strings"
	"time"

	campaignpb "github.com/louisbranch/duality-engine/api/gen/go/campaign/v1"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// CampaignCreateInput represents the MCP tool input for campaign creation.
type CampaignCreateInput struct {
	Name        string `json:"name" jsonschema:"campaign name"`
	GmMode      string `json:"gm_mode" jsonschema:"gm mode (HUMAN, AI, HYBRID)"`
	PlayerSlots int    `json:"player_slots" jsonschema:"number of player slots"`
	ThemePrompt string `json:"theme_prompt" jsonschema:"optional theme prompt"`
}

// CampaignCreateResult represents the MCP tool output for campaign creation.
type CampaignCreateResult struct {
	ID          string `json:"id" jsonschema:"campaign identifier"`
	Name        string `json:"name" jsonschema:"campaign name"`
	GmMode      string `json:"gm_mode" jsonschema:"gm mode"`
	PlayerSlots int    `json:"player_slots" jsonschema:"number of player slots"`
	ThemePrompt string `json:"theme_prompt" jsonschema:"theme prompt"`
}

// CampaignCreateTool defines the MCP tool schema for creating campaigns.
func CampaignCreateTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "campaign_create",
		Description: "Creates a new campaign metadata record",
	}
}

// CampaignCreateHandler executes a campaign creation request.
func CampaignCreateHandler(client campaignpb.CampaignServiceClient) mcp.ToolHandlerFor[CampaignCreateInput, CampaignCreateResult] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input CampaignCreateInput) (*mcp.CallToolResult, CampaignCreateResult, error) {
		runCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		response, err := client.CreateCampaign(runCtx, &campaignpb.CreateCampaignRequest{
			Name:        input.Name,
			GmMode:      gmModeFromString(input.GmMode),
			PlayerSlots: int32(input.PlayerSlots),
			ThemePrompt: input.ThemePrompt,
		})
		if err != nil {
			return nil, CampaignCreateResult{}, fmt.Errorf("campaign create failed: %w", err)
		}
		if response == nil || response.Campaign == nil {
			return nil, CampaignCreateResult{}, fmt.Errorf("campaign create response is missing")
		}

		result := CampaignCreateResult{
			ID:          response.Campaign.GetId(),
			Name:        response.Campaign.GetName(),
			GmMode:      gmModeToString(response.Campaign.GetGmMode()),
			PlayerSlots: int(response.Campaign.GetPlayerSlots()),
			ThemePrompt: response.Campaign.GetThemePrompt(),
		}

		return nil, result, nil
	}
}

func gmModeFromString(value string) campaignpb.GmMode {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "HUMAN":
		return campaignpb.GmMode_HUMAN
	case "AI":
		return campaignpb.GmMode_AI
	case "HYBRID":
		return campaignpb.GmMode_HYBRID
	default:
		return campaignpb.GmMode_GM_MODE_UNSPECIFIED
	}
}

func gmModeToString(mode campaignpb.GmMode) string {
	switch mode {
	case campaignpb.GmMode_HUMAN:
		return "HUMAN"
	case campaignpb.GmMode_AI:
		return "AI"
	case campaignpb.GmMode_HYBRID:
		return "HYBRID"
	default:
		return "UNSPECIFIED"
	}
}
