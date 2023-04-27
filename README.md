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

2) Via terminal, access the cloned repository;

3) Run the command to load the DB:
  * ``` docker-compose up -d ```

4) Run the command to start the API:
  * ``` go run main.go ```

5) Access the link through the browser:
  * ``` localhost:8080/cotacao ```


---

## Result

After accessing the localhost link, the result of the API will apear in the browser screen.

A txt file called `cotacao.txt` will be created in the folder containing the dollar quotation value in reais at the time you make the request.

If you want to see the result in the database, run the these commands in another terminal:

1) ```docker exec -it mysql bash ```

2) ```mysql -uroot -proot```

3) Check if the database `cambio_dolar` was successfully created:
  * ```show databases; ```

4) Access the database:  
  * ```use cambio_dolar;```

5) Check if the table `cotacoes` was successfully created:
  * ```show tables;```

6) Get the result of the table:
  * ```select * from cotacoes;```