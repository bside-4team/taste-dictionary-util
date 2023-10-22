package insertexternalrestaurant

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) InsertExternalCusine(externalUUID, x, y, placeUrl string, updatedAt time.Time, placeName string) error {
	stmt := `
	INSERT INTO
		public.external_restaurant_informations
		(external_uuid, "location", reference_link, updated_at, name)
	VALUES
		($1, ST_GeomFromText($2), $3, $4, $5)
`

	_, err := r.db.Exec(stmt, externalUUID, fmt.Sprintf("POINT(%s %s)", x, y), placeUrl, updatedAt, placeName)
	if err != nil {
		return err
	}

	return nil
}
