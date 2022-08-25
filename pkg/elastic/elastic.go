package elastic

import (
	"context"
	"fmt"
	els "github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"skframe/pkg/logger"
	"strings"
	"sync"
)

type elastic struct {
	Client *els.Client
	Url    string
}

var Client *elastic
var once sync.Once

func ConnectElastic(url string) {
	once.Do(func() {
		Client = NewElasticClient(url)
	})
}

func NewElasticClient(url string) *elastic {
	client, err := els.NewClient(els.SetSniff(false), els.SetErrorLog(logger.NewElasticLogger()), els.SetURL(url))
	if err != nil {
		logger.Error("elastic", zap.Error(err))
		return nil
	}
	return &elastic{
		Client: client,
		Url:    url,
	}
}

func (ptr *elastic) CreateById(value interface{}, index string, id interface{}) bool {
	_, err := ptr.Client.Index().Index(index).Id(cast.ToString(id)).BodyJson(value).Do(context.Background())
	return err == nil
}

func (ptr *elastic) CreateBatch(value []interface{}, index string, ids []interface{}) bool {
	bulkRequest := ptr.Client.Bulk()
	for _index, val := range value {
		req := els.NewBulkIndexRequest().Index(index).Doc(val).Id(cast.ToString(ids[_index]))
		bulkRequest.Add(req)
	}
	_, err := bulkRequest.Do(context.Background())
	return err == nil
}

func (ptr *elastic) GetById(index string, id interface{}) []byte {
	res, _ := ptr.Client.Get().Index(index).Id(cast.ToString(id)).Do(context.Background())
	return res.Source
}

func (ptr *elastic) UpdateById(index string, update map[string]interface{}, id interface{}) bool {
	_, err := ptr.Client.Update().Index(index).Id(cast.ToString(id)).Doc(update).Do(context.Background())
	return err == nil
}

func (ptr *elastic) UpdateByQuery(index string, query els.Query, params map[string]interface{}) bool {
	var scriptArr []string
	for key, val := range params {
		scriptArr = append(scriptArr, fmt.Sprintf("ctx._source.%s=%v", key, val))
	}
	//query els.Query
	//boolQ := els.NewBoolQuery()
	//boolQ.Must(els.NewMatchQuery("last_name", "smith"))
	scriptStr := els.NewScript(strings.Join(scriptArr, ";"))
	_, err := ptr.Client.UpdateByQuery(index).Query(query).Script(scriptStr).Refresh("true").Do(context.Background())
	return err == nil
}

func (ptr *elastic) UpdateBatch(update []map[string]interface{}, index string, ids []interface{}) bool {
	bulkRequest := ptr.Client.Bulk()
	for _index, val := range update {
		req := els.NewBulkUpdateRequest().Index(index).Doc(val).Id(cast.ToString(ids[_index]))
		bulkRequest.Add(req)
	}
	_, err := bulkRequest.Do(context.Background())
	return err == nil
}

func (ptr *elastic) DeleteById(index string, id interface{}) bool {
	_, err := ptr.Client.Delete().Index(index).Id(cast.ToString(id)).Do(context.Background())
	return err == nil
}
func (ptr *elastic) DeleteByQuery(index string, query els.Query) bool {
	_, err := ptr.Client.DeleteByQuery(index).Query(query).Refresh("true").Do(context.Background())
	return err == nil
}
func (ptr *elastic) Search(index string, query els.Query, offset, size int) (resData [][]byte) {
	//q := els.NewQueryStringQuery("last_name:Smith")
	//boolQ := els.NewBoolQuery()
	//boolQ.Must(els.NewMatchQuery("last_name", "smith"))
	//boolQ.Filter(els.NewRangeQuery("age").Gt(35))
	////matchPhraseQuery := els.NewMatchPhraseQuery("about", "rock climbing")
	res, _ := ptr.Client.Search(index).Query(query).Size(size).From(offset).Do(context.Background())
	for _, item := range res.Hits.Hits {
		resData = append(resData, item.Source)
	}
	return
}

func (ptr *elastic) SearchAll(index string, query els.Query) (resData [][]byte) {
	res, _ := ptr.Client.Search(index).Query(query).Do(context.Background())
	for _, item := range res.Hits.Hits {
		fmt.Printf(string(item.Source))
		resData = append(resData, item.Source)
	}
	//res.TotalHits()
	return
}
