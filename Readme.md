# Advertising

Запустить докер:
```
docker-compose up --build
```

Проверить с помощью golanci-lint:
```
golangci-lint run
```

# Endpoints
```
## Create advert
[POST] http://localhost:8000/api/adverts/  

Example Body: 
{
    "title": "Example",
    "description": "Example description",
    "photos": ["https://sun9-34.userapi.com/impg/3WOMrjs0H5io1nFdLaHv_NiOi5lxrz9qkk7RXg/-IRxizgUjRQ.jpg?size=960x1280&quality=95&sign=5f33c827a0d6a1ce884d82fe1202541d&type=album", "https://sun9-68.userapi.com/impg/18OF1APOug-EIq6K63oIjxqR2wYN43DifTF-zw/PLhXQGZHl4c.jpg?size=1200x1600&quality=95&sign=94a5e71255e270fe46ba5d8f68f02770&type=album"],
    "price": 100.00
}
```
```
## Get all adverts
[GET]  http://localhost:8000/api/adverts/  
```
```
Get advert by id
[GET]  http://localhost:8000/api/adverts/:id 
```


