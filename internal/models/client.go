// Package models contain the domais to persist in db
package models

type Client struct {
	ID    int
	Name  string
	Email string
	Phone string
}