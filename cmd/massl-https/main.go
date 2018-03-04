package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("usage: massl-https <command> [<args>]")
		fmt.Println("Commands are: ")
		fmt.Println(" server   Create a massl https server")
		fmt.Println(" client   Create a massl https client")
		return
	}

	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	listenAddrFlag := serverCmd.String("addr", "localhost:2223", "Listen address of your server")

	clientCmd := flag.NewFlagSet("client", flag.ExitOnError)
	targetURLFlag := clientCmd.String("url", "https://localhost:2223/echo", "Server url to post your content")

	switch os.Args[1] {
	case "server":
		serverCmd.Parse(os.Args[2:])
		log.Printf("listen: https://%s", *listenAddrFlag)

		server := NewServer(*listenAddrFlag)
		log.Fatal(server.Listen())

	case "client":
		clientCmd.Parse(os.Args[2:])
		log.Printf("url: %s", *targetURLFlag)

		client := NewClient(*targetURLFlag)
		err := client.Do()
		if err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}
