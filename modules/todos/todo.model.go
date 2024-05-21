package todos

type Todo struct {
	Id      string `bson:"_id" json:"id,omitempty"`
	Name    string `json:"name"`
	Details string `json:"details"`
}
