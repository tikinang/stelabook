package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"
)

const (
	TagQuotation = "bm.quotation"
	TagBookTitle = "doc.book-title"
	TagAuthors   = "ro.authors"
)

func main() {
	db, err := sql.Open("sqlite3", "books.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`
SELECT TagNames.TagName, Tags.Val, Tags.TimeEdt, Items.OID, Items.ParentID
FROM Tags
         JOIN Items ON Tags.ItemID = Items.OID
         JOIN TagNames ON Tags.TagID = TagNames.OID
WHERE TagNames.TagName IN (?,?,?)
`, TagQuotation, TagBookTitle, TagAuthors)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	bookMap := make(map[int64]*Book)
	for rows.Next() {
		var (
			tagName    string
			tagValue   string
			editedUnix int64
			itemId     int64
			parentId   sql.NullInt64
		)
		if err = rows.Scan(&tagName, &tagValue, &editedUnix, &itemId, &parentId); err != nil {
			panic(err)
		}

		var bookId int64
		if parentId.Valid {
			bookId = parentId.Int64
		} else {
			bookId = itemId
		}
		book := bookMap[bookId]
		if book == nil {
			book = &Book{
				Id: bookId,
			}
		}

		switch tagName {
		case TagQuotation:
			created := time.Unix(editedUnix, 0)
			var pbn PocketBookNote
			if err = json.Unmarshal([]byte(tagValue), &pbn); err != nil {
				panic(err)
			}
			book.Notes = append(book.Notes, Note{
				Created:  created,
				BeginRef: pbn.Begin,
				EndRef:   pbn.End,
				Content:  strings.ReplaceAll(pbn.Text, " ", "\n"),
			})
		case TagAuthors:
			book.Authors = tagValue
		case TagBookTitle:
			book.Title = tagValue
		}

		bookMap[book.Id] = book
	}

	var books []*Book
	for _, book := range bookMap {
		sort.Slice(book.Notes, func(i, j int) bool {
			return book.Notes[i].Created.Before(book.Notes[j].Created)
		})
		books = append(books, book)
	}
	sort.Slice(books, func(i, j int) bool {
		return books[i].Title < books[j].Title
	})

	mdFile, err := os.Create("books.md")
	if err != nil {
		panic(err)
	}
	defer mdFile.Close()

	mdTmpl, err := template.New("template").Parse(mdTemplateFile)
	if err != nil {
		panic(err)
	}
	err = mdTmpl.Execute(mdFile, books)
	if err != nil {
		panic(err)
	}
}

//go:embed md.tmpl
var mdTemplateFile string

type Book struct {
	Id      int64
	Title   string
	Authors string
	Notes   []Note
}

type Note struct {
	Created  time.Time
	BeginRef string
	EndRef   string
	Content  string
}

func (r *Note) CreatedFormatted() string {
	return r.Created.Format(time.DateTime)
}

type PocketBookNote struct {
	Begin string `json:"begin"`
	End   string `json:"end"`
	Text  string `json:"text"`
}
