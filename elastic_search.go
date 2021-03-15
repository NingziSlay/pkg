package pkg

import (
	"context"
	"github.com/olivere/elastic/v7"
)

// ESDoc es doc 对象
type ESDoc interface {
	Index() string
	ID() string
}

// ES simple wrapper of github.com/olivere/elastic/v7
type ES struct {
	client *elastic.Client
}

// NewES return a New ES instance
func NewES(debug bool, hosts ...string) (*ES, error) {
	opts := make([]elastic.ClientOptionFunc, 0)
	opts = append(opts, elastic.SetURL(hosts...))
	if debug {
		log := WithSampleLog()
		opts = append(opts, elastic.SetTraceLog(&log))
	} else {
		log := GetLogger()
		opts = append(opts, elastic.SetErrorLog(&log))
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &ES{client: client}, nil
}

// Client return an elastic.Client
func (es *ES) Client() *elastic.Client {
	return es.client
}

// BulkProcess start a bulk process
func (es *ES) BulkProcess(ctx context.Context, ch <-chan ESDoc, batch, worker int) (err error) {
	processor, err := es.client.BulkProcessor().BulkActions(batch).Workers(worker).Do(ctx)
	if err != nil {
		return
	}

	defer func() {
		err = processor.Flush()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case doc, ok := <-ch:
			// chan closed
			if !ok {
				return
			}
			bulk := elastic.NewBulkIndexRequest().Index(doc.Index()).Id(doc.ID()).RetryOnConflict(3).Doc(doc)
			processor.Add(bulk)
		}
	}
}
