package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type CzType struct {
	Type    string
	Message string
}

type CzCommit struct {
	Type           *CzType
	Scope          *string
	Subject        *string
	Body           *string
	BreakingChange *string
	Closes         *string
}

var StdinInput = bufio.NewReader(os.Stdin)

var (
	InputTypePrompt           = "选择或输入一个提交类型(必填): "
	InputScopePrompt          = "说明本次提交的影响范围(必填): "
	InputSubjectPrompt        = "对本次提交进行简短描述(必填): "
	InputBodyPrompt           = "对本次提交进行完整描述(选填): "
	InputBreakingChangePrompt = "如果当前代码版本与上一版本不兼容,对变动、变动的理由及迁移的方法进行描述(选填): "
	InputClosesPrompt         = "如果本次提交针对某个issue,列出关闭的issues(选填): "
)

var CzTypeList = []CzType{
	{
		Type:    ":tada:",
		Message: "初始化:	第一次提交",
	},
	{
		Type:    ":sparkles:",
		Message: "功能:	一个新的功能",
	},
	{
		Type:    ":bug:",
		Message: "修复:	修复一个Bug",
	},
	{
		Type:    ":memo:",
		Message: "文档:	变更的只有文档",
	},
	{
		Type:    ":art:",
		Message: "格式:	空格, 分号等格式修复'",
	},
	{
		Type:    ":hammer:",
		Message: "重构:	代码重构，注意和特性、修复区分开",
	},
	{
		Type:    ":zap:",
		Message: "性能:	提升性能",
	},
	{
		Type:    ":white_check_mark:",
		Message: "测试:	添加一个测试",
	},
	{
		Type:    "chore",
		Message: "工具:	开发工具变动(构建、脚手架工具等)",
	},
	{
		Type:    ":twisted_rightwards_arrows:",
		Message: "分支合并",
	},
}

func main() {
	amend := flag.Bool(
		"amend",
		false,
		"覆盖上次提交信息",
	)
	sign := flag.Bool("S", false, "对commit进行签名")
	author := flag.Bool(
		"author",
		false,
		"关于本软件开发者",
	)
	flag.Parse()
	if *author {
		Author()
		return
	}
	czCommit := &CzCommit{}
	czCommit.Type = InputType()
	czCommit.Scope = InputScope()
	czCommit.Subject = InputSubject()
	czCommit.Body = InputBody()
	czCommit.BreakingChange = InputBreakingChange()
	czCommit.Closes = InputCloses()
	commit := GenerateCommit(czCommit)
	if err := GitCommit(commit, *amend, *sign); err != nil {
		fmt.Println(err)
	}
}

func NewLine() {
	fmt.Println()
}

func GitCommit(commit string, amend bool, sign bool) (err error) {
	tempFile, err := os.CreateTemp("", "git_commit_")
	if err != nil {
		return
	}
	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
	}()
	if _, err = tempFile.WriteString(commit); err != nil {
		return
	}
	args := []string{"commit"}
	if amend {
		args = append(args, "--amend")
	}
	if sign {
		args = append(args, "-S")
	}
	args = append(args, "-F", tempFile.Name())
	cmd := exec.Command("git", args...)
	result, err := cmd.CombinedOutput()
	if err != nil && !strings.ContainsAny(err.Error(), "exit status") {
		return
	} else {
		fmt.Println(string(bytes.TrimSpace(result)))
	}
	return nil
}

func InputType() *CzType {
	typeNum := len(CzTypeList)
	for i := 0; i < typeNum; i++ {
		fmt.Printf("[%d] %s:\t%s\n", i+1, CzTypeList[i].Type, CzTypeList[i].Message)
	}
	fmt.Print(InputTypePrompt)
	text, _ := StdinInput.ReadString('\n')
	text = strings.TrimSpace(text)
	selectId, err := strconv.Atoi(text)
	if err == nil && (selectId > 0 && selectId <= typeNum) {
		NewLine()
		return &CzTypeList[selectId-1]
	}
	for i := 0; i < typeNum; i++ {
		if text == CzTypeList[i].Type {
			NewLine()
			return &CzTypeList[i]
		}
	}
	NewLine()
	return InputType()
}

func InputScope() *string {
	fmt.Print(InputScopePrompt)
	text, _ := StdinInput.ReadString('\n')
	text = strings.TrimSpace(text)
	if text != "" {
		NewLine()
		return &text
	}
	NewLine()
	return InputScope()
}

func InputSubject() *string {
	fmt.Print(InputSubjectPrompt)
	text, _ := StdinInput.ReadString('\n')
	text = strings.TrimSpace(text)
	if text != "" {
		NewLine()
		return &text
	}
	NewLine()
	return InputSubject()
}

func InputBody() *string {
	fmt.Print(InputBodyPrompt)
	text, _ := StdinInput.ReadString('\n')
	text = strings.TrimSpace(text)
	if text != "" {
		NewLine()
		return &text
	}
	NewLine()
	return nil
}

func InputBreakingChange() *string {
	fmt.Print(InputBreakingChangePrompt)
	text, _ := StdinInput.ReadString('\n')
	text = strings.TrimSpace(text)
	if text != "" {
		NewLine()
		return &text
	}
	NewLine()
	return nil
}

func InputCloses() *string {
	fmt.Print(InputClosesPrompt)
	text, _ := StdinInput.ReadString('\n')
	text = strings.TrimSpace(text)
	if text != "" {
		NewLine()
		return &text
	}
	NewLine()
	return nil
}

func GenerateCommit(czCommit *CzCommit) string {
	commit := fmt.Sprintf(
		"%s(%s): %s\n\n",
		czCommit.Type.Type,
		*czCommit.Scope,
		*czCommit.Subject,
	)
	if czCommit.Body != nil {
		commit += *czCommit.Body
		commit += "\n\n"
	}
	if czCommit.BreakingChange != nil {
		commit += "BREAKING CHANGE: " + *czCommit.BreakingChange
		commit += "\n\n"
	}
	if czCommit.Closes != nil {
		commit += "Closes fix " + *czCommit.Closes
	}
	return commit
}
