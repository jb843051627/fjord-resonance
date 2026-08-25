package store

type Sort struct {
	Field      string
	Descending bool
}

func (s Sort) Valid(allowed ...string) bool {
	for _, field := range allowed {
		if s.Field == field {
			return true
		}
	}
	return false
}

func (s Sort) SQL(defaultField string) string {
	field := s.Field
	if field == "" {
		field = defaultField
	}
	direction := "ASC"
	if s.Descending {
		direction = "DESC"
	}
	return field + " " + direction
}

func (s Sort) Reverse() Sort { s.Descending = !s.Descending; return s }
