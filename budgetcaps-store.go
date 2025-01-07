package main

import (
	"database/sql"
	"fmt"
)

type BudgetCapsStorage interface {
	Create(*BudgetCaps) (*BudgetCaps, error)
	Delete(int64) (*BudgetCaps, error)
	Update(int64, *BudgetCaps) (*BudgetCaps, error)
	UpdateActive(int64, *BudgetCaps) (*BudgetCaps, error)
	GetDataId(int64) (*BudgetCaps, error)
	GetData() ([]*BudgetCaps, error)
}

type BudgetCapsStore struct {
	mysql *MysqlDB
}

func NewBudgetCapsStorage(db *MysqlDB) *BudgetCapsStore {
	return &BudgetCapsStore{
		mysql: db,
	}
}

func (m *BudgetCapsStore) GetData() ([]*BudgetCaps, error) {
	query := `SELECT id, budgets_id, budget_posts_id, amount, created_at, updated_at FROM budget_caps`
	rows, err := m.mysql.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget caps: %w", err)
	}
	defer rows.Close()

	var budgetCapsList []*BudgetCaps
	for rows.Next() {
		budgetCaps := &BudgetCaps{}
		err := rows.Scan(
			&budgetCaps.ID,
			&budgetCaps.BudgetsID,
			&budgetCaps.BudgetPostsID,
			&budgetCaps.Amount,
			&budgetCaps.CreatedAt,
			&budgetCaps.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget cap: %w", err)
		}
		budgetCapsList = append(budgetCapsList, budgetCaps)
	}
	return budgetCapsList, nil
}

func (m *BudgetCapsStore) GetDataId(id int64) (*BudgetCaps, error) {
	query := `SELECT id, budgets_id, budget_posts_id, amount, created_at, updated_at FROM budget_caps WHERE id = ?`
	row := m.mysql.db.QueryRow(query, id)

	budgetCaps := &BudgetCaps{}
	err := row.Scan(&budgetCaps.ID, &budgetCaps.BudgetsID, &budgetCaps.BudgetPostsID, &budgetCaps.Amount, &budgetCaps.CreatedAt, &budgetCaps.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get budget cap by id: %w", err)
	}
	return budgetCaps, nil
}

// udpate table
// CREATE TABLE budget_caps (
//     id BIGINT NOT NULL AUTO_INCREMENT,
//     budgets_id BIGINT NOT NULL,
//     budget_posts_id BIGINT NOT NULL,
//     amount DOUBLE,
//     created_at DATETIME,
//     updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//     PRIMARY KEY (id),
//     FOREIGN KEY (budgets_id) REFERENCES budgets(id),
//     FOREIGN KEY (budget_posts_id) REFERENCES budget_posts(id)

// CREATE TABLE IF NOT EXISTS `mydb`.`budget_caps` (
//     `id` INT NOT NULL AUTO_INCREMENT,
//     `budgets_id` VARCHAR(45) NULL,
//     `budget_posts_id` VARCHAR(45) NULL,
//     `amount` DOUBLE NULL DEFAULT 0,
//     `created_at` DATETIME NULL,
//     `updated_at` DATETIME NULL DEFAULT NOW(),
//     PRIMARY KEY (`id`),
//     INDEX `fk_budget_caps_budgets_idx` (`budgets_id` ASC) VISIBLE,
//     INDEX `fk_budget_caps_posts_idx` (`budget_posts_id` ASC) VISIBLE,
//     CONSTRAINT `fk_budget_caps_budgets`
//       FOREIGN KEY (`budgets_id`)
//       REFERENCES `mydb`.`budgets` (`name`)
//       ON DELETE NO ACTION
//       ON UPDATE NO ACTION,
//     CONSTRAINT `fk_budget_caps_posts`
//       FOREIGN KEY (`budget_posts_id`)
//       REFERENCES `mydb`.`budget_posts` (`name`)
//       ON DELETE NO ACTION
//       ON UPDATE NO ACTION)
//   ENGINE = InnoDB;

// );

func (m *BudgetCapsStore) Create(budgetCaps *BudgetCaps) (*BudgetCaps, error) {

	query := `INSERT INTO budget_caps (budgets_id, budget_posts_id, amount, created_at, updated_at) VALUES (?, ?, ?, now(), now())`
	result, err := m.mysql.db.Exec(query, budgetCaps.BudgetsID, budgetCaps.BudgetPostsID, budgetCaps.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to insert budget cap: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	newBudgetCaps, err := m.GetDataId(lastInsertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get new budget cap: %w", err)
	}

	return newBudgetCaps, nil
}

func (m *BudgetCapsStore) Delete(id int64) (*BudgetCaps, error) {
	deletedBudgetCaps, _ := m.GetDataId(id)

	query := `DELETE FROM budget_caps WHERE id = ?`
	_, err := m.mysql.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete budget cap: %w", err)
	}

	return deletedBudgetCaps, nil
}

func (m *BudgetCapsStore) Update(id int64, budgetCaps *BudgetCaps) (*BudgetCaps, error) {
	query := `UPDATE budget_caps SET budgets_id = ?, budget_posts_id = ?, amount = ?, updated_at = now() WHERE id = ?`
	_, err := m.mysql.db.Exec(query, budgetCaps.BudgetsID, budgetCaps.BudgetPostsID, budgetCaps.Amount, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget cap: %w", err)
	}

	updatedBudgetCaps, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated budget cap: %w", err)
	}

	return updatedBudgetCaps, nil
}

func (m *BudgetCapsStore) UpdateActive(id int64, budgetCaps *BudgetCaps) (*BudgetCaps, error) {
	query := `UPDATE budget_caps SET amount = ?, updated_at = now() WHERE id = ?`
	_, err := m.mysql.db.Exec(query, budgetCaps.Amount, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget cap: %w", err)
	}

	updatedBudgetCaps, err := m.GetDataId(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated budget cap: %w", err)
	}

	return updatedBudgetCaps, nil
}
