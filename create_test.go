package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
)

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
