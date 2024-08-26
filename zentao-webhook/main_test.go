package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type fakeMsgApp struct{}

func (*fakeMsgApp) GetToken() (string, int64, error) {
	if _, err := base64.StdEncoding.DecodeString(encAesKey); err != nil {
		return "", 0, err
	}

	return "4432fe3c996lee3d9dldabe33bb6cdff", 7200, nil
}

func (*fakeMsgApp) SendTxtMsg(toUser, toDept, content string) error {
	if toUser != "Neo" {
		return errors.New("user not exist")
	}
	return nil
}

func TestSend(t *testing.T) {
	var fma *fakeMsgApp

	if err := send(fma, "Oracle", "You're cuter than I thought."); err != nil {
		if err.Error() != "user not exist" {
			t.Error(err)
		}
	}

	if err := send(fma, "Neo", "I can't go back, can I?"); err != nil {
		t.Error(err)
	}

	origin := os.Getenv("EncodingAESKey")
	defer func() {
		encAesKey = origin
	}()
	encAesKey = "abc"
	if err := send(fma, "Neo", "I can't go back, can I?"); err != nil {
		if !strings.HasPrefix(err.Error(), "illegal base64 data") {
			t.Error(err)
		}
	}

}

func TestHandler(t *testing.T) {
	var fma *fakeMsgApp
	handler := webhookHandlerWrapper(fma)

	testCases := []struct {
		in   ZenTaoWebHookData
		want int
	}{
		{
			ZenTaoWebHookData{
				ObjectType: "task",
				ObjectID:   "122",
				Product:    "1",
				Project:    "",
				Action:     "finished",
				Actor:      "Neo",
				Date:       "2024-05-28 21:02:05",
				Comment:    "@Trinity",
				Text:       "Neo finished the job [#122::Reboot system](http://172.20.40.124/zentao/task-view-122.html)",
			},
			400,
		},
		{
			ZenTaoWebHookData{
				ObjectType: "task",
				ObjectID:   "122",
				Product:    "1",
				Project:    "",
				Action:     "finished",
				Actor:      "Neo",
				Date:       "2024-05-28 21:02:05",
				Comment:    "@Trinity@Morpheus",
				Text:       "Neo finished the job [#122::Reboot system](http://172.20.40.124/zentao/task-view-122.html)",
			},
			400,
		},
		{
			ZenTaoWebHookData{
				ObjectType: "task",
				ObjectID:   "122",
				Product:    "1",
				Project:    "",
				Action:     "finished",
				Actor:      "Neo",
				Date:       "2024-05-28 21:02:05",
				Comment:    "",
				Text:       "Neo finished the job [#122::Reboot system](http://172.20.40.124/zentao/task-view-122.html)",
			},
			200,
		},
		{
			ZenTaoWebHookData{
				ObjectType: "task",
				ObjectID:   "123",
				Product:    "1",
				Project:    "",
				Action:     "opened",
				Actor:      "Trinity",
				Date:       "2024-05-28 21:02:06",
				Comment:    "@Neo",
				Text:       "Trinity created a job [#122::Make a deal with the machine](http://172.20.40.124/zentao/task-view-123.html)",
			},
			200,
		},
	}
	for index, tc := range testCases {
		t.Run(fmt.Sprintf("Case-%d\n", index), func(t *testing.T) {
			payloadBytes, err := json.Marshal(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest("POST", "http://172.20.40.129:8090", bytes.NewBuffer(payloadBytes))
			w := httptest.NewRecorder()
			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != tc.want {
				t.Errorf("%+v\n", resp)
			}
		})
	}
}
