package entity

// Root 返回的结果
type Root struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
