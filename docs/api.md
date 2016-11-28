# API

Here is the Admiral API documentation.

# Returns

In this documentation, endpoints have returns values. Here are these values descriptions:

Value | Description
----- | -----------
200 OK           | Everything worked as expected
401 Unauthorized | You are not authorized to access the given resource
404 Not Found    | The resource you are trying to get does not exist
409 Conflict     | The resource you are trying to add already exists

# Non-authenticated calls

URI | Method | Description | Body | Returns
--- | ------ | ----------- | ---- | -------
/           | GET  | Health check endpoint           |                                           | 200 OK
/events     | POST | Docker Registry events endpoint |                                           | 200 OK
/v1/version | GET  | Get admiral version             |                                           | 200 OK
/v1/user    | PUT  | Create a new user               | <pre>{"username":"", "password":""}</pre> | <ul><li>200 OK</li><li>409 Conflict</li></ul>

# Authenticated calls

Authenticated calls need the user to use HTTP basic authentication system in order to do these calls.

## Login

URI | Method | Description | Parameters | Returns
--- | ------ | ----------- | ---------- | -------
/v1/login | GET | Check if the user is able to login | | <ul><li>200 OK</li><li>401 Unauthorized</li></ul>

## Images

URI | Method | Description | Parameters | Returns
--- | ------ | ----------- | ---------- | -------
/v1/images               | GET    | Return the user's images                 |                | <ul><li>200 OK</li><li>401 Unauthorized</li></ul>
/v1/image/*image         | GET | Return the image | The image name | <ul><li>200 OK</li><li>401 Unauthorized</li><li>404 Not Found</li></ul>
/v1/image/*image         | DELETE | Remove the given image with all its tags | The image name | <ul><li>200 OK</li><li>401 Unauthorized</li><li>404 Not Found</li></ul>
/v1/image/public/*image  | PATCH  | Set the given image as public            | The image name | <ul><li>200 OK</li><li>401 Unauthorized</li><li>404 Not Found</li></ul>
/v1/image/private/*image | PATCH  | Set the given image as public            | The image name | <ul><li>200 OK</li><li>401 Unauthorized</li><li>404 Not Found</li></ul>

## Token

URI | Method | Description | Parameters | Returns
--- | ------ | ----------- | ---------- | -------
/v1/token | GET | Get a bearer token for the asked resource | | <ul><li>200 OK</li><li>401 Unauthorized</li></ul>

## Admin
URI | Method | Description | Parameters | Returns
--- | ------ | ----------- | ---------- | -------
/v1/admin/login | GET | Check if the admin is able to login | | <ul><li>200 OK</li><li>401 Unauthorized</li></ul>
/v1/admin/images               | GET    | Return all images                 |                | <ul><li>200 OK</li><li>401 Unauthorized</li></ul>
/v1/admin/image/*image         | GET | Return the image | The image name | <ul><li>200 OK</li><li>401 Unauthorized</li><li>404 Not Found</li></ul>
/v1/admin/image/*image         | DELETE | Remove the given image with all its tags | The image name | <ul><li>200 OK</li><li>401 Unauthorized</li><li>404 Not Found</li></ul>
/v1/admin/image/public/*image  | PATCH  | Set the given image as public            | The image name | <ul><li>200 OK</li><li>401 Unauthorized</li><li>404 Not Found</li></ul>
/v1/admin/image/private/*image | PATCH  | Set the given image as public            | The image name | <ul><li>200 OK</li><li>401 Unauthorized</li><li>404 Not Found</li></ul>
