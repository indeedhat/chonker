package report

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/indeedhat/chonker/internal/types"
	"github.com/indeedhat/dotenv"
)

const (
	pushEndpoint = "https://api.pushbullet.com/v2/pushes"
)

var (
	recipientEmail dotenv.String = "PB_EMAIL"
	authToken      dotenv.String = "PB_AUTH"
)

type pushMessage struct {
	Type  string `json:"type"`
	Email string `json:"email"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func ReportOnFeed(r types.Report) {
	msg := pushMessage{
		Type:  "note",
		Email: recipientEmail.Get(),
	}

	if r.Error() != nil {
		msg.Title = "Chonker Error"
		msg.Body = r.Error().Error()
	} else {
		msg.Title = r.Title()
		msg.Body = r.Message()
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	body := bytes.NewReader(data)
	req, err := http.NewRequest("POST", pushEndpoint, body)
	if err != nil {
		return
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Access-Token", authToken.Get())
}
