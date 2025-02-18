package repos

import (
	"crud/models"
	"database/sql"
)

type ItemRepo struct {
	db *PostgresDB
}

func NewItemRepo(db *PostgresDB) *ItemRepo {
	return &ItemRepo{db: db}
}

func (r *ItemRepo) Create(item models.Item) (int, error) {
	var id int
	err := r.db.db.QueryRow(
		"INSERT INTO items (name, count, price) VALUES ($1, $2, $3) RETURNING id",
		item.Name, item.Count, item.Price,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ItemRepo) GetByID(id int) (*models.Item, error) {
	item := &models.Item{}
	err := r.db.db.QueryRow(
		"SELECT id, name, count, price FROM items WHERE id = $1",
		id,
	).Scan(&item.Id, &item.Name, &item.Count, &item.Price)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *ItemRepo) GetAll() ([]models.Item, error) {
	rows, err := r.db.db.Query("SELECT id, name, count, price FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Id, &item.Name, &item.Count, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemRepo) Update(id int, item models.ItemEdit) error {
	result, err := r.db.db.Exec(
		"UPDATE items SET name = $1, count = $2, price = $3 WHERE id = $4",
		item.Name, item.Count, item.Price, id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ItemRepo) Delete(id int) error {
	result, err := r.db.db.Exec("DELETE FROM items WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ItemRepo) GetByMarketID(marketID int) ([]models.Item, error) {
	query := `
        SELECT i.id, i.name, i.count, i.price 
        FROM items i 
        JOIN market_items mi ON i.id = mi.item_id 
        WHERE mi.market_id = $1
    `

	rows, err := r.db.db.Query(query, marketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Id, &item.Name, &item.Count, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
