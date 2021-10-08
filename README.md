# GO-React starter
![](./assets/logo.png)

This is a basic example of a go web server with a react frontend.

It uses the [go fiber](https://github.com/gofiber/fiber) framework 

## Getting started

### Running locally
Clone this repository
Download and install [golang](https://golang.org)

Download and install [postgres](https://www.postgresql.org/download/)

setup your postgres database, enter your config secrets in the [.env](./server/.env)

- [A complete guide to PostgreSQL](https://prabhupant.medium.com/a-complete-guide-to-postgresql-e4d1cefb9866)

- [Installing PostgreSQL for Mac, Linux, and Windows](https://medium.com/@dan.chiniara/installing-postgresql-for-windows-7ec8145698e3)

```bash
cd server
go mod download
go run main.go
```

This will start the go server.

To start the react app navigate to the client directory

```bash
cd client
yarn install
yarn start
```
### Using docker
using docker compose 

```bash
docker-compose build
docker-compose up
```


## Endpoints
| endpoint      | method | body                                           | description       |
|---------------|--------|------------------------------------------------|-------------------|
| /api/session  | GET    |                                                | GET user session  |
| /api/login    | POST   | { email String, password String }              | login user      |
| /api/register | POST   | { email String, password String, name String } | register new user |
|               |        |                                                |                   |


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License
[MIT](https://choosealicense.com/licenses/mit/)
