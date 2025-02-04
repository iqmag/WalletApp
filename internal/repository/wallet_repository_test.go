package repository_test

import (
    "context"
    "database/sql"
    "testing"

    _ "github.com/lib/pq"
    "WalletApp/internal/repository"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("postgres","user=postgres password=cyball dbname=postgres sslmode=disable")
	if err != nil {
	    return nil, err
	}

	// Создание таблицы для тестов
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS wallets (id UUID PRIMARY KEY,balance BIGINT NOT NULL DEFAULT 0)")
	return db, err
}

func TestPostgresWalletRepository(t *testing.T) {
	db, err := setupTestDB()
	if err != nil{
	    t.Fatalf("could not connect to database :%v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresWalletRepository(db)

	// Создаём кошелёк
	walletID, err := repo.CreateWallet(context.Background())
	if err != nil {
		t.Fatalf("could not create wallet: %v", err)
	}

	balance, err := repo.GetBalance(context.Background(), walletID)
	if err != nil || balance != 0 {
	    t.Fatalf("expected balance to be 0 but got %d; error:%v", balance, err)
    }

	err = repo.UpdateBalance(context.Background(), walletID ,1000)
	if err != nil {
	    t.Fatalf("could not update balance :%v", err)
    }

	balance, err = repo.GetBalance(context.Background(), walletID)
	if err != nil || balance != 1000 {
	    t.Fatalf("expected balance to be 1000 but got %d; error:%v", balance, err)
    }
}