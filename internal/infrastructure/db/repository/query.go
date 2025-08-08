package repository

const (
	queryAddingNewCoin = `INSERT INTO public.watched_currencies (symbol) VALUES ($1) ON CONFLICT DO NOTHING;`
	queryDeleteCoin    = `DELETE FROM public.watched_currencies WHERE symbol = $1;`

	queryListRequest     = `SELECT symbol FROM public.watched_currencies;`
	queryUpdatePriceCoin = `INSERT INTO public.currency_prices (symbol, "timestamp", price, "precision", currency) VALUES ($1, $2, $3, $4, $5);`

	queryGetPriceForCoin = `(
  		SELECT symbol, price, "precision", currency, "timestamp"
  		FROM public.currency_prices
  		WHERE symbol = $1 AND "timestamp" = $2
  		LIMIT 1
	)
	UNION ALL
	(
  		SELECT symbol, price, "precision", currency, "timestamp"
		FROM public.currency_prices
		WHERE symbol = $1 AND "timestamp" <> $2
		ORDER BY ABS("timestamp" - $2)
		LIMIT 1
	)
	LIMIT 1;`
)
