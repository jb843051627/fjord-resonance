package store

type Page struct {
	Limit  int
	Offset int
}

func NewPage(limit, offset int) Page {
	return Page{Limit: NormalizeLimit(limit, 50, 200), Offset: maxPage(offset, 0)}
}

func (p Page) Next() Page { return Page{Limit: p.Limit, Offset: p.Offset + p.Limit} }

func (p Page) Previous() Page {
	offset := p.Offset - p.Limit
	if offset < 0 {
		offset = 0
	}
	return Page{Limit: p.Limit, Offset: offset}
}

func (p Page) Valid() bool { return p.Limit > 0 && p.Limit <= 200 && p.Offset >= 0 }

func maxPage(a, b int) int {
	if a > b {
		return a
	}
	return b
}
