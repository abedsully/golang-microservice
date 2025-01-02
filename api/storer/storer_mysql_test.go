package storer

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	st := NewMySqlStorer(db)

	p := &Product{
		Name: "IPhone 16",
		Image: "iphone16.png",
		Category: "Smartphone",
		Description: "Brand new smartphone designed by Apple Inc",
		Rating: 5,
		NumReviews: 10,
		Price: 100.0,
		CountInStock: 100,
	}

	mock.ExpectExec("INSERT INTO products(name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1, 1))

	cp, err := st.CreateProduct(context.Background(), p)

	require.NoError(t, err)
	require.Equal(t, int64(1), cp.ID)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}