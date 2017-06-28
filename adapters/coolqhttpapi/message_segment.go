package coolqhttpapi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/richardchien/botadapt/t"
	"github.com/richardchien/botadapt/ucbi"
)

var faces = map[int]string{
	14: "微笑", 1: "撇嘴", 2: "色", 3: "发呆", 4: "得意",
	5: "流泪", 6: "害羞", 7: "闭嘴", 8: "睡", 9: "大哭",
	10: "尴尬", 11: "发怒", 12: "调皮", 13: "呲牙", 0: "惊讶",
	15: "难过", 16: "酷", 96: "冷汗", 18: "抓狂", 19: "吐",
	20: "偷笑", 21: "可爱", 22: "白眼", 23: "傲慢", 24: "饥饿",
	25: "困", 26: "惊恐", 27: "流汗", 28: "憨笑", 29: "大兵",
	30: "奋斗", 31: "咒骂", 32: "疑问", 33: "嘘", 34: "晕",
	35: "折磨", 36: "衰", 37: "骷髅", 38: "敲打", 39: "再见",
	97: "擦汗", 98: "抠鼻", 99: "鼓掌", 100: "糗大了", 101: "坏笑",
	102: "左哼哼", 103: "右哼哼", 104: "哈欠", 105: "鄙视", 106: "委屈",
	107: "快哭了", 108: "阴险", 109: "亲亲", 110: "吓", 111: "可怜",
	172: "眨眼睛", 182: "笑哭", 179: "doge", 173: "泪奔", 174: "无奈",
	212: "托腮", 175: "卖萌", 178: "斜眼笑", 177: "喷血", 180: "惊喜",
	181: "骚扰", 176: "小纠结", 183: "我最美", 112: "菜刀", 89: "西瓜",
	113: "啤酒", 114: "篮球", 115: "乒乓", 171: "茶", 60: "咖啡",
	61: "饭", 46: "猪头", 63: "玫瑰", 64: "凋谢", 116: "示爱",
	66: "爱心", 67: "心碎", 53: "蛋糕", 54: "闪电", 55: "炸弹",
	56: "刀", 57: "足球", 117: "瓢虫", 59: "便便", 75: "月亮",
	74: "太阳", 69: "礼物", 49: "拥抱", 76: "强", 77: "弱",
	78: "握手", 79: "胜利", 118: "抱拳", 119: "勾引", 120: "拳头",
	121: "差劲", 122: "爱你", 123: "NO", 124: "OK", 42: "爱情",
	85: "飞吻", 43: "跳跳", 41: "发抖", 86: "怄火", 125: "转圈",
	126: "磕头", 127: "回头", 128: "跳绳", 129: "挥手", 130: "激动",
	131: "街舞", 132: "献吻", 133: "左太极", 134: "右太极", 136: "双喜",
	137: "鞭炮", 138: "灯笼", 140: "K歌", 144: "喝彩", 145: "祈祷",
	146: "爆筋", 147: "棒棒糖", 148: "喝奶", 151: "飞机", 158: "钞票",
	168: "药", 169: "手枪", 188: "蛋", 192: "红包", 184: "河蟹",
	185: "羊驼", 190: "菊花", 187: "幽灵", 193: "大笑", 194: "不开心",
	197: "冷漠", 198: "呃", 199: "好棒", 200: "拜托", 201: "点赞",
	202: "无聊", 203: "托脸", 204: "吃", 205: "送花", 206: "害怕",
	207: "花痴", 208: "小样儿", 210: "飙泪", 211: "我不看",
}

func (s MessageSegment) toUCBI(messageSource map[string]string) ucbi.MessageSegment {
	dst := ucbi.MessageSegment{Data: t.JSONObject{}}

	switch s.Type {
	case "text":
		dst.Type = "text"
		dst.Text = s.Data["text"]
		dst.Data = nil
	case "emoji":
		dst.Type = "text"
		emojiCode, _ := strconv.ParseInt(s.Data["id"], 10, 32)
		dst.Text = string(emojiCode)
		dst.Data = nil
	case "face":
		dst.Type = "*face"
		dst.Text = "[表情"
		faceID, _ := strconv.ParseInt(s.Data["id"], 10, 0)
		faceText, ok := faces[int(faceID)]
		if ok {
			dst.Text += ":" + faceText
		}
		dst.Text += "]"
		dst.Data["id"] = faceID
	case "bface":
		dst.Type = "*bface"
		dst.Text = "[原创表情"
		bfaceText, ok := s.Data["text"]
		if ok {
			dst.Text += ":" + bfaceText
		}
		dst.Text += "]"
		dst.Data["id"] = s.Data["id"]
		dst.Data["p"] = s.Data["p"]
	case "sface":
		dst.Type = "*sface"
		dst.Text = "[小表情]"
		dst.Data["id"] = s.Data["id"]
	case "image":
		dst.Type = "image"
		dst.Text = "[图片]"
		dst.Data["url"] = s.Data["url"]
		dst.Data["media_id"] = s.Data["file"]
	case "record":
		dst.Type = "audio"
		dst.Text = "[语音]"
		dst.Data["media_id"] = s.Data["file"]
	case "at":
		dst.Type = "at"
		userID := s.Data["qq"]
		dst.Data["user_id"] = userID
		userName, _ := getNickname(messageSource, userID)
		dst.Data["user_name"] = userName
		dst.Data["user"] = dst.Data["user_name"]
		dst.Text = "@"
		if userName != "" {
			dst.Text += userName
		} else {
			dst.Text += "未知用户"
		}
	case "rps":
		dst.Type = "*rps"
		dst.Text = "[猜拳"
		rpsType, _ := strconv.ParseInt(s.Data["type"], 10, 0)
		switch rpsType {
		case 1:
			dst.Text += ":石头]"
		case 2:
			dst.Text += ":剪刀]"
		case 3:
			dst.Text += ":布]"
		default:
			dst.Text += "]"
		}
		dst.Data["type"] = rpsType
	case "dice":
		dst.Type = "*dice"
		dst.Text = "[骰子:" + s.Data["type"] + "]"
		dst.Data["type"], _ = strconv.ParseInt(s.Data["type"], 10, 0)
	case "shake":
		dst.Type = "*shake"
		dst.Text = "[戳一戳]"
	case "share":
		dst.Type = "link"
		dst.Text = s.Data["url"]
		dst.Data["url"] = s.Data["url"]
		dst.Data["title"] = strings.TrimSpace(s.Data["title"])
		dst.Data["content"] = strings.TrimSpace(s.Data["content"])
		dst.Data["image"] = s.Data["image"]
	case "contact":
		dst.Type = "contact"
		switch s.Data["type"] {
		case "qq":
			dst.Data["type"] = "user"
			dst.Data["user_id"] = s.Data["id"]
			dst.Text = "[推荐联系人:"
		case "group":
			dst.Data["type"] = "group"
			dst.Data["group_id"] = s.Data["id"]
			dst.Text = "[推荐群:"
		}
		dst.Text += s.Data["id"] + "]"
	case "sign":
		dst.Type = "*sign"
		dst.Text = "[群签到]"
		dst.Data["location"] = strings.TrimSpace(s.Data["location"])
		dst.Data["title"] = strings.TrimSpace(s.Data["title"])
		dst.Data["image"] = s.Data["image"]
	case "show":
		dst.Type = "*show"
		dst.Text = "[厘米秀]"
		dst.Data["id"] = s.Data["id"]
		if qq, ok := s.Data["qq"]; ok {
			dst.Data["qq"] = qq
		}
	case "location":
		dst.Type = "location"
		description := strings.TrimSpace(s.Data["content"])
		dst.Text = fmt.Sprintf("[位置:%v]", description)
		dst.Data["latitude"], _ = strconv.ParseFloat(s.Data["lat"], 64)
		dst.Data["longitude"], _ = strconv.ParseFloat(s.Data["lon"], 64)
		dst.Data["description"] = description
		dst.Data["*title"] = strings.TrimSpace(s.Data["title"])
		dst.Data["*style"] = s.Data["style"]
	case "music":
		switch s.Data["type"] {
		case "qq":
			s.Data["url"] = fmt.Sprintf("http://i.y.qq.com/v8/playsong.html?songid=%v&souce=qqaio&_wv=1&ADTAG=aiodiange", s.Data["id"])
		case "163":
			s.Data["url"] = fmt.Sprintf("http://music.163.com/m/song?id=%v#?thirdfrom=qq", s.Data["id"])
		case "xiami":
			s.Data["url"] = fmt.Sprintf("https://h.xiami.com/song.html?f=&from=&disabled=&id=%v", s.Data["id"])
		}
		fallthrough
	case "rich":
		dst.Type = "rich"
		dst.Text = "[富媒体]"
		description, ok := s.Data["text"]
		if ok {
			description = strings.TrimSpace(description)
			locationPrefix := "位置分享"
			if strings.HasPrefix(description, locationPrefix) {
				// this is actually a location share message
				dst.Type = "location"
				dst.Data["description"] = strings.TrimSpace(strings.TrimPrefix(description, locationPrefix))
				dst.Text = fmt.Sprintf("[位置:%v]", dst.Data["description"])
				break
			} else {
				dst.Data["description"] = description
				dst.Text = fmt.Sprintf("[富媒体:%v]", description)
			}
		}
		url, ok := s.Data["url"]
		if ok {
			dst.Data["url"] = url
			dst.Text = url
		}
	default:
		// types that we don't know
		dst.Type = "*" + s.Type
		for k, v := range s.Data {
			dst.Data[k] = v
		}
	}

	return dst
}
