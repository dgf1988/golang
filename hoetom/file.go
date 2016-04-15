package hoetom

import (
	"io/ioutil"
	"os"
	"strings"
	"bufio"
	"io"
)

func FileLoad(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func FileLoadLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	lines := make([]string, 0)
	for {
		buf, isPrefix, err := r.ReadLine()
		if err == io.EOF {
			return lines, nil
		}
		if err != nil {
			return nil, err
		}
		if isPrefix {
			return lines, nil
		}
		lines = append(lines, string(buf))
	}
	return lines, nil
}

func FileSave(filename, datas string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(datas)
	if err != nil {
		return err
	}
	return nil
}

func FileSaveLines(filename string, lines []string) error {
	return FileSave(filename, strings.Join(lines, "\n"))
}


/*

var RootPath string = "d:\\HtmlStor\\hoetom\\"

//var RootPath string = ""
var PlayerListPath string = "playerlist\\"
var PlayerInfoPath string = "playerinfo\\"
var SgfPath string = "sgf\\"

func MkdirForPlayerList() error {
	var pathname string = RootPath + PlayerListPath
	fileinfo, err := os.Stat(pathname)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(pathname, 0666)
		} else {
			return err
		}
	}
	if !fileinfo.IsDir() {
		return os.MkdirAll(pathname, 0666)
	}
	return nil
}

func MkdirForPlayerInfo() error {
	var pathname string = RootPath + PlayerInfoPath
	fileinfo, err := os.Stat(pathname)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(pathname, 0666)
		} else {
			return err
		}
	}
	if !fileinfo.IsDir() {
		return os.MkdirAll(pathname, 0666)
	}
	return nil
}

func MkdirForSgf() error {
	var pathname string = RootPath + SgfPath
	sgfpath, err := os.Stat(pathname)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(pathname, 0666)
		} else {
			return err
		}
	}
	if !sgfpath.IsDir() {
		return os.MkdirAll(pathname, 0666)
	}
	return nil
}

func GetPath(filename string) string {
	repath := regexp.MustCompile(".*\\\\")
	filepath := repath.FindString(filename)
	return filepath
}

func MkdirAll(pathname string) error {
	path, err := os.Stat(pathname)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(pathname, 0666)
		}
		return err
	}
	if !path.IsDir() {
		return os.MkdirAll(pathname, 0666)
	}
	return nil
}

func SaveFile(filename string, datas string) error {
	MkdirAll(GetPath(RootPath + filename))
	f, err := os.Create(RootPath + filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(datas)
	if err != nil {
		return err
	}
	return nil
}

func LoadFile(filename string) (string, error) {
	_file, err := os.Open(RootPath + filename)
	if err != nil {
		return "", err
	}
	defer _file.Close()
	textbuf, err := ioutil.ReadAll(_file)
	if err != nil {
		return "", err
	}
	return string(textbuf), nil
}

*/