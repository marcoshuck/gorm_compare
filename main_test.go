// Results:
//goos: linux
//goarch: amd64
//pkg: github.com/marcoshuck/gorm_compare
//BenchmarkSingleQueryCreate
//BenchmarkSingleQueryCreate-4    	      20	  51741170 ns/op
//BenchmarkMultiQueryCreate
//BenchmarkMultiQueryCreate-4     	      24	  45862690 ns/op
//BenchmarkMultiQueryCreateTx
//BenchmarkMultiQueryCreateTx-4   	      56	  18557073 ns/op
//BenchmarkMultiQueryUpdate
//BenchmarkMultiQueryUpdate-4     	       1	1679886780 ns/op
//BenchmarkSingleQueryUpdate
//BenchmarkSingleQueryUpdate-4    	       1	1693897196 ns/op
//PASS
//
//Process finished with exit code 0
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
	db, err := gorm.Open("mysql", "root:1234@(172.17.0.2)/test?parseTime=true")
	if err != nil {
		return nil, err
	}
	db.DropTableIfExists(&Test{})
	db.AutoMigrate(&Test{})
	return db, nil
}

func populate(db *gorm.DB) {
	for i := 0; i < 100; i++ {
		db.Model(&Test{}).Create(&Test{
			Name:  fmt.Sprintf("Entity-%d", i),
		})
	}
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
			db.Save(&Test{
				Name:  fmt.Sprintf("Entity-%d", i),
			})
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
		b.StartTimer()
		err = db.Transaction(func(tx *gorm.DB) error {
			for i := 0; i < 3; i++ {
				tx = tx.Save(&Test{
					Name:  fmt.Sprintf("Entity-%d", i),
				})
			}
			return tx.Error
		})
		b.StopTimer()
		if err != nil {
			b.Skip()
		}
	}
}

func BenchmarkMultiQueryUpdate(b *testing.B) {
	db, err := setup()
	defer db.Close()
	if err != nil {
		b.Skip()
	}
	for n := 0; n < b.N; n++ {
		db.DropTableIfExists(&Test{})
		db.AutoMigrate(&Test{})
		populate(db)
		b.StartTimer()
		var t Test
		db.Model(&Test{}).Where("name = ?", "Entity-1").First(&t)
		t.Name = "Tested"
		db.Model(&Test{}).Where("name = ?", "Entity-1").Update(&t)
		b.StopTimer()
	}
}

func BenchmarkSingleQueryUpdate(b *testing.B) {
	db, err := setup()
	defer db.Close()
	if err != nil {
		b.Skip()
	}
	for n := 0; n < b.N; n++ {
		db.DropTableIfExists(&Test{})
		db.AutoMigrate(&Test{})
		populate(db)
		b.StartTimer()
		t := Test{
			Name: "Tested",
		}
		db.Model(&Test{}).Where("name = ?", "Entity-1").Update(&t)
		b.StopTimer()
	}
}