package tepidreload

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type DevTemplates struct {
	IsDev          bool
	DevTemplateMap map[string]string
}

func (dt *DevTemplates) MakeLocalDevTemplates(templateRootDir string, templateExtension string) {
	retMap := make(map[string]string)
	filepath.Walk(templateRootDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), templateExtension) {
			retMap[info.Name()] = path
		}
		return nil
	})

	dt.DevTemplateMap = retMap
}

func (dt *DevTemplates) GetLocalTemplate(templateName string) (string, bool) {
	path, found := dt.DevTemplateMap[templateName]
	return path, found
}
