package products_test

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/anwardh/meliProject/internal/products"
	"github.com/stretchr/testify/assert"
)

func Test_sqlRepository_GetOne_Mock(t *testing.T) {
	t.Run("deve buscar um produto com o id informado", func(t *testing.T) {
		db, mock := SetupMock(t)
		defer db.Close()

		// Aqui é como se criamos uma tabela
		columns := []string{"id", "name", "type", "count", "price"}
		rows := sqlmock.NewRows(columns)
		productId := 1
		rows.AddRow(productId, "", "", 0, 0.0)

		// esse rows funciona como se fosse o resultado, definimos as colunas desse resultado e os valores
		mock.ExpectQuery(products.GetOneProduct).WithArgs(productId).WillReturnRows(rows)

		// criamos o nosso repository com o db que foi criado no mock
		repository := products.NewMySqlRepository(db)

		// Verificamos se nao há produtos na base de dados com esse id
		// Fazemos as asserções de não ter erros nessa interação
		// e verificamos se o retorno é nil
		getResult, err := repository.GetOne(productId)
		assert.NoError(t, err)
		assert.Equal(t, productId, getResult.ID)
	})

	t.Run("deve buscar um produto inexistente", func(t *testing.T) {
		db, mock := SetupMock(t)
		defer db.Close()

		expectedError := errors.New("o produto com o id informado não existe")

		// Aqui é como se criamos uma tabela
		columns := []string{"id", "name", "type", "count", "price"}
		rows := sqlmock.NewRows(columns)
		productId := 1
		rows.AddRow(productId, "", "", 0, 0.0)

		// esse rows funciona como se fosse o resultado, definimos as colunas desse resultado e os valores
		mock.ExpectQuery(products.GetOneProduct).WithArgs(productId).WillReturnError(expectedError)

		// criamos o nosso repository com o db que foi criado no mock
		repository := products.NewMySqlRepository(db)

		// Verificamos se nao há produtos na base de dados com esse id
		// Fazemos as asserções de não ter erros nessa interação
		// e verificamos se o retorno é nil
		_, err := repository.GetOne(productId)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), expectedError.Error())
	})
}

func Test_sqlRepository_Store_Mock(t *testing.T) {
	mockProduct := products.Product{
		Name:     "batata",
		Category: "vegetais",
		Count:    20,
		Price:    3.99,
	}

	t.Run("deve criar um produto ao chamar o repository.Store", func(t *testing.T) {
		db, mock := SetupMock(t)
		defer db.Close()

		mock.ExpectPrepare(regexp.QuoteMeta(products.ProductStore))

		// estamos retornando 1,1 por conta do retorno precisar do ultimo id
		// inserido e do número de linhas afetadas
		mock.ExpectExec(regexp.QuoteMeta(products.ProductStore)).WithArgs(
			mockProduct.Name,
			mockProduct.Category,
			mockProduct.Count,
			mockProduct.Price,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		// Realizando os testes

		// criamos o nosso repository com o db que foi criado no mock
		repository := products.NewMySqlRepository(db)

		// Aqui estamos inserindo o produto uqe iremos buscar futuramente
		_, err := repository.Store(mockProduct)
		assert.NoError(t, err)
	})
}

func SetupMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	// Configurando o sql mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	return db, mock
}
