package entity

// ChapterEntity, 小说章节对象，定义了章节的基础信息
type ChapterEntity struct {
	Rank uint16			`json:"rank"`			// 章节的序
	ChapterName string	`json:"chapter_name"`	// 章节名
	ChapterURL string	`json:"chapter_url"`	// 章节的资源url
	Content string		`json:"content"`		// 章节的正文内容
	Continuous uint8	`json:"continuous"`		// 防止一个章节的content过长，分几条数据存储,为 0 表示无后续
}

// 章节rank值对应的区间范围
var CHAPTER_RANGE uint64 = 10
// 章节content长度的限制
var COUNT = 4990

// GetRankRange, 函数名
// 返回章节rank值对应的区间（处于哪10个章节中，如：1-10、11-20等）
// @param rank: 章节的序
func GetRankRange(rank uint64) (rankRange []uint64) {
	ra := (rank - 1) / CHAPTER_RANGE * CHAPTER_RANGE
	for i := 1; uint64(i) <= CHAPTER_RANGE; i++ {
		rankRange = append(rankRange, ra + uint64(i))
	}
	return rankRange
}

// ContentPaging, 函数名
//
func ContentPaging(content string, chapter ChapterEntity, list *[]ChapterEntity) {
	contentRune := []rune(content)
	for left, right := 0, COUNT; len(contentRune) > left; left, right = left +COUNT, right +COUNT {
		if len(contentRune) <= right {
			chapter.Content = string(contentRune[left:])
			if left != 0 {
				chapter.Continuous = uint8(right / COUNT)
			}
			*list = append(*list, chapter)
			break
		} else {
			chapter.Content = string(contentRune[left:right])
			chapter.Continuous = uint8(right / COUNT)
			*list = append(*list, chapter)
		}
	}
}
