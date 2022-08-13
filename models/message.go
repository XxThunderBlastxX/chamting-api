package models

import (
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"log"
	"time"
)

const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

// channels for communicating with ProcessMessage go routine
var (
	Cli     = make(chan Client)
	PayLoad = make(chan []byte)
)

// Client holds the structure of a single client
type Client struct {
	Id   string          // client id fetched from query or if not passed then generated automatically
	Conn *websocket.Conn // websocket connection for each client
}

// Subscription holds array of clients subscribed to a topic
type Subscription struct {
	Topic   string    // string od topic
	Clients *[]Client // array of clients subscribed to the topic
}

// Server holds the array of subscriptions
type Server struct {
	Subscriptions []Subscription // array of all subscriptions
}

// Message holds the structure of JSON message send via websocket
type Message struct {
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
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

// RemoveClient is method to remove client
func (s *Server) RemoveClient(client *Client) {
	for _, sub := range s.Subscriptions {
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

// ProcessMessage is the method to process the message and execute different func depending on the action given
// action is subscribe then execute Subscribe func
// action is unsubscribe then execute Unsubscribe func
// action is publish then execute Publish func
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
			s.Publish(m.Topic, []byte(m.Message))
			break

		case subscribe:
			s.Subscribe(&client, m.Topic)
			break

		case unsubscribe:
			s.Unsubscribe(&client, m.Topic)
			break

		default:
			s.Send(&client, "Server: Action unrecognized")
			break
		}
	}
}

// Publish is a method to broadcast message to all the clients which are subscribed with the given topic
func (s *Server) Publish(topic string, message []byte) {
	var clients []Client

	for _, sub := range s.Subscriptions {
		if sub.Topic == topic {
			clients = append(clients, *sub.Clients...)
		}
	}

	for _, client := range clients {
		s.Send(&client, string(message))
	}
}

// Subscribe is a method to subscribe to a given topic by any client
func (s *Server) Subscribe(client *Client, topic string) {
	exist := false

	for _, sub := range s.Subscriptions {
		if sub.Topic == topic {
			exist = true
			*sub.Clients = append(*sub.Clients, *client)
		}
	}

	if !exist {
		newClient := &[]Client{*client}

		newSub := &Subscription{
			Topic:   topic,
			Clients: newClient,
		}

		s.Subscriptions = append(s.Subscriptions, *newSub)
	}
}

// Unsubscribe is a method to unsubscribe to a given topic by any client
func (s *Server) Unsubscribe(client *Client, topic string) {
	// Read all topics
	for _, sub := range s.Subscriptions {
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
