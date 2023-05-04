# CambioDolar 

## Chalange Client-Server API

This little API follows the requirements of the task suggested by `FullCycle` in the `GoExpert` course.

The task is:

```
Você precisará nos entregar dois sistemas em Go:
- client.go
- server.go

Os requisitos para cumprir este desafio são:

O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.

O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida deverá retornar no formato JSON o resultado para o cliente.

Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.

O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON). Utilizando o package "context", o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.

O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}

O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e a porta a ser utilizada pelo servidor HTTP será a 8080.

Ao finalizar, envie o link do repositório para correção.
```

---

## Requirements
* You need to have installed locally:
  * Docker and docker-compose
  * Golang (version > 1.19)

---

## Run the API

1) Clone/download the repository to a local folder;

2) Via terminal, access the cloned repository folder;

3) Run this command to load the MySql image DB:
  * ``` docker-compose up -d ```

4) Run this command to start the API:
  * ``` go run main.go ```
  * **DETAIL**: if everything happens well, a file named `cambio_dolar.db` from SqLite apears in the curent folder.

5) Access this link through the browser  to see the result:
  * ``` localhost:8080/cotacao ```
  * **DETAIL**: if the API return a log with `context deadline exceeded`, reload the browser page.


---

## Result

After accessing the localhost link, the result of the API will apear in the browser screen.

* Two files will be created in the main project folder if all goes well:
  * A txt file called `cotacao.txt`, containing the dollar quotation value in reais at the time you make the request.
  * A db file called `cambio_dolar.db`, containing the Sqlite database records.
  
### Mysql DB:
If you want to see the result in this database, run the these commands in another terminal:

1) ``` docker exec -it mysql bash ```

2) ``` mysql -uroot -proot ```

3) Check if the database `cambio_dolar` was successfully created:
  * ``` show databases; ```

4) Access the database:  
  * ``` use cambio_dolar; ```

5) Check if the table `cotacoes` was successfully created:
  * ``` show tables; ```

6) Get the result of the table:
  * ``` select * from cotacoes; ```


### SqLite DB:
If you want to see the result in Sqlite database, run the these commands in another terminal:

1) You need to have the `sqlite3` installed to run this command:
  * ``` sqlite3 ```

2) Check if the table `cotacoes` was successfully created running the following command:
  * ``` .tables ```

3) Get the result of the table cotacoes running the following command:
  * ``` select * from cotacoes; ```