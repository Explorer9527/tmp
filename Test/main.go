package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
)

const str string = "123" //必须赋值

const ( //枚举类型
	A string = "a"
	B string = "b"
)

type Flg int

const ( //枚举类型
	_ Flg = iota //iota是从0开始自增的,跳过
	O            //1
	N            //2
)

var a int = 1 //全局变量允许声明后不使用
var (
	b = true
	c = "123"
)

// 请求实体
type People struct {
	name string `from:"name"`
	info string `from:"info"`
}

// 结构体返回
type Response struct {
	a string
	b bool
}

// 函数绑定类型
type Bin struct {
	s int
}

/*
启动服务：go run main.go
*/
func main() {
	dataTypes()                 //1)数据类型测试
	fmt.Println("123")          //
	bin := &Bin{s: 1}           //2)闭包测试,这里注意也可以使用bin := Bin{s: 1}形式但是涉及隐式转换
	closure := bin.functions(3) //  调用绑定函数,返回闭包
	s1, s2 := closure(5)        //  调用闭包函数
	fmt.Println(s1, s2)         //  5  aaaaaa
	loop()                      //3)循环测试
	var no Flg = 1              //  枚举类型测试
	res := no.judge()           //  调用判断方法
	fmt.Println(res)            //  old
	goroutines()                //4)goroutine、defer、context、chanel详解
	engine := gin.Default()     //启动http服务
	//不同请求返回
	request1(engine)
	request2(engine)
	responseHTML(engine)

	if result := engine.Run(":8090"); result != nil {
		log.Fatal(result.Error())
	}
}

// 1.数据类型
func dataTypes() int {
	//变量定义
	var a int = 1  //完整定义
	var b int      //仅声明
	var c = 10     //无类型定义
	d := b + a + c //自动推断类型
	fmt.Println(d)
	//批量定义
	var e, f = 1, 2
	g, h := e+f, "dfsf"
	fmt.Println(h, g)
	//指针
	var p1 *int                //定义指针
	fmt.Println(p1)            //变量值为nil
	i1 := 1                    //定义一个数
	p1 = &i1                   //指针值赋值为i1的地址
	fmt.Println(p1, &i1)       //显示为i1内存地址：0xc00000b478
	fmt.Println(*p1)           //显示内存地址对应的值：1
	*p1 = 3                    //修改内存地址处的值
	fmt.Println(*p1, i1)       //打印的值为：3 3
	var p2 **int               //定义指针的指针
	fmt.Println(p2)            //变量值为nil
	p2 = &p1                   //p2值赋值为p1这个变量的地址
	fmt.Println(p2)            //p1内存地址：0xc00006e780
	fmt.Println(**p2)          //内存地址最终指向的值：3
	**p2 = 10                  //修改最终地址处的值为10
	fmt.Println(*p2, *p1, i1)  //0xc00000b478  10  10
	fmt.Println(**p2, *p1, i1) //10 10 10
	n := unsafe.Pointer(&i1)   //先将i1指针转换为pointer类型
	m := uintptr(n)            //再将pointer类型转换为整型数值
	fmt.Println(n, m)          //打印0xc00000b478 824633767032
	x := unsafe.Pointer(m + 1) //偏移地址值，并转换为pointer类型
	y := (*uint8)(x)           //整数值转换为地址，注意这里y值完全不确定边界、类型等，属于危险操作
	z := *y                    //获取新地址的值，注意这里z值完全不确定边界、类型等，属于危险操作
	fmt.Println(x, y, z)       //打印 0xc00000b479  0xc00000b479  0
	return 0
}

// 2.函数、闭包
func (bin *Bin) functions(int) func(int) (int, string) {
	return func(a int) (c int, d string) { //返回闭包
		c = a
		d = "aaaaaa"
		return
	}
}

// 3.循环
func loop() {
	var i int = 1
lable1:
	if i == 5 {
		fmt.Println("456")
		testAi()
	} else {
		for ; i < 10; i++ {
			fmt.Println(i)
			if i == 5 {
				goto lable1
			}
		}
	}
}

func testAi() {
	var a int8 = 1
	fmt.Printf("a: %d\n", a)
}

// 4.枚举判断新旧
func (f *Flg) judge() string {
	if *f == 1 {
		return "old"
	} else if *f == 2 {
		return "new"
	} else {
		return "err"
	}
}

// 5.多任务并发
func goroutines() {
	//1.仅仅同步管理，使用sync.Waitgroup
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)   //添加一个计数
		go func() { //go关键字开启一个协程
			defer wg.Done() //函数结束时调用Done
			fmt.Println(i)  //打印的值为：0 1 2 3 4 5 6 7 8 9
		}()
	}
	wg.Wait() //等待所有协程执行完毕
	//2.使用channel并发控制
	ch := make(chan bool, 10) //创建一个容量为10的channel
	for i := 0; i < 10; i++ {
		go func() {
			ch <- true //向channel中写入数据
			fmt.Println(i)
		}()
	}
	for i := 0; i < 10; i++ {
		<-ch //从channel中读取数据
	}
	close(ch) //关闭channel

	//3.sync.waitgroup和channel结合使用:汇总多个协程的数据
	var swg sync.WaitGroup
	cha := make(chan int, 5)
	swg.Add(1)
	go func() { //协程1放入数据
		defer swg.Done()
		for i := 0; i < 2; i++ {
			cha <- i
		}
	}()
	swg.Add(1)
	go func() { //协程2放入数据
		defer swg.Done()
		for i := 2; i < 5; i++ {
			cha <- i
		}
	}()
	swg.Wait()
	close(cha)
	for i := range cha {
		fmt.Println(i) //打印的值为：3 4 0 1
	}
	// 4.select和channel结合使用
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch3 := make(chan int, 1)
	for i := 0; i < 5; i++ {
		ch1 <- i
		ch2 <- i + 1
		ch3 <- i + 2
	}
	for i := 0; i < 5; i++ {
		select {
		case a := <-ch1: //接收1的数据，如果3路数据都没阻塞，则随机选择一个执行
			fmt.Println(a)
		case b := <-ch2: //接收2的数据，如果3路数据都没阻塞，则随机选择一个执行
			fmt.Println(b)
		case c := <-ch3: //接收3的数据，如果3路数据都没阻塞，则随机选择一个执行
			fmt.Println(c)
		case <-time.After(time.Second * 3): //超时处理机制，这里不会走到，但是放在这里
			fmt.Println("3s")
		default: //添加default处理，防止阻塞，这里注意default的优先级最高，因此超时机制不会生效
			fmt.Println("default")
		}
	}
}

// 6.指定路由分类处理
func request1(engine *gin.Engine) {
	//处理get请求: http://localhost:8090/hello?name=123&password=123
	engine.GET("/hello", func(c *gin.Context) {
		fmt.Printf("path:%s", c.FullPath())
		name := c.Query("name")
		password := c.Query("password")
		c.Writer.Write([]byte("nihao" + name + password)) //切片返回
	})
	//处理post请求
	engine.POST("/hello", func(c *gin.Context) {
		fmt.Printf("path:%s", c.FullPath())
		c.Writer.WriteString("你好") //string返回
	})
}

// 7.路由通用处理
func request2(engine *gin.Engine) {
	//处理get请求: http://localhost:8090/world?name=123&info=sjdkadsa
	engine.Handle("GET", "/world", func(c *gin.Context) {
		fmt.Printf("path:%s", c.FullPath())
		var people People
		err := c.ShouldBindQuery(&people) //结构体直接映射
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		fmt.Println(people.name)
		fmt.Println(people.info)

		var response Response = Response{
			a: "123",
			b: true,
		}
		c.JSON(http.StatusOK, &response) //结构体返回
	})
	//处理post请求
	engine.Handle("POST", "/world", func(c *gin.Context) {
		fmt.Printf("path:%s", c.FullPath())
		var people People
		// if err := c.ShouldBind(&people); err != nil {
		if err := c.BindJSON(&people); err != nil { //json数据格式映射请求数据
			log.Fatal(err.Error())
			return
		}
		fmt.Println(people.info)
		c.JSON(http.StatusOK, map[string]interface{}{ //map形式返回
			"name": people.name,
			"info": people.info,
		})
	})
}

// 8.返回静态资源
func responseHTML(engine *gin.Engine) {
	//请求：http://localhost:8090/index
	engine.LoadHTMLGlob("./html/*")   //加载全局HTML
	engine.Static("/static", "./img") //加载静态资源：请求路径和真实路径映射关系
	engine.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{ //模板添加数据
			"title": "测试",
			"path":  "模板测试",
		})
	})
}
