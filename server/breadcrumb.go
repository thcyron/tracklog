package server

type BreadcrumbItem struct {
	Title  string
	Path   string
	Active bool
}

type Breadcrumb struct {
	items []BreadcrumbItem
}

func (b *Breadcrumb) Add(title, path string, active bool) {
	b.items = append(b.items, BreadcrumbItem{
		Title:  title,
		Path:   path,
		Active: active,
	})
}

func (b *Breadcrumb) Items() []BreadcrumbItem {
	return b.items
}
