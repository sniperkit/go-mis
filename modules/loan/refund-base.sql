select loan.id as loan_id, investor.id as investor_id, account.id as investor_account_id, loan.plafond,
account."totalDebit", account."totalCredit", account."totalBalance"
from loan
join r_investor_product_pricing_loan on r_investor_product_pricing_loan."loanId" = loan.id
join investor on investor.id = r_investor_product_pricing_loan."investorId"
join r_account_investor on r_account_investor."investorId" = investor.id
join account on account.id = r_account_investor."accountId"
where loan.id = 31553