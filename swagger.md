


# api server
This is api server
  

## Informations

### Version

0.1

## Content negotiation

### URI Schemes
  * http
  * https

### Consumes
  * application/json

### Produces
  * application/json

## Access control

### Security Schemes

#### Bearer (header: Authorization)



> **Type**: apikey

## All endpoints

###  cocktail

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /api/cocktails | [popular cocktail list request](#popular-cocktail-list-request) | Get popular cocktail list |
  


###  login

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /api/google-authenticate | [google authenticate request](#google-authenticate-request) | Get access token. |
| GET | /api/google-login | [google login](#google-login) | Login with google OAuth2 |
  


## Paths

### <span id="google-authenticate-request"></span> Get access token. (*googleAuthenticateRequest*)

```
POST /api/google-authenticate
```

Use Code to exchange access token.

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| Body | `body` | [GoogleAuthenticateRequest](#google-authenticate-request) | `models.GoogleAuthenticateRequest` | |  | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#google-authenticate-request-201) | Created |  |  | [schema](#google-authenticate-request-201-schema) |

#### Responses


##### <span id="google-authenticate-request-201"></span> 201
Status: Created

###### <span id="google-authenticate-request-201-schema"></span> Schema
   
  

[GoogleAuthenticateResponse](#google-authenticate-response)

### <span id="google-login"></span> Login with google OAuth2 (*googleLogin*)

```
GET /api/google-login
```

Will redirect with authorization code.

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [307](#google-login-307) | Temporary Redirect | redirect |  | [schema](#google-login-307-schema) |

#### Responses


##### <span id="google-login-307"></span> 307 - redirect
Status: Temporary Redirect

###### <span id="google-login-307-schema"></span> Schema

### <span id="popular-cocktail-list-request"></span> Get popular cocktail list (*popularCocktailListRequest*)

```
POST /api/cocktails
```

Get popular cocktail list order by create date.

#### Security Requirements
  * Bearer: []

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| Body | `body` | [GetPopularCocktailListRequest](#get-popular-cocktail-list-request) | `models.GetPopularCocktailListRequest` | |  | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#popular-cocktail-list-request-200) | OK |  |  | [schema](#popular-cocktail-list-request-200-schema) |
| [400](#popular-cocktail-list-request-400) | Bad Request | bad request |  | [schema](#popular-cocktail-list-request-400-schema) |
| [401](#popular-cocktail-list-request-401) | Unauthorized | unauthorized |  | [schema](#popular-cocktail-list-request-401-schema) |
| [404](#popular-cocktail-list-request-404) | Not Found | item not found |  | [schema](#popular-cocktail-list-request-404-schema) |
| [500](#popular-cocktail-list-request-500) | Internal Server Error | internal error |  | [schema](#popular-cocktail-list-request-500-schema) |

#### Responses


##### <span id="popular-cocktail-list-request-200"></span> 200
Status: OK

###### <span id="popular-cocktail-list-request-200-schema"></span> Schema
   
  

[GetPopularCocktailListResponse](#get-popular-cocktail-list-response)

##### <span id="popular-cocktail-list-request-400"></span> 400 - bad request
Status: Bad Request

###### <span id="popular-cocktail-list-request-400-schema"></span> Schema

##### <span id="popular-cocktail-list-request-401"></span> 401 - unauthorized
Status: Unauthorized

###### <span id="popular-cocktail-list-request-401-schema"></span> Schema

##### <span id="popular-cocktail-list-request-404"></span> 404 - item not found
Status: Not Found

###### <span id="popular-cocktail-list-request-404-schema"></span> Schema

##### <span id="popular-cocktail-list-request-500"></span> 500 - internal error
Status: Internal Server Error

###### <span id="popular-cocktail-list-request-500-schema"></span> Schema

## Models

### <span id="get-popular-cocktail-list-request"></span> GetPopularCocktailListRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Page | int64 (formatted integer)| `int64` | ✓ | |  | `1` |
| PageSize | int64 (formatted integer)| `int64` | ✓ | |  | `10` |



### <span id="get-popular-cocktail-list-response"></span> GetPopularCocktailListResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| PopularCocktailList | [][PopularCocktailList](#popular-cocktail-list)| `[]*PopularCocktailList` |  | |  |  |
| Total | int64 (formatted integer)| `int64` |  | |  |  |



### <span id="google-authenticate-request"></span> GoogleAuthenticateRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | string| `string` |  | |  |  |



### <span id="google-authenticate-response"></span> GoogleAuthenticateResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Token | string| `string` |  | |  |  |



### <span id="popular-cocktail-list"></span> PopularCocktailList


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| CocktailID | int64 (formatted integer)| `int64` |  | |  |  |
| CreatedDate | string| `string` |  | |  |  |
| Photo | string| `string` |  | |  |  |
| Title | string| `string` |  | |  |  |



### <span id="response-data"></span> ResponseData


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Data | [interface{}](#interface)| `interface{}` |  | |  |  |
| ErrorCode | string| `string` |  | |  |  |
| ErrorMessage | string| `string` |  | |  |  |


