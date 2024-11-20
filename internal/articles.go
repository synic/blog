package internal

import (
	"encoding/json"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/synic/adamthings.me/internal/model"
)

func ParseArticles(filesystem fs.FS) ([]*model.Article, error) {
	var articles []*model.Article

	if Debug {
		log.Println("🐝 Debugging enabled, unpublished articles will be shown")
	}

	err := fs.WalkDir(filesystem, ".", func(name string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(name) != ".json" {
			return nil
		}

		var article model.Article
		data, err := fs.ReadFile(filesystem, name)

		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &article)

		if err != nil {
			return err
		}

		if Debug || article.IsPublished {
			articles = append(articles, &article)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return articles, nil
}
