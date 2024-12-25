package migrations

import "github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/model"

func Models() []any {
	return []any{

		// (*models.Users)(nil),
		// (*model.Products)(nil),
		// (*model.Category)(nil),
		(*model.SystemBank)(nil),
	}
}

func RawBeforeQueryMigrate() []string {
	return []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
	}
}

func RawAfterQueryMigrate() []string {
	return []string{}
}
