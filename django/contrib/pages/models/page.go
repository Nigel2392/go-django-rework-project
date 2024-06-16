package models

import (
	"context"
	"database/sql"
)

type PageNode struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Path     string `json:"path"`
	Depth    int64  `json:"depth"`
	Numchild int64  `json:"numchild"`
	PageID   int64  `json:"page_id"`
	Typehash string `json:"typehash"`
}

type DBQuerier interface {
	Querier
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Querier interface {
	WithTx(tx *sql.Tx) Querier
	DeleteNode(ctx context.Context, id int64) error
	GetChildren(ctx context.Context, path interface{}, depth interface{}) ([]PageNode, error)
	GetDescendants(ctx context.Context, path interface{}, depth int64) ([]PageNode, error)
	GetNodeByID(ctx context.Context, id int64) (PageNode, error)
	GetNodeByPath(ctx context.Context, path string) (PageNode, error)
	InsertNode(ctx context.Context, title string, path string, depth int64, numchild int64, page_id int64, typehash string) (int64, error)
	UpdateNode(ctx context.Context, title string, path string, depth int64, numchild int64, page_id int64, typehash string, iD int64) error
	UpdateNodePathAndDepth(ctx context.Context, path string, depth int64, iD int64) error
}

/// MoveNodeParams contains parameters for moving a node.
//ype MoveNodeParams struct {
//ID          int64
//NewPath     string
//NewDepth    int64
//NewParentID sql.NullInt64
//
/// MoveNode moves a node within the tree.
//unc MoveNode(ctx context.Context, q DBQuerier, params MoveNodeParams) error {
//// Start a transaction
//var tx *sql.Tx
//var err error
//if tx, err = q.BeginTx(ctx); err != nil {
//	return err
//}
//var queries = q.WithTx(tx)
//defer tx.Rollback()
//// Lock the table
//if _, err := tx.ExecContext(ctx, "LOCK TABLE PageNode IN EXCLUSIVE MODE"); err != nil {
//	return err
//}
//// Fetch the node and its children
//nodes, err := queries.GetNodeWithChildren(ctx, params.ID)
//if err != nil {
//	return err
//}
//// Calculate new paths and depths
//oldPath := nodes[0].Path.String
//newPath := params.NewPath
//pathPrefix := oldPath + "."
//for _, node := range nodes {
//	relativePath := node.Path.String[len(pathPrefix):]
//	newChildPath := newPath + "." + relativePath
//	newDepth := params.NewDepth + (node.Depth.Int64 - nodes[0].Depth.Int64)
//	var newPathSQL sql.NullString
//	var newDepthSQL sql.NullInt64
//	if node.ID == params.ID {
//		newPathSQL = sql.NullString{String: newPath, Valid: true}
//		newDepthSQL = sql.NullInt64{Int64: newDepth, Valid: true}
//	} else {
//		newPathSQL = sql.NullString{String: newChildPath, Valid: true}
//		newDepthSQL = sql.NullInt64{Int64: newDepth, Valid: true}
//	}
//	// Update node path and depth
//	if node.ID == params.ID {
//		if err := queries.UpdateNodePathAndDepth(ctx, node.ID, newPathSQL, newDepthSQL, params.NewParentID); err != nil {
//			return err
//		}
//	} else {
//		if err := queries.UpdateChildNode(ctx, node.ID, newPathSQL, newDepthSQL); err != nil {
//			return err
//		}
//	}
//}
//// Commit the transaction
//if err := tx.Commit(); err != nil {
//	return err
//}
//return nil
//
//
