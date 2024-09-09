package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func InitI18n() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// 加載翻譯文件
	loadTranslationFile("en.json")
	loadTranslationFile("zh-TW.json")
}

func loadTranslationFile(filename string) {
	// 獲取當前文件的絕對路徑
	_, currentFile, _, _ := runtime.Caller(0)
	// 獲取項目根目錄（假設 config 目錄直接位於項目根目錄下）
	rootDir := filepath.Dir(filepath.Dir(currentFile))
	// 構建到 locales 目錄的完整路徑
	path := filepath.Join(rootDir, "locales", filename)

	_, err := bundle.LoadMessageFile(path)
	if err != nil {
		log.Printf("加載翻譯文件 %s 時出錯: %v\n", filename, err)
		// 列出 locales 目錄內容以進行調試
		entries, _ := os.ReadDir(filepath.Join(rootDir, "locales"))
		log.Printf("locales 目錄內容: %v", getFileNames(entries))
	} else {
		log.Printf("成功加載翻譯文件: %s\n", filename)
	}
}

// 輔助函數：獲取目錄條目的文件名列表
func getFileNames(entries []os.DirEntry) []string {
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}
	return names
}

func Translate(msgID string, language string, params map[string]interface{}) string {
	if bundle == nil {
		return msgID
	}

	localizer := i18n.NewLocalizer(bundle, language)
	if localizer == nil {
		return msgID
	}

	translation, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    msgID,
		TemplateData: params,
	})
	if err != nil {
		log.Printf("翻譯錯誤: %v", err)
		return msgID
	}

	return translation
}
