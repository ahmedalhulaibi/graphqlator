package mysqlsubstance

/*GetCurrentDatabaseNameQuery used in GetCurrentDatabaseNamefunc*/
var GetCurrentDatabaseNameQuery = `SELECT DATABASE()`

/*DescribeDatabaseQuery used in DescribeDatabaseFunc*/
var DescribeDatabaseQuery = `SHOW TABLES`

/*DescribeTableQuery used in DescribeTableFunc*/
var DescribeTableQuery = `DESCRIBE %s`

/*DescribeTableRelationshipQuery used in DescribeTableRelationshipFunc*/
var DescribeTableRelationshipQuery = `SELECT 
TABLE_NAME,COLUMN_NAME,CONSTRAINT_NAME, REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME
FROM
INFORMATION_SCHEMA.KEY_COLUMN_USAGE
WHERE
REFERENCED_TABLE_SCHEMA = '%s' AND
REFERENCED_TABLE_NAME = ?;`

/*DescribeTableConstraintsQuery used in DescribeTableConstraintsFunc*/
var DescribeTableConstraintsQuery = `SELECT DISTINCT kcu.column_name as 'Column', tc.constraint_type as 'Constraint'
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE as kcu
JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS as tc on tc.constraint_name = kcu.constraint_name
WHERE kcu.table_name = ?
order by kcu.column_name, tc.constraint_type;`
