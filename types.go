package couchdb

import "strings"

const langJavaScript = "javascript"

type QueryParameters struct {
	Conflicts       *bool   `url:"conflicts,omitempty"`
	Descending      *bool   `url:"descending,omitempty"`
	Group           *bool   `url:"group,omitempty"`
	IncludeDocs     *bool   `url:"include_docs,omitempty"`
	Attachments     *bool   `url:"attachments,omitempty"`
	AttEncodingInfo *bool   `url:"att_encoding_info,omitempty"`
	InclusiveEnd    *bool   `url:"inclusive_end,omitempty"`
	Reduce          *bool   `url:"reduce,omitempty"`
	UpdateSeq       *bool   `url:"update_seq,omitempty"`
	GroupLevel      *int    `url:"group_level,omitempty"`
	Limit           *int    `url:"limit,omitempty"`
	Skip            *int    `url:"skip,omitempty"`
	Key             *string `url:"key,omitempty"`
	EndKey          *string `url:"endkey,comma,omitempty"`
	EndKeyDocID     *string `url:"end_key_doc_id,omitempty"`
	Stale           *string `url:"stale,omitempty"`
	StartKey        *string `url:"startkey,comma,omitempty"`
	StartKeyDocID   *string `url:"startkey_docid,omitempty"`
}

type ViewResponse struct {
	Offset    int   `json:"offset,omitempty"`
	Rows      []Row `json:"rows,omitempty"`
	TotalRows int   `json:"total_rows,omitempty"`
	UpdateSeq int   `json:"update_seq,omitempty"`
}

type Row struct {
	ID    string                 `json:"id"`
	Key   interface{}            `json:"key"`
	Value interface{}            `json:"value,omitempty"`
	Doc   map[string]interface{} `json:"doc,omitempty"`
}

// DesignDocument is a special type of CouchDB document that contains application code.
// http://docs.couchdb.org/en/latest/json-structure.html#design-document
type DesignDocument struct {
	Document
	Language string                        `json:"language,omitempty"`
	Views    map[string]DesignDocumentView `json:"views,omitempty"`
	Filters  map[string]string             `json:"filters,omitempty"`
}

// Name returns design document name without the "_design/" prefix
func (dd DesignDocument) Name() string {
	return strings.TrimPrefix(dd.ID, "_design/")
}

// DesignDocumentView contains map/reduce functions.
type DesignDocumentView struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`
}

// CouchDoc describes interface for every couchdb document.
type CouchDoc interface {
	GetID() string
	GetRev() string
}

// Document is base struct which should be embedded by any other couchdb document.
type Document struct {
	ID          string                `json:"_id,omitempty"`
	Rev         string                `json:"_rev,omitempty"`
	Attachments map[string]Attachment `json:"_attachments,omitempty"`
}

// Attachment describes attachments of a document.
// http://docs.couchdb.org/en/stable/api/document/common.html#attachments
// By using attachments you are also able to upload a document in multipart/related format.
// http://docs.couchdb.org/en/latest/api/document/common.html#creating-multiple-attachments
type Attachment struct {
	ContentType   string  `json:"content_type,omitempty"`
	Data          string  `json:"data,omitempty"`
	Digest        string  `json:"digest,omitempty"`
	EncodedLength float64 `json:"encoded_length,omitempty"`
	Encoding      string  `json:"encoding,omitempty"`
	Length        int64   `json:"length,omitempty"`
	RevPos        float64 `json:"revpos,omitempty"`
	Stub          bool    `json:"stub,omitempty"`
	Follows       bool    `json:"follows,omitempty"`
}

// GetID returns document id
func (d *Document) GetID() string {
	return d.ID
}

// GetRev returns document revision
func (d *Document) GetRev() string {
	return d.Rev
}

// DocumentResponse is response for multipart/related file upload.
type DocumentResponse struct {
	Ok  bool
	ID  string
	Rev string
}

type PurgeResponse struct {
	PurgeSeq float64 `json:"purge_seq"`
	Purged   map[string][]string
}

// BulkDoc describes POST /db/_bulk_docs request object.
// http://docs.couchdb.org/en/latest/api/database/bulk-api.html#post--db-_bulk_docs
type BulkDoc struct {
	AllOrNothing bool       `json:"all_or_nothing,omitempty"`
	NewEdits     bool       `json:"new_edits,omitempty"`
	Docs         []CouchDoc `json:"docs"`
}

type FindArgs struct {
	// http://docs.couchdb.org/en/stable/api/database/find.html
	Selector       map[string]interface{} `json:"selector,omitempty"`        // – 选择器， 查询条件参数 JSON object describing criteria used to select documents. More information provided in the section on selector syntax. Required
	Limit          int64                  `json:"limit,omitempty"`           //  – 查询条数，配合bookmark 使用可以达到分页效果  Maximum number of results returned. Default is 25. Optional
	Skip           int64                  `json:"skip,omitempty"`            // – Skip the first ‘n’ results, where ‘n’ is the value specified. Optional
	Sort           []interface{}          `json:"sort,omitempty"`            // – 排序 JSON array following sort syntax. Optional
	Fields         []interface{}          `json:"fields,omitempty"`          // – 过滤 JSON array specifying which fields of each object should be returned. If it is omitted, the entire object is returned. More information provided in the section on filtering fields. Optional
	UseIndex       string                 `json:"use_index,omitempty"`       //  索引( string|array)  – Instruct a query to use a specific index. Specified either as "<design_document>" or ["<design_document>", "<index_name>"]. Optional
	R              int64                  `json:"r,omitempty"`               // (number)   – Read quorum needed for the result. This defaults to 1, in which case the document found in the index is returned. If set to a higher value, each document is read from at least that many replicas before it is returned in the results. This is likely to take more time than using only the document stored locally with the index. Optional, default: 1
	Bookmark       string                 `json:"bookmark,omitempty"`        //  – A string that enables you to specify which page of results you require. Used for paging through result sets. Every query returns an opaque string under the bookmark key that can then be passed back in a query to get the next page of results. If any part of the selector query changes between requests, the results are undefined. Optional, default: null
	Update         bool                   `json:"update,omitempty"`          //(boolean) – Whether to update the index prior to returning the result. Default is true. Optional
	Stable         bool                   `json:"stable,omitempty"`          //(boolean)  – Whether or not the view results should be returned from a “stable” set of shards. Optional
	Stale          string                 `json:"stale,omitempty"`           // – Combination of update=false and stable=true options. Possible options: "ok", false (default). Optional
	ExecutionStats bool                   `json:"execution_stats,omitempty"` //  在查询响应中包含执行统计信息 (boolean)  – Include execution statistics in the query response. Optional, default: ``false``\
}

type CouchSelectorBody struct {
	// Docs           []CouchDoc     `json:"docs,omitempty"`            // – JSON object describing criteria used to select documents. More information provided in the section on selector syntax. Required
	BookMark       string         `json:"bookmark,omitempty"`        //  – 配合limit使用， 查询一次返回的， 下次再查询传入这个查询的就是下一页， 分页效果
	ExecutionStats ExecutionStats `json:"execution_stats,omitempty"` // —在查询响应中包含执行统计信息
	Warning        string         `json:"warning,omitempty"`         // – 异常信息
}

type ExecutionStats struct {
	ExecutionTimeMs         float64 `json:"execution_time_ms,omitempty"`          // 执行时间
	ResultsReturned         int     `json:"results_returned,omitempty"`           // 结果返回
	TotalDocsExamined       int     `json:"total_docs_examined,omitempty"`        // 已检查的文档总数
	TotalKeysExamined       int     `json:"total_keys_examined,omitempty"`        // 已检查的总键数
	TotalQuorumDocsExamined int     `json:"total_quorum_docs_examined,omitempty"` // 已审核的法定文档总数
}

/*
{
	"index": {
		"partial_filter_selector": {
			"Attribute.UserHospitalizedCases": {
				"$in": [{
					"Name": "刘"
				}]
			}
		},
		 "fields": ["UserKey"]
	},
	"name": "0x00ifhd72832jfsuajzi",
	"type": "json"
}
*/

type Index struct {
	Index interface{} `json:"index,omitempty"` //fields []string,  partial_filter_selector  interface{}               // 描述要创建的索引的json对象。(json) – JSON object describing the index to create.
	DDoc  string      `json:"ddoc,omitempty"`  // (string) – Name of the design document in which the index will be created. By default, each index will be created in its own design document. Indexes can be grouped into design documents for efficiency. However, a change to one index in a design document will invalidate all other indexes in the same document (similar to views). Optional
	Name  string      `json:"name,omitempty"`  // (string) – Name of the index. If no name is provided, a name will be generated automatically. Optional
	Type  string      `json:"type,omitempty"`  // (string) – Can be "json" or "text". Defaults to json. Geospatial indexes will be supported in the future. Optional Text indexes are supported via a third party library Optional
}

type CouchIndexBody struct {
	Result string `json:"result,omitempty"` //(string) – Flag to show whether the index was created or one already exists. Can be “created” or “exists”.
	Id     string `json:"id,omitempty"`     //(string) – Id of the design document the index was created in.
	Name   string `json:"name,omitempty"`   //(string) – Name of the index created.
}
