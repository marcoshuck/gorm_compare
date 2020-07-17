// Results:
// BenchmarkSingleQueryCreate
// BenchmarkSingleQueryCreate-4    	      22	  49663533 ns/op
//
// BenchmarkMultiQueryCreate
// BenchmarkMultiQueryCreate-4     	      24	  48899667 ns/op
//
// BenchmarkMultiQueryCreateTx
// BenchmarkMultiQueryCreateTx-4   	      58	  17333588 ns/op
package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"testing"
)

type Test struct {
	gorm.Model
	Name string
}

func setup() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:1234@(172.17.0.2)/test")
	if err != nil {
		return nil, err
	}
	db.DropTableIfExists(&Test{})
	db.AutoMigrate(&Test{})
	return db, nil
}

func BenchmarkSingleQueryCreate(b *testing.B) {
	db, err := setup()
	defer db.Close()
	if err != nil {
		b.Skip()
	}
	for n := 0; n < b.N; n++ {
		for i := 0; i < 3; i++ {
			q := db.Model(&Test{}).Save(&Test{
				Name:  fmt.Sprintf("Entity-%d", i),
			})
			err = q.Error
			if err != nil {
				b.Skip()
			}
		}
	}
}

func BenchmarkMultiQueryCreate(b *testing.B) {
	db, err := setup()
	defer db.Close()
	if err != nil {
		b.Skip()
	}
	for n := 0; n < b.N; n++ {
		for i := 0; i < 3; i++ {
			db = db.Save(&Test{
				Name:  fmt.Sprintf("Entity-%d", i),
			})
		}
		if err != nil {
			b.Skip()
		}
	}
}

func BenchmarkMultiQueryCreateTx(b *testing.B) {
	db, err := setup()
	defer db.Close()
	if err != nil {
		b.Skip()
	}
	for n := 0; n < b.N; n++ {
		err = db.Transaction(func(tx *gorm.DB) error {
			for i := 0; i < 3; i++ {
				tx = tx.Save(&Test{
					Name:  fmt.Sprintf("Entity-%d", i),
				})
			}
			return tx.Error
		})

		if err != nil {
			b.Skip()
		}
	}
}