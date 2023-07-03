package products

import "context"

// Adicionando a Estrutura Product e seus campos rotulados
type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Count    int     `json:"count"`
	Price    float64 `json:"price"`
}

type ProductFullDataResponse struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Category         string  `json:"category" db:"type"`
	Count            int     `json:"count"`
	Price            float64 `json:"price"`
	Warehouse        string  `json:"warehouse"`
	WarehouseAddress string  `json:"warehouse_address"`
}

/*
Criação da variável para guardar os produtos

	Corresponde a nossa Camada de Persistência de Dados
*/

// Criação da Interface e Declaração dos Métodos
type Repository interface {
	GetOne(id int) (Product, error)
	GetOneWithContext(ctx context.Context, id int) (Product, error)
	GetAll() ([]Product, error)
	Store(product Product) (Product, error)
	// LastID() (int, error)
	// Declaração do Método Update - que cuidará de atualizar um dado
	Update(product Product) (Product, error)

	// Declaração do Método UpdateName
	UpdateName(id int, name string) (Product, error)

	// Declaração do Método Delete
	Delete(id int) error
}
