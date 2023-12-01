package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
	// Crear una nueva instancia del cliente de MongoDB con la URI de conexión proporcionada
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Establecer un contexto con un límite de tiempo para controlar la duración del proceso de conexión
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Asegurar que se llame a cancel para liberar recursos relacionados con el contexto

	// Conectar el cliente al servidor de MongoDB usando el contexto creado
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Verificar si el cliente puede realizar un ping al servidor de MongoDB para asegurar la conectividad
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Error al conectar con MongoDB")
		return nil
	}
	fmt.Println("Conexión exitosa a MongoDB")
	return client // Devolver el cliente de MongoDB conectado
}

var Client *mongo.Client = DBSet() // Inicializa una variable global Client con el cliente de MongoDB obtenido de DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	// Obtiene la colección específica del cliente de MongoDB para la base de datos "Ecommerce"
	var collection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return collection // Devuelve la colección obtenida
}
