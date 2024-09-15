package langchain

import (
	"errors"
	"strings"

	"go.uber.org/zap"
)

/* Generates up to 2 topics starting from the given text. */

func GenerateBookTopics(pickContent string) ([]string, error) {
	promptString := `
		Generate up to 2 topics starting from this text: "{{.text}}".

		Requirements:
			- Each topic must be a discipline or field of study related to the text.
			- Each topic should be a single word or a compound word representing a discipline, such as 'design', 'psychology', 'quantum mechanics'.
			- Prioritize broader topics over more specific ones. For instance, prefer 'physics' over 'theoretical physics'.
			- Do not separate words with "_" or "-". If there is a space between words, use a space.

		Return the output as an object of type {"values": ["topic1", "topic2"]}.
	`

	response, err := NewLLMRequest(promptString, map[string]interface{}{"text": pickContent})
	if err != nil {
		return nil, err
	}

	topics, ok := response["values"].([]interface{})
	if !ok {
		return nil, errors.New("unable to parse response")
	}

	var topicsStr []string
	for _, topic := range topics {
		normalised := strings.ToLower(topic.(string))
		topicsStr = append(topicsStr, normalised)
	}

	return topicsStr, nil
}

/* Generate 30 keywords to perform semantic search for each pick */

func GeneratePickKeywords(pickContent string) ([]string, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	promptString := `
		Generate 5 keywords based on the following text: "{{.text}}". Ensure the keywords are relevant to the subject matter of the text.

		Requirements:
			- Exclude words that are directly taken from the text.
			- Each keyword should be a single word or a compound word.
			- Do not use underscores ("_") or hyphens ("-") to connect words; if a keyword consists of multiple words, use a space.
			- Exclude dates, numbers, and overly general words like 'innovation', 'technology', and 'science'.
			- Limit the inclusion of 'isms' like 'capitalism' and 'socialism' to the most pertinent ones.
			- Incorporate specific names of people's creations, such as 'Wassily Chair', if relevant.

		Return the output as an object of type {"keywords": ["keyword1", "keyword2", ...]}.
	`

	response, err := NewLLMRequest(promptString, map[string]interface{}{"text": pickContent})
	if err != nil {
		logger.Error("Error generating pick keywords", zap.Error(err))
		return nil, err
	}

	logger.Info("Generated pick keywords", zap.Any("response", response))

	keywords, ok := response["keywords"].([]interface{})
	if !ok {
		logger.Error("Error parsing response", zap.Any("response", response))
		return nil, errors.New("unable to parse response")
	}

	keywordsStr := make([]string, len(keywords))

	for i, keyword := range keywords {
		normalised := strings.ToLower(keyword.(string))
		keywordsStr[i] = normalised
	}

	return keywordsStr, nil
}

/* Enrich the pick content by correcting the text and adding more details. */
func EnrichPickContent(pickContent string) (string, error) {
	promptString := `
		Generate a sharp pick based on the following text: {{.text}}.
		A sharp pick is an enhanced version of the original text with added details to improve its depth and clarity.

		Requirements:
			- The text provided by the user may have no meaning or contain errors.
			- Don't be hallucinated by the text in the case of nonsense or errors.
			- Carefully analyze the context and provide an informative and well-structured response.
			- If the context is unclear or incomplete, return an empty string.
			- Do not try to make sense of the text by combining unrelated elements.
			- If the text is incorrect, just provide a corrected version without telling the user where the error is.
			- Write the response in the same language as the input text and keep the same level of simplicity/complexity.
			- Give user more insights about the topic in the text, be informative not just saying like "the theory provides a new perspective on the topic".
			- The enhanced text must not exceed 300 characters.

		Return the output as an object of type {"content": "enhanced_text"}.
	`

	response, err := NewLLMRequest(promptString, map[string]interface{}{"text": pickContent})
	if err != nil {
		return "", err
	}

	enrichedText, ok := response["content"].(string)
	if !ok {
		return "", errors.New("unable to parse response")
	}

	return enrichedText, nil
}

/* Generate a detailed explanation starting from a given keyword. */
func GenerateKeywordExplanation(keyword string) (map[string]interface{}, error) {
	promptString := `
		Provide a detailed explanation of the term "{{.text}}".

		Requirements:
			- The explanation should be informative and concise.
			- The explanation should be written in complete sentences and be grammatically correct.
			- The explanation length should not exceed 300 characters.
			- If the term is ambiguous or has multiple meanings, provide the most common or relevant definition.
			- If the term is not recognized or cannot be explained accurately, return an empty string.
			- Generate an array of up to 3 urls to relevant sources that support the enhanced text.
			- If the keyword is a person's name, provide a brief biography or description of their work.

		Return the output as an object of type {"content": "explanation" "sources": ["url1", "url2", "url3"]}.
	`

	response, err := NewLLMRequest(promptString, map[string]interface{}{"text": keyword})
	if err != nil {
		return nil, err
	}

	return response, nil
}

/* Translate a word or phrase into a different language. */
func TranslateWord(word, language string) (map[string]interface{}, error) {
	promptString := `
		Translate the word "{{.text}}" into the language: "{{.language}}".
		Then, write a brief explanation of the word in the target language.

		Requirements:
			- The translation should be accurate and reflect the meaning of the original word.
			- If the word has multiple meanings, provide the most common or relevant translation.
			- If the word is not recognized or cannot be translated accurately, return an empty string.
			- The explanation length should not exceed 150 characters.
			- Provide a url to a dictionary, in the language {{.language}}, for the translated word.

		Return the response as an object of type {"word": "translated_word", "explanation": "brief_explanation", "url": <dictionary_url"}.
	`

	response, err := NewLLMRequest(promptString, map[string]interface{}{
		"text":     word,
		"language": language,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
