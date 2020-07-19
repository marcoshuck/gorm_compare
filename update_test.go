package main

import (
	"fmt"
	"testing"
)

func BenchmarkMultiQueryUpdate(b *testing.B) {
	b.StopTimer()
	db, err := setup()
	defer db.Close()
	if err != nil {
		b.Skip()
	}
	for n := 0; n < b.N; n++ {
		for i := 0; i < 10; i++ {
			b.StopTimer()
			db.DropTableIfExists(&Test{})
			db.AutoMigrate(&Test{})
			populate(db)
			b.StartTimer()
			var t Test
			db.Model(&Test{}).Where("name = ?", fmt.Sprintf("Entity-%d", i)).First(&t)
			t.Name = "Tested"
			db.Model(&Test{}).Where("name = ?", fmt.Sprintf("Entity-%d", i)).Update(&t)
		}
	}
}

func BenchmarkSingleQueryUpdate(b *testing.B) {
	b.StopTimer()
	db, err := setup()
	defer db.Close()
	if err != nil {
		b.Skip()
	}
	for n := 0; n < b.N; n++ {
		for i := 0; i < 10; i++ {
			b.StopTimer()
			db.DropTableIfExists(&Test{})
			db.AutoMigrate(&Test{})
			populate(db)
			b.StartTimer()
			db.Model(&Test{}).Where("name = ?", fmt.Sprintf("Entity-%d", i)).Update(&Test{ Name: "Tested" })
		}
	}
}