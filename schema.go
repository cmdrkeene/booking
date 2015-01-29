package booking

import "database/sql"

// Schema is the aggregate of all table create statements, constraints, etc.
// Per-object constants, e.g. CalendarSchema live in their respective files
type Schema struct {
	DB *sql.DB `inject:""`
}

func (s *Schema) Load() error {
	queries := []string{
		CalendarSchema,
		GuestbookSchema,
		RegisterSchema,
	}
	for _, query := range queries {
		_, err := s.DB.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
