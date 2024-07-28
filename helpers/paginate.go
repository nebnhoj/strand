package helpers

type Paginate struct {
	Total int         `json:"total"`
	Data  interface{} `json:"data"`
}
