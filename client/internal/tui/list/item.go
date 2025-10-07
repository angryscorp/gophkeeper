package list

// item represents a record displayed in the list view.
type item struct {
	title    string
	kind     string
	info     string
	showInfo bool
}

// Title returns the item's title, shown as the main label in the list.
func (i item) Title() string {
	return i.title
}

// Description returns either the detailed info (if showInfo is true)
// or the item's kind as a short description.
func (i item) Description() string {
	if i.showInfo {
		return i.info
	}
	return i.kind
}

// FilterValue returns the value used for filtering and searching
// within the list, which is the item's title.
func (i item) FilterValue() string {
	return i.title
}
