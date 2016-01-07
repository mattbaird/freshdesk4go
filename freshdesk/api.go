package freshdesk

import (
	"fmt"
)

// either ApiKey or Username/Password
type FreshDeskClient struct {
	API
}

func NewClient(domain, username, password string, secure bool) FreshDeskClient {
	protocol := "http"
	if secure {
		protocol = "https"
	}
	api := NewAPI(protocol, domain, username, password)
	return FreshDeskClient{API: api}
}

/**
* POST to /contacts.json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X POST -d '{ "user": { "name":"Super Man", "email":"ram@freshdesk.com" }}' https://domain.freshdesk.com/contacts.json
* Response JSON
**/
func (client *FreshDeskClient) UserCreate(name, email string) (userResponse UserResponse, err error) {
	return UserResponse{}, nil
}

/**
* GET /contacts/[id].json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X GET https://domain.freshdesk.com/contacts/19.json
* Response JSON
**/
func (client *FreshDeskClient) UserView(id int) (userResponse UserResponse, err error) {
	requestUrl := client.BaseUrl() + fmt.Sprintf("/contacts/%v.json", id)
	err = client.DoWithResult(requestUrl, GET, &userResponse)
	return userResponse, err
}

/**
* DELETE /contacts/[id].json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X DELETE https://domain.freshdesk.com/contacts/1.json
* Response 200 OK
**/
func (client *FreshDeskClient) UserDelete(id int) (bool, error) {
	requestUrl := client.BaseUrl() + fmt.Sprintf("/contacts/%v.json", id)
	err := client.DoWithResult(requestUrl, DELETE, nil)
	return err == nil, err
}

/**
* POST to /contacts.json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X POST -d '{ "user": { "name":"Super Man", "email":"ram@freshdesk.com" }}' https://domain.freshdesk.com/contacts.json
* Response JSON
**/
func (client *FreshDeskClient) CustomerCreate(name, domains, description string) (CustomerResponse, error) {
	return CustomerResponse{}, nil
}

/**
* GET to /customers.json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X GET https://domain.freshdesk.com/customers.json
* Response JSON
**/
func (client *FreshDeskClient) CustomerList(name, domains, description string) ([]CustomerResponse, error) {
	return []CustomerResponse{}, nil
}

/**
* GET /customers/[id].json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X GET https://domain.freshdesk.com/contacts/19.json
* Response JSON
**/
func (client *FreshDeskClient) CustomerView(id int) (CustomerResponse, error) {
	return CustomerResponse{}, nil
}

/**
* DELETE /customers/[id].json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X DELETE https://domain.freshdesk.com/contacts/1.json
* Response 200 OK
**/
func (client *FreshDeskClient) CustomerDelete(id int) (bool, error) {
	requestUrl := client.BaseUrl() + fmt.Sprintf("/contacts/%v.json", id)
	err := client.DoWithResult(requestUrl, DELETE, nil)
	return err == nil, err
}
