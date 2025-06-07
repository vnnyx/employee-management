package constants

const (
	ErrWrapSqlxNamed             = "Sqlx.Named()"
	ErrWrapSqlxIn                = "Sqlx.In()"
	ErrWrapPgxscanGet            = "Pgxscan.Get()"
	ErrWrapPgxscanSelect         = "Pgxscan.Select()"
	ErrWrapPgxscanScanOne        = "Pgxscan.ScanOne()"
	ErrWrapPgxscanScanAll        = "Pgxscan.ScanAll()"
	ErrWrapPgxBatchQueue         = "Pgx.Batch.Queue()"
	ErrWrapPgxBatchQueueQueryRow = "Pgx.Batch.Queue().QueryRow()"
	ErrWrapDbIn                  = "DB.In()"
	ErrWrapDbExec                = "DB.Exec()"
	ErrWrapDbQuery               = "DB.Query()"
	ErrWrapDbQueryRowScan        = "DB.QueryRow().Scan()"
	ErrWrapDbScan                = "DB.Scan()"
	ErrWrapDbCopyFrom            = "DB.CopyFrom()"
	ErrWrapDbSendBatchClose      = "DB.SendBatch().Close()"
	ErrWrapUrlQueryUnescape      = "Url.QueryUnescape()"
	ErrWrapPgxQueryRow           = "Pgx.QueryRow()"
	ErrWrapPgxQuery              = "Pgx.Query()"
	ErrWrapJsonUnmarshal         = "Json.Unmarshal()"
)

// Issue code
const (
	AuthTokenExpired   = "AUTH_TOKEN_EXPIRED"
	AuthTokenInvalid   = "AUTH_TOKEN_INVALID"
	AuthTokenMalformed = "AUTH_TOKEN_MALFORMED"
	AuthTokenNotFound  = "AUTH_TOKEN_NOT_FOUND"
	ValidationError    = "VALIDATION_ERROR"

	InternalServerError = "INTERNAL_SERVER_ERROR"
)
