![ci](https://github.com/jaydee029/Verses/actions/workflows/ci.yml/badge.svg)
![cd](https://github.com/jaydee029/Verses/actions/workflows/cd.yml/badge.svg)
# Verse-S
A social media Backend API which follows a RESTful architecture, it has the following features:-
- `User signup/login`: User can Signup using an email and a password, a refresh token and an authentication token is returned, which then can be used for login, an authrntication token expires in one hour ie one session is one hour long, another authentication token can be generated using the refresh token whiich lasts 60 days, the refresh token can be revoked upon users request.
- `prose creation and profane`: users can create and access posts in the form of `prose`, it also provide features like `profane` which allows certain words to be muted.Prose can be deleted as well. 
- `Data Storage`: passwords are stored using JWT based encryption, users are given a unique uuid for identification. Data is stored on a production ready database ie postgreSQL.
- `admin tasks`: allows viewing the heath and basic metrics for the api using REST endpoints.
- `Deployment`: The backend is deployed on an AWS EC2 instance , the postgreSQL database is deployed using AWS RDS instance.
- `CI/CD`: A continuous integration and continuous delivery/deployment is set up using github actions, the ci workflow checks for the linting, formatting and security while the cd workflow, builds the code, dockerizes it, creates database migrations and finally deploys the image on the ec2 instance.

## API DOCUMENTATION

### Welcome Page
```
http://13.201.15.193/app
```

### Base URL
```
http://13.201.15.193/api
```

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
  - #### Health Check `GET` `/healthz`
    Description: Check the health status of the API.
    
- #### Administrative Tasks
    #### Base URL
    ```
    http://13.201.15.193/admin
    ```

  - #### Metrics `GET` `/metrics`
    Description: Retrieve server metrics.

