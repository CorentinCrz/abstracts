# abstracts
## Description

This is a school project. The requirements are to create
a restfull api, using [Go](https://golang.org/), to
manage some Books. We must
use a noSQL database and decided to use [Elasticsearch](https://www.elastic.co/fr/elasticsearch/).

A complete documentation with a clean repository is mandatory.

## Developers

- Carolyne FERNANDEZ PRADA
- Corentin CROIZAT
- Clément HALLER

## Table of Contents

- [Directory Structure](#directory-structure)
- [Model Structure](#model-structure)
- [Used Package](#used-package)
- [Api Documentation](#api-documentation)
- [Installation](#installation)

## Directory Structure

```
<Abstract>/
├─ api/
|   └─ server/
|       └─ router.go
|       └─ server.go
├─ controllers/
|   └─ controller.go
├─ documentation/
|   └─ swagger.json
├─ models/
|   └─ book.go
├─ service/
|   └─ book.go
|   └─ service.ho
├─ view/
├─ docker-compose.yml
├─ main.go
└─ README.md
```

## Model Structure

``` Go
type Book struct {
	Id    string
	Title string
	Author string
	Abstract string
}
```

## Used Package

* [mux](https://github.com/gorilla/mux) - HTTP request router and dispatcher
* [go-elasticsearch](https://github.com/elastic/go-elasticsearch) - The official Go client for [Elasticsearch](https://www.elastic.co/fr/elasticsearch/).
* [swagger](https://github.com/go-swagger/go-swagger) - Go implementation of [Swagger 2.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md)
* [CORS](https://github.com/rs/cors) - Go cors handler

## Pourquoi Elastic Search

Elasticsearch a été  conçue pour faire des requêtes full-text très rapides.Ce qui nous semblait très pertinent pour faire de la recherche sur des résumés titre ou des auteurs des livres.


## Api Documentation

Once docker-compose is up you can navigate to
`http://localhost:8085/`


## Api server
`http://localhost:8080`

## Installation

* Launch Elasticsearch, Kibana, swagger, and api
``` bash
docker-compose up -d --build
```
Wait for the Elasticsearch container to mount
