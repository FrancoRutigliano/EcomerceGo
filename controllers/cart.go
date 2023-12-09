package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/FrancoRutigliano/EcommerceGolang/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Definición de la estructura Application que contiene referencias a colecciones de MongoDB.
// Para que recordemos en mongo, las colecciones son un conjunto de documentos que almacena
// información de manera muy similar a las bases de datos relacionales.
type Application struct {
	prodCollection *mongo.Collection // Colección de productos
	userCollection *mongo.Collection // Colección de usuarios
}

// NewApplication es una función que actúa como constructor para la estructura Application.
// Crea una nueva instancia de Application con las colecciones proporcionadas.
func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection, // Asigna la colección de productos proporcionada al campo prodCollection
		userCollection: userCollection, // Asigna la colección de usuarios proporcionada al campo userCollection
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 'necesitamos checkear' si el id del producto esta en la base de datos, si existe
		// Por eso le pasamos una query de id al contexto.
		productQueryID := c.Query("id")
		// ahora checkeamos si el productQueryID esta vacio
		// Porque no podemos agregar un producto al carrito si no tenemos un id
		if productQueryID == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		// También debemos verificar si el id de usuario existe
		// Para: Integridad de los datos, Seguridad de control y acceso y para tener un registro de la actividad
		userQueryID := c.Query("userID")
		// Hacemos el check de si el usuario esta vacío, caso de que si, abortamos con error
		if userQueryID == "" {
			log.Println("user id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		}

		// El id de producto fue recibido
		// En este caso, primitive.ObjectIDFromHex()
		// Lo que hace es tomar una cadena Hexadecimal como argumento y trata de convertirla a ObjectID
		// Ya que quizas productQuerID = 5ff1e194b8576f48c2f8c7a1, esta funcion la convierte.
		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// Ahora ya deberíamos poder llamar a la funcion que conceta con la DB en database
		// para eso tenemos que pasarle el context
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)

		// si sucede algún error al momento de conectar a base de datos para agregar el producto
		if err != nil {
			// json identado, con un error de server
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		// status 200 se utiliza para saber que el proceso se terminó exitosamente
		// IndentedJson devuelve una estructura Json 'más cuidada', pero cuidado, debería ser unicamente por motivos de desarrollo, porque consume más cpu y banda ancha
		c.IndentedJSON(200, "Successfully Added to the cart")

	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Vamos a estar extrayendo el valor id de una solicitud http
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is invalid")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		// Extraer el valor de el userID de una solicitud http
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
			return
		}

		// Ahora necesito transformar lo que probablemente venga en formato hexadecimal desde la solicitud http.
		// Por ejemplo http://ejemplo.com/ruta?id=5ff1e194b8576f48c2f8c7a1
		ProductID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		// El contexto que vamos a declarar va a ser pasado a la funcion que hace conexion con la DB
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Ahora invocamos a la funcion que conecta y realiza los cambios en la base de datos
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
		// Deberíamos comprobar si la conexion salió bien
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		// si todo salio bien
		c.IndentedJSON(200, "Successfully removed from cart")
	}
}

func GetItemFromCart() gin.HandlerFunc {
	panic("Obtener un item del carrito")
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	panic("Func paraComprar del carrito")
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	panic("Func para Compra instantanea")
}
