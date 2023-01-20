package functions

import (
	"encoding/json"
	"net/http"
)

type InvokeRequest struct {
	Data     map[string]json.RawMessage `json:"data"`
	Metadata map[string]interface{}     `json:"metadata"`
}

type InvokeResponse struct {
	Outputs     map[string]interface{} `json:"outputs"`
	Logs        []string               `json:"logs"`
	ReturnValue interface{}            `json:"returnValue"`
}

// Helper to get a binding from the invocation request
func GetInvocationBinding(r http.Request, bindingName string) map[string]interface{} {
	invokeRequest := InvokeRequest{}

	d := json.NewDecoder(r.Body)
	d.Decode(&invokeRequest)

	var reqData map[string]interface{}
	json.Unmarshal(invokeRequest.Data[bindingName], &reqData)

	return reqData
}

// Helper to create a response with a single binding
func NewInvokeResponse(bindingName, msg string, body interface{}) InvokeResponse {
	outputs := make(map[string]interface{})
	outputs["message"] = msg
	outputs[bindingName] = map[string]interface{}{
		"body": body,
	}

	return InvokeResponse{
		Outputs:     outputs,
		Logs:        []string{},
		ReturnValue: "",
	}
}
