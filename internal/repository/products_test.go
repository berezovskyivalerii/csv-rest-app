package repository

import (
	"berezovskyivalerii/csv-rest-app/internal/domain"
	"context"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestProducts_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewProducts(db)

	type args struct {
		ctx     context.Context
		product domain.Product
	}
	type mockBehavior func(args args, id int64)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int64
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx: ctx,
				product: domain.Product{
					Name:  "Shaker",
					Price: 100,
				},
			},
			id: 2,
			mockBehavior: func(args args, id int64) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("insert into products").
					WithArgs(args.product.Name, args.product.Price).WillReturnRows(rows)
			},
		},
		{
			name: "SEcond",
			args: args{
				ctx: ctx,
				product: domain.Product{
					Name: "Rofl",
				},
			},
			mockBehavior: func(args args, id int64) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("insert into products").
					WithArgs(args.product.Name, args.product.Price).WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.Create(testCase.args.ctx, testCase.args.product)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}

		})
	}
}
