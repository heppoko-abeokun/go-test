package main

import (
	"net/http"
	"io"
	"os"
	"time"
	"html/template"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	renderer := &TemplateRenderer{
    	templates: template.Must(template.ParseGlob("views/*.html")),
    }
    e.Renderer = renderer

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", "")
	})

	// Route => handler
	e.GET("/page1", func(c echo.Context) error {
		return c.Render(http.StatusOK, "page1", map[string]interface{}{
		  "name": "Abe!",
		  "today" : time.Now().Format("2006/01/02"), 
        })
	}).Name = "page1"

	// Route => handler
	e.GET("/form-page", func(c echo.Context) error {
		return c.Render(http.StatusOK, "form-page", map[string]interface{}{
			"user_name": c.FormValue("user_name"),
		})
	})

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}