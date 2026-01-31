//go:build tcp
// +build tcp

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type GameClient struct {
	conn         net.Conn
	serverAddr   string
	connected    bool
}

func NewGameClient(addr string) *GameClient {
	return &GameClient{
		serverAddr: addr,
		connected:  false,
	}
}

func (c *GameClient) Connect() error {
	var err error
	maxRetries := 10
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		fmt.Printf("Attempting to connect to %s (attempt %d/%d)...\n", c.serverAddr, i+1, maxRetries)
		c.conn, err = net.Dial("tcp", c.serverAddr)
		if err == nil {
			c.connected = true
			fmt.Println("Successfully connected to Python bridge server!")
			return nil
		}
		fmt.Printf("Connection failed: %v. Retrying in %v...\n", err, retryDelay)
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("failed to connect after %d attempts: %v", maxRetries, err)
}

func (c *GameClient) SendRequest(request map[string]interface{}) (map[string]interface{}, error) {
	if !c.connected {
		return nil, fmt.Errorf("not connected to server")
	}

	// Send request
	encoder := json.NewEncoder(c.conn)
	if err := encoder.Encode(request); err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	// Receive response
	decoder := json.NewDecoder(c.conn)
	var response map[string]interface{}
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("error receiving response: %v", err)
	}

	return response, nil
}

func (c *GameClient) GetGameState() (map[string]interface{}, error) {
	request := map[string]interface{}{
		"type": "get_state",
	}
	return c.SendRequest(request)
}

func (c *GameClient) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.connected = false
		fmt.Println("Disconnected from server")
	}
}

func main() {
	client := NewGameClient("127.0.0.1:9999")

	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer client.Close()

	fmt.Println("\n=== Tesselbox Game Client ===")
	fmt.Println("Connected to Python bridge server")
	fmt.Println("Fetching game state...\n")

	// Test connection by getting game state
	state, err := client.GetGameState()
	if err != nil {
		log.Printf("Error getting game state: %v", err)
	} else {
		fmt.Println("Game State:")
		jsonData, _ := json.MarshalIndent(state, "", "  ")
		fmt.Println(string(jsonData))
	}

	// Keep connection alive and periodically fetch state
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	count := 0
	for range ticker.C {
		count++
		_, err := client.GetGameState()
		if err != nil {
			log.Printf("Error getting game state: %v", err)
			break
		}

		fmt.Printf("\r[Update #%d] Connected - Game active", count)
	}

	fmt.Println("\n\nClient shutting down...")
}