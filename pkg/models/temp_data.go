package models

// TemplateData holds data sent from handlers to template
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string //* Stands For Cross Site Request Forgery Token
	// *(Long String if random numbers) and we will generate it using Golang nosurf
	Flash   string //* FlashMessage to send to the user (({status:success}))
	Warning string
	Error   string
}
