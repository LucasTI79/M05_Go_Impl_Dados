package products

import (
	"database/sql"
	"log"
)

type mysqlRepository struct {
	db *sql.DB
}

func NewMySqlRepository(dbConn *sql.DB) Repository {
	return &mysqlRepository{
		db: dbConn,
	}
}

const (
	GetAllProducts    = "SELECT id, name, type, count, price FROM products"
	ProductStore      = "INSERT INTO products(name, type, count, price) VALUES( ?, ?, ?, ? )"
	GetOneProduct     = "SELECT * FROM products WHERE id = ?"
	UpdateProduct     = "UPDATE products SET name = ?, type = ?, count = ?, price = ? WHERE id = ?"
	UpdateProductName = "UPDATE products SET name = ? WHERE id = ?"
	DeleteProduct     = "DELETE FROM products WHERE id = ?"
)

func (r *mysqlRepository) Store(product Product) (Product, error) {
	// o banco é iniciado
	stmt, err := r.db.Prepare(ProductStore) // monta o  SQL
	if err != nil {
		log.Fatal(err)
	}
	// o defer vai ser a última coisa a ser executada na função Store
	defer stmt.Close() // a instrução fecha quando termina. Se eles permanecerem abertos, o consumo de memória é gerado

	var result sql.Result
	result, err = stmt.Exec(product.Name, product.Category, product.Count, product.Price) // retorna um sql.Result ou um error
	if err != nil {
		return Product{}, err
	}
	insertedId, _ := result.LastInsertId() // do sql.Result retornado na execução obtemos o Id inserido
	product.ID = int(insertedId)

	return product, nil
}

// JBDC -> ORM de java
// Gorm -> ORM de Go lang

func (r *mysqlRepository) GetOne(id int) Product {
	var product Product

	rows, err := r.db.Query(GetOneProduct, id)

	if err != nil {
		log.Println(err)
		return product
	}

	// 1 "bolo de cenoura" "doces" 1 25.00

	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Count, &product.Price)
		if err != nil {
			log.Println(err.Error())
			return product
		}
	}
	return product
}

// main -> routes -> controller <-> service <-> repository <-> db

func (r *mysqlRepository) Update(product Product) (Product, error) {
	stmt, err := r.db.Prepare(UpdateProduct)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Category, product.Count, product.Price, product.ID)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (r *mysqlRepository) Delete(id int) error {
	stmt, err := r.db.Prepare(DeleteProduct)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlRepository) GetAll() ([]Product, error) {
	var products []Product

	rows, err := r.db.Query(GetAllProducts)

	if err != nil {
		log.Println(err)
		return products, err
	}

	for rows.Next() {
		// id, name, type, count, price
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Count, &product.Price)
		if err != nil {
			return products, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *mysqlRepository) UpdateName(id int, name string) (Product, error) {
	var product Product

	stmt, err := r.db.Prepare(UpdateProductName)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.ID)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}
