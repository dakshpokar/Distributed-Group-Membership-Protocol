package utils

import (
	"MP2/types"
	"encoding/json"
)

/*
*

	RequestToJSONBytes() converts a Request struct to a JSON byte array

	Parameters:
		message: a Request struct

	Returns:
		[]byte: a JSON byte array

*
*/
func RequestToJSONBytes(message types.Request) []byte {
	data, _ := json.Marshal(message)
	data = append(data, '\r')
	return data
}

/*
*

	JSONBytesToRequest() converts a JSON byte array to a Request struct

	Parameters:
		data: a JSON byte array

	Returns:
		Request: a Request struct

*
*/
func JSONBytesToRequest(data []byte) types.Request {
	var req types.Request
	json.Unmarshal(data, &req)
	return req
}

/*
*

	ResponseToJSONBytes() converts a Response struct to a JSON byte array

	Parameters:
		message: a Response struct

	Returns:
		[]byte: a JSON byte array

*
*/
func ResponseToJSONBytes(message types.Response) []byte {
	data, _ := json.Marshal(message)
	data = append(data, '\r')
	return data
}

/*
*

	JSONBytesToResponse() converts a JSON byte array to a Response struct

	Parameters:
		data: a JSON byte array

	Returns:
		Response: a Response struct

*
*/
func JSONBytesToResponse(data []byte) types.Response {
	var resp types.Response
	json.Unmarshal(data, &resp)
	return resp
}
