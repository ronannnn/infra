package msg

type Message struct {
	Reason             string
	ReasonTemplateData any
	Msg                string
}

func New(reason string) (m *Message) {
	m = &Message{}
	m.Reason = reason
	return
}

func (m *Message) WithReason(reason string) *Message {
	m.Reason = reason
	return m
}

func (m *Message) WithReasonTemplateData(data any) *Message {
	m.ReasonTemplateData = data
	return m
}

func (m *Message) WithMsg(msg string) *Message {
	m.Msg = msg
	return m
}
