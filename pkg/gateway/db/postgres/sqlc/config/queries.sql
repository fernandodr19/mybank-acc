-- name: CreateAccount :one
INSERT INTO accounts (document, balance, available_credit)
VALUES (@document, @balance, @available_credit)
RETURNING id;

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id = @id;

-- name: Deposit :execrows
UPDATE accounts
SET balance = balance + @amount
WHERE id = @id;

-- name: Withdraw :execrows
UPDATE accounts
SET balance = balance - @amount
WHERE id = @id AND (balance >= @amount);

-- name: DecreaseAvailableCredit :execrows
UPDATE accounts
SET available_credit = available_credit - @amount
WHERE id = @id;