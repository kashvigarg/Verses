# Verse-S
A social media Backend APi

## API DOCUMENTATION

### Base URL

### Endpoints

- #### User Management
   - #### Create User `POST ``/users`
      Description: Create a new user account.
     
      Request Body:

      {
	"Name": "name"
	"Password":"password"
	"Email":"email"
}
     
      Response:
     
      Status Code: 201 Created
     
     Response Body: User object (JSON)
   - #### User Login `POST` `/login`
     Description: Authenticate a user and generate access and refresh tokens.

     Request Body:
     
     username (string, required): User's username.
     
     password (string, required): User's password.
     
     Response:
     
     Status Code: 200 OK
     
     Response Body: Access and refresh tokens (JSON)
   - #### Refresh Token `POST` `/refresh`
     Description: Refresh the access token using the refresh token.
     
     Request Body:
     
     refresh_token (string, required): Refresh token.
     
     Response:
     
     Status Code: 200 OK
  
     Response Body: New access token (JSON)
  - #### Revoke Token `POST` `/revoke`
    Description: Revoke the refresh token.
    
    Request Body:
    
    refresh_token (string, required): Refresh token to revoke.
    
    Response:
    
    Status Code: 204 No Content
  - #### Update User `PUT` `/users` 
    Description: Update user account information.
    
    Request Body:
    
    Updated user data (JSON)
    
    Response:
    
    Status Code: 200 OK
    
    Response Body: Updated user object (JSON)
- #### Prose Management
  - #### Create Prose `POST` `/prose`
    Description: Create a new prose (post).
    
    Request Body:
    
    content (string, required): Text content of the prose.
    
    Response:
    
    Status Code: 201 Created
    
    Response Body: Prose object (JSON)
  - #### Get Prose `GET` `/prose`
    Description: Get a list of prose.
    
    Response:
    
    Status Code: 200 OK
    
    Response Body: List of prose objects (JSON)
  - #### Get Prose by ID `GET` `/prose/{proseId}`
    Description: Get details of a prose by its ID.
    
    Response:
    
    Status Code: 200 OK
    
    Response Body: Prose object (JSON)
  - #### Delete Prose ` DELETE` `/prose/{proseId}`
    Description: Delete a prose by its ID.

    Response:

    Status Code: 204 No Content
- #### Administrative Tasks
  - #### Gold Webhooks `POST` `/gold/webhooks`
    Description: Receive webhook notifications related to gold.

    Request Body: Webhook data (JSON)

    Response:

    Status Code: 200 OK
  - #### Metrics `GET` `/metrics`
    Description: Retrieve server metrics.

    Response:

    Status Code: 200 OK

    Response Body: Metrics data (JSON)

  - #### Health Check `GET` `/healthz`
    Description: Check the health status of the API.

    Response:

    Status Code: 200 OK
    
    Response Body: "OK" (string)

