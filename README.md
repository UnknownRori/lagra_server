# Lagra Server

E-Commerce app for campus project

## Running

> [!IMPORTANT]
> It's need [`migrate`](https://github.com/golang-migrate/migrate) for migration.
> 
> And go version > 1.22.3

```sh
# Clone the project and enter the directory
git clone https://github.com/UnknownRori/lagra_server
cd lagra_server

# Run migration
migrate -path migrations/ -database "mysql://username:password@/lagra" up

# Run the server and wait client connection
go run .

# Or build the executable
go build
```
