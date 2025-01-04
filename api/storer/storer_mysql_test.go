package storer

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func withTestDB(t *testing.T, fn func(*sqlx.DB, sqlmock.Sqlmock)) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}

	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	fn(db, mock)
}

func TestCreateProduct(t *testing.T) {
	p := &Product{
		Name:         "IPhone 16",
		Image:        "iphone16.png",
		Category:     "Smartphone",
		Description:  "Brand new smartphone designed by Apple Inc",
		Rating:       5,
		NumReviews:   10,
		Price:        100.0,
		CountInStock: 100,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1, 1))
				cp, err := st.CreateProduct(context.Background(), p)
				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed inserting product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnError(fmt.Errorf("Error inserting product"))
				_, err := st.CreateProduct(context.Background(), p)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting last inserted ID",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("Error getting last inserted ID")))
				_, err := st.CreateProduct(context.Background(), p)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySqlStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestGetProduct(t *testing.T) {
	p := &Product{
		Name:         "IPhone 16",
		Image:        "iphone16.png",
		Category:     "Smartphone",
		Description:  "Brand new smartphone designed by Apple Inc",
		Rating:       5,
		NumReviews:   10,
		Price:        100.0,
		CountInStock: 100,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnRows(rows)

				gp, err := st.GetProduct(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, int64(1), gp.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("Error getting product"))
				_, err := st.GetProduct(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySqlStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestGetAllProducts(t *testing.T) {
	p := &Product{
		ID:           1,
		Name:         "IPhone 16",
		Image:        "iphone16.png",
		Category:     "Smartphone",
		Description:  "Brand new smartphone designed by Apple Inc",
		Rating:       5,
		NumReviews:   10,
		Price:        100.0,
		CountInStock: 100,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	} {
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)
				mock.ExpectQuery("SELECT * FROM products").WillReturnRows(rows)

				products, err := st.GetAllProducts(context.Background())
				require.NoError(t, err)
				require.Len(t, products, 1)
				err = mock.ExpectationsWereMet()

				require.NoError(t, err)
			},
		},
		{
			name: "failed querying products",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM products").WillReturnError(fmt.Errorf("error querying products"))
				_, err := st.GetAllProducts(context.Background())
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
			st := NewMySqlStorer(db)
			tc.test(t, st, mock)
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	p := &Product{
		ID:           1,
		Name:         "IPhone 16",
		Image:        "iphone16.png",
		Category:     "Smartphone",
		Description:  "Brand new smartphone designed by Apple Inc",
		Rating:       5,
		NumReviews:   10,
		Price:        100.0,
		CountInStock: 100,
	}

	new_p := &Product{
		ID:           1,
		Name:         "IPhone 15",
		Image:        "iphone15.png",
		Category:     "Smartphone",
		Description:  "Second latest smartphone designed by Apple Inc",
		Rating:       4,
		NumReviews:   8,
		Price:        90.0,
		CountInStock: 90,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1, 1))
				cp, err := st.CreateProduct(context.Background(), p)
				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)

				mock.ExpectExec("UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=? WHERE id=?").WillReturnResult(sqlmock.NewResult(1, 1))

				up, err := st.UpdateProduct(context.Background(), new_p)
				require.NoError(t, err)
				require.Equal(t, int64(1), up.ID)
				require.Equal(t, new_p.Name, up.Name)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed updating product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=? WHERE id=?").WillReturnError(fmt.Errorf("error updating product"))

				_, err := st.UpdateProduct(context.Background(), new_p)

				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySqlStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM products WHERE id=?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				err := st.DeleteProduct(context.Background(), 1)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM products WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("error deleting product"))

				err := st.DeleteProduct(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
			st := NewMySqlStorer(db)
			tc.test(t, st, mock)
		})
	}
}

func TestCreateOrder(t *testing.T) {
	ois := []OrderItem {
		{
			Name: "Tesla Car",
			Quantity: 1,
			Image: "tesla.png",
			Price: 99.99,
			ProductID: 1,
		},
		{
			Name: "IPhone 15",
			Quantity: 1,
			Image: "iphone15.png",
			Price: 55.99,
			ProductID: 1,
		},
	}

	o := &Order{
		PaymentMethod: "QRIS",
		TaxPrice: 10.0,
		ShippingPrice: 20.0,
		TotalPrice: 200.00,
		Items: ois,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	} {
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1,1))
				mock.ExpectExec("INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1,1))
				mock.ExpectExec("INSERT INTO order_items (name, quantity, image, price, product_id, order_id) VALUES (?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(2,1))
				mock.ExpectCommit()

				co, err := st.CreateOrder(context.Background(), o)
				require.NoError(t, err)
				require.Equal(t, int64(1), co.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed creating order",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price VALUES (?, ?, ?, ?)").WillReturnError(fmt.Errorf("error creating order"))
				mock.ExpectRollback()
				_, err := st.CreateOrder(context.Background(), o)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed creating order item",
			test: func (t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock)  {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price VALUES (?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1,1))
				mock.ExpectExec("INSERT INTO order_items (name, quantity, image, price, product_id, order_id VALUES (?, ?, ?, ?, ?, ?)").WillReturnError(fmt.Errorf("error creating order item"))
				mock.ExpectRollback()

				_, err := st.CreateOrder(context.Background(), o)
				require.Error(t, err)
				mock.ExpectationsWereMet()
				require.NoError(t, err)
				
			},
		},
	}

	for _, tc := range tcs {
		withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
			st := NewMySqlStorer(db)
			tc.test(t, st, mock)
		})
	}
}