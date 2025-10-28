package http_utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

const (
	ERROR   = "error"
	SUCCESS = "success"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type MethodHandler struct {
	get    HandlerFunc
	post   HandlerFunc
	put    HandlerFunc
	delete HandlerFunc
}

func NewMethodHandler() *MethodHandler {
	return &MethodHandler{}
}

func (methodHandler *MethodHandler) Get(handler HandlerFunc) *MethodHandler {
	methodHandler.get = handler
	return methodHandler
}

func (methodHandler *MethodHandler) Post(handler HandlerFunc) *MethodHandler {
	methodHandler.post = handler
	return methodHandler
}

func (methodHandler *MethodHandler) Put(handler HandlerFunc) *MethodHandler {
	methodHandler.put = handler
	return methodHandler
}

func (methodHandler *MethodHandler) Delete(handler HandlerFunc) *MethodHandler {
	methodHandler.delete = handler
	return methodHandler
}

func (methodHandler *MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var handler HandlerFunc

	switch r.Method {
	case GET:
		handler = methodHandler.get
	case POST:
		handler = methodHandler.post
	case PUT:
		handler = methodHandler.put
	case DELETE:
		handler = methodHandler.delete
	}

	if handler == nil {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	handler(w, r)

	log.Printf("← %s %s (%v)", r.Method, r.URL.Path, time.Since(start))
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	responseMessage := fmt.Sprintf("Error response: %s", message)
	jsonResponse(w, statusCode, map[string]string{"error": responseMessage}, ERROR)
}

func SendSuccessResponse(w http.ResponseWriter, data interface{}) {
	jsonResponse(w, http.StatusOK, data, SUCCESS)
}

func SendBadRequestResponse(w http.ResponseWriter, message string) {
	responseMessage := fmt.Sprintf("Bad Request: %s", message)
	jsonResponse(w, http.StatusBadRequest, map[string]string{"error": responseMessage}, ERROR)
}

func SendInternalServerErrorResponse(w http.ResponseWriter, message string) {
	responseMessage := fmt.Sprintf("Internal Server Error: %s", message)
	jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": responseMessage}, ERROR)
}

func SendUnauthorizedResponse(w http.ResponseWriter, message string) {
	responseMessage := fmt.Sprintf("Unauthorized: %s", message)
	jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": responseMessage}, ERROR)
}

func ParseJSONRequest(r *http.Request, dst interface{}) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		log.Printf("Failed to parse JSON request: %v", err)
		return false
	}
	return true
}

func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}, status string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := APIResponse{
		Status: status,
		Data:   data,
	}
	json.NewEncoder(w).Encode(response)
}
