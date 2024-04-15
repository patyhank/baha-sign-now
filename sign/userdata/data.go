package userdata

type Data struct {
	Userid     string `json:"userid"`
	Nickname   string `json:"nickname"`
	Gold       int    `json:"gold"`
	Gp         int    `json:"gp"`
	Avatar     string `json:"avatar"`
	AvatarS    string `json:"avatar_s"`
	Lv         int    `json:"lv"`
	Properties struct {
		Class1 string `json:"Class1"`
	} `json:"properties"`
}

type UserData struct {
	Data Data `json:"data"`
}
