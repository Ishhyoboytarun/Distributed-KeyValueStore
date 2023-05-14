# Distributed KeyValueStore using etcd in Golang
This project implements a distributed key-value store using the etcd library in Golang. The key-value store can be accessed by multiple clients concurrently over the network, and provides a consistent and fault-tolerant way of storing and retrieving data.

## Features
Distributed key-value store using etcd for distributed coordination and consensus.
Support for concurrent read and write operations by multiple clients.
Consistent and fault-tolerant data storage and retrieval.
Simple API for adding, updating, and deleting key-value pairs.
Command-line interface for interacting with the key-value store.

## Requirements
### To run this project, you need to have the following installed on your system:

```
Golang (version 1.16 or higher)
etcd (version 3.x)]
```

## Installation
To install and run the key-value store, follow these steps:

### Clone the repository using the following command:

```git clone https://github.com/your-username/distributed-key-value-store.git```

### Navigate to the project directory:

```cd distributed-key-value-store```
### Install the dependencies:

```go mod download```

### Start the key-value store server:
```go run main.go```

## Usage
To interact with the key-value store, you can use the following API endpoints:

GET /get/{key} - Get the value associated with the given key.
POST /set/{key}/{value} - Create a new key-value pair.
PUT /set/{key}/{value} - Update the value associated with the given key.
DELETE /delete/{key} - Delete the key-value pair associated with the given key.

You can use any HTTP client (such as curl or Postman) to make requests to these endpoints.
