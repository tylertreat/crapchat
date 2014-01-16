package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)


const PORT string = "6666"


type Client struct {
    incoming chan string
    outgoing chan string
    reader   *bufio.Reader
    writer   *bufio.Writer
}

func (client *Client) Read() {
    for {
        line, _ := client.reader.ReadString('\n')
        client.incoming <- line
    }
}

func (client *Client) Write() {
    for data := range client.outgoing {
        client.writer.WriteString(data)
        client.writer.Flush()
    }
}

func (client *Client) Listen() {
    go client.Read()
    go client.Write()
}

func NewClient(connection net.Conn) *Client {
    writer := bufio.NewWriter(connection)
    reader := bufio.NewReader(connection)

    client := &Client{
        incoming: make(chan string),
        outgoing: make(chan string),
        reader: reader,
        writer: writer,
    }

    client.Listen()

    return client
}

type Server struct {
    clients []*Client
    joins chan net.Conn
    incoming chan string
    outgoing chan string
}

func (server *Server) Broadcast(data string) {
    fmt.Println(strings.TrimRight(data, " \t\r\n"))
    for _, client := range server.clients {
        client.outgoing <- data
    }
}

func (server *Server) Join(connection net.Conn) {
    client := NewClient(connection)
    server.clients = append(server.clients, client)
    go func() { for { server.incoming <- <-client.incoming } }()
}

func (server *Server) Listen() {
    go func() {
        for {
            select {
            case data := <-server.incoming:
                server.Broadcast(data)
            case conn := <-server.joins:
                server.Join(conn)
            }
        }
    }()
}

func NewServer() *Server {
    server := &Server{
        clients: make([]*Client, 0),
        joins: make(chan net.Conn),
        incoming: make(chan string),
        outgoing: make(chan string),
    }

    server.Listen()

    return server
}

func main() {
    server := NewServer()

    listener, _ := net.Listen("tcp", ":" + PORT)
    fmt.Println("Chat server started on port " + PORT)

    for {
        conn, _ := listener.Accept()
        fmt.Println("User joined the chat!")
        server.joins <- conn
    }
}

