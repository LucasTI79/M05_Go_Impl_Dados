package products

import (
	"context"
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
	GetAllFullDataProducts = "SELECT products.id, products.name, products.count, products.type, products.price, warehouses.name, warehouses.address " +
		"FROM products INNER JOIN warehouses ON products.id_warehouse = warehouses.id " +
		"WHERE products.id = ?"
	GetAllProducts    = "SELECT id, name, type, count, price, FROM products"
	ProductStore      = "INSERT INTO products(name, type, count, price) VALUES(?,?,?,?)"
	GetOneProduct     = "SELECT p.id, p.name, p.type, p.count, p.price FROM products p WHERE id = ?"
	UpdateProduct     = "UPDATE products SET name = ?, type = ?, count = ?, price = ? WHERE id = ?"
	UpdateProductName = "UPDATE products SET name = ? WHERE id = ?"
	DeleteProduct     = "DELETE FROM products WHERE id = ?"
)

// storage -> database

// products
// warehouses

func (r *mysqlRepository) GetFullData(id int) ProductFullDataResponse {
	var productFullData ProductFullDataResponse
	rows, err := r.db.Query(GetAllFullDataProducts, id)
	if err != nil {
		log.Println(err)
		return productFullData
	}
	for rows.Next() {
		if err := rows.Scan(&productFullData.ID, &productFullData.Name, &productFullData.Count, &productFullData.Category, &productFullData.Price, &productFullData.Warehouse,
			&productFullData.WarehouseAddress); err != nil {
			log.Fatal(err)
			return productFullData
		}
	}
	return productFullData
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

// JBDC -> ORM de java
// Gorm -> ORM de Go lang

func (r *mysqlRepository) GetOne(id int) (Product, error) {
	var product Product

	rows, err := r.db.Query(GetOneProduct, id)

	if err != nil {
		log.Println(err)
		return product, err
	}

	// 1 "bolo de cenoura" "doces" 1 25.00

	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Count, &product.Price)
		if err != nil {
			log.Println(err.Error())
			return product, err
		}
	}

	return product, err
}

func (r *mysqlRepository) GetOneWithContext(ctx context.Context, id int) (Product, error) {

	// Caso quisermos simular uma demora na execução de instrução do banco, usamos o time.Sleep()
	// time.Sleep(time.Second * 2)
	// Caso quisermos simular o timeout com a consulta no abnco, também podemos utilizar essa instrução
	// getQuery := "SELECT SLEEP(30) FROM DUAL where 0 < ?"

	var product Product

	// db.Query não é mais usado, mas db.QueryContext
	rows, err := r.db.QueryContext(ctx, GetOneProduct, id)

	if err != nil {
		log.Println(err)
		return product, err
	}

	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Count, &product.Price); err != nil {
			log.Fatal(err)
			return product, err
		}
	}
	return product, nil
}

// main -> routes -> controller <-> service <-> repository <-> db

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

func (r *mysqlRepository) Delete(id int) error {
	// DELETE FROM products WHERE id = ?
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

// Prepare
// Exec
// Usamos quando queremos alterar/deletar/inserir algum dado

// Query
// Scan
// Usamos quando queremos visualizar os dados
