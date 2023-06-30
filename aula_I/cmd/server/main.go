package main

import (
	"log"
	"net/http"
	"os"

	"github.com/anwardh/meliProject/cmd/server/handler"
	"github.com/anwardh/meliProject/config"
	"github.com/anwardh/meliProject/docs"
	"github.com/anwardh/meliProject/internal/products"
	"github.com/anwardh/meliProject/pkg/store"
	"github.com/anwardh/meliProject/pkg/web"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func respondWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, web.NewResponse(code, nil, message))
}

func TokenAuthMiddleware() gin.HandlerFunc {
	requiredToken := os.Getenv("TOKEN")

	// Verificação do token
	if requiredToken == "" { // Se o valor do token estiver vazio
		log.Fatal("por favor, configure a variável de ambiente - token")
	}

	return func(c *gin.Context) {
		token := c.GetHeader("token")

		if token == "" { // Se token que estiver no Header for vazio
			respondWithError(c, http.StatusUnauthorized, "API token obrigatório")
			return
		}

		if token != requiredToken { // Se o token da Header for diferente

			respondWithError(c, http.StatusUnauthorized, "token do API inválido")
			return
		}
		c.Next()
	}
}

/*
Instanciamos cada camada do domínio Products e usaremos os métodos do controlador para cada endpoint.
*/
// @title MELI Bootcamp API
// @version 1.0
// @description This API Handle MELI Products.
// @termsOfService https://developers.mercadolibre.com.ar/es_ar/terminos-y-condiciones

// @contact.name API Support
// @contact.url https://developers.mercadolibre.com.ar/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	config.InitConfig()

	store := store.Factory("arquivo", "products.json")
	if store == nil {
		log.Fatal("Não foi possivel criar a store")
	}
	repo := products.NewJsonRepository(store)
	service := products.NewService(repo)
	p := handler.NewProduct(service)

	r := gin.Default()
	pr := r.Group("/products")
	{
		pr.Use(TokenAuthMiddleware())

		pr.POST("/", p.Store())
		pr.GET("/", p.GetAll())
		pr.PUT("/:id", p.Update())
		pr.PATCH("/:id", p.UpdateName())
		pr.DELETE("/:id", p.Delete())
	}

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
