package main

import (
	"database/sql"
	"fmt"
)

type product struct {
	ID       int    `json:"iD"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT id, name, quantity, price from products"

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	products := []product{}

	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT name, quantity, price from products where id=%v", p.ID)
	row := db.QueryRow(query)
	err := row.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *product) addProduct(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO products (id, name, quantity, price) VALUES(%v, '%v', '%v', '%v')", p.ID, p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(id)
	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE products SET name='%v', quantity=%v, price=%v where id=%v", p.Name, p.Quantity, p.Price, p.ID)
	_, err := db.Exec(query)
	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("DELETE FROM products WHERE id=%v", p.ID)
	_, err := db.Exec(query)
	return err
}
