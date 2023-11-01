package main

import (
	"context"
	"fmt"
	pb "grpc/hello-server/proto"
	"io/ioutil"
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

	r.GET("/method2_request", func (c *gin.Context) {
      var responses []string;
	  var wg sync.WaitGroup
	  responseCh := make(chan string)

	  for i := 0; i < 10; i++ {
		wg.Add(1) //increment the waitgroup counter
		go func ()  {
			defer wg.Done()
			resp, err := http.Post("http://3.105.180.131:3000/method2_response", "application/json", nil);
			if err != nil {
			  c.JSON(500, err)
			  return
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
			  c.JSON(500, err)
			  return
			}
			responseCh <- string(body)
			defer resp.Body.Close();


		}()
	  }

	  go func() {
		wg.Wait()
		close(responseCh)
	  }()


	  for i := 0; i < 10; i++ {
		response := <- responseCh
		// if  response != "" {
		// 	c.JSON(500, response)
		// 	return
		// }
		responses = append(responses, response)
	  }

	  c.JSON(200, gin.H {
		"FIRST": "Send request successfully",
		"SECOND": responses,
	  })
	})

	r.Use(CORSMiddleware())

	r.Run(":5000")
}
