package repos

import (
	"crud/models"
	_ "database/sql"
)

type MarketRepo struct {
	db *PostgresDB
}

func NewMarketRepo(db *PostgresDB) *MarketRepo {
	return &MarketRepo{db: db}
}

func (r *MarketRepo) Create(market models.MarketEdit) (int, error) {
	var id int
	err := r.db.db.QueryRow(
		"INSERT INTO markets (name, address, phone_number) VALUES ($1, $2, $3) RETURNING id",
		market.Name, market.Address, market.PhoneNumber,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	for _, itemID := range market.ItemIds {
		_, err = r.db.db.Exec(
			"INSERT INTO market_items (market_id, item_id) VALUES ($1, $2)",
			id, itemID,
		)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (r *MarketRepo) Get(id int) (*models.MarketView, error) {
	market := &models.MarketView{}
	err := r.db.db.QueryRow(
		"SELECT name, address, phone_number FROM markets WHERE id = $1",
		id,
	).Scan(&market.Name, &market.Address, &market.PhoneNumber)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.db.Query(`
        SELECT i.name, i.count, i.price 
        FROM items i 
        JOIN market_items mi ON i.id = mi.item_id 
        WHERE mi.market_id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.ItemEdit
		if err := rows.Scan(&item.Name, &item.Count, &item.Price); err != nil {
			return nil, err
		}
		market.Items = append(market.Items, item)
	}
	return market, nil
}

func (r *MarketRepo) Update(id int, market models.MarketEdit) error {
	tx, err := r.db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE markets SET name = $1, address = $2, phone_number = $3 WHERE id = $4",
		market.Name, market.Address, market.PhoneNumber, id,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM market_items WHERE market_id = $1", id)
	if err != nil {
		return err
	}

	for _, itemID := range market.ItemIds {
		_, err = tx.Exec(
			"INSERT INTO market_items (market_id, item_id) VALUES ($1, $2)",
			id, itemID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *MarketRepo) Delete(id int) error {
	tx, err := r.db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM market_items WHERE market_id = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM markets WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
