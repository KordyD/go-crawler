package main

import (
	"log"

	"github.com/kordyd/go-crawler/internal/db/mongodb"
	"github.com/kordyd/go-crawler/internal/db/redis"
	"github.com/kordyd/go-crawler/internal/services"
)

// func main() {

// 	client := db.Connect()
// 	ctx := context.Background()
// 	c := colly.NewCollector()

// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"task_queue", // name
// 		false,        // durable
// 		false,        // delete when unused
// 		false,        // exclusive
// 		false,        // no-wait
// 		nil,          // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	msgs, err := ch.Consume(
// 		q.Name, // queue
// 		"",     // consumer
// 		true,   // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	)
// 	failOnError(err, "Failed to register a consumer")

// 	var forever chan struct{}

// 	go func() {
// 		for d := range msgs {
// 			log.Printf("Received a message: %s", d.Body)
// 			c.Visit(string(d.Body))
// 			// dotCount := bytes.Count(d.Body, []byte("."))
// 			// t := time.Duration(dotCount)
// 			// time.Sleep(t * time.Second)
// 			log.Printf("Done")
// 		}
// 	}()

// 	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

// 	// Find and visit all links
// 	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
// 		// e.Request.Visit(e.Attr("href"))
// 		// fmt.Println("Visited", e.Attr("href"))
// 		set := client.Set(ctx, e.Attr("href"), 1, 0)
// 		if set.Err() != nil {
// 			log.Println(set.Err())
// 		}

// 	})

// 	// c.OnRequest(func(r *colly.Request) {
// 	// 	fmt.Println("Visiting", r.URL)
// 	// })

// 	c.OnResponse(func(r *colly.Response) {
// 		fmt.Println("Visited", r.Request.URL)
// 		set := client.Set(ctx, r.Request.URL.String(), r.Body, 0)
// 		if set.Err() != nil {
// 			log.Println(set.Err())
// 		}
// 	})

// 	c.OnError(func(r *colly.Response, err error) {
// 		log.Println(err)
// 	})

// 	<-forever

// }

// func main() {
// 	client := db.Connect()
// 	ctx := context.Background()

// 	urls := []string{"https://redis.uptrace.dev/", "https://go-colly.org/"}
// 	parsedUrls := make(chan string, 5)
// 	fetchedBody := make(chan string, 5)

// 	var wg sync.WaitGroup

// 	for _, url := range urls {
// 		wg.Add(1)
// 		go scrapper.Scrapper(url, parsedUrls, fetchedBody, &wg)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(parsedUrls)
// 		close(fetchedBody)
// 	}()

// 	for url := range parsedUrls {
// 		client.Set(ctx, url, 1, 0).Result()
// 	}

// 	for body := range fetchedBody {
// 		fmt.Println(body)
// 	}

// }

func main() {
	// https://go-colly.org/

	// base, err := url.Parse("http://go-colly.org/")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(base.Scheme)

	// u, err := url.Parse("/lol")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// test_url := base.ResolveReference(u)

	// fmt.Println(test_url)

	// normilized := purell.NormalizeURL(test_url, purell.FlagsUnsafeGreedy)

	// fmt.Println(normilized)

	mongoDBName := "Crawler"
	mongoCollectionName := "URLs"

	client, disconnect, err := mongodb.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	redisClient := redis.Connect()

	coll := client.Database(mongoDBName).Collection(mongoCollectionName)

	services.SaveParsedUrls(coll, redisClient)

}
