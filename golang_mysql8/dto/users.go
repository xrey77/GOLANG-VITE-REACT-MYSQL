package dto

type Users struct {
	Id          string  `json:"id"`
	Firstname   string  `json:"firstname"`
	Lastname    string  `json:"lastname"`
	Email       string  `json:"email"`
	Mobile      string  `json:"mobile"`
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	Roles       string  `json:"roles"`
	Isactivated int     `json:"isactivated"`
	Isblocked   int     `json:"isblocked"`
	Userpicture string  `json:"userpicture"`
	Mailtoken   int     `json:"mailtoken"`
	Qrcodeurl   *string `json:"qrcodeurl"`
	Secret      *string `json:"secret"`
}
