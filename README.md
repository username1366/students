# Students RESTful web service

## Run
cd students
export GOPATH=$GOPATH:`pwd`
go get github.com/jinzhu/gorm
go get github.com/jinzhu/gorm/dialects/sqlite
go get github.com/gorilla/mux
cd src
go run main.go

### Add student
`curl -v  -XPOST localhost:8000/ -d'{"name": "Fred", "age": 25, "rating": 50}'`

### Get student by id
`curl -v localhost:8000/1`

### Get all students
`curl -v localhost:8000/`

### Update student info
`curl -v -XPUT localhost:8000/1 -d'{"name": "Sergey", "age": 46, "rating": 75'}`

### Delete student
`curl -v -XDELETE localhost:8000/1`
