package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "os/exec"
    "runtime"
    "strings"
)


const PORT string = "6666"


type Client struct {
    incomingReader *bufio.Reader
    outgoingReader *bufio.Reader
    writer         *bufio.Writer
}

func outputMessage(message string) {
    fmt.Println(message)

    command := ""
    if runtime.GOOS == "darwin" {
        command = "say"
    } else if runtime.GOOS == "linux" {
        command = "espeak"
    }

    exec.Command(command, message).Run()
}

func (client *Client) ReadIncoming() {
    for {
        line, _ := client.incomingReader.ReadString('\n')
        line = strings.TrimRight(line, " \t\r\n")
        outputMessage(line)
    }
}

func (client *Client) ReadOutgoing() {
    for {
        line, _ := client.outgoingReader.ReadString('\n')
        client.writer.WriteString(line)
        client.writer.Flush()
    }
}

func (client *Client) Start() {
    go client.ReadIncoming()
    client.ReadOutgoing()
}

func NewClient(connection net.Conn) *Client {
    incomingReader := bufio.NewReader(connection)
    outgoingReader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(connection)

    client := &Client{
        incomingReader: incomingReader,
        outgoingReader: outgoingReader,
        writer: writer,
    }

    return client
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: ", os.Args[0], "host")
        os.Exit(1)
    }

    host := os.Args[1]
    conn, err := net.Dial("tcp", host + ":" + PORT)

    if err != nil {
        log.Fatal(err)
    }

    client := NewClient(conn)
    client.Start()
}

