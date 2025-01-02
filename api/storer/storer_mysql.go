package storer

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MySQLStorer struct {
	db *sqlx.DB
}

func NewMySqlStorer(db *sqlx.DB) *MySQLStorer {
	return &MySQLStorer{
		db: db,
	}
}

func (ms *MySQLStorer) CreateProduct(ctx context.Context, p *Product) (*Product, error) {
	res, err := ms.db.NamedExecContext(ctx, "INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (:name, :image, :category, :description, :rating, :num_reviews, :price, :count_in_stock)", p)

	if err != nil {
		return nil, fmt.Errorf("Error inserting product: %w", err)
	}
	
	id, err := res.LastInsertId()

	if err != nil {
		return nil, fmt.Errorf("Error getting last inserted id: %w", err)
	}

	p.ID = id

	return p, nil
}

func (ms *MySQLStorer) GetProduct(ctx context.Context, id int64) (*Product, error) {
	var p Product

	err := ms.db.GetContext(ctx, &p, "SELECT * FROM products WHERE id=?", id)
	if err != nil {
		return nil, fmt.Errorf("Error getting product: %w", err)
	}

	return &p, nil
}

func (ms *MySQLStorer) GetAllProducts(ctx context.Context) ([]*Product, error) {
	var products []*Product

	err := ms.db.SelectContext(ctx, &products, "SELECT * FROM products")

	if err != nil {
		return nil, fmt.Errorf("Error getting list of all products: %w", err)
	}

	return products, nil
}

func (ms *MySQLStorer) UpdateProduct(ctx context.Context, p *Product) (*Product, error) {
	_, err := ms.db.ExecContext(ctx, "UPDATE products SET name=:name, image=:image, category=:category, description=:description, rating=:rating, num_reviews=:num_reviews, price=:price, count_in_stock=:count_in_stock WHERE id=:id", p)

	if err != nil {
		return nil, fmt.Errorf("Error updating product: %w", err)
	}

	return p, nil
}

func (ms *MySQLStorer) DeleteProduct(ctx context.Context, id int64) error {
	_, err := ms.db.ExecContext(ctx, "DELETE FROM products WHERE id=?", id)

	if err != nil {
		return fmt.Errorf("Error deleting product: %w", err)
	}

	return nil
}