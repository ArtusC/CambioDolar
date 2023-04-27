package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
)

func ClientHandler(bid string) error {

	ctx := context.Background()
	select {
	case <-ctx.Done():
		return errors.New("client time out")
	case <-time.After(300 * time.Millisecond):
		file, err := os.Create("cotacao.txt")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error to create the file via client, error: %v\n", err)
			return err
		}
		defer file.Close()
		_, err = file.WriteString(fmt.Sprintf("Dólar: %s\n", bid))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error to write the file via client, error: %v\n", err)
			return err
		}
	}

	return nil
}
