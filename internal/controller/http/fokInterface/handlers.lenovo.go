package fokusov

import (
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"github.com/gin-gonic/gin"
)

func (fok *Fokusov) lenovoAlarm(c *gin.Context) {
	alarmType := c.PostForm("")
	alarmDate := c.PostForm("")
	alarmText := c.PostForm("")

	ipAddr := c.PostForm("")
	roomNameRus := c.PostForm("")
	regionRus := c.PostForm("")
	loginEmployee := c.PostForm("")
	vcsType := c.PostForm("")

	zabbixAlarm := &entity.Zabbix{}

	lenovoVcs := &entity.Vcs{}

	if alarmType == "NewAlarm" {

		err := fok.LenovoRest.newAlarm(zabbixAlarm, lenovoVcs)
		if err != nil {

		} else {

		}
	}

	if alarmType == "CloseAlarm" {

	}

}
