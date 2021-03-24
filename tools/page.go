package tools

import "math"

func RangeList(a, b int, c []int) []int {
	for i := a; i < b; i++ {
		c = append(c, i)
	}
	return c
}

func GetPageNum(count int64, pageSize int) int {
	pageNum := int(math.Ceil(float64(count) / float64(pageSize)))
	return pageNum
}

//总页数,当前页，前后页数
func GetPaginationData(numPages, currentPage, aroundCount int) (leftPages, rightPages []int, leftHasMore, rightHasMore bool) {

	if currentPage <= aroundCount+2 {
		leftPages = RangeList(1, currentPage, leftPages)
	} else {
		leftHasMore = true
		leftPages = RangeList(currentPage-aroundCount, currentPage, leftPages)
	}

	if currentPage >= numPages-aroundCount-1 {
		rightPages = RangeList(currentPage+1, numPages+1, rightPages)
	} else {
		rightHasMore = true
		rightPages = RangeList(currentPage+1, currentPage+aroundCount+1, rightPages)

	}

	return

}

func HasNext(page, numPages int) (hasPrevious, hasNext bool, previousPageNumber, nextPageNumber int) {
	if page == 1 && numPages == page {
		hasPrevious = false
		hasNext = false
		nextPageNumber = 1
		previousPageNumber = 1
	} else if page == 1 && numPages > page {
		hasPrevious = false
		hasNext = true
		previousPageNumber = page
		nextPageNumber = page + 1
	} else if page > 1 && numPages > page {
		hasPrevious = true
		hasNext = true
		previousPageNumber = page - 1
		nextPageNumber = page + 1
	} else if page > 1 && page == numPages {
		hasPrevious = true
		hasNext = false
		previousPageNumber = page - 1
		nextPageNumber = page
	}
	return
}
