package products_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/anwardh/meliProject/internal/products"
	"github.com/stretchr/testify/assert"
)

func Test_sqlRepository_Store_Mock(t *testing.T) {
	t.Run("deve buscar um produto com o id informado", func(t *testing.T) {
		db, mock := SetupMock(t)
		defer db.Close()

		// Aqui é como se criamos uma tabela
		columns := []string{"id", "name", "type", "count", "price"}
		rows := sqlmock.NewRows(columns)
		productId := 1
		rows.AddRow(productId, "", "", 0, 0.0)

		// estamos retornando 1,1 por conta do retorno precisar do ultimo id
		// inserido e do número de linhas afetadas
		// mock.ExpectExec("INSERT INTO products").WithArgs(sql.Named("id", productId), sql.Named("name", "batata")).WillReturnResult(sqlmock.NewResult(1, 1))

		// esse rows funciona como se fosse o resultado, definimos as colunas desse resultado e os valores
		mock.ExpectQuery("SELECT .* FROM products").WithArgs(sql.Named("id", productId)).WillReturnRows(rows)

		// Realizando os testes

		// criamos o nosso repository com o db que foi criado no mock
		repository := products.NewMySqlRepository(db)
		// ctx := context.TODO()

		// Aqui é o produto que iremos inserir primeiro para depois buscá-lo
		// product := products.Product{
		// 	ID:   productId,
		// 	Name: "batata",
		// }

		// Verificamos se nao há produtos na base de dados com esse id
		// Fazemos as asserções de não ter erros nessa interação
		// e verificamos se o retorno é nil
		getResult, err := repository.GetOne(productId)
		assert.NoError(t, err)
		assert.Equal(t, products.Product{}, getResult)

		// Aqui estamos inserindo o produto uqe iremos buscar futuramente
		// _, err = repository.Store(product)
		// assert.NoError(t, err)
		// // Busca o produto que acabamos de inserir
		// getResult, err = repository.GetOne(productId)
		// assert.NoError(t, err)
		// assert.NotNil(t, getResult)
		// assert.Equal(t, product.ID, getResult.ID)
		// assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func SetupMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	// Configurando o sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	return db, mock
}
