package dao

import (
	"Novel/entity"
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

type ContentDao struct {
	TableName string
}

func (c ContentDao) CreateNovelContentTable() error {
	sql := `
	CREATE TABLE IF NOT EXISTS 表名( 
    	id SMALLINT UNSIGNED AUTO_INCREMENT, 
		rank SMALLINT UNSIGNED NOT NULL COMMENT "章节序列号", 
	    content VARCHAR(5000) NOT NULL COMMENT "章节正文", 
	    continuous TINYINT UNSIGNED NOT NULL COMMENT "该章是否过长，需要续接", 
	    PRIMARY KEY (id)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="小说章节正文表";`
	sql = strings.Replace(sql, "表名", c.TableName, 1)
	_, err := DB.Exec(sql)
	if err != nil {
		fmt.Println("CreateTable error, error info: " + err.Error())
		return err
	}
	fmt.Println("CreateTable success")
	return nil
}

func (c ContentDao) HasContent(rank uint64) bool {
	sqlSelect := fmt.Sprintf("select rank from %s where rank=%d", c.TableName, rank)
	var reRank uint16
	if err := DB.QueryRow(sqlSelect).Scan(&reRank); err == sql.ErrNoRows {
		fmt.Println("未查询到该rank的小说正文内容，可爬取网页内容!")
		return true
	} else {
		fmt.Println("查询到该rank的小说正文内容!")
		return false
	}
}

func (c ContentDao) SelectContentByRank(rank uint64) ([]entity.ChapterEntity, error) {
	rankRange := entity.GetRankRange(rank)
	sqlSelect := fmt.Sprintf("select rank,content,continuous from %s where rank between %d and %d",
		c.TableName, rankRange[0], rankRange[len(rankRange) - 1])
	rows, err := DB.Query(sqlSelect)
	if err != nil {
		fmt.Println("查询到该rank的小说正文内容错误，error info：" + err.Error())
		return nil, err
	}
	var chapterList []entity.ChapterEntity
	for rows.Next() {
		var chapter entity.ChapterEntity
		err = rows.Scan(&chapter.Rank, &chapter.Content, &chapter.Continuous)
		if err != nil {
			fmt.Println("查询到该rank的小说正文内容错误，error info：" + err.Error())
			return nil, err
		}
		chapterList = append(chapterList, chapter)
	}
	return chapterList, nil
}

func (c ContentDao) InsertAllContent(chapterList []entity.ChapterEntity) error  {
	sqlInsert := `insert ignore into 表名(rank, content, continuous) values `
	sqlInsert = strings.Replace(sqlInsert, "表名", c.TableName, 1)
	var buffer bytes.Buffer
	buffer.WriteString(sqlInsert)
	for index, chapter := range chapterList {
		s := fmt.Sprintf(`("%d", "%s", "%d")`, chapter.Rank, chapter.Content, chapter.Continuous)
		if index == len(chapterList) - 1 {
			buffer.WriteString(s + ";")
		} else {
			buffer.WriteString(s + ",")
		}
	}
	str := buffer.String()
	_, err := DB.Exec(str)
	if err != nil {
		fmt.Println("小说content数据插入失败！error info： " + err.Error())
		return err
	}
	fmt.Println("小说content数据插入成功！")
	return nil
}
