package main
//
//import (
//	"BookAPI/internal"
//
//	"bytes"
//	"encoding/json"
//	"net/http"
//	"testing"
//
//	"go.mongodb.org/mongo-driver/bson/primitive"
//)
//
//const BASE_PATH = "http://localhost:8080"
//
//var client = new(http.Client)
//
//func GenerateUniqueId(books internal.Books) string {
//	id := primitive.NewObjectID()
//
//	for _, book := range books.Books {
//		if book.Id == id {
//			return GenerateUniqueId(books)
//		}
//	}
//
//	return id.Hex()
//}
//
//func TestHealthCheck_Expect200OnSuccess(t *testing.T) {
//	req, err := http.NewRequest(http.MethodGet, BASE_PATH + "/health", nil)
//	if err != nil {
//		t.Fatalf("Error creating GET request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing GET request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected status code 200, instead got: %v", resp.StatusCode)
//	}
//}
//
//func TestPost_Expect400OnBadInput(t *testing.T) {
//	req, err := http.NewRequest(http.MethodPost, BASE_PATH, bytes.NewBuffer([]byte("Bad Input")))
//	if err != nil {
//		t.Fatalf("Error creating POST request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing POST request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusBadRequest {
//		t.Errorf("Expected status code 400, instead got: %v", resp.StatusCode)
//	}
//}
//
//func TestPost_Expect201OnSuccess(t *testing.T) {
//	b := internal.Book{BookId: -1, Title: "Test Book", Author: "Tester", Year: -1}
//
//	s, err := json.Marshal(b)
//	if err != nil {
//		t.Fatalf("Error marshalling book: %v", err)
//	}
//
//	req, err := http.NewRequest(http.MethodPost, BASE_PATH, bytes.NewBuffer(s))
//	if err != nil {
//		t.Fatalf("Error creating POST request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing POST request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusCreated {
//		t.Errorf("Expected status code 201, instead got: %v", resp.StatusCode)
//	}
//
//	if resp.Header.Get("Location") == "" {
//		t.Error("No Location header found")
//	}
//}
//
//func TestGet_Expect400OnBadBookId(t *testing.T) {
//	req, err := http.NewRequest(http.MethodGet, BASE_PATH + "/?bookId=NaN", nil)
//	if err != nil {
//		t.Fatalf("Error creating GET request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing GET request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusBadRequest {
//		t.Errorf("Expected status code 400, instead got: %v", resp.StatusCode)
//	}
//}
//
//func TestGet_Expect400OnBadYear(t *testing.T) {
//	req, err := http.NewRequest(http.MethodGet, BASE_PATH + "/?year=NaN", nil)
//	if err != nil {
//		t.Fatalf("Error creating GET request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing GET request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusBadRequest {
//		t.Errorf("Expected status code 400, instead got: %v", resp.StatusCode)
//	}
//}
//
//func TestGet_Expect200OnSuccess(t *testing.T) {
//	req, err := http.NewRequest(http.MethodGet, BASE_PATH, nil)
//	if err != nil {
//		t.Fatalf("Error creating GET request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing GET request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected status code 200, instead got: %v", resp.StatusCode)
//	}
//}
//
//func TestGetById_Returns404OnNotFound(t *testing.T) {
//	req, err := http.NewRequest(http.MethodGet, BASE_PATH, nil)
//	if err != nil {
//		t.Fatalf("Error creating GET request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing GET request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Fatalf("Expected status code 200 from GET request, instead got: %v", resp.StatusCode)
//	}
//
//	var body internal.Books
//
//	err = json.NewDecoder(resp.Body).Decode(&body)
//	if err != nil {
//		t.Fatalf("Error decoding GET response body: %v", err)
//	}
//
//	id := GenerateUniqueId(body)
//
//	req, err = http.NewRequest(http.MethodGet, BASE_PATH + "/"+ id, nil)
//	if err != nil {
//		t.Fatalf("Error creating GET request: %v", err)
//	}
//
//	resp, err = client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing GET request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusNotFound {
//		t.Errorf("Expected status code 404, instead got: %v", resp.StatusCode)
//	}
//}
//
//func TestGetById_Returns200OnSuccess(t *testing.T) {
//	b := internal.Book{BookId: -1, Title: "Test Book", Author: "Tester", Year: -1}
//
//	s, err := json.Marshal(b)
//	if err != nil {
//		t.Fatalf("Error marshalling book: %v", err)
//	}
//
//	req, err := http.NewRequest(http.MethodPost, BASE_PATH, bytes.NewBuffer(s))
//	if err != nil {
//		t.Fatalf("Error creating POST request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing POST request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusCreated {
//		t.Fatalf("Expected status code 201 from POST request, instead got: %v", resp.StatusCode)
//	}
//
//	id := resp.Header.Get("Location")
//
//	req, err = http.NewRequest(http.MethodGet, BASE_PATH + id, nil)
//	if err != nil {
//		t.Fatalf("Error creating GET request: %v", err)
//	}
//
//	resp, err = client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing GET request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected status code 200, instead got: %v", resp.StatusCode)
//	}
//}
//
//func TestPut_Expect400OnBadInput(t *testing.T) {
//	req, err := http.NewRequest(http.MethodPut, BASE_PATH + "/test", bytes.NewBuffer([]byte("Bad Input")))
//	if err != nil {
//		t.Fatalf("Error creating PUT request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing PUT request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusBadRequest {
//		t.Errorf("Expected status code 400, instead got: %v", resp.StatusCode)
//	}
//}
//
//func TestPut_Expect200OnUpdateSuccess(t *testing.T) {
//	b := internal.Book{BookId: -1, Title: "Test Book", Author: "Tester", Year: -1}
//
//	s, err := json.Marshal(b)
//	if err != nil {
//		t.Fatalf("Error marshalling book: %v", err)
//	}
//
//	req, err := http.NewRequest(http.MethodPost, BASE_PATH, bytes.NewBuffer(s))
//	if err != nil {
//		t.Fatalf("Error creating POST request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing POST request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusCreated {
//		t.Fatalf("Expected status code 201 from POST request, instead got: %v", resp.StatusCode)
//	}
//
//	id := resp.Header.Get("Location")
//
//	b = internal.Book{BookId: -2, Title: "Replacement Book", Author: "Tester", Year: -1}
//
//	s, err = json.Marshal(b)
//	if err != nil {
//		t.Fatalf("Error marshalling book: %v", err)
//	}
//
//	req, err = http.NewRequest(http.MethodPut, BASE_PATH + id, bytes.NewBuffer(s))
//	if err != nil {
//		t.Fatalf("Error creating PUT request: %v", err)
//	}
//
//	resp, err = client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing PUT request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected status code 200, instead got: %v", resp.StatusCode)
//	}
//
//	if resp.Header.Get("Location") == "" {
//		t.Error("No Location header found")
//	}
//}
//
//func TestPut_Expect201OnCreateSuccess(t *testing.T) {
//	req, err := http.NewRequest(http.MethodGet, BASE_PATH, nil)
//	if err != nil {
//		t.Fatalf("Error creating GET request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing GET request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Fatalf("Expected status code 200 from GET request, instead got: %v", resp.StatusCode)
//	}
//
//	var body internal.Books
//
//	err = json.NewDecoder(resp.Body).Decode(&body)
//	if err != nil {
//		t.Fatalf("Error decoding GET response body: %v", err)
//	}
//
//	id := GenerateUniqueId(body)
//
//	b := internal.Book{BookId: -1, Title: "Test Book", Author: "Tester", Year: -1}
//
//	s, err := json.Marshal(b)
//	if err != nil {
//		t.Fatalf("Error marshalling book: %v", err)
//	}
//
//	req, err = http.NewRequest(http.MethodPut, BASE_PATH + "/" + id, bytes.NewBuffer(s))
//	if err != nil {
//		t.Fatalf("Error creating PUT request: %v", err)
//	}
//
//	resp, err = client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing PUT request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusCreated {
//		t.Errorf("Expected status code 201, instead got: %v", resp.StatusCode)
//	}
//
//	if resp.Header.Get("Location") == "" {
//		t.Error("No Location header found")
//	}
//}
//
//func TestPatch_Expect400OnBadInput(t *testing.T) {
//	req, err := http.NewRequest(http.MethodGet, BASE_PATH, nil)
//	if err != nil {
//		t.Fatalf("Error creating GET request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing GET request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Fatalf("Expected status code 200 from GET request, instead got: %v", resp.StatusCode)
//	}
//
//	var body internal.Books
//
//	err = json.NewDecoder(resp.Body).Decode(&body)
//	if err != nil {
//		t.Fatalf("Error decoding GET response body: %v", err)
//	}
//
//	id := GenerateUniqueId(body)
//
//	req, err = http.NewRequest(http.MethodPatch, BASE_PATH + "/" + id, bytes.NewBuffer([]byte("Bad Input")))
//	if err != nil {
//		t.Fatalf("Error creating PATCH request: %v", err)
//	}
//
//	resp, err = client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing PATCH request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusBadRequest {
//		t.Errorf("Expected status code 400, instead got: %v", resp.StatusCode)
//	}
//}
//
//func TestPatch_Expect404OnNotFound(t *testing.T) {
//	req, err := http.NewRequest(http.MethodGet, BASE_PATH, nil)
//	if err != nil {
//		t.Fatalf("Error creating GET request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing GET request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Fatalf("Expected status code 200 from GET request, instead got: %v", resp.StatusCode)
//	}
//
//	var body internal.Books
//
//	err = json.NewDecoder(resp.Body).Decode(&body)
//	if err != nil {
//		t.Fatalf("Error decoding GET response body: %v", err)
//	}
//
//	id := GenerateUniqueId(body)
//
//	b := internal.Book{BookId: -1, Title: "Test Book", Author: "Tester", Year: -1}
//
//	s, err := json.Marshal(b)
//	if err != nil {
//		t.Fatalf("Error marshalling book: %v", err)
//	}
//
//	req, err = http.NewRequest(http.MethodPatch, BASE_PATH + "/" + id, bytes.NewBuffer(s))
//	if err != nil {
//		t.Fatalf("Error creating PATCH request: %v", err)
//	}
//
//	resp, err = client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing PATCH request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusNotFound {
//		t.Errorf("Expected status code 404, instead got: %v", resp.StatusCode)
//	}
//}
//
//func TestPatch_Expect200OnSuccess(t *testing.T) {
//	b := internal.Book{BookId: -1, Title: "Test Book", Author: "Tester", Year: -1}
//
//	s, err := json.Marshal(b)
//	if err != nil {
//		t.Fatalf("Error marshalling book: %v", err)
//	}
//
//	req, err := http.NewRequest(http.MethodPost, BASE_PATH, bytes.NewBuffer(s))
//	if err != nil {
//		t.Fatalf("Error creating POST request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing POST request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusCreated {
//		t.Fatalf("Expected status code 201 from POST request, instead got: %v", resp.StatusCode)
//	}
//
//	id := resp.Header.Get("Location")
//
//	u := make(map[string]interface{})
//	u["bookId"] = -2
//
//	s, err = json.Marshal(u)
//	if err != nil {
//		t.Fatalf("Error marshalling update: %v", err)
//	}
//
//	req, err = http.NewRequest(http.MethodPatch, BASE_PATH + id, bytes.NewBuffer(s))
//	if err != nil {
//		t.Fatalf("Error creating PATCH request: %v", err)
//	}
//
//	resp, err = client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing PATCH request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected status code 200, instead got: %v", resp.StatusCode)
//	}
//
//	if resp.Header.Get("Location") == "" {
//		t.Error("No Location header found")
//	}
//}
//
//func TestDelete_Expect404OnNotFound(t *testing.T) {
//	req, err := http.NewRequest(http.MethodGet, BASE_PATH, nil)
//	if err != nil {
//		t.Fatalf("Error creating GET request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing GET request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Fatalf("Expected status code 200 from GET request, instead got: %v", resp.StatusCode)
//	}
//
//	var body internal.Books
//
//	err = json.NewDecoder(resp.Body).Decode(&body)
//	if err != nil {
//		t.Fatalf("Error decoding GET response body: %v", err)
//	}
//
//	id := GenerateUniqueId(body)
//
//	req, err = http.NewRequest(http.MethodDelete, BASE_PATH + "/" + id, nil)
//	if err != nil {
//		t.Fatalf("Error creating DELETE request: %v", err)
//	}
//
//	resp, err = client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing DELETE request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusNoContent {
//		t.Errorf("Expected status code 404, instead got: %v", resp.StatusCode)
//	}
//}
//
//func TestDelete_Expect200OnSuccess(t *testing.T) {
//	b := internal.Book{BookId: -1, Title: "Test Book", Author: "Tester", Year: -1}
//
//	s, err := json.Marshal(b)
//	if err != nil {
//		t.Fatalf("Error marshalling book: %v", err)
//	}
//
//	req, err := http.NewRequest(http.MethodPost, BASE_PATH, bytes.NewBuffer(s))
//	if err != nil {
//		t.Fatalf("Error creating POST request: %v", err)
//	}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing POST request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusCreated {
//		t.Fatalf("Expected status code 201 from POST request, instead got: %v", resp.StatusCode)
//	}
//
//	id := resp.Header.Get("Location")
//
//	req, err = http.NewRequest(http.MethodDelete, BASE_PATH + id, nil)
//	if err != nil {
//		t.Fatalf("Error creating DELETE request: %v", err)
//	}
//
//	resp, err = client.Do(req)
//	if err != nil {
//		t.Fatalf("Error performing DELETE request: %v", err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected status code 200, instead got: %v", resp.StatusCode)
//	}
//}
