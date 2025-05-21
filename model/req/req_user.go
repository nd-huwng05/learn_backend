package req

type ReqUpdateUser struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}
