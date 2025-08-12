package main

import (
	"log"

	"github.com/gjcms/taxi_service/routes"

	"github.com/gjcms/taxi_service/database"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

const gopherDraw = `
       ,_---~~~~~----._
_,,_,*^____      _____''*g*\"*,          Welcome to your app!
/ __/ /'     ^.  /      \ ^@q   f     /
[  @f | @))    |  | @))   l  0 _/    /
\'/   \~____ / __ \_____/    \      /
|           _l__l_           I     /
}          [______]           I   /
]            | | |            |
]             ~ ~             |
|                            |
|                           |`

func main() {
	app := fiber.New()
	database.ConnectDb()

	routes.SetupRoutes(app)

	log.Println(gopherDraw)

	app.Listen(":3000")
}
