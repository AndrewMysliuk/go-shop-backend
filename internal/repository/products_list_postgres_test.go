package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/AndrewMislyuk/go-shop-backend/internal/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestProductsListPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	r := NewProductsListPostgres(db)

	type args struct {
		productId string
		item      domain.CreateProductInput
		createdAt time.Time
	}
	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				productId: "34c8d3e6-b8d7-43dc-847e-5764c4114856",
				item: domain.CreateProductInput{
					Title:        "Твидовый кардиган из хлопка",
					Image:        "w1.webp",
					Price:        749000,
					Sale:         0,
					SaleOldPrice: 0,
					Category:     "Женщинам",
					Type:         "Одежда",
					Subtype:      "Старые-коллекции",
					Description:  "",
				},
				createdAt: time.Now(),
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(args.productId)
				prep := mock.ExpectPrepare("INSERT INTO products")
				prep.ExpectQuery().
					WithArgs(args.productId, args.item.Title, args.item.Image, args.item.Price, args.item.Sale, args.item.SaleOldPrice, args.item.Category, args.item.Type, args.item.Subtype, args.item.Description, args.createdAt).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
		},

		{
			name: "Empty Fields",
			args: args{
				productId: "",
				item: domain.CreateProductInput{
					Title:        "",
					Image:        "",
					Price:        0,
					Sale:         0,
					SaleOldPrice: 0,
					Category:     "",
					Type:         "",
					Subtype:      "",
					Description:  "",
				},
				createdAt: time.Now(),
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(args.productId).RowError(0, errors.New("insert error"))
				prep := mock.ExpectPrepare("INSERT INTO products")
				prep.ExpectQuery().
					WithArgs(args.productId, args.item.Title, args.item.Image, args.item.Price, args.item.Sale, args.item.SaleOldPrice, args.item.Category, args.item.Type, args.item.Subtype, args.item.Description, args.createdAt).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.Create(testCase.args.item, testCase.args.productId, testCase.args.createdAt)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.args.productId, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductsListPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	r := NewProductsListPostgres(db)

	testTable := []struct {
		name    string
		mock    func()
		want    []domain.ProductsList
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "image", "price", "sale", "sale_old_price", "category", "type", "subtype", "description", "created_at"}).
					AddRow("453b4f0f-1f56-4c57-b43d-7b79792450a7", "Твидовый кардиган из хлопка", "w1.webp", 749000, 0, 0, "Женщинам", "Одежда", "Старые-коллекции", "", time.Date(2022, 01, 12, 13, 8, 21, 32963, time.Local)).
					AddRow("b07221f8-4133-4688-b2d6-d677f41f5b74", "Объемный водоотталкивающий тренч", "w2.webp", 499000, 50, 999000, "Женщинам", "Одежда", "Старые-коллекции", "", time.Date(2022, 01, 12, 13, 10, 18, 882593, time.Local)).
					AddRow("96a7193a-403d-4e01-94e6-c02c5bcb61f1", "Хлопковая рубашка в полоску", "w4.webp", 359000, 0, 0, "Женщинам", "Одежда", "Вышевка", "", time.Date(2022, 01, 12, 13, 16, 55, 558842, time.Local))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM products")).WillReturnRows(rows)
			},
			want: []domain.ProductsList{
				{Id: "453b4f0f-1f56-4c57-b43d-7b79792450a7", Title: "Твидовый кардиган из хлопка", Image: "w1.webp", Price: 749000, Sale: 0, SaleOldPrice: 0, Category: "Женщинам", Type: "Одежда", Subtype: "Старые-коллекции", Description: "", CreatedAt: time.Date(2022, 01, 12, 13, 8, 21, 32963, time.Local)},
				{Id: "b07221f8-4133-4688-b2d6-d677f41f5b74", Title: "Объемный водоотталкивающий тренч", Image: "w2.webp", Price: 499000, Sale: 50, SaleOldPrice: 999000, Category: "Женщинам", Type: "Одежда", Subtype: "Старые-коллекции", Description: "", CreatedAt: time.Date(2022, 01, 12, 13, 10, 18, 882593, time.Local)},
				{Id: "96a7193a-403d-4e01-94e6-c02c5bcb61f1", Title: "Хлопковая рубашка в полоску", Image: "w4.webp", Price: 359000, Sale: 0, SaleOldPrice: 0, Category: "Женщинам", Type: "Одежда", Subtype: "Вышевка", Description: "", CreatedAt: time.Date(2022, 01, 12, 13, 16, 55, 558842, time.Local)},
			},
		},
		{
			name: "No Records",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "image", "price", "sale", "sale_old_price", "category", "type", "subtype", "description", "created_at"})

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM products")).WillReturnRows(rows)
			},
			want: []domain.ProductsList{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.GetAll()
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductsListPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Fatal(err)
	}

	r := NewProductsListPostgres(db)

	type args struct {
		productId string
	}

	testTable := []struct {
		name    string
		mock    func()
		args    args
		want    domain.ProductsList
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "image", "price", "sale", "sale_old_price", "category", "type", "subtype", "description", "created_at"}).
					AddRow("453b4f0f-1f56-4c57-b43d-7b79792450a7", "Твидовый кардиган из хлопка", "w1.webp", 749000, 0, 0, "Женщинам", "Одежда", "Старые-коллекции", "", time.Date(2022, 01, 12, 13, 8, 21, 32963, time.Local))

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM products WHERE id = $1")).WithArgs("453b4f0f-1f56-4c57-b43d-7b79792450a7").WillReturnRows(rows)
			},
			args: args{
				productId: "453b4f0f-1f56-4c57-b43d-7b79792450a7",
			},
			want: domain.ProductsList{
				Id: "453b4f0f-1f56-4c57-b43d-7b79792450a7", Title: "Твидовый кардиган из хлопка", Image: "w1.webp", Price: 749000, Sale: 0, SaleOldPrice: 0, Category: "Женщинам", Type: "Одежда", Subtype: "Старые-коллекции", Description: "", CreatedAt: time.Date(2022, 01, 12, 13, 8, 21, 32963, time.Local),
			},
		},

		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "image", "price", "sale", "sale_old_price", "category", "type", "subtype", "description", "created_at"})

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM products WHERE id = $1")).WithArgs("453b4f0f-1f56-4c57-b43d-7b79792450a7").WillReturnRows(rows)
			},
			args: args{
				productId: "453b4f0f-1f56-4c57-b43d-7b79792450a7",
			},
			want: domain.ProductsList{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.GetById(testCase.args.productId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductsListPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	r := NewProductsListPostgres(db)

	type args struct {
		productId string
	}

	testTable := []struct {
		name    string
		mock    func()
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = $1")).
					WithArgs("453b4f0f-1f56-4c57-b43d-7b79792450a7").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				productId: "453b4f0f-1f56-4c57-b43d-7b79792450a7",
			},
		},

		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM products WHERE id = $1")).
					WithArgs("453b4f0f-1f56-4c57-b43d-7b79792450a7").WillReturnError(sql.ErrNoRows)
			},
			args: args{
				productId: "453b4f0f-1f56-4c57-b43d-7b79792450a7",
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			err := r.Delete(testCase.args.productId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProductsListPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	r := NewProductsListPostgres(db)

	type args struct {
		productId string
		item      domain.UpdateProductInput
	}

	testTable := []struct {
		name    string
		mock    func()
		args    args
		wantErr bool
	}{
		{
			name: "OK_AllFields",
			mock: func() {
				mock.ExpectExec("UPDATE products SET (.+) WHERE (.+)").
					WithArgs("new title", "new image", 1000, 1000, 100, "new category", "new type", "new subtype", "new description", "453b4f0f-1f56-4c57-b43d-7b79792450a7").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				productId: "453b4f0f-1f56-4c57-b43d-7b79792450a7",
				item: domain.UpdateProductInput{
					Title:        stringPointer("new title"),
					Image:        stringPointer("new image"),
					Price:        uintPointer(1000),
					Sale:         uintPointer(1000),
					SaleOldPrice: uintPointer(100),
					Category:     stringPointer("new category"),
					Type:         stringPointer("new type"),
					Subtype:      stringPointer("new subtype"),
					Description:  stringPointer("new description"),
				},
			},
		},
		{
			name: "OK_NoInputFields",
			mock: func() {
				mock.ExpectExec("UPDATE products SET WHERE (.+)").
					WithArgs("453b4f0f-1f56-4c57-b43d-7b79792450a7").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				productId: "453b4f0f-1f56-4c57-b43d-7b79792450a7",
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			err := r.Update(testCase.args.productId, testCase.args.item)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func uintPointer(b uint) *uint {
	return &b
}
