package clientcli

import (
	"bufio"
	"context"
	"fmt"
	"github.com/xpwu/go-cmd/cmd"
	"github.com/xpwu/go-x/exe"
	"log"
	"os"
)

func init() {
	cmd.RegisterCmdNoArgs("client", "client cli mode", func() {
		client()
	})
}

func client() {

	log.SetOutput(os.Stdout)

	write, err := ChanFromServer(context.TODO())
	if err != nil {
		fmt.Print(fmt.Sprintf("Connection to service(%s) failed. The service may not have started yet.", exe.Name))
		os.Exit(1)
	}

	fmt.Println("CLIENT CLI has started. The client cli mode exit does not affect the operation of the server.")

	response := make(chan string)

	write <- Request{
		Content:  "hello",
		Response: response,
	}
	out, ok := <-response
	if !ok {
		fmt.Println("error! Maybe server stopped")
		os.Exit(1)
	}
	fmt.Println(out)
	fmt.Print("> ")

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if line == "q" {
			break
		}
		if line == "" {
			fmt.Print("> ")
			continue
		}

		write <- Request{
			Content:  line,
			Response: response,
		}
		out, ok = <-response
		if !ok {
			fmt.Println("error! Maybe server stopped")
			os.Exit(1)
		}
		fmt.Println(out)
		fmt.Print("> ")
	}
}
