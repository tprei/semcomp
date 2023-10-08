# Let' Go: Do nil ao infinito

## O que é esse repositório?

Nesse repositório, temos os códigos e recursos utilizados no minicurso de Go da Semcomp de 2023

## Como executar os códigos?

Foi utilizado Go 1.19 para desenvolver os exemplos desse minicurso. Utilize qualquer versão após e deve haver retrocompatibilidade
Dentro da pasta examples/pokeapi, há um exemplo de chamadas de API paralelizadas. Como esse código não utiliza go modules, para executá-lo, você precisará primeiro compilar todos os arquivos juntos:

`go build *.go -o pokeapi`

E assim poderá executar:

`./pokeapi`

Há também alguns códigos em Python 3. No caso dos códigos de Python 3, talvez seja necessária a instalação das bibliotecas `requests` e `aiohttp`

## Projeto song search

Além do(s) exemplo(s), na pasta search há um exemplo de um repositório estruturado de forma mais próxima àquilo que se encontra na comunidade open source e também em empresas que utilizam a linguagem. Na pasta cmd há alguns entrypoints diferentes. Um deles, chamado de scraper, é utilizado para coletar (em paralelo) músicas e suas letras de diversos gêneros, salvando tudo em um arquivo songs.json de 6MB (que por comodidade está versionado nesse repositório).

Já o outro entrypoint (server) pode ser utilizado para executar um servidor HTTP na porta 8080. Ele é utilizado em conjunto com um simples front (feito com html/css/js puro), para ajudar na visualização do projeto final.

O desafio de live coding proposto é o de implementar uma função que seja capaz de pesquisar nas letras das músicas contidas na "base de dados" (em json)

Esse repositório ainda está em construção e podem haver discrepâncias em relação ao que está descrito nesse README

## Slides utilizados

[Aqui estão os slides utilizados: https://t.ly/semcomp-go](https://t.ly/semcomp-go)

## Contribuições

Qualquer contribuição é bem vinda. Simplesmente submeta um pull request diretamente para esse repositório ou através de um fork.
