# movie-service

Netflix-type service

Movie-service с CRUD операциями, NoSQL базой данных, аутентификацией, метриками, рассылкой почты и т.д. (production ready). 

**Docker**:

`docker run -d --name mongodb -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo:4.4.3`

**Run**:

`make run`

