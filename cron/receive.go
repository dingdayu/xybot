package cron

import (
	"fmt"
	"github.com/dingdayu/wxbot/types"
	"github.com/dingdayu/wxbot/utils"
	"strings"
	"os/user"
	"log"
)

func GroupChange(Msg types.Message) Message {
	if strings.Contains(Msg.Content, "加入了群聊") || strings.Contains(Msg.Content, "移出了群聊") {
		if strings.Contains(Msg.Content, "邀请你") {
			// INVITE 邀请你加入新群
			msg := Message{
				MsgId:        Msg.MsgId,
				UUID:         user.uuid,
				Event:        RedPackets,
				FromUserName: Msg.FromUserName,
				FromNickName: Contact.NickName,
				ToUserName:   Msg.ToUserName,
				Content:      Msg.Content,
				SendTime:     Msg.CreateTime,
			}
			checkGroupMsg(&msg, user.LoginUser.UserName)

			log.Println(msg)
		}
		if strings.Contains(Msg.Content, "加入了群聊") || strings.Contains(Msg.Content, "分享的二维码加入群聊") {
			// ADD 新人入群
			name := utils.PregMatch(`邀请"(.+)"加入了群聊`, Msg.Content)
			if len(name) <= 0 {
				name = utils.PregMatch(`"(.+)"通过扫描.+分享的二维码加入群聊`, Msg.Content)
			}
			fmt.Println(name)
		}
		if strings.Contains(Msg.Content, "移出了群聊") {
			// REMOVE 被移除
		}
		if strings.Contains(Msg.Content, "改群名为") {
			// RENAME 群名修改
			name := utils.PregMatch(`改群名为“(.+)”`, Msg.Content)
			fmt.Println(name)
		}
		if strings.Contains(Msg.Content, "移出群聊") {
			// BE_REMOVE 被移除群聊
		}
	}
}
