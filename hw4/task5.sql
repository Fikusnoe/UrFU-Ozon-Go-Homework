SELECT o.item, o.price
FROM orders AS o
WHERE o.price = (
    SELECT MIN(price)
    FROM orders
    WHERE price > 10000
)