# dbo-test project

![GitHub stars](https://img.shields.io/github/stars/FawwazAF/dbo-test)
![GitHub forks](https://img.shields.io/github/forks/FawwazAF/dbo-test)
![GitHub issues](https://img.shields.io/github/issues/FawwazAF/dbo-test)
![GitHub license](https://img.shields.io/github/license/FawwazAF/dbo-test)

## Table of Contents
1. [Introduction](#introduction)
2. [Entity-Relationship Diagram (ERD)](#entity-relationship-diagram-erd)
3. [Data Definition Language (DDL)](#data-definition-language-ddl)
4. [Setup Instructions](#setup-instructions)
5. [API Documentation](#api-documentation)
6. [Contributing](#contributing)
7. [License](#license)

## Introduction
Welcome to **dbo-test project**! This project is a demonstration of API for customer and order management features. It's built with Golang with Gin-Gonic Frameworks, PostgreSQL, and Docker.

## Entity-Relationship Diagram (ERD)
Below is the Entity-Relationship Diagram (ERD) for the database schema used in this project:

![ERD Diagram](https://github.com/user-attachments/assets/761019fd-f921-44fb-9572-57bc9e67f431)

## Data Definition Language (DDL)
Here is the DDL script to create the database schema:
- [Link](https://github.com/FawwazAF/dbo-test/blob/main/init/01-create-schema.sql)

## Setup Instructions
Follow these steps to get the project up and running on your local machine.

Prerequisites
- Docker
- Docker Compose
- Golang

1. Build and run the Docker containers:
```
docker-compose up --build
```
2. Run script for Create Table and populate data
give script execute permission : 
```
chmod +x run-scripts.sh
```
then run the script 
```
./run-scripts.sh
```

## API Documentation
Here is the link to the API Documentation. I'm using postman to create the API Docs.
- [Link](https://github.com/FawwazAF/dbo-test/blob/main/postman/dbo-test%20Service%20API.postman_collection.json)
