package models

type Message struct {
	From             string `json:"from"`
	To               string `json:"to"`
	Text             string `json:"text"`
	IsPriority       bool   `json:"priority"`
	IsUrgent         bool   `json:"urgent"`
	IsGuarantee      bool   `json:"guaranteed"`
	Callback         string `json:"callback"`

	OriginalCallback string `json:"-"`
	// our message uuid
	Id   string `json:"-"`
}
