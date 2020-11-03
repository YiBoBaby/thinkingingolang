package areatree

import (
	"fmt"
	"gorm.io/gorm"
	"mysql"
	"testing"
	"time"
)

func Test_AreaTree(t *testing.T) {
	// connect mysql
	db, err := mysql.BuildMysqlConnect("192.168.9.102", 31004, "golang", "gPMg#W9fdBA%tsd9", "area_tree")
	if err != nil {
		fmt.Printf("connect mysql fail, cause:%v", err)
		return
	}

	// clean record
	CleanTable(db)

	// create M=width, height=depth tree
	width := 100
	depth := 10
	CreateTestTreeData(width, depth, db)

	rc := new(ResultCollect)
	rc.Init()

	// 测试项: 添加并删除一个 height+1 的节点 100 次
	for i := 0; i < 100; i++ {
		appendNodeId := int64(1000*(depth+1) + 1)
		// 添加一个 height+1 的节点
		node := &AreaInfo{AreaId: appendNodeId, ParentId: int64(1000*depth + 1)}
		s1 := time.Now().Unix()
		InsertAreaInfo(node, db)
		s2 := time.Now().Unix()
		rc.AddTestResult("add_tree_node", s2-s1)

		// 删除添加的节点
		s1 = time.Now().Unix()
		DeleteAreaInfo(appendNodeId, db)
		s2 = time.Now().Unix()
		rc.AddTestResult("del_tree_node", s2-s1)
	}

	// 测试项: 获取根节点后代节点 100 次
	for i := 0; i < 100; i++ {
		s1 := time.Now().Unix()
		GetDescendant(1001, db)
		s2 := time.Now().Unix()
		rc.AddTestResult("get_descendant_tree_node", s2-s1)
	}

	// 测试项: 获取祖先节点100次
	lastLevId := int64(1000*depth + 1)
	for i := 0; i < 100; i++ {
		s1 := time.Now().Unix()
		GetAncestor(lastLevId, db)
		s2 := time.Now().Unix()
		rc.AddTestResult("get_ancestor_tree_node", s2-s1)
	}

	rc.PrintStatisticsResult()
}

// 生成测试数据
func CreateTestTreeData(width, depth int, db *gorm.DB) {
	// 插入
	root := &AreaInfo{AreaId: 1001}
	InsertAreaInfo(root, db)
	BFS(root, width, 1, depth, db)

	// 删除区域节点
	DeleteAreaInfo(11, db)
}

func BFS(parent *AreaInfo, width, depth, limit int, db *gorm.DB) {
	if depth >= limit {
		return
	}
	base := 1000 * (depth + 1)

	var first *AreaInfo
	for i := 1; i <= width; i++ {
		tmp := &AreaInfo{AreaId: int64(base + i), ParentId: parent.AreaId}
		if i == 1 {
			first = tmp
		}
		InsertAreaInfo(tmp, db)
	}
	BFS(first, width, depth+1, limit, db)
}

// 清空表
func CleanTable(db *gorm.DB) {
	txDb := db.Begin()
	commitFlag := false
	defer func() {
		if commitFlag {
			txDb.Commit()
		} else {
			txDb.Rollback()
		}
	}()
	_ = txDb.Table("t_area_info").Where("1=1").Delete(nil)
	_ = txDb.Table("t_area_closure_info").Where("1=1").Delete(nil)
	commitFlag = true
}

type ResultCollect struct {
	cache map[string][]int64
}

func (rc *ResultCollect) Init() {
	rc.cache = make(map[string][]int64)
}

func (rc *ResultCollect) AddTestResult(key string, val int64) {
	if v, ok := rc.cache[key]; ok {
		rc.cache[key] = append(v, val)
	} else {
		rc.cache[key] = []int64{val}
	}
}

func (rc *ResultCollect) PrintStatisticsResult() {
	for key, val := range rc.cache {
		l := len(val)
		var t, max, min int64
		for i, ele := range val {
			if i == 0 {
				max = ele
				min = ele
			} else {
				if ele > max {
					max = ele
				}
				if ele < min {
					min = ele
				}
			}
			t += ele
		}
		fmt.Printf("测试项: %s, 次数: %d, 最大值: %d, 最小值: %d, 平均值: %f \n", key, l, max, min, float64(t)/float64(l))
	}
}
