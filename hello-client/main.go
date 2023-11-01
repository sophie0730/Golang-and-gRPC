package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "grpc/hello-server/proto"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Allow-Credentials", "true")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
    r := gin.Default()
	r.GET("/method1", func (c *gin.Context)  {
		conn, err := grpc.Dial("http://3.105.180.131:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
	
		//建立連接
		client := pb.NewSayHelloClient(conn)

		var responses []string
		responseCh := make(chan string)

		for i := 0; i < 10; i++ {
			go func ()  {
				resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "sophie"})
				responseMsg := resp.GetResponseMsg()
				fmt.Println(responseMsg)

				responseCh <- responseMsg
			}()
		}
		for i := 0; i < 10; i++ {
			response := <- responseCh
			// if  response != "" {
			// 	c.JSON(500, response)
			// 	return
			// }
			responses = append(responses, response)
		  }
	
		c.JSON(200, gin.H{
			"response": responses,
		})
	})

	// r.GET("/method2_request", func (c *gin.Context) {
    //   var responses []string;
	//   var wg sync.WaitGroup
	//   responseCh := make(chan string)

	//   for i := 0; i < 10; i++ {
	// 	wg.Add(1) //increment the waitgroup counter
	// 	go func ()  {
	// 		defer wg.Done()
	// 		resp, err := http.Post("http://3.105.180.131:3000/method2_response", "application/json", nil);
	// 		if err != nil {
	// 		  c.JSON(500, err)
	// 		  return
	// 		}

	// 		body, err := ioutil.ReadAll(resp.Body)
	// 		if err != nil {
	// 		  c.JSON(500, err)
	// 		  return
	// 		}
	// 		responseCh <- string(body)
	// 		defer resp.Body.Close();


	// 	}()
	//   }

	//   go func() {
	// 	wg.Wait()
	// 	close(responseCh)
	//   }()


	//   for i := 0; i < 10; i++ {
	// 	response := <- responseCh
	// 	// if  response != "" {
	// 	// 	c.JSON(500, response)
	// 	// 	return
	// 	// }
	// 	responses = append(responses, response)
	//   }

	//   c.JSON(200, gin.H {
	// 	"FIRST": "Send request successfully",
	// 	"SECOND": responses,
	//   })
	// })
	r.GET("/method2_request", func(c *gin.Context) {
        var wg sync.WaitGroup
        var mu sync.Mutex

        // 创建一个通道来收集 JSON 响应
        jsonResponseCh := make(chan map[string]interface{}, 10)

        for i := 0; i < 10; i++ {
            wg.Add(1)

            go func() {
                defer wg.Done()

                res, err := http.Post("http://3.105.180.131:3000/method2_response", "application/json", nil)
                if err != nil {
                    c.JSON(500, gin.H{
                        "error": err.Error(),
                    })
                    return
                }

                defer res.Body.Close()

                var jsonResponse map[string]interface{}
                decoder := json.NewDecoder(res.Body)

                if err := decoder.Decode(&jsonResponse); err != nil {
                    c.JSON(500, gin.H{
                        "error": err.Error(),
                    })
                    return
                }

                // 使用互斥锁以确保安全地向通道发送响应
                mu.Lock()
                jsonResponseCh <- jsonResponse

                // 如果已经收集了足够的响应，关闭通道
                if len(jsonResponseCh) == 10 {
                    close(jsonResponseCh)
                }
                mu.Unlock()
            }()
        }

        // 等待所有请求完成
        wg.Wait()

        // 从通道中获取所有 JSON 响应
        var allJSONResponses []map[string]interface{}
        for jsonResponse := range jsonResponseCh {
            allJSONResponses = append(allJSONResponses, jsonResponse)
        }

        // 返回所有 JSON 响应
        c.JSON(200, allJSONResponses)
    })



	r.Use(CORSMiddleware())
	r.Run(":5000")
}
