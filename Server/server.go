package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	client "github.com/ArtusC/cambioDolar/Client"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
)

var (
	resJson []byte
	result  *CambioDolarStruct
)

const (
	cotacoesTable = `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id varchar(255),
		code varchar(255) DEFAULT NULL,
		codein varchar(50) DEFAULT NULL,
		name varchar(50) DEFAULT NULL,
		bid varchar(20) DEFAULT NULL
		);
`
)

func ServerHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	log.Println("[Mysql] Starting server")
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/cambio_dolar")
	if err != nil {
		panic(err)
	}
	log.Println("[Mysql] Server started with success!")
	defer db.Close()

	log.Println("[SqLite] Starting server")
	dbSqlite, err := sql.Open("sqlite3", "cambio_dolar.db")
	if err != nil {
		panic(err)
	}
	var version string
	err = dbSqlite.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[SqLite(%s)] Server started with success!\n", version)
	defer dbSqlite.Close()

	currentTime := time.Now()
	timeStampString := currentTime.Format("2006-01-02 15:04:05")
	msgStart := fmt.Sprintf("Request startetd, %v\n\nWaiting response...\n", timeStampString)
	log.Print(msgStart)
	w.Write([]byte(msgStart))

	log.Println("Request processed with success!")
	result, err = getResult(ctx)
	if err != nil {
		fmt.Printf("Error to get result: %v\n", err)
		return
	}

	resJson, err = json.Marshal(result)
	if err != nil {
		fmt.Printf("Error to Marshal result: %s\n", err)
		return
	}

	msgSuccess := fmt.Sprintf("\nRequest processed with success, result:\n%v\n", string(resJson))
	w.Write([]byte(msgSuccess))

	cotacao := NewCotacaoDolar(result.USDBRL.Code, result.USDBRL.Codein, result.USDBRL.Name, result.USDBRL.Bid)

	log.Println("Saving record on sqlite DB...")
	err = insertCotacaoSqlite(ctx, dbSqlite, cotacao)
	if err != nil {
		panic(fmt.Sprintf("error to insert the record on DB, error: %s", err.Error()))
	}
	log.Println("Record saved with success in sqlite!")

	log.Println("Saving record on mysql DB...")
	err = insertCotacaoMysql(db, cotacao)
	if err != nil {
		panic(fmt.Sprintf("error to insert the record on DB, error: %s", err.Error()))
	}
	log.Println("Record saved with success in mysql!")

	log.Println("Creating txt file...")
	errClient := client.ClientHandler(cotacao.Bid)
	if errClient != nil {
		log.Printf("Error to create the txt file, error: %v\n", errClient)
		return
	}
	log.Println("File created with success!")

	defer log.Println("Request finished with success!")
	defer w.Write([]byte("\nBye :)"))

}

func getResult(ctx context.Context) (*CambioDolarStruct, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(fmt.Sprintf("error to create the new request %s", err.Error()))
	}

	reqRes, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("error to sends the request %s", err.Error()))
	}
	defer reqRes.Body.Close()

	res, err := io.ReadAll(reqRes.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in read the response: %v\n", err)
		return nil, err
	}

	err = json.Unmarshal(res, &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error to unmarshal the response: %v\n", err)
		return nil, err
	}
	fmt.Println("Code: ", result.USDBRL.Code)
	log.Println("Response: ", string(res))

	return result, err
}

func insertCotacaoSqlite(ctx context.Context, db *sql.DB, cotacao *CambioDolarSqlStruct) error {

	select {
	case <-time.After(10 * time.Millisecond):

		log.Printf("Creating table cotacoes in sqlite database\n")
		_, err := db.Exec(cotacoesTable)
		if err != nil {
			log.Printf("Error to create cotacoes table in sqlite database, error: %s\n", err)
			return err
		}

		log.Printf("[SqLite] Inserting record in table cotacoes with values: %s, %s, %s, %s, %s\n", cotacao.ID, cotacao.Codein, cotacao.Code, cotacao.Name, cotacao.Bid)
		stmt, err := db.Prepare("insert into cotacoes(id, code, codein, name, bid) values(?, ?, ?, ?, ?)")
		if err != nil {
			log.Printf("[SqLite] Error to prepare statement to insert record, error: %s\n", err)
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(cotacao.ID, cotacao.Codein, cotacao.Code, cotacao.Name, cotacao.Bid)
		if err != nil {
			log.Printf("[SqLite] Error to insert record in the server, error: %s\n", err)
			return err
		}

	case <-ctx.Done():
		log.Println("[SqLite] Request canceled by client.")
	}

	return nil

}

func insertCotacaoMysql(db *sql.DB, cotacao *CambioDolarSqlStruct) error {
	log.Printf("[Mysql] Inserting record in table cotacoes with values: %s, %s, %s, %s, %s\n", cotacao.ID, cotacao.Codein, cotacao.Code, cotacao.Name, cotacao.Bid)
	stmt, err := db.Prepare("insert into cotacoes(id, code, codein, name, bid) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[Mysql] Error to prepare statement to insert record, error: %s\n", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cotacao.ID, cotacao.Codein, cotacao.Code, cotacao.Name, cotacao.Bid)
	if err != nil {
		log.Printf("[Mysql] Error to insert record in the server, error: %s\n", err)
		return err
	}

	return nil

}

func NewCotacaoDolar(code, codein, name, bid string) *CambioDolarSqlStruct {
	return &CambioDolarSqlStruct{
		ID:     uuid.New().String(),
		Code:   code,
		Codein: codein,
		Name:   name,
		Bid:    bid,
	}
}
