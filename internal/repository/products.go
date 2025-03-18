package repository

import (
	"context"
	"database/sql"
	"fmt"
	"berezovskiyvalerii/csv-rest-app/internal/domain"
	database "berezovskiyvalerii/csv-rest-app/pkg/db"
)

type Products struct {
	db *sql.DB
}

func NewProducts(db *sql.DB) *Products {
	return &Products{
		db: db,
	}
}

func (p *Products) Create(ctx context.Context, product domain.Product) (int64, error) {
	var id int64
	query := fmt.Sprintf("insert into %s(name, price) values($1, $2) RETURNING id", database.ProductsTable)
	row := p.db.QueryRowContext(ctx, query, product.Name, product.Price)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (p *Products) GetAll(ctx context.Context) ([]domain.Product, error){
	query := fmt.Sprintf("SELECT id, name, price FROM %s", database.ProductsTable)
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *Products) Update(ctx context.Context, id int64, product domain.ProductUpdate) (int64, error) {
	query := fmt.Sprintf("UPDATE %s SET %s = $1, %s = $2 WHERE id = $3",
		database.ProductsTable, "name", "price")
	result, err := p.db.Exec(query, &product.Name, &product.Price, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("record not found")
	}

	return id, nil
}

func (p *Products) Delete(ctx context.Context, id int64) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", database.ProductsTable)
	result, err := p.db.Exec(query, id)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error getting rows affected: %v\n", err)
		return 0, err
	}

	if rowsAffected == 0 {
		fmt.Printf("No rows affected. Record with id=%d not found.\n", id)
		return 0, fmt.Errorf("record not found")
	}

	return id, nil
}
