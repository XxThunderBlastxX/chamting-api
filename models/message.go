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
	initialize  = "initialize"
	//getJson     = "getJson"
)

var (
	Cli       = make(chan Client)    // Cli is client channel for communicating with go routine
	PayLoad   = make(chan []byte)    // PayLoad channel contains the payload sent from client
	RdbClient *redis.Client          // RdbClient is redis database instance client to store client info
	RdbChat   *redis.Client          // RdbChat is redis database instance client to store chat message
	RJson     *rejson.Handler        //RJson is a redis json instance of RdbChat client
	ctx       = context.Background() // ctx is used as context passed to redis
)

// Client holds the structure of a single client instance
type Client struct {
	Id   string          `json:"id"` // client id fetched from query or if not passed then generated automatically
	Conn *websocket.Conn // websocket connection for each client
}

// Online holds array of clients which are online and connected to websocket
type Online struct {
	Topic   string    // string of topic
	Clients *[]Client // array of clients subscribed to the topic
}

// Subscription is the list of all online users
type Subscription struct {
	Topic   string   // string of topic
	Clients []string // array of clients subscribed to the topic
}

// Server holds the array of subscriptions
type Server struct {
	Online       []Online       // array of all online clients
	Subscription []Subscription // array of all the subscribers
}

// Message holds the structure of JSON message send via websocket. If Time and MessageId is not sent from frontend then it is explicitly created here at backend
type Message struct {
	Action    string    `json:"action"`              // which action to perform with the message
	Topic     string    `json:"topic"`               // topic of the message sent
	MessageId string    `json:"messageId,omitempty"` // unique message id for each message
	Msg       string    `json:"message,omitempty"`   // message string that is sent
	Time      time.Time `json:"time,omitempty"`      // time at which message is sent
	SendBy    string    `json:"sendBy,omitempty"`    // client id of the client
}

// JsonMessage is used to send message data to redis
type JsonMessage struct {
	Topic string    `json:"topic"`   // topic of the conversation
	Msg   []Message `json:"message"` // array of conversation of type Message
}

// Send is a method to write message to the websocket
func (s *Server) Send(client *Client, message string) {
	err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Println("websocket error:", err)
		}
		return // closes the connection on error
	}

}

// StoreMessage is a method to store chat message as json object to redis
func (s *Server) StoreMessage(message string, topic string, msgId string, sendBy string, sentTime time.Time) {
	topicExist := RdbChat.Exists(ctx, "message:"+topic)

	if sentTime.IsZero() {
		sentTime = time.Now()
	}
	if msgId == "" {
		msgId = strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	if topicExist.Val() != 1 {
		jsonData := JsonMessage{
			Topic: topic,
			Msg: []Message{
				{
					MessageId: msgId,
					Msg:       message,
					Time:      sentTime,
					SendBy:    sendBy,
				},
			},
		}
		if _, err := RJson.JSONSet("message:"+topic, ".", jsonData); err != nil {
			log.Println("error occurred while storing json to redis !!!")
		}
	} else {
		jsonData := Message{
			MessageId: msgId,
			Msg:       message,
			Time:      sentTime,
			SendBy:    sendBy,
		}
		if _, err := RJson.JSONArrAppend("message:"+topic, ".message", jsonData); err != nil {
			log.Println("error occurred while storing json to redis !!!")
		}
	}
}

// RemoveClient is method to remove client
func (s *Server) RemoveClient(client *Client) {
	for _, sub := range s.Online {
		for i := 0; i < len(*sub.Clients); i++ {
			if client.Id == (*sub.Clients)[i].Id {
				if i == len(*sub.Clients)-1 {
					*sub.Clients = (*sub.Clients)[:len(*sub.Clients)-1]
				} else {
					*sub.Clients = append((*sub.Clients)[:i], (*sub.Clients)[i+1:]...)
					i--
				}
			}
		}
	}
}

// ProcessMessage is the method to process the message and execute different func depending on the action given.
// action is subscribe then execute Subscribe func.
// action is unsubscribe then execute Unsubscribe func.
// action is publish then execute Publish func.
// action is initialize execute InitServer func.
func (s *Server) ProcessMessage() {
	time.Sleep(time.Millisecond * 1)
	for {
		client := <-Cli
		payload := <-PayLoad
		m := Message{}
		if err := json.Unmarshal(payload, &m); err != nil {
			s.Send(&client, "Server: Invalid payload")
		}

		switch m.Action {
		case publish:
			s.Publish(m.Topic, []byte(m.Msg), m.MessageId, m.SendBy, m.Time)
			break

		case subscribe:
			s.Subscribe(&client, m.Topic)
			break

		case unsubscribe:
			s.Unsubscribe(&client, m.Topic)
			break

		case initialize:
			s.InitServer(m.Topic)
			break

		//case getJson:
		//	s.getJson(m.Topic)
		//	break

		default:
			s.Send(&client, "Server: Action unrecognized")
			break
		}
	}
}

// Publish is a method to broadcast message to all the clients which are subscribed with the given topic
func (s *Server) Publish(topic string, message []byte, msgId string, sendBy string, sentTime time.Time) {
	var subClients []string
	var onlineClient []Client

	for _, sub := range s.Subscription {
		if sub.Topic == topic {
			subClients = append(subClients, sub.Clients...)
		}
	}

	for _, on := range s.Online {
		if on.Topic == topic {
			onlineClient = append(onlineClient, *on.Clients...)
		}
	}

	for _, online := range onlineClient {
		if online.Id == sendBy {
			s.StoreMessage(string(message), topic, msgId, sendBy, sentTime)
		}
		for _, sub := range subClients {
			if sub == online.Id {
				s.Send(&online, string(message))
			}
		}
	}

}

// Subscribe is a method to subscribe to a given topic by any client
func (s *Server) Subscribe(client *Client, topic string) {
	exist := false

	for _, sub := range s.Online {
		if sub.Topic == topic {
			exist = true
			*sub.Clients = append(*sub.Clients, *client)
			RdbClient.SAdd(ctx, topic, client.Id)
		}
	}

	if !exist {
		newClient := &[]Client{*client}

		newOnline := &Online{
			Topic:   topic,
			Clients: newClient,
		}
		s.Online = append(s.Online, *newOnline)
		RdbClient.SAdd(ctx, topic, client.Id)
	}
	s.InitServer(topic)
}

// Unsubscribe is a method to unsubscribe to a given topic by any client
func (s *Server) Unsubscribe(client *Client, topic string) {
	// Read all topics
	for _, sub := range s.Online {
		if sub.Topic == topic {
			// Read all topics' client
			for i := 0; i < len(*sub.Clients); i++ {
				if client.Id == (*sub.Clients)[i].Id {
					// If found, remove client
					if i == len(*sub.Clients)-1 {
						// if it's stored as the last element, crop the array length
						*sub.Clients = (*sub.Clients)[:len(*sub.Clients)-1]
					} else {
						// if it's stored in between elements, overwrite the element and reduce iterator to prevent out-of-bound
						*sub.Clients = append((*sub.Clients)[:i], (*sub.Clients)[i+1:]...)
						i--
					}
				}
			}
		}
	}
}

// InitServer is a method to get all the client id for the topic from db
func (s *Server) InitServer(topic string) {
	result := RdbClient.SMembers(ctx, topic)

	exist := false

	for _, sub := range s.Subscription {
		if sub.Topic == topic {
			exist = true
			return
		}
	}

	if !exist {
		newSub := &Subscription{
			Topic:   topic,
			Clients: result.Val(),
		}
		s.Subscription = append(s.Subscription, *newSub)
	}
}

//func (s *Server) getJson(topic string) {
//	res, _ := RJson.JSONGet("message:"+topic, ".")
//	buffer, _ := res.([]byte)
//	jsonMsg := JsonMessage{}
//
//	_ = json.Unmarshal(buffer, &jsonMsg)
//
//	log.Println(jsonMsg)
//}
