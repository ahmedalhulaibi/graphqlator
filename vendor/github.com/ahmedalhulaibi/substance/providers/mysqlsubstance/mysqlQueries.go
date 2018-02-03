package mysqlsubstance

/*DescribeDatabaseQuery used in DescribeDatabaseFunc*/
var DescribeDatabaseQuery string = `SHOW TABLES`

/*DescribeTableQuery used in DescribeTableFunc*/
var DescribeTableQuery string = `DESCRIBE %s`

/*DescribeTableRelationshipQuery used in DescribeTableRelationshipFunc*/
var DescribeTableRelationshipQuery string = `SELECT 
TABLE_NAME,COLUMN_NAME,CONSTRAINT_NAME, REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME
FROM
INFORMATION_SCHEMA.KEY_COLUMN_USAGE
WHERE
REFERENCED_TABLE_SCHEMA = '%s' AND
REFERENCED_TABLE_NAME = '%s';`

/*DescribeTableConstraintsQuery used in DescribeTableConstraintsFunc*/
var DescribeTableConstraintsQuery string = `SELECT DISTINCT kcu.column_name as 'Column', tc.constraint_type as 'Constraint'
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE as kcu
JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS as tc on tc.constraint_name = kcu.constraint_name
WHERE kcu.table_name = '%s'
order by kcu.column_name, tc.constraint_type;`
