


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

###  login

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /api/google-login | [google login](#google-login) | Login with google OAuth2 |
  


## Paths

### <span id="google-login"></span> Login with google OAuth2 (*googleLogin*)

```
POST /api/google-login
```

todo

#### Security Requirements
  * Bearer

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#google-login-201) | Created |  | ✓ | [schema](#google-login-201-schema) |

#### Responses


##### <span id="google-login-201"></span> 201
Status: Created

###### <span id="google-login-201-schema"></span> Schema

###### Response headers

| Name | Type | Go type | Separator | Default | Description |
|------|------|---------|-----------|---------|-------------|
| Authorization | string | `string` |  |  | jwt token |

## Models

### <span id="response-data"></span> ResponseData


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Data | [interface{}](#interface)| `interface{}` |  | |  |  |
| ErrorCode | string| `string` |  | |  |  |
| ErrorMessage | string| `string` |  | |  |  |


