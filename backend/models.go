// Struct models
package main

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Todos    []Todo
}

// type Todo struct {
// 	ID        uint   `gorm:"primaryKey"`
// 	Title     string
// 	Completed bool
//     UserID uint `json:"userID"`
// 	IsDeleted bool `gorm:"default:false"`
// }
type Todo struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	UserID    uint   `json:"userID"`
	IsDeleted bool   `gorm:"default:false" json:"isDeleted"`
}