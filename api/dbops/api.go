package dbops

import (
	"database/sql"
	"log"
	"rushflow/api/defs"
	"rushflow/api/utils"
	"time"
)

// 增加用户
func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO users(login_name, pwd) VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

// 查询用户
func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE login_name=?")
	if err != nil {
		log.Printf("GetUserCredential error:%s", err)
		return "", err
	}
	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	defer stmtOut.Close()
	return pwd, nil
}

// 删除用户
func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM users WHERE login_name=? AND pwd=?")
	if err != nil {
		log.Printf("DeleteUser error:%s", err)
		return err
	}
	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

// 新增视频
func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	// create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}
	t := time.Now()
	cTime := t.Format("Jan 02 2006,15:04:05")
	stmtIns, err := dbConn.Prepare("INSERT INTO video_info (id,author_id,name,display_ctime) VALUES(?,?,?,?)")
	if err != nil {
		return nil, err
	}
	_, err = stmtIns.Exec(vid, aid, name, cTime)
	if err != nil {
		return nil, err
	}
	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: cTime}
	defer stmtIns.Close()
	return res, nil
}

// 查询视频
func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT author_id, name, display_ctime FROM video_info WHERE id=?")
	var aid int
	var dct string
	var name string
	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	defer stmtOut.Close()
	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: dct}
	return res, nil
}

// 删除视频
func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

// 新增评论
func AddNewComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}
	stmtIns, err := dbConn.Prepare("INSERT INTO comments (id,video_id,author_id,content) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

// 列出视频下所有评论
func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	// uers join comments -> author_id , video_id -> comment
	stmtOut, err := dbConn.Prepare(`SELECT comments.id, users.login_name, comments.content FROM comments
		INNER JOIN users ON comments.author_id = users.id
		WHERE comments.video_id=? AND comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?)`)
	var res []*defs.Comment
	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var id, authorName, content string
		if err := rows.Scan(&id, &authorName, &content); err != nil {
			return res, err
		}
		c := &defs.Comment{Id: id, VideoId: vid, Author: authorName, Content: content}
		res = append(res, c)
	}
	defer stmtOut.Close()
	return res, nil
}