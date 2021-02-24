package dao

import (
	"Novel/entity"
	"bytes"
	"fmt"
	"strings"
)

type DirectoryDao struct {
	TableName string
}

func (d DirectoryDao) CreateNovelDirectoryTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS 表名( 
    	id SMALLINT UNSIGNED AUTO_INCREMENT, 
		rank SMALLINT UNSIGNED NOT NULL UNIQUE COMMENT "章节序列号", 
	    chapter_name VARCHAR(50) NOT NULL COMMENT "章节名", 
	    chapter_url VARCHAR(100) NOT NULL COMMENT "章节资源url", 
	    PRIMARY KEY (id)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="小说章节目录表";`
	sql = strings.Replace(sql, "表名", d.TableName, 1)
	_, err := DB.Exec(sql)
	if err != nil {
		fmt.Println("CreateTable error, error info: " + err.Error())
		return err
	}
	fmt.Println("CreateTable success")
	return nil
}

func (d DirectoryDao) InsertAllData(chapterList []entity.ChapterEntity) error {
	sqlInsert := `insert ignore into 表名(rank, chapter_name, chapter_url) values `
	sqlInsert = strings.Replace(sqlInsert, "表名", d.TableName, 1)
	var buffer bytes.Buffer
	buffer.WriteString(sqlInsert)
	for index, chapter := range chapterList {
		s := fmt.Sprintf(`("%d", "%s", "%s")`, chapter.Rank, chapter.ChapterName, chapter.ChapterURL)
		if index == len(chapterList) - 1 {
			buffer.WriteString(s + ";")
		} else {
			buffer.WriteString(s + ",")
		}
	}
	str := buffer.String()
	_, err := DB.Exec(str)
	if err != nil {
		fmt.Println("小说目录数据插入失败！error info： " + err.Error())
		return err
	}
	fmt.Println("小说目录数据插入成功！")
	return nil
}

func (d DirectoryDao) SelectDIRByRank(rank uint64) ([]entity.ChapterEntity, error) {
	rankRange := entity.GetRankRange(rank)
	sqlSelect := fmt.Sprintf(`select rank, chapter_name, chapter_url from %s where rank between %d and %d`,
		d.TableName, rankRange[0], rankRange[len(rankRange) - 1])
	rows, err := DB.Query(sqlSelect)
	if err != nil {
		fmt.Println("查询到该rank的小说目录错误，error info：" + err.Error())
		return nil, err
	}
	var chapterList []entity.ChapterEntity
	for rows.Next() {
		var chapter entity.ChapterEntity
		err = rows.Scan(&chapter.Rank, &chapter.ChapterName, &chapter.ChapterURL)
		if err != nil {
			fmt.Println("查询到该rank的小说目录错误，error info：" + err.Error())
			return nil, err
		}
		chapterList = append(chapterList, chapter)
	}
	return chapterList, nil
}
