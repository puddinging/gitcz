package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/peterh/liner"
)

type CzCommit struct {
	Type           *CzType
	Scope          *string
	Subject        *string
	Body           *string
	BreakingChange *string
	Closes         *string
}

var (
	InputTypePrompt           = "选择或输入一个提交类型(必填): "
	InputScopePrompt          = "说明本次提交的影响范围(必填): "
	InputSubjectPrompt        = "对本次提交进行简短描述(必填): "
	InputBodyPrompt           = "对本次提交进行完整描述(选填): "
	InputBreakingChangePrompt = "如果当前代码版本与上一版本不兼容,对变动、变动的理由及迁移的方法进行描述(选填): "
	InputClosesPrompt         = "如果本次提交针对某个issue,列出关闭的issues(选填): "
)

func GitCz() {
	Init()
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	defer func(line *liner.State) {
		err := line.Close()
		if err != nil {
			return
		}
	}(line)
	czCommit := UserOperate(line)
	commit := GenerateCommit(&czCommit)
	if err := GitCommit(commit); err != nil {
		fmt.Println(err)
	}
}

func NewLine() {
	fmt.Println()
}

func GitCommit(commit string) (err error) {
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

func InputType(line *liner.State) *CzType {
	typeNum := len(CzTypeList)
	for i, czType := range CzTypeList {
		fmt.Printf("[%2d] %-30s: %s\n", i+1, czType.Type, czType.Message)
	}
	text, err := line.Prompt(InputTypePrompt)
	if err != nil {
		if errors.Is(err, liner.ErrPromptAborted) {
			fmt.Println("\nAborted")
			os.Exit(0)
		}
	}
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
	return InputType(line)
}

func InputScope(line *liner.State) *string {
	text, err := line.Prompt(InputScopePrompt)
	if err != nil {
		if errors.Is(err, liner.ErrPromptAborted) {
			fmt.Println("\nAborted")
			os.Exit(0)
		}
	}
	text = strings.TrimSpace(text)
	if text != "" {
		NewLine()
		return &text
	}
	NewLine()
	return InputScope(line)
}

func InputSubject(line *liner.State) *string {
	text, err := line.Prompt(InputSubjectPrompt)
	if err != nil {
		if errors.Is(err, liner.ErrPromptAborted) {
			fmt.Println("\nAborted")
			os.Exit(0)
		}
	}
	text = strings.TrimSpace(text)
	if text != "" {
		NewLine()
		return &text
	}
	NewLine()
	return InputSubject(line)
}

func InputBody(line *liner.State) *string {
	text, err := line.Prompt(InputBodyPrompt)
	if err != nil {
		if errors.Is(err, liner.ErrPromptAborted) {
			fmt.Println("\nAborted")
			os.Exit(0)
		}
	}
	text = strings.TrimSpace(text)
	if text != "" {
		NewLine()
		return &text
	}
	NewLine()
	return nil
}

func InputBreakingChange(line *liner.State) *string {
	text, err := line.Prompt(InputBreakingChangePrompt)
	if err != nil {
		if errors.Is(err, liner.ErrPromptAborted) {
			fmt.Println("\nAborted")
			os.Exit(0)
		}
	}
	text = strings.TrimSpace(text)
	if text != "" {
		NewLine()
		return &text
	}
	NewLine()
	return nil
}

func InputCloses(line *liner.State) *string {
	text, err := line.Prompt(InputClosesPrompt)
	if err != nil {
		if errors.Is(err, liner.ErrPromptAborted) {
			fmt.Println("\nAborted")
			os.Exit(0)
		}
	}
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

func UserOperate(line *liner.State) CzCommit {
	czCommit := &CzCommit{}
	czCommit.Type = InputType(line)
	czCommit.Scope = InputScope(line)
	czCommit.Subject = InputSubject(line)
	czCommit.Body = InputBody(line)
	czCommit.BreakingChange = InputBreakingChange(line)
	czCommit.Closes = InputCloses(line)
	return *czCommit
}

/*
*
打印帮助信息
*/
func help() {
	print("打印帮助信息")
	print("打印帮助信息")
}
