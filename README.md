# UrlShortener

## Fast application starting

### 0. Clone project
```bash
$ git clone https://github.com/RomanLevBy/UrlShortener.git && cd UrlShortener
```

### 1. Run project
```bash
make up
```

## Work with application

After application starting the app url is http://localhost:8087

### Authentication
Application use basic auth for create and delete a short url. The username and password can de set in ./config/local.yaml.
Default username is **basic_user** and default password is **basic_pass**


### The following API point is available:

### 1. Create short url
```http
GET /api/url
```
This is a POST request, submitting data to an API via the request body. This request submits JSON data.

A successful POST request typically returns a 200 OK.

Body raw (json)

```
{
    "url": "https://google.com",
    "alias" : "google"
}
```

### 2. Delete url 
```http
DELETE /api/url/{alias}
```
This is a DELETE request, and it is used to delete data that was previously created via a POST request.
You typically identify the entity being deleting by including an identifier in the URL (/api/url/google).

### 3. Redirect url
```http
GET /{alias}
```
This is a GET request and it is used to redirect to url from the short url. 
There is no request body for a GET request, but you should specify the url alias you want (e.g., http://localhost:8087/google).