package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var petEmojis = []string{"🐶", "🐱", "🐰", "🐹", "🐦", "🐟", "🐢", "🦜", "🦊", "🐻"}

var breedOptions = map[string][]string{
	"🐶": {"金毛", "拉布拉多", "柯基", "柴犬", "哈士奇", "泰迪", "边牧", "萨摩耶", "博美", "其他犬种"},
	"🐱": {"英短", "美短", "布偶", "暹罗", "橘猫", "狸花", "波斯", "无毛猫", "其他猫种"},
	"🐰": {"侏儒兔", "垂耳兔", "荷兰兔", "安哥拉兔", "其他兔种"},
	"🐹": {"金丝熊", "三线", "一线", "其他仓鼠"},
	"🐦": {"虎皮", "玄凤", "牡丹", "其他鹦鹉"},
	"🐟": {"金鱼", "锦鲤", "斗鱼", "热带鱼"},
	"🐢": {"巴西龟", "陆龟", "蛋龟"},
	"🦜": {"金刚鹦鹉", "葵花鹦鹉", "灰鹦鹉"},
	"🦊": {"赤狐", "北极狐", "耳廓狐"},
	"🐻": {"其他宠物"},
}

func GetBreeds(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"petEmojis":    petEmojis,
		"breedOptions": breedOptions,
	})
}
