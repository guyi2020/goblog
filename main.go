package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "html/template"
	_ "io/ioutil"
	"log"
	"net/http"
	_ "strings"
	"time"
)

type Info struct {
	Phone string `form:"phone" uri:"phone" binding:"required"`
	Hobby []string `form:"hobby[]"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

type User struct {
	Username string `form:"username" uri:"username" binding:"required"`
	Info Info
}

func GetUser(c *gin.Context)  {
	var user User
	c.Bind(&user)
	c.JSON(200, gin.H{
		"phone": user.Info,
		"username": user.Username,
	})
}

func indexHandler(c *gin.Context) {
	c.HTML(200, "form.html", nil)
}

func formHandler(c *gin.Context)  {
	var info Info
	c.ShouldBind(&info)
	c.JSON(200, gin.H{
		"hobby": info.Hobby,
	})
}

func registerUser(c *gin.Context)  {
	var user User
	if c.ShouldBind(&user) != nil {
		log.Println("error...")
		log.Println(user.Username)
		log.Println(user.Info.Phone)
		log.Println(user.Info.Birthday)
	}
	c.String(200, "Success")
}

func bindUri(c *gin.Context)  {
	var user User
	if err := c.ShouldBindUri(&user); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"username": user.Username, "phone": user.Info.Phone})
}


/*
func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
*/

// 使用中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// 请求之前执行
		c.Set("example", "使用中间件")
		c.Next()
		// 请求之后执行
		latency := time.Since(t)
		log.Print(latency)

		status := c.Writer.Status()
		log.Println(status)
	}
}

func main() {

	r := gin.Default()
	r.Use(Logger())

	// 9. 配置打印日志格式
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 1. 转Ascii码
	r.GET("/someJson", func(c *gin.Context) {
		data := map[string]interface{} {
			"name": "李好",
			"age": 18,
		}
		c.AsciiJSON(http.StatusOK, data)
	})

	// 2. 嵌套赋值参数
	r.GET("/user", GetUser)

	// 3.表单checkbox提交
	r.LoadHTMLGlob("views/*")
	r.GET("/form", indexHandler)
	r.POST("/form", formHandler)

	// 4.绑定query 参数 或者 post 数据
	r.GET("/register", registerUser)

	//5. 绑定uri, 注意结构体设置 uri，binding
	r.GET("/bind/:username/:phone", bindUri)

	/*
		6. 构建简单的二进制tmpl文件
			a. 安装包 go get github.com/jessevdk/go-assets-builder
			b. go-assets-builder html -o assets.go
			c. go build -o assets-in-binary
			d. 运行 ./assets-in-binary
	 */
	/*
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "/html/index.tmpl", gin.H{
			"Foo": "World",
		})
	})
	r.GET("/bar", func(c *gin.Context) {
		c.HTML(http.StatusOK, "/html/bar.tmpl", gin.H{
			"Bar": "World",
		})
	})
	*/

	//7. 开启 控制台打印彩色的日志， 关闭使用 gin.DisableConsoleColor()
	gin.ForceConsoleColor()



	 // 8. 配置开启服务参数 默认可使用 r.Run()
	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 10. 使用中间件
	r.GET("/test", func(c *gin.Context) {
		example := c.MustGet("example").(string)
		log.Println(example)
	})

	s.ListenAndServe()

}
