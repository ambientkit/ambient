// Package ambient is a modular web application framework.
package ambient

// // IAppLogger represents an application logger.
// type IAppLogger interface {
// 	Debug(format string, v ...interface{})
// 	Info(format string, v ...interface{})
// 	Warn(format string, v ...interface{})
// 	Error(format string, v ...interface{})
// }

// // IAppRenderer represents an application template enginer.
// type IAppRenderer interface {
// 	Page(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
// 	PageContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
// 	Post(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
// 	PostContent(w http.ResponseWriter, r *http.Request, content string, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
// 	Error(w http.ResponseWriter, r *http.Request, content string, statusCode int, fm template.FuncMap, vars map[string]interface{}) (status int, err error)
// }

// // IAppRouter represents an application request router.
// type IAppRouter interface {
// 	Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
// 	Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
// 	Patch(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
// 	Put(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
// 	Head(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
// 	Options(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
// 	Delete(path string, fn func(http.ResponseWriter, *http.Request) (int, error))
// 	Error(status int, w http.ResponseWriter, r *http.Request)
// 	Param(r *http.Request, param string) string
// 	ServeHTTP(w http.ResponseWriter, r *http.Request)
// }
