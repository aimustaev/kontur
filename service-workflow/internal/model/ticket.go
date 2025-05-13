package model

type Ticket struct {
	VerticalId string
	UserId     string
	Assign     string
	SkillId    string
}

type Message struct {
	ID      string
	From    string
	To      string
	Subject string
	Body    string
	Tags    []string
	Channel string // email, telegram
}
