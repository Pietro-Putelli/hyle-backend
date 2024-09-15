package domain

import "github.com/google/uuid"

type BookPick struct {
	TimestapModel

	ID   uint      `gorm:"primaryKey;autoIncrement;column:id"`
	Guid uuid.UUID `gorm:"type:uuid;unique;not null;column:guid;default:uuid_generate_v4()"`

	Book   *Book `gorm:"foreignKey:BookID;references:id;constraint:OnDelete:CASCADE"`
	BookID uint  `gorm:"column:book_id;not null"`

	User   *User `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE"`
	UserID uint  `gorm:"column:user_id;not null"`

	/* Pick content as JSON */
	Content string `gorm:"column:content;not null"`
	/* Pick content as Plain Text for searching into and preview */
	ContentText string `gorm:"column:content_text;not null"`
	/* Pick title that can be null */
	Title string `gorm:"column:title"`

	Index uint `gorm:"column:index;not null"`
}

type BookPickSearchKeyword struct {
	Pick   *BookPick `gorm:"foreignKey:PickID;references:id;constraint:OnDelete:CASCADE"`
	PickID uint      `gorm:"column:pick_id;not null"`
	Keywod string    `gorm:"column:keyword;not null"`
}

//----------------------------------------------
// Request DTOs
//----------------------------------------------

type DeleteBookPickPath struct {
	BookID string `path:"bookId" validate:"required,uuid4"`
	PickID string `path:"pickId" validate:"required,uuid4"`
}

//----------------------------------------------
// Response DTOs
//----------------------------------------------

type BookPickResponse struct {
	Guid    uuid.UUID `json:"guid"`
	Content string    `json:"content"`
	Index   uint      `json:"index"`
	Title   string    `json:"title"`
}

type BookPickPreviewResponse struct {
	Guid        uuid.UUID `json:"guid"`
	ContentText string    `json:"content"`
}
