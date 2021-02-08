package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const landing = `<!DOCTYPE html>
<html lang="en">

<head>
<title>the best tracks that an artist has to offer</title>
</head>

<body>
<h1>Listen to the best tracks that an artist has to offer</h1>


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

