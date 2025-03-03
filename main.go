package main

import (
	"log"
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	// Replace the connection details (user, dbname, password, host) with your own
	db, err := sqlx.Connect("postgres", "user=roa-uat-application dbname=port sslmode=disable password='penR40S!CP@B' host=localhost port=5432")
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

	queryMemoryUsage(db, "queryShort", queryShort)
	queryMemoryUsage(db, "queryStar", queryStar)
	select {}
}
func queryStar(db *sqlx.DB) {
	defer timeTrack(time.Now(), "queryStar")
	rows, _ := db.Query("SELECT * FROM fund_profile_txn LIMIT 50000")
	defer rows.Close()
}
func queryShort(db *sqlx.DB) {
	defer timeTrack(time.Now(), "queryShort")
	rows, _ := db.Query("SELECT id,fund_profile_txn_id FROM fund_profile_txn LIMIT 50000")
	defer rows.Close()
}
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
func queryMemoryUsage(db *sqlx.DB, name string, queryFunc func(*sqlx.DB)) {
	var memStart, memEnd runtime.MemStats

	runtime.GC() // Force GC before measuring memory
	runtime.ReadMemStats(&memStart)

	start := time.Now()
	queryFunc(db)
	elapsed := time.Since(start)

	runtime.ReadMemStats(&memEnd)

	log.Printf("%s took %s | Allocated Memory: %d KB | Total Allocations: %d KB | Heap In-Use: %d KB",
		name, elapsed,
		(memEnd.Alloc-memStart.Alloc)/1024,
		(memEnd.TotalAlloc-memStart.TotalAlloc)/1024,
		(memEnd.HeapInuse-memStart.HeapInuse)/1024,
	)
}
