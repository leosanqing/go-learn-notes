package main

import (
	"github.com/jinzhu/gorm"

	//
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Id       int `gorm:"primary_key;auto_increment"`
	Name     string
	Password string
	AddTime  int64
	Status   int
	Mobile   string
	Avatar   string
}

type Advert struct {
	Id int `gorm:"primary_key;auto_increment"`
	Title string
	SubTitle string
	ChannelId int
	Img string
	Sort string
	AddTime int64
	Url string
	Status int
}
func main() {
	db, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/fyouku?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	db.SingularTable(true)

	// 迁移 schema
	//db.AutoMigrate(&Advert{})


	//db.Create(&Advert{Name: "D42"})

	// Read
	var user []Advert
	db.First(&user) // 根据整形主键查找

	//// Update - 将 user 的 price 更新为 200
	//db.Model(&user).Update("Price", 200)
	//// Update - 更新多个字段
	//db.Model(&user).Updates(User{Price: 200, Code: "F42"}) // 仅更新非零值字段
	//db.Model(&user).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 user
	db.Delete(&user, 1)
}
