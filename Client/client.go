package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func ClientHandler(bid string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Printf("Error to create the file via client, error: %v\n", err)
		return err
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s\n", bid))
	if err != nil {
		log.Printf("Error to write the file via client, error: %v\n", err)
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("Error: client timed out\n")
		panic(ctx.Err())
	}

	return nil
}
