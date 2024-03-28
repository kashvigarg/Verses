# Verse-S
A social media Backend APi

## API DOCUMENTATION

### Base URL
### Endpoints

- #### User Management
   - #### Create User `/users`
     
      Method: POST
     
      Description: Create a new user account.
     
      Request Body:
     
      username (string, required): User's username.
     
      password (string, required): User's password.
     
      Response:
     
      Status Code: 201 Created
     
     Response Body: User object (JSON)
   - #### User Login `/login`
     Method: POST

     Description: Authenticate a user and generate access and refresh tokens.

     Request Body:
     
     username (string, required): User's username.
     
     password (string, required): User's password.
     
     Response:
     
     Status Code: 200 OK
     
     Response Body: Access and refresh tokens (JSON)
   - #### Refresh Token `/refresh`
     Method: POST
     
     Description: Refresh the access token using the refresh token.
     
     Request Body:
     
     refresh_token (string, required): Refresh token.
     
     Response:
     
     Status Code: 200 OK
  
     Response Body: New access token (JSON)
  - #### Revoke Token `/revoke`
    Method: POST
    
    Description: Revoke the refresh token.
    
    Request Body:
    
    refresh_token (string, required): Refresh token to revoke.
    
    Response:
    
    Status Code: 204 No Content
  - #### Update User `/users`
    Method: PUT
    
    Description: Update user account information.
    
    Request Body:
    
    Updated user data (JSON)
    
    Response:
    
    Status Code: 200 OK
    
    Response Body: Updated user object (JSON)
- #### Prose Management
  - #### Create Prose `/prose`
    Method: POST
    
    Description: Create a new prose (post).
    
    Request Body:
    
    content (string, required): Text content of the prose.
    
    Response:
    
    Status Code: 201 Created
    
    Response Body: Prose object (JSON)
  - #### Get Prose `/prose`
    Method: GET
    
    Description: Get a list of prose.
    
    Response:
    
    Status Code: 200 OK
    
    Response Body: List of prose objects (JSON)
  - #### Get Prose by ID `/prose/{proseId}`
    Method: GET
    
    Description: Get details of a prose by its ID.
    
    Response:
    
    Status Code: 200 OK
    
    Response Body: Prose object (JSON)
  - #### Delete Prose `/prose/{proseId}`
    Method: DELETE
Description: Delete a prose by its ID.
Response:
Status Code: 204 No Content
Administrative Tasks
Gold Webhooks
URL: /gold/webhooks
Method: POST
Description: Receive webhook notifications related to gold.
Request Body: Webhook data (JSON)
Response:
Status Code: 200 OK
Metrics
URL: /metrics
Method: GET
Description: Retrieve server metrics.
Response:
Status Code: 200 OK
Response Body: Metrics data (JSON)
Health Check
URL: /healthz
Method: GET
Description: Check the health status of the API.
Response:
Status Code: 200 OK
Response Body: "OK" (string)

