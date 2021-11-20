package main

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type submission struct {
	gorm.Model
	SubmissionID  string `gorm:"unique"`
	URL           string
	Title         string
	Author        string
	AgreeCount    int `gorm:"default:0"`
	DisagreeCount int `gorm:"default:0"`
}

func random_submission(db *gorm.DB) *submission {
	sub := &submission{}
	log.Println("got random submission", db.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "RANDOM()", Vars: []interface{}{[]int{1, 2, 3}}, WithoutParentheses: true},
	}).Take(sub).Error)
	return sub
}

func agree(id int) {}

func disagree(id int) {}
