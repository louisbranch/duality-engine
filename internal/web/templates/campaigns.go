// File campaigns.go defines view data for campaign templates.
package templates

// CampaignRow holds formatted campaign data for display.
type CampaignRow struct {
	// Name is the display name of the campaign.
	Name string
	// GMMode is the display label for the GM mode.
	GMMode string
	// ParticipantCount is the formatted number of participants.
	ParticipantCount string
	// CharacterCount is the formatted number of characters.
	CharacterCount string
	// ThemePrompt is the truncated theme prompt text.
	ThemePrompt string
	// CreatedDate is the formatted creation date.
	CreatedDate string
}
