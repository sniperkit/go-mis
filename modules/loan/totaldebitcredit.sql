select sum(account_transaction_debit.amount)
from account_transaction_debit
join r_account_transaction_debit on r_account_transaction_debit."accountTransactionDebitId" = account_transaction_debit.id
where r_account_transaction_debit."accountId" = 2470

select sum(account_transaction_credit.amount)
from account_transaction_credit
join r_account_transaction_credit on r_account_transaction_credit."accountTransactionCreditId" = account_transaction_credit.id
where r_account_transaction_credit."accountId" = 2470