package models

import (
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"log"
)

const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

type Client struct {
	Id   string
	Conn *websocket.Conn
}

type Subscription struct {
	Topic   string
	Clients *[]Client
}

type Server struct {
	Subscriptions []Subscription
}

type Message struct {
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

func (s *Server) Send(client *Client, message string) {
	err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Println("websocket error:", err)
		}
		return // Calls the deferred function, i.e. closes the connection on error
	}
}

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

func (s *Server) ProcessMessage(client Client, messageType int, payload []byte) *Server {
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

	return s
}

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

		newTopic := &Subscription{
			Topic:   topic,
			Clients: newClient,
		}

		s.Subscriptions = append(s.Subscriptions, *newTopic)
	}
}

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
