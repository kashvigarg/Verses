# Verse-S
A social media Backend APi

## API DOCUMENTATION

### Base URL

### Endpoints

- #### User Management
   - #### Create User `POST ` `/users`
     Description: Create a new user account.
     
     Request Body:
     ```json
     {
         "Name": "name",
         "Password": "password",
         "Email": "email"
     }
     ```

     Response:
      ```json
     {
         "id": "UUID",
         "email": "string",
         "name": "string",
         "is_gold": "bool"
     }
     ```
     
     Status Code: 201 Created
   - #### User Login `POST` `/login`
     Description: Authenticate a user and generate access and refresh tokens.

     Request Body:
      ```json
     {
      "password": "string",
      "email": "string"
      }
      ```
     
     Response:
      ```json
     {
       "Email": "email",
       "Token": "Token",
       "Refresh_token": "Refresh_token"
      }
      ```
     Status Code: 200 OK
   - #### Refresh Token `POST` `/refresh`
     Description: Refresh the access token using the refresh token.
     
     Request Body:
     
     Header: `Bearer "refresh_token"`
     
     Response:
     ```json
     {
       "Token": "auth_token"
      }
     ```
     Status Code: 200 OK

  - #### Revoke Token `POST` `/revoke`
    Description: Revoke the refresh token.
    
    Request Body:
    
    Header: `Bearer "reresh_token"`
    
    Response:
    
    `"Token Revoked"`
    
    Status Code: 204 No Content
  - #### Update User `PUT` `/users` 
    Description: Update user account information.
    
    Request Body:
    
    Header : `Bearer "auth_token"`
    
    ```json
     {
         "Name": "name",
         "Password": "password",
         "Email": "email"
     }
     ```
    
    Response:
    ```json
     {
         "id": "",
         "email": "string",
         "name": "string",
         "is_gold": ""
     }
     ```
    
    Status Code: 200 OK
- #### Prose Management
  - #### Create Prose `POST` `/prose`
    Description: Create a new prose (post).
    
    Request Body:
    Header: `Bearer "auth_token"`
    ```json
    {
      "Body":"body"
     }
    ```
    
    Response:
    ```json
    {
     "ID":       "uuid",
	"Body":     "string",
	"Created_at":"time",
	"Updated_at":"time"
	}
    ```
    Status Code: 201 Created

  - #### Get Prose `GET` `/prose`
    Description: Get a list of prose.
 
    Request Body:
    
    Header: `Bearer "auth_token"`
    
    Response:
    ```json
    [
     {
	"ID": "uuid"
	"Body": "string"    
	"CreatedAt": "time"
	"UpdatedAt": "time"
     },
    ]
    ```

    Status Code: 200 OK

  - #### Get Prose by ID `GET` `/prose/{proseId}`
    Description: Get details of a prose by its ID.
    
    Request Body:
 
    Header: `Bearer "auth_token"`
    
    Response:
       ```json
    {
	"ID": "uuid"
	"Body": "string"    
	"CreatedAt": "time"
	"UpdatedAt": "time"
     }
    ```
    
    Status Code: 200 OK
  - #### Delete Prose ` DELETE` `/prose/{proseId}`
    Description: Delete a prose by its ID.
 
    Request Body:
 
    Header: `Bearer "auth_token"`
    
    Response:
    
    `" Prose Deleted"`
    
    Status Code: 204 No Content
- #### Administrative Tasks
  - #### Gold Webhooks `POST` `/gold/webhooks`
    Description: Receive webhook notifications related to gold.

    Request Body: 

    Response:

    Status Code: 200 OK
  - #### Metrics `GET` `/metrics`
    Description: Retrieve server metrics.

  - #### Health Check `GET` `/healthz`
    Description: Check the health status of the API.

