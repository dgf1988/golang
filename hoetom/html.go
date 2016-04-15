package hoetom

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func HtmlTitle(text string) string {
	re, err := regexp.Compile("<title>(.*?)</title>")
	if err != nil {
		return err.Error()
	}
	m := re.FindStringSubmatch(text)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}

//从Html查找棋谱ID
func HtmlAllSgfid(text string) []int64 {
	resgfid := regexp.MustCompile(`(matchinfor_2011\.jsp\?id=|matchviewer_html_2011\.jsp\?id=)(\d+)`)
	sgfidsmatch := resgfid.FindAllStringSubmatch(text, -1)
	allid := make([]int64, len(sgfidsmatch))
	for i := range sgfidsmatch {
		if len(sgfidsmatch[i]) != 3 {
			continue
		}
		id, err := strconv.ParseInt(sgfidsmatch[i][2], 10, 64)
		if err != nil {
			continue
		}
		allid[i] = id
	}
	return allid
}

func HtmlAllPlayerid(text string) []int64 {
	re_playerid := regexp.MustCompile(`playerinfor_2011\.jsp\?id=(\d+)`)
	ms_playerid := re_playerid.FindAllStringSubmatch(text, -1)
	list_playerid := make([]int64, 0)
	for i := range ms_playerid {
		findid, err := strconv.ParseInt(ms_playerid[i][1], 10, 64)
		if err != nil {
			return list_playerid
		}
		list_playerid = append(list_playerid, findid)
	}
	return list_playerid
}

var (
	ErrHtmlFindPlayer      error = errors.New("Html: 找不到棋手资料")
	ErrHtmlFindPlayerName  error = errors.New("Html: 找不到棋手名字或ID")
	ErrHtmlFIndPlayerItems error = errors.New("Html: 找不到棋手数据")
)

func HtmlFindPlayer(text string) ([]string, error) {
	re_p, err := regexp.Compile(`<table id="table1" summary="([^"]+)"[\S\s]*?</table>`)
	if err != nil {
		return nil, err
	}
	re_name, err := regexp.Compile(`>姓名<(?s:.*)id=(\d+)[^>]*?>([^<]+)`)
	if err != nil {
		return nil, err
	}
	re_items, err := regexp.Compile(`"#DEDFDE">[^<]+</th>[\s]*<td><b>([^<]*)</b></td>`)
	if err != nil {
		return nil, err
	}
	res := make([]string, 7)

	m_p := re_p.FindStringSubmatch(text)
	if len(m_p) != 2 {
		return nil, ErrHtmlFindPlayer
	}

	m_name := re_name.FindStringSubmatch(m_p[0])
	if len(m_name) != 3 {
		return nil, ErrHtmlFindPlayerName
	}
	res[0] = m_name[1]
	res[1] = m_name[2]
	m_items := re_items.FindAllStringSubmatch(m_p[0], -1)
	if len(m_items) != 5 {
		return nil, ErrHtmlFIndPlayerItems
	}
	for i := range m_items {
		res[i+2] = m_items[i][1]
	}
	return res, nil
}

//错误定义
var (
	ErrHtmlFindSgf              error = errors.New("Html: 找不到棋谱")
	ErrHtmlFindSgfItems         error = errors.New("Html: 找不到棋谱的项目")
	ErrHtmlFindSgfItemDatas     error = errors.New("Html: 找不到棋谱项目的数据")
	ErrHtmlFindSgfId            error = errors.New("Html: 找不到棋谱ID")
	ErrHtmlFindSgfSteps         error = errors.New("Html: 找不到棋谱的下棋数据")
	ErrHtmlFindSgfJsonSteps     error = errors.New("Html: json棋谱解析错误")
	ErrHtmlFindSgfJsonStepsType error = errors.New("Html: json棋谱数据类型断言失败")
)

//提取数据
func HtmlFindSgf(text string) (*Sgf, error) {
	//匹配table
	retable := regexp.MustCompile(`<table width="300" id=table2>\s+((<tr>(?s:.*?)</tr>\s+){8})`)
	tablematch := retable.FindStringSubmatch(text)
	if len(tablematch) != 3 {
		return nil, ErrHtmlFindSgf
	}

	//数据保存
	var sgf Sgf

	//匹配ID
	resgfid := regexp.MustCompile(`var id = (\d+)`)
	sgfidmatch := resgfid.FindStringSubmatch(text)
	if len(sgfidmatch) == 2 {
		id, err := strconv.Atoi(sgfidmatch[1])
		if err != nil {
			return nil, err
		}
		//保存ID
		sgf.Sgfid = int64(id)
	} else {
		return nil, ErrHtmlFindSgfId
	}

	//匹配项目
	retrs := regexp.MustCompile(`<tr>\s+(<td.*?</td>\s+<td>(?s:.*?)</td>)\s+</tr>`)
	trsmatch := retrs.FindAllString(tablematch[1], -1)
	if len(trsmatch) != 8 {
		return nil, ErrHtmlFindSgfItems
	}

	//匹配项目数据
	datas := make([]string, 8)
	redata := regexp.MustCompile(`<td>([^<]*)`)
	for i := range trsmatch {
		data := redata.FindStringSubmatch(trsmatch[i])
		if len(data) != 2 {
			return nil, ErrHtmlFindSgfItemDatas
		}
		datas[i] = data[1]
	}
	sgf.Event = datas[0]
	sgf.Black = datas[1]
	sgf.White = datas[2]
	sgf.Rule = datas[3]
	sgf.Result = datas[4]
	if len(datas[6]) == 10 {
		datatime, err := time.Parse("2006-01-02", datas[6])
		if err != nil {
			sgf.Time = time.Time{}
		} else {
			sgf.Time = datatime
		}
	} else {
		sgf.Time = time.Time{}
	}
	sgf.Place = datas[7]

	//匹配棋谱
	resgf := regexp.MustCompile(`var steps = (\[.*\]);`)
	sgfmatch := resgf.FindStringSubmatch(text)
	if len(sgfmatch) != 2 {
		return &sgf, ErrHtmlFindSgfSteps
	}
	steps := make([]map[string]interface{}, 0)
	stepitems := make([]string, 0)
	err := json.Unmarshal([]byte(sgfmatch[1]), &steps)
	if err != nil {
		return &sgf, ErrHtmlFindSgfJsonSteps
	}
	for i := range steps {
		pass, ok := steps[i]["pass"].(bool)
		if !ok {
			return &sgf, ErrHtmlFindSgfJsonStepsType
		}
		c, ok := steps[i]["comment"].(string)
		if !ok {
			return &sgf, ErrHtmlFindSgfJsonStepsType
		}
		x, ok := steps[i]["x"].(float64)
		if !ok {
			return &sgf, ErrHtmlFindSgfJsonStepsType
		}
		y, ok := steps[i]["y"].(float64)
		if !ok {
			return &sgf, ErrHtmlFindSgfJsonStepsType
		}
		w, ok := steps[i]["white"].(bool)
		if !ok {
			return &sgf, ErrHtmlFindSgfJsonStepsType
		}
		b, ok := steps[i]["black"].(bool)
		if !ok {
			return &sgf, ErrHtmlFindSgfJsonStepsType
		}
		var step Step
		step.C = c
		step.X = int(x)
		step.Y = int(y)
		if w {
			step.P = 2
		} else if b {
			step.P = 1
		} else {
			step.P = 0
		}
		if pass {
			step.X = -1
			step.Y = -1
		}
		stepitems = append(stepitems, step.ToSgf())
	}
	sgf.Steps = strings.Join(stepitems, "")
	return &sgf, nil
}
