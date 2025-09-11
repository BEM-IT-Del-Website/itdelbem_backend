package utils

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type MetadataResponse struct {
	Status   string      `json:"status"`
	Message  string      `json:"message"`
	Metadata interface{} `json:"metadata"`
	Data     interface{} `json:"data"`
}

type PaginationLinks struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}

type PaginationMetadata struct {
	CurrentPage int             `json:"current_page"`
	PerPage     int             `json:"per_page"`
	TotalItems  int             `json:"total_items"`
	TotalPages  int             `json:"total_pages"`
	Links       PaginationLinks `json:"links"`
}

func ResponseHandler(status, message string, data interface{}) Response {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return response
}

func MetadataFormatResponse(status string, message string, metadata interface{}, data interface{}) MetadataResponse {
	response := MetadataResponse{
		Status:   status,
		Message:  message,
		Metadata: metadata,
		Data:     data,
	}
	return response
}
