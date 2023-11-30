package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Creamos una colección User para mongoDB
// ID es un identificador único para objetos, probablemente proporcionado por una base
// de datos NoSQL como MongoDB.
// First_Name y Last_Name, Password... almacenan el info del usuario como punteros a strings.
// Created_At guarda la fecha y hora en la que se creó este registro utilizando el tipo time.Time.
// Updated_At almacena la fecha y hora de la última actualización en este registro utilizando el tipo time.Time.
// UserCart es una lista o array que contiene objetos de tipo ProductUser, posiblemente
// almacenando los productos que el usuario tiene en su carrito.
// Address_Details es una lista o array que contiene objetos de tipo Address, probablemente
// almacenando la información de direcciones asociadas al usuario.
// Order_Status es una lista o array que contiene objetos de tipo Order,
// representando el historial o estado de los pedidos realizados por el usuario.
type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name      *string            `json:"first_name" validate:"required,min=2,max=30"`
	Last_Name       *string            `json:"last_name" validate:"required,min=2,max=30"`
	Password        *string            `json:"password" validate:"required,min=6"`
	Email           *string            `json:"email" validate:"email, required"`
	Phone           *string            `json:"phone" validate:"required"`
	Token           *string            `json:"token"`
	Refresh_Token   *string            `json:"refresh_token"`
	Created_At      time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updated_at"`
	User_ID         string             `json:"user_id"`
	UserCart        []ProductUser      `json:"usercart" bson:"usercart"`
	Address_Details []Address          `json:"address_details" bson:"address"`
	Order_Status    []Order            `json:"order_status" bson:"orders"`
}

// Coleccion Products para MongoDB
type Products struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name"`
	Price        *uint64            `json:"price"`
	Rating       *uint8             `json:"rating"`
	Image        *string            `json:"image"`
}

// Coleccion de ProductUser para MongoDB
type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        int                `json:"price" bson:"price"`
	Rating       *uint8             `json:"rating" bson:"rating"`
	Image        *string            `json:"image" bson:"image"`
}

// Coleccion de Address para MongoDB
type Address struct {
	Address_id primitive.ObjectID `bson:"_id"`
	House      *string            `json:"house_name" bson:"house_name"`
	Street     *string            `json:"street_name" bson:"street_name"`
	City       *string            `json:"city_name" bson:"city_name"`
	Pincode    *string            `json:"pin_code" bson:"pin_code"`
}

// Coleccion de Order para MongoDB
type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_list" bson:"order_list"`
	Ordered_at     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Price          int                `json:"total_price" bson:"total_price"`
	Discount       *int               `json:"discount" bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

// COD = Cash on Delivery
type Payment struct {
	Digital bool
	COD     bool
}
