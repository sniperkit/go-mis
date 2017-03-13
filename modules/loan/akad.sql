select loan."agreementType", loan."updatedAt", loan.purpose, loan.plafond, loan.tenor, loan.installment, cif_investor.name as investor, cif_borrower.name as borrower, "group"."name" as "group"
from loan
join r_investor_product_pricing_loan on r_investor_product_pricing_loan."loanId" = loan.id
join investor on investor.id = r_investor_product_pricing_loan."investorId"
join r_cif_investor on r_cif_investor."investorId" = investor.id
join (
	select * from cif where cif."deletedAt" is null
) as cif_investor on cif_investor.id = r_cif_investor."cifId"
join r_loan_borrower on r_loan_borrower."loanId" = loan.id
join borrower on borrower.id  = r_loan_borrower."borrowerId"
join r_cif_borrower on r_cif_borrower."borrowerId" = borrower.id
join (
	select * from cif where cif."deletedAt" is null
) as cif_borrower on cif_borrower.id = r_cif_borrower."cifId"
join r_loan_group on r_loan_group."loanId" = loan.id
join "group" on "group".id = r_loan_group."groupId"
where loan."deletedAt" is null and loan.id = 21417