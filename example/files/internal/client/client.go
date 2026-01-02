package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client manages communication with the Petstore API
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

// NewClient creates a new Petstore API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
		APIKey: apiKey,
	}
}

// Category represents a pet category
type Category struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Tag represents a pet tag
type Tag struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Pet represents a pet in the store
type Pet struct {
	ID        int64     `json:"id,omitempty"`
	Name      string    `json:"name"`
	Category  *Category `json:"category,omitempty"`
	PhotoURLs []string  `json:"photoUrls"`
	Tags      []Tag     `json:"tags,omitempty"`
	Status    string    `json:"status,omitempty"`
}

// User represents a user
type User struct {
	ID         int64  `json:"id,omitempty"`
	Username   string `json:"username,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Phone      string `json:"phone,omitempty"`
	UserStatus int32  `json:"userStatus,omitempty"`
}

// Order represents a store order
type Order struct {
	ID       int64     `json:"id,omitempty"`
	PetID    int64     `json:"petId,omitempty"`
	Quantity int32     `json:"quantity,omitempty"`
	ShipDate time.Time `json:"shipDate,omitempty"`
	Status   string    `json:"status,omitempty"`
	Complete bool      `json:"complete,omitempty"`
}

// doRequest performs an HTTP request with proper headers
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("api_key", c.APIKey)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.Errorf("API request failed with status %d: %s", res.StatusCode, string(body))
	}

	return body, nil
}

// CreatePet creates a new pet
func (c *Client) CreatePet(pet *Pet) (*Pet, error) {
	data, err := json.Marshal(pet)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/pet", c.BaseURL), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result Pet
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetPet retrieves a pet by ID
func (c *Client) GetPet(id int64) (*Pet, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/pet/%d", c.BaseURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var pet Pet
	err = json.Unmarshal(body, &pet)
	if err != nil {
		return nil, err
	}

	return &pet, nil
}

// UpdatePet updates an existing pet
func (c *Client) UpdatePet(pet *Pet) (*Pet, error) {
	data, err := json.Marshal(pet)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/pet", c.BaseURL), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result Pet
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// DeletePet deletes a pet
func (c *Client) DeletePet(id int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/pet/%d", c.BaseURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

// CreateUser creates a new user
func (c *Client) CreateUser(user *User) (*User, error) {
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user", c.BaseURL), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result User
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetUser retrieves a user by username
func (c *Client) GetUser(username string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/%s", c.BaseURL, username), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user
func (c *Client) UpdateUser(username string, user *User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/user/%s", c.BaseURL, username), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

// DeleteUser deletes a user
func (c *Client) DeleteUser(username string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/user/%s", c.BaseURL, username), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

// CreateOrder creates a new store order
func (c *Client) CreateOrder(order *Order) (*Order, error) {
	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/store/order", c.BaseURL), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var result Order
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOrder retrieves an order by ID
func (c *Client) GetOrder(id int64) (*Order, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/store/order/%d", c.BaseURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var order Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// DeleteOrder deletes an order
func (c *Client) DeleteOrder(id int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/store/order/%d", c.BaseURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}
