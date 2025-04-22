// campay.go

package campay

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "fmt"
	"io"
	"net/http"
	"os"
)

// Define structure of the json response I expect from Campay after making a request
type PaymentResponse struct {
	Reference string `json:"reference"`
	UssdCode  string `json:"ussd_code"`
}

type CheckStatus struct {
	Reference string `json:"reference"`
	Status string `json:"status"`
}

// This function takes the user's mobile number, how much to charge them,
// a unique identifier for the transaction, and a description
// This function returns a pointer to paymentResponse, and an error if anything happens.

func SendPaymentRequest(apikey string, from string, amount string, reference string, description string) (*PaymentResponse, error) {
	from ="237" + from
	payload := map[string]interface{}{
		"from":               from,
		"amount":             amount,
		"description":        description,
		"external_reference": reference,
	}
	// Converts what campay expects and receives to JSON
	jsonData, err := json.Marshal(payload) 
	if err != nil {
		return nil, err
	}

	// Create a new request
	req, err := http.NewRequest("POST", "https://demo.campay.net/api/collect/", bytes.NewBuffer(jsonData))
	fmt.Println(req)
	if err != nil {
		return nil, err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+os.Getenv("API_KEY"))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read the response
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Parse the JSON response
	// Takes the JSON body and fills the PaymentResponse struct with te reference and the ussd_code

	var pr PaymentResponse
	err = json.Unmarshal(body, &pr)
	if err != nil {
		return nil, err
	}

	// Returns the pointer to the parsed response and no error if successful
	return &pr, nil
}


// Function to check transaction status
func GetStatus(apiKey string, reference string) CheckStatus{
	client := &http.Client{}
	url1 := fmt.Sprintf("https://demo.campay.net/api/transaction/%s/", reference) // Insert reference here

	req1, err := http.NewRequest("GET", url1, nil)
	
	if err != nil {
		fmt.Println("Error creating GET request:", err)
		// return
	}

	req1.Header.Set("Authorization", "Token "+apiKey)
	req1.Header.Set("Content-Type", "application/json")
	

	resp1, err := client.Do(req1)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		// return
	}
	
	defer resp1.Body.Close()

	// Print the body response for degugging
	fmt.Println("Raw Status Response:", resp1)

	var status CheckStatus
	json.NewDecoder(resp1.Body).Decode(&status)

	return status

}