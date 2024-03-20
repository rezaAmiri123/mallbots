package agent

// import (
// 	"context"


// 	"github.com/jackc/pgx/v4/pgxpool"
// 	"github.com/rezaAmiri123/edatV2/di"
// 	edatpgx "github.com/rezaAmiri123/edatV2/pgx"
	
// )

// func (a *Agent) setupDatabase() (err error) {
// 	var pgConn *pgxpool.Pool
// 	pgConn, err = pgxpool.Connect(context.Background(), a.config.Postgres.Conn)
// 	if err!= nil{
// 		return err
// 	}

// 	// 1. Outbox: Use session client which will fetch a transaction from the context
// 	pgTxConn := edatpgx.NewSessionClient()



// 	a.container.AddSingleton(constants.DatabaseKey, func(c di.Container) (any, error) {
// 		return pgConn, nil
// 	})

// 	a.container.AddSingleton(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
// 		return pgTxConn, nil
// 	})

// 	return nil
// }

// // //
// // var data []byte
// // query :=fmt.Sprintf("SELECT data FROM %s WHERE id = 'a3856544-b8cb-4adb-8fc6-7aafe7cbec18' LIMIT 1", constants.OutboxTableName)
// // dbConn.QueryRowContext(context.Background(), query).Scan(&data)
// // // fmt.Println("xxxxxxxxx ",string(data))
// // var Addresses consumersV1.ConsumerRegistered
// // // if err = proto.Unmarshal(data, &Addresses);err!= nil{
// // // 	fmt.Println("yyyyyyyyyyy ",err.Error())
// // // }
// // if err = json.Unmarshal(data, &Addresses); err != nil {
// // 	fmt.Println("yyyyyyyyyyy ",err.Error())
// // }
// // fmt.Println("zzzzzzzzzzzz ",Addresses.Id)
// // fmt.Println("zzzzzzzzzzzz2 ",Addresses.Name)

// // //
