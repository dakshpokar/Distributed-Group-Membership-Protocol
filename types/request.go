package types

// Request is a struct that stores the request type, data, file and byte data
type Request struct {
	Req_type  string `json:"req_type"`
	Data      string `json:"data"`
	File      string `json:"file"`
	Byte_Data []byte `json:"byte_data"`
	Meta_Data string `json:"meta_data"`
}

// Response is a struct that stores the response type, data and byte data
type Response struct {
	Resp_type string `json:"resp_type"`
	Data      string `json:"data"`
	Byte_Data []byte `json:"byte_data"`
}
