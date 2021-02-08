package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const landing = `<!DOCTYPE html>
<html lang="en">

<head>
<title>The Very Best</title>
</head>

<body>
<h1>The Very Best Of</h1>
<h2>Listen to the best that an artist has to offer</h2>


<form action="/artist" method="GET">
Artist:
<input type="text" name="Artist">
<br/>
<input type="submit" value="Submit">
</form>

</body>`

func Landing(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(landing))
}

