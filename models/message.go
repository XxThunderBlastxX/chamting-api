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
	//Topic   string    // string of topic
	clients *Client // online clients
}

//// Subscription is the list of all online users
//type Subscription struct {
//	Topic   string   // string of topic
//	Clients []string // array of clients subscribed to the topic
//}

// Server holds the array of subscriptions
type Server struct {
	online []Online // array of all online clients
	//Subscription []Subscription // array of all the subscribers
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

// RemoveClient is method to remove client
//func (s *Server) RemoveClient(client *Client) {
//	for _, sub := range s.Online {
//		for i := 0; i < len(*sub.Clients); i++ {
//			if client.Id == (*sub.Clients)[i].Id {
//				if i == len(*sub.Clients)-1 {
//					*sub.Clients = (*sub.Clients)[:len(*sub.Clients)-1]
//				} else {
//					*sub.Clients = append((*sub.Clients)[:i], (*sub.Clients)[i+1:]...)
//					i--
//				}
//			}
//		}
//	}
//}

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
			s.Send(&client, m)
		}

		switch m.Action {
		case publish:
			s.Publish(m)
			break

		case subscribe:
			s.Subscribe(&client, m.Topic)
			break

		//case unsubscribe:
		//	s.Unsubscribe(&client, m.Topic)
		//	break

		//case initialize:
		//	s.InitServer(m.Topic)
		//	break

		//case getJson:
		//	s.getJson(m.Topic)
		//	break

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

	//
	//for _, sub := range s.Subscription {
	//	if sub.Topic == msg.Topic {
	//		subClients = append(subClients, sub.Clients...)
	//	}
	//}
	//
	//for _, on := range s.Online {
	//	if on.Topic == msg.Topic {
	//		onlineClient = append(onlineClient, *on.Clients...)
	//	}
	//}
	//
	//for _, online := range onlineClient {
	//	if online.Id == msg.SendBy {
	//		s.StoreMessage(msg.Msg, msg.Topic, msg.MessageId, msg.SendBy, msg.Time)
	//	}
	//	for _, sub := range subClients {
	//		if sub == online.Id {
	//			s.Send(&online, msg)
	//		}
	//	}
	//}

}

// Subscribe is a method to subscribe to a given topic by any client
func (s *Server) Subscribe(client *Client, topic string) {
	//exist := false
	RdbClient.SAdd(ctx, topic, client.Id)
	//for _, sub := range s.Online {
	//	if sub.Topic == topic {
	//		exist = true
	//		*sub.Clients = append(*sub.Clients, *client)
	//		RdbClient.SAdd(ctx, topic, client.Id)
	//	}
	//}

	//if !exist {
	//	newClient := &[]Client{*client}
	//
	//	newOnline := &Online{
	//		Topic:   topic,
	//		Clients: newClient,
	//	}
	//	s.Online = append(s.Online, *newOnline)
	//	RdbClient.SAdd(ctx, topic, client.Id)
	//}
	//s.InitServer(topic)

}

// OnlineClient is a method use to tell the server that the client is online.
func (s *Server) OnlineClient(client *Client) {
	newOnline := Online{clients: client}
	s.online = append(s.online, newOnline)
}

// Unsubscribe is a method to unsubscribe to a given topic by any client
//func (s *Server) Unsubscribe(client *Client, topic string) {
//	// Read all topics
//	for _, sub := range s.Online {
//		if sub.Topic == topic {
//			// Read all topics' client
//			for i := 0; i < len(*sub.Clients); i++ {
//				if client.Id == (*sub.Clients)[i].Id {
//					// If found, remove client
//					if i == len(*sub.Clients)-1 {
//						// if it's stored as the last element, crop the array length
//						*sub.Clients = (*sub.Clients)[:len(*sub.Clients)-1]
//					} else {
//						// if it's stored in between elements, overwrite the element and reduce iterator to prevent out-of-bound
//						*sub.Clients = append((*sub.Clients)[:i], (*sub.Clients)[i+1:]...)
//						i--
//					}
//				}
//			}
//		}
//	}
//}

//// InitServer is a method to get all the client id for the topic from db
//func (s *Server) InitServer(topic string) {
//	result := RdbClient.SMembers(ctx, topic)
//
//	exist := false
//
//	for _, sub := range s.Subscription {
//		if sub.Topic == topic {
//			exist = true
//			return
//		}
//	}
//
//	if !exist {
//		newSub := &Subscription{
//			Topic:   topic,
//			Clients: result.Val(),
//		}
//		s.Subscription = append(s.Subscription, *newSub)
//	}
//}

//func (s *Server) getJson(topic string) {
//	res, _ := RJson.JSONGet("message:"+topic, ".")
//	buffer, _ := res.([]byte)
//	jsonMsg := JsonMessage{}
//
//	_ = json.Unmarshal(buffer, &jsonMsg)
//
//	log.Println(jsonMsg)
//}
