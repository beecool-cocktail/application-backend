package elastic

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/fatih/structs"
	"github.com/olivere/elastic/v7"
	"reflect"
	"strconv"
)

type cocktailsInfo struct {
	Title       string `structs:"title"`
	Description string `structs:"description"`
}

type cocktailElasticSearchRepository struct {
	es *elastic.Client
}

func NewElasticSearchCocktailRepository(es *elastic.Client) domain.CocktailElasticSearchRepository {
	return &cocktailElasticSearchRepository{es}
}

func (c *cocktailElasticSearchRepository) Index(ctx context.Context, co *domain.CocktailElasticSearch) error {

	_, err := c.es.Index().
		Index(domain.CocktailsIndex).
		Id(strconv.FormatInt(co.CocktailID, 10)).
		BodyJson(co).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (c *cocktailElasticSearchRepository) Search(ctx context.Context, text string,
	from, size int) ([]domain.CocktailElasticSearch, int64, error) {

	var result []domain.CocktailElasticSearch

	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewMatchQuery("description", text))
	res, _ := c.es.Search().
		Index(domain.CocktailsIndex).
		Query(boolQuery).
		From(from - 1).Size(size).
		Do(ctx)

	total := res.TotalHits()

	for _, item := range res.Each(reflect.TypeOf(domain.CocktailElasticSearch{})) {
		if t, ok := item.(domain.CocktailElasticSearch); ok {
			result = append(result, t)
		}
	}

	return result, total, nil
}

func (c *cocktailElasticSearchRepository) Update(ctx context.Context, co *domain.CocktailElasticSearch) error {

	updateField := cocktailsInfo{
		Title:       co.Title,
		Description: co.Description,
	}

	_, err := c.es.Update().
		Index(domain.CocktailsIndex).
		Id(strconv.FormatInt(co.CocktailID, 10)).
		Doc(structs.Map(updateField)).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (c *cocktailElasticSearchRepository) Delete(ctx context.Context, id int64) error {

	_, err := c.es.Delete().
		Index(domain.CocktailsIndex).
		Id(strconv.FormatInt(id, 10)).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}
