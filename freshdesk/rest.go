package freshdesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var debug = false

const (
	CONTENT_TYPE_APPLICATION_JSON = "application/json"
	CONTENT_TYPE_MULTIPART_FORM   = "application/x-www-form-urlencoded"
	GET                           = "GET"
	POST                          = "POST"
	DELETE                        = "DELETE"
)

// should be used everywhere we build rest apis
type JsonEnvelope struct {
	Status        ApiStatus   `json:"status"`
	Response      interface{} `json:"response"` // depending on ApiStatus, it's either an ApiError or the return type
	TimeToProcess string      `json:"time_to_process,omitempty"`
}

func (env *JsonEnvelope) Json() string {
	b, _ := json.Marshal(env)
	return string(b)
}

type ApiStatus struct {
	Code        int    `json:"code,omitempty"`
	I18NMessage string `json:"i18n_message,omitempty"`
	Message     string `json:"message,omitempty"`
}

var NotFound = &ApiError{Err: "Not Found", Stack: "", Message: "Not Found", I18NMessage: "not.found", DeveloperMessage: "", Code: 404, MoreInfo: nil}

type ApiError struct {
	Err              string              `json:"error,omitempty"`
	Stack            string              `json:"stack,omitempty"`
	Message          string              `json:"message,omitempty"`
	I18NMessage      string              `json:"i18n_message,omitempty"`
	DeveloperMessage string              `json:"developer_message,omitempty"`
	Code             int                 `json:"code,omitempty"`
	MoreInfo         map[string][]string `json:"more_info,omitempty"`
}

func (apiError ApiError) Error() string {
	return fmt.Sprintf("%v:%s", apiError.Code, apiError.Message)
}

func Success(msg string, i18msg string, obj interface{}) JsonEnvelope {
	status := ApiStatus{Code: 0, I18NMessage: i18msg, Message: msg}
	//	bytes, _ := json.Marshal(obj)
	retval := JsonEnvelope{Status: status, Response: obj}
	return retval
}

func SimpleFailure(code int, msg string) JsonEnvelope {
	return Failure(code, msg, "", "", "", "", "", "")
}

func Failure(code int, msg string, i18msg string, err string, stack string,
	errMsg string, devMsg string, url string) JsonEnvelope {
	status := ApiStatus{Code: code, I18NMessage: i18msg, Message: msg}
	moreinfo := make(map[string][]string)
	if len(url) > 0 {
		moreinfo["url"] = []string{url}
	}
	responseError := ApiError{Err: err, Stack: stack, Message: errMsg, I18NMessage: i18msg, DeveloperMessage: devMsg, Code: code, MoreInfo: moreinfo}
	//	bytes, _ := json.Marshal(responseError)
	retval := JsonEnvelope{Status: status, Response: responseError}
	return retval
}

func Unauthorized() JsonEnvelope {
	status := ApiStatus{Code: http.StatusUnauthorized, I18NMessage: "unauthorized", Message: "Unauthorized"}
	retval := JsonEnvelope{Status: status, Response: nil}
	return retval
}

type API struct {
	Protocol  string            `json:"protocol,omitempty"`
	Domain    string            `json:"domain,omitempty"`
	Port      int               `json:"port,omitempty"`
	Transport http.RoundTripper `json:"-"`
	Username  string            `json:"-"`
	Password  string            `json:"-"`
}

func NewAPI(protocol string, domain string, username, password string) API {
	return API{Protocol: protocol, Domain: domain, Username: username, Password: password}
}

func (api *API) BaseUrl() string {
	return fmt.Sprintf("%s://%s.freshdesk.com", api.Protocol, api.Domain)
}

func (api *API) GetBody(requestUrl string, cTimeout time.Duration, rwTimeout time.Duration) (string, error) {
	client := NewTimeoutClient(cTimeout, rwTimeout)
	var req *http.Request
	var httpErr error
	req, httpErr = http.NewRequest(GET, strings.TrimSpace(requestUrl), nil)
	if httpErr != nil {
		return "", &ApiError{Err: httpErr.Error(), Stack: "", Message: httpErr.Error(), I18NMessage: "", DeveloperMessage: "", Code: 0, MoreInfo: nil}
	}
	req.SetBasicAuth(api.Username, api.Password)
	req.Close = true
	resp, httpErr := client.Do(req)
	if httpErr != nil {
		errCode := 500
		if resp != nil {
			errCode = resp.StatusCode
		}
		return "", &ApiError{Err: httpErr.Error(), Stack: "", Message: httpErr.Error(), I18NMessage: "", DeveloperMessage: "", Code: errCode, MoreInfo: nil}
	}
	defer resp.Body.Close()
	body, readBodyError := ioutil.ReadAll(resp.Body)
	if readBodyError != nil {
		return "", &ApiError{Err: readBodyError.Error(), Stack: "", Message: readBodyError.Error(), I18NMessage: "", DeveloperMessage: "", Code: resp.StatusCode, MoreInfo: nil}
	}
	if resp.StatusCode == 404 {
		return "", NotFound
	}
	if resp.StatusCode >= 300 {
		// marshal to jsonEnvelope which is what the API will return
		var jsonEnvelope JsonEnvelope
		unmarshalError := json.Unmarshal(body, &jsonEnvelope)
		if unmarshalError == nil {
			var apiError *ApiError
			b, marshalError := json.Marshal(jsonEnvelope.Response)
			if marshalError != nil {
				return "", &ApiError{Err: fmt.Sprintf("Unable to Marshal jsonEnvelope:%v", marshalError), Stack: "", Message: fmt.Sprintf("Unable to Marshal jsonEnvelope :%v", marshalError), I18NMessage: "", DeveloperMessage: "", Code: resp.StatusCode, MoreInfo: nil}
			}
			apiErrorUnmarshalError := json.Unmarshal(b, &apiError)
			if apiErrorUnmarshalError != nil {
				return "", &ApiError{Err: fmt.Sprintf("Unable to Unmarshal ApiError:%v", apiErrorUnmarshalError), Stack: "", Message: fmt.Sprintf("Unable to Unmarshal ApiError:%v", apiErrorUnmarshalError), I18NMessage: "", DeveloperMessage: "", Code: resp.StatusCode, MoreInfo: nil}
			} else {
				return "", apiError
			}
		} else {
			if debug {
				fmt.Printf("DoWithResult called url [%s] returned %s\n", requestUrl, string(body))
			}
			return "", &ApiError{Err: fmt.Sprintf("Unable to Unmarshal jsonEnvelope:%v", unmarshalError), Stack: "", Message: fmt.Sprintf("Unable to Unmarshal jsonEnvelope:%v", unmarshalError), I18NMessage: "", DeveloperMessage: "", Code: resp.StatusCode, MoreInfo: nil}
		}
	} // else unmarshall to the result type specified by caller
	return string(body), nil
}

func (api *API) DoWithResultEx(requestUrl string, method string, payload string, result interface{},
	cTimeout time.Duration, rwTimeout time.Duration, contentType string) error {
	client := NewTimeoutClient(cTimeout, rwTimeout)
	var req *http.Request
	if len(payload) > 0 {
		var httpErr error
		req, httpErr = http.NewRequest(strings.TrimSpace(method), strings.TrimSpace(requestUrl), bytes.NewBufferString(payload))
		if httpErr != nil {
			return &ApiError{Err: httpErr.Error(), Stack: "", Message: httpErr.Error(), I18NMessage: "", DeveloperMessage: "", Code: 0, MoreInfo: nil}
		}
		req.SetBasicAuth(api.Username, api.Password)
		req.Header.Add("Content-Type", contentType)
		req.Header.Add("Content-Length", strconv.Itoa(len(payload)))
	} else {
		var httpErr error
		req, httpErr = http.NewRequest(strings.TrimSpace(method), strings.TrimSpace(requestUrl), nil)
		if httpErr != nil {
			return &ApiError{Err: httpErr.Error(), Stack: "", Message: httpErr.Error(), I18NMessage: "", DeveloperMessage: "", Code: 0, MoreInfo: nil}
		}
		req.SetBasicAuth(api.Username, api.Password)
	}
	if strings.Index(requestUrl, "%") > 0 {
		fmt.Printf("******** REQUEST WITH PERCENT ENCODED PARAM FOUND ************:%v\n", requestUrl)
	}
	req.Close = true
	var httpErr error
	resp, httpErr := client.Do(req)
	if httpErr != nil {
		errCode := 500
		if resp != nil {
			errCode = resp.StatusCode
		}
		return &ApiError{Err: httpErr.Error(), Stack: "", Message: httpErr.Error(), I18NMessage: "", DeveloperMessage: "", Code: errCode, MoreInfo: nil}
	}
	defer resp.Body.Close()

	body, readBodyError := ioutil.ReadAll(resp.Body)
	if readBodyError != nil {
		return &ApiError{Err: readBodyError.Error(), Stack: "", Message: readBodyError.Error(), I18NMessage: "", DeveloperMessage: "", Code: resp.StatusCode, MoreInfo: nil}
	}

	if debug {
		fmt.Printf("URL: %s, Method: %s, HTTP Status Code: %d, Body: %s\n", requestUrl, method, resp.StatusCode, body)
	}
	if resp.StatusCode == 404 {
		return NotFound
	}
	if resp.StatusCode >= 300 {
		// marshal to jsonEnvelope which is what the API will return
		var jsonEnvelope JsonEnvelope
		unmarshalError := json.Unmarshal(body, &jsonEnvelope)
		if unmarshalError == nil {
			var apiError *ApiError
			b, marshalError := json.Marshal(jsonEnvelope.Response)
			if marshalError != nil {
				return &ApiError{Err: fmt.Sprintf("unable to unmarshal jsonEnvelope:%v", marshalError), Stack: "", Message: fmt.Sprintf("unable to unmarshal jsonEnvelope:%v", marshalError), I18NMessage: "", DeveloperMessage: "", Code: resp.StatusCode, MoreInfo: nil}
			}
			apiErrorUnmarshalError := json.Unmarshal(b, &apiError)
			if apiErrorUnmarshalError != nil {
				return &ApiError{Err: fmt.Sprintf("unable to unmarshal apiError:%v", apiErrorUnmarshalError), Stack: "", Message: fmt.Sprintf("unable to unmarshal apiError:%v", apiErrorUnmarshalError), I18NMessage: "", DeveloperMessage: "", Code: resp.StatusCode, MoreInfo: nil}
			} else {
				return apiError
			}
		} else {
			if debug {
				fmt.Printf("DoWithResult called url [%s] with method [%s] and returned %s\n", requestUrl, method, string(body))
			}
			return &ApiError{Err: fmt.Sprintf("unable to unmarshal apiError:%v", unmarshalError), Stack: "", Message: fmt.Sprintf("unable to unmarshal apiError:%v", unmarshalError), I18NMessage: "", DeveloperMessage: "", Code: resp.StatusCode, MoreInfo: nil}
		}
	} // else unmarshall to the result type specified by caller
	var jsonEnvelope JsonEnvelope
	jsonEnvelopeUnmarshalError := json.Unmarshal(body, &jsonEnvelope)
	if jsonEnvelopeUnmarshalError != nil {
		return &ApiError{Err: jsonEnvelopeUnmarshalError.Error(), Stack: "", Message: jsonEnvelopeUnmarshalError.Error(), I18NMessage: "", DeveloperMessage: "", Code: 0, MoreInfo: nil}
	}
	if result != nil {
		b, responseMashalError := json.Marshal(jsonEnvelope.Response)
		if responseMashalError != nil {
			return &ApiError{Err: responseMashalError.Error(), Stack: "", Message: responseMashalError.Error(), I18NMessage: "", DeveloperMessage: "", Code: 0, MoreInfo: nil}
		}
		resultUnmarshalError := json.Unmarshal(b, &result)
		if resultUnmarshalError != nil {
			return &ApiError{Err: resultUnmarshalError.Error(), Stack: "", Message: resultUnmarshalError.Error(), I18NMessage: "", DeveloperMessage: "", Code: 0, MoreInfo: nil}
		}
	}
	return nil
}

func (api *API) DoWithResult(requestUrl string, method string, result interface{}) (err error) {
	return api.DoWithResultEx(requestUrl, method, "", result, connectTimeOut, readWriteTimeout, CONTENT_TYPE_APPLICATION_JSON)
}

func GetPayload(reader io.Reader, obj *interface{}) ([]byte, error) {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}
	if obj != nil {
		return body, json.Unmarshal(body, obj)
	}
	return body, nil
}

type RestError struct {
	StatusCode int
	Message    string
	I18nMsg    string
	Stack      string
	DevMsg     string
	Url        string
}

func (e *RestError) Error() string {
	return e.Message
}

func NewRestError(msg string, statusCode int, url string) *RestError {
	return &RestError{Message: msg, StatusCode: statusCode, Url: url}
}
