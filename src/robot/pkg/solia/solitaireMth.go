package solia

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// 空结构体
var Exists = struct{}{}

//Set
type Set struct {
	// struct为结构体类型的变量
	m map[interface{}]struct{}
}

func (s *Set) Add(items ...interface{}) error {
	for _, item := range items {
		s.m[item] = Exists
	}
	return nil
}

func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]
	return ok
}

// 清空set
func (s *Set) Clear() {
	s.m = make(map[interface{}]struct{})
}

//多人游戏对象
type MpSolia struct {
	rd   []string          // 存储成语数据
	Mp   map[string]*Solia //游戏者信息
	lock sync.RWMutex
}

type Solia struct {
	StrSet Set    //已经使用过的成语
	tryNum int    // 尝试次数
	nowStr string //当前接龙的成语
}

//多人游戏开始
func (ms *MpSolia) ReadStart(userID string) (string, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	if ms.Mp[userID] != nil { //已经有游戏者信息了，说明已经开始游戏了
		return "", errors.New(ReStart)
	} else {
		ms.Mp[userID] = &Solia{}
		return ms.Mp[userID].readStart(userID, ms)
	}
}

//单人游戏开始获取第一个成语
func (s *Solia) readStart(userID string, ms *MpSolia) (string, error) {
	//打开文件
	s.tryNum = 3
	s.StrSet.m = make(map[interface{}]struct{})
	rand.Seed(time.Now().Unix())
	n := rand.Intn(9999) + 1
	str, err := ms.readLineNum(n)
	s.StrSet.Add(str)
	s.nowStr = str
	return str, err
}

func (ms *MpSolia) readLineNum(lineNum int) (string, error) {
	//成语数据为空，则重新获取
	if ms.rd == nil {
		err := ms.getFiles()
		if err != nil {
			return "", err
		}
	}
	if lineNum > 0 && lineNum <= len(ms.rd) {
		for i := lineNum - 1; i < len(ms.rd); i++ {
			line := ms.rd[i]
			if utf8.RuneCountInString(line) == 4 {
				return line, nil
			}
		}
	}
	return "人山人海", nil
}

func (ms *MpSolia) ReadStr(content string, userId string) (string, error) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	str, err, b := ms.Mp[userId].readStr(content, ms)
	if b { //游戏结束
		delete(ms.Mp, userId)
	}
	if err != nil {
		return "", err
	}
	return str, err
}

func (s *Solia) readStr(content string, ms *MpSolia) (string, error, bool) {
	strs := strings.Split(content, ">")
	content = strings.TrimSpace(strs[len(strs)-1])
	if s.StrSet.Contains(content) { //判断是否重复的成语
		if s.tryNum > 1 {
			s.tryNum--
			return "", errors.New(fmt.Sprintf(ContainsNotOver, s.tryNum)), false
		} else {
			return "", errors.New(fmt.Sprintf(ContainsOver3)), true
		}
	}
	s1 := string([]rune(s.nowStr)[utf8.RuneCountInString(s.nowStr)-1:])
	if s1 != string([]rune(content)[:1]) { //判断首字是否和上一个成语尾字一样
		if s.tryNum > 1 {
			s.tryNum--
			return "", errors.New(fmt.Sprintf(SameNotOver, s.tryNum)), false
		} else {
			return "", errors.New(fmt.Sprintf(SameNotOver3)), true
		}
	}
	var res string
	var flag bool
	if len([]rune(content)) >= 4 {
		str1 := string([]rune(content)[3:])

		if ms.rd == nil {
			err := ms.getFiles()
			if err != nil {
				return "", err, false
			}
		}
		for i := 0; i < len(ms.rd); i++ {
			line := ms.rd[i]
			if str1 == string([]rune(line)[:1]) && !s.StrSet.Contains(line) {
				res = line
			}
			if line == content {
				flag = true
			}
			if flag && res != "" {
				break
			}
		}
	}
	if !flag {
		if s.tryNum > 1 {
			s.tryNum--
			return "", errors.New(fmt.Sprintf(ErrorNotOver, s.tryNum)), false
		} else {
			return "", errors.New(fmt.Sprintf(ErrorOver3)), true
		}
	} else {
		s.StrSet.Add(content)
		s.nowStr = content
	}
	if res == "" {
		return "", errors.New(fmt.Sprintf(Success)), true
	} else {
		s.StrSet.Add(res)
		s.nowStr = res
	}
	return res, nil, false
}

func (ms *MpSolia) getFiles() error {
	strPath, _ := os.Getwd()
	strPath = strPath[:strings.LastIndex(strPath, "robot")+5]
	file, err := os.Open(strPath + "./idiom.txt") //只是用来读的时候，用os.Open。绝对路径，获取robot路径
	if err != nil {
		fmt.Printf("打开文件失败,err:%v\n", err)
		return err
	}
	defer file.Close()                                 //关闭文件,为了避免文件泄露和忘记写关闭文件
	decoder := mahonia.NewDecoder("gbk")               //转码，避免中文字符乱码
	reader := bufio.NewReader(decoder.NewReader(file)) //创建新的读的对象
	var chunks []string
	for {
		line, err := reader.ReadString('\n')
		line = strings.Replace(line, "\r\n", "", -1)
		if err == io.EOF {
			fmt.Println("文件读完了")
			break
		}
		if err != nil { //错误处理
			fmt.Printf("读取文件失败,错误为:%v", err)
			return err
		}

		chunks = append(chunks, line)
	}
	ms.rd = chunks
	return nil
}

func (ms *MpSolia) GetMapValue(userId string) *Solia {
	ms.lock.RLock()
	defer ms.lock.RUnlock()
	return ms.Mp[userId]
}

func (ms *MpSolia) DeleteMapValue(userId string) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	delete(ms.Mp, userId)
}
