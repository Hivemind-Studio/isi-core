package response

type WebResponse struct {
	Status     int         `json:"-"`
	Message    string      `json:"message"`
	Error      interface{} `json:"-"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination,omitempty"`
}
