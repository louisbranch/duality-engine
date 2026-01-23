package service

import (
	"github.com/louisbranch/duality-engine/api/gen/go/duality/v1"
	"github.com/louisbranch/duality-engine/internal/mcp/domain"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerDualityTools(mcpServer *mcp.Server, client dualityv1.DualityServiceClient) {
	mcp.AddTool(mcpServer, domain.ActionRollTool(), domain.ActionRollHandler(client))
	mcp.AddTool(mcpServer, domain.DualityOutcomeTool(), domain.DualityOutcomeHandler(client))
	mcp.AddTool(mcpServer, domain.DualityExplainTool(), domain.DualityExplainHandler(client))
	mcp.AddTool(mcpServer, domain.DualityProbabilityTool(), domain.DualityProbabilityHandler(client))
	mcp.AddTool(mcpServer, domain.RulesVersionTool(), domain.RulesVersionHandler(client))
	mcp.AddTool(mcpServer, domain.RollDiceTool(), domain.RollDiceHandler(client))
}
