// yourproject/entity/user.go
package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Number   string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"default:'user'"`
	Status   string `gorm:"default:'active'"`
	// ReferalDetails ReferalDetails
	CreatedAt time.Time
}

type ReferalDetails struct {
	gorm.Model
	UserID      uint   `json:"userID" gorm:"index;foreignKey:UserID"`
	Count       uint   `json:"count"`
	ReferalCode string `json:"referalCode" gorm:"unique_index"`
}

type TokenUser struct {
	Users        Claims
	AccessToken  string
	RefreshToken string
}

type UserAddress struct {
	gorm.Model
	UserID      uint   `json:"userID"  gorm:"index;foreignKey:UserID"`
	Street      string `gorm:"not null"  json:"street" form:"street" binding:"required"`
	City        string `gorm:"not null" json:"city" form:"city" binding:"required"`
	State       string `gorm:"not null" json:"state" form:"state" binding:"required"`
	PostalCode  string `gorm:"not null" json:"postalcode" form:"postalcode" binding:"required"`
	Country     string `gorm:"not null" json:"country" form:"country" binding:"required"`
	IsPrimary   string `gorm:"default:'false'"`
	PhoneNumber string `json:"phone_number" gorm:"not null" form:"phonenumber" binding:"required"`
}

type UserDetail struct {
	UserName    string `form:"username" binding:"required"`
	Email       string `form:"email" binding:"required"`
	PhoneNumber string `form:"mobile" binding:"required"`
}
type UpdatePassword struct {
	Password        string `form:"password" binding:"required"`
	NewPassword     string `form:"newpassword" binding:"required"`
	ConfirmPassword string `form:"confirmpassword" binding:"required"`
}

type Compare struct {
	ID       uint
	Password string
	Role     string
	Username string
	Email    string
	Status   string
}

// type Product struct {
// 	Product_ID   primitive.ObjectID `bson:"_id"`
// 	Product_Name *string            `json:"product_name"`
// 	Seller_ID
// 	Category_ID

// 	Price *uint64 `json:"price"`
// 	Discount_Price
// 	Rating *uint8  `json:"rating"`
// 	Image  *string `json:"image"`
// }

//	type Category struct {
//		gorm.Model
//		CategoryName string `gorm:"unique;not null"`
//		Status       string `gorm:"default:'listed'"`
//		CreatedAt    time.Time
//	}

type OrderReq struct {
	CartID        string `json:"cartID"`
	AddressID     string `json:"addressID"`
	PaymentMethod string `json:"paymentMethod"`
	CouponCode    string `json:"couponcode"`
}
type Updatecart struct {
	ID       string `json:"id"`
	Quantity string `json:"quantity"`
}

type UserRequest struct {
	ID      int    `json:"id"`
	Request string `json:"request"`
}
