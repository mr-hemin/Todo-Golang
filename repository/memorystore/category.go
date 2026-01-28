package memorystore

import (
	"todo/entity"
)

type Category struct {
	categories []entity.Category
}

func (c Category) DoesThisUesrHaveThisCategoryID(userID, categoryID uint) bool {
	isFound := false

	for _, ctg := range c.categories {
		if ctg.ID == categoryID && ctg.UserID == userID {
			isFound = true

			break
		}
	}

	return isFound
}
