package models

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/nitishm/go-rejson/v4"
	"log"
	"strings"
	"time"
)

// Used for accessing the action sent by payload
const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

var (
	Cli       = make(chan Client)    // Cli is client channel for communicating with go routine
	PayLoad   = make(chan []byte)    // PayLoad channel contains the payload sent from client
	RdbClient *redis.Client          // RdbClient is redis database instance client to store client info
	RdbChat   *redis.Client          // RdbChat is redis database instance client to store chat message
	RJson     *rejson.Handler        // RJson is a redis json instance of RdbChat client
	ctx       = context.Background() // ctx is used as context passed to redis
)

// Client holds the structure of a single client instance
type Client struct {
	Id   string          `json:"id"` // client id fetched from query or if not passed then generated automatically
	Conn *websocket.Conn // websocket connection for each client
}

// Online holds array of clients which are online and connected to websocket
type Online struct {
	clients *Client // online clients
}

// Server holds the array of subscriptions
type Server struct {
	online []Online // array of all online clients
}

// Message holds the structure of JSON message send via websocket. If Time and MessageId is not sent from frontend then it is explicitly created here at backend
type Message struct {
	Action    string    `json:"action,omitempty"`  // which action to perform with the message
	Topic     string    `json:"topic,omitempty"`   // topic of the message sent
	MessageId string    `json:"messageId"`         // unique message id for each message
	Msg       string    `json:"message,omitempty"` // message string that is sent
	Time      time.Time `json:"time,omitempty"`    // time at which message is sent
	SendBy    string    `json:"sendBy,omitempty"`  // client id of the client
}

// JsonMessage is used to send message data to redis
type JsonMessage struct {
	Topic string    `json:"topic"`   // topic of the conversation
	Msg   []Message `json:"message"` // array of conversation of type Message
}

// Send is a method to write message to the websocket
func (s *Server) Send(client *Client, msg Message) {
	msgData := Message{
		Topic:     msg.Topic,
		MessageId: msg.MessageId,
		Msg:       msg.Msg,
		Time:      msg.Time,
		SendBy:    msg.SendBy,
	}
	jsonData, _ := json.Marshal(msgData)

	err := client.Conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Println("websocket error:", err)
		}
		return
	}
}

func (s *Server) SendErr(client *Client, msg string) {
	if err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Println("websocket error:", err)
		}
		return
	}
}

// StoreMessage is a method to store chat message as json object to redis
func (s *Server) StoreMessage(msg Message) {
	// Check if the topic already exist.
	//
	// If topicExist.Val() == 1 then topic exist else does not exist.
	topicExist := RdbChat.Exists(ctx, "message:"+msg.Topic)

	if topicExist.Val() != 1 {
		jsonData := JsonMessage{
			Topic: msg.Topic,
			Msg: []Message{
				{
					MessageId: msg.MessageId,
					Msg:       msg.Msg,
					Time:      msg.Time,
					SendBy:    msg.SendBy,
				},
			},
		}
		if _, err := RJson.JSONSet("message:"+msg.Topic, ".", jsonData); err != nil {
			log.Println("error occurred while storing json to redis !!!")
		}
	} else {
		jsonData := Message{
			MessageId: msg.MessageId,
			Msg:       msg.Msg,
			Time:      msg.Time,
			SendBy:    msg.SendBy,
		}
		if _, err := RJson.JSONArrAppend("message:"+msg.Topic, ".message", jsonData); err != nil {
			log.Println("Error occurred while storing message to redis !!!")
		}
	}
}

// ProcessMessage is the method to process the message and execute different func depending on the action given.
// action is subscribe then execute Subscribe func.
// action is unsubscribe then execute Unsubscribe func.
// action is publish then execute Publish func.
func (s *Server) ProcessMessage() {
	time.Sleep(time.Millisecond * 1)
	for {
		client := <-Cli
		payload := <-PayLoad
		m := Message{}
		if err := json.Unmarshal(payload, &m); err != nil {
			s.Send(&client, m)
		}

		switch m.Action {
		case publish:
			s.Publish(m)
			break

		case subscribe:
			s.Subscribe(&client, m.Topic)
			break

		case unsubscribe:
			s.Unsubscribe(&client, m.Topic)
			break

		default:
			s.SendErr(&client, "Wrong action passed")
			break
		}
	}
}

// Publish is a method to broadcast message to all the clients which are subscribed with the given topic
func (s *Server) Publish(msg Message) {
	// Just a temp var to check if the message was sent
	var c int32 = 0

	// Checks if time is provided.
	//
	// If not provided then new current time is generated.
	if msg.Time.IsZero() {
		msg.Time = time.Now()
	}

	// Checks if messageId is provided.
	//
	// If messageId is not provided then new messageId is generated.
	if msg.MessageId == "" {
		msg.MessageId = strings.Replace(uuid.New().String(), "-", "", -1)
	}

	var subClients []string
	//var onlineClient Online

	// Gets all the subscribers.
	subClients, _ = RdbClient.SMembers(ctx, msg.Topic).Result()

	// Sends to all the online subscribers
	for _, sub := range subClients {
		for _, on := range s.online {
			if sub == on.clients.Id {
				s.Send(on.clients, msg)
				c++
			}
		}
	}

	// Store message to redis
	if c > 0 {
		s.StoreMessage(msg)
	}
}

// Subscribe is a method to subscribe to a given topic by any client
func (s *Server) Subscribe(client *Client, topic string) {
	RdbClient.SAdd(ctx, topic, client.Id)
}

// Unsubscribe is a method to unsubscribe to a given topic by any client
func (s *Server) Unsubscribe(client *Client, topic string) {
	RdbClient.SRem(ctx, topic, client.Id)
}

// OnlineClient is a method use to tell the server that the client is online.
func (s *Server) OnlineClient(client *Client) {
	newOnline := Online{clients: client}
	s.online = append(s.online, newOnline)
}
