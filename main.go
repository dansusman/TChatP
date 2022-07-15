package main

import (
	"log"
	"net"
)

// spin up the server and wait for messages to arrive
func main() {
    server := newServer()
    go server.run()

    ln, err := net.Listen("tcp", ":3333")
    handleErr("unable to start server!", err)

    defer ln.Close()

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("unable to accept connection: %s", err.Error())
            continue
        }
        client := server.newClient(conn)
        go client.readActions()
    }
}


func handleErr(msg string, err error) {
    if err != nil {
        log.Fatal(msg + " " + err.Error() + "\n")
    }
}
