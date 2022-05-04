package repository

const (
	sqlGetOrder = `
	SELECT t.orderid, t.orderdate, t.netamount, t.tax, t.totalamount,
	t.prod_id, p.title, p.price, t.quantity
	FROM products p INNER JOIN
		(SELECT o.*, ol.prod_id, ol.quantity
		FROM orders o INNER JOIN orderlines ol
		ON o.orderid = ol.orderid) t
	ON p.prod_id = t.prod_id
	WHERE orderid=$1
	`
	sqlGetCustomerOrders = `
	SELECT t.orderid, t.orderdate, t.netamount, t.tax, t.totalamount,
	t.prod_id, p.title, p.price, t.quantity
	FROM products p INNER JOIN
		(SELECT o.*, ol.prod_id, ol.quantity
		FROM orders o INNER JOIN orderlines ol
		ON o.orderid = ol.orderid) t
	ON p.prod_id = t.prod_id
	WHERE customerid=$1
	`
	sqlAddOrderSelectProducts = `
	SELECT i.prod_id, i.quan_in_stock, p.price, p.title
	FROM inventory i INNER JOIN products p
	ON i.prod_id = p.prod_id
	WHERE i.prod_id = ANY($1)
	`
	sqlAddOrder = `
	INSERT INTO orders (orderdate, customerid, netamount, tax, totalamount) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING orderid
	`
	sqlAddOrderOrderlines = `
	INSERT INTO orderlines (orderlineid, orderid, prod_id, quantity, orderdate) 
	VALUES ($1, $2, $3, $4, $5)
	`

	sqlGetAllProducts = `
	SELECT p.prod_id, p.title, p.price, i.quan_in_stock 
	FROM products p INNER JOIN inventory i
	ON p.prod_id = i.prod_id
	LIMIT $1
	`
	sqlGetProduct = `
	SELECT p.prod_id, p.title, p.price, i.quan_in_stock 
	FROM products p INNER JOIN inventory i
	ON p.prod_id = i.prod_id
	WHERE p.prod_id = $1
	`
	// I use only 2 columns from sample database to simplify the project logic
	sqlAddProduct = `
	INSERT INTO products (category, title, actor, price, special, common_prod_id)
	VALUES (-1, $1, '', $2, -1, -1)
	RETURNING prod_id
	`

	// I use only 3 columns from sample database to simplify the project logic
	sqlAddCustomer = `
	INSERT INTO customers (
		firstname,
		lastname,
		address1,
		address2,
		city,
		state,
		zip,
		country,
		region,
		email,
		phone,
		creditcardtype,
		creditcard,
		creditcardexpiration,
		username,
		password,
		age,
		income,
		gender
	)
	VALUES (
		$1,
		$2,
		'',
		'',
		'',
		'',
		-1,
		'',
		-1,
		'',
		'',
		-1,
		'',
		'',
		'',
		'',
		$3,
		-1,
		''
	)
	RETURNING customerid
	`
)
