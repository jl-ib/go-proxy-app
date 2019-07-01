
package middleware

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/kataras/iris"
	"strconv"
	"sort"
)

// Queue
type Queue struct {
	Domain string
	Weight int
	Priority int
}

var Que []string

// // Queues declaration
// var Qs []*Queue

type Repository interface {
	Read() []*Queue
}

func (q *Queue) Read() []*Queue {
	path, _ := filepath.Abs("")
	file, err := os.Open(path + "/api/middleware/domain.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	queue := []*Queue{};
	objcnt := 0;
	cnt := 0;
	obj := &Queue{};

	for scanner.Scan() {
		if scanner.Text() == "" {
			queue = append(queue, obj)
			fmt.Println("OUT", scanner.Text())
			objcnt++
			obj = &Queue{}
			continue
		}
		cstr := strings.Split(scanner.Text(), ":")
		if (len(cstr) > 1) {
			if (cnt % 3 == 1) {
				objw, _ := strconv.Atoi(cstr[1]);
				obj.Weight = objw;
			} else {
				objp, _ := strconv.Atoi(cstr[1]);
				obj.Priority = objp;
			}
		} else {
			if (cnt % 3 == 0) {
				objD := cstr[0];
				obj.Domain = objD;
			}
		}

		cnt++;
		fmt.Println("IN", scanner.Text())
	}
	
	prioritizeQueue(queue)
	
	fmt.Println("print", queue[0].Domain, queue[1].Priority, queue[2].Weight, queue[3])
	return queue;
}


func prioritizeQueue(queue []*Queue) {
	sort.Slice(queue, func(i, j int) bool {
		return queue[i].Weight > queue[j].Weight
	})
	sort.Slice(queue, func(i, j int) bool {
		return queue[i].Priority > queue[j].Priority
	})
}

// func MockQueue() []*Queue {
// 	return []*Queue {
// 		{
// 			Domain: "alpha",
// 			Priority: 5,
// 			Weight: 5,
// 		},
// 		{
// 			Domain: "beta",
// 			Priority: 1,
// 			Weight: 5,
// 		},
// 		{
// 			Domain: "alpha",
// 			Priority: 5,
// 			Weight: 1,
// 		},
// 		{
// 			Domain: "beta",
// 			Priority: 3,
// 			Weight: 5,
// 		},
// 	}
// }

// ProxyMiddleware should queue our incoming requests
func ProxyMiddleware(c iris.Context) {
	domain := c.GetHeader("domain")
	if len(domain) == 0 {
		c.JSON(iris.Map{"status": 400, "result": "error"})
		return
	}
	var repo Repository
	repo = &Queue{}
	fmt.Println("FROM HEADER", domain)

	for _, row := range repo.Read() {
		fmt.Println("FROM SOURCE", row.Domain)
	}
	Que = append(Que, domain)

	c.Next()
}