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
func (client *FreshDeskClient) UserCreate(name, email string) (User, error) {
	var userResponse UserResponse = UserResponse{}
	userResponse.User = User{Name: name, Email: email}
	requestUrl := client.BaseUrl() + fmt.Sprintf("/contacts.json")
	err := client.DoWithResultEx(requestUrl, POST, userResponse.Json(), &userResponse, connectTimeOut, readWriteTimeout, CONTENT_TYPE_APPLICATION_JSON)
	if err != nil {
		return User{}, err
	}

	return userResponse.User, err
}

/**
* GET /contacts/[id].json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X GET https://domain.freshdesk.com/contacts/19.json
* Response JSON
**/
func (client *FreshDeskClient) UserView(id int) (User, error) {
	var userResponse UserResponse = UserResponse{}
	requestUrl := client.BaseUrl() + fmt.Sprintf("/contacts/%v.json", id)
	err := client.DoWithResult(requestUrl, GET, &userResponse)
	if err != nil {
		return User{}, err
	}
	return userResponse.User, err
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
* POST to /customer.json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X POST -d '{ "user": { "name":"Super Man", "email":"ram@freshdesk.com" }}' https://domain.freshdesk.com/contacts.json
* Response JSON
**/
func (client *FreshDeskClient) CustomerCreate(name, domains, description string) (Customer, error) {
	var customerResponse CustomerResponse = CustomerResponse{}
	customerResponse.Customer = Customer{Name: name, Domains: domains, Description: description}
	requestUrl := client.BaseUrl() + fmt.Sprintf("/customer.json")
	err := client.DoWithResultEx(requestUrl, POST, customerResponse.Json(), &customerResponse, connectTimeOut, readWriteTimeout, CONTENT_TYPE_APPLICATION_JSON)
	if err != nil {
		return Customer{}, err
	}
	return customerResponse.Customer, err
}

/**
* GET to /customers.json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X GET https://domain.freshdesk.com/customers.json
* Response JSON
**/
func (client *FreshDeskClient) CustomerList(filter string) ([]Customer, error) {
	var customerResponses []CustomerResponse
	requestUrl := client.BaseUrl() + fmt.Sprintf("/customers.json")
	if filter != "" {
		requestUrl = requestUrl + fmt.Sprintf("?letter=%s", filter)
	}
	err := client.DoWithResult(requestUrl, GET, &customerResponses)
	var customers []Customer
	if err != nil {
		return customers, err
	}
	for _, response := range customerResponses {
		customers = append(customers, response.Customer)
	}
	return customers, err

}

/**
* GET /customers/[id].json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X GET https://domain.freshdesk.com/contacts/19.json
* Response JSON
**/
func (client *FreshDeskClient) CustomerView(id int) (Customer, error) {
	var customerResponses CustomerResponse
	requestUrl := client.BaseUrl() + fmt.Sprintf("/customers/%v.json", id)
	err := client.DoWithResult(requestUrl, GET, &customerResponses)
	if err != nil {
		return Customer{}, err
	}
	return customerResponses.Customer, err
}

/**
* DELETE /customers/[id].json
* curl -v -u user@yourcompany.com:test -H "Content-Type: application/json" -X DELETE https://domain.freshdesk.com/contacts/1.json
* Response 200 OK
**/
func (client *FreshDeskClient) CustomerDelete(id int) (bool, error) {
	requestUrl := client.BaseUrl() + fmt.Sprintf("/customers/%v.json", id)
	err := client.DoWithResult(requestUrl, DELETE, nil)
	return err == nil, err
}
