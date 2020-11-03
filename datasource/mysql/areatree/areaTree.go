package areatree

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

// 邻接表
type AreaInfo struct {
	Id         int64
	AreaId     int64
	ParentId   int64
	AreaName   string
	CreateTime *time.Time
	UpdateTime *time.Time
}

// 闭包表
type AreaClosure struct {
	Id         int64
	Ancestor   int64
	Descendant int64
	Depth      int
}

// 插入区域节点
func InsertAreaInfo(areaInfo *AreaInfo, db *gorm.DB) {
	txDb := db.Begin()
	commitFlag := false
	defer func() {
		if commitFlag {
			txDb.Commit()
		} else {
			txDb.Rollback()
		}
	}()

	tx1 := txDb.Table("t_area_info").Select("area_id", "parent_id", "area_name").Create(areaInfo)
	sql := fmt.Sprintf("INSERT INTO t_area_closure_info (ancestor, descendant,depth)" +
		" SELECT t.ancestor, ?, t.depth+1 FROM t_area_closure_info AS t" +
		" WHERE t.descendant = ?" +
		" UNION ALL" +
		" SELECT ?, ?, 1")
	tx2 := txDb.Exec(sql, areaInfo.AreaId, areaInfo.ParentId, areaInfo.AreaId, areaInfo.AreaId)
	if tx1.RowsAffected > 0 && tx2.RowsAffected > 0 {
		commitFlag = true
	}
}

// 删除区域节点
func DeleteAreaInfo(areaId int64, db *gorm.DB) {
	txDb := db.Begin()
	commitFlag := false
	defer func() {
		if commitFlag {
			txDb.Commit()
		} else {
			txDb.Rollback()
		}
	}()

	tx1 := txDb.Table("t_area_info").Where("area_id IN (SELECT descendant FROM t_area_closure_info WHERE ancestor = ?)", areaId).Delete(nil)
	tx2 := txDb.Exec("DELETE t1 FROM `t_area_closure_info` t1, `t_area_closure_info` t2 where t1.id = t2.id and t2.ancestor = ?", areaId)
	if tx1.RowsAffected > 0 && tx2.RowsAffected > 0 {
		commitFlag = true
	}
}

// 获取当前节点的后代节点
func GetDescendant(areaId int64, db *gorm.DB) {
	rs := make([]int64, 0, 0)
	_ = db.Table("t_area_closure_info").Select("descendant").Where("ancestor = ?", areaId).Find(&rs)
}

// 获取当前节点的祖先节点
func GetAncestor(areaId int64, db *gorm.DB) {
	rs := make([]int64, 0, 0)
	_ = db.Table("t_area_closure_info").Select("DISTINCT(ancestor)").Where("descendant = ?", areaId).Find(&rs)
}

// 待移动节点 移动=删除+添加
func MoveNodeTo() {

}
