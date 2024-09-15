package utility

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type VolumeInfo struct {
	Authors []string `json:"authors"`
}

type Item struct {
	VolumeInfo VolumeInfo `json:"volumeInfo"`
}

type Response struct {
	Items []Item `json:"items"`
}

func GetAuthorFromTitle(bookTitle string) (string, error) {
	url := "https://www.googleapis.com/books/v1/volumes"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("q", bookTitle)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get author from title")
	}

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if len(data.Items) > 0 {
		book := data.Items[0]

		if len(book.VolumeInfo.Authors) > 0 {
			authors := strings.Join(book.VolumeInfo.Authors, ", ")
			return authors, nil
		}
		return "", nil
	}

	return "", nil
}
