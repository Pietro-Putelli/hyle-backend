package domain

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	TimestapModel

	ID     uint      `gorm:"primaryKey;autoIncrement;column:id"`
	Guid   uuid.UUID `gorm:"type:uuid;unique;not null;column:guid;default:uuid_generate_v4()"`
	User   *User     `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE"`
	UserID uint      `gorm:"column:user_id;not null"`

	Title  string `gorm:"column:title;not null"`
	Author string `gorm:"column:author;not null"`
}

type Topic struct {
	Topic string `gorm:"column:topic;not null" json:"topic"`
	Color string `gorm:"column:color;not null" json:"color"`

	UserID uint  `gorm:"column:user_id"`
	User   *User `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE"`
}

type BookTopic struct {
	Book   *Book `gorm:"foreignKey:BookID;references:id;constraint:OnDelete:CASCADE"`
	BookID uint  `gorm:"column:book_id;not null"`

	Topic   *Topic `gorm:"foreignKey:TopicID;references:topic;constraint:OnDelete:CASCADE"`
	TopicID string `gorm:"column:topic_id;not null"`
}

//----------------------------------------------
// Request DTOs
//----------------------------------------------

type CreateBookPickBody struct {
	Content     string `json:"content" validate:"required"`
	ContentText string `json:"contentText" validate:"required"`
	Index       *uint  `json:"index"`
}

type CreateBookBody struct {
	BookID uuid.UUID `json:"bookId"`
	Title  string    `json:"title" omitempty:"true"`
	Author string    `json:"author"`

	Pick *CreateBookPickBody `json:"pick" validate:"required"`
}

// BookListParams is a struct to get the list of books
type BookListParams struct {
	Offset int    `json:"offset" validate:"gte=0" minimum:"0" query:"offset" required:"true"`
	Limit  int    `json:"limit" validate:"gte=0" minimum:"0" query:"limit" required:"true"`
	Type   string `json:"type" query:"type" enum:"short,long" oneof:"short long" default:"short"`

	/* Search string to search in books */
	Search string `json:"search" query:"search"`

	/* Topics to filter books */
	Topics string `json:"topics" query:"topics"`
}

//----------------------------------------------
// Response DTOs
//----------------------------------------------

// ShortBookResponse is a struct to return a short version of a book
type ShortBookResponse struct {
	Guid  uuid.UUID `json:"guid"`
	Title string    `json:"title"`
}

// ShortBookFromBook converts a Book to a ShortBookResponse
func ShortBookFromBook(book *Book) *ShortBookResponse {
	return &ShortBookResponse{
		Guid:  book.Guid,
		Title: book.Title,
	}
}

// ShortBooksFromBookList converts a list of Book to a list of ShortBookResponse
func ShortBooksFromBookList(books *[]Book) *[]ShortBookResponse {
	var shortBooks []ShortBookResponse

	for _, book := range *books {
		shortBooks = append(shortBooks, *ShortBookFromBook(&book))
	}
	return &shortBooks
}

type BookTopicResponse struct {
	Topic string `json:"topic"`
	Color string `json:"color"`
}

type BookTopicListResponse struct {
	Topic string `json:"topic"`
	Color string `json:"color"`
	Count uint   `json:"count"`
}

type BookResponse struct {
	Guid       uuid.UUID `json:"guid"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	UpdatedAt  time.Time `json:"updatedAt"`
	CreatedAt  time.Time `json:"createdAt"`
	PicksCount int64     `json:"picksCount"`

	/* The Preview is a random choosen pick to display when loaded books */
	Preview *BookPickPreviewResponse `json:"preview"`

	Topics *[]Topic            `json:"topics"`
	Picks  *[]BookPickResponse `json:"picks"`
}

type CreateBookResponse interface {
	BookResponse | BookPick
}

type GetPicksParams struct {
	BookID  string `json:"bookId" validate:"required,uuid4"`
	Offset  int    `json:"offset" validate:"gte=0"`
	Limit   int    `json:"limit" validate:"gte=0"`
	OrderBy string `json:"orderBy"`

	UntilPickID string `json:"untilPickId" validate:"omitempty,uuid4"`
}

type EditBookTopicParams struct {
	Topic string `json:"topic" validate:"required"`
	Color string `json:"color" validate:"required"`
}

type EditBookOrderParams struct {
	Guid  uuid.UUID `json:"guid" validate:"required,uuid4"`
	Index uint      `json:"index" validate:"required"`
}

// EditBookBody is a struct to edit a book
type EditBookBody struct {
	BookID string                `json:"bookId" validate:"required,uuid4"`
	Title  string                `json:"title"`
	Author string                `json:"author"`
	Topics []EditBookTopicParams `json:"topics"`
	Picks  []EditBookOrderParams `json:"picks"`
}

type EditBookPickBody struct {
	BookId  string  `json:"bookId" validate:"required"`
	PickId  string  `json:"pickId" validate:"required"`
	Content string  `json:"content"`
	Text    string  `json:"text"`
	Title   *string `json:"title"`
}

// BookPickSearchKeywordMessage used as a model when send messages through SQS for generating search keywords
type BookPickSearchKeywordMessage struct {
	PickID      uint      `json:"pick_id"`
	PickContent string    `json:"content"`
	UserGuid    uuid.UUID `json:"user_guid"`
}

type PickSearchKeyword struct {
	PickID  uint   `json:"pick_id"`
	Keyword string `json:"keyword"`
	UserID  uint   `json:"user_id"`
}

// SearchGetParams Used as model for get params in semantic search
type SearchGetParams struct {
	Query  string `json:"query" validate:"required"`
	Offset int    `json:"offset" validate:"required,gte=0"`
	Limit  int    `json:"limit" validate:"required,gte=0"`

	/* The book id to search in */
	BookID string `json:"bookId" validate:"omitempty,uuid4"`
}

type SemanticSearchResponse struct {
	BookID      uuid.UUID `json:"book_id"`
	BookTitle   string    `json:"book_title"`
	PickID      uuid.UUID `json:"pick_id"`
	PickContent string    `json:"pick_content"`
	PickIndex   uint      `json:"pick_index"`
	PickTitle   string    `json:"pick_title"`
}

type SearchPickInBookResponse struct {
	PickID      uuid.UUID `json:"pick_id"`
	PickTitle   string    `json:"pick_title"`
	PickContent string    `json:"pick_content"`
	PickIndex   uint      `json:"pick_index"`
}

type SaveBookBody struct {
	BookID uuid.UUID `json:"bookId" validate:"required,uuid4"`
}

type SharpPickParams struct {
	Text string `json:"text" validate:"required"`
}

type GenerateKeywordDetailParams struct {
	Keyword string `json:"keyword" validate:"required"`
}

type TranslateWordParams struct {
	Word string `json:"word" validate:"required"`
	Lang string `json:"lang" validate:"required"`
}

type BookPickPushNotification struct {
	BookID  uuid.UUID `json:"book_id"`
	Content string    `json:"content"`
}
