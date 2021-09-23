package utils

import "gorm.io/gorm"

func Paginate(meta *MetaData) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case meta.Limit > 100:
			meta.Limit = 100
		case meta.Limit <= 0:
			meta.Limit = 10
		}

		offset := (meta.CurrentPage - 1) * meta.Limit
		return db.Offset(offset).Limit(meta.Limit)
	}
}
