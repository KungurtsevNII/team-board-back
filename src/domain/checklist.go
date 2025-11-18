package domain

type Checklist struct {
	Title string
	Items []ChecklistItem
}

type ChecklistItem struct{
	Title string
	Completed bool
}

func NewChecklist(title string, ChecklistItem []ChecklistItem) Checklist {
	return Checklist{
		Title: title,
		Items: ChecklistItem,
	}
}

func NewChecklistItem(title string, completed bool) ChecklistItem {
	return ChecklistItem{
		Title: title,
		Completed: completed,
	}
}