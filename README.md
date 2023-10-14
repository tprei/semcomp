# Let' Go: Do nil ao infinito

## O que é esse repositório?

Nesse repositório, temos os códigos e recursos utilizados no minicurso de Go da Semcomp de 2023

## Como executar os códigos?

### Exemplo 1: PokeAPI

Foi utilizado Go 1.19 para desenvolver os exemplos desse minicurso. Utilize qualquer versão após e deve haver retrocompatibilidade

**Dentro da pasta** `examples/pokeapi` há um exemplo de chamadas de API paralelizadas. Como esse código **não utiliza go modules**, para executá-lo, você precisará primeiro compilar todos os arquivos juntos:

`go build *.go -o pokeapi`

E assim poderá executar:

`./pokeapi`

Há também alguns códigos em Python 3. No caso dos códigos de Python 3, talvez seja necessária a instalação das bibliotecas `requests` e `aiohttp`


## Modules

### Module scraper

Além do(s) exemplo(s), na pasta search há um module `scraper`, estruturado de forma mais próxima àquilo que se encontra na comunidade open source e também em empresas que utilizam a linguagem. 

Esse module utiliza o poder das goroutines e channels para coletar músicas e suas letras de diversos gêneros, salvando tudo em um arquivo songs.json de ~12MB (que por comodidade está versionado nesse repositório).

### Module search

Paralelo ao scraper, podemos utilizar o module `search` para executar um servidor HTTP na porta 8080, cuja finalidade é servir como [Backend](https://pt.wikipedia.org/wiki/Front-end_e_back-end) em conjunto com um simples Frontend (feito sem dependências externas, apenas `html/css/js`) para ajudar na visualização do projeto final.

O desafio de live coding proposto é o de implementar uma função (e ferramental necessário) capaz de pesquisar nas letras das músicas contidas na "base de dados" (em JSON)

Toda a comunicação HTTP, leitura/parsing de json e produção da resposta foram abstraídas em packages internos `server` e `database/json`. Isso foi feito propositalmente de forma a simplificar o desafio para que os alunos do curso foquem na prática da sintaxe da linguagem e não fiquem presos a conceitos de web. Adicionalmente, caso haja tempo e os alunos tenham interesse, podemos explorar esses conceitos =).

Esse repositório ainda está em construção e podem haver discrepâncias em relação ao que está descrito nesse README.

## Go workspace

Esse repositório contém mais de um go module (`search` e `scraper`) e portanto possui um arquivo `go.work`. Esse arquivo é similar ao que conhecemos por `go.mod`, porém ao contrário de organizar as dependências de um module, ele organiza um repositório de **múltiplos módulos**, que podem ou não depender uns dos outros.

Normalmente, rodar o comando `go run` fora de um module resulta em um erro. Em um workspace, adicionamos um nível a mais de hierarquia. Isso pode ser especialmente útil quando queremos nos referir a um module dentro de outro, sem precisar a cache criada pelo `go get`. No caso desse repositório, ele foi utilizado apenas para separar as dependências dos modules e permitir que o projeto `search` rode apenas com dependências nativas.

Mais informações sobre go workspaces podem ser lidas [aqui](https://go.dev/doc/tutorial/workspaces).
## Slides utilizados

[Aqui estão os slides utilizados: https://bit.ly/semcomp-go](https://bit.ly/semcomp-go)

## Contribuições

Qualquer contribuição é bem vinda. Simplesmente submeta um pull request diretamente para esse repositório ou através de um fork.
