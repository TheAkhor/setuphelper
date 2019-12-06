Setup Helper is a scaffold that helps you kickstart your golang programs. It is designed using echo, mongo and java webtokens.

It's recommended that you install "gin" and run the program using "gin run server.go" while doing development.

Currenlty this project is configured to connect to a local db. You can change the connection by editing api->utilities->database.go

clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")


Create a JWT 
`
curl -X POST \
  http://localhost:3000/login  -H 'Content-Type: application/json' -d '{"userName": "Test2", "password": "Password"}'
`

Create a USER

`
curl -X POST \
  http://localhost:3000/restricted/users \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9uIFNub3ciLCJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTc1NzU1Njk4fQ.Ifunqlmwmy_6x5CAvzWenO8dOxohiJrcHYPxnOawM0Y' \
  -H 'Content-Type: application/json' \
  -d '{
    "user": {
        "firstName": "blaine",
        "lastName": "Everingham",
        "userName": "Test2",
        "password": "Password"
    },
    "contact": {
        "firstName": "Nested",
        "lastName": "Nested",
        "name": "test",
        "birthdate": "20090908",
        "securitylist": [
            {"level":"admin"}
        ],
        "phoneList": [
            {
                "type": "mobile",
                "number": "4035543965"
            }
        ],
        "emailList": [
            {
                "email": "test@test.com"
            }
        ],
        "photoList": [
            {
                "file": "/tmp/file/directory"
            }
        ]
    }
}'
`

