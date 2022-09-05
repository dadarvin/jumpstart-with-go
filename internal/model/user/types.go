package user

// User struct for user data
type User struct {
	Id       int    `json:"id" form:"id"`
	UserName string `json:"username" form:"username"`
	NickName string `json:"nickname" form:"nickname"`
	Password string `json:"password" form:"password"`
}

// UserPicture struct for containing basecode64 from client to convert into base64 format
type UserPicture struct {
	Id             int    `json:"id"`
	UserName       string `json:"username"`
	ProfilePicture string `json:"profilepicture"`
}
