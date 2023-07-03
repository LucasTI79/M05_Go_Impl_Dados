package products_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/DATA-DOG/go-txdb"
	"github.com/anwardh/meliProject/internal/products"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_sqlRepository_Store(t *testing.T) {
	db := InitDatabase(t)
	defer db.Close()
	repository := products.NewMySqlRepository(db)

	invalidProductId := 2
	product := products.Product{
		Name:     "batata",
		Category: "vegetais",
		Count:    20,
		Price:    3.99,
	}
	product, err := repository.Store(product)
	assert.NoError(t, err)

	// aqui estamos testando o caso de um id que não exista
	getResult, err := repository.GetOne(invalidProductId)
	assert.NoError(t, err)
	assert.Equal(t, products.Product{}, getResult)
	getResult, err = repository.GetOne(product.ID)
	assert.NoError(t, err)
	assert.NotNil(t, getResult)
	assert.Equal(t, product.Name, getResult.Name)
}

func InitDatabase(t *testing.T) *sql.DB {
	t.Helper()
	txdb.Register("txdb", "mysql", "root:root@/storage")
	db, err := sql.Open("txdb", uuid.New().String())
	assert.NoError(t, err)
	return db
}

func TestGetOneWithContext(t *testing.T) {
	db := InitDatabase(t)

	id := 9
	product := products.Product{
		Name: "teste",
	}
	myRepo := products.NewMySqlRepository(db)
	// cria um context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	productResult, err := myRepo.GetOneWithContext(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, product.Name, productResult.Name)
}
