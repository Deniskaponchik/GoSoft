// models.article.go

package main

import "errors"

type article struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AnomStruct []anomaly
}
type anomaly struct {
	ApName     string
	Date       string
	AnomalyStr []string
}

var anomStr1 = []string{"USER_LOW_PHY_RATE", "USER_HIGH_TCP_LATENCY"}
var anomStr2 = []string{"USER_DNS_TIMEOUT", "USER_HIGH_WIFI_RETRIES"}

var anomSlice = []anomaly{
	anomaly{ApName: "KRA-2FL-RECEPTION", Date: "2023-11-09", AnomalyStr: anomStr1},
	anomaly{ApName: "KRA-2FL-RECEPTION", Date: "2023-11-10", AnomalyStr: anomStr2},
}

// For this demo, we're storing the article list in memory
// In a real application, this list will most likely be fetched
// from a database or from static files
var articleList = []article{
	article{ID: 1, Title: "Article 1", Content: "Article 1 body", AnomStruct: anomSlice},
	article{ID: 2, Title: "Article 2", Content: "Article 2 body"},
}

// Return a list of all the articles
func getAllArticles() []article {
	return articleList
}

// Fetch an article based on the ID supplied
func getArticleByID(id int) (*article, error) {
	for _, a := range articleList {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("Article not found")
}

// Create a new article with the title and content provided
func createNewArticle(title, content string) (*article, error) {
	// Set the ID of a new article to one more than the number of articles
	a := article{ID: len(articleList) + 1, Title: title, Content: content}

	// Add the article to the list of articles
	articleList = append(articleList, a)

	return &a, nil
}
