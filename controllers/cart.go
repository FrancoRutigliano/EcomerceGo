package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/FrancoRutigliano/EcommerceGolang/database"
	"github.com/FrancoRutigliano/EcommerceGolang/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	return func(c *gin.Context) {
		user_id := c.Query("id")
		// comprobamos si el id que devuelve el query esta vacio
		if user_id == "" {
			// al querer obtener un item del carrito vamos a devolver al header un content/type
			// Encabezado de la solicitud http, con informacion extra de la solicitud en cuestion
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			// abortamos funcion
			c.Abort()
			return
		}

		// de lo que nos devuelve la base de datos probablemente en formato hexadecimal, lo tendremos que convertir para despues pasarlo a la funcion que llama a la base de datos

		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		// vamos a crear un contexto que va a ser creado unicamente para la funcion que llame a la base de datos
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.User

		// BSON.D es una representacion ordenada de un BSON
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledCart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "id not found")
			return
		}

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usertcart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "usertcart.price"}}}}}}
		// Ahora vamos a ejecutar una operacion de agregacion a una colección
		// Estas operaciones de agregación a la base de datos son filter_match, unwind y grouping
		PointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})

		if err != nil {
			log.Println(err)
		}

		var listing []bson.M

		if err = PointCursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledCart.UserCart)
		}
		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == "" {
			log.Println("user ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}

		// Ahora debemos crear un context y una cancelacion del contexto. Todo esto para pasarselo a la funcion que llama a la base de datos
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Vamos a llamar a la funcion que hace conexion con la base de datos
		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		// caso de que haya un problema en la conexion, damos un aviso del error
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully Placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserQueryID := c.Query("userid")
		// En el caso de que el user que devuelve el query este vacio
		if UserQueryID == "" {
			log.Println("User ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		}

		ProductQueryID := c.Query("pid")
		if ProductQueryID == "" {
			log.Println("Product id is empty")
			c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
		}
		// Transformamos  de Hex a lo que devuelve mongo como parametro de la url http
		productID, err := primitive.ObjectIDFromHex(ProductQueryID)
		// si no se pudo hacer la conversion
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Invocamos a la funcion que se va a conectar con la base de datos
		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, UserQueryID)
		// debemos corroborar si el error no esta vacio
		// ya que si esta vacio pudo haber algún problema en la conexion a base de datos
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "Successfully placed the order")
	}
}
