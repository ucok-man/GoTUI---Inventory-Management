package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

type Item struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

type InventoryModel struct {
	srcfile     string
	inventories []Item
}

func NewInventoryModel(srcfile string) (*InventoryModel, error) {
	model := &InventoryModel{
		srcfile:     srcfile,
		inventories: []Item{},
	}

	var patherr *fs.PathError
	_, err := os.Stat(srcfile)
	if err != nil && !errors.As(err, &patherr) {
		return nil, fmt.Errorf("error acessing %s file: %v", srcfile, err)
	}

	if err != nil && errors.As(err, &patherr) {
		dir, _ := path.Split(srcfile)
		if dir != "" {
			if err := os.MkdirAll(dir, 0744); err != nil {
				return nil, fmt.Errorf("error creating %v directory: %v", dir, err)
			}
		}

		_, err = os.OpenFile(srcfile, os.O_CREATE, 0644)
		if err != nil {
			return nil, fmt.Errorf("error creating %v file: %v", srcfile, err)
		}

		if err := model.Save(); err != nil {
			return nil, err
		}
	}

	if err := model.Load(); err != nil {
		return nil, err
	}

	return model, nil
}

func (m *InventoryModel) Load() error {
	data, err := os.ReadFile(m.srcfile)
	if err != nil {
		return fmt.Errorf("err reading %v file: %v", m.srcfile, err)
	}

	if err := json.Unmarshal(data, &m.inventories); err != nil {
		return fmt.Errorf("err unmarshal %v file: %v", m.srcfile, err)
	}
	return nil
}

func (m *InventoryModel) Save() error {
	data, err := json.MarshalIndent(&m.inventories, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}
	return os.WriteFile(m.srcfile, data, 0644)
}

func (m *InventoryModel) Get() []Item {
	return m.inventories
}

func (m *InventoryModel) Add(item Item) error {
	if item.Stock < 0 {
		return fmt.Errorf("error stock must be positive number")
	}

	m.inventories = append(m.inventories, item)
	return m.Save()
}

// index start from 1
func (m *InventoryModel) Delete(index int) error {
	if index < 0 || index > len(m.inventories) {
		return fmt.Errorf("error invalid index parameter")
	}
	m.inventories = append(m.inventories[:index-1], m.inventories[index:]...)
	return m.Save()
}
