package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const PORT string = "8080"

// Un message contient l'expéditeur et le texte
type message struct {
	text string
	from string
}

func loopInfini(listener net.Listener, clients map[net.Conn]string, messages chan message) {
	// loop infini (loop sans condition):
	for {
		// 4. Accepter une nouvelle connexion
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		// 5. Gérer le nouveau client dans une goroutine dédiée
		go func(connection net.Conn) {
			defer connection.Close()
			
			fmt.Fprint(connection, "Entrez votre nom: ")
			nameInput, _ := bufio.NewReader(connection).ReadString('\n')
			name := strings.TrimSpace(nameInput)
			
			clients[connection] = name // Ajouter au map des clients
			fmt.Printf("%s a rejoint le chat.\n", name)

			// Boucle de lecture des messages du client
			scanner := bufio.NewScanner(connection)
			for scanner.Scan() {
				msgText := scanner.Text()
				messages <- message{text: msgText, from: name}
			}

			// Nettoyage à la déconnexion
			delete(clients, connection)
			fmt.Printf("%s a quitté le chat.\n", name)
		}(conn)
	}
}

func main() {
	// 1. Écouter sur le port 8080
	listener, _ := net.Listen("tcp", ":"+PORT)
	defer listener.Close() // s'exécute à la toute fin/sortie de la fonction

	// 2. Channels pour gérer les clients et les messages
	clients := make(map[net.Conn]string)
	messages := make(chan message)

	fmt.Println("Serveur de chat démarré sur :8080")

	// 3. Goroutine de "Broadcasting" (Le hub)
	// Elle tourne en boucle et distribue chaque message reçu à tous les clients
	go func() {
		for msg := range messages {
			for conn := range clients {
				fmt.Fprintf(conn, "[%s]: %s\n", msg.from, msg.text)
			}
		}
	}()

	// Ceci DOIT impérativement être exécuté à la fin.
	// Parce que c'est un loop infini (contient un `for` sans condition).
	loopInfini(listener, clients, messages)
}
