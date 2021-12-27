# User Service

Microservice to manager users.

### Assumptions & Limitations

1. This API aims to be used by a admins, not for the user itself.
   * My concern here is regarding the API structure and password. I would have a endpoint `/v1/user` and provide a mechanism to change the password (requesting the current and the new one). For admins purpose in general you shouldn't have the current password, you'd be allowed to just replace it.
2. I'm using soft-delete for user removal. As far as GDPR is concerned, the soft-delete could also anonymise the personal data.
3. I've created a dummy cache layer (`stub.UserStorageCache`) just to show how I would extend the current codebase adding more functionality if necessary.
4. I've created a dummy publisher (`stub.EventService`) just to make it clear where I would have implemented.
5. The pagination by cursor (moving forward) wasn't implemented only page (limit + offset).
6. I'm using `github.com/rs/xid` to generate the ids which uses the Mongo Object ID algorithm. Main reason here is that I can sort by id and the result will be order of creation as it has a timestamp component. It'll also be the source of my cursor.
7. I'm using `https://github.com/golang/mock` for mock generation.
8. I decided to go only for HTTP API as I'm more confortable with it.
9. In a high traffic environment I'd have change the approach of storing and publishing an event within the same service. Other than that, I'd have a an API which will receive the request (POST/PUT) and publish an event. Then I'd have a different service (or even the same with some flags) where will consume this event and save in the database and it'll also provide the GET methods as well. Drawback would be the user will not be available right after the POST/PUT but only after the consumption of the events.
10. I've worked in the past with TDD, but nowadays I do not anymore. I have focused do do the testing over the business logic, no database test is implemented.

### How to run?

TLDR; `make build up`

I'm using `docker-compose` which will run a instance of mysql 8, and the User Service. The User Service will be exposed on localhost under the port 80, you can change it editing the `docker-compose.yml` file under the section *ports*. If you need to point the service for a different mysql instance it also can be adjusted in the yaml file.

Endpoints available are:
* **POST http://localhost/v1/users**: to create a new user
* **GET http://localhost/v1/users**: to retrieve a list of all users
  * Some query string are accept, like `country`, `per_page` and `page`
* **GET http://localhost/v1/users/{id}**: to retrieve a specific user
* **PUT http://localhost/v1/users/{id}**: to update a specific user
* **DELETE http://localhost/v1/users/{id}**: to delete a specific user

The json format accept is:
```json
{
    "id": "c74tbdnblarkcprj54f0",
    "first_name": "Guilherme",
    "last_name": "S.",
    "nickname": "xguiga",
    "email": "email@gmail.com",
    "country": "DE"
}
```

### How to develop?

This project uses two build stage, the first stage is the builder, which uses a golang docker image to build the project, the final image will be based in a `alpine:3.11`.

For development is suggested that you create a `docker-compose.override.yml` using the command:

```sh
$ cp docker-composer.{dev,override}.yml
```

This new docker-compose will be merged with the main one, but only the builder image will be built, and the filesystem of the project will be available inside of the container under the folder `/go/src/github.com/guilherme-santos/user`. Note: the service will not be running by default, you can type `make exec` or `make exec cmd=/bin/sh` to enter the container, after that you can use any makefile target available e.g. `run`, `test`, `go-generate`.
