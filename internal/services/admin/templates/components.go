package templates

// Breadcrumb represents a single breadcrumb navigation item.
// If URL is empty the item is rendered as plain text (current page).
type Breadcrumb struct {
	// Label is the display text for this breadcrumb.
	Label string
	// URL is the navigation target. Empty means current (final) page.
	URL string
}

// PageHeading provides data for the PageHeader component.
type PageHeading struct {
	// Breadcrumbs is the breadcrumb trail. Empty means no breadcrumbs.
	Breadcrumbs []Breadcrumb
	// Title is the page heading text.
	Title string
	// ActionURL is the optional action button target. Empty means no button.
	ActionURL string
	// ActionLabel is the display text for the action button.
	ActionLabel string
}
