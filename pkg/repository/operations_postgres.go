package repository

import (
	"CRUD_API"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB interface {
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type ProductPostgres struct {
	DB DB
}

func NewProductPostgres(DB *pgxpool.Pool) *ProductPostgres {
	return &ProductPostgres{DB: DB}
}

func (r *ProductPostgres) Create(name, description string) (CRUD_API.Products, error) {
	var products CRUD_API.Products
	var err error
	tx, err := r.DB.Begin(context.Background()) //транзакция-должно быть некс запросов-_-. они выполняются вместе,по отдельности-не выполнятся и отзовут изменения
	if err != nil {
		return products, fmt.Errorf("error starting transaction: %w", err)
	}
	request := "INSERT INTO products (name, description) VALUES ($1, $2) RETURNING id,name,description" //INSERT INTO-в какую таблицу нужно вставить новую запись
	dorequest := tx.QueryRow(context.Background(), request, name, description)
	/*контекст-срок выполнения операции ,тут дефолт(нет ограничений)*/
	if err := dorequest.Scan(&products.ID, &products.Name, &products.Description); err != nil { //присваиваем id наш result,чтобы в дальнейшем считывать его
		tx.Rollback(context.Background()) //откатываем изменения
		return products, fmt.Errorf("error in scan (func Create) bracho:(:%w", err)
	}
	if err := tx.Commit(context.Background()); err != nil { //все гуд-подтвердаем изменения
		tx.Rollback(context.Background())
		return products, fmt.Errorf("error committing transaction: %w", err)
	}
	return products, err
}

func (r *ProductPostgres) ReadAll() ([]CRUD_API.Products, error) {
	var products []CRUD_API.Products
	request := "SELECT * FROM products ORDER BY id ASC;"
	/*ORDER BY - сортировка, id ASC-поле id по возр(убывание-DESC)*/
	dorequest, err := r.DB.Query(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("err in request in getall:%s", err)
	}
	defer dorequest.Close() //закрываем в конце для предотвращения утечки памяти.
	/*цикл for для того,чтобы данные из запроса лежали в читаемой структуре,а затем присвоить нужной переменной.*/
	for dorequest.Next() { //next перебирает все строчки,пока они не кончатся.
		var product CRUD_API.Products                                                            //пустая структура
		if err := dorequest.Scan(&product.ID, &product.Name, &product.Description); err != nil { //в пустую struct передаем уже лежащие знач,для дальнейшего их считывания
			return nil, fmt.Errorf("Error in scan in getall:<: %s\n", err)
		}
		products = append(products, product) //добавляем наши полученные значения
	}
	if len(products) == 0 {
		fmt.Println("The products have not been created yet")
		return products, sql.ErrNoRows
	}
	return products, err
}

func (r *ProductPostgres) ReadById(id int) (CRUD_API.Products, error) {
	var products CRUD_API.Products
	var err error
	request := "SELECT id, name, description FROM products WHERE ID=$1"
	dorequest := r.DB.QueryRow(context.Background(), request, id)
	if err := dorequest.Scan(&products.ID, &products.Name, &products.Description); err != nil {
		return products, fmt.Errorf("err in readbyid:%s", err)
	}
	return products, err
}

func (r *ProductPostgres) Update(name, description string, id int) (CRUD_API.Products, error) {
	var product CRUD_API.Products
	var err error
	request := "UPDATE products SET name=$1, description=$2 WHERE ID=$3 RETURNING id, name, description" //set-конкретное новое значение,кот.нужно присвоить
	dorequest := r.DB.QueryRow(context.Background(), request, name, description, id)
	if err := dorequest.Scan(&product.ID, &product.Name, &product.Description); err != nil {
		return product, fmt.Errorf("err in update :%s", err)
	}
	return product, err
}

func (r *ProductPostgres) Delete(id int) error {
	request := "DELETE FROM products WHERE ID=$1"
	res, err := r.DB.Exec(context.Background(), request, id) //exec не возвращает строк результата
	if err != nil {
		return fmt.Errorf("err in delete:%s", err)
	}
	if res.RowsAffected() == 0 { //количество удаленных строк
		fmt.Println("id didn`t exists bro(")
		return sql.ErrNoRows
	}
	return err
}
