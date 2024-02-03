package repository

import (
	"CRUD_API"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v4/pgxpool"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestAll(t *testing.T) {
	t.Run("Create", TestProductPostgres_Create)
	t.Run("ReadAll", TestProductPostgres_ReadAll)
	t.Run("ReadById", TestProductPostgres_ReadById)
	t.Run("Update", TestProductPostgres_Update)
	t.Run("Delete", TestProductPostgres_Delete)
}

func TestProductPostgres_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mock.Close()
	r := ProductPostgres{mock}
	type args struct {
		name        string
		description string
	}
	type Mock func(args args)
	tests := []struct {
		name    string
		mock    Mock
		args    args
		want    CRUD_API.Products
		wantErr bool //по умолчанию-false
	}{
		{
			name: "OK",
			args: args{
				name:        "TestName",
				description: "TestDescription",
			},
			mock: func(args args) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO products").
					WithArgs(args.name, args.description).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description"}). //newRows создает новые столбцы в моковой бд
														AddRow(1, args.name, args.description)) //заполняет столбцы заданными аргументами
				mock.ExpectCommit()

			},

			want: CRUD_API.Products{
				ID:          1,
				Name:        "TestName",
				Description: "TestDescription",
			},
		},
		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		{
			name: "empty fields",
			args: args{
				name:        "",
				description: "",
			},
			mock: func(args args) {

				mock.ExpectQuery("INSERT INTO products").
					WithArgs(args.name, args.description).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description"}).
						AddRow(2, args.name, args.description).RowError(0, errors.New("some err in sec/test")))

			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { //запускает тесты по именам
			tt.mock(tt.args) //выполняем нашу моковую func
			got, err := r.Create(tt.args.name, tt.args.description)
			if tt.wantErr {
				assert.Error(t, err) //есть ли ошибка по факту
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet()) //проверяет,что все моки выполнились
		})
	}
}

func TestProductPostgres_ReadAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	r := ProductPostgres{mock}
	type args struct {
		name        string
		description string
	}
	type Mock func(args args)
	tests := []struct {
		name    string
		args    args
		mock    Mock
		wantErr bool
		want    []CRUD_API.Products
	}{
		{
			name: "OK",
			args: args{
				name:        "TestName",
				description: "TestDescription",
			},
			mock: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM products ORDER BY id ASC;")).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description"}).
						AddRow(1, args.name, args.description))
			},
			want: []CRUD_API.Products{
				{
					ID:          1,
					Name:        "TestName",
					Description: "TestDescription",
				},
			},
		},
		{
			name: "No Records",
			mock: func(args args) {
				//regexp.QuoteMeta для того,чтобы код понял символ ;
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM products ORDER BY id ASC;")).
					WillReturnRows(mock.NewRows([]string{"id", "name", "description"}))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)
			got, err := r.ReadAll()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductPostgres_ReadById(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	r := ProductPostgres{mock}
	type args struct {
		name        string
		description string
		id          int
	}
	type Mock func(args args)
	tests := []struct {
		mock    Mock
		name    string
		args    args
		want    CRUD_API.Products
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				name:        "ok_test",
				description: "ok-description",
			},
			mock: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description FROM products WHERE ID=$1")).
					WithArgs(args.id).
					WillReturnRows(mock.NewRows([]string{"id", "name", "description"}).
						AddRow(1, args.name, args.description))
			},
			want: CRUD_API.Products{
				ID:          1,
				Name:        "ok_test",
				Description: "ok-description",
			},
		},
		{
			name: "non-exist id",
			args: args{
				id: 321,
			},
			mock: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description FROM products WHERE ID=$1")).
					WithArgs(args.id).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)
			got, err := r.ReadById(tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())

		})

	}
}
func TestProductPostgres_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	type args struct {
		id          int
		name        string
		description string
	}
	r := ProductPostgres{mock}
	type Mock func(args args)
	tests := []struct {
		name    string
		args    args
		mock    Mock
		want    CRUD_API.Products
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				id:          1,
				name:        "OK",
				description: "OK",
			},
			mock: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE products SET name=$1, description=$2 WHERE ID=$3 RETURNING id, name, description")).
					WithArgs(args.name, args.description, args.id).
					WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description"}).
						AddRow(args.id, args.name, args.description))
			},
			want: CRUD_API.Products{
				ID:          1,
				Name:        "OK",
				Description: "OK",
			},
		},
		{
			name: "no records",
			mock: func(args args) {
				mock.ExpectQuery(regexp.QuoteMeta("UPDATE products SET name=$1, description=$2 WHERE ID=$3 RETURNING id, name, description")).
					WithArgs(args.name, args.description, args.id).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)
			got, err := r.Update(tt.args.name, tt.args.description, tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, tt.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestProductPostgres_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal()
	}
	defer mock.Close()
	r := ProductPostgres{mock}
	type args struct {
		id int
	}
	type Mock func(args args)
	tests := []struct {
		name    string
		args    args
		mock    Mock
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				id: 1,
			},
			mock: func(args args) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE ID=$1")).
					WithArgs(args.id).
					WillReturnResult(pgxmock.NewResult("DELETE", 1)) //имитация удаления строк
			},
		},
		{
			name: "no records found",
			args: args{
				id: 2,
			},
			mock: func(args args) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE ID=$1")).
					WithArgs(args.id).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.args)
			err := r.Delete(tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
