package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 全てのオリジンを許可する
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// WebSocketのアップグレードを行う
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// クライアントを登録する
	clients[conn] = true

	// WebSocketの受信ループ
	for {
		// クライアントからメッセージを受信する
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			delete(clients, conn)
			break
		}

		// 受信したメッセージを全てのクライアントに送信する
		broadcast <- message
	}

}

func handleMessages() {
	for {
		// ブロードキャストチャネルからメッセージを取得する
		message := <-broadcast

		// 全てのクライアントにメッセージを送信する
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	go handleMessages()
	fmt.Println("WebSocket server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
