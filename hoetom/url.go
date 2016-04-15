package hoetom

import (
	"fmt"
)

func UrlPlayer(playerid int64) string {
	return fmt.Sprint("http://www.hoetom.com/playerinfor_2011.jsp?id=", playerid)
}

func UrlSgf(sgfid int64) string {
	return fmt.Sprint("http://www.hoetom.com/matchviewer_html_2011.jsp?id=", sgfid)
}

func UrlSgfListLast(page int64) string {
	return fmt.Sprint("http://www.hoetom.com/matchlatest_2011.jsp?pn=", page)
}


var playerlistquery string = "0abcdefghijklmnopqrstuvwxyz*"

func UrlPlayerList(index int) string {
	if index > 27 || index <= 0 {
		return "http://www.hoetom.com/playeranking_2011.jsp"
	}
	queryindex := playerlistquery[index]
	queryurl := "http://www.hoetom.com/playerlist_2011.jsp?nid=-1&pc=1&ln=" + string(queryindex)
	return queryurl
}