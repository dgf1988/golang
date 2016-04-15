package hoetom

import (
	"encoding/json"
	"fmt"
	"strings"
)

const ConfigFilename = "config.ini"
const ConfigDefault = `
htmlcoding=gb18030
nettimeout=30
nettrytime=3
nettrydelayed=3
dbdriver=mysql
dbusername=root
dbpassword=guofeng001
dbhostname=localhost
dbport=3306
dbname=weiqi_hoetom
dbcharset=utf8
`

//类型
type Config map[string]string

//内容维护一个配置变量
var config Config

//初始化
func init() {
	ConfigInit()
}

func ConfigInit() {
	if err := ConfigLoad(ConfigFilename); err != nil {
		fmt.Println(err.Error())
		fmt.Println("config.ini not found, load and save default config.ini.")
		ConfigSetDefault()
		ConfigSave(ConfigFilename)
	}
}

//创建新的配置变量
func ConfigNew() Config {
	return make(Config)
}

func ConfigGet(key string) string {
	v, ok := config[key]
	if ok {
		return v
	} else {
		return ""
	}
}

func ConfigGetDefault(key, def string) string {
	value := ConfigGet(key)
	if value == "" {
		return def
	}
	return value
}

func ConfigHas(key string) bool {
	_, ok := config[key]
	return ok
}

func ConfigSet(key, value string) {
	config[key] = value
}

func ConfigSetBy(line string) {
	s := strings.SplitN(line, "=", 2)
	if len(s) != 2 {
		return
	}
	if len(s[0]) == 0 {
		return
	}
	s[0] = strings.TrimSpace(s[0])
	s[1] = strings.TrimSpace(s[1])
	ConfigSet(s[0], s[1])
}

func ConfigSetDefault() {
	config = ConfigNew()
	lines := strings.Split(ConfigDefault, "\n")
	for _, ln := range lines {
		ConfigSetBy(ln)
	}
}

func ConfigSave(filename string) {
	FileSaveLines(filename, ConfigToLines())
}

func ConfigLoad(filename string) error {
	config = ConfigNew()
	lines, err := FileLoadLines(filename)
	if err != nil {
		return err
	}
	for _, line := range lines {
		ConfigSetBy(line)
	}
	return nil
}

func ConfigToJson() string {
	x, err := json.Marshal(config)
	if err != nil {
		return ""
	}
	return string(x)
}

func ConfigToLines() []string {
	lines := make([]string, 0)
	for k, v := range config {
		lines = append(lines, fmt.Sprint(k, "=", v))
	}
	return lines
}

func ConfigLen() int {
	return len(config)
}

func ConfigPrint() {
	for _, l := range ConfigToLines() {
		fmt.Println(l)
	}
}

//初始化
//获取
//
