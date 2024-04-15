package chatlist

type ChatSubData struct {
	Desc         string `json:"desc"`
	Creator      string `json:"creator"`
	CreatorNick  string `json:"creatorNick"`
	Link         string `json:"link,omitempty"`
	TradeType    string `json:"tradeType,omitempty"`
	Time         string `json:"time,omitempty"`
	Room         string `json:"room,omitempty"`
	RoomPassword string `json:"roomPassword,omitempty"`
}
type ChatLastMsg struct {
	StanzaId int64  `json:"stanzaId"`
	Text     string `json:"text"`
}

type ChatRoom struct {
	Jid             string      `json:"jid"`
	Id              string      `json:"id"`
	Uid             string      `json:"uid,omitempty"`
	Flag            string      `json:"flag,omitempty"`
	Type            int         `json:"type"`
	PrimaryData     interface{} `json:"primaryData,omitempty"`
	Name            string      `json:"name"`
	Avatar          string      `json:"avatar"`
	LastMessageTime int64       `json:"lastMessageTime"`
	SubCount        int         `json:"subCount,omitempty"`
	JoinCount       int         `json:"joinCount,omitempty"`
	SubType         int         `json:"subType,omitempty"`
	OfficialType    int         `json:"officialType,omitempty"`
	Bsn             string      `json:"bsn,omitempty"`
	LastMsg         ChatLastMsg `json:"lastMsg,omitempty"`
	Pin             bool        `json:"pin"`
	SubData         ChatSubData `json:"subData,omitempty"`
	AvatarL         string      `json:"avatar_l,omitempty"`
	Gsn             string      `json:"gsn,omitempty"`
	Contribute      bool        `json:"contribute,omitempty"`
	Sub             string      `json:"sub,omitempty"`
	Ask             string      `json:"ask,omitempty"`
}

type Data struct {
	Pin  []string   `json:"pin"`
	List []ChatRoom `json:"list"`
}

type ChatList struct {
	Data `json:"data"`
}
