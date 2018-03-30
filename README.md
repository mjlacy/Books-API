# Book API

Queries the book collection in a mongo database

## Routes
- GET /

- GET /{id}

- GET /health

- POST /

- PUT /{id}

- DELETE /{id}

### Example Calls

- GET localhost:8080/

- GET localhost:8080/3

- POST localhost:8080/
    - Body (3 Fields):
    ```
    {
        "BookId" : 2,
        "Name" : "War and Peace",
        "Author" : "Leo Tolstoy",
        "Year" : 1869
    }
    ```

- PUT localhost:8080/2
    - Body (3 Fields):
    ```
    {
        "BookId" : 2,
        "Name" : "War and Peace",
        "Author" : "Leo Tolstoy",
        "Year" : 1870
    }
    ```

- DELETE localhost:8080/2

- GET localhost:8080/health

##### Example Response for `GET /2` (5 Fields)

```
{
    "_id": "5a8ddeef55235ccd5c0e5699",
    "BookId": 2,
    "Name": "War and Peace",
    "Author": "Leo Tolstoy",
    "Year": 1869
}
```