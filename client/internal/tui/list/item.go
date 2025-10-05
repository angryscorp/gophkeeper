package list

type Item struct {
	title    string
	kind     string
	info     string
	showInfo bool
}

func (i Item) Title() string {
	return i.title
}

func (i Item) Description() string {
	if i.showInfo {
		return i.info
	}
	return i.kind
}

func (i Item) FilterValue() string {
	return i.title
}
