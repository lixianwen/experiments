package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	ydapp "github.com/youduim/EntAppSdkGo"
)

// 禅道WebHook数据
type ZenTaoWebHookData struct {
	// 对象类型
	ObjectType string `json:"objectType"`
	// 对象ID
	ObjectID string `json:"objectID"`
	// 产品ID
	Product string `json:"product"`
	// 项目ID
	Project string `json:"project"`
	// 动作
	Action string `json:"action"`
	// 操作者
	Actor string `json:"actor"`
	// 操作时间
	Date string `json:"date"`
	// 备注
	Comment string `json:"comment"`
	// 操作内容
	Text string `json:"text"`
}

type APP interface {
	GetToken() (string, int64, error)
	SendTxtMsg(string, string, string) error
}

const Buin = 123456

var (
	appId     string
	encAesKey string
)

func init() {
	var ok bool
	if appId, ok = os.LookupEnv("AppID"); !ok {
		panic("请设置环境变量 'AppID'")
	}
	if encAesKey, ok = os.LookupEnv("EncodingAESKey"); !ok {
		panic("请设置环境变量 'EncodingAESKey'")
	}
}

func webhookHandlerWrapper(app APP) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload ZenTaoWebHookData
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Printf("Unmarshal err, %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("zenTaoData: %+v\n", payload)

		re := regexp.MustCompile(`[^@]+`)
		referTo := re.FindAllStringSubmatch(payload.Comment, -1)
		// // Action
		// re = regexp.MustCompile(`(.*)(创建|变更|编辑|激活|指派|评论|确认|完成|解决|关闭)`)
		// actor := re.FindStringSubmatch(payload.Text)
		// // match actor
		// if len(actor) == 3 && actor[1] != "" {
		// 	referTo = append(referTo, []string{actor[1]})
		// }
		log.Println("referTo", referTo)

		sl := make([]string, 0, len(referTo))
		for _, submatch := range referTo {
			for _, match := range submatch {
				match = strings.TrimSpace(match)
				if match != "" {
					sl = append(sl, match)
				}
			}
		}
		toUser := strings.Join(sl, "|")
		fmt.Printf("##%s\n", toUser)
		if toUser == "" {
			log.Println("No receiver.")
			return
		}

		// send message
		message := payload.Text
		if payload.Comment != "" {
			message += "备注：" + payload.Comment
		}
		if err := send(app, toUser, message); err != nil {
			log.Printf("Failed to send message, %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func send(app APP, to, message string) error {
	token, expire, err := app.GetToken()
	if err != nil {
		return err
	}
	log.Printf("Token: %s, Expire: %d", token, expire) //expire为过期的时间戳，单位秒

	if err := app.SendTxtMsg(to, "", message); err != nil {
		return err
	}
	return nil
}

func main() {
	ydapp.Server_Addr = "http://1.1.1.1:8090" //设置服务器地址
	app, err := ydapp.NewMsgApp(Buin, appId, encAesKey)
	if err != nil {
		log.Println("New app error:", err)
		return
	}

	http.HandleFunc("/", webhookHandlerWrapper(app))
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
