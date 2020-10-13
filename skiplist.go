package skiplist

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

/**
 * SkipNode
 * l  4  30  -------------------->  nil
 * e  3  30  ---> 50  ----------->  nil
 * v  2  30  ---> 50  ---->70 ----> nil
 * e  1  30 ->40->50 ->60->70 ->90->nil
 * l
 * skipList本身是排好序的
 **/
/**
 */
type SkipNode struct {
	Key     uint32
	Val     interface{}
	Forward []*SkipNode // 当前节点的下一个节点(每一层的下一节点)，例如30的forward就是 40, 50, 50, nil
	Level   int         // 所在层数
}

// NewNode 创建新节点 searchKey: key, value: val, createLevel: 当前节点所有在level, maxLevel: 整个skipList的最大深度
func NewNode(searchKey uint32, value interface{}, createLevel int, maxLevel int) *SkipNode {
	// every forward prepare a maxLevel empty point first
	forwardEmpty := make([]*SkipNode, maxLevel)
	for i := 0; i <= maxLevel-1; i++ {
		forwardEmpty[i] = nil
	}
	return &SkipNode{Key: searchKey, Val: value, Forward: forwardEmpty, Level: createLevel}
}

type SkipList struct {
	Header      *SkipNode
	MaxLevel    int // 最大深度
	Propobility float32
	Level       int // 当前跳表的所有层数
}

const (
	DefaultMaxLevel    int     = 15   // maximal level allow to create in the skiplist
	DefaultPropobility float32 = 0.25 // default propobility
)

// NewSkipList 初始化 SkipList
func NewSkipList() *SkipList {
	newList := &SkipList{Header: NewNode(0, "header", 1, DefaultMaxLevel), Level: 1}
	newList.MaxLevel = DefaultMaxLevel
	newList.Propobility = DefaultPropobility
	return newList
}

// 随机生成一个propobility, [0,1）
func randomP() float32 {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Float32()
}

// 可以由用户指定最大层数
func (sl *SkipList) SetMaxLevel(maxLevel int) {
	sl.MaxLevel = maxLevel
}

// RandomLevel 随机生成跳表的层数，因为自己也不知道这个跳表要多少层合适
func (sl *SkipList) RandomLevel() int {
	level := 1
	for randomP() < sl.Propobility && level < sl.MaxLevel {
		level++
	}
	return level
}

// Search a element by search key and return the interface{}
func (sl *SkipList) Search(searchKey uint32) (interface{}, error) {
	currentNode := sl.Header

	// 从最上层开始往下找
	for i := sl.Level - 1; i >= 0; i-- {
		for currentNode.Forward[i] != nil && currentNode.Forward[i].Key <= searchKey {
			if currentNode.Forward[i].Key == searchKey {
				// 在第i层已经找到
				return currentNode.Forward[i].Val, nil
			} else {
				// 在第i层继续右移
				currentNode = currentNode.Forward[i]
			}
		}
	}
	// 到最下一层了
	if currentNode.Forward[0] != nil && currentNode.Forward[0].Key == searchKey {
		return currentNode.Forward[0].Val, nil
	}
	return nil, errors.New("Not Found.")
}

// Insert: Insert a search key and its value which can by interface{}
func (sl *SkipList) Insert(searchKey uint32, value interface{}) {
	// level := sl.RandomLevel()
	// if level > sl.Level {
	// 	sl.Level = level
	// }
	updateList := make([]*SkipNode, sl.MaxLevel)
	currentNode := sl.Header

	// QuickSearch in forward list
	for i := sl.Header.Level - 1; i >= 0; i-- {
		for currentNode.Forward[i] != nil && currentNode.Forward[i].Key < searchKey {
			currentNode = currentNode.Forward[i]
		}
		updateList[i] = currentNode
	}

	//Step to next node. (which is the target insert location)
	currentNode = currentNode.Forward[0]

	if currentNode != nil && currentNode.Key == searchKey {
		currentNode.Val = value
	} else {
		newLevel := sl.RandomLevel()
		if newLevel > sl.Level {
			for i := sl.Level + 1; i <= newLevel; i++ {
				updateList[i-1] = sl.Header
			}
			sl.Level = newLevel //This is not mention in cookbook pseudo code
			sl.Header.Level = newLevel
		}

		newNode := NewNode(searchKey, value, newLevel, sl.MaxLevel) //New node
		for i := 0; i <= newLevel-1; i++ {                          //zero base
			newNode.Forward[i] = updateList[i].Forward[i]
			updateList[i].Forward[i] = newNode
		}
	}
}

//Delete: Delete element by search key
func (b *SkipList) Delete(searchKey uint32) error {
	updateList := make([]*SkipNode, b.MaxLevel)
	currentNode := b.Header

	//Quick search in forward list
	for i := b.Header.Level - 1; i >= 0; i-- {
		for currentNode.Forward[i] != nil && currentNode.Forward[i].Key < searchKey {
			currentNode = currentNode.Forward[i]
		}
		updateList[i] = currentNode
	}

	//Step to next node. (which is the target delete location)
	currentNode = currentNode.Forward[0]

	if currentNode.Key == searchKey {
		for i := 0; i <= currentNode.Level-1; i++ {
			if updateList[i].Forward[i] != nil && updateList[i].Forward[i].Key != currentNode.Key {
				break
			}
			updateList[i].Forward[i] = currentNode.Forward[i]
		}

		for currentNode.Level > 1 && b.Header.Forward[currentNode.Level] == nil {
			currentNode.Level--
		}

		//free(currentNode)  //no need for Golang because GC
		currentNode = nil
		return nil
	}
	return errors.New("Not found")
}

//DisplayAll: Display current SkipList content in console, will also print out the linked pointer.
func (b *SkipList) DisplayAll() {
	fmt.Printf("\nhead->")
	currentNode := b.Header

	//Draw forward[0] base
	for {
		fmt.Printf("[key:%d][val:%v]->", currentNode.Key, currentNode.Val)
		if currentNode.Forward[0] == nil {
			break
		}
		currentNode = currentNode.Forward[0]
	}
	fmt.Printf("nil\n")

	fmt.Println("---------------------------------------------------------")
	currentNode = b.Header
	//Draw all data node.
	for {
		fmt.Printf("[node:%d], val:%v, level:%d ", currentNode.Key, currentNode.Val, currentNode.Level)

		if currentNode.Forward[0] == nil {
			break
		}

		for j := currentNode.Level - 1; j >= 0; j-- {
			fmt.Printf(" fw[%d]:", j)
			if currentNode.Forward[j] != nil {
				fmt.Printf("%d", currentNode.Forward[j].Key)
			} else {
				fmt.Printf("nil")
			}
		}
		fmt.Printf("\n")
		currentNode = currentNode.Forward[0]
	}
	fmt.Printf("\n")
}
