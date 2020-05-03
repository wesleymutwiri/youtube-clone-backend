package seed

import (
	"api/youtube-clone-backend/api/models"
	"log"

	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Username: "Steven Victor",
		Email:    "steven@email.com",
		Password: "password",
	},
	models.User{
		Username: "Martin Luther",
		Email:    "luther@email.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

	}
}
