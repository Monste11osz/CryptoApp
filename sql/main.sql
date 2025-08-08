-- Таблица с ценами криптовалют
CREATE TABLE public.currency_prices (
                                        id serial4 NOT NULL,
                                        symbol varchar(50) NOT NULL,
                                        "timestamp" int8 NOT NULL,
                                        price int8 NOT NULL,
                                        "precision" int2 DEFAULT 8 NOT NULL,
                                        currency text DEFAULT 'USD'::text NULL,
                                        CONSTRAINT currency_prices_pkey PRIMARY KEY (id),
                                        CONSTRAINT currency_prices_symbol_timestamp_key UNIQUE (symbol, "timestamp")
);

CREATE INDEX idx_currency_prices_symbol ON public.currency_prices USING btree (symbol);
CREATE INDEX idx_timestamp ON public.currency_prices USING btree ("timestamp");

-- Таблица отслеживаемых валют
CREATE TABLE public.watched_currencies (
                                           id serial4 NOT NULL,
                                           symbol text NOT NULL,
                                           added_at timestamp DEFAULT now() NULL,
                                           CONSTRAINT watched_currencies_pkey PRIMARY KEY (id),
                                           CONSTRAINT watched_currencies_symbol_key UNIQUE (symbol)
);
