package permission

type Scope struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}
