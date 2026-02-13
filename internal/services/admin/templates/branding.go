package templates

import "github.com/louisbranch/fracturing.space/internal/platform/branding"

// AppName returns the canonical product name.
func AppName() string {
	return branding.AppName
}
