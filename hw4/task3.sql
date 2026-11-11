SELECT u.name, COALESCE(SUM(o.price), 0) AS total_cost
FROM users AS u
LEFT JOIN orders AS o ON u.id = o.user_id
GROUP BY u.id, u.name
ORDER BY total_cost DESC
