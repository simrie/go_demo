# go_demo.git

This is a demo of a password hashing API using only libraries that are part of Golang, and no persistent database.  

## Usage

### Build

```
go build
```

## Test the Project's Functions

```
go test -v ./hasher

go test -v ./store

```

### Run the Server

```
.\go_demo.git
```

### cUrl commands

The following instructions assume that the server is accessible by cUrl commands at "localhost:8080".

#### Get Stats

```
curl http://localhost:8080/stats
```

Expected Result looks like this:

```
{
  "total": 3,
  "average": 5000000000
}
```

#### Post a Password to Encode

Replace the "angryMonkey" in the cUrl below with the string you would like to encode. (Or leave it as it is to encode the string "angryMonkey").

```
curl --data 'password=angryMonkey' http://localhost:8080/hash
```

Expected result is an integer representing the key to retrieve the encoded password after five seconds have passed.

#### Retrieve an Encoded Password

Repace the {id} with an integer received from posting a string to hash.

```
curl http://localhost:8080/hash/{id}
```

If five seconds have passed, expect a string that represents a SHA512 encrypted, base64 encoded hash value.

For example, if you encoded the string "angryMonkey" and received the integer you are now passing as the id, you should see:

```
"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
```

If 5 seconds have no passed since submitting the string to hash, an error will result.

```
Error": Unable to retrieve by id
```


#### Shutdown the Server

```
curl http://localhost:8080/stats
```

The server terminal should display a message indicating that the server is shutting down.


## go mod

A go.mod file was produced with go mod init.

```
go mod init
```

The go.mod file produced shows the version of Goused in development.

```
go mod tidy
```

Running go mod tidy does not produce any go.sum,  probably because there are no dependencies.


