package dao

import (
	"Novel/entity"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
)


func CreateNovelInfoTable() {
	sql := `
	CREATE TABLE IF NOT EXISTS novel_info( 
    	id MEDIUMINT UNSIGNED AUTO_INCREMENT, 
		novel_id VARCHAR(15) NOT NULL UNIQUE COMMENT "小说ID", 
	    title VARCHAR(20) NOT NULL COMMENT "小说名", 
	    author VARCHAR(20) NOT NULL COMMENT "小说作者", 
	    category VARCHAR(10) DEFAULT "" COMMENT "小说种类",
	    status VARCHAR(10) DEFAULT "" COMMENT "小说状态",
	    description VARCHAR(200) DEFAULT "" COMMENT "小说描述",
	    update_info VARCHAR(50) DEFAULT "" COMMENT "最近更新",
	    book_img VARCHAR(100) DEFAULT "" COMMENT "小说图片url",
	    sources VARCHAR(5) DEFAULT "001" COMMENT "小说资源标识符",
	    source_url VARCHAR(100) DEFAULT "" COMMENT "小说资源url",
	    create_time DATETIME,
	    update_time DATETIME,
	    PRIMARY KEY (id)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT="小说基础信息表";`

	_, err := DB.Exec(sql)
	if err != nil {
		fmt.Println("CreateTable error, error info: " + err.Error())
		os.Exit(0)
	}
	fmt.Println("CreateTable success")
}

var InfoDao *novelInfoDao

func GetInfoDaoInstance() *novelInfoDao {
	if InfoDao == nil {
		InfoDao = &novelInfoDao{}
	}
	return InfoDao
}

type novelInfoDao struct {
}

func (p *novelInfoDao) HasNovelInfo(novelID string) bool {
	var id string
	err := DB.QueryRow("select novel_id from novel_info where novel_id=?", novelID).Scan(&id)
	if err == sql.ErrNoRows {
		fmt.Println("未查询到该小说信息，可插入novel_info表中!")
		return true
	} else {
		fmt.Println("查询到该小说信息，可更新数据!")
		return false
	}
}

func (p *novelInfoDao) QuerySourceByNovelID(novelID string) (source, sourceURL string, err error) {
	row := DB.QueryRow("select sources, source_url from novel_info where novel_id=?", novelID)
	if err = row.Scan(&source, &sourceURL); err != nil {
		fmt.Println("未查询到该小说信息!")
		return "", "", err
	}
	return source, sourceURL, nil
}

func (p *novelInfoDao) InsertNovelInfo(info entity.NovelInfoEntity) int64 {
	sqlInsert := `
	INSERT INTO 
	    novel_info(novel_id, title, author, category, status, description, update_info, book_img, sources, source_url, create_time) 
		VALUE(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := DB.Exec(sqlInsert, info.ID, info.Title, info.Author, info.Category, info.Status, info.Description, info.Update, info.BookImg, info.Source, info.SourceURL, time.Now())
	if err != nil {
		fmt.Println("Insert novel_info error, error info: " + err.Error())
		return 0
	}
	if id, err := result.LastInsertId(); err != nil {
		fmt.Println("Insert novel_info error, error info: " + err.Error())
		return 0
	} else {
		return id
	}
}

func (*novelInfoDao) UpdateNovelInfo(info entity.NovelInfoEntity) int64 {
	sqlUpdate := `
	UPDATE 表名 
	    SET status=?, update_info=?, source_url=?, update_time=? 
	    WHERE novel_id=?`
	sqlUpdate = strings.Replace(sqlUpdate, "表名", "novel_info", 1)
	result, err := DB.Exec(sqlUpdate, info.Status, info.Update, info.SourceURL, time.Now(), info.ID)
	if err != nil {
		fmt.Println("update novel_info error, error info: " + err.Error())
		return 0
	}
	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("update novel_info error, error info: " + err.Error())
		return 0
	} else {
		fmt.Println("RowsAffected: ", rows)
		return rows
	}
}

