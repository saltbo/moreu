package bind

type QueryUser struct {
	QueryPage
	Name string `form:"name"`
}

type BodyUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BodyUserPatch struct {
	Token    string `json:"token"`
	Password string `json:"password"`
	Enabled  bool   `json:"enabled"`
}
