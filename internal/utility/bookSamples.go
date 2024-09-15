package utility

import (
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CopyBookSamplesToUser(userID uint, tx *gorm.DB) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	/* 1. Get all book samples */
	books := []domain.Book{}

	err := tx.Model(&domain.Book{}).Where("user_id IS NULL").Find(&books).Error
	if err != nil {
		logger.Error("failed to get book samples", zap.Error(err))
		return err
	}

	for _, book := range books {
		newBook := domain.Book{
			UserID: userID,
			Title:  book.Title,
			Author: book.Author,
		}

		/* 1.1 Copy books */
		if err := tx.Create(&newBook).Error; err != nil {
			logger.Error("failed to copy book", zap.Error(err))
			break
		}

		picks := []domain.BookPick{}
		tx.Model(&domain.BookPick{}).Where("book_id = ?", book.ID).Order("index DESC").Find(&picks)

		picksCopy := make([]domain.BookPick, len(picks))

		for i, pick := range picks {
			picksCopy[i] = domain.BookPick{
				UserID:      userID,
				BookID:      newBook.ID,
				Content:     pick.Content,
				ContentText: pick.ContentText,
				Title:       pick.Title,
				Index:       pick.Index,
			}
		}

		/* 1.2 Copy picks */
		if err := tx.Create(&picksCopy).Error; err != nil {
			logger.Error("failed to copy picks", zap.Error(err))
			break
		}

		var bookTopics []domain.Topic
		tx.Table("topics").
			Select("topics.topic, topics.color").
			Joins("JOIN book_topics ON topics.id = book_topics.topic_id").
			Where("book_topics.book_id = ?", book.ID).
			Scan(&bookTopics)

		/* 1.3 Insert into topics the topics that don't exist yet */
		for _, topic := range bookTopics {
			err := tx.Exec("INSERT INTO topics (user_id, topic, color) VALUES (?, ?, ?) ON CONFLICT (user_id, topic) DO NOTHING", userID, topic.Topic, topic.Color).Error
			if err != nil {
				logger.Error("failed to insert into topics", zap.Error(err))
				break
			}
		}

		/* 1.4 Insert into book_topics the book_id, the topic_id and the user_id */
		for _, topic := range bookTopics {
			err := tx.Exec("INSERT INTO book_topics (book_id, topic_id) VALUES (?, (SELECT id FROM topics WHERE topic = ? AND user_id = ?))", newBook.ID, topic.Topic, userID).Error
			if err != nil {
				logger.Error("failed to insert into book_topics", zap.Error(err))
				break
			}
		}

		/* 1.5. Copy for each pick the keywords */
		pickIDs := make([]uint, len(picksCopy))
		for i, pick := range picksCopy {
			pickIDs[i] = pick.ID
		}

		/* 1.6 Get all keywords for the picks */
		keywords := []domain.PickSearchKeyword{}
		tx.Table("pick_search_keywords").Where("pick_id IN ?", pickIDs).Find(&keywords)

		if len(keywords) != 0 {
			keywordsCopy := make([]domain.PickSearchKeyword, len(keywords))
			for i, keyword := range keywords {
				keywordsCopy[i] = domain.PickSearchKeyword{
					PickID:  pickIDs[i],
					Keyword: keyword.Keyword,
					UserID:  userID,
				}
			}

			if err := tx.Table("pick_search_keywords").Create(&keywordsCopy).Error; err != nil {
				logger.Error("failed to copy keywords", zap.Error(err))
				break
			}
		}
	}

	return err
}
