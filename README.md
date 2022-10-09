# <center> 📫 Chamting App - API  </center>

> Made with Go for Chamting App

### 👨‍💻 Technical Stack

- Go
- Fiber
- Redis
- MongoDB

### 🛠 Prerequisite before running the application

> Create .env in root dir of the project and copy all the variables from .env.example and past in .env and assign the values of all the variables.

### ⚙ How to run the application

- To run the application in development mode:

```shell
go run .
```

- To make a build of the application:

```shell
go build -o api
```

- To make a build of the application for windows:

```shell
go build -o api.exe
```

### 📲 How User SignUp/SignIn

#### User Model ->

```go
type User struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	UserName  string             `json:"username" bson:"username"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
```

Request send for Signup:

```json
{
  "username": "Username",
  "email": "Email",
  "name": "Name",
  "password": "password"
}
```

Response received when SignUp:

```json
{
  "data": {
    "id": "Unique MongoId",
    "email": "Email",
    "password": "Cryptic Password",
    "username": "Username",
    "name": "Name",
    "created_at": "Time "
  },
  "error": "",
  "success": true,
  "token": "JWT Token"
}
```

Request send for SignIn:

```json
{
  "email": "Email",
  "password": "Password"
}
```

Response received when SignIn:

```json
{
  "data": {
    "id": "Unique MongoId",
    "email": "Email",
    "password": "Cryptic Password",
    "username": "Username",
    "name": "Name",
    "created_at": "Time "
  },
  "error": "",
  "success": true,
  "token": "JWT Token"
}
```

### 💻 How to connect to the websocket

- To connect to the websocket you can use `ws://localhost:3200/ws?id={user_id}`<br></br>
  user_id is the same as the id received as response when a user signin or signup.

#### Message model :

```go
// Message holds the structure of JSON message send via websocket.
// If Time and MessageId is not sent from frontend then it is explicitly created here at backend
type Message struct {
    Action    string    `json:"action,omitempty"`  // which action to perform with the message
    Topic     string    `json:"topic"`             // topic of the message sent
    MessageId string    `json:"messageId"`         // unique message id for each message
    Msg       string    `json:"message,omitempty"` // message string that is sent
    Time      time.Time `json:"time,omitempty"`    // time at which message is sent
    SendBy    string    `json:"sendBy,omitempty"`  // client id of the client
}
```

> Since the websocket messaging is developed in the Pub/Sub Pattern so every user has to subscribe to a topic to send and receive message from that topic.

### 🔔 How to Subscribe to a topic

To subscribe to a topic you need to send json object to websocket:

```json
{
  "action": "subscribe",
  "topic": "Topic"
}
```

**_Topic_** send to the server is unique. Even if you try to give same topic to **_subscribe_** again it will not add new topic to database rather override the existing topic.
To chat privately all the users have to **_subscribe_** to same topic.

### 📤 How to Send/Publish message to a topic

To send message to a topic you need to send json object to websocket:

```json
{
  "action": "publish",
  "topic": "Topic",
  "message": "Message",
  "messageId": "Message Id",
  "time": "Time",
  "sendBy": "UserId of the sender"
}
```

Note :- If you don't set messageId and time in json object it will be automatically be generated with unique and random messageId and current time in the server.

Receive a response for **_publish_** action to all the subscribed users:

```json
{
  "topic": "Topic",
  "messageId": "Unique MessageID",
  "message": "Message",
  "time": "Time",
  "sendBy": "UserId of the Sender"
}
```

After sending this json object value of **_message_** will automatically be sent to all the users subscribed to that topic. To receive message all the users subscribed also need to be online and connected to the server.

### 🧪 How to unit-test each route
```sh
go test -v ./routes
```
> Note :- Websocket unit test is not written.

## 🙍‍♂️ Author

- 👦 [ThunderBlast](https://github.com/XxThunderBlastxX)

## 📃 Licence

Copyright © 2022 [ThunderBlast](https://github.com/xXThunderBlastxX).<br />
This project is [MIT](LICENCE) licensed.
