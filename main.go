package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	server = newServer()
)

func hello(c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	roomID := c.QueryParam("room_id")
	fmt.Printf("room id: %+v\n", roomID)
	if roomID != "" {
		server.addClientToRoomByID(roomID, conn)
	}
	return nil
}

func processData(c echo.Context) error {
	roomID := c.QueryParam("room_id")
	msg := c.QueryParam("msg")

	fmt.Printf("msg get param: %s\n", msg)
	//TODO: add error msg if room is empty
	if roomID != "" {
		server.broadcastToRoomByID(roomID, msg)
	}
	return c.String(http.StatusOK, "")
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "./public")
	e.GET("/ws", hello)
	e.GET("/data", processData)

	e.Logger.Fatal(e.Start(":1323"))

}
