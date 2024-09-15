package book

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/sqs"
	"github.com/pietro-putelli/feynman-backend/internal/user"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
	"github.com/pietro-putelli/feynman-backend/langchain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ Service = (*serviceImpl)(nil)

type Service interface {
	// GetBookByGuid Get only book entry by guid
	GetBookByGuid(guid uuid.UUID) (*domain.Book, error)

	// GetCompleteBookByGuid Get complete book by guid (with the initial picks) and user
	GetCompleteBookByGuid(userID, bookID uuid.UUID) (*domain.BookResponse, error)

	DeleteBook(userID, bookID uuid.UUID) error

	DeleteBookPick(userID uuid.UUID, params *domain.DeleteBookPickPath) (bool, error)

	CreateBookPick(userID uuid.UUID, data *domain.CreateBookBody) (any, error)

	// GetBooks Get all user's books (UI:Homepage)
	GetBooks(userID uuid.UUID, params *domain.BookListParams) ([]domain.BookResponse, error)

	// GetShortBooksList Get short formatted books for search list
	GetShortBooksList(userID uuid.UUID, params *domain.BookListParams) ([]domain.ShortBookResponse, error)

	// GetPicksByBook Get picks by bookID
	GetPicksByBook(params *domain.GetPicksParams) ([]domain.BookPickResponse, error)

	// GetUserBooksTopics Get all topics related to the user's books
	GetUserBooksTopics(userID uuid.UUID) ([]domain.BookTopicListResponse, error)

	// EditBook Edit book's properties: (title, author, topics, pick's order)
	EditBook(userID uuid.UUID, params *domain.EditBookBody) error

	// EditBookPick Edit book pick properties
	EditBookPick(userID uuid.UUID, body *domain.EditBookPickBody) error

	// AddPickKeywords Add pick's keywords in database
	AddPickKeywords(userID uuid.UUID, pickID uint, keywords []string) error

	// SemanticSearch Semantic Search across all picks
	SemanticSearch(userID uuid.UUID, params *domain.SearchGetParams) ([]domain.SemanticSearchResponse, error)

	// Search pick in a specific book
	SearchPickInBook(params *domain.SearchGetParams) ([]domain.SearchPickInBookResponse, error)

	// SaveBook Save book to userID's account with book's guid
	SaveBook(userID uuid.UUID, book *domain.SaveBookBody) (*domain.BookResponse, error)
}

type serviceImpl struct {
	db          *gorm.DB
	userService user.Service
}

// NewService creates a new book service
func NewService(db *gorm.DB, userService user.Service) Service {
	return &serviceImpl{
		db:          db,
		userService: userService,
	}
}

//---------------------------------------------------------------------
// Service Implementation
//---------------------------------------------------------------------

func (service *serviceImpl) GetBookByGuid(guid uuid.UUID) (*domain.Book, error) {
	book := domain.Book{}

	if err := service.db.Where("guid = ?", guid).First(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (service *serviceImpl) GetCompleteBookByGuid(userID, bookID uuid.UUID) (*domain.BookResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	completeBook := domain.BookResponse{}

	if err := service.db.Transaction(func(tx *gorm.DB) error {
		book := domain.Book{}

		if err := tx.Model(&domain.Book{}).Where("guid = ?", bookID).First(&book).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Debug(fmt.Sprintf("Book for user %s with guid %s not found", userID, bookID))
				return err
			}
			return err
		}

		var topics []domain.Topic
		tx.Table("topics").
			Select("topics.topic, topics.color").
			Joins("JOIN book_topics ON topics.id = book_topics.topic_id").
			Where("book_topics.book_id = ?", book.ID).
			Scan(&topics)

		var picks []domain.BookPickResponse
		tx.Model(&domain.BookPick{}).Where("book_id = ?", book.ID).Order("index DESC").Limit(3).Find(&picks)

		preview := domain.BookPickPreviewResponse{}
		tx.Model(&domain.BookPick{}).Where("book_id = ?", book.ID).Order("RANDOM()").Limit(1).First(&preview) // why random?

		var picksCount int64
		tx.Model(&domain.BookPick{}).Where("book_id = ?", book.ID).Count(&picksCount)

		completeBook = domain.BookResponse{
			Guid:       book.Guid,
			Title:      book.Title,
			Author:     book.Author,
			UpdatedAt:  book.UpdatedAt,
			CreatedAt:  book.CreatedAt,
			PicksCount: picksCount,
			Preview:    &preview,
			Topics:     &topics,
			Picks:      &picks,
		}

		return nil
	}); err != nil {
		logger.Error("Failed to get complete book", zap.Error(err))
		return nil, err
	}

	return &completeBook, nil
}

func (service *serviceImpl) DeleteBook(userID, bookID uuid.UUID) error {

	err := service.db.Table("books").Where("user_id = (SELECT id FROM users WHERE guid = ?) AND guid = ?", userID, bookID).Delete(&domain.Book{}).Error
	return err
}

func (service *serviceImpl) DeleteBookPick(userID uuid.UUID, params *domain.DeleteBookPickPath) (bool, error) {
	bookID := params.BookID
	pickID := params.PickID

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var isLastPick bool = false

	err := service.db.Transaction(func(tx *gorm.DB) error {
		book := domain.Book{}
		service.db.Model(&domain.Book{}).Where("guid = ?", bookID).First(&book)

		var picksCount int64
		service.db.Model(&domain.BookPick{}).Where("book_id = ?", book.ID).Count(&picksCount)

		isLastPick = picksCount == 1

		var err error

		if isLastPick {
			err = service.db.Model(&domain.Book{}).Where("id = ?", book.ID).Delete(&book).Error
		} else {
			/* Delete the pick */
			pickToDelete := domain.BookPick{}

			err := service.db.Model(&domain.BookPick{}).Where("book_id = ? AND guid = ?", book.ID, pickID).First(&pickToDelete).Delete(&pickToDelete).Error
			if err != nil {
				return err
			}

			/* Update all indexes */
			service.db.Model(&domain.BookPick{}).Where("book_id = ? AND index > ?", book.ID, pickToDelete.Index).Update("index", gorm.Expr("index - 1"))
		}

		return err
	})

	return isLastPick, err
}

func (service *serviceImpl) CreateBookPick(userID uuid.UUID, data *domain.CreateBookBody) (any, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var response any = nil

	user, _ := service.userService.GetUserByGuid(userID)

	err := service.db.Transaction(func(tx *gorm.DB) error {
		/* 1. Check if BookID is empty, if so create a new book to which associate the pick to */

		if data.BookID == uuid.Nil {
			newBook := domain.Book{
				UserID: user.ID,
				Title:  data.Title,
				Author: data.Author,
			}

			/* Get the author from the title using Google Books API */
			if data.Author == "" {
				author, err := utility.GetAuthorFromTitle(data.Title)
				if err != nil {
					logger.Error("Failed to get author from title (Google Book API)", zap.Error(err))
				}

				newBook.Author = author
			}

			if err := tx.Create(&newBook).Error; err != nil {
				return err
			}

			newPick := domain.BookPick{
				BookID:      newBook.ID,
				Content:     data.Pick.Content,
				ContentText: data.Pick.ContentText,
				Index:       0,
				UserID:      user.ID,
			}

			if err := tx.Create(&newPick).Error; err != nil {
				return err
			}

			message := domain.BookPickSearchKeywordMessage{
				PickID:      newPick.ID,
				PickContent: newPick.ContentText,
				UserGuid:    userID,
			}

			err := sqs.SendMessage(sqs.QueueNames.PickKeywords, message)
			if err != nil {
				return err
			}

			topics, err := langchain.GenerateBookTopics(newPick.ContentText)
			if err != nil {
				logger.Error("Failed to generate book topics", zap.Error(err))
				return err
			}

			/* 1. Get random colors for the topics */
			topicsLen := len(topics)
			colors := []string{}

			err = tx.Table("topic_colors").Select("color").Order("RANDOM()").Limit(topicsLen).Pluck("color", &colors).Error
			if err != nil {
				logger.Error("Failed to get colors", zap.Error(err))
				return err
			}

			/* 2. Insert topics into topics table if they don't exist yet */
			for index, topic := range topics {
				err := tx.Exec("INSERT INTO topics (user_id, topic, color) VALUES (?, ?, ?) ON CONFLICT (user_id, topic) DO NOTHING", user.ID, topic, colors[index]).Error
				if err != nil {
					logger.Error("Failed to insert topic", zap.Error(err))
					return err
				}
			}

			/* 3. Select from topics all entries that match the topics list */
			var topicIDs []uint

			err = tx.Table("topics").Select("id").Where("user_id = ? AND topic IN (?)", user.ID, topics).Pluck("id", &topicIDs).Error
			if err != nil {
				logger.Error("Failed to get topic IDs", zap.Error(err))
				return err
			}

			/* 4. Insert into book_topics the book_id and the topic_id */
			for _, topicID := range topicIDs {
				err := tx.Exec("INSERT INTO book_topics (book_id, topic_id) VALUES (?, ?)", newBook.ID, topicID).Error
				if err != nil {
					logger.Error("Failed to insert into book_topics", zap.Error(err))
					return err
				}
			}

			/* 5. Get all topics for the book */
			bookTopics := []domain.Topic{}

			err = tx.Model(&domain.Topic{}).Select("topic, color").Joins("JOIN book_topics ON topics.id = book_topics.topic_id").Where("book_topics.book_id = ?", newBook.ID).Scan(&bookTopics).Error
			if err != nil {
				logger.Error("Failed to get book topics", zap.Error(err))
				return err
			}

			response = &domain.BookResponse{
				Guid:       newBook.Guid,
				Title:      newBook.Title,
				Author:     newBook.Author,
				UpdatedAt:  newBook.UpdatedAt,
				CreatedAt:  newBook.CreatedAt,
				PicksCount: 1,
				Preview:    &domain.BookPickPreviewResponse{Guid: newPick.Guid, ContentText: newPick.ContentText},
				Picks: &[]domain.BookPickResponse{
					{
						Guid:    newPick.Guid,
						Content: newPick.Content,
						Index:   0,
					},
				},
				Topics: &bookTopics,
			}
		} else {
			/* 2. If the book exists and is owned by the userID, just add the pick to it */

			book := domain.Book{}

			if err := tx.Model(&domain.Book{}).Where("guid = ? AND user_id = ?", data.BookID, user.ID).First(&book).Error; err != nil {
				return err
			}

			var index uint

			var lastPickIndex uint64
			tx.Model(&domain.BookPick{}).Where("book_id = ?", book.ID).Order("index desc").Limit(1).Pluck("index", &lastPickIndex)
			/* Add 1 because the new item is not yet in the database */
			lastPickIndex++

			if data.Pick.Index == nil {
				index = uint(lastPickIndex)
			} else {
				index = *data.Pick.Index
			}

			newPick := domain.BookPick{
				BookID:      book.ID,
				Content:     data.Pick.Content,
				ContentText: data.Pick.ContentText,
				/* This is the new index that the pick must assume */
				Index:  index,
				UserID: user.ID,
			}

			message := domain.BookPickSearchKeywordMessage{
				PickID:      newPick.ID,
				PickContent: newPick.ContentText,
				UserGuid:    user.Guid,
			}

			err := sqs.SendMessage(sqs.QueueNames.PickKeywords, message)
			if err != nil {
				return err
			}

			// /* If the index is not the last one, we need to update the indexes of the following picks to keep track of their order */

			// if lastPickIndex != uint64(index) {}
			tx.Model(&domain.BookPick{}).Where("book_id = ? AND index >= ?", book.ID, index).Update("index", gorm.Expr("index + 1"))

			if err := tx.Create(&newPick).Error; err != nil {
				logger.Error("Failed to create pick", zap.Error(err))
				return err
			}

			var topics []domain.BookTopicResponse
			tx.Model(&domain.Topic{}).
				Select("topic, color").
				Joins("JOIN book_topics ON topics.id = book_topics.topic_id").
				Where("book_topics.book_id = ?", book.ID).
				Scan(&topics)

			/* Update book's updated_at */
			tx.Model(&domain.Book{}).Where("id = ?", book.ID).Update("updated_at", time.Now())

			response = &domain.BookPickResponse{
				Guid:    newPick.Guid,
				Content: newPick.Content,
				Index:   index,
			}
		}

		return nil
	})

	if err != nil {
		logger.Error("Failed to create book or pick (1)", zap.Error(err))
	}

	return response, err
}

func (service *serviceImpl) GetShortBooksList(userID uuid.UUID, params *domain.BookListParams) ([]domain.ShortBookResponse, error) {
	currentUser, err := service.userService.GetUserByGuid(userID)
	if err != nil {
		return nil, err
	}

	searchString := params.Search

	query := service.db.Model(&domain.Book{}).Where("user_id = ?", currentUser.ID)

	var books []domain.Book

	/* Use the GET for searching into the list */
	if searchString != "" {
		query = query.Where("title ILIKE ?", "%"+searchString+"%").Order("updated_at desc")
	} else {
		query = query.Order("updated_at desc")
	}

	query.Offset(params.Offset).Limit(params.Limit).Find(&books)

	return *domain.ShortBooksFromBookList(&books), nil
}

/*
	GetBooks returns a list of books for the user's homepage.
*/

func (service *serviceImpl) GetBooks(userID uuid.UUID, params *domain.BookListParams) ([]domain.BookResponse, error) {
	var response []domain.BookResponse

	currentUser, err := service.userService.GetUserByGuid(userID)
	if err != nil {
		return nil, err
	}

	commaSeparated := params.Topics

	err = service.db.Transaction(func(tx *gorm.DB) error {
		var books []domain.Book

		/* Get All currentUser's books filtered by topics */
		query := tx.Model(&domain.Book{}).
			Select("DISTINCT books.*").
			Joins("JOIN book_topics ON books.id = book_topics.book_id").Where("books.user_id = ?", currentUser.ID).
			Joins("JOIN topics ON topics.id = book_topics.topic_id")

		/* If at least one topic is provided, filter by them */
		if commaSeparated != "" && !strings.Contains(commaSeparated, "all") {
			topics := strings.Split(commaSeparated, ",")
			query.Where("topics.topic IN (?)", topics)
		}

		/* Then perform the query */
		query.Offset(params.Offset).
			Limit(params.Limit).
			Order("books.updated_at desc").
			Find(&books)

		bookIDs := make([]uint, len(books))
		for i, book := range books {
			bookIDs[i] = book.ID
		}

		/* Build domain.BookResponse for each book */

		for _, book := range books {

			picks := []domain.BookPickResponse{}
			tx.Model(&domain.BookPick{}).Where("book_id = ?", book.ID).Order("index DESC").Limit(3).Find(&picks)

			topics := []domain.Topic{}
			tx.Table("topics").
				Select("topics.topic, topics.color").
				Joins("JOIN book_topics ON topics.id = book_topics.topic_id").
				Where("book_topics.book_id = ?", book.ID).
				Scan(&topics)

			preview := domain.BookPickPreviewResponse{}
			tx.Model(&domain.BookPick{}).Where("book_id = ?", book.ID).Order("RANDOM()").Limit(1).First(&preview)

			var picksCount int64
			tx.Model(&domain.BookPick{}).Where("book_id = ?", book.ID).Count(&picksCount)

			response = append(response, domain.BookResponse{
				Guid:       book.Guid,
				Title:      book.Title,
				Author:     book.Author,
				UpdatedAt:  book.UpdatedAt,
				CreatedAt:  book.CreatedAt,
				PicksCount: picksCount,
				Preview:    &preview,
				Topics:     &topics,
				Picks:      &picks,
			})
		}

		return nil
	})

	return response, err
}

/*
	Get picks by bookID, offset and limit
*/

func (service *serviceImpl) GetPicksByBook(params *domain.GetPicksParams) ([]domain.BookPickResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	bookID, err := uuid.Parse(params.BookID)
	if err != nil {
		return nil, err
	}

	book, err := service.GetBookByGuid(bookID)
	if err != nil {
		return nil, err
	}

	var picks []domain.BookPickResponse

	orderBy := params.OrderBy
	orderByString := "index DESC"

	if orderBy == "asc" {
		orderByString = "index ASC"
	}

	err = service.db.Transaction(func(tx *gorm.DB) error {
		if params.UntilPickID != "" {
			var picksCount int64

			var pickToSearch domain.BookPick
			err := tx.Model(&domain.BookPick{}).Where("book_id = ? AND guid = ?", book.ID, params.UntilPickID).First(&pickToSearch).Error
			if err != nil {
				return err
			}

			tx.Model(&domain.BookPick{}).Where("book_id = ? AND index > ?", book.ID, pickToSearch.Index).Count(&picksCount)

			limit := params.Limit
			upperLimit := int(picksCount) + (limit-(int(picksCount)%limit))%limit

			err = tx.Model(&domain.BookPick{}).Where("book_id = ?", book.ID).Order(orderByString).Offset(params.Offset).Limit(upperLimit).Find(&picks).Error
			if err != nil {
				return err
			}

			return nil
		}

		err = tx.Model(&domain.BookPick{}).Where("book_id = ?", book.ID).Order(orderByString).Offset(params.Offset).Limit(params.Limit).Find(&picks).Error

		return err
	})

	return picks, err
}

/* Get books' topics */

func (service *serviceImpl) GetUserBooksTopics(userID uuid.UUID) ([]domain.BookTopicListResponse, error) {
	topics := []domain.BookTopicListResponse{}

	user, err := service.userService.GetUserByGuid(userID)
	if err != nil {
		return nil, err
	}

	var results []domain.BookTopicListResponse

	service.db.Table("topics").
		Select("topics.topic, topics.color, COUNT(DISTINCT books.id) AS count").
		Joins("JOIN book_topics ON book_topics.topic_id = topics.id").
		Joins("JOIN books ON books.id = book_topics.book_id").
		Where("books.user_id = ?", user.ID).
		Group("topics.topic, topics.color").
		Order("topics.topic").
		Scan(&results)

	var booksCount int64
	service.db.Model(&domain.Book{}).Where("user_id = ?", user.ID).Count(&booksCount)

	topics = append(topics, domain.BookTopicListResponse{
		Topic: "all",
		Count: uint(booksCount),
	})

	topics = append(topics, results...)

	return topics, nil
}

/* Edit book */

func (service *serviceImpl) EditBook(userID uuid.UUID, body *domain.EditBookBody) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	return service.db.Transaction(func(tx *gorm.DB) error {

		user, err := service.userService.GetUserByGuid(userID)
		if err != nil {
			return err
		}

		query := tx.Model(&domain.Book{}).Where("guid = ? AND user_id = ?", body.BookID, user.ID)

		if body.Title != "" {
			err := query.Update("title", body.Title).Error
			if err != nil {
				return err
			}
		}

		if body.Author != "" {
			err := query.Update("author", body.Author).Error
			if err != nil {
				return err
			}
		}

		if len(body.Topics) > 0 {
			book := domain.Book{}
			tx.Model(&domain.Book{}).Where("guid = ?", body.BookID).First(&book)

			/* 1. Delete all BookTopics relations in order to overwrite the new ones */
			tx.Table("book_topics").Where("book_id = ?", book.ID).Delete(&domain.BookTopic{})

			for _, topic := range body.Topics {
				/* 2. Add new Topics if they don't already exist */
				tx.Table("topics").Exec("INSERT INTO topics (user_id, topic, color) VALUES (?, ?, ?) ON CONFLICT (user_id, topic) DO UPDATE SET color = ?", user.ID, topic.Topic, topic.Color, topic.Color)

				/* 3. Create the association in the BookTopics bridge table */
				tx.Table("book_topics").
					Exec("INSERT INTO book_topics (book_id, topic_id) VALUES (?, (SELECT id FROM topics WHERE topic = ? AND user_id = ?))", book.ID, topic.Topic, user.ID)
			}

			/* 4. If no topics has at least one entry in the bridge table, delete the topic */
			tx.Table("topics").
				Exec("DELETE FROM topics WHERE user_id = ? AND id NOT IN (SELECT topic_id FROM book_topics)", user.ID)
		}

		if len(body.Picks) > 0 {
			changedPickIds := make([]uuid.UUID, len(body.Picks))
			for i, pick := range body.Picks {
				changedPickIds[i] = pick.Guid
			}

			/* Get all changed picks */
			picks := []domain.BookPick{}
			tx.Model(&domain.BookPick{}).Where("guid IN ?", changedPickIds).Find(&picks)

			// For each pick in picks update the index to the new one searching it using guid in body.Picks

			logger, _ := zap.NewProduction()
			defer logger.Sync()

			for index, pick := range picks {
				for _, newPick := range body.Picks {
					if pick.Guid == newPick.Guid {
						picks[index].Index = newPick.Index
					}
				}
			}

			for _, pick := range picks {
				if err := tx.Model(&domain.BookPick{}).Where("guid = ?", pick.Guid).Update("index", pick.Index).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

/* Edit Book Pick */
func (service *serviceImpl) EditBookPick(userID uuid.UUID, body *domain.EditBookPickBody) error {
	pickData := map[string]interface{}{}

	if body.Content != "" {
		pickData["content"] = body.Content
	}

	if body.Text != "" {
		pickData["content_text"] = body.Text
	}

	if body.Title != nil {
		pickData["title"] = *body.Title
	}

	pick := domain.BookPick{}

	return service.db.Transaction(func(tx *gorm.DB) error {

		bookId := body.BookId

		query := tx.Model(&domain.BookPick{}).Where("book_id = (SELECT id FROM books WHERE guid = ?) AND guid = ?", bookId, body.PickId)

		err := query.First(&pick).Error
		if err != nil {
			return err
		}

		/* If at least 20% of the content has changed, re-generate pick's search keywords */
		if utility.AtLeast20PercentChanged(pick.ContentText, body.Text) {
			message := domain.BookPickSearchKeywordMessage{
				PickID:      pick.ID,
				PickContent: body.Text,
				UserGuid:    userID,
			}

			err = sqs.SendMessage(sqs.QueueNames.PickKeywords, message)
			if err != nil {
				return err
			}
		}

		err = tx.Model(&domain.BookPick{}).Where("guid = ?", body.PickId).Updates(&pickData).Error
		if err != nil {
			return err
		}

		return nil
	})
}

/* Add pick's keywords in database */
func (service *serviceImpl) AddPickKeywords(userID uuid.UUID, pickID uint, keywords []string) error {
	return service.db.Transaction(func(tx *gorm.DB) error {

		user, err := service.userService.GetUserByGuid(userID)
		if err != nil {
			return err
		}

		/* 1. Insert all keyword into pick_search_keywords */
		pickKeywords := make([]domain.PickSearchKeyword, len(keywords))
		for i, keyword := range keywords {
			pickKeywords[i] = domain.PickSearchKeyword{
				PickID:  pickID,
				Keyword: keyword,
				UserID:  user.ID,
			}
		}

		err = tx.Table("pick_search_keywords").Create(&pickKeywords).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (service *serviceImpl) SemanticSearch(userID uuid.UUID, params *domain.SearchGetParams) ([]domain.SemanticSearchResponse, error) {
	user, err := service.userService.GetUserByGuid(userID)
	if err != nil {
		return nil, err
	}

	/*
		Perform a sintactic search across all picks, on books' title and them join with the picks_search_keywords table and search the query even in the keywords.
	*/

	response := []domain.SemanticSearchResponse{}

	query := params.Query

	err = service.db.Transaction(func(tx *gorm.DB) error {
		return tx.Raw(`
			WITH matched_picks AS (
				SELECT bp.id AS pick_id
				FROM book_picks bp
				WHERE bp.user_id = ? AND ((bp.content_text ILIKE '%' || ? || '%') OR (bp.title ILIKE '%' || ? || '%'))
			),

			matched_keywords AS (
				SELECT ps.pick_id
				FROM pick_search_keywords ps
				WHERE ps.keyword ILIKE '%' || ? || '%' AND ps.user_id = ?
			)

			SELECT DISTINCT b.guid AS book_id, b.title AS book_title, bp.guid AS pick_id, bp.content_text AS pick_content, bp.index AS pick_index, bp.title AS pick_title
			FROM book_picks bp
			JOIN books b ON b.id = bp.book_id
			WHERE bp.id IN (SELECT pick_id FROM matched_picks)
			OR bp.id IN (SELECT pick_id FROM matched_keywords)
			OR (b.title ILIKE '%' || ? || '%' AND b.user_id = ?)
			LIMIT ? OFFSET ?
	`, user.ID, query, query, query, user.ID, query, user.ID, params.Limit, params.Offset).Scan(&response).Error
	})

	return response, err
}

func (service *serviceImpl) SearchPickInBook(params *domain.SearchGetParams) ([]domain.SearchPickInBookResponse, error) {
	query := params.Query

	response := []domain.SearchPickInBookResponse{}

	bookGuid, err := uuid.Parse(params.BookID)
	if err != nil {
		return nil, err
	}

	book, err := service.GetBookByGuid(bookGuid)
	if err != nil {
		return nil, err
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	err = service.db.Raw(`
			SELECT bp.guid AS pick_id, bp.title AS pick_title, bp.content_text AS pick_content, bp.index AS pick_index FROM book_picks AS bp WHERE book_id = ? AND ((content_text ILIKE '%' || ? || '%') OR (title ILIKE '%' || ? || '%')) LIMIT ? OFFSET ?
		`, book.ID, query, query, params.Limit, params.Offset).Scan(&response).Error

	return response, err
}

func (service *serviceImpl) SaveBook(userID uuid.UUID, body *domain.SaveBookBody) (*domain.BookResponse, error) {
	user, err := service.userService.GetUserByGuid(userID)
	if err != nil {
		return nil, err
	}

	book, err := service.GetBookByGuid(body.BookID)
	if err != nil {
		return nil, err
	}

	response := domain.BookResponse{}

	err = service.db.Transaction(func(tx *gorm.DB) error {
		newBook := domain.Book{
			UserID: user.ID,
			Title:  book.Title,
			Author: book.Author,
		}

		/* 1. Create a copy of the book */
		if err := tx.Create(&newBook).Error; err != nil {
			return err
		}

		/* 2. Get all picks from the book */
		var picks []domain.BookPick
		pickQuery := tx.Model(&domain.BookPick{}).Where("book_id = ?", book.ID)
		pickQuery.Find(&picks)

		/* 3. Create a copy of each pick */
		picksCopy := make([]domain.BookPick, len(picks))

		for i, pick := range picks {
			picksCopy[i] = domain.BookPick{
				UserID:      user.ID,
				BookID:      newBook.ID,
				Content:     pick.Content,
				ContentText: pick.ContentText,
				Index:       pick.Index,
			}
		}

		/* 4. Bulk insert the picks */
		if err := tx.Create(&picksCopy).Error; err != nil {
			return err
		}

		/* 5.1 Get  all topics for the book */
		var bookTopics []domain.Topic
		tx.Table("topics").
			Select("topics.topic, topics.color").
			Joins("JOIN book_topics ON topics.id = book_topics.topic_id").
			Where("book_topics.book_id = ?", book.ID).
			Scan(&bookTopics)

		/* 5.2 Insert into topics the topics that don't exist yet */
		for _, topic := range bookTopics {
			err := tx.Exec("INSERT INTO topics (user_id, topic, color) VALUES (?, ?, ?) ON CONFLICT (user_id, topic) DO NOTHING", user.ID, topic.Topic, topic.Color).Error
			if err != nil {
				return err
			}
		}

		/* 5.3 Insert into book_topics the book_id, the topic_id and the user_id */
		for _, topic := range bookTopics {
			err := tx.Exec("INSERT INTO book_topics (book_id, topic_id) VALUES (?, (SELECT id FROM topics WHERE topic = ? AND user_id = ?))", newBook.ID, topic.Topic, user.ID).Error
			if err != nil {
				return err
			}
		}

		/* 5.4 Get all topics for the book */
		newTopics := []domain.Topic{}
		tx.Table("topics").
			Select("topics.topic, topics.color").
			Joins("JOIN book_topics ON topics.id = book_topics.topic_id").
			Where("book_topics.book_id = ?", newBook.ID).
			Scan(&newTopics)

		/* 6. Get all picks for the book */
		newPicks := []domain.BookPickResponse{}
		tx.Model(&domain.BookPick{}).Where("book_id = ?", newBook.ID).Order("index DESC").Limit(3).Find(&newPicks)

		/* 7. Get a random pick for the preview */
		preview := domain.BookPickPreviewResponse{}
		tx.Model(&domain.BookPick{}).Where("book_id = ?", newBook.ID).Order("RANDOM()").Limit(1).First(&preview)

		response = domain.BookResponse{
			Guid:       newBook.Guid,
			Title:      newBook.Title,
			Author:     newBook.Author,
			UpdatedAt:  newBook.UpdatedAt,
			CreatedAt:  newBook.CreatedAt,
			PicksCount: int64(len(newPicks)),
			Preview:    &preview,
			Topics:     &newTopics,
			Picks:      &newPicks,
		}

		/* 8. Copy for each pick the keywords */
		pickIDs := make([]uint, len(picksCopy))
		for i, pick := range picksCopy {
			pickIDs[i] = pick.ID
		}

		/* 8.1 Get all keywords for the picks */
		keywords := []domain.PickSearchKeyword{}
		tx.Table("pick_search_keywords").Where("pick_id IN ?", pickIDs).Find(&keywords)

		if len(keywords) != 0 {
			keywordsCopy := make([]domain.PickSearchKeyword, len(keywords))
			for i, keyword := range keywords {
				keywordsCopy[i] = domain.PickSearchKeyword{
					PickID:  pickIDs[i],
					Keyword: keyword.Keyword,
					UserID:  user.ID,
				}
			}

			if err := tx.Table("pick_search_keywords").Create(&keywordsCopy).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return &response, err
}
