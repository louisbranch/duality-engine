package web

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	campaignv1 "github.com/louisbranch/duality-engine/api/gen/go/campaign/v1"
	"github.com/louisbranch/duality-engine/internal/web/templates"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// campaignsRequestTimeout caps the gRPC request time for campaigns.
	campaignsRequestTimeout = 2 * time.Second
	// campaignThemePromptLimit caps the number of characters shown in the table.
	campaignThemePromptLimit = 80
)

// Handler routes web requests for the UI.
type Handler struct {
	campaignClient campaignv1.CampaignServiceClient
}

// NewHandler builds the HTTP handler for the web server.
func NewHandler(campaignClient campaignv1.CampaignServiceClient) http.Handler {
	handler := &Handler{campaignClient: campaignClient}
	return handler.routes()
}

// routes wires the HTTP routes for the web handler.
func (h *Handler) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", templ.Handler(templates.Home()))
	mux.Handle("/campaigns", http.HandlerFunc(h.handleCampaigns))
	return mux
}

// handleCampaigns returns the first page of campaign rows for HTMX.
func (h *Handler) handleCampaigns(w http.ResponseWriter, r *http.Request) {
	if h.campaignClient == nil {
		h.renderCampaignTable(w, r, nil, "Campaign service unavailable.")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), campaignsRequestTimeout)
	defer cancel()

	response, err := h.campaignClient.ListCampaigns(ctx, &campaignv1.ListCampaignsRequest{})
	if err != nil {
		log.Printf("list campaigns: %v", err)
		h.renderCampaignTable(w, r, nil, "Campaigns unavailable.")
		return
	}

	campaigns := response.GetCampaigns()
	if len(campaigns) == 0 {
		h.renderCampaignTable(w, r, nil, "No campaigns yet.")
		return
	}

	rows := buildCampaignRows(campaigns)
	h.renderCampaignTable(w, r, rows, "")
}

// renderCampaignTable renders a campaign table with optional rows and message.
func (h *Handler) renderCampaignTable(w http.ResponseWriter, r *http.Request, rows []templates.CampaignRow, message string) {
	templ.Handler(templates.CampaignsTable(rows, message)).ServeHTTP(w, r)
}

// buildCampaignRows formats campaign rows for the table.
func buildCampaignRows(campaigns []*campaignv1.Campaign) []templates.CampaignRow {
	rows := make([]templates.CampaignRow, 0, len(campaigns))
	for _, campaign := range campaigns {
		if campaign == nil {
			continue
		}
		rows = append(rows, templates.CampaignRow{
			Name:             campaign.GetName(),
			GMMode:           formatGmMode(campaign.GetGmMode()),
			ParticipantCount: strconv.FormatInt(int64(campaign.GetParticipantCount()), 10),
			CharacterCount:   strconv.FormatInt(int64(campaign.GetCharacterCount()), 10),
			ThemePrompt:      truncateText(campaign.GetThemePrompt(), campaignThemePromptLimit),
			CreatedDate:      formatCreatedDate(campaign.GetCreatedAt()),
		})
	}
	return rows
}

// formatGmMode returns a display label for a GM mode enum.
func formatGmMode(mode campaignv1.GmMode) string {
	switch mode {
	case campaignv1.GmMode_HUMAN:
		return "Human"
	case campaignv1.GmMode_AI:
		return "AI"
	case campaignv1.GmMode_HYBRID:
		return "Hybrid"
	default:
		return "Unspecified"
	}
}

// formatCreatedDate returns a YYYY-MM-DD string for a timestamp.
func formatCreatedDate(createdAt *timestamppb.Timestamp) string {
	if createdAt == nil {
		return ""
	}
	return createdAt.AsTime().Format("2006-01-02")
}

// truncateText shortens text to a maximum length with an ellipsis.
func truncateText(text string, limit int) string {
	if limit <= 0 || text == "" {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	return string(runes[:limit]) + "..."
}
