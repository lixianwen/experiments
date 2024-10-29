package models

import (
	"fmt"
	"gdemo/hash"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string     `json:"name" gorm:"unique"`
	Email      string     `json:"email" binding:"email"`
	Age        int        `json:"age" binding:"min=0,max=100" gorm:"default:18"`
	Password   string     `json:"password,omitempty" gorm:"size:255"`
	CreditCard CreditCard `json:"card" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreditCard struct {
	gorm.Model
	Number string `json:"cid"`
	UserID uint   `json:"-"`
}

// All returns all the users in the table users.
func (user *User) All() ([]*User, error) {
	var sl []*User
	if result := db.Joins("CreditCard").Find(&sl); result.Error != nil {
		return nil, result.Error
	}

	return sl, nil
}

// Create creates a user instance.
func (user *User) Create() error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	if result := db.Create(user); result.Error != nil {
		return result.Error
	} else if result.RowsAffected != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", result.RowsAffected)
	}

	return nil
}

func (user *User) Get(id int) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	// db.InnerJoins("CreditCard", db.Where("users.id = ?", id)).First(user)
	if result := db.Joins("CreditCard").First(user, id); result.Error != nil {
		return result.Error
	}

	return nil
}

func (user *User) GetByName(name string) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	if err := db.First(user, "name = ?", name).Error; err != nil {
		return err
	}

	return nil
}

func (user *User) Update(id int) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	var oldUser User
	// has one association
	db.Joins("CreditCard").First(&oldUser, id)

	user.ID = uint(id)
	// critical step for update association
	user.CreditCard.ID = oldUser.CreditCard.ID
	if err := db.Model(user).Select("*").Omit("created_at", "password").Updates(user).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Select("*").Omit("created_at", "password").Updates(user).Error; err != nil {
		return err
	}

	return nil
}

func (user *User) Delete(id int) error {
	user.ID = uint(id)
	if err := db.Unscoped().Select("CreditCard").Delete(user).Error; err != nil {
		return err
	}

	return nil
}
