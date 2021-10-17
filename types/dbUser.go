package types

type DBUser struct {
	Id         string `json:"_id" bson:"_id"`
	Name       string `json:"name"`
	Job        string `json:"job"`
	EmployedAt string `json:"employedat"`
	Parent     string `json:"parent"`
}
