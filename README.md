# Book API

Queries the book collection in a mongo database

## Routes
- GET /

- GET /{id}

- GET /health

- POST /

- PUT /{id}

- PATCH /{id}

- DELETE /{id}

### Example Calls

- GET localhost:8080/

- GET localhost:8080/5a80868574fdd6de0f4fa438

- POST localhost:8080/
    - Body (4 Fields):
    ```
    {
        "BookId": 1,
        "Name" : "War and Peace",
        "Author" : "Leo Tolstoy",
        "Year" : 1869
    }
    ```

- PUT localhost:8080/5a80868574fdd6de0f4fa439
    - Body (4 Fields):
    ```
    {
        "BookId": 1,
        "Name" : "War and Peace",
        "Author" : "Leo Tolstoy",
        "Year" : 1870
    }
    ```
    
- PATCH localhost:8080/5a80868574fdd6de0f4fa439
    - Body (1-4 Fields):
    ```
    {
        "BookId": 2
    }
    ```

- DELETE localhost:8080/5aa5841a740db1970dff3248

- GET localhost:8080/health

##### Example Response for `GET /5a8ddeef55235ccd5c0e5699` (4 Fields)

```
{
    "_id": "5a8ddeef55235ccd5c0e5699",
    "Name": "War and Peace",
    "Author": "Leo Tolstoy",
    "Year": 1869
}
```
