// handlers.client.go
package fokusov

/*
func getClient(c *gin.Context) {
	// Check if the client hostname is valid
	//if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
	clientHostname := c.Param("article_id")

	// Check if the client exists
	if client, err := getArticleByID(articleID); err == nil {
		// Call the render function with the title, article and the name of the
		// template
		render(c, gin.H{
			"title":   article.Title,
			"payload": article}, "article.html")

	} else {
		// If the article is not found, abort with an error
		c.AbortWithError(http.StatusNotFound, err)
	}

}

func showIndexPage(c *gin.Context) {
	//articles := getAllArticles()

	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Home Page"}, "index.html")
	//"payload": articles}, "index.html")
}
*/
/*
func showArticleCreationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Create New Article"}, "create-article.html")
}

func createArticle(c *gin.Context) {
	// Obtain the POSTed title and content values
	title := c.PostForm("title")
	content := c.PostForm("content")

	if a, err := createNewArticle(title, content); err == nil {
		// If the article is created successfully, show success message
		render(c, gin.H{
			"title":   "Submission Successful",
			"payload": a}, "submission-successful.html")
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
*/
