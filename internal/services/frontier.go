package services

import (
	"context"
	"fmt"
	"math"
	"net/url"

	"github.com/kordyd/go-crawler/internal/entities"
	"github.com/redis/go-redis/v9"
)

func Prioritizer(url entities.Url, queueDb *redis.Client, numberOfQueues int) error {

	queueIndex := int(math.Floor(url.Rank * float64(numberOfQueues)))

	if queueIndex >= numberOfQueues {
		queueIndex = numberOfQueues - 1
	}
	queueName := fmt.Sprintf("queue_%d", queueIndex)

	err := queueDb.LPush(context.Background(), queueName, url).Err()
	if err != nil {
		return fmt.Errorf("failed to push URL into queue: %w", err)
	}

	return nil

}

func Router(urlToRoute entities.Url, queueDb *redis.Client) error {
	host, err := host(urlToRoute)
	if err != nil {
		return fmt.Errorf("failed to get host: %w", err)
	}

	queueName := fmt.Sprintf("queue_%s", host)

	err = queueDb.LPush(context.Background(), queueName, urlToRoute).Err()
	if err != nil {
		return fmt.Errorf("failed to push URL into queue: %w", err)
	}

	return nil
}

func host(urlToGetHost entities.Url) (string, error) {
	parsedURL, err := url.Parse(urlToGetHost.Link)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}
	return parsedURL.Hostname(), nil
}
