package pgsqlsubstance

/*DescribeDatabaseQuery used in DescribeDatabaseFunc*/
var DescribeDatabaseQuery string = `select * from pg_catalog.pg_tables where schemaname not in ('pg_catalog','information_schema');`

/*DescribeTableQuery used in DescribeTableFunc*/
var DescribeTableQuery string = `select
att.attrelid as "classId",
class.relname as "Table",
att.attname as "Field",
dsc.description as "description",
typ.typname as "Type",
att.attnum as "num",
att.attnotnull as "isNotNull",
att.atthasdef as "hasDefault"
from
pg_catalog.pg_attribute as att
left join pg_catalog.pg_description as dsc on dsc.objoid = att.attrelid
and dsc.objsubid = att.attnum
left join pg_type as typ on typ.oid = att.atttypid
left join pg_catalog.pg_class as class on class.oid = att.attrelid
where
att.attrelid in (
	select
		rel.oid as "id"
	from
		pg_catalog.pg_class as rel
		left join pg_catalog.pg_description as dsc on dsc.objoid = rel.oid
		and dsc.objsubid = 0
	where
		class.relname = $1
		and rel.relpersistence in ('p')
		and rel.relkind in ('r', 'v', 'm', 'c', 'f')
)
and att.attnum > 0
and not att.attisdropped
order by
att.attrelid,
att.attnum;`

/*DescribeTableRelationshipQuery used in DescribeTableRelationshipFunc*/
var DescribeTableRelationshipQuery string = `select 
tc.table_name as "table_name",
kcu.column_name as "column",
class.relname as "ref_table",
con.confkey as "ref_columnNum"
  from
	pg_catalog.pg_constraint as con
	left join information_schema.table_constraints as tc on tc.constraint_name = con.conname
	left join information_schema.key_column_usage as kcu on kcu.constraint_name = con.conname
	left join pg_catalog.pg_class as class on class.oid = con.confrelid
  where
	tc.table_name = $1 and
	con.contype = 'f'
	  ;`

/*DescribeTableConstraintsQuery used in DescribeTableConstraintsFunc*/
var DescribeTableConstraintsQuery string = `select distinct on (con.conrelid, con.conkey, con.confrelid, con.confkey)
	tc.table_name,
	kcu.column_name as "column",
	contype
  from
	pg_catalog.pg_constraint as con
	left join information_schema.table_constraints as tc on tc.constraint_name = con.conname
	left join information_schema.key_column_usage as kcu on kcu.constraint_name = con.conname
	left join pg_catalog.pg_class as class on class.oid = con.confrelid
  where
		tc.table_name = $1
  ;`
