package service

import (
	"context"
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/olivere/elastic/v7"
)

func newElasticSearch(configure *Configure) (*elastic.Client, error) {
	if configure.Elastic == nil {
		return nil, errors.New("elastic search configure is not initialed")
	}

	client, err := elastic.NewClient(elastic.SetURL(configure.Elastic.Urls...))
	if err != nil {
		return nil, err
	}

	exist, err := client.IndexExists(domain.CocktailsIndex).Do(context.Background())
	if err != nil {
		return nil, err
	}

	if !exist {
		indexCreate, err := client.CreateIndex(domain.CocktailsIndex).
			BodyString(domain.CocktailsMapping).
			Do(context.Background())
		if err != nil {
			return nil, err
		}

		if !indexCreate.Acknowledged {
			return nil, errors.New("elastic search index create failed")
		}
	}

	return client, nil
}
