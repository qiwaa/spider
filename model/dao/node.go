package dao

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	NodeNormal = iota
	NodeSeed
)

type NodeInfo struct {
	ID     int64
	Addr   string
	NodeId string
	Type   int
}

func (info *NodeInfo) TableName() string {
	return "node"
}

func (info *NodeInfo) Add() error {

	sqlstr := fmt.Sprintf("REPLACE INTO %s (addr,node_id) VALUES (?,?)", info.TableName())

	stmt, err := Prepare(sqlstr)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(info.Addr, info.NodeId); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (info *NodeInfo) List() ([]*NodeInfo, error) {
	sqlstr := fmt.Sprintf("SELECT id, addr, node_id, type From %s", info.TableName())

	stmt, err := Prepare(sqlstr)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.Wrap(err, "get data error")
	}
	defer rows.Close()

	items := make([]*NodeInfo, 0)
	for rows.Next() {
		var item NodeInfo
		err := rows.Scan(&item.ID, &item.Addr, &item.NodeId, &item.Type)
		if err != nil {
			logger.Errorf("get data error: %v\n", err)
			continue
		}
		items = append(items, &item)
	}

	err = rows.Err()
	if err != nil {
		logger.Errorf("get rows data error: %v\n", err)
	}

	return items, nil
}
