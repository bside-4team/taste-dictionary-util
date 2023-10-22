package insertexternalrestaurant

import (
	"fmt"
	"testing"
	"time"

	"github.com/Taehoya/go-utils/pq"
	"github.com/stretchr/testify/assert"
)

func TestInsertExternalCusine(t *testing.T) {
	db, err := pq.InitTestDB()
	fmt.Println(db)
	assert.NoError(t, err)
	defer db.Close()

	repository := NewRepository(db)

	t.Run("Successfully save user", func(t *testing.T) {
		defer pq.SetUp(db, "./teardown_test.sql")
		pq.SetUp(db, "./teardown_test.sql")

		externalUUID := "123"
		x := "123.1235"
		y := "123.1236"
		placeUrl := "google.com"
		updatedAt := time.Now()
		name := "test=name"

		err = repository.InsertExternalCusine(externalUUID, x, y, placeUrl, updatedAt, name)
		assert.NoError(t, err)
	})
}
