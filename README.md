# Bun microservice

Start the service

```cmd
  air
```
The service will be started on http://localhost:3002

---

## API Endpoints

### Retrieve All Messages

Retrieve all messages stored in the database.

- **Endpoint:** `GET /data`
- **Description:** This endpoint allows you to retrieve all messages that have been stored in the database.
- **Response:** Returns a JSON array containing the stored messages.
- **Example Request:**
  ```http
  GET /data
  ```
- **Example Response:**
  ```json
  {
    "success": true,
    "data": [
      {
          "ID": 1,
          "message": "Hello, world!"
      },
      {
          "ID": 2,
          "message": "How are you?"
      }
    ]
  }
  ```
---
Certainly, you can modify the POST request to send a JSON body instead of `x-www-form-urlencoded`. Here's the updated documentation for the "Store a Message" endpoint:

### Store a Message

Store a new message in the database.

- **Endpoint:** `POST /data`
- **Description:** This endpoint allows you to store a new message in the database.
- **Request Payload:** The request payload must be in JSON format with a single key "message" and a string value.
- **Example Request:**
  ```http
  POST /data
  Content-Type: application/json

  {
    "message": "This is a new message"
  }
  ```
- **Example Response:**
  ```json
  {
    "data": null,
    "success": true
  }
  ```
- **Error Response:**
  - `400 Bad Request`: If the request payload is missing the "message" key or if the value is not a string.

---
